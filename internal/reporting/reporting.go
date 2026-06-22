package reporting

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"reconx/internal/engine"
	"reconx/internal/logging"
)

type ReportData struct {
	Domain      string
	Timestamp   string
	Duration    string
	Subdomains  int
	Alive       int
	HTTPAlive   int
	Ports       int
	Secrets     int
	Cloud       int
	OSINT       int
	Vulns       int
	Screenshots int
	SecretsFile string
	VulnsFile   string
	HasDNS      bool
	HasHTTP     bool
	HasContent  bool
	HasHistory  bool
	HasJS       bool
}

type ReportingTask struct{}

func (t *ReportingTask) Name() string {
	return "Reporting Engine"
}

func (t *ReportingTask) RequiredTools() []string {
	return nil
}

func (t *ReportingTask) Execute(ctx context.Context, session *engine.Session) error {
	outDir := session.Config.OutputDir
	logging.Log.Info("Generating reports", "output", outDir)

	data := collectData(session.Target, outDir)

	if err := generateMarkdown(outDir, data); err != nil {
		logging.Log.Warn("markdown report failed", "error", err)
	}
	if err := generateHTML(outDir, data); err != nil {
		logging.Log.Warn("html report failed", "error", err)
	}
	if err := generateJSON(outDir, data); err != nil {
		logging.Log.Warn("json report failed", "error", err)
	}

	logging.Log.Info("Reports generated")
	return nil
}

func collectData(domain, outDir string) ReportData {
	d := ReportData{
		Domain:    domain,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	d.Subdomains = countLines(filepath.Join(outDir, "subdomains.txt"))
	d.Alive = countLines(filepath.Join(outDir, "alive.txt"))
	d.HTTPAlive = countLines(filepath.Join(outDir, "http_alive.txt"))
	d.Ports = countLines(filepath.Join(outDir, "ports.json"))

	if _, err := os.Stat(filepath.Join(outDir, "dns.json")); err == nil {
		d.HasDNS = true
	}
	if _, err := os.Stat(filepath.Join(outDir, "http.json")); err == nil {
		d.HasHTTP = true
	}
	if _, err := os.Stat(filepath.Join(outDir, "crawl.json")); err == nil {
		d.HasContent = true
	}
	if _, err := os.Stat(filepath.Join(outDir, "historical.txt")); err == nil {
		d.HasHistory = true
	}
	if _, err := os.Stat(filepath.Join(outDir, "javascript.txt")); err == nil {
		d.HasJS = true
	}

	d.Secrets = fileSize(filepath.Join(outDir, "secrets.json"))
	d.Cloud = fileSize(filepath.Join(outDir, "cloud.json"))
	d.OSINT = fileSize(filepath.Join(outDir, "osint.json"))
	d.Vulns = fileSize(filepath.Join(outDir, "nuclei.json"))
	d.SecretsFile = filepath.Join(outDir, "secrets.json")
	d.VulnsFile = filepath.Join(outDir, "nuclei.json")

	screenshotDir := filepath.Join(outDir, "screenshots")
	if entries, err := os.ReadDir(screenshotDir); err == nil {
		d.Screenshots = len(entries)
	}

	return d
}

func countLines(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	if len(data) == 0 {
		return 0
	}
	return strings.Count(string(data), "\n")
}

func fileSize(path string) int {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return int(info.Size())
}

func generateMarkdown(outDir string, data ReportData) error {
	path := filepath.Join(outDir, "report.md")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, `# ReconX Report - %s

**Generated:** %s

## Summary

| Metric | Count |
|--------|-------|
| Subdomains Discovered | %d |
| Alive Hosts | %d |
| HTTP Services | %d |
| Open Ports | %d |
| Screenshots | %d |
| Vulnerabilities | %d |
`, data.Domain, data.Timestamp, data.Subdomains, data.Alive, data.HTTPAlive, data.Ports, data.Screenshots, data.Vulns)

	if data.Secrets > 0 {
		fmt.Fprintf(f, "\n⚠️ Secrets were found! See `%s`\n", data.SecretsFile)
	}
	if data.Vulns > 0 {
		fmt.Fprintf(f, "\n🔴 Vulnerabilities were found! See `%s`\n", data.VulnsFile)
	}

	fmt.Fprintf(f, "\n## Detailed Status\n\n")
	fmt.Fprintf(f, "- [x] Asset Discovery (%d subdomains)\n", data.Subdomains)
	fmt.Fprintf(f, "- [%s] DNS Resolution\n", check(data.HasDNS))
	fmt.Fprintf(f, "- [%s] Port Discovery (%d ports)\n", check(data.Ports > 0), data.Ports)
	fmt.Fprintf(f, "- [%s] HTTP Discovery (%d alive)\n", check(data.HTTPAlive > 0), data.HTTPAlive)
	fmt.Fprintf(f, "- [%s] Content Discovery\n", check(data.HasContent))
	fmt.Fprintf(f, "- [%s] Historical Collection\n", check(data.HasHistory))
	fmt.Fprintf(f, "- [%s] JS Intelligence\n", check(data.HasJS))
	fmt.Fprintf(f, "- [%s] Cloud Enumeration (%d bytes)\n", check(data.Cloud > 0), data.Cloud)
	fmt.Fprintf(f, "- [%s] OSINT (%d bytes)\n", check(data.OSINT > 0), data.OSINT)
	fmt.Fprintf(f, "- [%s] Visual Reconnaissance (%d screenshots)\n", check(data.Screenshots > 0), data.Screenshots)
	fmt.Fprintf(f, "- [%s] Vulnerability Discovery (%d findings)\n", check(data.Vulns > 0), data.Vulns)

	return nil
}

func check(ok bool) string {
	if ok {
		return "x"
	}
	return " "
}

func generateHTML(outDir string, data ReportData) error {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><title>ReconX Report - {{.Domain}}</title>
<style>
body{font-family:system-ui,sans-serif;max-width:800px;margin:2em auto;padding:0 1em}
table{border-collapse:collapse;width:100%}
td,th{border:1px solid #ccc;padding:8px;text-align:left}
th{background:#f5f5f5}
.ok{color:green}.warn{color:orange}.fail{color:red}
</style></head>
<body>
<h1>ReconX Report - {{.Domain}}</h1>
<p>Generated: {{.Timestamp}}</p>
<h2>Summary</h2>
<table>
<tr><th>Metric</th><th>Count</th></tr>
<tr><td>Subdomains</td><td>{{.Subdomains}}</td></tr>
<tr><td>Alive Hosts</td><td>{{.Alive}}</td></tr>
<tr><td>HTTP Services</td><td>{{.HTTPAlive}}</td></tr>
<tr><td>Open Ports</td><td>{{.Ports}}</td></tr>
<tr><td>Screenshots</td><td>{{.Screenshots}}</td></tr>
<tr><td>Vulnerabilities</td><td>{{.Vulns}}</td></tr>
</table>
</body></html>`

	t, err := template.New("report").Parse(tmpl)
	if err != nil {
		return err
	}

	path := filepath.Join(outDir, "report.html")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, data)
}

func generateJSON(outDir string, data ReportData) error {
	path := filepath.Join(outDir, "report.json")
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}
