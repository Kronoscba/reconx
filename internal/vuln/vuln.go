package vuln

import (
	"context"
	"os/exec"
	"reconx/internal/engine"
	"reconx/internal/logging"
)

type VulnerabilityTask struct{}

func (t *VulnerabilityTask) Name() string {
	return "Vulnerability Discovery"
}

func (t *VulnerabilityTask) RequiredTools() []string {
	return []string{"nuclei"}
}

func (t *VulnerabilityTask) Execute(ctx context.Context, session *engine.Session) error {
	httpAliveFile := session.Config.GetPath("http_alive.txt")
	nucleiJsonFile := session.Config.GetPath("nuclei.json")

	logging.Log.Info("Running Vulnerability Discovery", "input", httpAliveFile)

	// nuclei -l <input> -severity low,medium,high,critical -json-export <output>
	cmd := exec.CommandContext(ctx, "nuclei", "-l", httpAliveFile, "-severity", "low,medium,high,critical", "-json-export", nucleiJsonFile)

	if err := cmd.Run(); err != nil {
		logging.Log.Warn("nuclei failed", "error", err)
	}

	logging.Log.Info("Vulnerability discovery complete", "output", nucleiJsonFile)
	return nil
}
