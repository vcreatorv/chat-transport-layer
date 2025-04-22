package app

import (
	"TransportLayer/internal/config"
	"TransportLayer/internal/delivery/http"
	"TransportLayer/internal/delivery/kafka"
	"TransportLayer/internal/server"
	"TransportLayer/internal/usecase/service"
	"github.com/gorilla/mux"
	"log"
)

func Init(cfg *config.Config) (*server.Server, *kafka.Consumer, *kafka.Producer) {
	kafkaConsumer, err := kafka.NewConsumer(cfg.Kafka.Consumer)
	if err != nil {
		log.Fatalf("ошибка создания кафра консюмера: %w", err)
	}
	kafkaProducer, err := kafka.NewProducer(cfg.Kafka.Producer)
	if err != nil {
		log.Fatalf("ошибка созданию кафка продюсера: %w", err)
	}

	msgUC := service.NewMessageService(cfg.Segment)
	msgHandler := http.NewMessageHandler(msgUC, cfg.Kafka, *kafkaProducer)

	srv := server.NewServer(cfg)
	srv.SetupRoutes(func(r *mux.Router) {
		msgHandler.Configure(r)
	})

	return srv, kafkaConsumer, kafkaProducer
}
