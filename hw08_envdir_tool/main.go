package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal("Incorrect input\ngo-envdir PATH COMMAND arg1 arg2...\n")
	}

	path := args[0]
	cmd := args[1:]
	env, err := ReadDir(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	returnCode := RunCmd(cmd, env)
	os.Exit(returnCode)
}
