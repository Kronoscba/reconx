package osint

import (
	"context"
	"os/exec"
	"reconx/internal/engine"
	"reconx/internal/logging"
)

type OSINTTask struct{}

func (t *OSINTTask) Name() string {
	return "OSINT Expansion"
}

func (t *OSINTTask) RequiredTools() []string {
	return []string{"uncover"}
}

func (t *OSINTTask) Execute(ctx context.Context, session *engine.Session) error {
	osintJsonFile := session.Config.GetPath("osint.json")

	logging.Log.Info("Running OSINT Expansion", "domain", session.Target)

	// uncover -q <domain> -json -o <output>
	// Use -q instead of -d
	cmd := exec.CommandContext(ctx, "uncover", "-q", session.Target, "-json", "-o", osintJsonFile)

	if err := cmd.Run(); err != nil {
		logging.Log.Warn("uncover failed or not installed", "error", err)
		// Not returning error to avoid breaking the pipeline if OSINT fails (often requires API keys)
	}

	logging.Log.Info("OSINT expansion complete", "output", osintJsonFile)
	return nil
}
