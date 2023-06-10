package models

type CreateServerRequest struct {
	Name               string                 `validate:"required,min=1" json:"name"`
	IP                 string                 `validate:"required,ip" json:"ip"`
	Port               int                    `validate:"required,min=1,max=65535" json:"port"`
	Capacity           int                    `validate:"required" json:"capacity"`
	CurrentPlayerCount int                    `validate:"required" json:"current_player_count"`
	CustomData         map[string]interface{} `json:"custom_data"`
}

type UpdateServerRequest struct {
	Name               string                 `validate:"required,min=1" json:"name"`
	IP                 string                 `validate:"required,ip" json:"ip"`
	Port               int                    `validate:"required,min=1,max=65535" json:"port"`
	Capacity           int                    `validate:"required" json:"capacity"`
	CurrentPlayerCount int                    `validate:"required" json:"current_player_count"`
	CustomData         map[string]interface{} `json:"custom_data"`
}

type SearchRequest struct {
	Field string `validate:"required" json:"field"`
	Query struct {
		Operator string      `validate:"required" json:"operator"`
		Value    interface{} `json:"value"`
	} `validate:"required" json:"query"`
}
