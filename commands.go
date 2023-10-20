package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"runtime"

	"github.com/charmbracelet/log"
	"github.com/google/shlex"
	"github.com/pkg/errors"
)

const (
	Native        = "native"
	Sh            = "sh"
	Bash          = "bash"
	PowerShell    = "powershell"
	CommandPrompt = "command_prompt"
	Shell         = "shell"
)

var (
	ShShim                = []string{"sh", "-c", "%s"}
	BashShim              = []string{"bash", "-c", "%s"}
	CommandPromptShim     = []string{"cmd", "/c", "%s"}
	WindowsPowerShellShim = []string{"powershell", "-ExecutionPolicy", "Bypass", "-Command", "%s"}
	PosixPowerShellShim   = []string{"pwsh", "-Command", "%s"}
)

var (
	commandShims = map[string][]string{
		Sh:            ShShim,
		Bash:          BashShim,
		CommandPrompt: CommandPromptShim,
		Shell:         GetShellCommandShim(),
		PowerShell:    GetPowerShellCommandShim(),
	}
)

// Command stores information about a command to be executed on a host. If the command type is "", the default command type for the current platform will be used. If the command type is "native", the command will be executed without wrapping it in a shell command (e.g. using `bash -c`).
type Command struct {
	Command     string `json:"command"`
	CommandType string `json:"command_type"`
}

// Run executes a command, blocks until it completes, and returns information about the subprocess that was executed.
func (c Command) Run(opts *ExecuteCommandOptions) (*ExecutedCommand, error) {
	command := c.Command
	commandType := c.CommandType

	var err error
	if commandType == "" {
		commandType = GetDefaultCommandCommandType()
	}
	if commandType != Native {
		command, err = WrapCommand(command, commandType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to wrap command")
		}
	}
	p, err := ExecuteCommand(command, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute command")
	}
	return &ExecutedCommand{
		Time:      time.Now(),
		Command:   c.Command,
		Processes: []Process{*p},
	}, nil
}

// ExecutedCommand stores information about a command that has been executed on a host.
type ExecutedCommand struct {
	Id        string    `json:"id"`
	Time      time.Time `json:"time"`
	Command   string    `json:"command"`
	Error     string    `json:"error,omitempty"`
	Processes []Process `json:"processes,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
}

// ExecuteCommand executes a command, blocks until it completes, and returns information about the subprocess that was executed.
func ExecuteCommand(command string, opts *ExecuteCommandOptions) (*Process, error) {
	if command == "" {
		return nil, errors.New("command cannot be empty")
	}
	args, err := shlex.Split(command)
	if err != nil {
		return nil, errors.Wrap(err, "failed to split command")
	}
	p, err := executeShellCommand(args, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute command")
	}
	return p, nil
}

// GetPowerShellCommandShim returns the shim for executing PowerShell commands on the current platform (i.e. using either powershell.exe or pwsh).
func GetPowerShellCommandShim() []string {
	if runtime.GOOS == "windows" {
		return WindowsPowerShellShim
	} else {
		return PosixPowerShellShim
	}
}

// GetShellCommandShim returns the shim for executing shell commands on the current platform (i.e. using either cmd.exe or sh).
func GetShellCommandShim() []string {
	if runtime.GOOS == "windows" {
		return CommandPromptShim
	} else {
		return ShShim
	}
}

// GetDefaultCommandCommandType returns the default command type for the current platform.
func GetDefaultCommandCommandType() string {
	if runtime.GOOS == "windows" {
		return PowerShell
	} else {
		return Sh
	}
}

// WrapCommand wraps a command in a shell command.
func WrapCommand(command, commandType string) (string, error) {
	if commandType == "" {
		commandType = GetDefaultCommandCommandType()
	}
	shim, ok := commandShims[commandType]
	if !ok {
		return "", errors.New("invalid shim type")
	}

	// Replace newlines with the appropriate separator for running multiple commands.
	n := strings.Count(command, "\n")
	s := "\n"
	switch commandType {
	case CommandPrompt:
		s = " & "
	case PowerShell:
		s = "; "
	case Sh:
		s = "; "
	}
	command = strings.Replace(command, "\n", s, n-1)

	// Substitute %s for the command to execute.
	for i, token := range shim {
		if strings.Contains(token, "%s") {
			shim[i] = fmt.Sprintf(token, command)
		}
	}
	return strings.Join(shim, " "), nil
}

// ExecuteShellCommand wraps a command, executes it, blocks until it completes, and returns information about the subprocess that was executed.
func ExecuteShellCommand(command, commandType string, opts *ExecuteCommandOptions) (*Process, error) {
	w, err := WrapCommand(command, commandType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wrap command")
	}
	return ExecuteCommand(w, opts)
}

func executeShellCommand(args []string, opts *ExecuteCommandOptions) (*Process, error) {
	command := strings.Join(args, " ")
	log.Infof("Executing command: `%s`", command)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.SysProcAttr = getSysProcAttrs()
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	// Execute the command.
	err := cmd.Start()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start command")
	}

	// Gather information about the subprocess.
	var fileOpts *FileOptions
	if opts == nil {
		fileOpts = GetDefaultFileOptions()
	} else {
		fileOpts = opts.FileOptions
	}
	subprocess, err := GetProcess(cmd.Process.Pid, fileOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup process")
	}
	pid := subprocess.PID
	ppid := subprocess.PPID

	// Wait for the command to complete.
	log.Infof("Waiting for command to complete: `%s` (PID: %d, PPID: %d)", command, pid, ppid)
	err = cmd.Wait()
	if err != nil {
		log.Errorf("Failed to wait for command to complete: `%s` - %s", command, err)
	}
	exitCode := cmd.ProcessState.ExitCode()
	subprocess.Stdout = stdout.String()
	subprocess.Stderr = stderr.String()
	subprocess.ExitCode = &exitCode

	log.Infof("Command completed: `%s` (PID: %d, PPID: %d, exit code: %d)", command, pid, ppid, exitCode)
	return subprocess, nil
}
