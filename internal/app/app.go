package app

import (
	"TransportLayer/internal/config"
	"TransportLayer/internal/server"
)

func Init(cfg *config.Config) *server.Server {
	return server.NewServer(cfg)
}
