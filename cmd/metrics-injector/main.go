package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/MGYOSBEL/pathfinder/internal/processor"
	"github.com/MGYOSBEL/pathfinder/pkg/mqtt"
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
		Server: "localhost:1883",
		Topic:  "home/#",
		QoS:    0,
	}
	client := mqtt.NewClient(opts)
	if err := client.Connect(); err != nil {
		panic(err)
	}
	defer client.Disconnect()

	p := processor.New(client, client, processor.Options{
		InputTopic:  "home/#",
		OutputTopic: "forwarded",
	})
	p.Process()

	// Wait for cancel
	<-ctx.Done()
	return nil
}
