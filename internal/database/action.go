package database

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/sam/makeathing/internal/actions"
	"github.com/sam/makeathing/internal/characters"
)

// Action represents an action in the database
type Action struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	Type                 string `json:"type"`
	Description          string `json:"description"`
	MinLevel             int    `json:"min_level"`
	RequiredStatsJSON    string `json:"required_stats_json"`
	RequiredSkillsJSON   string `json:"required_skills_json"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// ActionRepository provides methods for working with actions
type ActionRepository struct {
	db *DB
}

// NewActionRepository creates a new action repository
func NewActionRepository(db *DB) *ActionRepository {
	return &ActionRepository{db: db}
}

// Create creates a new action
func (r *ActionRepository) Create(action *Action) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO actions (name, type, description, min_level, required_stats_json, required_skills_json, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		action.Name,
		action.Type,
		action.Description,
		action.MinLevel,
		action.RequiredStatsJSON,
		action.RequiredSkillsJSON,
		action.CreatedAt,
		action.UpdatedAt,
	)
	return err
}

// GetByName retrieves an action by its name
func (r *ActionRepository) GetByName(name string) (*Action, error) {
	row := r.db.QueryRow(`
		SELECT id, name, type, description, min_level, required_stats_json, required_skills_json, created_at, updated_at
		FROM actions
		WHERE name = ?
	`, name)

	action := &Action{}
	err := row.Scan(
		&action.ID,
		&action.Name,
		&action.Type,
		&action.Description,
		&action.MinLevel,
		&action.RequiredStatsJSON,
		&action.RequiredSkillsJSON,
		&action.CreatedAt,
		&action.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return action, nil
}

// GetByID retrieves an action by its ID
func (r *ActionRepository) GetByID(id int) (*Action, error) {
	row := r.db.QueryRow(`
		SELECT id, name, type, description, min_level, required_stats_json, required_skills_json, created_at, updated_at
		FROM actions
		WHERE id = ?
	`, id)

	action := &Action{}
	err := row.Scan(
		&action.ID,
		&action.Name,
		&action.Type,
		&action.Description,
		&action.MinLevel,
		&action.RequiredStatsJSON,
		&action.RequiredSkillsJSON,
		&action.CreatedAt,
		&action.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return action, nil
}

// Update updates an action
func (r *ActionRepository) Update(action *Action) error {
	stmt, err := r.db.Prepare(`
		UPDATE actions
		SET name = ?, type = ?, description = ?, min_level = ?, required_stats_json = ?, required_skills_json = ?, updated_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		action.Name,
		action.Type,
		action.Description,
		action.MinLevel,
		action.RequiredStatsJSON,
		action.RequiredSkillsJSON,
		time.Now(),
		action.ID,
	)
	return err
}

// Delete deletes an action by ID
func (r *ActionRepository) Delete(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM actions WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// GetAll retrieves all actions
func (r *ActionRepository) GetAll() ([]*Action, error) {
	rows, err := r.db.Query(`
		SELECT id, name, type, description, min_level, required_stats_json, required_skills_json, created_at, updated_at
		FROM actions
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []*Action
	for rows.Next() {
		action := &Action{}
		err := rows.Scan(
			&action.ID,
			&action.Name,
			&action.Type,
			&action.Description,
			&action.MinLevel,
			&action.RequiredStatsJSON,
			&action.RequiredSkillsJSON,
			&action.CreatedAt,
			&action.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}

	return actions, nil
}

// Helper functions to convert between JSON action and database action

// ActionFromJSONAction converts a JSON action to a database action
func ActionFromJSONAction(jsonAction *actions.Action) (*Action, error) {
	// Convert required stats to JSON
	requiredStatsJSON, err := json.Marshal(jsonAction.Requirements.RequiredStats)
	if err != nil {
		return nil, err
	}

	// Convert required skills to JSON
	requiredSkillsJSON, err := json.Marshal(jsonAction.Requirements.RequiredSkills)
	if err != nil {
		return nil, err
	}

	return &Action{
		Name:               jsonAction.Name,
		Type:               string(jsonAction.Type),
		Description:        jsonAction.Description,
		MinLevel:           jsonAction.Requirements.MinLevel,
		RequiredStatsJSON:  string(requiredStatsJSON),
		RequiredSkillsJSON: string(requiredSkillsJSON),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}, nil
}

// ToJSONAction converts a database action to a JSON action
func (a *Action) ToJSONAction() (*actions.Action, error) {
	// Convert required stats JSON to struct
	var requiredStats characters.Stats
	if a.RequiredStatsJSON != "" {
		if err := json.Unmarshal([]byte(a.RequiredStatsJSON), &requiredStats); err != nil {
			return nil, err
		}
	}

	// Convert required skills JSON to slice
	var requiredSkills []string
	if a.RequiredSkillsJSON != "" {
		if err := json.Unmarshal([]byte(a.RequiredSkillsJSON), &requiredSkills); err != nil {
			return nil, err
		}
	}

	return &actions.Action{
		Name:        a.Name,
		Type:        characters.SkillType(a.Type),
		Description: a.Description,
		Requirements: actions.ActionRequirements{
			MinLevel:       a.MinLevel,
			RequiredStats:  requiredStats,
			RequiredSkills: requiredSkills,
		},
	}, nil
}