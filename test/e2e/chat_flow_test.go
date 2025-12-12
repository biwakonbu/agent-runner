package e2e_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/biwakonbu/agent-runner/internal/chat"
	"github.com/biwakonbu/agent-runner/internal/ide"
	"github.com/biwakonbu/agent-runner/internal/meta"
	"github.com/biwakonbu/agent-runner/internal/orchestrator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockEventEmitter captures events for verification
type MockEventEmitter struct {
	Events []struct {
		Name string
		Data any
	}
}

func (m *MockEventEmitter) Emit(eventName string, data any) {
	m.Events = append(m.Events, struct {
		Name string
		Data any
	}{Name: eventName, Data: data})
}

func TestChatFlow(t *testing.T) {
	// 1. Setup specific test workspace
	tempDir, err := os.MkdirTemp("", "multiverse-chat-flow-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	wsStore := ide.NewWorkspaceStore(filepath.Join(tempDir, "workspaces"))
	projectRoot := filepath.Join(tempDir, "project")

	// Ensure project root exists
	err = os.MkdirAll(projectRoot, 0755)
	require.NoError(t, err)

	ws := &ide.Workspace{
		ProjectRoot: projectRoot,
		Version:     "1.0",
		DisplayName: "Chat Flow Test",
	}
	err = wsStore.SaveWorkspace(ws)
	require.NoError(t, err)

	wsID := wsStore.GetWorkspaceID(projectRoot)
	wsDir := wsStore.GetWorkspaceDir(wsID)

	// 2. Initialize Components
	taskStore := orchestrator.NewTaskStore(wsDir)
	sessionStore := chat.NewChatSessionStore(wsDir)
	metaClient := meta.NewMockClient()
	eventEmitter := &MockEventEmitter{}

	// Initialize ChatHandler
	handler := chat.NewHandler(metaClient, taskStore, sessionStore, wsID, projectRoot, nil, eventEmitter)

	// 3. Create Session
	ctx := context.Background()
	session, err := handler.CreateSession(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, session.ID)

	// 4. Send Message (Trigger Decompose)
	// The mock client ignores the message content and returns hardcoded tasks
	resp, err := handler.HandleMessage(ctx, session.ID, "Create a hello world task")
	require.NoError(t, err)

	// 5. Verify Response
	assert.NotNil(t, resp)
	assert.Len(t, resp.GeneratedTasks, 2, "Expected 2 mocked tasks")

	// Verify Task Content (from Mock Decompose)
	// Task 1: "Mock概念設計タスク"
	// Task 2: "Mock実装タスク"

	var task1, task2 orchestrator.Task

	// Find tasks by title
	for _, task := range resp.GeneratedTasks {
		if task.Title == "Mock概念設計タスク" {
			task1 = task
		} else if task.Title == "Mock実装タスク" {
			task2 = task
		}
	}

	assert.NotEmpty(t, task1.ID)
	assert.NotEmpty(t, task2.ID)
	assert.Equal(t, 1, task1.WBSLevel)
	assert.Equal(t, 3, task2.WBSLevel)

	// Verify Dependencies
	// Task 2 depends on Task 1
	require.Len(t, task2.Dependencies, 1)
	assert.Equal(t, task1.ID, task2.Dependencies[0])

	// 6. Verify Persistence in TaskStore
	savedTask1, err := taskStore.LoadTask(task1.ID)
	require.NoError(t, err)
	assert.Equal(t, task1.Title, savedTask1.Title)

	savedTask2, err := taskStore.LoadTask(task2.ID)
	require.NoError(t, err)
	assert.Equal(t, task2.Title, savedTask2.Title)

	// 7. Verify Real-time Events
	// We expect 'chat:progress' events and 'task:created' events

	var taskCreatedCount int
	var progressCount int

	for _, event := range eventEmitter.Events {
		if event.Name == orchestrator.EventTaskCreated {
			taskCreatedCount++
			data, ok := event.Data.(orchestrator.TaskCreatedEvent)
			if ok {
				assert.NotEmpty(t, data.Task.ID)
			}
		} else if event.Name == orchestrator.EventChatProgress {
			progressCount++
		}
	}

	assert.Equal(t, 2, taskCreatedCount, "Expected 2 task:created events")
	assert.Greater(t, progressCount, 0, "Expected chat:progress events")
}
