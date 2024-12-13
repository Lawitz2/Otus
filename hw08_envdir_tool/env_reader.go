package main

import (
	"bufio"
	"bytes"
	"log/slog"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
//func ReadDir(dirPath string) (Environment, error) {
//	dir, err := os.ReadDir(dirPath)
//	if err != nil {
//		slog.Error(err.Error())
//		return nil, err
//	}
//	os.Chdir(dirPath)
//
//	env := make(Environment, len(dir))
//	for _, entry := range dir {
//		ev := EnvValue{}
//		f, _ := os.Open(entry.Name())
//		if strings.Contains(f.Name(), "=") {
//			continue
//		}
//
//		fStat, _ := f.Stat()
//		if fStat.Size() == 0 {
//			ev.NeedRemove = true
//		}
//
//		scanner := bufio.NewScanner(f)
//		scanner.Scan()
//
//		ev.Value = string(bytes.ReplaceAll([]byte(scanner.Text()), []byte{0}, []byte{10}))
//		ev.Value = strings.TrimRight(ev.Value, "\t ")
//
//		env[f.Name()] = ev
//		f.Close()
//	}
//
//	//for key, val := range env {
//	//	if val.NeedRemove {
//	//		delete(env, key)
//	//	}
//	//}
//
//	return env, nil
//}

func ReadDir(dir string) (Environment, error) {
	fsDir, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(dir))

	for _, entry := range fsDir {
		if strings.Contains(entry.Name(), "=") {
			continue // "=" не должно быть в имени файла
		}

		f, err := os.Open(dir + "/" + entry.Name())
		if err != nil {
			slog.Error("Не могу открыть файл", "error", err)
			continue
		}

		ev := EnvValue{}

		fStat, _ := f.Stat()
		if fStat.Size() == 0 {
			ev.NeedRemove = true
			env[entry.Name()] = ev
			continue
		}

		scanner := bufio.NewScanner(f)
		scanner.Scan()

		ev.Value = string(bytes.ReplaceAll([]byte(scanner.Text()), []byte{0}, []byte{10}))
		ev.Value = strings.TrimRight(ev.Value, "\t ")

		env[entry.Name()] = ev
	}

	return env, nil
}
