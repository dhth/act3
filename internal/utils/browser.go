package utils

import (
	"os/exec"
)

func OpenURLsInBrowser(urls []string, goos string) error {
	var openCmd string
	switch goos {
	case "darwin":
		openCmd = "open"
	default:
		openCmd = "xdg-open"
	}
	c := exec.Command(openCmd, urls...)
	return c.Run()
}
