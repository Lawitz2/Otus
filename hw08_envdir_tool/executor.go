package main

import (
	"log/slog"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	returnCode = 1
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Env = os.Environ()
	envSlice := make([]string, 0, len(env))
	for key, e := range env {
		envSlice = append(envSlice, key+"="+e.Value)
	}
	c.Env = append(c.Env, envSlice...)
	err := c.Start()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	c.Wait()
	returnCode = 0
	return
}
