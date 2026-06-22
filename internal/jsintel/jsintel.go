package jsintel

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"reconx/internal/engine"
	"reconx/internal/logging"
	"strings"
)

type JSIntelTask struct{}

func (t *JSIntelTask) Name() string {
	return "JS Intelligence"
}

func (t *JSIntelTask) RequiredTools() []string {
	return []string{"xnLinkFinder"}
}

func (t *JSIntelTask) Execute(ctx context.Context, session *engine.Session) error {
	crawlJsonFile := session.Config.GetPath("crawl.json")
	jsFile := session.Config.GetPath("javascript.txt")
	
	logging.Log.Info("Running JS Intelligence", "input", crawlJsonFile)
	
	jsList, err := extractJSUrls(crawlJsonFile)
	if err != nil {
		return fmt.Errorf("failed to extract JS URLs: %w", err)
	}
	defer os.Remove(jsList)

	cmd := exec.CommandContext(ctx, "xnLinkFinder", "-i", jsList, "-o", jsFile)
	if err := cmd.Run(); err != nil {
		return err
	}
	
	logging.Log.Info("JS intelligence complete", "output", jsFile)
	return nil
}

func extractJSUrls(crawlFile string) (string, error) {
	jsListFile := crawlFile + ".js_list"
	file, err := os.Open(crawlFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	outFile, err := os.Create(jsListFile)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ".js") {
			start := strings.Index(line, "https://")
			if start == -1 {
				start = strings.Index(line, "http://")
			}
			if start != -1 {
				end := strings.Index(line[start:], "\"")
				if end != -1 {
					url := line[start : start+end]
					fmt.Fprintln(outFile, url)
					count++
				}
			}
		}
	}
	
	logging.Log.Debug("Extracted JS files", "count", count, "file", jsListFile)
	return jsListFile, nil
}
