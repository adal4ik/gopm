package cmd

import (
	"fmt"
	"gopm/internal/archiver"
	"gopm/internal/config"
	"gopm/internal/files"
	"log"

	"github.com/spf13/cobra"
)

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

		fmt.Printf("Packet Name: %s, Version: %s\n", cfg.Name, cfg.Version)
		fmt.Println("------------------------------------")
		fmt.Println("Searching for files to pack...")

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

		fmt.Printf("Total files to be packed: %d\n", len(allFilesToPack))
		fmt.Println("------------------------------------")

		archiveName := fmt.Sprintf("%s-%s.tar.gz", cfg.Name, cfg.Version)
		fmt.Printf("Creating archive: %s\n", archiveName)

		if err := archiver.Create(archiveName, allFilesToPack); err != nil {
			log.Fatalf("Failed to create archive: %v", err)
		}

		fmt.Println("Archive created successfully!")
		fmt.Println("------------------------------------")
		fmt.Println("TODO: Upload this archive to the SSH server.")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
