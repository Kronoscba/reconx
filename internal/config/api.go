package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SetupAPI() error {
	providers := []string{
		"shodan",
		"censys",
		"fofa",
		"hunter",
		"zoomeye",
		"netlas",
		"criminalip",
	}

	configDir := filepath.Join(os.Getenv("HOME"), ".config", "uncover")
	configFile := filepath.Join(configDir, "provider-config.yaml")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	var configBuilder strings.Builder
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("--- ReconX API Configuration ---")
	fmt.Println("Leave blank if you don't have the key. Press Enter to skip.")

	for _, p := range providers {
		fmt.Printf("Enter API Key for %s: ", p)
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)
		if key != "" {
			configBuilder.WriteString(fmt.Sprintf("%s:\n  - \"%s\"\n", p, key))
		}
	}

	return os.WriteFile(configFile, []byte(configBuilder.String()), 0644)
}
