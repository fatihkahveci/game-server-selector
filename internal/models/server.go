package models

import "time"

type ServerModel struct {
	ID                 string                 `json:"id"`
	Name               string                 `json:"name"`
	IP                 string                 `json:"ip"`
	Port               int                    `json:"port"`
	Capacity           int                    `json:"capacity"`
	CurrentPlayerCount int                    `json:"current_player_count"`
	CustomData         map[string]interface{} `json:"custom_data"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}
type Server struct {
	Name               string                 `json:"name"`
	IP                 string                 `json:"ip"`
	Port               int                    `json:"port"`
	Capacity           int                    `json:"capacity"`
	CurrentPlayerCount int                    `json:"current_player_count"`
	CustomData         map[string]interface{} `json:"custom_data"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

func (s *ServerModel) ToServer() Server {
	return Server{
		Name:               s.Name,
		IP:                 s.IP,
		Port:               s.Port,
		Capacity:           s.Capacity,
		CurrentPlayerCount: s.CurrentPlayerCount,
		CustomData:         s.CustomData,
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
	}
}
