package utils

import (
	"context"
	"i2/pkg/api"
	"os"
	"strings"

	"github.com/1password/onepassword-sdk-go"
)

func GetLocalIP(ips []string) string {
	for _, ip := range ips {
		if strings.HasPrefix(ip, "192.168") {
			return ip
		}
	}
	return ""
}

func ReadSecretFrom1Password(ctx context.Context, secretKey string) (string, error) {
	token := os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")

	client, err := onepassword.NewClient(
		ctx,
		onepassword.WithServiceAccountToken(token),
		onepassword.WithIntegrationInfo("My 1Password i2 Integration", api.Version),
	)
	if err != nil {
		return "", err
	}
	return client.Secrets.Resolve(ctx, secretKey)
}
