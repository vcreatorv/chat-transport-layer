package kafka

import (
	"TransportLayer/internal/config"
	"TransportLayer/internal/entity"
	"encoding/json"
	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	cfg      config.KafkaProducerConfig
}

func NewProducer(cfg config.KafkaProducerConfig) (*Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.RequiredAcks(cfg.RequiredAcks)
	saramaConfig.Producer.Retry.Max = cfg.RetryMax
	saramaConfig.Producer.Return.Successes = cfg.ReturnSuccess

	producer, err := sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		cfg:      cfg,
	}, nil
}

func (p *Producer) WriteToKafka(segment entity.Segment) error {
	segmentBytes, err := json.Marshal(segment)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.cfg.Topic,
		Value: sarama.StringEncoder(segmentBytes),
	}

	_, _, err = p.producer.SendMessage(msg)
	return err
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
