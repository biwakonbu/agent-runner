package ipc

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFilesystemQueue(t *testing.T) {
	queue := NewFilesystemQueue("/test/workspace")

	if queue.WorkspaceDir != "/test/workspace" {
		t.Errorf("expected WorkspaceDir /test/workspace, got %s", queue.WorkspaceDir)
	}
}

func TestGetQueueDir(t *testing.T) {
	queue := NewFilesystemQueue("/test/workspace")

	dir := queue.GetQueueDir("codegen")
	expected := filepath.Join("/test/workspace", "ipc", "queue", "codegen")

	if dir != expected {
		t.Errorf("expected %s, got %s", expected, dir)
	}
}

func TestGetProcessingDir(t *testing.T) {
	queue := NewFilesystemQueue("/test/workspace")

	dir := queue.GetProcessingDir("codegen")
	expected := filepath.Join("/test/workspace", "ipc", "processing", "codegen")

	if dir != expected {
		t.Errorf("expected %s, got %s", expected, dir)
	}
}

func TestEnqueue(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ipc_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	queue := NewFilesystemQueue(tmpDir)

	job := &Job{
		ID:      "job-1",
		TaskID:  "task-1",
		PoolID:  "codegen",
		Payload: map[string]string{"key": "value"},
	}

	if err := queue.Enqueue(job); err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	// ファイルが作成されたことを確認
	expectedPath := filepath.Join(tmpDir, "ipc", "queue", "codegen", "job-1.json")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected job file to exist at %s", expectedPath)
	}

	// ファイル内容を確認
	data, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("failed to read job file: %v", err)
	}

	var loadedJob Job
	if err := json.Unmarshal(data, &loadedJob); err != nil {
		t.Fatalf("failed to unmarshal job: %v", err)
	}

	if loadedJob.ID != job.ID {
		t.Errorf("expected ID %s, got %s", job.ID, loadedJob.ID)
	}
	if loadedJob.TaskID != job.TaskID {
		t.Errorf("expected TaskID %s, got %s", job.TaskID, loadedJob.TaskID)
	}
	if loadedJob.PoolID != job.PoolID {
		t.Errorf("expected PoolID %s, got %s", job.PoolID, loadedJob.PoolID)
	}
}

func TestListJobs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ipc_list_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	queue := NewFilesystemQueue(tmpDir)

	// 空のキューの場合
	jobs, err := queue.ListJobs("codegen")
	if err != nil {
		t.Fatalf("ListJobs failed on empty queue: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(jobs))
	}

	// 複数の Job を追加
	for i := 1; i <= 3; i++ {
		job := &Job{
			ID:     "job-" + string(rune('0'+i)),
			TaskID: "task-1",
			PoolID: "codegen",
		}
		if err := queue.Enqueue(job); err != nil {
			t.Fatalf("Enqueue failed: %v", err)
		}
	}

	jobs, err = queue.ListJobs("codegen")
	if err != nil {
		t.Fatalf("ListJobs failed: %v", err)
	}
	if len(jobs) != 3 {
		t.Errorf("expected 3 jobs, got %d", len(jobs))
	}

	// 別の Pool は空のまま
	otherJobs, err := queue.ListJobs("test")
	if err != nil {
		t.Fatalf("ListJobs failed for test pool: %v", err)
	}
	if len(otherJobs) != 0 {
		t.Errorf("expected 0 jobs in test pool, got %d", len(otherJobs))
	}
}

func TestDequeue(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ipc_dequeue_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	queue := NewFilesystemQueue(tmpDir)

	// 空のキューからDequeue
	job, err := queue.Dequeue("codegen")
	if err != nil {
		t.Fatalf("Dequeue failed on empty queue: %v", err)
	}
	if job != nil {
		t.Errorf("expected nil job from empty queue, got %+v", job)
	}

	// Job を追加
	originalJob := &Job{
		ID:      "job-1",
		TaskID:  "task-1",
		PoolID:  "codegen",
		Payload: map[string]string{"action": "build"},
	}
	if err := queue.Enqueue(originalJob); err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	// Dequeue
	dequeuedJob, err := queue.Dequeue("codegen")
	if err != nil {
		t.Fatalf("Dequeue failed: %v", err)
	}
	if dequeuedJob == nil {
		t.Fatal("expected non-nil job from Dequeue")
	}

	if dequeuedJob.ID != originalJob.ID {
		t.Errorf("expected ID %s, got %s", originalJob.ID, dequeuedJob.ID)
	}

	// Queue からは削除されていることを確認
	jobs, _ := queue.ListJobs("codegen")
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs in queue after dequeue, got %d", len(jobs))
	}

	// Processing ディレクトリに移動していることを確認
	procPath := filepath.Join(tmpDir, "ipc", "processing", "codegen", "job-1.json")
	if _, err := os.Stat(procPath); os.IsNotExist(err) {
		t.Errorf("expected job file in processing directory at %s", procPath)
	}
}

func TestComplete(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ipc_complete_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	queue := NewFilesystemQueue(tmpDir)

	// Job を追加して Dequeue
	job := &Job{
		ID:     "job-1",
		TaskID: "task-1",
		PoolID: "codegen",
	}
	if err := queue.Enqueue(job); err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	_, err = queue.Dequeue("codegen")
	if err != nil {
		t.Fatalf("Dequeue failed: %v", err)
	}

	// Complete
	if err := queue.Complete("job-1", "codegen"); err != nil {
		t.Fatalf("Complete failed: %v", err)
	}

	// Processing ディレクトリから削除されていることを確認
	procPath := filepath.Join(tmpDir, "ipc", "processing", "codegen", "job-1.json")
	if _, err := os.Stat(procPath); !os.IsNotExist(err) {
		t.Errorf("expected job file to be removed from processing directory")
	}
}

func TestCompleteNonExistent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ipc_complete_nonexistent_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	queue := NewFilesystemQueue(tmpDir)

	// 存在しない Job を Complete してもエラーにならない
	err = queue.Complete("nonexistent-job", "codegen")
	if err != nil {
		t.Errorf("Complete should not fail for non-existent job: %v", err)
	}
}

func TestJobStruct(t *testing.T) {
	// Job 構造体の基本的な検証
	job := Job{
		ID:      "test-id",
		TaskID:  "test-task",
		PoolID:  "test-pool",
		Payload: map[string]interface{}{"key": "value"},
	}

	if job.ID != "test-id" {
		t.Errorf("expected ID test-id, got %s", job.ID)
	}
	if job.TaskID != "test-task" {
		t.Errorf("expected TaskID test-task, got %s", job.TaskID)
	}
	if job.PoolID != "test-pool" {
		t.Errorf("expected PoolID test-pool, got %s", job.PoolID)
	}
}

func TestEnqueueMultiplePools(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ipc_multi_pool_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	queue := NewFilesystemQueue(tmpDir)

	// 複数の Pool に Job を追加
	pools := []string{"codegen", "test", "default"}
	for _, poolID := range pools {
		job := &Job{
			ID:     "job-" + poolID,
			TaskID: "task-1",
			PoolID: poolID,
		}
		if err := queue.Enqueue(job); err != nil {
			t.Fatalf("Enqueue failed for pool %s: %v", poolID, err)
		}
	}

	// 各 Pool に 1 つずつ Job があることを確認
	for _, poolID := range pools {
		jobs, err := queue.ListJobs(poolID)
		if err != nil {
			t.Fatalf("ListJobs failed for pool %s: %v", poolID, err)
		}
		if len(jobs) != 1 {
			t.Errorf("expected 1 job in pool %s, got %d", poolID, len(jobs))
		}
	}
}
