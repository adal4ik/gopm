package cmd

import (
	"fmt"
	"gopm/internal/archiver"
	"gopm/internal/config"
	"gopm/internal/files"
	"gopm/internal/ssh"
	"log"

	"github.com/spf13/cobra"
)

var sshHost, sshUser, sshPass, sshDir string
var sshPort int

var createCmd = &cobra.Command{
	Use:   "create [path_to_packet.json]",
	Short: "Packs files into an archive and uploads it to a remote server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		cfg, err := config.LoadPacketConfig(filePath)
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		var allFilesToPack []string
		for _, target := range cfg.Targets {
			foundFiles, err := files.FindFilesByTarget(target)
			if err != nil {
				log.Fatalf("Error finding files for target '%s': %v", target.Path, err)
			}
			allFilesToPack = append(allFilesToPack, foundFiles...)
		}

		if len(allFilesToPack) == 0 {
			log.Println("No files found to pack. Exiting.")
			return
		}

		archiveName := fmt.Sprintf("%s-%s.tar.gz", cfg.Name, cfg.Version)
		if err := archiver.Create(archiveName, allFilesToPack); err != nil {
			log.Fatalf("Failed to create archive: %v", err)
		}
		fmt.Printf("Archive created successfully: %s\n", archiveName)
		fmt.Println("------------------------------------")

		fmt.Println("Connecting to SSH server...")
		sshClient, err := ssh.NewClient(sshHost, sshUser, sshPass, sshPort)
		if err != nil {
			log.Fatalf("Failed to connect to SSH server: %v", err)
		}
		defer sshClient.Close()

		fmt.Printf("Uploading %s to %s on %s...\n", archiveName, sshDir, sshHost)
		if err := sshClient.UploadFile(archiveName, sshDir); err != nil {
			log.Fatalf("Failed to upload archive: %v", err)
		}

		fmt.Println("------------------------------------")
		fmt.Println("Package created and uploaded successfully!")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&sshHost, "host", "H", "", "SSH server host (e.g., localhost)")
	createCmd.Flags().StringVarP(&sshUser, "user", "u", "", "SSH user")
	createCmd.Flags().StringVarP(&sshPass, "pass", "p", "", "SSH password")
	createCmd.Flags().StringVarP(&sshDir, "dir", "d", "/upload", "Remote directory to upload packages")
	createCmd.Flags().IntVarP(&sshPort, "port", "P", 22, "SSH server port")

	createCmd.MarkFlagRequired("host")
	createCmd.MarkFlagRequired("user")
	createCmd.MarkFlagRequired("pass")
}
