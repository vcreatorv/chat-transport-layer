package main

import (
	"TransportLayer/internal/app"
	"TransportLayer/internal/config"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig("configs/config.yml")
	if err != nil {
		log.Fatal(err)
		return
	}

	srv, kafkaConsumer, kafkaProducer := app.Init(cfg)

	go func() {
		if err := kafkaConsumer.ReadFromKafka(); err != nil {
			fmt.Println(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-quit
		fmt.Printf("Останавливаем сервер. Пришел сигнал: %v\n", sig)

		if err := kafkaProducer.Close(); err != nil {
			log.Fatalf("Не удалось остановить кафка продюсера: %v\n", err)
		}
		if err := kafkaConsumer.Close(); err != nil {
			log.Fatalf("Не удалось остановить кафка консюмера: %v\n", err)
		}

		if err := srv.GracefulStop(); err != nil {
			log.Fatalf("Не удалось остановить сервер: %v\n", err)
		}
	}()

	fmt.Printf("Запуск сервера на порту %s\n", cfg.HTTP.Port)
	if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Не удалось запустить сервер: %v\n", err)
	}
}
