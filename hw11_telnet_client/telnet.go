package main

import (
	"bufio"
	"context"
	"io"
	"log/slog"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetCl struct {
	address string
	timeout time.Duration
	in      chan string
	out     io.Writer
	conn    net.Conn
	ctx     context.Context
	ctxCanc func()
}

func (t *TelnetCl) Connect() error {
	tcpConn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	t.conn = tcpConn
	return nil
}

func (t *TelnetCl) Close() error {
	t.ctxCanc()
	err := t.conn.Close()
	if err != nil {
		slog.Error("error closing connection", "err", err.Error())
		return err
	}
	return nil
}

func (t *TelnetCl) Send() error {
	defer wg.Done()
	for {
		select {
		case <-t.ctx.Done():
			return nil
		case data, ok := <-t.in:
			if !ok {
				t.Close()
				return nil
			}
			_, err := t.conn.Write([]byte(data))
			if err != nil {
				slog.Error("error writing to socket", "err", err.Error())
				t.Close()
				return err
			}
		}
	}
}

func (t *TelnetCl) Receive() error {
	defer wg.Done()
	scanner := bufio.NewScanner(t.conn)
	for {
		select {
		case <-t.ctx.Done():
			return nil
		default:
			if !scanner.Scan() {
				return scanner.Err()
			}
			t.out.Write([]byte(scanner.Text() + "\n"))
		}
	}
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	ctx, fn := context.WithCancel(context.Background())
	return &TelnetCl{
		address: address,
		timeout: timeout,
		in:      stdin(in, address),
		out:     out,
		ctx:     ctx,
		ctxCanc: fn,
	}
}

func stdin(input io.ReadCloser, addr string) chan string {
	in := make(chan string, 1)
	go func() {
		defer close(in)
		sc := bufio.NewScanner(input)
		for sc.Scan() {
			in <- sc.Text()
		}
		if !sc.Scan() && sc.Err() == nil {
			slog.Info("terminating connection to " + addr)
		}
	}()
	return in
}
