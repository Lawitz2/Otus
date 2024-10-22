package main

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint:goconst
func TestCopy(t *testing.T) {
	from = ".\\testdata\\input.txt"
	to = ".\\testdata\\output_test.txt"
	inInfo, err := os.Stat(from)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	t.Run("just a file copy", func(t *testing.T) {
		err = Copy(from, to, 0, 0)
		require.Nil(t, err)

		outInfo, err := os.Stat(to)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		require.Equal(t, outInfo.Size(), inInfo.Size())
		os.Remove(to)
	})

	t.Run("offset and limit", func(t *testing.T) {
		offset = 150
		limit = 500
		inInfo, err = os.Stat(from)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		err = Copy(from, to, offset, 0)
		require.Nil(t, err)

		outInfo, err := os.Stat(to)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		require.Equal(t, outInfo.Size(), inInfo.Size()-offset)
		os.Remove(to)

		err = Copy(from, to, offset, limit)
		require.Nil(t, err)

		outInfo, _ = os.Stat(to)
		require.Equal(t, outInfo.Size(), min(limit, inInfo.Size()-offset))
		os.Remove(to)

		limit = inInfo.Size() * 2 // limit higher than file size
		err = Copy(from, to, offset, limit)
		require.Nil(t, err)

		outInfo, _ = os.Stat(to)
		require.Equal(t, outInfo.Size(), min(limit, inInfo.Size()-offset))
		os.Remove(to)
	})

	t.Run("comparison to pre-made test files", func(t *testing.T) {
		testA := ".\\testdata\\out_offset0_limit0.txt"
		testB := ".\\testdata\\out_offset0_limit10.txt"
		testC := ".\\testdata\\out_offset0_limit1000.txt"
		testD := ".\\testdata\\out_offset0_limit10000.txt"
		testE := ".\\testdata\\out_offset100_limit1000.txt"
		testF := ".\\testdata\\out_offset6000_limit1000.txt"

		err = Copy(from, to, 0, 0) // to out_offset0_limit0.txt
		require.Nil(t, err)
		outInfo, _ := os.Stat(from)
		exampleInfo, _ := os.Stat(testA)
		os.SameFile(outInfo, exampleInfo)
		os.Remove(to)

		err = Copy(from, to, 0, 10) // to out_offset0_limit0.txt
		require.Nil(t, err)
		outInfo, _ = os.Stat(from)
		exampleInfo, _ = os.Stat(testB)
		os.SameFile(outInfo, exampleInfo)
		os.Remove(to)

		err = Copy(from, to, 0, 1000) // to out_offset0_limit1000.txt
		require.Nil(t, err)
		outInfo, _ = os.Stat(from)
		exampleInfo, _ = os.Stat(testC)
		os.SameFile(outInfo, exampleInfo)
		os.Remove(to)

		err = Copy(from, to, 0, 10000) // to out_offset0_limit10000.txt
		require.Nil(t, err)
		outInfo, _ = os.Stat(from)
		exampleInfo, _ = os.Stat(testD)
		os.SameFile(outInfo, exampleInfo)
		os.Remove(to)

		err = Copy(from, to, 0, 10) // to out_offset100_limit1000.txt
		require.Nil(t, err)
		outInfo, _ = os.Stat(from)
		exampleInfo, _ = os.Stat(testE)
		os.SameFile(outInfo, exampleInfo)
		os.Remove(to)

		err = Copy(from, to, 0, 10) // to out_offset6000_limit1000.txt
		require.Nil(t, err)
		outInfo, _ = os.Stat(from)
		exampleInfo, _ = os.Stat(testF)
		os.SameFile(outInfo, exampleInfo)
		os.Remove(to)
	})

	t.Run("incorrect requests", func(t *testing.T) {
		offset = 150
		limit = 500

		err = Copy(from+"randomstuff", to, offset, limit) // trying to copy a non-existing file
		require.Equal(t, err, ErrUnsupportedFile)

		err = Copy(from, to, inInfo.Size()*2, limit) // offset higher than the whole file
		require.Equal(t, err, ErrOffsetExceedsFileSize)

		err = Copy(from, to, -5, limit) // negative offset value
		require.Equal(t, err, ErrNegativeOffset)

		err = Copy(from, to, offset, -400) // negative limit value
		require.Equal(t, err, ErrNegativeLimit)

		// confirm that none of the incorrect requests created a new file
		_, err = os.Stat(to)
		require.True(t, os.IsNotExist(err))
	})
}
