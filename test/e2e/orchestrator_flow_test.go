package e2e_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/biwakonbu/agent-runner/internal/ide"
	"github.com/biwakonbu/agent-runner/internal/orchestrator"
	"github.com/biwakonbu/agent-runner/internal/orchestrator/ipc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrchestratorFlow(t *testing.T) {
	// 1. Setup specific test workspace
	tempDir, err := os.MkdirTemp("", "multiverse-e2e-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	t.Logf("Test Workspace: %s", tempDir)

	// 2. Build binaries (Orchestrator)
	// Assuming running from test/e2e
	wd, _ := os.Getwd()
	rootDir := filepath.Dir(filepath.Dir(wd)) // ../../

	orchBin := filepath.Join(tempDir, "multiverse-orchestrator")
	buildCmd := exec.Command("go", "build", "-o", orchBin, "./cmd/multiverse-orchestrator")
	buildCmd.Dir = rootDir
	out, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Failed to build orchestrator: %s", string(out))

	// Mock Agent Runner
	mockRunner := filepath.Join(rootDir, "test/e2e/mock_runner.sh")

	// 3. Initialize Backend Components (Simulating App logic)
	// We use the internal packages directly to ensure the flow works from data perspective
	wsStore := ide.NewWorkspaceStore(filepath.Join(tempDir, "workspaces"))

	// Create a workspace
	projectRoot := filepath.Join(tempDir, "project")
	ws := &ide.Workspace{
		ProjectRoot: projectRoot,
		Version:     "1.0",
		DisplayName: "E2E Test Project",
	}
	err = wsStore.SaveWorkspace(ws)
	require.NoError(t, err)

	wsID := wsStore.GetWorkspaceID(projectRoot)
	wsDir := wsStore.GetWorkspaceDir(wsID)
	taskStore := orchestrator.NewTaskStore(wsDir)
	queue := ipc.NewFilesystemQueue(wsDir)
	scheduler := orchestrator.NewScheduler(taskStore, queue)

	// 4. Create Task (Simulating IDE creates task)
	task := &orchestrator.Task{
		ID:     "e2e-task-1",
		Title:  "E2E Test Task",
		Status: orchestrator.TaskStatusPending,
		PoolID: "default",
	}
	err = taskStore.SaveTask(task)
	require.NoError(t, err)

	// 5. Schedule Task (Simulating IDE schedules task)
	err = scheduler.ScheduleTask(task.ID)
	require.NoError(t, err)

	// Verify task is READY
	loadedTask, err := taskStore.LoadTask(task.ID)
	require.NoError(t, err)
	assert.Equal(t, orchestrator.TaskStatusReady, loadedTask.Status)

	// 6. Start Orchestrator Process
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, orchBin,
		"--workspace", wsDir,
		"--agent-runner", mockRunner,
		"--pool", "default")

	// Capture output for debugging
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	require.NoError(t, err)

	// 7. Wait for completion (Polling)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(500 * time.Millisecond):
				t, err := taskStore.LoadTask(task.ID)
				if err == nil {
					if t.Status == orchestrator.TaskStatusSucceeded {
						done <- true
						return
					}
					if t.Status == orchestrator.TaskStatusFailed {
						// Fail fast
						return
					}
				}
			}
		}
	}()

	select {
	case <-done:
		t.Log("Task succeeded!")
	case <-ctx.Done():
		t.Fatal("Timeout waiting for task success")
	}

	// Cleanup process
	_ = cmd.Process.Kill()
}
