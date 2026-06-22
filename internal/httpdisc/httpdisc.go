package httpdisc

import (
	"context"
	"os/exec"
	"reconx/internal/engine"
	"reconx/internal/logging"
)

type HTTPDiscoveryTask struct{}

func (t *HTTPDiscoveryTask) Name() string {
	return "HTTP Discovery"
}

func (t *HTTPDiscoveryTask) RequiredTools() []string {
	return []string{"httpx"}
}

func (t *HTTPDiscoveryTask) Execute(ctx context.Context, session *engine.Session) error {
	aliveFile := session.Config.GetPath("alive.txt")
	httpAliveFile := session.Config.GetPath("http_alive.txt")
	httpJsonFile := session.Config.GetPath("http.json")

	logging.Log.Info("Running HTTP Discovery", "input", aliveFile)

	cmd := exec.CommandContext(ctx, "httpx", "-l", aliveFile, "-silent", "-o", httpAliveFile)
	if err := cmd.Run(); err != nil {
		return err
	}

	exec.CommandContext(ctx, "httpx", "-l", aliveFile, "-json", "-o", httpJsonFile).Run()

	logging.Log.Info("HTTP discovery complete", "output", httpAliveFile)
	return nil
}
