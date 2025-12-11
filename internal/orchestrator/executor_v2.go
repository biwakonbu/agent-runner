package orchestrator

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	"github.com/biwakonbu/agent-runner/internal/orchestrator/persistence"
	"github.com/google/uuid"
)

// ExecutorV2 defines the interface for executing tasks with new models
type ExecutorV2 interface {
	Execute(ctx context.Context, task persistence.TaskState) error
}

// executorV2Impl implements ExecutorV2
type executorV2Impl struct {
	AgentRunnerPath string
	ProjectRoot     string
	Repo            persistence.WorkspaceRepository
	Logger          *slog.Logger
}

func NewExecutorV2(agentRunnerPath, projectRoot string, repo persistence.WorkspaceRepository, logger *slog.Logger) ExecutorV2 {
	return &executorV2Impl{
		AgentRunnerPath: agentRunnerPath,
		ProjectRoot:     projectRoot,
		Repo:            repo,
		Logger:          logger,
	}
}

func (e *executorV2Impl) Execute(ctx context.Context, task persistence.TaskState) error {
	e.Logger.Info("ExecutorV2: starting task execution", "task_id", task.TaskID)

	taskYAML := e.generateTaskYAML(task)

	// Create Attempt Action
	attemptID := uuid.New().String()
	actionStart := &persistence.Action{
		ID:          uuid.New().String(),
		At:          time.Now(),
		Kind:        "task.attempt_started",
		WorkspaceID: "TODO-WS-ID", // propagate this
		Payload: map[string]interface{}{
			"task_id":    task.TaskID,
			"attempt_id": attemptID,
		},
	}
	_ = e.Repo.History().AppendAction(actionStart)

	cmd := exec.CommandContext(ctx, e.AgentRunnerPath)
	cmd.Dir = e.ProjectRoot
	cmd.Stdin = strings.NewReader(taskYAML)

	// Capture output
	out, err := cmd.CombinedOutput()

	finishedAt := time.Now()

	// Determine success/failure
	success := err == nil

	// Create Result Action
	kind := "task.succeeded"
	if !success {
		kind = "task.failed"
	}

	actionResult := &persistence.Action{
		ID:          uuid.New().String(),
		At:          finishedAt,
		Kind:        kind,
		WorkspaceID: "TODO-WS-ID",
		Payload: map[string]interface{}{
			"task_id":    task.TaskID,
			"attempt_id": attemptID,
			"output":     string(out), // Careful with size
			"error":      fmt.Sprintf("%v", err),
		},
	}
	if err := e.Repo.History().AppendAction(actionResult); err != nil {
		e.Logger.Error("failed to append result action", "err", err)
	}

	// Update Task State in Repository (Status, Outputs)
	// Currently Scheduler logic refreshes from memory or relies on events?
	// The PRD says Executor updates result.
	// We need to Read-Update-Write TaskState here or rely on Scheduler to pick up events?
	// "Executor ... results based on state/tasks.json ... auto update".

	// Let's update the task state directly here for MVP simplicity to "succeeded" or "failed"
	// Race condition warning: Scheduler might be updating "running" status logic.
	// But Scheduler is single-threaded logic mostly.

	currentTasks, repoErr := e.Repo.State().LoadTasks()
	if repoErr == nil {
		for i, t := range currentTasks.Tasks {
			if t.TaskID == task.TaskID {
				if success {
					currentTasks.Tasks[i].Status = "succeeded"
				} else {
					currentTasks.Tasks[i].Status = "failed"
				}
				currentTasks.Tasks[i].UpdatedAt = finishedAt
				_ = e.Repo.State().SaveTasks(currentTasks)
				break
			}
		}
	}

	return err
}

func (e *executorV2Impl) generateTaskYAML(task persistence.TaskState) string {
	// Extract inputs
	goal, _ := task.Inputs["goal"].(string)

	// Construct prompt
	promptText := fmt.Sprintf("Execute task: %s\n\nGoal:\n%s", task.TaskID, goal)

	if constraints, ok := task.Inputs["constraints"].([]interface{}); ok {
		promptText += "\n\nConstraints:"
		for _, c := range constraints {
			promptText += fmt.Sprintf("\n- %v", c)
		}
	}

	promptTextIndented := ""
	for _, line := range strings.Split(promptText, "\n") {
		promptTextIndented += fmt.Sprintf("      %s\n", line)
	}

	return fmt.Sprintf(`version: "1"
task:
  id: %s
  title: "Task %s"
  repo: "."
  prd:
    text: |
%srunner:
  max_loops: 5
  worker:
    cli: "codex"
`, task.TaskID, task.TaskID, promptTextIndented)
}
