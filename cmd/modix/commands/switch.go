package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/promacanthus/modix/internal/config"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch [model-name]",
	Short: "Switch to a different model",
	Long: `Switch to a different model.

Examples:
  modix switch "claude-official"
  modix switch "deepseek-reasoner"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		return runSwitch(modelName)
	},
}

func runSwitch(modelName string) error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Find the model by searching all vendors for this model name
	modelInfos := modixConfig.GetAllModelInfos()
	var foundVendor string

	for _, modelInfo := range modelInfos {
		for _, model := range modelInfo.Models {
			if model == modelName {
				foundVendor = modelInfo.Vendor
				break
			}
		}
		if foundVendor != "" {
			break
		}
	}

	if foundVendor != "" {
		if err := modixConfig.SetCurrentVendorAndModel(foundVendor, modelName); err != nil {
			return err
		}

		if err := config.SaveConfig(modixConfig); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		// Update Claude configuration based on the selected model
		if err := config.UpdateClaudeEnvConfig(modelName, foundVendor); err != nil {
			return fmt.Errorf("failed to update Claude configuration: %w", err)
		}

		fmt.Printf("Switched to model: %s\n", modelName)
	} else {
		return fmt.Errorf("model '%s' not found in any vendor configuration", modelName)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(switchCmd)
}