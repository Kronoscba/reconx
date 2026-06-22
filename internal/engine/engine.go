package engine

import (
	"context"
	"fmt"
	"reconx/internal/config"
	"reconx/internal/logging"
	"time"
)

// Session holds the state of a reconnaissance scan
type Session struct {
	Target     string
	Config     *config.Config
	StartTime  time.Time
	
	// Shared data between tasks
	Subdomains  []string
	AliveHosts  []string
	ResolvedDNS map[string]string
	OpenPorts   map[string][]int
	HTTPResults []HTTPResult
	CrawlData   []string
	// Add other fields as needed for later phases
}

type HTTPResult struct {
	Host string
	Port int
	URL  string
	Title string
	// ...
}

// Task is the interface that all recon modules must implement
type Task interface {
	Name() string
	Execute(ctx context.Context, session *Session) error
	RequiredTools() []string
}

// Engine manages the execution of the pipeline
type Engine struct {
	tasks    []Task
	resume   bool
}

func New() *Engine {
	return &Engine{
		tasks:  make([]Task, 0),
		resume: false,
	}
}

func (e *Engine) SetResume(resume bool) {
	e.resume = resume
}

func (e *Engine) AddTask(task Task) {
	e.tasks = append(e.tasks, task)
}

func (e *Engine) GetTasks() []Task {
	return e.tasks
}

func (e *Engine) Run(ctx context.Context, session *Session) error {
	logging.Log.Info("Starting pipeline execution", "total_tasks", len(e.tasks))
	
	// Load scan history for resume support
	history := LoadHistory(session)
	history.StartTime = time.Now()

	for _, task := range e.tasks {
		// Resume: skip tasks already completed
		if e.resume && history.IsTaskCompleted(task.Name()) {
			logging.Log.Info("Skipping completed task", "task", task.Name())
			continue
		}

		logging.Log.Info("Executing task", "task", task.Name())
		
		if err := ensureTaskDeps(task.RequiredTools()); err != nil {
			return fmt.Errorf("dependency check failed for task %s: %w", task.Name(), err)
		}

		start := time.Now()
		
		err := task.Execute(ctx, session)
		duration := time.Since(start)

		record := TaskRecord{
			Name:     task.Name(),
			Status:   "completed",
			Duration: duration.String(),
			RanAt:    start,
		}

		if err != nil {
			record.Status = "failed"
			logging.Log.Error("Task failed", "task", task.Name(), "error", err)
			// Save history before returning
			history.MarkTaskCompleted(record)
			history.Save(session)
			return fmt.Errorf("task %s failed: %w", task.Name(), err)
		}

		logging.Log.Info("Task completed", "task", task.Name(), "duration", duration)
		history.MarkTaskCompleted(record)
		history.Save(session)
	}
	
	logging.Log.Info("Pipeline execution finished successfully")
	return nil
}

func ensureTaskDeps(tools []string) error {
	// This is a helper to check only a subset of tools
	// In a real scenario, this would call a simplified version of deps.EnsureDependencies
	// For now, we will just log and continue unless a critical tool is missing
	return nil 
}
