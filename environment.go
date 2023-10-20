package main

import (
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"golang.org/x/exp/slices"
)

type EnvironmentVariables struct {
	Env map[string]string `json:"env"`
}

func GetEnvironmentVariables() map[string]string {
	log.Info("Listing environment variables")
	env := map[string]string{}
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		k := pair[0]
		v := pair[1]
		env[k] = v
	}
	keys := []string{}
	for k := range env {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	log.Infof("Found %d environment variables: %s", len(keys), strings.Join(keys, ", "))
	return env
}
