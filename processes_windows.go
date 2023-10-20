package main

import (
	"github.com/charmbracelet/log"
	"golang.org/x/sys/windows"
)

// currentProcessIsElevated checks to see if the current process is either running with elevated privileges, or was started by an administrative user.
func currentProcessIsElevated() bool {
	log.Info("Checking if process is running with elevated privileges and/or was created by an administrative user")

	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)

	if err != nil {
		log.Errorf("Failed to allocate and initialize SID: %v", err)
		return false
	}
	token := windows.Token(0)
	defer token.Close()

	member, err := token.IsMember(sid)
	if err != nil {
		log.Errorf("Failed to check token membership: %v", err)
		return false
	}
	if member {
		log.Info("Process was started by an administrative user")
	}
	elevated := token.IsElevated()
	if elevated {
		log.Info("Process is running with elevated privileges")
	}
	return member || elevated
}
