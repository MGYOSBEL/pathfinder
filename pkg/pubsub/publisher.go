package pubsub

import "github.com/MGYOSBEL/pathfinder/pkg/message"

type Publisher interface {
	Publish(topic string, m message.Message) error
}
