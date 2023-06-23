package main

import (
	"game-server-selector/internal/domain"
	serverhttp "game-server-selector/internal/server/http"
	"game-server-selector/internal/services"
	"game-server-selector/internal/storage"

	"net/http"
)

func main() {
	configService := services.NewConfigService()
	//TODO: after added new storage support we need to check config and init storage
	inMemoryStorage := storage.NewInMemoryStorage()
	metricService := services.NewMetricService()
	metricService.Register()
	serverService := services.NewServerService(inMemoryStorage, metricService)
	serverDomain := domain.NewServerDomain(serverService, configService)

	if configService.IsPrometheusEnabled() {
		handler := metricService.Handler()
		http.Handle("/metrics", handler)
		go http.ListenAndServe(configService.GetPrometheusPort(), nil)
	}
	fiber := serverhttp.NewFiberServer(serverDomain, configService, metricService)
	err := fiber.Start()
	if err != nil {
		panic(err)
	}
}
