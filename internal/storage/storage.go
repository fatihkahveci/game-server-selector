package storage

import "game-server-selector/internal/models"

type Storage interface {
	Add(server models.ServerModel) error
	Get(key string) (models.ServerModel, error)
	Update(key string, server models.ServerModel) error
	Delete(key string) error
	List() ([]models.ServerModel, error)
	Search(req []models.SearchRequest) ([]models.ServerModel, error)
}
