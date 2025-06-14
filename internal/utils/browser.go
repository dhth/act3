package utils

import (
	"os/exec"
	"runtime"
)

func OpenURLsInBrowser(urls []string) error {
	var openCmd string
	switch runtime.GOOS {
	case "darwin":
		openCmd = "open"
	default:
		openCmd = "xdg-open"
	}
	c := exec.Command(openCmd, urls...)
	return c.Run()
}
