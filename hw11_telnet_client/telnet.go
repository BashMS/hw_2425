package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error

	io.Closer

	Send() error

	Receive() error
}

// TClient - структура для хранения состояния клиента + имплиментируем методы основного интерфейса

// По другому не придумал как адрес, таймаут и коннект передавть между методами.

type TClient struct {
	in *io.ReadCloser

	out *io.Writer

	address string

	timeout time.Duration

	conn *net.Conn
}

func (t *TClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	t.conn = &conn
	return nil
}

func (t *TClient) Close() error {
	return (*t.conn).Close()
}

func (t *TClient) Send() error {
	_, err := io.Copy(*t.conn, *t.in)
	return err
}

func (t *TClient) Receive() error {
	_, err := io.Copy(*t.out, *t.conn)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	t := TClient{
		in:      &in,
		out:     &out,
		address: address,
		timeout: timeout,
	}
	return &t
}
