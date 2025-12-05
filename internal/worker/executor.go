package worker

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/biwakonbu/agent-runner/internal/core"
	"github.com/biwakonbu/agent-runner/internal/logging"
	"github.com/biwakonbu/agent-runner/pkg/config"
)

type Executor struct {
	Config      config.WorkerConfig
	Sandbox     SandboxProvider
	RepoPath    string
	containerID string // 持続的なコンテナを保持
	logger      *slog.Logger
}

func NewExecutor(cfg config.WorkerConfig, repoPath string) (*Executor, error) {
	sb, err := NewSandboxManager()
	if err != nil {
		return nil, err
	}
	return &Executor{
		Config:      cfg,
		Sandbox:     sb,
		RepoPath:    repoPath,
		containerID: "", // 未初期化
		logger:      logging.WithComponent(slog.Default(), "worker-executor"),
	}, nil
}

// SetLogger sets a custom logger for the executor
func (e *Executor) SetLogger(logger *slog.Logger) {
	e.logger = logging.WithComponent(logger, "worker-executor")
}

func (e *Executor) RunWorker(ctx context.Context, prompt string, env map[string]string) (*core.WorkerRunResult, error) {
	logger := logging.WithTraceID(e.logger, ctx)

	// Verify container is running
	if e.containerID == "" {
		logger.Error("container not started")
		return nil, fmt.Errorf("container not started: call Start() first")
	}

	// Determine timeout
	timeoutSec := e.Config.MaxRunTimeSec
	if timeoutSec <= 0 {
		timeoutSec = 1800 // Default 30 minutes
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	// Construct Command
	// "codex exec --sandbox workspace-write --json --cwd /workspace/project '<prompt>'"
	cmd := []string{
		"codex", "exec",
		"--sandbox", "workspace-write",
		"--json",
		"--cwd", "/workspace/project",
		prompt,
	}

	logger.Info("executing worker command",
		slog.String("container_id", e.containerID[:12]),
		slog.Int("prompt_length", len(prompt)),
		slog.Int("timeout_sec", timeoutSec),
	)
	logger.Debug("worker command details",
		slog.String("prompt", prompt),
		slog.Any("cmd", cmd),
	)

	// Exec on persistent container
	start := time.Now()
	exitCode, output, err := e.Sandbox.Exec(ctx, e.containerID, cmd)
	finish := time.Now()

	res := &core.WorkerRunResult{
		ID:         fmt.Sprintf("run-%d", start.Unix()),
		StartedAt:  start,
		FinishedAt: finish,
		ExitCode:   exitCode,
		RawOutput:  output,
		Summary:    "Worker executed", // Could parse JSON output if needed
		Error:      err,
	}

	durationMs := float64(finish.Sub(start).Milliseconds())
	if err != nil {
		logger.Error("worker execution failed",
			slog.Int("exit_code", exitCode),
			slog.Float64("duration_ms", durationMs),
			slog.Any("error", err),
		)
	} else {
		logger.Info("worker execution completed",
			slog.Int("exit_code", exitCode),
			slog.Int("output_length", len(output)),
			slog.Float64("duration_ms", durationMs),
		)
		logger.Debug("worker output", slog.String("output", output))
	}

	return res, nil
}

// Start starts a persistent container for the task
func (e *Executor) Start(ctx context.Context) error {
	logger := logging.WithTraceID(e.logger, ctx)

	if e.containerID != "" {
		logger.Warn("container already started", slog.String("container_id", e.containerID[:12]))
		return fmt.Errorf("container already started (ID: %s)", e.containerID)
	}

	image := e.Config.DockerImage
	if image == "" {
		image = "ghcr.io/biwakonbu/agent-runner-codex:latest" // Default
	}

	// Resolve RepoPath to absolute path
	repoPath := e.RepoPath
	if repoPath == "" {
		repoPath = "."
	}
	absRepo, err := filepath.Abs(repoPath)
	if err != nil {
		logger.Error("failed to get absolute path", slog.String("repo_path", repoPath), slog.Any("error", err))
		return fmt.Errorf("failed to get absolute path for %s: %w", repoPath, err)
	}
	repoPath = absRepo

	logger.Info("starting container",
		slog.String("image", image),
		slog.String("repo_path", repoPath),
	)

	start := time.Now()
	containerID, err := e.Sandbox.StartContainer(ctx, image, repoPath, nil)
	if err != nil {
		logger.Error("failed to start container",
			slog.String("image", image),
			slog.Any("error", err),
			logging.LogDuration(start),
		)
		return fmt.Errorf("failed to start container: %w", err)
	}

	e.containerID = containerID
	logger.Info("container started",
		slog.String("container_id", containerID[:12]),
		logging.LogDuration(start),
	)
	return nil
}

// Stop stops the persistent container
func (e *Executor) Stop(ctx context.Context) error {
	logger := logging.WithTraceID(e.logger, ctx)

	if e.containerID == "" {
		logger.Warn("no container to stop")
		return fmt.Errorf("no container to stop")
	}

	// Store containerID before clearing (for error message)
	containerID := e.containerID
	logger.Info("stopping container", slog.String("container_id", containerID[:12]))

	// Clear containerID first to prevent resource leak
	// even if StopContainer fails
	e.containerID = ""

	start := time.Now()
	err := e.Sandbox.StopContainer(ctx, containerID)
	if err != nil {
		logger.Error("failed to stop container",
			slog.String("container_id", containerID[:12]),
			slog.Any("error", err),
			logging.LogDuration(start),
		)
		return fmt.Errorf("failed to stop container: %w", err)
	}

	logger.Info("container stopped",
		slog.String("container_id", containerID[:12]),
		logging.LogDuration(start),
	)
	return nil
}
