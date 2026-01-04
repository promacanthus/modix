package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Show version information, build details, and environment information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runVersion()
	},
}

func runVersion() error {
	bold.Println("Modix Version Information")
	fmt.Println()

	// Get version information
	info := GetVersionInfo()

	// Version information
	blue.Printf("Version: %s\n", info.Version)
	blue.Printf("Build: %s\n", info.Build)
	blue.Printf("Commit: %s\n", info.Commit)
	blue.Printf("Branch: %s\n", info.Branch)

	fmt.Println()

	// Environment information
	bold.Println("Environment Information")
	blue.Printf("Go Version: %s\n", runtime.Version())
	blue.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	blue.Printf("Compiler: %s\n", runtime.Compiler)

	return nil
}
