package service

import (
	"TransportLayer/internal/config"
	"TransportLayer/internal/entity"
	"TransportLayer/internal/usecase"
)

type MessageService struct {
	cfg config.SegmentConfig
}

func NewMessageService(cfg config.SegmentConfig) usecase.MessageService {
	return &MessageService{
		cfg: cfg,
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
