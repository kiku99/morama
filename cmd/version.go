/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.0" // 기본값

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "버전 정보를 출력합니다",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("morama version", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
