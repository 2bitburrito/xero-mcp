package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenURL opens the specified URL in the default web browser of the user
// using OS-specific commands.
func OpenURL(url string) error {
	fmt.Println("Visiting this URL to authorize: ", url)
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin": // macOS
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open" // Common on many Linux distributions
		args = []string{url}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return exec.Command(cmd, args...).Start()
}
