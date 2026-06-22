package discovery

import (
	"os"
	"reconx/internal/logging"
)

func SaveResults(filename string, data []string) error {
	content := ""
	for _, item := range data {
		content += item + "\n"
	}
	return os.WriteFile(filename, []byte(content), 0644)
}

func RunSubfinder(domain string) ([]string, error) {
	logging.Log.Info("Running Subfinder", "domain", domain)
	return RunTool("subfinder", []string{"-d", domain, "-silent"})
}

func RunAssetfinder(domain string) ([]string, error) {
	logging.Log.Info("Running Assetfinder", "domain", domain)
	return RunTool("assetfinder", []string{"--subs-only", domain})
}
