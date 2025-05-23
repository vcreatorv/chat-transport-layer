package usecase

import (
	"TransportLayer/internal/entity"
)

type MessageService interface {
	SegmentMessage(message entity.SendRequest) ([]*entity.Segment, error)
	AddSegment(segment *entity.Segment)
	SendCompletedMessages(sender func(body entity.ReceiveRequest))
}
