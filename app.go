package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/biwakonbu/agent-runner/internal/ide"
	"github.com/biwakonbu/agent-runner/internal/logging"
	"github.com/biwakonbu/agent-runner/internal/orchestrator"
	"github.com/biwakonbu/agent-runner/internal/orchestrator/ipc"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	workspaceStore *ide.WorkspaceStore
	taskStore      *orchestrator.TaskStore
	scheduler      *orchestrator.Scheduler
	executor       *orchestrator.Executor
	currentWS      *ide.Workspace
	logger         *slog.Logger
}

// NewApp creates a new App application struct
func NewApp() *App {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	baseDir := fmt.Sprintf("%s/.multiverse/workspaces", homeDir)
	logger := logging.NewLogger(logging.DefaultConfig())
	logger = logging.WithComponent(logger, "ide-app")

	return &App{
		workspaceStore: ide.NewWorkspaceStore(baseDir),
		logger:         logger,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.logger.Info("IDE application started")
}

// SelectWorkspace opens a directory selection dialog and loads the workspace.
func (a *App) SelectWorkspace() string {
	a.logger.Info("opening workspace selection dialog")
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Project Root",
	})
	if err != nil {
		a.logger.Error("workspace selection dialog failed", slog.Any("error", err))
		return ""
	}
	if selection == "" {
		a.logger.Info("workspace selection canceled")
		return ""
	}

	a.logger.Info("workspace selected", slog.String("path", selection))
	return a.LoadWorkspace(selection)
}

// LoadWorkspace loads a workspace from a given path.
func (a *App) LoadWorkspace(path string) string {
	a.logger.Info("loading workspace", slog.String("path", path))

	id := a.workspaceStore.GetWorkspaceID(path)
	ws, err := a.workspaceStore.LoadWorkspace(id)
	if err != nil {
		a.logger.Info("creating new workspace", slog.String("path", path))
		// Create new workspace if not exists
		ws = &ide.Workspace{
			Version:     "1.0",
			ProjectRoot: path,
			DisplayName: path, // Simplified for now
		}
		if err := a.workspaceStore.SaveWorkspace(ws); err != nil {
			a.logger.Error("failed to save workspace", slog.Any("error", err))
			runtime.LogErrorf(a.ctx, "Failed to save workspace: %v", err)
			return ""
		}
		a.logger.Info("new workspace created", slog.String("id", a.workspaceStore.GetWorkspaceID(path)))
	}

	a.currentWS = ws

	// Initialize TaskStore and Scheduler for this workspace
	wsDir := a.workspaceStore.GetWorkspaceDir(id)
	a.taskStore = orchestrator.NewTaskStore(wsDir)
	queue := ipc.NewFilesystemQueue(wsDir)
	a.scheduler = orchestrator.NewScheduler(a.taskStore, queue)

	// Initialize Executor with path to agent-runner binary
	// In production, this would be configured or detected
	a.executor = orchestrator.NewExecutor("./agent-runner", a.taskStore)

	a.logger.Info("workspace loaded successfully",
		slog.String("id", id),
		slog.String("workspace_dir", wsDir),
	)
	return id
}

// ScheduleTask adds a task to the execution queue
func (a *App) ScheduleTask(taskID string) error {
	a.logger.Info("scheduling task", slog.String("task_id", taskID))
	if a.scheduler == nil {
		a.logger.Error("scheduler not initialized")
		return fmt.Errorf("scheduler not initialized")
	}
	err := a.scheduler.ScheduleTask(taskID)
	if err != nil {
		a.logger.Error("failed to schedule task", slog.String("task_id", taskID), slog.Any("error", err))
		return err
	}
	a.logger.Info("task scheduled successfully", slog.String("task_id", taskID))
	return nil
}

// GetWorkspace returns the workspace details.
func (a *App) GetWorkspace(id string) *ide.Workspace {
	a.logger.Debug("getting workspace", slog.String("id", id))
	ws, err := a.workspaceStore.LoadWorkspace(id)
	if err != nil {
		a.logger.Warn("workspace not found", slog.String("id", id), slog.Any("error", err))
		return nil
	}
	return ws
}

// ListTasks returns all tasks in the current workspace.
// Note: In a real app, we might want pagination or filtering.
// For now, we'll just list all jsonl files in the tasks dir.
func (a *App) ListTasks() []orchestrator.Task {
	a.logger.Debug("listing tasks")
	if a.taskStore == nil {
		a.logger.Warn("task store not initialized")
		return []orchestrator.Task{}
	}

	dir := a.taskStore.GetTaskDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		a.logger.Warn("failed to read task directory", slog.String("dir", dir), slog.Any("error", err))
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
	a.logger.Info("tasks listed", slog.Int("count", len(tasks)))
	return tasks
}

// CreateTask creates a new task.
func (a *App) CreateTask(title string, poolID string) *orchestrator.Task {
	a.logger.Info("creating task", slog.String("title", title), slog.String("pool_id", poolID))
	if a.taskStore == nil {
		a.logger.Error("task store not initialized")
		return nil
	}

	taskID := fmt.Sprintf("task-%s", uuid.New().String()[:8])
	task := &orchestrator.Task{
		ID:     taskID,
		Title:  title,
		Status: orchestrator.TaskStatusPending,
		PoolID: poolID,
	}

	if err := a.taskStore.SaveTask(task); err != nil {
		a.logger.Error("failed to save task", slog.Any("error", err))
		runtime.LogErrorf(a.ctx, "Failed to save task: %v", err)
		return nil
	}
	a.logger.Info("task created", slog.String("task_id", taskID))
	return task
}

// RunTask executes a task using the AgentRunner Core.
func (a *App) RunTask(taskID string) error {
	a.logger.Info("running task", slog.String("task_id", taskID))
	if a.executor == nil {
		a.logger.Error("executor not initialized")
		return fmt.Errorf("executor not initialized")
	}

	task, err := a.taskStore.LoadTask(taskID)
	if err != nil {
		a.logger.Error("failed to load task", slog.String("task_id", taskID), slog.Any("error", err))
		return fmt.Errorf("failed to load task: %w", err)
	}

	// Create trace ID for this task execution
	traceID := uuid.New().String()
	ctx := logging.ContextWithTraceID(a.ctx, traceID)
	a.logger.Info("starting task execution",
		slog.String("task_id", taskID),
		slog.String("trace_id", traceID),
	)

	// Execute task in background
	go func() {
		logger := logging.WithTraceID(a.logger, ctx)
		_, err := a.executor.ExecuteTask(ctx, task)
		if err != nil {
			logger.Error("task execution failed", slog.String("task_id", taskID), slog.Any("error", err))
			runtime.LogErrorf(a.ctx, "Task execution failed: %v", err)
		} else {
			logger.Info("task execution completed", slog.String("task_id", taskID))
		}
	}()

	return nil
}
