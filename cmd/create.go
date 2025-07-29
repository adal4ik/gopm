package cmd

import (
	"fmt"
	"gopm/internal/config"
	"gopm/internal/files" // <-- Импортируем новый пакет
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

			if len(foundFiles) > 0 {
				fmt.Printf("Found %d file(s) for target '%s':\n", len(foundFiles), target.Path)
				for _, f := range foundFiles {
					fmt.Printf("  - %s\n", f)
				}
				allFilesToPack = append(allFilesToPack, foundFiles...)
			} else {
				fmt.Printf("No files found for target '%s'\n", target.Path)
			}
		}

		fmt.Println("------------------------------------")
		fmt.Printf("Total files to be packed: %d\n", len(allFilesToPack))
		fmt.Println("TODO: Archive these files into a .tar.gz package.")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
