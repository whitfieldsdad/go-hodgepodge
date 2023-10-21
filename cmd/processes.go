package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/whitfieldsdad/go-hodgepodge/hodgepodge"
)

type listProcesses struct{}

func (*listProcesses) Help() string {
	return "List processes"
}

func (*listProcesses) Run(args []string) int {
	opts := &hodgepodge.FileOptions{
		IncludeFileHashes:     true,
		IncludeFileTraits:     true,
		IncludeFileTimestamps: true,
	}
	processes, err := hodgepodge.ListProcesses(opts)
	if err != nil {
		log.Errorf("Failed to list processes: %s", err)
		return 1
	}
	for _, process := range processes {
		blob, err := json.Marshal(process)
		if err != nil {
			log.Errorf("Failed to marshal process: %s", err)
			return 1
		}
		fmt.Println(string(blob))
	}
	return 0
}

func (o *listProcesses) Synopsis() string {
	return o.Help()
}
