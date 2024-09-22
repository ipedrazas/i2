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
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	sshHost          string
	bucketVMS        string
	bucketContainers string
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
		listContainers()
	},
}

func init() {
	rootCmd.AddCommand(csCmd)
	csCmd.Flags().StringVarP(&sshHost, "ssh", "s", "", "SSH connection string")
	csCmd.Flags().BoolVarP(&all, "all", "a", false, "return all containers in all nodes")
	bucket := viper.GetString("nats.bucket")
	bucketVMS = bucket + "-vms"
	bucketContainers = bucket + "-containers"
}

func listContainers() {
	var dc *dckr.DockerClient
	var err error
	if all {
		doAll()
		return
	}

	if sshHost != "" {
		if !strings.HasPrefix(sshHost, "ssh://") {
			// get host data from NATS
			ctx := context.Background()
			conf := getDefaultNatsConf()
			st, err := store.NewStore(ctx, &conf)
			if err != nil {
				log.Fatalf("Error creating store: %v", err)
			}
			defer st.Close()

			jvm, err := store.GetKV(ctx, sshHost, bucketVMS, st.NatsConn)
			if err != nil {
				log.Fatalf("Error getting VM: %v", err)
			}
			vm := prxmx.Node{}
			err = json.Unmarshal(jvm, &vm)
			if err != nil {
				log.Fatalf("Error unmarshalling VM: %v", err)
			}
			sshHost = "ssh://ivan@" + utils.GetLocalIP(vm.IP)

		}
		dc, err = dckr.NewDockerClientWithSSH(sshHost)
		if err != nil {
			fmt.Println("Error creating Docker client:", err)
			return
		}
	} else {
		dc, err = dckr.NewDockerClient()
		if err != nil {
			fmt.Println("Error creating Docker client:", err)
			return
		}
	}

	defer dc.Close()
	containers, err := dc.ListContainers()
	if err != nil {
		fmt.Println("Error listing containers:", err)
		return
	}
	fmt.Println("Container ID \tName \t\tImage")
	for _, container := range containers {
		// Print container ID and name
		fmt.Printf("%s \t%s \t%s\n", container.ID[:12], container.Names[0][1:], container.Image)
	}
}

func doAll() {
	vms := []prxmx.Node{}
	ctx := context.Background()
	conf := getDefaultNatsConf()
	st, err := store.NewStore(ctx, &conf)
	if err != nil {
		log.Fatalf("Error creating store: %v", err)
	}
	defer st.Close()
	keys, _ := store.GetKeys(ctx, bucketVMS, st.NatsConn)
	for _, key := range keys {
		jvm, err := store.GetKV(ctx, key, bucketVMS, st.NatsConn)
		if err != nil {
			log.Fatalf("Error getting VM: %v", err)
		}
		vm := prxmx.Node{}
		err = json.Unmarshal(jvm, &vm)
		if err != nil {
			log.Fatalf("Error unmarshalling VM: %v", err)
		}
		vms = append(vms, vm)
	}
	for _, vm := range vms {
		if vm.Running {
			containers := []types.Container{}
			var err error

			// get containers from NATS
			cs, _ := store.GetKV(ctx, vm.Name, bucketContainers, st.NatsConn)

			if len(cs) > 0 {
				err = json.Unmarshal(cs, &containers)
				if err != nil {
					log.Fatalf("Error unmarshalling containers: %v", err)
				}
			} else {
				containers, err = getRemoteContainers(vm)
				if err != nil {
					log.Fatalf("Error listing containers: %v", err)
				}
			}

			printContainers(vm, containers)
			if len(cs) == 0 {
				saveContainers(ctx, vm, containers, st)
			}
		}
	}
}

func getRemoteContainers(vm prxmx.Node) ([]types.Container, error) {
	dc, err := dckr.NewDockerClientWithSSH("ssh://ivan@" + utils.GetLocalIP(vm.IP))
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}
	defer dc.Close()
	return dc.ListContainers()
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

func printContainers(vm prxmx.Node, containers []types.Container) {

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
	stop := styleTop.Render("VM: " + vm.Name + "\t\tIP: " + utils.GetLocalIP(vm.IP))

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
		Headers("ID", "Name", "Image", "Ports").Width(120)

	for _, container := range containers {
		t.Row(container.ID[:12], container.Names[0][1:], container.Image, dckr.PortsAsString(container.Ports))
	}

	group := lipgloss.JoinVertical(
		lipgloss.Left,
		stop,
		t.Render(),
	)

	fmt.Println(group)

}
