package pubsub

type Subscriber interface {
	Subscribe(topic string, handler Handler) error
}
