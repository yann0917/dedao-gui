package utils

import (
	"os/exec"
	"strconv"
)

func configurePDFProcess(cmd *exec.Cmd) {
	cmd.Cancel = func() error {
		// Chrome has helper processes; terminate the isolated process tree too.
		if err := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid)).Run(); err != nil {
			return cmd.Process.Kill()
		}
		return nil
	}
}
