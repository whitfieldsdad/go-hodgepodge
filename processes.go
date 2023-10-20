package main

import (
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
	"github.com/pkg/errors"
)

type Process struct {
	Name             string     `json:"name,omitempty"`
	WorkingDirectory string     `json:"working_directory,omitempty"`
	PID              int        `json:"pid,omitempty"`
	PPID             int        `json:"ppid,omitempty"`
	File             File       `json:"file,omitempty"`
	CommandLine      string     `json:"command_line,omitempty"`
	Argv             []string   `json:"argv,omitempty"`
	Argc             int        `json:"argc,omitempty"`
	Stdout           string     `json:"stdout,omitempty"`
	Stderr           string     `json:"stderr,omitempty"`
	ExitCode         *int       `json:"exit_code,omitempty"`
	StartTime        *time.Time `json:"start_time,omitempty"`
	ExitTime         *time.Time `json:"exit_time,omitempty"`
}

// GetParentProcessId returns the PPID of a process.
func GetParentProcessId(pid int) (int, error) {
	p, err := sysinfo.Process(pid)
	if err != nil {
		return -1, errors.Wrap(err, "failed to lookup process")
	}
	info, err := p.Info()
	if err != nil {
		return -1, errors.Wrap(err, "failed to get process info")
	}
	return info.PPID, nil
}

// GetProcess looks up a process by PID.
func GetProcess(pid int, opts *FileOptions) (*Process, error) {
	p, err := sysinfo.Process(pid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup process")
	}
	process, err := getProcess(p, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get process")
	}
	return process, nil
}

// ListProcesses returns a list of processes.
func ListProcesses(opts *FileOptions) ([]Process, error) {
	log.Infof("Listing processes")
	processes, err := sysinfo.Processes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list processes")
	}
	var rows []Process
	for _, process := range processes {
		row, err := getProcess(process, opts)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get process")
		}
		rows = append(rows, *row)
	}
	log.Infof("Found %d processes", len(rows))
	return rows, nil
}

func getProcess(p types.Process, opts *FileOptions) (*Process, error) {
	info, err := p.Info()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get process info")
	}
	process := &Process{
		Name:             info.Name,
		WorkingDirectory: info.CWD,
		PID:              info.PID,
		PPID:             info.PPID,
		CommandLine:      strings.Join(info.Args, " "),
		Argv:             info.Args,
		Argc:             len(info.Args),
		StartTime:        &info.StartTime,
	}
	path := info.Exe
	file, err := GetFile(path, opts)
	if err != nil {
		log.Warnf("Failed to get file metadata: %s (path: %s)", err, path)
	} else {
		process.File = *file
	}
	return process, nil
}

// GetPidMap returns a map of PID -> PPID.
func GetPidMap() (map[int]int, error) {
	processes, err := sysinfo.Processes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list processes")
	}
	var m = make(map[int]int)
	for _, proc := range processes {
		pid := proc.PID()
		parent, err := proc.Parent()
		if err != nil {
			log.Warnf("Failed to lookup PPID while building PID map: %s (PID: %d)", err, pid)
			continue
		}
		ppid := parent.PID()
		m[pid] = ppid
	}
	return m, nil
}

func CurrentProcessIsElevated() (bool, error) {
	return currentProcessIsElevated()
}

func KillProcess(pid int) error {
	log.Infof("Killing process: %d", pid)
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = p.Kill()
	if err != nil {
		return err
	}
	log.Infof("Killed process: %d", pid)
	return nil
}
