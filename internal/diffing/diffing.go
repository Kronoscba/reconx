package diffing

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"reconx/internal/engine"
	"reconx/internal/logging"
)

type DiffResult struct {
	NewHosts []string `json:"new_hosts"`
	Removed  []string `json:"removed_hosts"`
	NewTech  []string `json:"new_technologies"`
	NewPorts int      `json:"new_ports"`
}

type DiffingTask struct{}

func (t *DiffingTask) Name() string {
	return "Attack Surface Diffing"
}

func (t *DiffingTask) RequiredTools() []string {
	return nil
}

func (t *DiffingTask) Execute(ctx context.Context, session *engine.Session) error {
	currentDir := session.Config.OutputDir
	parentDir := filepath.Dir(currentDir)
	domain := session.Target

	// Look for a previous scan snapshot in the parent directory
	backupDir := filepath.Join(parentDir, domain+"_previous")

	logging.Log.Info("Running Attack Surface Diffing", "current", currentDir, "previous", backupDir)

	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		logging.Log.Info("No previous scan found, saving current as baseline", "backup", backupDir)
		return saveSnapshot(currentDir, backupDir)
	}

	diff := compareScans(currentDir, backupDir)
	if err := saveDiff(currentDir, diff); err != nil {
		return err
	}

	// Update baseline for next comparison
	if err := saveSnapshot(currentDir, backupDir); err != nil {
		return err
	}

	logging.Log.Info("Diff complete", "new_hosts", len(diff.NewHosts), "removed", len(diff.Removed))
	return nil
}

func saveSnapshot(src, dst string) error {
	os.RemoveAll(dst)
	os.MkdirAll(dst, 0755)

	for _, f := range []string{"alive.txt", "http_alive.txt", "ports.json", "nuclei.json"} {
		data, err := os.ReadFile(filepath.Join(src, f))
		if err != nil {
			continue
		}
		os.WriteFile(filepath.Join(dst, f), data, 0644)
	}
	return nil
}

func compareScans(current, previous string) DiffResult {
	var diff DiffResult

	diff.NewHosts = diffFiles(filepath.Join(current, "alive.txt"), filepath.Join(previous, "alive.txt"))
	diff.Removed = diffFiles(filepath.Join(previous, "alive.txt"), filepath.Join(current, "alive.txt"))

	diff.NewTech = diffFiles(filepath.Join(current, "http_alive.txt"), filepath.Join(previous, "http_alive.txt"))

	return diff
}

func diffFiles(new, old string) []string {
	newSet := readLines(new)
	oldSet := readLines(old)

	var diff []string
	for item := range newSet {
		if !contains(oldSet, item) {
			diff = append(diff, item)
		}
	}
	return diff
}

func readLines(path string) map[string]struct{} {
	f, err := os.Open(path)
	if err != nil {
		return map[string]struct{}{}
	}
	defer f.Close()

	set := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			set[line] = struct{}{}
		}
	}
	return set
}

func contains(set map[string]struct{}, item string) bool {
	_, ok := set[item]
	return ok
}

func saveDiff(outDir string, diff DiffResult) error {
	path := filepath.Join(outDir, "diff.md")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "# Attack Surface Diff Report\n\n")
	fmt.Fprintf(f, "## Summary\n")
	fmt.Fprintf(f, "- New hosts: %d\n", len(diff.NewHosts))
	fmt.Fprintf(f, "- Removed hosts: %d\n", len(diff.Removed))
	fmt.Fprintf(f, "- New technologies: %d\n\n", len(diff.NewTech))

	if len(diff.NewHosts) > 0 {
		fmt.Fprintf(f, "## New Hosts\n")
		for _, h := range diff.NewHosts {
			fmt.Fprintf(f, "- %s\n", h)
		}
		fmt.Fprintln(f)
	}

	if len(diff.Removed) > 0 {
		fmt.Fprintf(f, "## Removed Hosts\n")
		for _, h := range diff.Removed {
			fmt.Fprintf(f, "- %s\n", h)
		}
		fmt.Fprintln(f)
	}

	return nil
}
