package orchestrator

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// Executor wraps AgentRunner Core execution.
type Executor struct {
	AgentRunnerPath string // Path to agent-runner binary
	TaskStore       *TaskStore
}

// NewExecutor creates a new Executor.
func NewExecutor(agentRunnerPath string, ts *TaskStore) *Executor {
	return &Executor{
		AgentRunnerPath: agentRunnerPath,
		TaskStore:       ts,
	}
}

// ExecuteTask runs the agent-runner for a given task.
func (e *Executor) ExecuteTask(ctx context.Context, task *Task) (*Attempt, error) {
	// Create new attempt
	attempt := &Attempt{
		ID:        uuid.New().String(),
		TaskID:    task.ID,
		Status:    AttemptStatusRunning,
		StartedAt: time.Now(),
	}

	// Save attempt
	if err := e.TaskStore.SaveAttempt(attempt); err != nil {
		return nil, fmt.Errorf("failed to save attempt: %w", err)
	}

	// Update task status to RUNNING
	task.Status = TaskStatusRunning
	now := time.Now()
	task.StartedAt = &now
	if err := e.TaskStore.SaveTask(task); err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	// Generate task YAML for agent-runner
	taskYAML := e.generateTaskYAML(task)

	// Execute agent-runner
	cmd := exec.CommandContext(ctx, e.AgentRunnerPath)
	cmd.Dir = filepath.Dir(e.TaskStore.WorkspaceDir)

	// Pass task YAML via stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return e.handleExecutionError(attempt, task, err)
	}

	go func() {
		defer func() { _ = stdin.Close() }()
		_, _ = stdin.Write([]byte(taskYAML))
	}()

	output, err := cmd.CombinedOutput()
	finishedAt := time.Now()
	attempt.FinishedAt = &finishedAt

	if err != nil {
		attempt.Status = AttemptStatusFailed
		attempt.ErrorSummary = fmt.Sprintf("Execution failed: %s\nOutput: %s", err.Error(), string(output))
		task.Status = TaskStatusFailed
		task.DoneAt = &finishedAt
	} else {
		attempt.Status = AttemptStatusSucceeded
		task.Status = TaskStatusSucceeded
		task.DoneAt = &finishedAt
	}

	// Save updated attempt and task
	if err := e.TaskStore.SaveAttempt(attempt); err != nil {
		return attempt, fmt.Errorf("failed to update attempt: %w", err)
	}
	if err := e.TaskStore.SaveTask(task); err != nil {
		return attempt, fmt.Errorf("failed to update task: %w", err)
	}

	return attempt, nil
}

func (e *Executor) handleExecutionError(attempt *Attempt, task *Task, err error) (*Attempt, error) {
	now := time.Now()
	attempt.FinishedAt = &now
	attempt.Status = AttemptStatusFailed
	attempt.ErrorSummary = err.Error()

	task.Status = TaskStatusFailed
	task.DoneAt = &now

	_ = e.TaskStore.SaveAttempt(attempt)
	_ = e.TaskStore.SaveTask(task)

	return attempt, err
}

func (e *Executor) generateTaskYAML(task *Task) string {
	// Simple task YAML for agent-runner
	return fmt.Sprintf(`version: "1"
task:
  id: %s
  title: %s
  repo: "."
  prd:
    text: "Execute task: %s"
runner:
  max_loops: 5
  worker:
    cli: "codex"
`, task.ID, task.Title, task.Title)
}
