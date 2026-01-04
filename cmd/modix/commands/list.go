package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/promacanthus/modix/internal/config"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured models",
	Long:  `List all configured models with their status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runList()
	},
}

func runList() error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Get all models and sort them by vendor and model name for consistent output
	modelInfos := modixConfig.GetAllModelInfos()
	sort.Slice(modelInfos, func(i, j int) bool {
		if modelInfos[i].Vendor == modelInfos[j].Vendor {
			return modelInfos[i].Models[0] < modelInfos[j].Models[0]
		}
		return modelInfos[i].Vendor < modelInfos[j].Vendor
	})

	// Print header
	bold.Printf("%-35s %-15s %-15s %-10s %-10s\n",
		"MODEL", "COMPANY", "VENDOR", "ENDPOINT", "API_KEY")
	fmt.Printf("%s %s %s %s %s\n",
		strings.Repeat("-", 35),
		strings.Repeat("-", 15),
		strings.Repeat("-", 15),
		strings.Repeat("-", 10),
		strings.Repeat("-", 10))

	// Print models
	for _, modelInfo := range modelInfos {
		for _, modelName := range modelInfo.Models {
			// Show endpoint status - special handling for Anthropic
			endpointDisplay := ""
			if strings.ToLower(modelInfo.Vendor) == "anthropic" {
				endpointDisplay = blue.Sprint("[ - ]")
			} else if modelInfo.Endpoint == "" {
				endpointDisplay = red.Sprint("[ N ]")
			} else {
				endpointDisplay = green.Sprint("[ Y ]")
			}

			// Show API key status - special handling for Anthropic
			apiKeyDisplay := ""
			if strings.ToLower(modelInfo.Vendor) == "anthropic" {
				apiKeyDisplay = blue.Sprint("[ - ]")
			} else if !modelInfo.HasAPIKey {
				apiKeyDisplay = red.Sprint("[ N ]")
			} else {
				apiKeyDisplay = green.Sprint("[ Y ]")
			}

			// Highlight current model
			modelDisplay := modelName
			if modixConfig.CurrentVendor == modelInfo.Vendor && modixConfig.CurrentModel == modelName {
				modelDisplay = yellow.Sprint(modelName)
			} else {
				modelDisplay = modelName
			}

			// Highlight current vendor
			vendorDisplay := modelInfo.Vendor
			if modixConfig.CurrentVendor == modelInfo.Vendor {
				vendorDisplay = yellow.Sprint(modelInfo.Vendor)
			} else {
				vendorDisplay = modelInfo.Vendor
			}

			blue.Printf("%-35s %-15s %-15s %-10s %-10s\n",
				modelDisplay,
				modelInfo.Company,
				vendorDisplay,
				endpointDisplay,
				apiKeyDisplay)
		}
	}

	// Show summary information
	totalModels := len(modelInfos)
	configuredModels := 0
	for _, modelInfo := range modelInfos {
		if modelInfo.HasAPIKey && modelInfo.HasEndpoint {
			configuredModels++
		}
	}

	currentModelInfo := "None"
	if currentModel, _, exists := modixConfig.GetCurrentModel(); exists {
		currentModelInfo = fmt.Sprintf("%s@%s", modixConfig.CurrentVendor, *currentModel)
	}

	fmt.Println()
	bold.Println("--- Summary ---")
	blue.Printf("%-20s %s\n", "Total models:", yellow.Sprint(totalModels))
	blue.Printf("%-20s %s\n", "Configured models:", green.Sprint(configuredModels))
	blue.Printf("%-20s %s\n", "Current model:", yellow.Sprint(currentModelInfo))

	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)
}