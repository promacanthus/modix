package llm

import (
	"github.com/promacanthus/modix/cmd/modix/commands"
	"github.com/spf13/cobra"
)

// LLMCmd represents the llm command family (legacy)
var LLMCmd = &cobra.Command{
	Use:        "llm",
	Short:      "Legacy LLM commands (deprecated)",
	Long:       `Legacy LLM commands. Use top-level 'vendor' and 'model' commands instead.`,
	Deprecated: "Use 'vendor' and 'model' commands instead",
	Hidden:     false,
}

func init() {
	// Add vendor and model as subcommands for backward compatibility
	commands.AddLLMSubcommands(LLMCmd)
}
