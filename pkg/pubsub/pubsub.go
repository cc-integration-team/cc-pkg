package pubsub

type Publisher interface {
	Publish(topic string, message []byte) error
}

type Subscriber interface {
	Subscribe(topic string, handler func(message []byte)) error
}
