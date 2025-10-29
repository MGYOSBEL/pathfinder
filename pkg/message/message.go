package message

type Message struct {
	Payload []byte
	Metadata
}

type Metadata struct {
	Topic string
}
