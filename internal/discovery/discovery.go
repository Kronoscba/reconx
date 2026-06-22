package discovery

import (
	"bufio"
	"context"
	"os/exec"
	"reconx/internal/engine"
	"reconx/internal/logging"
	"strings"
)

type DiscoveryTask struct{}

func (t *DiscoveryTask) Name() string {
	return "Asset Discovery"
}

func (t *DiscoveryTask) RequiredTools() []string {
	return []string{"subfinder", "assetfinder"}
}

func (t *DiscoveryTask) Execute(ctx context.Context, session *engine.Session) error {
	logging.Log.Info("Running Asset Discovery", "domain", session.Target)

	resultsChan := make(chan []string, 2)
	errChan := make(chan error, 2)

	go func() {
		res, err := RunSubfinder(session.Target)
		if err != nil {
			errChan <- err
			return
		}
		resultsChan <- res
	}()

	go func() {
		res, err := RunAssetfinder(session.Target)
		if err != nil {
			errChan <- err
			return
		}
		resultsChan <- res
	}()

	var allSubdomains []string
	for i := 0; i < 2; i++ {
		select {
		case res := <-resultsChan:
			allSubdomains = append(allSubdomains, res...)
		case err := <-errChan:
			logging.Log.Warn("Tool failed", "error", err)
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	session.Subdomains = Deduplicate(allSubdomains)

	subsFile := session.Config.GetPath("subdomains.txt")
	if err := SaveResults(subsFile, session.Subdomains); err != nil {
		return err
	}

	logging.Log.Info("Discovery complete", "total_subdomains", len(session.Subdomains), "file", subsFile)
	return nil
}

func RunTool(name string, args []string) ([]string, error) {
	logging.Log.Debug("Running tool", "tool", name, "args", args)

	cmd := exec.Command(name, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var results []string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			logging.Log.Debug("Found", "domain", line)
			results = append(results, line)
		}
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}

func Deduplicate(lists ...[]string) []string {
	unique := make(map[string]struct{})
	for _, list := range lists {
		for _, item := range list {
			unique[item] = struct{}{}
		}
	}

	var final []string
	for item := range unique {
		final = append(final, item)
	}
	return final
}
