package models

import (
	"context"
	"errors"

	"os"
	"time"

	"github.com/1password/onepassword-sdk-go"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

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
	Timeout  int    `mapstructure:"timeout"`
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

func NewConfig(options ...func(*Config)) *Config {
	conf := &Config{}
	var err error
	defer func() {
		if r := recover(); r != nil {
			if errMsg, ok := r.(string); ok {
				err = errors.New(errMsg)
			}
		}
	}()
	err = viper.Unmarshal(conf)
	if err != nil {
		log.Errorf("unable to decode into config struct, %v", err)
		return nil
	}
	for _, o := range options {
		o(conf)
	}
	return conf
}

func WithOnePassword(ctx context.Context) func(*Config) {
	return func(conf *Config) {

		key := conf.OnePassword.CloudFlareToken
		cfToken, err := ReadSecretFrom1Password(ctx, key)
		if err != nil {
			log.Errorf("unable to read Cloudflare token, %v", err)
			panic("1Password token or key not found")
		}
		conf.CloudFlare.ApiToken = cfToken

		key = conf.OnePassword.ProxmoxToken
		proxmoxToken, err := ReadSecretFrom1Password(ctx, key)
		if err != nil {
			log.Errorf("unable to read proxmox token, %v", err)
			panic("1Password token or key not found")
		}
		conf.Proxmox.Pass = proxmoxToken

		key = conf.OnePassword.NatsToken
		natsToken, err := ReadSecretFrom1Password(ctx, key)
		if err != nil {
			log.Errorf("unable to read NATS token, %v", err)
			panic("1Password token or key not found")
		}
		conf.Nats.Password = natsToken
	}
}

func ReadSecretFrom1Password(ctx context.Context, secretKey string) (string, error) {
	token := os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")

	client, err := onepassword.NewClient(
		ctx,
		onepassword.WithServiceAccountToken(token),
		onepassword.WithIntegrationInfo("My 1Password i2 Integration", "v0.1.0"),
	)
	if err != nil {
		return "", err
	}
	return client.Secrets.Resolve(ctx, secretKey)
}
