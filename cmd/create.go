package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [path_to_packet.json]",
	Short: "Packs files into an archive and uploads it to a remote server",
	Long: `Reads a packet manifest file (e.g., packet.json), finds files based on
the specified targets, creates a .tar.gz archive, and uploads it via SFTP.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Executing 'create' command...")
		filePath := args[0]
		fmt.Printf("Config file specified: %s\n", filePath)
		fmt.Println("TODO: Implement parsing, archiving, and uploading logic here.")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
