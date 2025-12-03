package atoms

import (
	c_at "aigents-base/internal/common/atoms"
	d "aigents-base/internal/agents/domain"
	"database/sql"
	"github.com/gin-gonic/gin"
)


func GetAgentByUUID(gctx *gin.Context, agentUUID string) *d.Agent {
	query := `
	SELECT
		a.agent_uuid,
		a.name,
		a.description,
		a.image_url,
		a.auth_uuid,
		a.created_at,
		a.updated_at,
		COALESCE(a.deleted_at, TIMESTAMP '0001-01-01 00:00:00'),
		ac.category_id,
		ac.category_name,
		acfg.agent_config_uuid,
		acfg.category_preset_enabled,
		asys.agent_system_uuid,
		asys.system_preset
	FROM agents a
	INNER JOIN agents_config acfg ON a.agent_config_uuid = acfg.agent_config_uuid
	INNER JOIN agent_categories ac ON acfg.category_id = ac.category_id
	INNER JOIN agent_systems asys ON acfg.agent_system_uuid = asys.agent_system_uuid
	WHERE a.agent_uuid = $1 AND a.deleted_at IS NULL;
	`

	var data d.Agent
	var systemPresetJSON []byte

	err := r.db.QueryRow(query, agentUUID).Scan(
		&data.AgentUUID,
		&data.Name,
		&data.Description,
		&data.ImageURL,
		&data.AuthUUID,
		&data.CreatedAt,
		&data.UpdatedAt,
		&data.DeletedAt,
		&data.AgentConfig.Category.CategoryID,
		&data.AgentConfig.Category.CategoryName,
		&data.AgentConfig.AgentConfigUUID,
		&data.AgentConfig.CategoryPresetEnabled,
		&data.AgentConfig.AgentSystem.AgentSystemUUID,
		&systemPresetJSON,
	)

	if err == sql.ErrNoRows {
		err = c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusNotFound,
			"(R) Agent not found.",
			fmt.Sprintf("Agent with UUID %s not found", data.AgentUUID))
		return err
	}

	if err != nil {
		err = c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusInternalServerError,
			"(R) Could not get agent.",
			fmt.Sprintf("Failed to get agent: %s", err.Error()))
		return err
	}

	if err := json.Unmarshal(systemPresetJSON, &data.AgentConfig.AgentSystem.SystemPreset); err != nil {
		err = c_at.AbortAndBuildErrLogAtom(
			gctx,
			http.StatusInternalServerError,
			"(R) Could not parse agent system preset.",
			fmt.Sprintf("Failed to unmarshal system_preset: %s", err.Error()))
		return err
	}

	return nil
}
