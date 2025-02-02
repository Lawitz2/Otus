package main

import (
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var timeout time.Duration

var wg sync.WaitGroup

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "set connection timeout. default value is 15s")
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 2 {
		slog.Error("incorrect input, use format `go-telnet --timeout=10s host port` (timeout flag is optional)")
		return
	}

	addr := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	tc := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)
	err := tc.Connect()
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
