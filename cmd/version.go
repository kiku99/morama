/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"     // default value
	commit  = "unknown" // Git commit hash
	date    = "unknown" // build date
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("morama version %s\n", version)
		fmt.Printf("Git commit: %s\n", commit)
		fmt.Printf("Built: %s\n", date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
