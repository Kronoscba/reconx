package visual

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"reconx/internal/engine"
	"reconx/internal/logging"
)

type VisualReconTask struct{}

func (t *VisualReconTask) Name() string {
	return "Visual Reconnaissance"
}

func (t *VisualReconTask) RequiredTools() []string {
	return []string{"gowitness"}
}

func (t *VisualReconTask) Execute(ctx context.Context, session *engine.Session) error {
	httpAliveFile := session.Config.GetPath("http_alive.txt")
	screenshotDir := session.Config.GetPath("screenshots")
	jsonlFile := session.Config.GetPath("gowitness.jsonl")

	logging.Log.Info("Running Visual Reconnaissance", "input", httpAliveFile)

	if err := os.MkdirAll(screenshotDir, 0755); err != nil {
		return fmt.Errorf("failed to create screenshots directory: %w", err)
	}

	// gowitness v3: gowitness scan file -f <input> --screenshot-path <dir> --write-jsonl --write-jsonl-file <jsonl>
	cmd := exec.CommandContext(ctx, "gowitness", "scan", "file", "-f", httpAliveFile, "--screenshot-path", screenshotDir, "--write-jsonl", "--write-jsonl-file", jsonlFile)

	if err := cmd.Run(); err != nil {
		logging.Log.Warn("gowitness failed", "error", err)
	}

	logging.Log.Info("Visual reconnaissance complete", "screenshots", screenshotDir, "data", jsonlFile)
	return nil
}
