package main

import (
	"bufio"
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
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (t *TelnetCl) Connect() error {
	tcpConn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		slog.Error("couldn't connect", "err", err.Error())
		return err
	}
	t.conn = tcpConn
	return nil
}

func (t *TelnetCl) Close() error {
	defer wg.Done()
	err := t.conn.Close()
	if err != nil {
		slog.Error("error closing connection", "err", err.Error())
		return err
	}
	return nil
}

func (t *TelnetCl) Send() error {
	slog.Info("sending data")
	scanner := bufio.NewScanner(t.in)
	var err error
	for scanner.Scan() {
		_, err = t.conn.Write([]byte(scanner.Text()))
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	}
	return nil
}

func (t *TelnetCl) Receive() error {
	slog.Info("reading connection")
	scanner := bufio.NewScanner(t.conn)
	var err error
	for scanner.Scan() {
		_, err = t.out.Write([]byte(scanner.Text()))
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetCl{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
