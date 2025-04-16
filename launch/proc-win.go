//go:build windows
// +build windows

package launch

import (
	"os/exec"
)

func setProcessGroupID(cmd *exec.Cmd) {
	// Windows doesn't support process groups in the same way
}

func killProcessGroup(cmd *exec.Cmd) {
	cmd.Process.Kill()
}
