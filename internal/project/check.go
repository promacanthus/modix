package project

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// Tool represents a tool to check
type Tool struct {
	Name        string
	Binary      string
	Description string
}

// Check checks if required tools are installed
func Check(format string) error {
	tools := []Tool{
		{Name: "git", Binary: "git", Description: "Version control system"},
		{Name: "claude-code", Binary: "claude", Description: "Claude Code CLI"},
		{Name: "codex-cli", Binary: "codex", Description: "OpenAI Codex CLI"},
		{Name: "gemini-cli", Binary: "gemini", Description: "Google Gemini CLI"},
	}

	results := make([]map[string]interface{}, 0)
	allInstalled := true

	for _, tool := range tools {
		result := checkTool(tool)
		results = append(results, result)
		if !result["installed"].(bool) {
			allInstalled = false
		}
	}

	// Output result
	if format == "json" {
		output := map[string]interface{}{
			"success": allInstalled,
			"tools":   results,
		}
		jsonData, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(jsonData))
	} else {
		fmt.Println("Dependency Check Results:")
		fmt.Println("=========================")
		for _, result := range results {
			installed := result["installed"].(bool)
			name := result["name"].(string)
			desc := result["description"].(string)
			version := result["version"].(string)

			if installed {
				fmt.Printf("✓ %-15s %s (version: %s)\n", name, desc, version)
			} else {
				fmt.Printf("✗ %-15s %s (not installed)\n", name, desc)
			}
		}
		fmt.Println()
		if allInstalled {
			fmt.Println("✓ All required tools are installed")
		} else {
			fmt.Println("✗ Some required tools are missing")
		}
	}

	return nil
}

// checkTool checks if a tool is installed and returns its version
func checkTool(tool Tool) map[string]interface{} {
	result := map[string]interface{}{
		"name":        tool.Name,
		"description": tool.Description,
		"installed":   false,
		"version":     "unknown",
	}

	// Check if binary exists in PATH
	cmd := exec.Command(tool.Binary, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Try alternative command for some tools
		if tool.Name == "git" {
			cmd = exec.Command(tool.Binary, "version")
			output, err = cmd.CombinedOutput()
		}
	}

	if err == nil {
		result["installed"] = true
		version := string(output)
		// Clean up version string
		if len(version) > 0 {
			// Remove newlines and trim
			version = version[:len(version)-1]
			result["version"] = version
		}
	}

	return result
}
