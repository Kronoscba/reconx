package deps

import (
	"fmt"
	"os/exec"
	"reconx/internal/logging"
)

type Tool struct {
	Name     string
	Install  string
	Optional bool
}

var RequiredTools = []Tool{
	{"subfinder", "go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest", false},
	{"assetfinder", "go install -v github.com/owisp/assetfinder@latest", false},
	{"dnsx", "go install -v github.com/projectdiscovery/dnsx/cmd/dnsx@latest", false},
	{"naabu", "go install -v github.com/projectdiscovery/naabu/cmd/naabu@latest", false},
	{"httpx", "go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest", false},
	{"katana", "go install -v github.com/projectdiscovery/katana/cmd/katana@latest", false},
	{"gau", "go install -v github.com/lc/gau/v2/cmd/gau@latest", false},
	{"waymore", "go install -v github.com/xortalj/waymore@latest", true},
	{"xnLinkFinder", "go install -v github.com/tomnomnom/linkfinder@latest", true}, // Note: might need rename
	{"trufflehog", "", false}, // TruffleHog is usually a binary download, not go install
	{"cloud-enum", "", true},  // Python based, install manually via uv
	{"uncover", "go install -v github.com/projectdiscovery/uncover@latest", false},
	{"gowitness", "go install -v github.com/sensepost/gowitness@latest", false},
	{"nuclei", "go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest", false},
}

func EnsureDependencies() error {
	logging.Log.Info("Checking dependencies...")

	for _, tool := range RequiredTools {
		_, err := exec.LookPath(tool.Name)
		if err != nil {
			if tool.Optional {
				logging.Log.Debug("Optional tool missing, skipping", "tool", tool.Name)
				continue
			}

			logging.Log.Warn("Missing required tool", "tool", tool.Name)
			if tool.Install == "" {
				return fmt.Errorf("tool %s is required but has no automated install command. Please install it manually", tool.Name)
			}

			logging.Log.Info("Attempting to install", "tool", tool.Name)
			if err := installTool(tool.Install); err != nil {
				return fmt.Errorf("failed to install %s: %w", tool.Name, err)
			}
		}
	}

	logging.Log.Info("All dependencies are satisfied")
	return nil
}

func installTool(installCmd string) error {
	cmd := exec.Command("sh", "-c", installCmd)
	return cmd.Run()
}
