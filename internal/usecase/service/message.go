package service

import (
	"TransportLayer/internal/config"
	"TransportLayer/internal/entity"
	"TransportLayer/internal/repository"
	"TransportLayer/internal/usecase"
	"fmt"
	"time"
)

type MessageService struct {
	cfg           config.SegmentConfig
	msgRepository repository.MessageRepository
}

func NewMessageService(cfg config.SegmentConfig, msgRepository repository.MessageRepository) usecase.MessageService {
	return &MessageService{
		cfg:           cfg,
		msgRepository: msgRepository,
	}
}

func (s *MessageService) SegmentMessage(msg entity.SendRequest) ([]*entity.Segment, error) {
	segments := splitString(msg.Data, s.cfg.MaxSegmentSize)

	result := make([]*entity.Segment, len(segments))

	for i, seg := range segments {
		result[i] = &entity.Segment{
			SegmentNumber:  i + 1,
			TotalSegments:  len(segments),
			Username:       msg.Username,
			SendTime:       msg.SendTime,
			SegmentPayload: seg,
		}
	}

	return result, nil
}

func (s *MessageService) AddSegment(segment *entity.Segment) {
	s.msgRepository.AddSegment(segment)
}

func (s *MessageService) SendCompletedMessages(sender func(body entity.ReceiveRequest)) {
	messages := s.msgRepository.GetAllMessages()

	for sendTime, message := range messages {
		if message.Received < message.Total &&
			time.Since(message.Last) <= s.cfg.AssemblyPeriod+time.Second {
			continue
		}

		payload := entity.ReceiveRequest{
			Username: message.Username,
			SendTime: sendTime,
		}

		if message.Received == message.Total {
			payload.Data = assembleMessage(message)
			fmt.Printf("sent message: %+v\n", payload)
		} else {
			payload.Error = entity.ErrSegmentLost
			fmt.Printf("sent error: %+v\n", payload)
		}
		go sender(payload)
		s.msgRepository.DeleteMessage(sendTime)
	}
}

func splitString(s string, maxSegmentSize int) []string {
	var chunks []string
	for i := 0; i < len(s); i += maxSegmentSize {
		end := i + maxSegmentSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

func assembleMessage(msg *entity.Message) string {
	result := ""
	for _, segment := range msg.Segments {
		result += segment
	}
	return result
}
