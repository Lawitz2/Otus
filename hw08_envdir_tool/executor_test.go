package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	// successful execution should return exit code 0
	t.Run("bash", func(t *testing.T) {
		cmd := []string{"bash"}
		env, _ := ReadDir(`.\testdata\env\`)
		exitCode := RunCmd(cmd, env)
		require.Equal(t, 0, exitCode)
	})

	// exec file not found should return exit code -1
	t.Run("exec file not found", func(t *testing.T) {
		cmd := []string{"tasklist_bad_name"}
		env, _ := ReadDir(".\\testdata\\env")
		exitCode := RunCmd(cmd, env)
		require.Equal(t, -1, exitCode)
	})
}
