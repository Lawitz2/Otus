package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("correct path", func(t *testing.T) {
		// folder contains 5 env files, one of which is empty and should be excluded
		// end result is expected to be 4 env variables in a map with no error
		t.Parallel()
		input := "./testdata/env"
		dir, err := ReadDir(input)
		require.Nil(t, err)
		require.Equal(t, 4, len(dir))
		require.Equal(t, dir["HELLO"].Value, `"hello"`)
		require.Equal(t, dir["UNSET"].Value, "")
	})

	t.Run("bad path", func(t *testing.T) {
		// bad path test
		// end result is expected to be nil, with a non-nil error value
		t.Parallel()
		dir, err := ReadDir(".\\badInput")
		require.Error(t, err)
		require.Nil(t, dir)
	})
}
