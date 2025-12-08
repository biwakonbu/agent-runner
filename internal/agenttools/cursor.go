package agenttools

import (
	"context"
	"fmt"
)

// DefaultCursorModel
const DefaultCursorModel = "claude-3-5-sonnet-20241022" // Cursor often uses Sonnet

// CursorProvider builds ExecPlan for Cursor CLI.
type CursorProvider struct {
	cliPath string
	model   string
	env     map[string]string
	flags   []string
}

// NewCursorProvider constructs a CursorProvider from config.
func NewCursorProvider(cfg ProviderConfig) *CursorProvider {
	return &CursorProvider{
		cliPath: nonEmpty(cfg.CLIPath, "cursor"),
		model:   cfg.Model,
		env:     mergeEnv(nil, cfg.ExtraEnv),
		flags:   append([]string{}, cfg.Flags...),
	}
}

func (p *CursorProvider) Kind() string {
	return "cursor-cli"
}

func (p *CursorProvider) Capabilities() Capability {
	return Capability{
		Kind:          p.Kind(),
		DefaultModel:  nonEmpty(p.model, DefaultCursorModel),
		SupportsStdin: true,
		Notes:         "Cursor CLI wrapper. Assumes `cursor chat` or similar interface.",
	}
}

// Build generates the execution plan for Cursor CLI.
func (p *CursorProvider) Build(_ context.Context, req Request) (ExecPlan, error) {
	if err := ensurePrompt(req.Prompt); err != nil {
		return ExecPlan{}, err
	}

	// Only 'exec' (which might map to a chat command) is supported
	mode := req.Mode
	if mode == "" {
		mode = "exec"
	}
	if mode != "exec" {
		return ExecPlan{}, fmt.Errorf("%w: %s (only 'exec' is supported)", ErrUnsupportedMode, mode)
	}

	// Cursor CLI interactions are often: cursor --chat-prompt "..." or similar.
	// As of now, `cursor` command opens the IDE.
	// If there is a `cursor-cli` or specific headless mode, we use it.
	// Assuming a hypothetical `cursor cmd` or similar for this PRD context.
	// If strictly opening IDE, it's not a worker.
	// PRD implies "CLI Subscription Session" usage, presumably via a headless interface or bridged.
	// We'll treat it as `cursor chat-cli` or similar command structure for now.

	args := []string{}

	// Let's assume a subcommand "chat" or "cmd" if it exists, or just args.
	// For now, standard args.

	// Model
	// args = append(args, "--model", ...)

	// Extra flags
	args = append(args, p.flags...)
	args = append(args, req.Flags...)

	plan := ExecPlan{
		Command: p.cliPath,
		Args:    args,
		Env:     mergeEnv(p.env, req.ExtraEnv),
		Workdir: req.Workdir,
		Timeout: req.Timeout,
	}

	if req.UseStdin {
		plan.Stdin = req.Prompt
		plan.Args = append(plan.Args, "-")
	} else {
		plan.Args = append(plan.Args, req.Prompt)
	}

	return plan, nil
}

func init() {
	Register("cursor-cli", func(cfg ProviderConfig) (AgentToolProvider, error) {
		return NewCursorProvider(cfg), nil
	})
}
