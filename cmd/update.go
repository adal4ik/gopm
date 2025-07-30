package cmd

import (
	"fmt"
	"gopm/internal/config"
	"log"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [path_to_packages.json]",
	Short: "Downloads and unpacks packages from a remote server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		fmt.Printf("Reading update config file: %s\n", filePath)

		cfg, err := config.LoadUpdateConfig(filePath)
		if err != nil {
			log.Fatalf("Error loading update config: %v", err)
		}

		fmt.Println("------------------------------------")
		fmt.Println("Update config loaded successfully!")
		fmt.Printf("Found %d packages to process:\n", len(cfg.Packages))
		for i, pkg := range cfg.Packages {
			versionConstraint := pkg.Version
			if versionConstraint == "" {
				versionConstraint = "latest"
			}
			fmt.Printf("  %d. Package: '%s', Version constraint: '%s'\n", i+1, pkg.Name, versionConstraint)
		}
		fmt.Println("------------------------------------")
		fmt.Println("TODO: Connect to SSH, find best versions, download, and extract.")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
