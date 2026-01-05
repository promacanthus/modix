package commands

import "github.com/spf13/cobra"

// Exported commands for main package to use
var (
	VendorCmd *cobra.Command
	ModelCmd  *cobra.Command
)

func init() {
	VendorCmd = vendorCmd
	ModelCmd = modelCmd
}

// AddLLMSubcommands adds LLM-related subcommands to the llm command
func AddLLMSubcommands(llmCmd *cobra.Command) {
	llmCmd.AddCommand(vendorCmd)
	llmCmd.AddCommand(modelCmd)
}
