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
	"fmt"
	"i2/pkg/dckr"
	"i2/pkg/dfiles"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	filePath       string
	containerID    string
	privateKeyPath string
	remoteUser     string
	remoteDir      string
	scpHost        string
)

// cpCmd represents the cp command
var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if scpHost != "" {
			copyFileRemote()
		} else {
			copyFileToContainer(filePath, containerID)
		}
	},
}

func init() {
	rootCmd.AddCommand(cpCmd)

	cpCmd.Flags().StringVarP(&filePath, "file", "f", "", "File to copy")
	cpCmd.Flags().StringVarP(&containerID, "container", "c", "", "Container ID")
	cpCmd.Flags().StringVarP(&privateKeyPath, "private-key", "p", "", "Private key path")
	cpCmd.Flags().StringVarP(&scpHost, "scp", "n", "", "Scp remote host with the file: remote_username@HOST:/remote/directory")
}

func parseScpHost() (string, string, string) {
	parts := strings.Split(scpHost, "@")
	if len(parts) != 2 {
		log.Fatalf("Invalid scp host format. Expected format: username@host")
	}
	username := parts[0]
	host := parts[1]
	return username, host, remoteDir
}

func copyFileRemote() error {
	pkp := viper.GetString("ssh.privateKeyPath")
	if scpHost != "" {
		remoteUser, sshHost, remoteDir = parseScpHost()
	}
	return dfiles.SCPTransfer(filePath, remoteUser, sshHost, remoteDir, pkp)
}

func copyFileToContainer(filePath, containerID string) error {
	// Create a new Docker client
	dc, err := dckr.NewDockerClient()
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer dc.Close()

	// Open the local file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer file.Close()

	// Get file info
	_, err = file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Copy the file to the container
	err = dc.CopyToContainer(context.Background(), containerID, "/", file, container.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
		CopyUIDGID:                true,
	})
	if err != nil {
		return fmt.Errorf("failed to copy file to container: %w", err)
	}

	fmt.Printf("File %s copied successfully to container %s\n", filePath, containerID)
	return nil

}
