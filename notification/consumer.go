package notification

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

type Consumer interface {
	Start(ctx context.Context, sub <-chan Message)
}

type Dispatcher struct {
	sender Sender
	store  Storage
}

func NewDispatcher(sender Sender, store Storage) *Dispatcher {
	return &Dispatcher{
		sender: sender,
		store:  store,
	}
}

func (d *Dispatcher) Start(ctx context.Context, sub <-chan Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-sub:
			var n interface{}
			if err := json.Unmarshal(msg.Data, &n); err != nil {
				log.Println("❌ Unmarshal error:", err)
				continue
			}
			d.store.Save(msg.Id, StatusRecord{
				NotificationID: msg.Id,
				Status:         StatusSending,
				UpdatedAt:      time.Now(),
			})
			if res := d.sender.Send(ctx, n); !res.Success {
				log.Println("<UNK> Sender error:", res)
				d.store.Save(msg.Id, StatusRecord{
					NotificationID: msg.Id,
					Status:         StatusFailed,
					UpdatedAt:      time.Now(),
				})
			} else {
				log.Println("<UNK> Sender success:", res)
				d.store.Save(msg.Id, StatusRecord{
					NotificationID: msg.Id,
					Status:         StatusSent,
					UpdatedAt:      time.Now(),
					Error:          res.Err,
				})
			}

		}
	}
}
