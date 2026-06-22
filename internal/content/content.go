package content

import (
	"context"
	"os/exec"
	"reconx/internal/engine"
	"reconx/internal/logging"
)

type ContentDiscoveryTask struct{}

func (t *ContentDiscoveryTask) Name() string {
	return "Content Discovery"
}

func (t *ContentDiscoveryTask) RequiredTools() []string {
	return []string{"katana"}
}

func (t *ContentDiscoveryTask) Execute(ctx context.Context, session *engine.Session) error {
	httpAliveFile := session.Config.GetPath("http_alive.txt")
	crawlJsonFile := session.Config.GetPath("crawl.json")

	logging.Log.Info("Running Content Discovery", "input", httpAliveFile)

	cmd := exec.CommandContext(ctx, "katana", "-list", httpAliveFile, "-jsonl", "-o", crawlJsonFile)
	if err := cmd.Run(); err != nil {
		return err
	}

	logging.Log.Info("Content discovery complete", "output", crawlJsonFile)
	return nil
}
