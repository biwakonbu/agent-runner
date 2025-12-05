package orchestrator

import (
	"os"
	"testing"
	"time"
)

func TestTaskStore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "task_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewTaskStore(tmpDir)

	task := &Task{
		ID:        "task-1",
		Title:     "Test Task",
		Status:    TaskStatusPending,
		PoolID:    "pool-1",
		CreatedAt: time.Now(),
	}

	if err := store.SaveTask(task); err != nil {
		t.Fatalf("SaveTask failed: %v", err)
	}

	loadedTask, err := store.LoadTask("task-1")
	if err != nil {
		t.Fatalf("LoadTask failed: %v", err)
	}

	if loadedTask.Title != task.Title {
		t.Errorf("expected Title %s, got %s", task.Title, loadedTask.Title)
	}
	if loadedTask.Status != TaskStatusPending {
		t.Errorf("expected Status %s, got %s", TaskStatusPending, loadedTask.Status)
	}

	// Update task
	task.Status = TaskStatusRunning
	if err := store.SaveTask(task); err != nil {
		t.Fatalf("SaveTask update failed: %v", err)
	}

	loadedTask2, err := store.LoadTask("task-1")
	if err != nil {
		t.Fatalf("LoadTask failed: %v", err)
	}

	if loadedTask2.Status != TaskStatusRunning {
		t.Errorf("expected Status %s, got %s", TaskStatusRunning, loadedTask2.Status)
	}
}

func TestAttemptStore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "attempt_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewTaskStore(tmpDir)

	attempt := &Attempt{
		ID:        "attempt-1",
		TaskID:    "task-1",
		Status:    AttemptStatusRunning,
		StartedAt: time.Now(),
	}

	if err := store.SaveAttempt(attempt); err != nil {
		t.Fatalf("SaveAttempt failed: %v", err)
	}

	loadedAttempt, err := store.LoadAttempt("attempt-1")
	if err != nil {
		t.Fatalf("LoadAttempt failed: %v", err)
	}

	if loadedAttempt.Status != attempt.Status {
		t.Errorf("expected Status %s, got %s", attempt.Status, loadedAttempt.Status)
	}
}
