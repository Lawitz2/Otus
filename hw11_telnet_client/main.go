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
	go tcpserver()

	tc := NewTelnetClient("localhost:8103", time.Second*15, os.Stdin, os.Stdout)
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

func tcpserver() {
	listener, err := net.Listen("tcp", "localhost:8103")
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
				conn.Write([]byte("your message was received by the server\n"))
				slog.Info("server received " + strconv.Itoa(n) + " bytes")
				slog.Debug(string(buf))
				if err != nil {
					slog.Error("server error", "err", err.Error())
					return
				}
				time.Sleep(time.Millisecond * 100)
				conn.Close()
			}
		}()
	}
}
