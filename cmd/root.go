package cmd

import (
	"os"

	"github.com/kiku99/morama/internal/config"
	"github.com/kiku99/morama/internal/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "morama",
	Short: "A CLI tool for managing your watched movies and dramas",
	Long: `Morama is a command-line application for tracking and managing 
your watched movies and dramas. You can add new entries with ratings 
and comments, and view them in a beautiful table format.

Examples:
  morama add "언젠가는 슬기로울 전공의생활" --drama
  morama add "인셉션" --movie  
  morama list'`,
}

// Execute initializes config, logger, and runs the CLI
func Execute() {
	// Initialize configuration
	config.GetConfig()

	// Initialize logger from config
	if err := utils.InitLoggerFromConfig(); err != nil {
		utils.Warning("Failed to initialize logger: %v", err)
	}

	utils.Info("Morama CLI started")
	utils.LogUserAction("app_started", "application launched")

	err := rootCmd.Execute()
	if err != nil {
		utils.Error("Command execution failed: %v", err)
		os.Exit(1)
	}

	utils.Info("Morama CLI finished")
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
