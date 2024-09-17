package prxmx

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/spf13/viper"
)

var cluster *Cluster

func setup() {

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user home directory: %v", err)
	}

	// Search config in home directory with name ".i2" (without extension).
	viper.AddConfigPath(home)
	viper.AddConfigPath(home + "/.config/i2")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".i2")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	ApiURL := viper.GetString("proxmox.url")
	User := viper.GetString("proxmox.user")
	Pass := viper.GetString("proxmox.pass")

	cluster = NewCluster(ApiURL, User, Pass)
}

func TestCluster_GetClusterNodes(t *testing.T) {

	setup()

	// cluster := NewCluster(ApiURL, User, Pass)
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{name: "test", want: []string{"beelink02", "beelink03", "nuc02", "x1", "len01"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := cluster.GetClusterNodes()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cluster.GetClusterNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("Cluster.GetClusterNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCluster_GetVMs(t *testing.T) {
	setup()
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{name: "test", want: 23, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := cluster.GetVMs(false)

			if (err != nil) != tt.wantErr {
				t.Errorf("Cluster.GetVMs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Cluster.GetVMs() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
