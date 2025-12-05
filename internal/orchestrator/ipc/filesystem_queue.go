package ipc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Job represents a unit of work in the queue.
type Job struct {
	ID      string `json:"id"`
	TaskID  string `json:"taskId"`
	PoolID  string `json:"poolId"`
	Payload any    `json:"payload"`
}

// FilesystemQueue handles file-based IPC queue operations.
type FilesystemQueue struct {
	WorkspaceDir string
}

// NewFilesystemQueue creates a new FilesystemQueue.
func NewFilesystemQueue(workspaceDir string) *FilesystemQueue {
	return &FilesystemQueue{WorkspaceDir: workspaceDir}
}

// GetQueueDir returns the directory for a specific pool's queue.
func (q *FilesystemQueue) GetQueueDir(poolID string) string {
	return filepath.Join(q.WorkspaceDir, "ipc", "queue", poolID)
}

// Enqueue adds a job to the queue.
func (q *FilesystemQueue) Enqueue(job *Job) error {
	dir := q.GetQueueDir(job.PoolID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create queue directory: %w", err)
	}

	path := filepath.Join(dir, job.ID+".json")
	data, err := json.MarshalIndent(job, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write job file: %w", err)
	}

	return nil
}

// Dequeue removes a job from the queue (simulated by listing and removing).
// In a real scenario, the worker would pick this up.
// This function is mainly for testing or for the orchestrator to check status.
func (q *FilesystemQueue) ListJobs(poolID string) ([]string, error) {
	dir := q.GetQueueDir(poolID)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}

	var jobIDs []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			jobIDs = append(jobIDs, entry.Name()[:len(entry.Name())-5])
		}
	}
	return jobIDs, nil
}
