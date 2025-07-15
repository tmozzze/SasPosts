package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

type PubSub interface {
	Publish(ctx context.Context, channel string, message interface{}) error
	Subscribe(ctx context.Context, channel string) (<-chan []byte, func())
}

type redisPubSub struct {
	client *redis.Client
}

func NewPubSub(client *redis.Client) PubSub {
	return &redisPubSub{client: client}
}

func (p *redisPubSub) Publish(ctx context.Context, channel string, message interface{}) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	if err := p.client.Publish(ctx, channel, payload).Err(); err != nil {
		fmt.Printf("failed to publish message to redis: %v\n", err)
	}
	return nil
}

func (p *redisPubSub) Subscribe(ctx context.Context, channel string) (<-chan []byte, func()) {
	pubsub := p.client.Subscribe(ctx, channel)

	ch := make(chan []byte)
	var once sync.Once

	closeFunc := func() {
		once.Do(func() {
			pubsub.Close()
			close(ch)
		})
	}
	go func() {
		<-ctx.Done()
		closeFunc()
	}()

	go func() {
		defer closeFunc()
		for msg := range pubsub.Channel() {
			ch <- []byte(msg.Payload)
		}
	}()

	return ch, closeFunc
}
