package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gopm",
	Short: "gopm is a simple package manager for transferring files over SSH",
	Long: `A command-line tool to pack files into an archive based on a manifest file,
upload it to a remote server via SSH, and vice versa`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
