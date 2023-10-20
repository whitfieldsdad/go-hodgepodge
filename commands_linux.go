//go:build !windows && !js && !darwin
// +build !windows,!js,!darwin

package main

import "syscall"

func getSysProcAttrs() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}
