package repositories

import (
	d "aigents-base/internal/agents/domain"
	agitf "aigents-base/internal/agents/interfaces"
	c_at "aigents-base/internal/common/atoms"
	"fmt"

	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type AgentRepository struct {
	db *sql.DB
}

func NewAgentRepository(db *sql.DB) agitf.AgentRepositoryITF {
	return &AgentRepository{db: db}
}


func (r *AgentRepository) Create(gctx *gin.Context, data *d.Agent) error {
	query := `
	WITH ins_system AS (
		INSERT INTO agent_systems (category_system_preset)
		VALUES ($1)
		RETURNING agent_system_uuid
	),
	ins_config AS (
		INSERT INTO agents_config (
			category_id,
			category_preset_enabled,
			agent_system_uuid
		)
		VALUES ($2, $3, (SELECT agent_system_uuid FROM ins_system))
		RETURNING agent_config_uuid
	)
	INSERT INTO agents (
		name,
		description,
		image_url,
		agent_config_uuid,
		auth_uuid
	)
	VALUES ($4, $5, $6, (SELECT agent_config_uuid FROM ins_config), $7)
	RETURNING agent_uuid, created_at, updated_at, COALESCE(deleted_at, NULL);
	`

	err := r.db.QueryRow(
		query,
		data.AgentConfig.AgentSystem.SystemPreset, // $1
		data.AgentConfig.Category.CategoryID,      // $2
		data.AgentConfig.CategoryPresetEnabled,    // $3
		data.Name,                                 // $4
		data.Description,                          // $5
		data.ImageURL,                             // $6
		data.AuthUUID,                             // $7
	).Scan(
		&data.AgentUUID,
		&data.CreatedAt,
		&data.UpdatedAt,
		&data.DeletedAt,
	)

	if err != nil {
		err = c_at.BuildErrLogAtom(gctx, fmt.Sprintf("(R) Failed to create agent: %s", err.Error()))
		c_at.AbortRespAtom(
			gctx,
			http.StatusInternalServerError,
			fmt.Sprintf("(R) Failed to create agent: %s", err.Error()),
		)

		return err
	}

	return nil
}

func (a *AgentRepository) GetByID(gctx *gin.Context, data *d.Agent) error {

	return nil
}

func (r *AgentRepository) Fetch(gctx *gin.Context, limit, offset uint64) ([]d.Agent, error) {
	query := `
	SELECT
		a.agent_uuid,
		a.name,
		a.description,
		a.image_url,
		a.auth_uuid,
		ac.category_id,
		ac.category_name,
		a.created_at,
		a.updated_at,
                COALESCE(deleted_at, TIMESTAMP '0001-01-01 00:00:00')
	FROM agents a
	INNER JOIN agents_config acfg ON a.agent_config_uuid = acfg.agent_config_uuid
	INNER JOIN agent_categories ac ON acfg.category_id = ac.category_id
	WHERE a.deleted_at IS NULL
	ORDER BY a.created_at DESC
	LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		err = c_at.BuildErrLogAtom(gctx, fmt.Sprintf("(R) Failed to fetch agents: %s", err.Error()))
		c_at.AbortRespAtom(
			gctx,
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to fetch agents: %s", err.Error()),
		)
		return nil, err
	}
	defer rows.Close()

	var agents []d.Agent

	for rows.Next() {
		var agent d.Agent

		err := rows.Scan(
			&agent.AgentUUID,
			&agent.Name,
			&agent.Description,
			&agent.ImageURL,
			&agent.AuthUUID,
			&agent.AgentConfig.Category.CategoryID,
			&agent.AgentConfig.Category.CategoryName,
			&agent.CreatedAt,
			&agent.UpdatedAt,
			&agent.DeletedAt,
		)
		if err != nil {
			err = c_at.BuildErrLogAtom(gctx, fmt.Sprintf("(R) Failed to scan agent: %s", err.Error()))
			c_at.AbortRespAtom(
				gctx,
				http.StatusInternalServerError,
				fmt.Sprintf("(R) Failed to scan agent: %s", err.Error()),
			)
			return nil, err
		}

		agents = append(agents, agent)
	}

	if err = rows.Err(); err != nil {
		err = c_at.BuildErrLogAtom(gctx, fmt.Sprintf("Row iteration error: %s", err.Error()))
		c_at.AbortRespAtom(
			gctx,
			http.StatusInternalServerError,
			fmt.Sprintf("(R) Row iteration error: %s", err.Error()),
		)

		return nil, err
	}

	return agents, nil
}


func (a *AgentRepository) FetchWithFilter(gctx *gin.Context, flags []string, limit, offset uint64) ([]d.Agent, error) {
	return []d.Agent{}, nil
}

func (a *AgentRepository) Update(gctx *gin.Context, data *d.Agent) error {
	return nil
}

func (a *AgentRepository) Delete(gctx *gin.Context, data *d.Agent) error {
	return nil
}
