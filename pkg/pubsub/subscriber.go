package pubsub

type Subscriber interface {
	Subscribe(handler Handler) error
}
