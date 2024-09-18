package dckr

import (
	"os"
	"path/filepath"
	"strings"
)

func GetDockerSocketPath() string {
	defaultSocket := "unix:///var/run/docker.sock"
	homeDir, err := os.UserHomeDir()
	if err == nil {
		customSocket := filepath.Join(homeDir, ".docker", "run", "docker.sock")
		if _, err := os.Stat(customSocket); err != nil {
			return defaultSocket // Default Unix socket
		}
		return "unix:///" + customSocket
	}
	return defaultSocket
}

// GetContainerName returns the first name of a container without the leading slash
func GetContainerName(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return strings.TrimPrefix(names[0], "/")
}
