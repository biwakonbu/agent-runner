package main

import (
	"context"
	"fmt"
	"os"

	"github.com/biwakonbu/agent-runner/internal/ide"
	"github.com/biwakonbu/agent-runner/internal/orchestrator"
	"github.com/biwakonbu/agent-runner/internal/orchestrator/ipc"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	workspaceStore *ide.WorkspaceStore
	taskStore      *orchestrator.TaskStore
	scheduler      *orchestrator.Scheduler
	currentWS      *ide.Workspace
}

// NewApp creates a new App application struct
func NewApp() *App {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	baseDir := fmt.Sprintf("%s/.multiverse/workspaces", homeDir)
	return &App{
		workspaceStore: ide.NewWorkspaceStore(baseDir),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// SelectWorkspace opens a directory selection dialog and loads the workspace.
func (a *App) SelectWorkspace() string {
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Project Root",
	})
	if err != nil {
		return ""
	}
	if selection == "" {
		return ""
	}

	id := a.workspaceStore.GetWorkspaceID(selection)
	ws, err := a.workspaceStore.LoadWorkspace(id)
	if err != nil {
		// Create new workspace if not exists
		ws = &ide.Workspace{
			Version:     "1.0",
			ProjectRoot: selection,
			DisplayName: selection, // Simplified for now
		}
		if err := a.workspaceStore.SaveWorkspace(ws); err != nil {
			runtime.LogErrorf(a.ctx, "Failed to save workspace: %v", err)
			return ""
		}
	}

	a.currentWS = ws

	// Initialize TaskStore and Scheduler for this workspace
	wsDir := a.workspaceStore.GetWorkspaceDir(id)
	a.taskStore = orchestrator.NewTaskStore(wsDir)
	queue := ipc.NewFilesystemQueue(wsDir)
	a.scheduler = orchestrator.NewScheduler(a.taskStore, queue)

	return id
}

// GetWorkspace returns the workspace details.
func (a *App) GetWorkspace(id string) *ide.Workspace {
	ws, err := a.workspaceStore.LoadWorkspace(id)
	if err != nil {
		return nil
	}
	return ws
}

// ListTasks returns all tasks in the current workspace.
// Note: In a real app, we might want pagination or filtering.
// For now, we'll just list all jsonl files in the tasks dir.
func (a *App) ListTasks() []orchestrator.Task {
	if a.taskStore == nil {
		return []orchestrator.Task{}
	}

	dir := a.taskStore.GetTaskDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []orchestrator.Task{}
	}

	var tasks []orchestrator.Task
	for _, entry := range entries {
		if !entry.IsDir() && len(entry.Name()) > 6 && entry.Name()[len(entry.Name())-6:] == ".jsonl" {
			id := entry.Name()[:len(entry.Name())-6]
			task, err := a.taskStore.LoadTask(id)
			if err == nil {
				tasks = append(tasks, *task)
			}
		}
	}
	return tasks
}

// CreateTask creates a new task.
func (a *App) CreateTask(title string, poolID string) *orchestrator.Task {
	if a.taskStore == nil {
		return nil
	}

	task := &orchestrator.Task{
		ID:     fmt.Sprintf("task-%d", os.Getpid()), // Simplified ID generation
		Title:  title,
		Status: orchestrator.TaskStatusPending,
		PoolID: poolID,
	}
	// In reality, use a better ID generator (UUID)

	if err := a.taskStore.SaveTask(task); err != nil {
		runtime.LogErrorf(a.ctx, "Failed to save task: %v", err)
		return nil
	}
	return task
}

// RunTask schedules a task for execution.
func (a *App) RunTask(taskID string) error {
	if a.scheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}
	return a.scheduler.ScheduleTask(taskID)
}
