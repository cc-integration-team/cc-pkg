package pubsub

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type redisPublisher struct {
	client *redis.Client
}

func NewRedisPublisher(client *redis.Client) *redisPublisher {
	return &redisPublisher{
		client: client,
	}
}

func (p *redisPublisher) Publish(ctx context.Context, topic string, message []byte) error {
	return p.client.Publish(ctx, topic, message).Err()
}

type redisSubscriber struct {
	client *redis.Client
}

func NewRedisSubscriber(client *redis.Client) *redisSubscriber {
	return &redisSubscriber{
		client: client,
	}
}

func (s *redisSubscriber) Subscribe(ctx context.Context, topic string, handler HandlerFunc) error {
	pubsub := s.client.Subscribe(ctx, topic)

	// Ping the subscription to ensure it's ready
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return err
	}

	// Start a goroutine to handle messages
	go func() {
		defer pubsub.Close()
		select {
		case msg, ok := <-pubsub.Channel():
			if !ok {
				return
			}
			handler(ctx, []byte(msg.Payload))
		case <-ctx.Done():
			return
		}
	}()

	return nil
}
