package cmd

import (
	"fmt"
	"gopm/internal/config"
	"gopm/internal/ssh"
	"gopm/internal/versioning"
	"log"

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

		fmt.Println("------------------------------------")
		fmt.Println("Connecting to SSH server...")
		sshClient, err := ssh.NewClient(updateSshHost, updateSshUser, updateSshPass, updateSshPort)
		if err != nil {
			log.Fatalf("Failed to connect to SSH server: %v", err)
		}
		defer sshClient.Close()

		fmt.Printf("Fetching file list from %s...\n", updateSshDir)
		remoteFiles, err := sshClient.ListFiles(updateSshDir)
		if err != nil {
			log.Fatalf("Failed to get file list from server: %v", err)
		}
		fmt.Printf("Found %d files on server.\n", len(remoteFiles))
		fmt.Println("------------------------------------")

		// Проходим по каждому пакету из нашего конфига
		for _, pkg := range cfg.Packages {
			versionConstraint := pkg.Version
			if versionConstraint == "" {
				versionConstraint = "latest"
			}
			fmt.Printf("Processing package '%s', version '%s'...\n", pkg.Name, versionConstraint)

			// Ищем лучший файл на сервере
			bestFile, err := versioning.FindBestMatch(remoteFiles, pkg.Name, pkg.Version)
			if err != nil {
				log.Printf("  - WARN: %v\n", err) // Не падаем, а просто предупреждаем
				continue
			}
			fmt.Printf("  - Found best match: %s\n", bestFile)
			fmt.Printf("  - TODO: Download %s\n", bestFile)
		}
		fmt.Println("------------------------------------")
		fmt.Println("TODO: Unpack downloaded archives.")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Добавляем флаги для SSH
	updateCmd.Flags().StringVar(&updateSshHost, "host", "localhost", "SSH server host")
	updateCmd.Flags().IntVar(&updateSshPort, "port", 22, "SSH server port")
	updateCmd.Flags().StringVar(&updateSshUser, "user", "", "SSH user")
	updateCmd.Flags().StringVar(&updateSshPass, "pass", "", "SSH password")
	updateCmd.Flags().StringVar(&updateSshDir, "dir", "/upload", "Remote directory with packages")

	updateCmd.MarkFlagRequired("user")
	updateCmd.MarkFlagRequired("pass")
}
