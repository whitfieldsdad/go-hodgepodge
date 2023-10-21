//go:build !windows

package hodgepodge

import "path/filepath"

func isHidden(path string) (bool, error) {
	filename := filepath.Base(path)
	return filename[0] == '.', nil
}
