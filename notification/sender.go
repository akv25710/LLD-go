package notification

import (
	"context"
	"time"
)

type Sender interface {
	Send(ctx context.Context, data interface{}) Result
}

func NewSender() Sender {
	return NewFirebaseSender()
}

type FirebaseSender struct {
}

func NewFirebaseSender() *FirebaseSender {
	return &FirebaseSender{}
}

func (p *FirebaseSender) Send(ctx context.Context, data interface{}) Result {
	start := time.Now()
	time.Sleep(50 * time.Millisecond)
	return Result{
		Err:     nil,
		Latency: time.Since(start),
		Success: true,
	}
}
