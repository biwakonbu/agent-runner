package worker

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"syscall"
)

// LocalSandbox implements SandboxProvider but runs commands locally on the host.
// WARNING: This provides NO isolation. Use only for trusted CLI tools or testing.
type LocalSandbox struct {
	Workdir string
}

func NewLocalSandbox(workdir string) *LocalSandbox {
	return &LocalSandbox{Workdir: workdir}
}

// StartContainer acts as a no-op setup for LocalSandbox
func (s *LocalSandbox) StartContainer(ctx context.Context, image string, repoPath string, env map[string]string) (string, error) {
	// For LocalSandbox, we don't start a container. return a dummy ID.
	return "local-host", nil
}

// Exec runs the command locally using os/exec
func (s *LocalSandbox) Exec(ctx context.Context, containerID string, cmd []string, stdin io.Reader) (int, string, error) {
	if len(cmd) == 0 {
		return 0, "", fmt.Errorf("empty command")
	}

	c := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	c.Dir = s.Workdir
	c.Stdin = stdin

	// Combine stdout and stderr
	output, err := c.CombinedOutput()

	exitCode := 0
	if err != nil {
		// try to get exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			if ws, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = ws.ExitStatus()
			} else {
				exitCode = 1
			}
		} else {
			// Other errors (e.g. command not found)
			return 1, string(output), err
		}
	}

	return exitCode, string(output), nil
}

// StopContainer acts as a no-op teardown for LocalSandbox
func (s *LocalSandbox) StopContainer(ctx context.Context, containerID string) error {
	// Nothing to stop (unless we tracked background processes, but Exec waits)
	return nil
}
