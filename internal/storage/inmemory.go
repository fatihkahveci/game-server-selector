package storage

import (
	"errors"
	"game-server-selector/internal/models"
	"game-server-selector/internal/pkg/search"
	"sync"
	"time"
)

var ErrNotFound = errors.New("not found")

type InMemoryStorage struct {
	servers map[string]models.ServerModel
	mu      sync.RWMutex
	search  search.SearchService
}

func NewInMemoryStorage() Storage {
	ss := search.NewSearchService()
	s := &InMemoryStorage{
		search:  ss,
		servers: map[string]models.ServerModel{},
	}
	//go s.tick()
	return s
}

func (s *InMemoryStorage) Get(key string) (models.ServerModel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	server, ok := s.servers[key]
	if !ok {
		return models.ServerModel{}, ErrNotFound
	}

	return server, nil
}

func (s *InMemoryStorage) Update(key string, server models.ServerModel) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.servers[key] = server

	return nil
}

func (s *InMemoryStorage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.servers, key)

	return nil
}

func (s *InMemoryStorage) List() ([]models.ServerModel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	servers := []models.ServerModel{}
	for _, server := range s.servers {
		servers = append(servers, server)
	}

	return servers, nil
}

func (s *InMemoryStorage) Search(req []models.SearchRequest) ([]models.ServerModel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	servers := []models.ServerModel{}

	for _, server := range s.servers {
		ok, err := s.search.SearchOne(req, server.ToServer())
		if err != nil {
			continue
		}
		if ok {
			servers = append(servers, server)
		}
	}

	return servers, nil
}

func (s *InMemoryStorage) Add(server models.ServerModel) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.servers[server.ID] = server

	return nil
}

func (s *InMemoryStorage) tick() {
	for {
		for _, v := range s.servers {
			if v.UpdatedAt.Unix() < time.Now().Add(-time.Minute*5).Unix() {
				s.Delete(v.ID)
			}
		}
	}
}
