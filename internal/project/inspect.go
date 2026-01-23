package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Inspect inspects the modix project configuration
func Inspect(format string) error {
	modixDir := ".modix"

	// Check if .modix directory exists
	if _, err := os.Stat(modixDir); os.IsNotExist(err) {
		return fmt.Errorf("modix project not initialized: .modix directory not found")
	}

	// Read all configuration files
	shells, err := readJSONFile(filepath.Join(modixDir, "shells.json"))
	if err != nil {
		return fmt.Errorf("failed to read shells.json: %w", err)
	}

	brains, err := readJSONFile(filepath.Join(modixDir, "brains.json"))
	if err != nil {
		return fmt.Errorf("failed to read brains.json: %w", err)
	}

	agents, err := readJSONFile(filepath.Join(modixDir, "agents.json"))
	if err != nil {
		return fmt.Errorf("failed to read agents.json: %w", err)
	}

	runtimes, err := readJSONFile(filepath.Join(modixDir, "runtimes.json"))
	if err != nil {
		return fmt.Errorf("failed to read runtimes.json: %w", err)
	}

	projects, err := readJSONFile(filepath.Join(modixDir, "projects.json"))
	if err != nil {
		return fmt.Errorf("failed to read projects.json: %w", err)
	}

	state, err := readJSONFile(filepath.Join(modixDir, "state.json"))
	if err != nil {
		return fmt.Errorf("failed to read state.json: %w", err)
	}

	version, err := readJSONFile(filepath.Join(modixDir, "version.json"))
	if err != nil {
		return fmt.Errorf("failed to read version.json: %w", err)
	}

	// Output result
	if format == "json" {
		output := map[string]interface{}{
			"success":  true,
			"shells":   shells,
			"brains":   brains,
			"agents":   agents,
			"runtimes": runtimes,
			"projects": projects,
			"state":    state,
			"version":  version,
		}
		jsonData, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(jsonData))
	} else {
		fmt.Println("Modix Project Inspection")
		fmt.Println("========================")
		fmt.Println()

		// Version
		if versionMap, ok := version.(map[string]interface{}); ok {
			fmt.Println("Version Information:")
			fmt.Printf("  Version: %v\n", versionMap["version"])
			fmt.Printf("  Config Version: %v\n", versionMap["configVersion"])
			fmt.Printf("  Created At: %v\n", versionMap["createdAt"])
			fmt.Println()
		}

		// Shells
		if shellsMap, ok := shells.(map[string]interface{}); ok {
			fmt.Println("Shell Configuration:")
			if shellsData, ok := shellsMap["shells"].(map[string]interface{}); ok {
				for name, shell := range shellsData {
					if shellMap, ok := shell.(map[string]interface{}); ok {
						fmt.Printf("  %s:\n", name)
						fmt.Printf("    Name: %v\n", shellMap["name"])
						fmt.Printf("    Binary: %v\n", shellMap["binary"])
						fmt.Printf("    Status: %v\n", shellMap["status"])
					}
				}
			}
			fmt.Println()
		}

		// Brains
		if brainsMap, ok := brains.(map[string]interface{}); ok {
			fmt.Println("Brain Configuration:")
			if brainsData, ok := brainsMap["brains"].(map[string]interface{}); ok {
				for name, brain := range brainsData {
					if brainMap, ok := brain.(map[string]interface{}); ok {
						fmt.Printf("  %s:\n", name)
						fmt.Printf("    Provider: %v\n", brainMap["provider"])
						fmt.Printf("    Model: %v\n", brainMap["model"])
						fmt.Printf("    Status: %v\n", brainMap["status"])
					}
				}
			}
			fmt.Println()
		}

		// Agents
		if agentsMap, ok := agents.(map[string]interface{}); ok {
			fmt.Println("Agent Configuration:")
			if agentsData, ok := agentsMap["agents"].(map[string]interface{}); ok {
				for name, agent := range agentsData {
					if agentMap, ok := agent.(map[string]interface{}); ok {
						fmt.Printf("  %s:\n", name)
						fmt.Printf("    Role: %v\n", agentMap["role"])
						fmt.Printf("    Description: %v\n", agentMap["description"])
						fmt.Printf("    Manifest Version: %v\n", agentMap["manifestVersion"])
					}
				}
			}
			fmt.Println()
		}

		// Runtimes
		if runtimesMap, ok := runtimes.(map[string]interface{}); ok {
			fmt.Println("Runtime Configuration:")
			if runtimesData, ok := runtimesMap["runtimes"].(map[string]interface{}); ok {
				for name, runtime := range runtimesData {
					if runtimeMap, ok := runtime.(map[string]interface{}); ok {
						fmt.Printf("  %s:\n", name)
						fmt.Printf("    Agent: %v\n", runtimeMap["agent"])
						fmt.Printf("    Shell: %v\n", runtimeMap["shell"])
						fmt.Printf("    Brain: %v\n", runtimeMap["brain"])
						fmt.Printf("    Status: %v\n", runtimeMap["status"])
					}
				}
			}
			fmt.Println()
		}

		// Projects
		if projectsMap, ok := projects.(map[string]interface{}); ok {
			fmt.Println("Project Configuration:")
			if projectsData, ok := projectsMap["projects"].(map[string]interface{}); ok {
				for name, project := range projectsData {
					if projectMap, ok := project.(map[string]interface{}); ok {
						fmt.Printf("  %s:\n", name)
						fmt.Printf("    Agents: %v\n", projectMap["agents"])
						fmt.Printf("    Runtimes: %v\n", projectMap["runtimes"])
						fmt.Printf("    Created At: %v\n", projectMap["createdAt"])
					}
				}
			}
			fmt.Println()
		}

		// State
		if stateMap, ok := state.(map[string]interface{}); ok {
			fmt.Println("State Information:")
			fmt.Printf("  Last Updated: %v\n", stateMap["lastUpdated"])
			if counts, ok := stateMap["counts"].(map[string]interface{}); ok {
				fmt.Println("  Counts:")
				for key, value := range counts {
					fmt.Printf("    %s: %v\n", key, value)
				}
			}
			fmt.Println()
		}
	}

	return nil
}

// readJSONFile reads and parses a JSON file
func readJSONFile(path string) (interface{}, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	return data, nil
}
