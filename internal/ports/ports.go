package ports

import (
	"context"
	"os/exec"
	"reconx/internal/engine"
	"reconx/internal/logging"
)

type PortScanTask struct{}

func (t *PortScanTask) Name() string {
	return "Port Discovery"
}

func (t *PortScanTask) RequiredTools() []string {
	return []string{"naabu"}
}

func (t *PortScanTask) Execute(ctx context.Context, session *engine.Session) error {
	aliveFile := session.Config.GetPath("alive.txt")
	portsJsonFile := session.Config.GetPath("ports.json")

	logging.Log.Info("Running Port Discovery", "input", aliveFile)

	cmd := exec.CommandContext(ctx, "naabu", "-hL", aliveFile, "-json", "-o", portsJsonFile)
	if err := cmd.Run(); err != nil {
		return err
	}

	logging.Log.Info("Port discovery complete", "output", portsJsonFile)
	return nil
}
