package services

import (
	"game-server-selector/internal/models"
	"game-server-selector/internal/storage"
)

type ServerService interface {
	GetServer(ID string) (models.ServerModel, error)
	UpdateServer(ID string, server models.ServerModel) error
	DeleteServer(ID string) error
	ListServers() ([]models.ServerModel, error)
	SearchServers(search []models.SearchRequest) ([]models.ServerModel, error)
	CreateServer(server models.ServerModel) error
}

type serverService struct {
	storage       storage.Storage
	metricService MetricService
}

func NewServerService(s storage.Storage, m MetricService) ServerService {
	return &serverService{
		storage:       s,
		metricService: m,
	}
}

func (s *serverService) GetServer(ID string) (models.ServerModel, error) {
	return s.storage.Get(ID)
}

func (s *serverService) UpdateServer(ID string, server models.ServerModel) error {
	return s.storage.Update(ID, server)
}

func (s *serverService) DeleteServer(ID string) error {
	s.metricService.DecServerCount()
	return s.storage.Delete(ID)
}

func (s *serverService) ListServers() ([]models.ServerModel, error) {
	return s.storage.List()
}

func (s *serverService) SearchServers(search []models.SearchRequest) ([]models.ServerModel, error) {
	return s.storage.Search(search)
}

func (s *serverService) CreateServer(server models.ServerModel) error {
	s.metricService.IncServerCount()
	return s.storage.Add(server)
}
