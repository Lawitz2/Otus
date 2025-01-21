package main

import (
	"bytes"
	"log/slog"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) int {
	var stderr, stdout bytes.Buffer
	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	c.Env = os.Environ()
	c.Stderr = &stderr
	c.Stdout = &stdout
	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
			continue
		}
		os.Setenv(k, v.Value)
	}
	err := c.Run()
	if err != nil {
		slog.Error(err.Error())
	}
	return c.ProcessState.ExitCode()
}
