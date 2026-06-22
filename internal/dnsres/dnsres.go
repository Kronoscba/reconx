package dnsres

import (
	"context"
	"os/exec"
	"reconx/internal/engine"
	"reconx/internal/logging"
)

type DNSResTask struct{}

func (t *DNSResTask) Name() string {
	return "DNS Resolution"
}

func (t *DNSResTask) RequiredTools() []string {
	return []string{"dnsx"}
}

func (t *DNSResTask) Execute(ctx context.Context, session *engine.Session) error {
	subsFile := session.Config.GetPath("subdomains.txt")
	aliveFile := session.Config.GetPath("alive.txt")
	dnsJsonFile := session.Config.GetPath("dns.json")

	logging.Log.Info("Running DNS Resolution", "input", subsFile)

	cmd := exec.CommandContext(ctx, "dnsx", "-l", subsFile, "-silent", "-o", aliveFile)
	if err := cmd.Run(); err != nil {
		return err
	}

	exec.CommandContext(ctx, "dnsx", "-l", subsFile, "-json", "-o", dnsJsonFile).Run()

	logging.Log.Info("DNS resolution complete", "output", aliveFile)
	return nil
}
