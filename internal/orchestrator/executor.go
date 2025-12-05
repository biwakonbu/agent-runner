package orchestrator

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/biwakonbu/agent-runner/internal/logging"
	"github.com/google/uuid"
)

// Executor wraps AgentRunner Core execution.
type Executor struct {
	AgentRunnerPath string // Path to agent-runner binary
	TaskStore       *TaskStore
	logger          *slog.Logger
}

// NewExecutor creates a new Executor.
func NewExecutor(agentRunnerPath string, ts *TaskStore) *Executor {
	return &Executor{
		AgentRunnerPath: agentRunnerPath,
		TaskStore:       ts,
		logger:          logging.WithComponent(slog.Default(), "orchestrator-executor"),
	}
}

// SetLogger sets a custom logger for the executor
func (e *Executor) SetLogger(logger *slog.Logger) {
	e.logger = logging.WithComponent(logger, "orchestrator-executor")
}

// ExecuteTask runs the agent-runner for a given task.
func (e *Executor) ExecuteTask(ctx context.Context, task *Task) (*Attempt, error) {
	logger := logging.WithTraceID(e.logger, ctx)
	start := time.Now()

	// Create new attempt
	attempt := &Attempt{
		ID:        uuid.New().String(),
		TaskID:    task.ID,
		Status:    AttemptStatusRunning,
		StartedAt: time.Now(),
	}

	logger.Info("starting task execution",
		slog.String("task_id", task.ID),
		slog.String("task_title", task.Title),
		slog.String("attempt_id", attempt.ID),
	)

	// Save attempt
	if err := e.TaskStore.SaveAttempt(attempt); err != nil {
		logger.Error("failed to save attempt", slog.Any("error", err))
		return nil, fmt.Errorf("failed to save attempt: %w", err)
	}

	// Update task status to RUNNING
	task.Status = TaskStatusRunning
	now := time.Now()
	task.StartedAt = &now
	if err := e.TaskStore.SaveTask(task); err != nil {
		logger.Error("failed to update task status", slog.Any("error", err))
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}
	logger.Info("task status updated to RUNNING")

	// Generate task YAML for agent-runner
	taskYAML := e.generateTaskYAML(task)
	logger.Debug("generated task YAML", slog.Int("yaml_length", len(taskYAML)))

	// Execute agent-runner
	logger.Info("executing agent-runner", slog.String("binary_path", e.AgentRunnerPath))
	cmd := exec.CommandContext(ctx, e.AgentRunnerPath)
	cmd.Dir = filepath.Dir(e.TaskStore.WorkspaceDir)

	// Pass task YAML via stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Error("failed to create stdin pipe", slog.Any("error", err))
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
		logger.Error("agent-runner execution failed",
			slog.Any("error", err),
			slog.Int("output_length", len(output)),
			logging.LogDuration(start),
		)
		logger.Debug("agent-runner output", slog.String("output", string(output)))
	} else {
		attempt.Status = AttemptStatusSucceeded
		task.Status = TaskStatusSucceeded
		task.DoneAt = &finishedAt
		logger.Info("agent-runner execution succeeded",
			slog.Int("output_length", len(output)),
			logging.LogDuration(start),
		)
		logger.Debug("agent-runner output", slog.String("output", string(output)))
	}

	// Save updated attempt and task
	if err := e.TaskStore.SaveAttempt(attempt); err != nil {
		logger.Error("failed to update attempt", slog.Any("error", err))
		return attempt, fmt.Errorf("failed to update attempt: %w", err)
	}
	if err := e.TaskStore.SaveTask(task); err != nil {
		logger.Error("failed to update task", slog.Any("error", err))
		return attempt, fmt.Errorf("failed to update task: %w", err)
	}

	logger.Info("task execution completed",
		slog.String("final_status", string(attempt.Status)),
		logging.LogDuration(start),
	)
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
