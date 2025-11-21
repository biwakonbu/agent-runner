package mock

import (
	"context"

	"github.com/biwakonbu/agent-runner/internal/core"
)

type WorkerExecutor struct {
	RunWorkerFunc func(ctx context.Context, prompt string, env map[string]string) (*core.WorkerRunResult, error)
	StartFunc     func(ctx context.Context) error
	StopFunc      func(ctx context.Context) error
}

func (w *WorkerExecutor) RunWorker(ctx context.Context, prompt string, env map[string]string) (*core.WorkerRunResult, error) {
	if w.RunWorkerFunc != nil {
		return w.RunWorkerFunc(ctx, prompt, env)
	}
	return nil, nil
}

func (w *WorkerExecutor) Start(ctx context.Context) error {
	if w.StartFunc != nil {
		return w.StartFunc(ctx)
	}
	return nil
}

func (w *WorkerExecutor) Stop(ctx context.Context) error {
	if w.StopFunc != nil {
		return w.StopFunc(ctx)
	}
	return nil
}

// NewMockWorkerExecutor creates a mock WorkerExecutor with default behavior
func NewMockWorkerExecutor() *WorkerExecutor {
	return &WorkerExecutor{}
}
