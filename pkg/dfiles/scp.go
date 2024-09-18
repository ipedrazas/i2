package dfiles

import (
	"context"
	"fmt"
	"os"

	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

// SCPTransfer copies a file to a remote server using SCP
func SCPTransfer(filePath, remoteUser, remoteHost, remoteDir, privateKeyPath string) error {
	// Read private key
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key: %v", err)
	}

	// Create signer
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	// Create client config
	clientConfig := &ssh.ClientConfig{
		User: remoteUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Create a new SCP client
	client := scp.NewClient(fmt.Sprintf("%s:22", remoteHost), clientConfig)

	// Connect to the remote server
	err = client.Connect()
	if err != nil {
		return fmt.Errorf("couldn't establish a connection to the remote server: %v", err)
	}
	defer client.Close()

	// Open the local file
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer f.Close()

	// Get file info
	fileInfo, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	// Copy the file to the remote server
	remoteFilePath := fmt.Sprintf("%s/%s", remoteDir, fileInfo.Name())
	err = client.CopyFile(context.Background(), f, remoteFilePath, fileInfo.Mode().String())
	if err != nil {
		return fmt.Errorf("error while copying file: %v", err)
	}

	fmt.Printf("File transferred successfully to %s:%s\n", remoteHost, remoteFilePath)
	return nil
}
