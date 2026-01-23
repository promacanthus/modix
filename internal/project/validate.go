package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Validate validates the modix project configuration
func Validate(format string) error {
	modixDir := ".modix"

	// Check if .modix directory exists
	if _, err := os.Stat(modixDir); os.IsNotExist(err) {
		return fmt.Errorf("modix project not initialized: .modix directory not found")
	}

	// List of required configuration files
	requiredFiles := []string{
		"shells.json",
		"brains.json",
		"agents.json",
		"runtimes.json",
		"projects.json",
		"state.json",
		"version.json",
	}

	results := make([]map[string]interface{}, 0)
	allValid := true

	for _, filename := range requiredFiles {
		result := validateFile(filepath.Join(modixDir, filename), filename)
		results = append(results, result)
		if !result["valid"].(bool) {
			allValid = false
		}
	}

	// Output result
	if format == "json" {
		output := map[string]interface{}{
			"success": allValid,
			"files":   results,
		}
		jsonData, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(jsonData))
	} else {
		fmt.Println("Configuration Validation Results:")
		fmt.Println("==================================")
		for _, result := range results {
			valid := result["valid"].(bool)
			filename := result["filename"].(string)
			errorMsg := result["error"].(string)

			if valid {
				fmt.Printf("✓ %s\n", filename)
			} else {
				fmt.Printf("✗ %s: %s\n", filename, errorMsg)
			}
		}
		fmt.Println()
		if allValid {
			fmt.Println("✓ All configuration files are valid")
		} else {
			fmt.Println("✗ Some configuration files have errors")
		}
	}

	return nil
}

// validateFile validates a single configuration file
func validateFile(path string, filename string) map[string]interface{} {
	result := map[string]interface{}{
		"filename": filename,
		"valid":    false,
		"error":    "",
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result["error"] = "File not found"
		return result
	}

	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		result["error"] = fmt.Sprintf("Failed to read file: %v", err)
		return result
	}

	// Validate JSON format
	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		result["error"] = fmt.Sprintf("Invalid JSON format: %v", err)
		return result
	}

	// Validate required fields based on file type
	switch filename {
	case "shells.json":
		if !validateShellsJSON(data) {
			result["error"] = "Invalid shells.json format"
			return result
		}
	case "brains.json":
		if !validateBrainsJSON(data) {
			result["error"] = "Invalid brains.json format"
			return result
		}
	case "agents.json":
		if !validateAgentsJSON(data) {
			result["error"] = "Invalid agents.json format"
			return result
		}
	case "runtimes.json":
		if !validateRuntimesJSON(data) {
			result["error"] = "Invalid runtimes.json format"
			return result
		}
	case "projects.json":
		if !validateProjectsJSON(data) {
			result["error"] = "Invalid projects.json format"
			return result
		}
	case "state.json":
		if !validateStateJSON(data) {
			result["error"] = "Invalid state.json format"
			return result
		}
	case "version.json":
		if !validateVersionJSON(data) {
			result["error"] = "Invalid version.json format"
			return result
		}
	}

	result["valid"] = true
	return result
}

// validateShellsJSON validates shells.json format
func validateShellsJSON(data interface{}) bool {
	m, ok := data.(map[string]interface{})
	if !ok {
		return false
	}
	_, hasVersion := m["version"]
	_, hasShells := m["shells"]
	return hasVersion && hasShells
}

// validateBrainsJSON validates brains.json format
func validateBrainsJSON(data interface{}) bool {
	m, ok := data.(map[string]interface{})
	if !ok {
		return false
	}
	_, hasVersion := m["version"]
	_, hasBrains := m["brains"]
	return hasVersion && hasBrains
}

// validateAgentsJSON validates agents.json format
func validateAgentsJSON(data interface{}) bool {
	m, ok := data.(map[string]interface{})
	if !ok {
		return false
	}
	_, hasVersion := m["version"]
	_, hasAgents := m["agents"]
	return hasVersion && hasAgents
}

// validateRuntimesJSON validates runtimes.json format
func validateRuntimesJSON(data interface{}) bool {
	m, ok := data.(map[string]interface{})
	if !ok {
		return false
	}
	_, hasVersion := m["version"]
	_, hasRuntimes := m["runtimes"]
	return hasVersion && hasRuntimes
}

// validateProjectsJSON validates projects.json format
func validateProjectsJSON(data interface{}) bool {
	m, ok := data.(map[string]interface{})
	if !ok {
		return false
	}
	_, hasVersion := m["version"]
	_, hasProjects := m["projects"]
	return hasVersion && hasProjects
}

// validateStateJSON validates state.json format
func validateStateJSON(data interface{}) bool {
	m, ok := data.(map[string]interface{})
	if !ok {
		return false
	}
	_, hasLastUpdated := m["lastUpdated"]
	_, hasCounts := m["counts"]
	_, hasHistory := m["history"]
	return hasLastUpdated && hasCounts && hasHistory
}

// validateVersionJSON validates version.json format
func validateVersionJSON(data interface{}) bool {
	m, ok := data.(map[string]interface{})
	if !ok {
		return false
	}
	_, hasVersion := m["version"]
	_, hasConfigVersion := m["configVersion"]
	_, hasCreatedAt := m["createdAt"]
	return hasVersion && hasConfigVersion && hasCreatedAt
}
