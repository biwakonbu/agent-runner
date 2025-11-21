package worker

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/biwakonbu/agent-runner/internal/core"
	"github.com/biwakonbu/agent-runner/pkg/config"
)

type Executor struct {
	Config      config.WorkerConfig
	Sandbox     SandboxProvider
	RepoPath    string
	containerID string // 持続的なコンテナを保持
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
	}, nil
}

func (e *Executor) RunWorker(ctx context.Context, prompt string, env map[string]string) (*core.WorkerRunResult, error) {
	// Verify container is running
	if e.containerID == "" {
		return nil, fmt.Errorf("container not started: call Start() first")
	}

	// Construct Command
	// "codex exec --sandbox workspace-write --json --cwd /workspace/project '<prompt>'"
	cmd := []string{
		"codex", "exec",
		"--sandbox", "workspace-write",
		"--json",
		"--cwd", "/workspace/project",
		prompt,
	}

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

	return res, nil
}

// Start starts a persistent container for the task
func (e *Executor) Start(ctx context.Context) error {
	if e.containerID != "" {
		return fmt.Errorf("container already started (ID: %s)", e.containerID)
	}

	image := e.Config.DockerImage
	if image == "" {
		image = "ghcr.io/biwakonbu/agent-runner-codex:latest" // Default
	}

	repoPath := e.RepoPath
	if repoPath == "" {
		absRepo, err := filepath.Abs(".")
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}
		repoPath = absRepo
	}

	containerID, err := e.Sandbox.StartContainer(ctx, image, repoPath, nil)
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	e.containerID = containerID
	return nil
}

// Stop stops the persistent container
func (e *Executor) Stop(ctx context.Context) error {
	if e.containerID == "" {
		return fmt.Errorf("no container to stop")
	}

	// Store containerID before clearing (for error message)
	containerID := e.containerID

	// Clear containerID first to prevent resource leak
	// even if StopContainer fails
	e.containerID = ""

	err := e.Sandbox.StopContainer(ctx, containerID)
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	return nil
}
