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

	server := app.Init(cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-quit
		fmt.Printf("Останавливаем сервер. Пришел сигнал: %v\n", sig)
		if err := server.GracefulStop(); err != nil {
			log.Fatalf("Не удалось остановить сервер: %v", err)
		}
	}()

	fmt.Printf("Запуск сервера на порту %s\n", cfg.HTTP.Port)
	if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
