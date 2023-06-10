package domain

import (
	"game-server-selector/internal/models"
	"game-server-selector/internal/services"
	"log"
	"time"

	"github.com/google/uuid"
)

type ServerDomain interface {
	CreateServer(req models.CreateServerRequest) (models.ServerModel, error)
	GetServer(ID string) (models.Server, error)
	UpdateServer(ID string, req models.UpdateServerRequest) (models.Server, error)
	DeleteServer(ID string) error
	ListServers() ([]models.Server, error)
	SearchServers(search []models.SearchRequest) ([]models.Server, error)
}

type serverDomain struct {
	serverService services.ServerService
	configService services.ConfigService
}

func NewServerDomain(s services.ServerService, c services.ConfigService) ServerDomain {
	return &serverDomain{
		serverService: s,
		configService: c,
	}
}

func (s *serverDomain) CreateServer(req models.CreateServerRequest) (models.ServerModel, error) {
	now := time.Now()
	server := models.ServerModel{
		ID:                 uuid.New().String(),
		Name:               req.Name,
		Port:               req.Port,
		IP:                 req.IP,
		CustomData:         req.CustomData,
		Capacity:           req.Capacity,
		CurrentPlayerCount: req.CurrentPlayerCount,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	err := s.serverService.CreateServer(server)
	if err != nil {
		return models.ServerModel{}, err
	}
	return server, nil
}

func (s *serverDomain) GetServer(ID string) (models.Server, error) {
	serverModel, err := s.serverService.GetServer(ID)
	if err != nil {
		return models.Server{}, err
	}
	return serverModel.ToServer(), nil
}

func (s *serverDomain) UpdateServer(ID string, req models.UpdateServerRequest) (models.Server, error) {
	server, err := s.serverService.GetServer(ID)
	if err != nil {
		return models.Server{}, err
	}
	server.Name = req.Name
	server.Port = req.Port
	server.IP = req.IP
	server.CustomData = req.CustomData
	server.Capacity = req.Capacity
	server.UpdatedAt = time.Now()
	server.CurrentPlayerCount = req.CurrentPlayerCount
	err = s.serverService.UpdateServer(ID, server)
	if err != nil {
		return models.Server{}, err
	}
	return server.ToServer(), nil
}

func (s *serverDomain) DeleteServer(ID string) error {
	return s.serverService.DeleteServer(ID)
}

func (s *serverDomain) ListServers() ([]models.Server, error) {
	serverModelList, err := s.serverService.ListServers()
	if err != nil {
		return nil, err
	}
	serverList := []models.Server{}
	for _, serverModel := range serverModelList {
		if s.isServerDirty(serverModel) {
			err = s.DeleteServer(serverModel.ID)
			if err != nil {
				log.Default().Println("Error deleting server: ", err)
			}
			continue
		}
		serverList = append(serverList, serverModel.ToServer())
	}
	return serverList, nil
}

func (s *serverDomain) SearchServers(search []models.SearchRequest) ([]models.Server, error) {
	serverModelList, err := s.serverService.SearchServers(search)
	if err != nil {
		return nil, err
	}
	serverList := []models.Server{}
	for _, serverModel := range serverModelList {
		if s.isServerDirty(serverModel) {
			err = s.DeleteServer(serverModel.ID)
			if err != nil {
				log.Default().Println("Error deleting server: ", err)
			}
			continue
		}
		serverList = append(serverList, serverModel.ToServer())
	}
	return serverList, nil
}

func (s *serverDomain) isServerDirty(server models.ServerModel) bool {
	dirtySecond := s.configService.GetServerDirtySeconds()
	now := time.Now()
	if now.Unix() >= server.UpdatedAt.Add(time.Second*time.Duration(dirtySecond)).Unix() {
		return true
	}
	return false
}
