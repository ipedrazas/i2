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

	"i2/pkg/dckr"
	"i2/pkg/prxmx"
	"i2/pkg/store"
	"i2/pkg/utils"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	all              bool
	bucketVMS        string
	bucketContainers string
	user             string
	print            bool
)

// csCmd represents the cs command
var csCmd = &cobra.Command{
	Use:     "containers",
	Aliases: []string{"cs"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		bucket := viper.GetString("nats.bucket")
		bucketVMS = bucket + "-vms"
		bucketContainers = bucket + "-containers"
		user = viper.GetString("ssh.user")

		// get host data from NATS
		ctx := context.Background()
		conf := getDefaultNatsConf()
		if conf.Url == "" {
			log.Fatalf("NATS_URL is not set")
		}
		st, err := store.NewStore(ctx, &conf)
		if err != nil {
			log.Fatalf("Error creating store: %v", err)
		}
		defer st.Close()

		if all {
			csInVms := doAll(st, ctx)
			totalContainers := 0

			for k, containers := range csInVms {
				totalContainers += len(containers)
				if print {
					nameIp := strings.Split(k, "-")
					printContainers(nameIp[0], nameIp[1], containers)
				}
			}
			log.Info("Total VMs: %d Total Containers: %d\n", len(csInVms), totalContainers)
			return
		}

		// Local containers
		if len(args) == 0 {
			containers := listLocalContainers()
			printContainers("Localhost", "127.0.0.1", containers)
			return
		}
		// remote containers
		if len(args[0]) > 0 {
			if !strings.HasPrefix(args[0], "ssh://") {
				vm, err := getVMFromNATS(args[0], st, ctx)
				if err != nil {
					log.Fatalf("Error getting VM: %v", err)
				}
				sshHost := "ssh://" + user + "@" + utils.GetLocalIP(vm.IP)
				containers := listRemoteContainers(st, sshHost, ctx)
				printContainers(vm.Name, utils.GetLocalIP(vm.IP), containers)
				return
			}
			containers := listRemoteContainers(st, args[0], ctx)
			iphost := strings.Split(args[0], "@")
			printContainers("Remote", iphost[1], containers)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(csCmd)

	csCmd.Flags().BoolVarP(&all, "all", "a", false, "return containers in all nodes")
	csCmd.Flags().BoolVarP(&print, "print", "p", false, "print containers to stdout")

}

func listLocalContainers() []types.Container {
	dc, err := dckr.NewDockerClient()
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}
	defer dc.Close()
	containers, _ := dc.ListContainers()
	return containers
}

func getVMFromNATS(hostname string, st *store.Store, ctx context.Context) (prxmx.Node, error) {
	jvm, err := store.GetKV(ctx, hostname, bucketVMS, st.NatsConn)
	if err != nil {
		log.Fatalf("Error getting VM From NATS: %v %s %s", err, hostname, bucketVMS)
	}
	vm := prxmx.Node{}
	err = json.Unmarshal(jvm, &vm)
	return vm, err
}

func listRemoteContainers(st *store.Store, sshHost string, ctx context.Context) []types.Container {
	var dc *dckr.DockerClient
	var err error

	if sshHost != "" {
		if !strings.HasPrefix(sshHost, "ssh://") {
			vm, err := getVMFromNATS(sshHost, st, ctx)
			if err != nil {
				log.Fatalf("Error getting VM: %v", err)
				return []types.Container{}
			}
			user := viper.GetString("ssh.user")
			sshHost = "ssh://" + user + "@" + utils.GetLocalIP(vm.IP)
		}
		dc, err = dckr.NewDockerClientWithSSH(sshHost)
		if err != nil {
			log.Info("Error creating Docker client:", err, sshHost)
			return []types.Container{}
		}
		defer dc.Close()
		containers, err := dc.ListContainers()
		if err != nil {
			log.Info("Error listing containers:", err)
			return []types.Container{}
		}
		return containers
	}
	return []types.Container{}
}

func doAll(st *store.Store, ctx context.Context) map[string][]types.Container {
	vms := []prxmx.Node{}
	allContainers := make(map[string][]types.Container)
	keys, _ := store.GetKeys(ctx, bucketVMS, st.NatsConn)
	for _, key := range keys {
		vm, err := getVMFromNATS(key, st, ctx)
		if err != nil {
			log.Fatalf("Error getting VM: %v", err)
			continue
		}
		vms = append(vms, vm)
	}
	for _, vm := range vms {
		if vm.Running {
			log.Info("Processing VM", vm.Name, "IP", vm.IP)
			user := viper.GetString("ssh.user")
			sshHost := "ssh://" + user + "@" + utils.GetLocalIP(vm.IP)
			containers := listRemoteContainers(st, sshHost, ctx)
			if len(containers) > 0 {
				saveContainers(ctx, vm, containers, st)
			}
			lip := vm.Name + "-" + utils.GetLocalIP(vm.IP)
			allContainers[lip] = containers
		}
	}
	return allContainers
}

func saveContainers(ctx context.Context, vm prxmx.Node, containers []types.Container, st *store.Store) {
	bcontainers, err := json.Marshal(containers)
	if err != nil {
		log.Fatalf("Error marshalling container: %v", err)
	}
	// Print container ID and name

	err = store.SetKV(ctx, vm.Name, bucketContainers, bcontainers, st.NatsConn)
	if err != nil {
		log.Fatalf("Error storing container: %v", err)
	}
}

func printContainers(hostname, ip string, containers []types.Container) {

	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Foreground(lipgloss.Color("252")).Bold(true)
	rowStyle := baseStyle.Foreground(lipgloss.Color("250"))

	var styleTop = lipgloss.NewStyle().
		Bold(true).
		MarginTop(1).
		PaddingTop(1).
		PaddingLeft(2).
		Foreground(lipgloss.Color("#c6a0f1")).
		Background(lipgloss.Color("#190e27")).
		Width(120)
	stop := styleTop.Render("VM: " + hostname + "\t\tIP: " + ip)

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return headerStyle
			default:
				return rowStyle
			}
		}).
		Headers("ID", "Name", "Image", "Ports", "Days Old").Width(120)

	for _, container := range containers {
		daysOld := int(time.Since(time.Unix(container.Created, 0)).Hours() / 24)
		t.Row(container.ID[:12], container.Names[0][1:], container.Image, dckr.PortsAsString(container.Ports), strconv.Itoa(daysOld))
	}

	group := lipgloss.JoinVertical(
		lipgloss.Left,
		stop,
		t.Render(),
	)

	fmt.Println(group)

}
