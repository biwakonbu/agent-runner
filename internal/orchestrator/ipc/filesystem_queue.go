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

// Dequeue claims the next available job from the queue.
// It moves the job file from the queue directory to a processing directory.
func (q *FilesystemQueue) Dequeue(poolID string) (*Job, error) {
	queueDir := q.GetQueueDir(poolID)
	entries, err := os.ReadDir(queueDir)
	if os.IsNotExist(err) {
		return nil, nil // Queue not created yet, empty
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read queue directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			// Found a job
			filename := entry.Name()
			srcPath := filepath.Join(queueDir, filename)

			// Create processing directory if not exists
			procDir := q.GetProcessingDir(poolID)
			if err := os.MkdirAll(procDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create processing directory: %w", err)
			}

			destPath := filepath.Join(procDir, filename)

			// Move file to processing (Claim)
			// Using Rename as atomic-ish operation on same filesystem
			if err := os.Rename(srcPath, destPath); err != nil {
				// Potential race condition or error, skip to next or return error
				// For now, return error to retry
				return nil, fmt.Errorf("failed to claim job (move): %w", err)
			}

			// Read and unmarshal
			data, err := os.ReadFile(destPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read claimed job file: %w", err)
			}

			var job Job
			if err := json.Unmarshal(data, &job); err != nil {
				return nil, fmt.Errorf("failed to unmarshal job: %w", err)
			}

			return &job, nil
		}
	}

	return nil, nil // No jobs found
}

// Complete removes a job from the processing directory, marking it as done.
func (q *FilesystemQueue) Complete(jobID, poolID string) error {
	procDir := q.GetProcessingDir(poolID)
	path := filepath.Join(procDir, jobID+".json")

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil // Already gone
		}
		return fmt.Errorf("failed to remove completed job file: %w", err)
	}

	return nil
}

// GetProcessingDir returns the directory for a specific pool's processing jobs.
func (q *FilesystemQueue) GetProcessingDir(poolID string) string {
	return filepath.Join(q.WorkspaceDir, "ipc", "processing", poolID)
}

// ListJobs returns all job IDs in the queue (pending).
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
