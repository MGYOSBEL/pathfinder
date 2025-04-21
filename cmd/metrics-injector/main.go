package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/MGYOSBEL/pathfinder/internal/mqtt"
)

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
	client := mqtt.NewClient()
	if err := client.Connect(); err != nil {
		panic(err)
	}
	client.Subscribe()

	// Wait for cancel
	<-ctx.Done()
	return nil
}
