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
	"i2/pkg/api"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		version()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func version() {

	var styleLeft0 = lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("#3c017d")).
		PaddingTop(1).
		Width(2)
	var styleLeft1 = lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("#280154")).
		PaddingTop(1).
		Width(4)
	var styleRight = lipgloss.NewStyle().
		Bold(true).
		PaddingTop(1).
		PaddingLeft(2).
		Foreground(lipgloss.Color("#c6a0f1")).
		Background(lipgloss.Color("#190e27")).
		Width(48)

	vstr := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		styleLeft0.Render("\n\n\n\n\n "),
		styleLeft1.Render("\n\n\n\n\n "),
		styleRight.Render("i2 - Ivan's Homelab Platform cli\n\n\tVersion: "+api.Version+"\n\tBuild Date: "+api.BuildDate+"\n\tGit Commit: "+api.GitCommit+"\n"),
	)

	fmt.Println(vstr)

}
