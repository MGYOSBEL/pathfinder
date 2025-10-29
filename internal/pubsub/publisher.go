package pubsub

type Publisher interface {
	Publish(topic string, message any) error
}
