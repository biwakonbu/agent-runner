package e2e_test

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/biwakonbu/agent-runner/internal/worker"
	"github.com/biwakonbu/agent-runner/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockSandbox implements worker.SandboxProvider for testing
type MockSandbox struct{}

func (m *MockSandbox) StartContainer(ctx context.Context, image string, repoPath string, env map[string]string) (string, error) {
	return "mock-container-id", nil
}

func (m *MockSandbox) StopContainer(ctx context.Context, containerID string) error {
	return nil
}

func (m *MockSandbox) Exec(ctx context.Context, containerID string, cmd []string, stdin io.Reader) (int, string, error) {
	return 0, "mock output", nil
}

// TestCLISessionCheck verifies that the Worker Executor correctly enforces session requirements.
// This supports FR-P4-001 and AC-P4-07.
func TestCLISessionCheck(t *testing.T) {
	// 1. Setup temporary directory for HOME to simulate no auth.json
	tempHome, err := os.MkdirTemp("", "cli-session-test-home-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempHome)

	// Save original HOME
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	// Set HOME to temp dir
	os.Setenv("HOME", tempHome)

	// Clear CODEX_API_KEY
	os.Unsetenv("CODEX_API_KEY")

	// 2. Initialize Executor with Mock Sandbox
	cfg := config.WorkerConfig{
		Kind: "codex-cli",
	}
	executor := &worker.Executor{
		Config:   cfg,
		Sandbox:  &MockSandbox{},
		RepoPath: tempHome, // Dummy path
	}
	// Note: We are manually constructing Executor to inject mocks.
	// Real constructor NewExecutor establishes real SandboxManager.

	// 3. Call Start (Should Fail due to missing session)
	// Start calls verifyCodexSession
	ctx := context.Background()
	err = executor.Start(ctx)

	// 4. Assertions
	require.Error(t, err, "Expected Start to fail when no credentials/session are present")
	assert.Contains(t, err.Error(), "Codex CLI", "Error message should mention 'Codex CLI'")
	assert.Contains(t, err.Error(), "codex login", "Error message should suggest 'codex login'")

	t.Logf("Got expected error: %v", err)
}
