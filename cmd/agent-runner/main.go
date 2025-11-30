package main

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/biwakonbu/agent-runner/internal/cli"
	"github.com/biwakonbu/agent-runner/internal/core"
	"github.com/biwakonbu/agent-runner/internal/meta"
	"github.com/biwakonbu/agent-runner/internal/note"
	"github.com/biwakonbu/agent-runner/internal/worker"
	"github.com/biwakonbu/agent-runner/pkg/config"
	"gopkg.in/yaml.v3"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if err := Run(context.Background(), os.Stdin, os.Stdout, os.Stderr, logger); err != nil {
		slog.Error("application failed", "err", err)
		os.Exit(1)
	}
}

// Run is the main entry point for the application, extracted for testing.
func Run(ctx context.Context, stdin io.Reader, _, _ io.Writer, logger *slog.Logger) error {
	// 1. Parse CLI flags
	flags, err := cli.ParseFlags(os.Args[1:], os.Stderr)
	if err != nil {
		return err
	}

	// 2. Read YAML from stdin
	bytes, err := io.ReadAll(stdin)
	if err != nil {
		return err
	}

	var cfg config.TaskConfig
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return err
	}

	// 3. Initialize Components
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		logger.Warn("OPENAI_API_KEY not set, using mock mode")
	}

	// Resolve Meta Model ID
	metaModel := cli.ResolveMetaModel(flags.MetaModel, cfg.Runner.Meta.Model)
	logger.Info("resolved meta model", "model", metaModel)

	metaClient := meta.NewClient(cfg.Runner.Meta.Kind, apiKey, metaModel, cfg.Runner.Meta.SystemPrompt)

	workerExecutor, err := worker.NewExecutor(cfg.Runner.Worker, cfg.Task.Repo)
	if err != nil {
		return err
	}

	noteWriter := note.NewWriter()

	runner := core.NewRunner(&cfg, metaClient, workerExecutor, noteWriter)

	// 4. Run
	logger.Info("starting task", "title", cfg.Task.Title, "id", cfg.Task.ID)

	result, err := runner.Run(ctx)
	if err != nil {
		return err
	}

	logger.Info("task completed", "state", result.State)
	return nil
}
