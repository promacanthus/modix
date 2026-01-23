package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Init initializes a new modix project
func Init(format string) error {
	// Check if .modix directory already exists
	modixDir := ".modix"
	if _, err := os.Stat(modixDir); err == nil {
		return fmt.Errorf("modix project already initialized in current directory")
	}

	// Create .modix directory
	if err := os.MkdirAll(modixDir, 0755); err != nil {
		return fmt.Errorf("failed to create .modix directory: %w", err)
	}

	// Create configuration files with default values
	now := time.Now().Format(time.RFC3339)

	shells := shellsConfig{
		Version: "v1",
		Shells:  make(map[string]shellEntry),
	}

	brains := brainsConfig{
		Version: "v1",
		Brains:  make(map[string]brainEntry),
	}

	agents := agentsConfig{
		Version: "v1",
		Agents:  make(map[string]agentEntry),
	}

	runtimes := runtimesConfig{
		Version:  "v1",
		Runtimes: make(map[string]runtimeEntry),
	}

	projects := projectsConfig{
		Version:  "v1",
		Projects: make(map[string]projectEntry),
	}

	state := stateConfig{
		LastUpdated: now,
		Counts: map[string]int{
			"shells":   0,
			"brains":   0,
			"agents":   0,
			"runtimes": 0,
			"projects": 0,
		},
		History: make([]historyEntry, 0),
	}

	version := versionConfig{
		Version:       "1.0.0",
		ConfigVersion: "v1",
		CreatedAt:     now,
	}

	files := map[string]interface{}{
		"shells.json":   shells,
		"brains.json":   brains,
		"agents.json":   agents,
		"runtimes.json": runtimes,
		"projects.json": projects,
		"state.json":    state,
		"version.json":  version,
	}

	for filename, config := range files {
		if err := writeConfigFile(filepath.Join(modixDir, filename), config); err != nil {
			return fmt.Errorf("failed to create %s: %w", filename, err)
		}
	}

	// Output result
	if format == "json" {
		result := map[string]interface{}{
			"success": true,
			"message": "Modix project initialized successfully",
			"files":   []string{"shells.json", "brains.json", "agents.json", "runtimes.json", "projects.json", "state.json", "version.json"},
		}
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonData))
	} else {
		fmt.Println("✓ Modix project initialized successfully")
		fmt.Println("  Created .modix/ directory with configuration files:")
		fmt.Println("    - shells.json")
		fmt.Println("    - brains.json")
		fmt.Println("    - agents.json")
		fmt.Println("    - runtimes.json")
		fmt.Println("    - projects.json")
		fmt.Println("    - state.json")
		fmt.Println("    - version.json")
	}

	return nil
}

// writeConfigFile writes a config struct to a JSON file
func writeConfigFile(path string, config interface{}) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Configuration structs
type shellsConfig struct {
	Version string                 `json:"version"`
	Shells  map[string]shellEntry  `json:"shells"`
}

type shellEntry struct {
	Name            string   `json:"name"`
	Binary          string   `json:"binary"`
	RequiredVersion string   `json:"requiredVersion"`
	Capabilities    []string `json:"capabilities"`
	Status          string   `json:"status"`
	RegisteredAt    string   `json:"registeredAt"`
}

type brainsConfig struct {
	Version string                 `json:"version"`
	Brains  map[string]brainEntry  `json:"brains"`
}

type brainEntry struct {
	Provider string                 `json:"provider"`
	Model    string                 `json:"model"`
	Endpoint string                 `json:"endpoint"`
	Auth     map[string]string      `json:"auth"`
	Params   map[string]interface{} `json:"params"`
	Status   string                 `json:"status"`
	CreatedAt string                `json:"createdAt"`
}

type agentsConfig struct {
	Version string                 `json:"version"`
	Agents  map[string]agentEntry  `json:"agents"`
}

type agentEntry struct {
	Role            string            `json:"role"`
	Description     string            `json:"description"`
	ManifestVersion string            `json:"manifestVersion"`
	Defaults        map[string]string `json:"defaults"`
	Capabilities    []string          `json:"capabilities"`
	DefinedAt       string            `json:"definedAt"`
}

type runtimesConfig struct {
	Version  string                      `json:"version"`
	Runtimes map[string]runtimeEntry     `json:"runtimes"`
}

type runtimeEntry struct {
	Agent      string                 `json:"agent"`
	Shell      string                 `json:"shell"`
	Brain      string                 `json:"brain"`
	Status     string                 `json:"status"`
	Validation map[string]string      `json:"validation"`
	ComposedAt string                 `json:"composedAt"`
}

type projectsConfig struct {
	Version  string                      `json:"version"`
	Projects map[string]projectEntry     `json:"projects"`
}

type projectEntry struct {
	Agents    []string `json:"agents"`
	Runtimes  []string `json:"runtimes"`
	CreatedAt string   `json:"createdAt"`
}

type stateConfig struct {
	LastUpdated string                 `json:"lastUpdated"`
	Counts      map[string]int         `json:"counts"`
	History     []historyEntry         `json:"history"`
}

type historyEntry struct {
	Event     string `json:"event"`
	Target    string `json:"target"`
	Timestamp string `json:"timestamp"`
}

type versionConfig struct {
	Version       string `json:"version"`
	ConfigVersion string `json:"configVersion"`
	CreatedAt     string `json:"createdAt"`
}

