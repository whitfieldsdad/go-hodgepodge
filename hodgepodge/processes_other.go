//go:build !windows
// +build !windows

package hodgepodge

import "os"

func currentProcessIsElevated() (bool, error) {
	return os.Geteuid() == 0, nil
}
