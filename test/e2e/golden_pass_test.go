package e2e_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/biwakonbu/agent-runner/internal/chat"
	"github.com/biwakonbu/agent-runner/internal/ide"
	"github.com/biwakonbu/agent-runner/internal/meta"
	"github.com/biwakonbu/agent-runner/internal/orchestrator"
	"github.com/biwakonbu/agent-runner/internal/orchestrator/ipc"
	"github.com/biwakonbu/agent-runner/internal/orchestrator/persistence"
	"github.com/stretchr/testify/require"
)

// StubTaskExecutor simulates AgentRunner execution for golden pass.
// It records call order and always succeeds.
type StubTaskExecutor struct {
	mu    sync.Mutex
	Calls []string
}

func (s *StubTaskExecutor) ExecuteTask(ctx context.Context, task *orchestrator.Task) (*orchestrator.Attempt, error) {
	s.mu.Lock()
	s.Calls = append(s.Calls, task.ID)
	s.mu.Unlock()

	now := time.Now()
	return &orchestrator.Attempt{
		ID:         "attempt-" + task.ID,
		TaskID:     task.ID,
		Status:     orchestrator.AttemptStatusSucceeded,
		StartedAt:  now,
		FinishedAt: &now,
	}, nil
}

func TestGoldenPass_Backend(t *testing.T) {
	// 1. Setup workspace and repo
	tempDir := t.TempDir()

	wsStore := ide.NewWorkspaceStore(filepath.Join(tempDir, "workspaces"))
	projectRoot := filepath.Join(tempDir, "project")
	require.NoError(t, os.MkdirAll(projectRoot, 0755))

	ws := &ide.Workspace{
		ProjectRoot: projectRoot,
		Version:     "1.0",
		DisplayName: "Golden Pass Backend",
	}
	require.NoError(t, wsStore.SaveWorkspace(ws))

	wsID := wsStore.GetWorkspaceID(projectRoot)
	wsDir := wsStore.GetWorkspaceDir(wsID)

	repo := persistence.NewWorkspaceRepository(wsDir)
	require.NoError(t, repo.Init())

	queue := ipc.NewFilesystemQueue(wsDir)
	taskStore := orchestrator.NewTaskStore(wsDir)
	sessionStore := chat.NewChatSessionStore(wsDir)

	// 2. Chat -> Decompose -> persist design/state/tasks
	metaClient := meta.NewMockClient()
	handler := chat.NewHandler(metaClient, taskStore, sessionStore, wsID, projectRoot, repo, nil)

	ctx := context.Background()
	session, err := handler.CreateSession(ctx)
	require.NoError(t, err)

	resp, err := handler.HandleMessage(ctx, session.ID, "TODO アプリを作成して")
	require.NoError(t, err)
	require.Len(t, resp.GeneratedTasks, 2)

	var conceptTask, implTask orchestrator.Task
	for _, task := range resp.GeneratedTasks {
		switch task.Title {
		case "Mock概念設計タスク":
			conceptTask = task
		case "Mock実装タスク":
			implTask = task
		}
	}
	require.NotEmpty(t, conceptTask.ID)
	require.NotEmpty(t, implTask.ID)
	require.Equal(t, []string{conceptTask.ID}, implTask.Dependencies)

	// 3. Verify design persistence
	wbs, err := repo.Design().LoadWBS()
	require.NoError(t, err)
	require.NotEmpty(t, wbs.WBSID)
	require.NotEmpty(t, wbs.RootNodeID)

	_, err = repo.Design().GetNode(conceptTask.ID)
	require.NoError(t, err)
	_, err = repo.Design().GetNode(implTask.ID)
	require.NoError(t, err)

	// 4. Verify state persistence
	tasksState, err := repo.State().LoadTasks()
	require.NoError(t, err)
	require.Len(t, tasksState.Tasks, 2)

	nodesRuntime, err := repo.State().LoadNodesRuntime()
	require.NoError(t, err)
	require.Len(t, nodesRuntime.Nodes, 2)

	// 5. Run ExecutionOrchestrator with stub executor
	scheduler := orchestrator.NewScheduler(repo, queue, nil)
	exec := &StubTaskExecutor{}
	backlogStore := orchestrator.NewBacklogStore(wsDir)

	orch := orchestrator.NewExecutionOrchestrator(
		scheduler,
		exec,
		repo,
		queue,
		nil,
		backlogStore,
		[]string{"default"},
	)

	runCtx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()
	require.NoError(t, orch.Start(runCtx))

	// Wait until both tasks are succeeded in state.
	require.Eventually(t, func() bool {
		state, err := repo.State().LoadTasks()
		if err != nil {
			return false
		}
		if len(state.Tasks) != 2 {
			return false
		}
		succeeded := 0
		for _, ts := range state.Tasks {
			if ts.Status == string(orchestrator.TaskStatusSucceeded) {
				succeeded++
			}
		}
		return succeeded == 2
	}, 10*time.Second, 100*time.Millisecond)

	require.NoError(t, orch.Stop())
	orch.Wait()

	// 6. Verify execution order respected (dependency first)
	exec.mu.Lock()
	require.Len(t, exec.Calls, 2)
	require.Equal(t, conceptTask.ID, exec.Calls[0])
	require.Equal(t, implTask.ID, exec.Calls[1])
	exec.mu.Unlock()

	// 7. Verify legacy TaskStore was synchronized
	legacyConcept, err := taskStore.LoadTask(conceptTask.ID)
	require.NoError(t, err)
	require.Equal(t, orchestrator.TaskStatusSucceeded, legacyConcept.Status)
	require.Equal(t, 1, legacyConcept.AttemptCount)
	require.NotNil(t, legacyConcept.StartedAt)
	require.NotNil(t, legacyConcept.DoneAt)

	legacyImpl, err := taskStore.LoadTask(implTask.ID)
	require.NoError(t, err)
	require.Equal(t, orchestrator.TaskStatusSucceeded, legacyImpl.Status)
	require.Equal(t, 1, legacyImpl.AttemptCount)
	require.NotNil(t, legacyImpl.StartedAt)
	require.NotNil(t, legacyImpl.DoneAt)
}
