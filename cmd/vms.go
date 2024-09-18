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
	"fmt"
	"i2/pkg/prxmx"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	all      bool
	asTable  bool
	cluster  *prxmx.Cluster
	selected string
)

// vmsCmd represents the vms command
var vmsCmd = &cobra.Command{
	Use:   "vms",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		proxmoxURL := viper.GetString("proxmox.url")
		proxmoxUser := viper.GetString("proxmox.user")
		proxmoxPass := viper.GetString("proxmox.pass")

		cluster = prxmx.NewCluster(proxmoxURL, proxmoxUser, proxmoxPass)

		if asTable {
			getVMs()
		} else {
			getVMsList()
			// get the IP of the selected VM
			tokens := strings.Split(selected, "192")
			ip := "ssh ivan@192" + tokens[1]
			fmt.Println(ip)
		}
	},
}

func getVMsList() {
	vms, err := cluster.GetVMs(all)
	if err != nil {
		log.Fatalf("Error getting VMs: %v", err)
	}
	items := []list.Item{}
	running := 0
	for _, vm := range vms {
		if vm.Running {
			running++
			items = append(items, item{title: "🖥️    " + vm.Name + "   " + getLocalIP(vm.IP), description: "     " + vm.Uptime.ToStringShort()})
		}
	}
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Proxmox VMs"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
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

func init() {
	rootCmd.AddCommand(vmsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// vmsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	vmsCmd.Flags().BoolVarP(&all, "all", "a", false, "Return all VMs")
	vmsCmd.Flags().BoolVarP(&asTable, "table", "t", false, "Return a table")
}

func getVMs() {

	vms, err := cluster.GetVMs(all)
	if err != nil {
		log.Fatalf("Error getting VMs: %v", err)
	}

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
			// case row%2 == 0:
			// 	return evenRowStyle
			default:
				return oddRowStyle
			}
		}).
		Headers("Name", "IP", "Uptime")

	for _, vm := range vms {
		if vm.Running {
			running++
			t.Row("🖥️  "+vm.Name, " "+getLocalIP(vm.IP), vm.Uptime.ToStringShort())
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

func getLocalIP(ips []string) string {
	for _, ip := range ips {
		if strings.HasPrefix(ip, "192.168") {
			return ip
		}
	}
	return ""
}
