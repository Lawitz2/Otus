package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeOffset        = errors.New("negative offset value")
	ErrNegativeLimit         = errors.New("negative limit value")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileDescr, err := os.Stat(from) // getting filesize.
	if err != nil {
		return ErrUnsupportedFile
	}

	// check for input validity
	if fileDescr.Size() < offset {
		return ErrOffsetExceedsFileSize
	}
	if offset < 0 {
		return ErrNegativeOffset
	}
	if limit < 0 {
		return ErrNegativeLimit
	}

	targetFile, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer targetFile.Close()

	copyFile, err := os.Create(toPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer copyFile.Close()

	_, err = targetFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(targetFile)
	writer := bufio.NewWriter(copyFile)

	// the following math allows me to have nil error after successful copy.
	var size int64
	if limit == 0 {
		size = fileDescr.Size() - offset
	} else {
		size = min(limit, fileDescr.Size()-offset)
	}

	prBar := newProgressBar(size)

	// Copying contents of original file into a multiwriter.
	// One part of the multiwriter is writing data into the new file.
	// The other one feeds data into progress bar for progress tracking.
	_, err = io.CopyN(io.MultiWriter(writer, prBar), reader, size)
	if err != nil {
		return err
	}

	writer.Flush()

	fmt.Println("\ndone!")

	return err
}

// progress bar struct
// Progress represents how much was already copied
// Capacity represents total amount that needs to be copied.
type pBar struct {
	Progress int64
	Capacity int64
}

func newProgressBar(capacity int64) *pBar {
	return &pBar{
		Progress: 0,
		Capacity: capacity,
	}
}

// implementation of io.Writer interface
// shows "bytes copied/bytes total", a progress bar with 20 segments (each representing 5%)
// and the percentage of copying progress.
func (p *pBar) Write(b []byte) (n int, err error) {
	p.Progress += int64(len(b))
	done := int(20 * p.Progress / p.Capacity)
	notDone := 20 - done
	percentage := 100 * p.Progress / p.Capacity
	//nolint:lll
	fmt.Printf("\rcopying %d/%d [%s%s], %d%%", p.Progress, p.Capacity, strings.Repeat("█", done), strings.Repeat("-", notDone), percentage)
	return len(b), nil
}
