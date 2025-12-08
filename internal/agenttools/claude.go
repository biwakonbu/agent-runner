package agenttools

import (
	"context"
	"fmt"
)

// DefaultClaudeModel defines the default model for Claude Code.
// 参照: https://docs.anthropic.com/en/docs/claude-code
const DefaultClaudeModel = "claude-3-5-sonnet-20241022"

// ClaudeProvider builds ExecPlan for Claude Code CLI.
// Wrapper for `claude-code` or `claude` CLI.
type ClaudeProvider struct {
	cliPath string
	model   string
	env     map[string]string
	flags   []string
}

// NewClaudeProvider constructs a ClaudeProvider from config.
func NewClaudeProvider(cfg ProviderConfig) *ClaudeProvider {
	return &ClaudeProvider{
		cliPath: nonEmpty(cfg.CLIPath, "claude"),
		model:   cfg.Model,
		env:     mergeEnv(nil, cfg.ExtraEnv),
		flags:   append([]string{}, cfg.Flags...),
	}
}

func (p *ClaudeProvider) Kind() string {
	return "claude-code"
}

func (p *ClaudeProvider) Capabilities() Capability {
	return Capability{
		Kind:          p.Kind(),
		DefaultModel:  nonEmpty(p.model, DefaultClaudeModel),
		SupportsStdin: true,
		Notes:         "Claude Code CLI wrapper. Assumes `claude [prompt]` interface.",
	}
}

// Build generates the execution plan for Claude Code CLI.
func (p *ClaudeProvider) Build(_ context.Context, req Request) (ExecPlan, error) {
	if err := ensurePrompt(req.Prompt); err != nil {
		return ExecPlan{}, err
	}

	mode := req.Mode
	if mode == "" {
		mode = "exec"
	}

	// claude-code usually handles conversation or single shot.
	// We map 'exec' to single shot or piped input.
	if mode != "exec" {
		return ExecPlan{}, fmt.Errorf("%w: %s (only 'exec' is supported)", ErrUnsupportedMode, mode)
	}

	args := []string{}

	// JSON output support check
	// Note: claude-code might not support --json for all commands, checking docs or assuming future support.
	// For now, if requested, we add it.
	jsonOutput := true
	if v, ok := req.ToolSpecific["json_output"].(bool); ok {
		jsonOutput = v
	}
	// TODO: Verify if claude cli supports --json. Currently assuming standard behavior or ignoring if unknown.
	// Let's assume it doesn't support --json by default unless we know specific flag.
	// Removing --json for safety unless explicitly needed or verified.
	_ = jsonOutput

	// Model specification
	if p.model != "" || req.Model != "" {
		model := nonEmpty(req.Model, p.model)
		args = append(args, "--model", model)
	}

	// Temperature mapping (if supported)
	// Claude Code CLI flags are TBD.
	// TODO: Add temperature flag when Claude Code CLI supports it
	_ = req.Temperature // Currently unused, pending CLI support

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

	// Prompt handling
	if req.UseStdin {
		// Piped input
		// echo "prompt" | claude
		plan.Stdin = req.Prompt
		// args usually don't need "-" for simplified CLIs, but depends.
		// If claude-code detects stdin, it uses it.
	} else {
		plan.Args = append(plan.Args, req.Prompt)
	}

	return plan, nil
}

func init() {
	Register("claude-code", func(cfg ProviderConfig) (AgentToolProvider, error) {
		return NewClaudeProvider(cfg), nil
	})
}
