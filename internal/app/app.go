package app

import (
	"TransportLayer/internal/config"
	"TransportLayer/internal/delivery/http"
	"TransportLayer/internal/delivery/kafka"
	"TransportLayer/internal/repository/inmemory"
	"TransportLayer/internal/server"
	"TransportLayer/internal/usecase/service"
	"TransportLayer/internal/utils"
	"github.com/gorilla/mux"
	"log"
	"time"
)

func Init(cfg *config.Config) (*server.Server, *kafka.Consumer, *kafka.Producer) {
	msgRepository := inmemory.NewMessageRepository()

	msgUC := service.NewMessageService(cfg.Segment, msgRepository)

	kafkaConsumer, err := kafka.NewConsumer(cfg.Kafka.Consumer, msgUC)
	if err != nil {
		log.Fatalf("ошибка создания kafka consumer: %w", err)
	}
	kafkaProducer, err := kafka.NewProducer(cfg.Kafka.Producer)
	if err != nil {
		log.Fatalf("ошибка созданию kafka producer: %w", err)
	}

	msgHandler := http.NewMessageHandler(msgUC, cfg.Kafka, *kafkaProducer)

	go func() {
		ticker := time.NewTicker(cfg.Segment.AssemblyPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				msgUC.SendCompletedMessages(utils.SendMessage)
			}
		}
	}()

	srv := server.NewServer(cfg)
	srv.SetupRoutes(func(r *mux.Router) {
		msgHandler.Configure(r)
	})

	return srv, kafkaConsumer, kafkaProducer
}
