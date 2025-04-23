package repository

import (
	"TransportLayer/internal/entity"
	"time"
)

type MessageRepository interface {
	AddSegment(segment *entity.Segment)
	GetAllMessages() map[time.Time]*entity.Message
	DeleteMessage(sendTime time.Time)
}
