package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/MGYOSBEL/pathfinder/internal/mqtt"
)

// TODO:
// - Start with the service config
// - Read it from DataPlatform to extract the specs
// - Can it be improved some way???
//
//TODO: Start with the tests

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	fmt.Println("Hello Mars!!!")
	opts := mqtt.Options{
		Server: "192.168.0.16:1883",
		Topic:  "#",
		QoS:    0,
	}
	client := mqtt.NewClient(opts)
	if err := client.Connect(); err != nil {
		panic(err)
	}
	defer client.Disconnect()

	client.Subscribe(func(message []byte) {
		fmt.Printf("%s\n", message)
	})

	// Wait for cancel
	<-ctx.Done()
	return nil
}
