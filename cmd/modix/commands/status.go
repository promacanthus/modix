package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/promacanthus/modix/internal/config"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current model status",
	Long:  `Show the currently active model and its configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStatus()
	},
}

func runStatus() error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if currentModel, modelConfig, exists := modixConfig.GetCurrentModel(); exists {
		bold.Println("Current model status:")
		fmt.Println()
		color.Cyan("Current model: %s\n", *currentModel)
		color.Cyan("Current vendor: %s\n", modixConfig.CurrentVendor)
		color.Cyan("Company: %s\n", modelConfig.Company)
		color.Cyan("API Endpoint: %s\n", modelConfig.APIEndpoint)
	} else {
		red.Println("No current model configured")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(statusCmd)
}