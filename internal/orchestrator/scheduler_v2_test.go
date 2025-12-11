package orchestrator

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"sync"

	"github.com/biwakonbu/agent-runner/internal/orchestrator/persistence"
)

// MockExecutorV2 for testing
type MockExecutorV2 struct {
	mu            sync.Mutex
	executedTasks []string
	wg            sync.WaitGroup
}

func (m *MockExecutorV2) Execute(ctx context.Context, task persistence.TaskState) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.executedTasks = append(m.executedTasks, task.TaskID)
	m.wg.Done()
	return nil
}

func TestSchedulerV2_CheckAndSchedule(t *testing.T) {
	// Setup Temp Repo
	tmpDir, err := os.MkdirTemp("", "scheduler-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	repo := persistence.NewWorkspaceRepository(tmpDir)
	if err := repo.Init(); err != nil {
		t.Fatal(err)
	}

	// Setup Initial State
	// 1. WBS Node
	repo.Design().SaveNode(&persistence.NodeDesign{
		NodeID:       "node-1",
		Dependencies: []string{}, // No deps
	})

	// 2. Task Pending
	repo.State().SaveTasks(&persistence.TasksState{
		Tasks: []persistence.TaskState{
			{
				TaskID: "task-1",
				NodeID: "node-1",
				Status: "pending",
			},
		},
	})

	// 3. Agent Available
	repo.State().SaveAgents(&persistence.AgentsState{
		Agents: []persistence.AgentState{
			{
				AgentID:      "agent-1",
				MaxParallel:  1,
				RunningTasks: []string{},
			},
		},
	})

	// 4. Runtime (Empty is fine for no deps)
	repo.State().SaveNodesRuntime(&persistence.NodesRuntime{
		Nodes: []persistence.NodeRuntime{},
	})

	// Setup Scheduler
	mockExec := &MockExecutorV2{}
	logger := slog.Default()
	scheduler := NewSchedulerV2(repo, mockExec, logger)

	// Expect 1 task execution
	mockExec.wg.Add(1)

	// Execute Schedule Step
	if err := scheduler.CheckAndSchedule(context.Background()); err != nil {
		t.Fatalf("CheckAndSchedule failed: %v", err)
	}

	// Wait for go routine (Execute is called in go routine)
	mockExec.wg.Wait()

	// Verify Dispatch
	mockExec.mu.Lock() // Start lock for verification
	defer mockExec.mu.Unlock()

	if len(mockExec.executedTasks) != 1 {
		t.Errorf("Expected 1 executed task, got %d", len(mockExec.executedTasks))
	} else if mockExec.executedTasks[0] != "task-1" {
		t.Errorf("Expected task-1, got %s", mockExec.executedTasks[0])
	}

	// Verify State Updates
	tasks, _ := repo.State().LoadTasks()
	if tasks.Tasks[0].Status != "running" {
		t.Errorf("Expected task status running, got %s", tasks.Tasks[0].Status)
	}
	if tasks.Tasks[0].AssignedAgent != "agent-1" {
		t.Errorf("Expected assigned agent agent-1, got %s", tasks.Tasks[0].AssignedAgent)
	}

	agents, _ := repo.State().LoadAgents()
	if len(agents.Agents[0].RunningTasks) != 1 {
		t.Errorf("Expected agent to have 1 running task")
	}

	// Verify History
	actions, _ := repo.History().ListActions(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour))
	if len(actions) == 0 {
		t.Error("Expected action history to be non-empty")
	}
	foundStart := false
	for _, a := range actions {
		if a.Kind == "task.started" {
			foundStart = true
			if a.Payload["task_id"] != "task-1" {
				t.Errorf("Action payload mismatch")
			}
		}
	}
	if !foundStart {
		t.Error("Did not find task.started action")
	}
}
