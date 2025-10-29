package pubsub

import "github.com/MGYOSBEL/pathfinder/pkg/message"

type Handler func(m message.Message)
