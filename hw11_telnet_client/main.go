package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var timeout string

var wg sync.WaitGroup

func init() {
	flag.StringVar(&timeout, "timeout", "15s", "set connection timeout. default value is 15s")
}

func main() {
	flag.Parse()

	if len(os.Args) < 3 || len(os.Args) > 4 {
		slog.Error("incorrect input, use format `go-telnet --timeout=10s host port` (timeout flag is optional)")
		return
	}

	host := os.Args[len(os.Args)-2]
	port := os.Args[len(os.Args)-1]
	addr := host + ":" + port

	t, err := time.ParseDuration(timeout)
	if err != nil {
		slog.Error("couldn't parse the timeout value", "err", err.Error())
		return
	}

	tc := NewTelnetClient(addr, t, os.Stdin, os.Stdout)
	err = tc.Connect()
	if err != nil {
		slog.Error("couldn't connect", "err", err.Error())
		return
	}

	done := make(chan os.Signal, 1)
	go func() {
		signal.Notify(done, syscall.SIGINT)
		<-done
		tc.Close()
	}()

	wg.Add(2)
	go tc.Send()
	go tc.Receive()

	wg.Wait()
}
