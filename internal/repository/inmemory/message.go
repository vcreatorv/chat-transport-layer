package inmemory

import (
	"TransportLayer/internal/entity"
	"TransportLayer/internal/repository"
	"sync"
	"time"
)

type MessageRepository struct {
	storage map[time.Time]*entity.Message
	mu      sync.Mutex
}

func NewMessageRepository() repository.MessageRepository {
	return &MessageRepository{
		storage: make(map[time.Time]*entity.Message),
	}
}

func (r *MessageRepository) AddSegment(segment *entity.Segment) {
	r.mu.Lock()
	defer r.mu.Unlock()

	sendTime := segment.SendTime
	if _, exists := r.storage[sendTime]; !exists {
		r.storage[sendTime] = &entity.Message{
			Received: 0,
			Total:    segment.TotalSegments,
			Last:     time.Now().UTC(),
			Username: segment.Username,
			Segments: make([]string, segment.TotalSegments),
		}
	}

	msg := r.storage[sendTime]
	msg.Received++
	msg.Last = time.Now().UTC()
	msg.Segments[segment.SegmentNumber-1] = segment.SegmentPayload
	r.storage[sendTime] = msg
}

func (r *MessageRepository) GetAllMessages() map[time.Time]*entity.Message {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.storage
}

func (r *MessageRepository) DeleteMessage(sendTime time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.storage, sendTime)
}
