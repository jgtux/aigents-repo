package domain

import (
	"time"
)


package domain

import (
	"time"

	"github.com/google/uuid"
)



type AgentSystem struct {
	AgentSystemUUID     string        `json:"agent_system_uuid"`
	CategorySystemPreset map[string]any   `json:"category_system_preset"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

type AgentCategory struct {
	CategoryID            int        `json:"category_id"`
	CategoryName          string     `json:"category_name"`
	AgentSystemUUIDPreset map[string]AgentSystem  `json:"agent_system_uuid_preset"`
	CreatedAt             time.Time  `json:"created_at"`
}


type AgentConfig struct {
	AgentConfigUUID     uuid.UUID `json:"agent_config_uuid"`
	CategoryID          map[int]AgentCategory `json:"category_id"`
	CategoryPresetEnabled bool    `json:"category_preset_enabled"`
	AgentSystemUUID     map[string]AgentSystem `json:"agent_system_uuid"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}



type Agent struct {
	AgentUUID       string `json:"agent_uuid"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	AgentConfigUUID map[string]AgentConfig `json:"agent_config_uuid"`
	AuthUUID        string `json:"auth_uuid"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	DeletedAt             time.Time `json:"deleted_at"`
}


