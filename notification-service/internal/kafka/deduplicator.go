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

	// Проверяем обработано ли сообщение за последние 30 с
	if lastProcessed, exists := d.messages[eventID]; exists && time.Since(lastProcessed) < 30*time.Second {
		return true
	}
	d.messages[eventID] = time.Now()
	return false
}
