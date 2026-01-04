package commands

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Force colored output for all terminals
	bold   = color.New(color.FgCyan, color.Bold)
	blue   = color.New(color.FgBlue)
	green  = color.New(color.FgGreen, color.Bold)
	yellow = color.New(color.FgYellow, color.Bold)
	red    = color.New(color.FgRed, color.Bold)
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "modix",
	Short: "CLI tool for managing and switching between Claude API backends and other LLMs",
	Long: `Modix is a CLI tool for managing and switching between Claude API backends and other LLMs.

This tool helps you manage different AI model configurations and easily switch between them.`,
	Version: "0.1.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Add subcommands
	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(checkCmd)
	RootCmd.AddCommand(initCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(statusCmd)
	RootCmd.AddCommand(switchCmd)
	RootCmd.AddCommand(pathCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}