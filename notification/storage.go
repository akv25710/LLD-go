package notification

import (
	"sync"
	"time"
)

type StatusRecord struct {
	NotificationID string
	Status         Status
	UpdatedAt      time.Time
	Error          error
}

type Storage interface {
	Save(id string, record StatusRecord)
	GetStatus(notificationID string) (StatusRecord, bool)
}

func NewStorage() Storage {
	return &InMemoryStorage{
		mu:   sync.RWMutex{},
		data: map[string]StatusRecord{},
	}
}

type InMemoryStorage struct {
	mu   sync.RWMutex
	data map[string]StatusRecord
}

func (s *InMemoryStorage) Save(id string, status StatusRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = status
}

func (s *InMemoryStorage) GetStatus(id string) (StatusRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[id]
	return val, ok
}
