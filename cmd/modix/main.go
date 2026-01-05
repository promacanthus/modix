package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/promacanthus/modix/cmd/modix/commands"
	"github.com/promacanthus/modix/cmd/modix/commands/agent"
	"github.com/promacanthus/modix/cmd/modix/commands/config"
	"github.com/promacanthus/modix/cmd/modix/commands/llm"
	"github.com/spf13/cobra"
)

var (
	// Force colored output for all terminals
	bold   = color.New(color.FgCyan, color.Bold)
	blue   = color.New(color.FgBlue)
	green  = color.New(color.FgGreen, color.Bold)
	yellow = color.New(color.FgYellow, color.Bold)
	red    = color.New(color.FgRed, color.Bold)

	// Version information injected by GoReleaser
	Version = "dev"
	Build   = "dev"
	Commit  = "none"
	Branch  = "main"
)

// VersionInfo holds the version information
type VersionInfo struct {
	Version string
	Build   string
	Commit  string
	Branch  string
}

// GetVersionInfo returns the current version information
func GetVersionInfo() VersionInfo {
	return VersionInfo{
		Version: Version,
		Build:   Build,
		Commit:  Commit,
		Branch:  Branch,
	}
}

// init initializes the color package to force colored output
func init() {
	color.NoColor = false
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "modix",
	Short: "CLI tool for managing LLM models, vendors, and coding agents",
	Long: `Modix is a CLI tool for managing LLM models, vendors, and coding agents.

Manage LLM providers (vendors) and their models, switch between different
LLM backends, and configure coding agents like Claude Code.

Available command groups:
  vendor   Manage LLM vendors (add, remove, update, list, show, model)
  model    Manage and switch LLM models (list, switch, status)
  llm      Legacy LLM commands (deprecated)
  agent    Manage coding agents (Claude Code, Gemini CLI, etc.)
  config   Manage modix configuration

Examples:
  # Initialize configuration
  modix config init

  # Add a vendor
  modix vendor add deepseek --company "DeepSeek" --endpoint "https://api.deepseek.com/v1" --api-key "sk-xxx"

  # Add a model to vendor
  modix vendor model add deepseek deepseek-reasoner

  # List all models
  modix model list

  # Switch to a model
  modix model switch deepseek-reasoner

  # Show current model status
  modix model status

  # Configure coding agent
  modix agent add claude-code
  modix agent config claude-code

  # Check configuration
  modix config check`,
	Version: "0.2.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Add new command groups
	RootCmd.AddCommand(commands.VendorCmd)
	RootCmd.AddCommand(commands.ModelCmd)
	RootCmd.AddCommand(llm.LLMCmd)
	RootCmd.AddCommand(agent.AgentCmd)
	RootCmd.AddCommand(config.ConfigCmd)

	// Add version command
	RootCmd.AddCommand(versionCmd)

	// Add legacy commands for backward compatibility (optional)
	// You can remove these later once users migrate
	addLegacyCommands()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// addLegacyCommands adds backward-compatible commands
// TODO: Consider removing these after migration period
func addLegacyCommands() {
	// Create aliases for common operations
	legacyCmd := &cobra.Command{
		Use:   "legacy",
		Short: "Legacy commands (deprecated)",
		Long:  `Legacy commands for backward compatibility. Use new commands instead.`,
		Hidden: true,
	}
	RootCmd.AddCommand(legacyCmd)
}

func main() {
	// Handle version flag
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		info := GetVersionInfo()
		fmt.Printf("modix version %s\n", info.Version)
		if info.Build != "dev" {
			fmt.Printf("Build: %s\n", info.Build)
			fmt.Printf("Commit: %s\n", info.Commit)
			fmt.Printf("Branch: %s\n", info.Branch)
		}
		os.Exit(0)
	}

	if err := Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
