package utils

import (
	"path/filepath"
	"strings"
)

func ExpandTilde(path string, homeDir string) string {
	pathWithoutTilde, found := strings.CutPrefix(path, "~/")
	if !found {
		return path
	}
	return filepath.Join(homeDir, pathWithoutTilde)
}
