package cmd

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
)

func Main() {
	c := cli.NewCLI("hodgepodge", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"processes list": func() (cli.Command, error) {
			return &listProcesses{}, nil
		},
	}
	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}
	os.Exit(exitStatus)
}
