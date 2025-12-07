package orchestrator

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	"github.com/biwakonbu/agent-runner/internal/logging"
	"github.com/google/uuid"
)

// TaskExecutor defines the interface for executing tasks
type TaskExecutor interface {
	ExecuteTask(ctx context.Context, task *Task) (*Attempt, error)
}

// Executor wraps AgentRunner Core execution.
type Executor struct {
	AgentRunnerPath string // Path to agent-runner binary
	ProjectRoot     string // Root directory of the project
	TaskStore       *TaskStore
	logger          *slog.Logger
}

// NewExecutor creates a new Executor.
func NewExecutor(agentRunnerPath string, projectRoot string, ts *TaskStore) *Executor {
	return &Executor{
		AgentRunnerPath: agentRunnerPath,
		ProjectRoot:     projectRoot,
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
	cmd.Dir = e.ProjectRoot

	// Pass task YAML via stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Error("failed to create stdin pipe", slog.Any("error", err))
		return e.handleExecutionError(attempt, task, err)
	}

	go func() {
		defer func() { _ = stdin.Close() }()
		select {
		case <-ctx.Done():
			return // Context canceled, close stdin (via defer) and exit
		default:
			// Write task YAML
			_, _ = stdin.Write([]byte(taskYAML))
		}
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
	return attempt, err
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
	// Construct the prompt text with Description and AcceptanceCriteria
	promptText := fmt.Sprintf("Execute task: %s", task.Title)
	if task.Description != "" {
		promptText += fmt.Sprintf("\n\nDescription:\n%s", task.Description)
	}
	if len(task.AcceptanceCriteria) > 0 {
		promptText += "\n\nAcceptance Criteria:"
		for _, ac := range task.AcceptanceCriteria {
			promptText += fmt.Sprintf("\n- %s", ac)
		}
	}

	// Simple task YAML for agent-runner
	// Using literal style Block Scalar (|) for prd.text to handle multi-line strings safely
	// Indentation must be correct (4 spaces for the text content)
	promptTextIndented := ""
	for _, line := range strings.Split(promptText, "\n") {
		promptTextIndented += fmt.Sprintf("    %s\n", line)
	}

	return fmt.Sprintf(`version: "1"
task:
  id: %s
  title: %s
  repo: "."
  prd:
    text: |
%srunner:
  max_loops: 5
  worker:
    cli: "codex"
`, task.ID, task.Title, promptTextIndented)
}
