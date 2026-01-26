package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/promacanthus/modix/internal/tui"
)

var (
	// Version is the version of the application
	Version = "dev"
	// ConfigFormat is the output format for config commands
	ConfigFormat string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mx",
	Short: "Modix - A tool for managing AI coding assistants and multi-agent orchestration",
	Long: `Modix is a command-line tool that unifies and manages multiple Large Language Model (LLM) vendors.
It simplifies the complexity of switching between different AI models and orchestrates multi-agent workflows.`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&ConfigFormat, "format", "f", "human", "Output format (human, json)")

	// TUI command
	RootCmd.AddCommand(&cobra.Command{
		Use:   "tui",
		Short: "Launch the interactive TUI application",
		Long:  `Launch the interactive Terminal User Interface for managing Modix projects and configurations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui.Run()
		},
	})
}
