package pubsub

import "github.com/MGYOSBEL/pathfinder/internal/message"

type Handler func(m message.Message)
