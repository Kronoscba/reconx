package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type TaskRecord struct {
	Name     string    `json:"name"`
	Status   string    `json:"status"` // completed, skipped, failed
	Duration string    `json:"duration"`
	RanAt    time.Time `json:"ran_at"`
	Outputs  []string  `json:"outputs"`
}

type ScanHistory struct {
	Domain    string       `json:"domain"`
	StartTime time.Time    `json:"start_time"`
	Tasks     []TaskRecord `json:"tasks"`
}

func LoadHistory(session *Session) *ScanHistory {
	path := filepath.Join(session.Config.OutputDir, ".scan_history.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return &ScanHistory{Domain: session.Target, Tasks: []TaskRecord{}}
	}

	var h ScanHistory
	if err := json.Unmarshal(data, &h); err != nil {
		return &ScanHistory{Domain: session.Target, Tasks: []TaskRecord{}}
	}
	return &h
}

func (h *ScanHistory) Save(session *Session) error {
	path := filepath.Join(session.Config.OutputDir, ".scan_history.json")
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (h *ScanHistory) IsTaskCompleted(name string) bool {
	for _, t := range h.Tasks {
		if t.Name == name && t.Status == "completed" {
			return true
		}
	}
	return false
}

func (h *ScanHistory) MarkTaskCompleted(record TaskRecord) {
	for i, t := range h.Tasks {
		if t.Name == record.Name {
			h.Tasks[i] = record
			return
		}
	}
	h.Tasks = append(h.Tasks, record)
}
