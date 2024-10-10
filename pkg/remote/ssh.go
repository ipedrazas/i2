package remote

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// ExecuteSSHCommand executes a command over SSH
func ExecuteSSHCommand(user, host, command, privateKeyPath string) (string, error) {
	// Read the SSH private key
	key, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", fmt.Errorf("unable to parse private key: %v", err)
	}

	// Create SSH client config
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the remote server
	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return "", fmt.Errorf("unable to connect: %v", err)
	}
	defer client.Close()

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("unable to create session: %v", err)
	}
	defer session.Close()

	// Execute the command
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	err = session.Run(command)
	if err != nil {
		return "", fmt.Errorf("command execution error: %v\nStderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
