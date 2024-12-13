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
	//envSlice := make([]string, 0, len(env))
	//for key, e := range env {
	//	envSlice = append(envSlice, key+"="+e.Value)
	//}
	//c.Env = append(c.Env, envSlice...)
	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
			continue
		}
		os.Setenv(k, v.Value)
	}
	err := c.Run()
	//fmt.Printf("command output:\n%s", stdout.String())
	if err != nil {
		//fmt.Printf("error output:\n%s", stderr.String())
		slog.Error(err.Error())
	}
	return c.ProcessState.ExitCode()
}
