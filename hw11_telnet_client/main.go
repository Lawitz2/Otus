package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	wg.Add(1)
	go tcpserver()

	tc := NewTelnetClient("localhost:81", time.Second*15, os.Stdin, os.Stdout)
	wg.Add(1)
	tc.Connect()
	go tc.Send()
	go tc.Receive()

	done := make(chan os.Signal, 1)
	go func() {
		signal.Notify(done, syscall.SIGINT)
		<-done
		slog.Info("received interrupt signal")
		tc.Close()
		os.Exit(0)
	}()

	wg.Wait()
	//time.Sleep(time.Second * 2)
}

func tcpserver() {
	defer wg.Done()

	listener, err := net.Listen("tcp", `localhost:81`)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error(err.Error())
			return
		}
		go func() {
			slog.Info("received connection", "source", conn.RemoteAddr())
			buf := make([]byte, 4096)
			for {
				n, err := conn.Read(buf)
				slog.Info("server received " + strconv.Itoa(n) + " bytes")
				slog.Debug(string(buf))
				if err != nil {
					slog.Error("server error", "err", err.Error())
					return
				}
				time.Sleep(time.Second / 10)
				conn.Close()
			}
		}()
	}
}
