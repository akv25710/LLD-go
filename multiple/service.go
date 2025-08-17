package multiple

import (
	"context"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

type Manager struct {
	mu     sync.Mutex
	states map[string]State
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func NewManager() *Manager {
	return &Manager{
		mu:     sync.Mutex{},
		states: make(map[string]State),
	}
}

func (manager *Manager) GetState(id string) State {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	log.Printf("Getting state for: %s", id)
	return manager.states[id]
}

func (manager *Manager) SetState(id string, state State) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	log.Printf("Setting state for: %s", id)
	manager.states[id] = state
}

func (manager *Manager) Start(ctx context.Context, count int) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	manager.cancel = cancel

	for i := 0; i < count; i++ {
		manager.wg.Add(1)
		id := uuid.New()
		go func(id string) {
			defer func() {
				manager.wg.Done()
				if r := recover(); r != nil {
					manager.SetState(id, Failed)
				}
			}()

			manager.startGoRoutine(id)
		}(id.String())
	}
	manager.wg.Wait()
	manager.Stop()
}

func (manager *Manager) Stop() {
	log.Printf("Stopping Manager")
	manager.wg.Add(1)
	go func() {
		defer manager.wg.Done()
		manager.cancel()
	}()
	manager.wg.Wait()
	log.Printf("Manager stopped")
}

func (manager *Manager) startGoRoutine(id string) {
	log.Printf("Starting Go Routine: %s", id)
	manager.states[id] = Created

	manager.SetState(id, Created)
	if manager.work(id) != nil {
		manager.SetState(id, Completed)
	} else {
		manager.SetState(id, Failed)
	}
}

func (manager *Manager) work(id string) error {
	log.Printf("Working on %s", id)
	manager.SetState(id, Running)
	time.Sleep(1 * time.Second)
	return nil
}
