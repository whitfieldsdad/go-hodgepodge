package main

import (
	"os"
	"strings"
	"time"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// ProcessTree is a tree of processes.
type ProcessTree struct {
	CreateTime time.Time
	UpdateTime time.Time
	processMap map[int]Process // Processes by ID
	pidToPpid  map[int]int     // Map of PID -> PPID
}

// NewProcessTree returns an empty process tree.
func NewProcessTree(processes []Process) *ProcessTree {
	now := time.Now()
	processMap := make(map[int]Process)
	pidMap := make(map[int]int)
	for _, p := range processes {
		processMap[p.Pid] = p
		pidMap[p.Pid] = p.Ppid
	}
	return &ProcessTree{
		CreateTime: now,
		UpdateTime: now,
		processMap: processMap,
		pidToPpid:  pidMap,
	}
}

// GetProcessTree returns a process tree populated with processes.
func GetProcessTree() (*ProcessTree, error) {
	processes, err := ListProcesses(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list processes")
	}
	t := NewProcessTree(processes)
	return t, nil
}

// AddProcess adds a process to the process tree and updates the last update time.
func (t *ProcessTree) AddProcess(pid, ppid int) {
	t.addProcess(pid, ppid)
	t.UpdateTime = time.Now()
}

// addProcess adds a process to the process tree without updating the last update time.
func (t *ProcessTree) addProcess(pid, ppid int) {
	t.pidToPpid[pid] = ppid
}

func (t *ProcessTree) RemoveProcesses(pids ...int) {
	for _, pid := range pids {
		delete(t.pidToPpid, pid)
	}
	t.UpdateTime = time.Now()
}

// GetAncestorPids returns a list of ancestor PIDs for the given PID.
func (t ProcessTree) GetAncestorPids(pid int) []int {
	ancestorPids := []int{}
	for {
		ppid, ok := t.pidToPpid[pid]
		if !ok {
			break
		}
		ancestorPids = append(ancestorPids, ppid)
		pid = ppid
	}
	return ancestorPids
}

// GetDescendantPids returns a list of descendant PIDs for the given PID.
func (t ProcessTree) GetDescendantPids(pid int) []int {
	g := t.ToDiGraph()
	descendants := []int{}
	_ = graph.DFS(g, pid, func(value int) bool {
		descendants = append(descendants, value)
		return false
	})
	return descendants
}

func (t ProcessTree) GetParentPid(pid int) (int, error) {
	p, ok := t.pidToPpid[pid]
	if !ok {
		return -1, errors.Errorf("PID %d not found", pid)
	}
	return p, nil
}

// GetChildProcesses
func (t ProcessTree) GetChildPids(pid int) ([]int, error) {
	pids := []int{}
	for process, parent := range t.pidToPpid {
		if pid == parent {
			pids = append(pids, process)
		}
	}
	return pids, nil
}

// ToDAG returns a DAG organized by PPID -> PID.
func (t ProcessTree) ToDiGraph() graph.Graph[int, int] {
	g := graph.New(graph.IntHash)
	for pid, ppid := range t.pidToPpid {
		_ = g.AddVertex(pid)
		_ = g.AddVertex(ppid)
		_ = g.AddEdge(ppid, pid)
	}
	return g
}

type ProcessHashFunction func(Process) string

func (t ProcessTree) ToDiGraphWithLabels(vertexLabel string) (*graph.Graph[int, int], error) {
	g := graph.New(graph.IntHash)
	for pid, ppid := range t.pidToPpid {

		processLabel, err := t.getProcessLabelByPid(pid, vertexLabel)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get process label")
		}
		parentProcessLabel, err := t.getProcessLabelByPid(ppid, vertexLabel)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get parent process label")
		}
		_ = g.AddVertex(pid, graph.VertexAttribute("label", processLabel))
		_ = g.AddVertex(ppid, graph.VertexAttribute("label", parentProcessLabel))
		_ = g.AddEdge(ppid, pid)
	}
	return &g, nil
}

func (t ProcessTree) getProcessLabelByPid(pid int, key string) (string, error) {
	process, ok := t.processMap[pid]
	if !ok {
		return "", errors.Errorf("PID %d not found in process map", pid)
	}
	var processData map[string]interface{}
	err := mapstructure.Decode(process, &processData)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode process data")
	}
	label := processData[key].(string)
	if label == "" {
		return "", errors.Errorf("label cannot be empty")
	}
	return label, nil
}

// RenderProcessTree saves the process tree to a DOT file.
func (t ProcessTree) DrawProcessTree(path string) error {
	if strings.HasSuffix(path, ".dot") || strings.HasSuffix(path, ".gv") {
		g := t.ToDiGraph()
		file, _ := os.Create(path)
		_ = draw.DOT(g, file)
		return nil
	}
	return errors.New("invalid file extension")
}
