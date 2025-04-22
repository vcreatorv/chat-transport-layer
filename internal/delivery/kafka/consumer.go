package kafka

import (
	"TransportLayer/internal/config"
	"TransportLayer/internal/entity"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
	cfg      config.KafkaConsumerConfig
}

func NewConsumer(cfg config.KafkaConsumerConfig) (*Consumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = cfg.ReturnErrors

	consumer, err := sarama.NewConsumer(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		cfg:      cfg,
	}, nil
}

func (c *Consumer) ReadFromKafka() error {
	partitionConsumer, err := c.consumer.ConsumePartition(c.cfg.Topic, 0, int64(c.cfg.AutoOffsetReset))
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var segment entity.Segment
			if err := json.Unmarshal(msg.Value, &segment); err != nil {
				fmt.Printf("ошибка при десериализации сообщения из consumer: %v", err)
			}
			fmt.Printf("%+v\n", segment)
		case err := <-partitionConsumer.Errors():
			fmt.Printf("ошибка при чтении из consumer: %s\n", err.Error())
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
