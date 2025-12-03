package domain

import (
	"time"
)

type AgentSystem struct {
	AgentSystemUUID     string        `json:"agent_system_uuid"`
	SystemPreset map[string]any   `json:"system_preset"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

type AgentCategory struct {
	CategoryID            uint64        `json:"category_id"`
	CategoryName          string     `json:"category_name"`
	AgentSystemPreset AgentSystem  `json:"agent_system_preset"`
	CreatedAt             time.Time  `json:"created_at"`
}


type AgentConfig struct {
	AgentConfigUUID    string `json:"agent_config_uuid"`
	Category          AgentCategory `json:"agent_category"`
	CategoryPresetEnabled bool    `json:"category_preset_enabled"`
	AgentSystem     AgentSystem `json:"agent_system"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}



type Agent struct {
	AgentUUID       string `json:"agent_uuid"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
        ImageURL        string   `json:"image_url"`
	AgentConfig     AgentConfig `json:"agent_config"`
	AuthUUID        string `json:"auth_uuid"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	DeletedAt             time.Time `json:"deleted_at"`
}


