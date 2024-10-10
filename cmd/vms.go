/*
Copyright © 2024 Ivan Pedrazas <ipedrazas@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"i2/pkg/prxmx"
	"i2/pkg/store"
	"i2/pkg/utils"

	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.design/x/clipboard"
)

var (
	asTable  bool
	cluster  *prxmx.Cluster
	selected string
	sync     bool
	vms      []prxmx.Node
)

type errMsg error

type spinnerModel struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func initialSpinnerModel() spinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return spinnerModel{spinner: s}
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(vms) > 0 {
		m.quitting = true
		return m, tea.Quit
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m spinnerModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s Loading VMS...press q to quit\n\n", m.spinner.View())
	if m.quitting {
		return str + "\n"
	}
	return str
}

// vmsCmd represents the vms command
var vmsCmd = &cobra.Command{
	Use:   "vms",
	Short: "List all your ProxMox VMs",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		bucket := viper.GetString("nats.bucket")
		bucketVMS = bucket + "-vms"
		bucketContainers = bucket + "-containers"

		proxmoxURL := viper.GetString("proxmox.url")
		proxmoxUser := viper.GetString("proxmox.user")
		proxmoxPass := viper.GetString("proxmox.pass")

		cluster = prxmx.NewCluster(proxmoxURL, proxmoxUser, proxmoxPass)

		ctx := context.Background()
		conf := getDefaultNatsConf()
		st, err := store.NewStore(ctx, &conf)
		if err != nil {
			log.Errorf("Error creating store: %v", err)
		}
		defer st.Close()

		if sync {
			log.Info("Fetching VMs from Proxmox")
			vms, err = cluster.GetVMs()
			if err != nil {
				log.Errorf("Error getting VMs: %v", err)
			}
			log.Info("%d VMs found\n", len(vms))
			err = syncVMS(vms)
			if err != nil {
				log.Errorf("Error syncing VMs: %v", err)
			}
			return
		}
		keys, _ := store.GetKeys(ctx, bucketVMS, st.NatsConn)

		if len(keys) == 0 {
			log.Info("Fetching VMs from Proxmox")
			go func() {
				vms, err = cluster.GetVMs()
				if err != nil {
					log.Errorf("Error getting VMs: %v", err)
				}
			}()

			p := tea.NewProgram(initialSpinnerModel())
			if _, err := p.Run(); err != nil {
				log.Errorf("Error running spinner: %v", err)
			}

			if viper.GetBool("sync.enabled") || sync {
				err = syncVMS(vms)
				if err != nil {
					log.Errorf("Error syncing VMs: %v", err)
				}
			}
		} else {
			for _, key := range keys {
				jvm, err := store.GetKV(ctx, key, bucketVMS, st.NatsConn)
				if err != nil {
					log.Errorf("Error getting VM: %v", err)
				}
				vm := prxmx.Node{}
				err = json.Unmarshal(jvm, &vm)
				if err != nil {
					log.Errorf("Error unmarshalling VM: %v", err)
				}
				vms = append(vms, vm)
			}
		}

		if asTable {
			getVMs(vms)
		} else {
			getVMsList(vms)
			// get the IP of the selected VM
			if selected != "" {
				tokens := strings.Split(selected, "192")
				ip := "ssh ivan@192" + tokens[1]
				err := clipboard.Init()
				if err != nil {
					panic(err)
				}
				clipboard.Write(clipboard.FmtText, []byte(ip))
				log.Info("ssh command copied to clipboard")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(vmsCmd)

	vmsCmd.Flags().BoolVarP(&asTable, "table", "t", false, "Return a table")
	vmsCmd.Flags().BoolVarP(&sync, "sync", "s", false, "Sync VMs with NATS")
}

func getDefaultNatsConf() store.NatsConf {
	return store.NatsConf{
		Url:      viper.GetString("nats.url"),
		User:     viper.GetString("nats.user"),
		Password: viper.GetString("nats.password"),
		Replicas: viper.GetInt("nats.replicas"),
		Bucket:   viper.GetString("nats.bucket"),
		Stream:   viper.GetString("nats.stream"),
	}
}

func syncVMS(vms []prxmx.Node) error {
	log.Info("Syncing VMs with NATS", bucketVMS)
	conf := getDefaultNatsConf()
	ctx := context.Background()
	st, err := store.NewStore(ctx, &conf)
	if err != nil {
		return err
	}
	defer st.Close()
	for _, vm := range vms {
		err = store.SetKV(ctx, vm.Name, bucketVMS, vm.ToBytes(), st.NatsConn)
		if err != nil {
			return err
		}
	}
	return nil
}

func getVMsList(vms []prxmx.Node) {

	items := []list.Item{}
	running := 0
	for _, vm := range vms {
		if vm.Running {
			running++
			items = append(items, item{title: "🖥️    " + vm.Name + "   " + utils.GetLocalIP(vm.IP), description: "     " + vm.Uptime.ToStringShort()})
		}
	}
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Proxmox VMs"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Error("Error running program:", err)
		os.Exit(1)
	}
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }

func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", " ":
			selected = m.list.SelectedItem().(item).title
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func getVMs(vms []prxmx.Node) {

	running := 0

	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Foreground(lipgloss.Color("252")).Bold(true)
	// selectedStyle := baseStyle.Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
	oddRowStyle := baseStyle.Foreground(lipgloss.Color("250"))
	// evenRowStyle := baseStyle.Foreground(lipgloss.Color("250"))
	docStyle := lipgloss.NewStyle().Padding(0, 0)
	sleepingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	width := 120

	t := table.New().
		Width(width).
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return headerStyle
			default:
				return oddRowStyle
			}
		}).
		Headers("Name", "IP", "Uptime")

	for _, vm := range vms {
		if vm.Running {
			running++
			t.Row("🖥️  "+vm.Name, " "+utils.GetLocalIP(vm.IP), vm.Uptime.ToStringShort())
		} else {
			t.Row("💤  "+sleepingStyle.Render(vm.Name), " ", vm.Uptime.ToStringShort())
		}

	}
	fmt.Println(t.Render())

	doc := strings.Builder{}

	// Status Bar.

	statusNugget := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Padding(0, 1)

	statusBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle := lipgloss.NewStyle().
		Inherit(statusBarStyle).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#FF5F87")).
		Padding(0, 1).
		MarginRight(1)

	encodingStyle := statusNugget.
		Background(lipgloss.Color("#A550DF")).
		Align(lipgloss.Right)

	statusText := lipgloss.NewStyle().Inherit(statusBarStyle)

	fishCakeStyle := statusNugget.Background(lipgloss.Color("#6124DF"))

	w := lipgloss.Width

	statusKey := statusStyle.Render("ProxMox")
	total := encodingStyle.Render(fmt.Sprintf("Total VMs: %d", len(vms)))
	totalRunning := fishCakeStyle.Render(fmt.Sprintf("Running VMs: %d", running))

	statusVal := statusText.
		Width(width - w(statusKey) - w(total) - w(totalRunning)).
		Render("💻  - " + viper.GetString("proxmox.url"))

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		total,
		totalRunning,
	)
	tw := w(t.Render())

	doc.WriteString(statusBarStyle.Width(tw).Render(bar))

	fmt.Println(docStyle.Render(doc.String()))

}
