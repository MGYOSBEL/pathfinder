package main

import (
	"context"

	// Register topicparser plugin
	_ "github.com/MGYOSBEL/pathfinder/pkg/benthos/processor/topicparser"

	// Import all standard Benthos components
	_ "github.com/redpanda-data/connect/public/bundle/free/v4"

	"github.com/redpanda-data/benthos/v4/public/service"
)

func main() {
	service.RunCLI(context.Background())
}
