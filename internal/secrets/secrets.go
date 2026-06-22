package secrets

import (
	"context"
	"os"
	"os/exec"
	"reconx/internal/engine"
	"reconx/internal/logging"
)

type SecretsTask struct{}

func (t *SecretsTask) Name() string {
	return "Secret Discovery"
}

func (t *SecretsTask) RequiredTools() []string {
	return []string{"trufflehog"}
}

func (t *SecretsTask) Execute(ctx context.Context, session *engine.Session) error {
	histFile := session.Config.GetPath("historical.txt")
	secretsJsonFile := session.Config.GetPath("secrets.json")

	logging.Log.Info("Running Secret Discovery", "input", histFile)

	cmd := exec.CommandContext(ctx, "trufflehog", "filesystem", "--no-update", "-j", histFile)

	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 183 {
			// Secrets found
		} else {
			return err
		}
	}

	err = os.WriteFile(secretsJsonFile, out, 0644)
	if err != nil {
		return err
	}

	logging.Log.Info("Secret discovery complete", "output", secretsJsonFile)
	return nil
}
