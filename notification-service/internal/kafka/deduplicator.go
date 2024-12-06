package kafka

import (
	"sync"
	"time"
)

type MessageDeduplicator struct {
	mu       sync.Mutex
	messages map[string]time.Time
}

func NewMessageDeduplicator() *MessageDeduplicator {
	return &MessageDeduplicator{messages: make(map[string]time.Time)}
}

func (d *MessageDeduplicator) IsDuplicate(eventID string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if lastProcessed, exists := d.messages[eventID]; exists && time.Since(lastProcessed) < 30*time.Second {
		return true
	}
	d.messages[eventID] = time.Now()
	return false
}

// Удаляем устаревшие записи из памяти.
func (d *MessageDeduplicator) Cleanup(expiry time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for id, timestamp := range d.messages {
		if time.Since(timestamp) > expiry {
			delete(d.messages, id)
		}
	}
}
