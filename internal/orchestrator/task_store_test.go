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

func TestListAttemptsByTaskID(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "list_attempts_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewTaskStore(tmpDir)

	// 空のディレクトリの場合
	attempts, err := store.ListAttemptsByTaskID("task-1")
	if err != nil {
		t.Fatalf("ListAttemptsByTaskID failed on empty dir: %v", err)
	}
	if len(attempts) != 0 {
		t.Errorf("expected 0 attempts, got %d", len(attempts))
	}

	// 複数のAttemptを保存
	attempt1 := &Attempt{
		ID:        "attempt-1",
		TaskID:    "task-1",
		Status:    AttemptStatusSucceeded,
		StartedAt: time.Now(),
	}
	attempt2 := &Attempt{
		ID:        "attempt-2",
		TaskID:    "task-1",
		Status:    AttemptStatusFailed,
		StartedAt: time.Now(),
	}
	attempt3 := &Attempt{
		ID:        "attempt-3",
		TaskID:    "task-2", // 別のタスク
		Status:    AttemptStatusRunning,
		StartedAt: time.Now(),
	}

	if err := store.SaveAttempt(attempt1); err != nil {
		t.Fatalf("SaveAttempt failed: %v", err)
	}
	if err := store.SaveAttempt(attempt2); err != nil {
		t.Fatalf("SaveAttempt failed: %v", err)
	}
	if err := store.SaveAttempt(attempt3); err != nil {
		t.Fatalf("SaveAttempt failed: %v", err)
	}

	// task-1のAttemptsを取得
	attempts, err = store.ListAttemptsByTaskID("task-1")
	if err != nil {
		t.Fatalf("ListAttemptsByTaskID failed: %v", err)
	}
	if len(attempts) != 2 {
		t.Errorf("expected 2 attempts for task-1, got %d", len(attempts))
	}

	// task-2のAttemptsを取得
	attempts, err = store.ListAttemptsByTaskID("task-2")
	if err != nil {
		t.Fatalf("ListAttemptsByTaskID failed: %v", err)
	}
	if len(attempts) != 1 {
		t.Errorf("expected 1 attempt for task-2, got %d", len(attempts))
	}
}

func TestGetPoolSummaries(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pool_summary_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewTaskStore(tmpDir)

	// 空のディレクトリの場合
	summaries, err := store.GetPoolSummaries()
	if err != nil {
		t.Fatalf("GetPoolSummaries failed on empty dir: %v", err)
	}
	if len(summaries) != 0 {
		t.Errorf("expected 0 summaries, got %d", len(summaries))
	}

	// 複数のタスクを保存
	tasks := []*Task{
		{ID: "task-1", Title: "Task 1", Status: TaskStatusRunning, PoolID: "codegen", CreatedAt: time.Now()},
		{ID: "task-2", Title: "Task 2", Status: TaskStatusRunning, PoolID: "codegen", CreatedAt: time.Now()},
		{ID: "task-3", Title: "Task 3", Status: TaskStatusPending, PoolID: "codegen", CreatedAt: time.Now()},
		{ID: "task-4", Title: "Task 4", Status: TaskStatusFailed, PoolID: "codegen", CreatedAt: time.Now()},
		{ID: "task-5", Title: "Task 5", Status: TaskStatusRunning, PoolID: "test", CreatedAt: time.Now()},
		{ID: "task-6", Title: "Task 6", Status: TaskStatusSucceeded, PoolID: "test", CreatedAt: time.Now()},
	}

	for _, task := range tasks {
		if err := store.SaveTask(task); err != nil {
			t.Fatalf("SaveTask failed: %v", err)
		}
	}

	summaries, err = store.GetPoolSummaries()
	if err != nil {
		t.Fatalf("GetPoolSummaries failed: %v", err)
	}
	if len(summaries) != 2 {
		t.Errorf("expected 2 pools, got %d", len(summaries))
	}

	// サマリを検証
	poolMap := make(map[string]PoolSummary)
	for _, s := range summaries {
		poolMap[s.PoolID] = s
	}

	codegen, ok := poolMap["codegen"]
	if !ok {
		t.Fatal("codegen pool not found")
	}
	if codegen.Running != 2 {
		t.Errorf("expected codegen running=2, got %d", codegen.Running)
	}
	if codegen.Queued != 1 {
		t.Errorf("expected codegen queued=1, got %d", codegen.Queued)
	}
	if codegen.Failed != 1 {
		t.Errorf("expected codegen failed=1, got %d", codegen.Failed)
	}
	if codegen.Total != 4 {
		t.Errorf("expected codegen total=4, got %d", codegen.Total)
	}

	testPool, ok := poolMap["test"]
	if !ok {
		t.Fatal("test pool not found")
	}
	if testPool.Running != 1 {
		t.Errorf("expected test running=1, got %d", testPool.Running)
	}
	if testPool.Total != 2 {
		t.Errorf("expected test total=2, got %d", testPool.Total)
	}
}
