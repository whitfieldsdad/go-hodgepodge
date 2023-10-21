package hodgepodge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecutePowerShellCommand(t *testing.T) {
	subprocess, err := ExecuteShellCommand("whoami", PowerShell, nil)
	if err != nil {
		t.Fatalf("Failed to execute command: %s\n", err)
	}
	assert.Equal(t, 0, *subprocess.ExitCode, "Exit code should be 0")
}

func TestExecuteShellCommand(t *testing.T) {
	subprocess, err := ExecuteShellCommand("whoami", Sh, nil)
	if err != nil {
		t.Fatalf("Failed to execute command: %s\n", err)
	}
	assert.Equal(t, 0, *subprocess.ExitCode, "Exit code should be 0")
}
