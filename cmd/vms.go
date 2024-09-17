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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var all bool

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
		getVMs()
	},
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
}

func getVMs() {
	proxmoxURL := viper.GetString("proxmox.url")
	proxmoxUser := viper.GetString("proxmox.user")
	proxmoxPass := viper.GetString("proxmox.pass")

	cluster := prxmx.NewCluster(proxmoxURL, proxmoxUser, proxmoxPass)

	vms, err := cluster.GetVMs(all)
	if err != nil {
		log.Fatalf("Error getting VMs: %v", err)
	}

	running := 0
	for _, vm := range vms {
		fmt.Println(vm.ToString())
		if vm.Running {
			running++
		}
	}
	fmt.Println("\nTotal VMs: ", len(vms))
	fmt.Println("Running VMs: ", running)
}
