package main

import (
	"context"
	"fmt"
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
	ctx context.Context

	in *io.ReadCloser

	out *io.Writer

	address string

	timeout time.Duration

	conn *net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Создаем контекст с таймаутом и отменой

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	t := TClient{
		ctx: ctx,

		in: &in,

		out: &out,

		address: address,

		timeout: timeout,
	}

	return &t
}

func (t *TClient) Connect() error {
	dialer := &net.Dialer{}

	conn, err := dialer.DialContext(t.ctx, "tcp", t.address)
	if err != nil {
		return fmt.Errorf("dialer.DialContext: %w", err)
	}

	defer conn.Close()

	t.conn = &conn

	return nil
}

func (t *TClient) Close() error {
	err := (*t.conn).Close()
	if err != nil {
		return fmt.Errorf("conn.Close: %w", err)
	}

	return nil
}

func (t *TClient) Send() error {
	_, err := io.Copy(*t.conn, *t.in)
	if err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	return nil
}

func (t *TClient) Receive() error {
	_, err := io.Copy(*t.out, *t.conn)
	if err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	return nil
}
