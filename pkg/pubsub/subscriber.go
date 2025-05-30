package pubsub

import "context"

type HandlerFunc func(ctx context.Context, message []byte)
type CloseFunc func() error

type Subscriber interface {
	Subscribe(ctx context.Context, topic string, handler HandlerFunc) (CloseFunc, error)
}
