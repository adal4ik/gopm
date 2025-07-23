package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [path_to_packages.json]",
	Short: "Downloads and unpacks packages from a remote server",
	Long: `Reads a list of required packages, finds the best matching versions on the
remote server, downloads, and extracts them.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Executing 'update' command...")
		filePath := args[0]
		fmt.Printf("Config file specified: %s\n", filePath)
		fmt.Println("TODO: Implement version searching, downloading, and extracting logic here.")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
