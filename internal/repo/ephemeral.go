package repo

import (
	"fmt"
	"sync"
	"time"

	ulid "github.com/oklog/ulid/v2"
)

const messageMaxAge = time.Second * 60

type Message struct {
	ID      ulid.ULID `db:"id"`
	Message string    `db:"message"`
}

type MessageRepo struct {
	mu       sync.RWMutex
	messages []Message
}

func NewMessageRepo() *MessageRepo {
	repo := &MessageRepo{
		messages: make([]Message, 0),
	}

	// Clean up old messages
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			err := repo.MessageTidy(messageMaxAge)
			if err != nil {
				fmt.Println("Error tidying messages:", err)
				continue
			}
		}
	}()

	return repo

}

func (e *MessageRepo) MessagesListAll() []Message {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.messages
}

func (e *MessageRepo) MessageAdd(msgContent string) ulid.ULID {
	id := ulid.Make()
	e.mu.Lock()
	defer e.mu.Unlock()
	e.messages = append(e.messages, Message{ID: id, Message: msgContent})
	return id
}
func (e *MessageRepo) MessageTidy(age time.Duration) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	cutoff := uint64(time.Now().Add(-age).UnixMilli())
	var filtered []Message
	for _, msg := range e.messages {
		// Keep message if ID timestamp is newer than cutoff
		if msg.ID.Time() > cutoff {
			filtered = append(filtered, msg)
		}
	}

	e.messages = filtered
	return nil
}
