package models

import "time"

type OnePassword struct {
	CloudFlareToken string `mapstructure:"cloudflare"`
	ProxmoxToken    string `mapstructure:"proxmox"`
	NatsToken       string `mapstructure:"nats"`
}

type Config struct {
	Proxmox     Proxmox     `mapstructure:"proxmox"`
	Nats        Nats        `mapstructure:"nats"`
	Api         Api         `mapstructure:"api"`
	Sync        Sync        `mapstructure:"sync"`
	SSH         SSHConfig   `mapstructure:"ssh"`
	PushGateway PushGateway `mapstructure:"push_gateway"`
	CloudFlare  CloudFlare  `mapstructure:"cloudflare"`
	GCP         GCP         `mapstructure:"gcp"`
	OnePassword OnePassword `mapstructure:"1password"`
}

type SSHConfig struct {
	User           string `mapstructure:"user"`
	PrivateKey     string `mapstructure:"private_key"`
	PublicKey      string `mapstructure:"public_key"`
	PrivateKeyFile string `mapstructure:"private_key_file"`
	PublicKeyFile  string `mapstructure:"public_key_file"`
	ConfigFile     string `mapstructure:"config_file"`
}

type Sync struct {
	Interval   time.Duration `mapstructure:"interval"`
	Timeout    time.Duration `mapstructure:"timeout"`
	VMS        bool          `mapstructure:"vms"`
	Containers bool          `mapstructure:"containers"`
	Enabled    bool          `mapstructure:"enabled"`
}

type Api struct {
	Port      int    `mapstructure:"port"`
	Host      string `mapstructure:"host"`
	PublicUrl string `mapstructure:"public_url"`
	Scheme    string `mapstructure:"scheme"`
	Version   string `mapstructure:"version"`
}

type Nats struct {
	URL      string `mapstructure:"url"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Bucket   string `mapstructure:"bucket"`
	Stream   string `mapstructure:"stream"`
	Replicas int    `mapstructure:"replicas"`
}

type Proxmox struct {
	URL  string `mapstructure:"url"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
}

type PushGateway struct {
	URL          string        `mapstructure:"url"`
	PushInterval time.Duration `mapstructure:"push_interval"`
}

type CloudFlare struct {
	ApiToken  string `mapstructure:"api_token"`
	IsDefault bool   `mapstructure:"is_default"`
}

type GCP struct {
	ProjectId       string `mapstructure:"project_id"`
	CredentialsFile string `mapstructure:"credentials_file"`
	IsDefault       bool   `mapstructure:"is_default"`
}
