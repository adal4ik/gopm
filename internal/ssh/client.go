package ssh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

// NewClient создает новое SSH-соединение и SFTP-сессию.
func NewClient(host, user, pass string, port int) (*Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ssh: %w", err)
	}

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create sftp client: %w", err)
	}

	return &Client{
		sshClient:  conn,
		sftpClient: sftpClient,
	}, nil
}

// Close корректно закрывает SFTP-сессию и SSH-соединение.
func (c *Client) Close() {
	c.sftpClient.Close()
	c.sshClient.Close()
}

// UploadFile загружает один файл (localPath) в удаленную директорию (remoteDir).
func (c *Client) UploadFile(localPath, remoteDir string) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file %s: %w", localPath, err)
	}
	defer localFile.Close()

	remoteFileName := filepath.Base(localPath)
	remotePath := filepath.Join(remoteDir, remoteFileName)

	remoteFile, err := c.sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file %s: %w", remotePath, err)
	}
	defer remoteFile.Close()

	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("failed to copy content to remote file: %w", err)
	}

	return nil
}

func (c *Client) ListFiles(remoteDir string) ([]string, error) {
	fileInfos, err := c.sftpClient.ReadDir(remoteDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read remote directory %s: %w", remoteDir, err)
	}

	var fileNames []string
	for _, info := range fileInfos {
		if !info.IsDir() {
			fileNames = append(fileNames, info.Name())
		}
	}

	return fileNames, nil
}

func (c *Client) DownloadFile(remotePath, localDir string) (string, error) {
	remoteFile, err := c.sftpClient.Open(remotePath)
	if err != nil {
		return "", fmt.Errorf("failed to open remote file %s: %w", remotePath, err)
	}
	defer remoteFile.Close()

	localFileName := filepath.Base(remotePath)
	localFilePath := filepath.Join(localDir, localFileName)

	localFile, err := os.Create(localFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file %s: %w", localFilePath, err)
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return "", fmt.Errorf("failed to copy content from remote file: %w", err)
	}

	return localFilePath, nil
}
