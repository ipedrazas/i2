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
	"i2/pkg/api"
	"i2/pkg/models"
	"i2/pkg/utils"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := readConfig()
		api.RunServer(conf)
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

func readConfig() *models.Config {
	ctx := context.Background()
	conf := &models.Config{}
	err := viper.Unmarshal(conf)
	if err != nil {
		log.Errorf("unable to decode into config struct, %v", err)
	}

	key := conf.OnePassword.CloudFlareToken
	cfToken, err := utils.ReadSecretFrom1Password(ctx, key)
	if err != nil {
		log.Errorf("unable to read Cloudflare token, %v", err)
	}
	conf.CloudFlare.ApiToken = cfToken

	key = conf.OnePassword.ProxmoxToken
	proxmoxToken, err := utils.ReadSecretFrom1Password(ctx, key)
	if err != nil {
		log.Errorf("unable to read proxmox token, %v", err)
		panic(err)
	}
	conf.Proxmox.Pass = proxmoxToken

	key = conf.OnePassword.NatsToken
	natsToken, err := utils.ReadSecretFrom1Password(ctx, key)
	if err != nil {
		log.Errorf("unable to read NATS token, %v", err)
		panic(err)
	}
	conf.Nats.Password = natsToken
	log.Infof("config %v", conf)
	return conf
}
