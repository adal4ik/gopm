package cmd

import (
	"fmt"
	"gopm/internal/archiver"
	"gopm/internal/config"
	"gopm/internal/ssh"
	"gopm/internal/versioning"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var updateSshHost, updateSshUser, updateSshPass, updateSshDir string
var updateSshPort int
var updateCmd = &cobra.Command{
	Use:   "update [path_to_packages.json]",
	Short: "Downloads and unpacks packages from a remote server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		cfg, err := config.LoadUpdateConfig(filePath)
		if err != nil {
			log.Fatalf("Error loading update config: %v", err)
		}

		sshClient, err := ssh.NewClient(updateSshHost, updateSshUser, updateSshPass, updateSshPort)
		if err != nil {
			log.Fatalf("Failed to connect to SSH server: %v", err)
		}
		defer sshClient.Close()

		remoteFiles, err := sshClient.ListFiles(updateSshDir)
		if err != nil {
			log.Fatalf("Failed to get file list from server: %v", err)
		}

		fmt.Printf("Found %d total files on server. Processing packages...\n", len(remoteFiles))
		fmt.Println("------------------------------------")

		var downloadedPackages int
		for _, pkg := range cfg.Packages {
			bestFile, err := versioning.FindBestMatch(remoteFiles, pkg.Name, pkg.Version)
			if err != nil {
				log.Printf("WARN: For package '%s': %v\n", pkg.Name, err)
				continue
			}

			fmt.Printf("  - Found best match for '%s': %s\n", pkg.Name, bestFile)

			fmt.Printf("  - Downloading...\n")
			remoteFilePath := filepath.Join(updateSshDir, bestFile)
			downloadedArchivePath, err := sshClient.DownloadFile(remoteFilePath, ".")
			if err != nil {
				log.Printf("  - ERROR: Failed to download %s: %v\n", bestFile, err)
				continue
			}

			fmt.Printf("  - Extracting %s...\n", downloadedArchivePath)
			if err := archiver.Extract(downloadedArchivePath, "."); err != nil {
				log.Printf("  - ERROR: Failed to extract %s: %v\n", downloadedArchivePath, err)
				os.Remove(downloadedArchivePath)
				continue
			}

			if err := os.Remove(downloadedArchivePath); err != nil {
				log.Printf("  - WARN: Failed to clean up archive file %s: %v\n", downloadedArchivePath, err)
			}

			fmt.Printf("  - Successfully installed package %s version %s\n", pkg.Name, "...") // TODO: extract version
			downloadedPackages++
		}

		fmt.Println("------------------------------------")
		fmt.Printf("Update finished. Successfully installed %d package(s).\n", downloadedPackages)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&updateSshHost, "host", "localhost", "SSH server host")

	updateCmd.Flags().IntVar(&updateSshPort, "port", 22, "SSH server port")

	updateCmd.Flags().StringVar(&updateSshUser, "user", "", "SSH user")
	updateCmd.Flags().StringVar(&updateSshPass, "pass", "", "SSH password")
	updateCmd.Flags().StringVar(&updateSshDir, "dir", "/upload", "Remote directory with packages")

	updateCmd.MarkFlagRequired("user")
	updateCmd.MarkFlagRequired("pass")
}
