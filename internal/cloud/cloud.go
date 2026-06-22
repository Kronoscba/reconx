package cloud

import (
	"context"
	"os/exec"
	"path/filepath"
	"reconx/internal/engine"
	"reconx/internal/logging"
)

type CloudEnumerationTask struct{}

func (t *CloudEnumerationTask) Name() string {
	return "Cloud Enumeration"
}

func (t *CloudEnumerationTask) RequiredTools() []string {
	return []string{"cloud-enum"}
}

func (t *CloudEnumerationTask) Execute(ctx context.Context, session *engine.Session) error {
	cloudJsonFile := session.Config.GetPath("cloud.json")

	logging.Log.Info("Running Cloud Enumeration", "domain", session.Target)

	projectPath := "/media/gabi/Data/Proyectos/reconx/internal/cloud/cloud_enum"
	scriptPath := filepath.Join(projectPath, "cloud_enum.py")

	// Use absolute path to the script to avoid 'No such file' errors with uv run
	cmd := exec.CommandContext(ctx, "uv", "run", "--project", projectPath, scriptPath, "-k", session.Target, "-l", cloudJsonFile, "-f", "json")

	if err := cmd.Run(); err != nil {
		logging.Log.Warn("cloud-enum failed", "error", err)
	}

	logging.Log.Info("Cloud enumeration complete", "output", cloudJsonFile)
	return nil
}
