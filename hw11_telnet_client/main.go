package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "timeout for connection")

	flag.Parse()

	args := flag.Args()

	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet --timeout=10s host port")

		os.Exit(1)
	}

	address := args[0]
	if strings.HasSuffix(address, ":") {
		address = fmt.Sprintf("%s%s", address, args[1])
	} else {
		address = fmt.Sprintf("%s:%s", address, args[1])
	}

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to server:", err)

		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)

	defer cancel()

	ctx, stop := context.WithTimeout(ctx, *timeout)

	defer stop()

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()

		writeRoutine(ctx, client)
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		readRoutine(ctx, client)
	}()

	wg.Wait()

	client.Close()
}

func readRoutine(ctx context.Context, client TelnetClient) {
	for {
		select {
		case <-ctx.Done():

			return
		default:

			err := client.Receive()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error Receive:", err)
				ctx.Done()
			}
		}
	}
}

func writeRoutine(ctx context.Context, client TelnetClient) {
	for {
		select {
		case <-ctx.Done():

			return
		default:

			err := client.Send()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error Send:", err)
				ctx.Done()
			}
		}
	}
}
