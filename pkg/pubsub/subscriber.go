package pubsub

import "context"

type HandlerFunc func(ctx context.Context, message []byte)

type Subscriber interface {
	Subscribe(ctx context.Context, topic string, handler HandlerFunc) error
}
