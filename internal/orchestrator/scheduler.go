package orchestrator

import (
	"fmt"
	"log/slog"

	"github.com/biwakonbu/agent-runner/internal/logging"
	"github.com/biwakonbu/agent-runner/internal/orchestrator/ipc"
)

// Scheduler manages task execution.
type Scheduler struct {
	TaskStore *TaskStore
	Queue     *ipc.FilesystemQueue
	logger    *slog.Logger
}

// NewScheduler creates a new Scheduler.
func NewScheduler(ts *TaskStore, q *ipc.FilesystemQueue) *Scheduler {
	return &Scheduler{
		TaskStore: ts,
		Queue:     q,
		logger:    logging.WithComponent(slog.Default(), "scheduler"),
	}
}

// ScheduleTask schedules a task for execution.
func (s *Scheduler) ScheduleTask(taskID string) error {
	task, err := s.TaskStore.LoadTask(taskID)
	if err != nil {
		return fmt.Errorf("failed to load task: %w", err)
	}

	if task.Status != TaskStatusPending && task.Status != TaskStatusFailed {
		return fmt.Errorf("task is not in a schedulable state: %s", task.Status)
	}

	// Update task status to READY
	task.Status = TaskStatusReady
	if err := s.TaskStore.SaveTask(task); err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	// Create a job for the queue
	job := &ipc.Job{
		ID:      fmt.Sprintf("job-%s-%d", task.ID, task.UpdatedAt.UnixNano()),
		TaskID:  task.ID,
		PoolID:  task.PoolID,
		Payload: map[string]string{"action": "run_task"},
	}

	if err := s.Queue.Enqueue(job); err != nil {
		s.logger.Error("failed to enqueue job",
			slog.String("job_id", job.ID),
			slog.String("task_id", task.ID),
			slog.Any("error", err),
		)
		return fmt.Errorf("failed to enqueue job: %w", err)
	}

	s.logger.Info("task scheduled",
		slog.String("task_id", task.ID),
		slog.String("job_id", job.ID),
		slog.String("pool_id", task.PoolID),
	)
	return nil
}
