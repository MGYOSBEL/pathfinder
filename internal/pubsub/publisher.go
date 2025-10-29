package pubsub

import "github.com/MGYOSBEL/pathfinder/internal/message"

type Publisher interface {
	Publish(topic string, m message.Message) error
}
