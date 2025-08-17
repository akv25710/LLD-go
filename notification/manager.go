package notification

import (
	"context"
	"time"
)

type Manager interface {
	Publish(ctx context.Context, topic string, data []byte)
	Subscribe(ctx context.Context, topic string, buffer int) <-chan Message
}

func NewManager() Manager {
	return NewPubSubManager()
}

type PubSubManager struct {
	topics map[string][]chan Message
}

func NewPubSubManager() PubSubManager {
	return PubSubManager{
		topics: make(map[string][]chan Message),
	}
}

func (p PubSubManager) Publish(ctx context.Context, topic string, data []byte) {
	for _, ch := range p.topics[topic] {
		select {
		case <-ctx.Done():
			return
		case ch <- Message{Topic: topic, Data: data}:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (p PubSubManager) Subscribe(ctx context.Context, topic string, buffer int) <-chan Message {
	ch := make(chan Message, buffer)
	p.topics[topic] = append(p.topics[topic], ch)
	return ch
}
