package historical

import (
	"context"
	"fmt"
	"os/exec"
	"reconx/internal/engine"
	"reconx/internal/logging"
)

type HistoricalTask struct{}

func (t *HistoricalTask) Name() string {
	return "Historical Collection"
}

func (t *HistoricalTask) RequiredTools() []string {
	return []string{"gau", "waymore"}
}

func (t *HistoricalTask) Execute(ctx context.Context, session *engine.Session) error {
	histFile := session.Config.GetPath("historical.txt")
	
	logging.Log.Info("Running Historical Collection", "domain", session.Target)
	
	// GAU
	shellGau := fmt.Sprintf("gau --subs %s --o %s", session.Target, histFile)
	cmdGau := exec.CommandContext(ctx, "sh", "-c", shellGau)
	if err := cmdGau.Run(); err != nil {
		logging.Log.Warn("GAU failed", "error", err)
	}

	// WAYMORE
	shellWay := fmt.Sprintf("waymore -i %s -o historical_waymore", session.Target)
	cmdWay := exec.CommandContext(ctx, "sh", "-c", shellWay)
	if err := cmdWay.Run(); err != nil {
		logging.Log.Debug("Waymore failed or not installed", "error", err)
	}
	
	logging.Log.Info("Historical collection complete", "output", histFile)
	return nil
}
