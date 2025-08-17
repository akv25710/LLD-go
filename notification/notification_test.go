package notification

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestNotification(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	topic := "notification"
	manager := NewManager()
	sender := NewSender()
	storage := NewStorage()
	consumer := NewDispatcher(sender, storage)

	sub := manager.Subscribe(ctx, topic, 10)
	go consumer.Start(ctx, sub)

	for i := 1; i < 10; i++ {
		n := Notification{
			Id:          uuid.New().String(),
			Message:     "Welcome back!",
			DeviceToken: "abc123_device_token",
		}
		go func(n Notification) {
			raw, _ := json.Marshal(n)
			manager.Publish(ctx, topic, raw)

			time.Sleep(5 * time.Millisecond)
			if record, ok := storage.GetStatus(n.Id); ok {
				log.Printf("📊 Status for %s: %s\n", record.NotificationID, record.Status)
			}
		}(n)
	}

	time.Sleep(100 * time.Second)
}
