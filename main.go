package main

import (
	"game-server-selector/internal/domain"
	"game-server-selector/internal/server/http"
	"game-server-selector/internal/services"
	"game-server-selector/internal/storage"
)

func main() {
	configService := services.NewConfigService()
	//TODO: after added new storage support we need to check config and init storage
	inMemoryStorage := storage.NewInMemoryStorage()
	serverService := services.NewServerService(inMemoryStorage)
	serverDomain := domain.NewServerDomain(serverService, configService)

	fiber := http.NewFiberServer(serverDomain, configService)
	err := fiber.Start()
	if err != nil {
		panic(err)
	}
}
