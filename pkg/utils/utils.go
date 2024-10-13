package utils

import (
	"strings"
)

func GetLocalIP(ips []string) string {
	for _, ip := range ips {
		if strings.HasPrefix(ip, "192.168") {
			return ip
		}
	}
	return ""
}
