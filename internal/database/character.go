package database

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/sam/makeathing/internal/characters"
	"github.com/sam/makeathing/internal/items"
)

// Character represents a character in the database
type Character struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Race        string    `json:"race"`
	Class       string    `json:"class"`
	StatsJSON   string    `json:"stats_json"`
	Health      int       `json:"health"`
	Mana        int       `json:"mana"`
	Experience  int       `json:"experience"`
	Level       int       `json:"level"`
	IsVendor    bool      `json:"is_vendor"`
	IsNpc       bool      `json:"is_npc"`
	InventoryJSON string   `json:"inventory_json"`
	SkillsJSON  string    `json:"skills_json"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CharacterRepository provides methods for working with characters
type CharacterRepository struct {
	db *DB
}

// NewCharacterRepository creates a new character repository
func NewCharacterRepository(db *DB) *CharacterRepository {
	return &CharacterRepository{db: db}
}

// Create creates a new character
func (r *CharacterRepository) Create(character *Character) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO characters (id, name, race, class, stats_json, health, mana, experience, level, is_vendor, is_npc, inventory_json, skills_json, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		character.ID,
		character.Name,
		character.Race,
		character.Class,
		character.StatsJSON,
		character.Health,
		character.Mana,
		character.Experience,
		character.Level,
		character.IsVendor,
		character.IsNpc,
		character.InventoryJSON,
		character.SkillsJSON,
		character.CreatedAt,
		character.UpdatedAt,
	)
	return err
}

// GetByID retrieves a character by its ID
func (r *CharacterRepository) GetByID(id string) (*Character, error) {
	row := r.db.QueryRow(`
		SELECT id, name, race, class, stats_json, health, mana, experience, level, is_vendor, is_npc, inventory_json, skills_json, created_at, updated_at
		FROM characters
		WHERE id = ?
	`, id)

	character := &Character{}
	err := row.Scan(
		&character.ID,
		&character.Name,
		&character.Race,
		&character.Class,
		&character.StatsJSON,
		&character.Health,
		&character.Mana,
		&character.Experience,
		&character.Level,
		&character.IsVendor,
		&character.IsNpc,
		&character.InventoryJSON,
		&character.SkillsJSON,
		&character.CreatedAt,
		&character.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return character, nil
}

// Update updates a character
func (r *CharacterRepository) Update(character *Character) error {
	stmt, err := r.db.Prepare(`
		UPDATE characters
		SET name = ?, race = ?, class = ?, stats_json = ?, health = ?, mana = ?, experience = ?, level = ?, is_vendor = ?, is_npc = ?, inventory_json = ?, skills_json = ?, updated_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		character.Name,
		character.Race,
		character.Class,
		character.StatsJSON,
		character.Health,
		character.Mana,
		character.Experience,
		character.Level,
		character.IsVendor,
		character.IsNpc,
		character.InventoryJSON,
		character.SkillsJSON,
		time.Now(),
		character.ID,
	)
	return err
}

// Delete deletes a character by ID
func (r *CharacterRepository) Delete(id string) error {
	stmt, err := r.db.Prepare("DELETE FROM characters WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// GetAll retrieves all characters
func (r *CharacterRepository) GetAll() ([]*Character, error) {
	rows, err := r.db.Query(`
		SELECT id, name, race, class, stats_json, health, mana, experience, level, is_vendor, is_npc, inventory_json, skills_json, created_at, updated_at
		FROM characters
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*Character
	for rows.Next() {
		character := &Character{}
		err := rows.Scan(
			&character.ID,
			&character.Name,
			&character.Race,
			&character.Class,
			&character.StatsJSON,
			&character.Health,
			&character.Mana,
			&character.Experience,
			&character.Level,
			&character.IsVendor,
			&character.IsNpc,
			&character.InventoryJSON,
			&character.SkillsJSON,
			&character.CreatedAt,
			&character.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, nil
}

// Helper functions to convert between JSON character and database character

// CharacterFromJSONCharacter converts a JSON character to a database character
func CharacterFromJSONCharacter(jsonCharacter *characters.CharacterJSON) (*Character, error) {
	// Convert stats to JSON
	statsJSON, err := json.Marshal(jsonCharacter.Stats)
	if err != nil {
		return nil, err
	}

	// Convert inventory to JSON
	inventoryJSON, err := json.Marshal(jsonCharacter.Inventory)
	if err != nil {
		return nil, err
	}

	// Convert skills to JSON
	skillsJSON, err := json.Marshal(jsonCharacter.Skills)
	if err != nil {
		return nil, err
	}

	return &Character{
		ID:            jsonCharacter.ID,
		Name:          jsonCharacter.Name,
		Race:          jsonCharacter.Race,
		Class:         jsonCharacter.Class,
		StatsJSON:     string(statsJSON),
		Health:        jsonCharacter.Health,
		Mana:          jsonCharacter.Mana,
		Experience:    jsonCharacter.Experience,
		Level:         jsonCharacter.Level,
		IsVendor:      jsonCharacter.IsVendor,
		IsNpc:         jsonCharacter.IsNpc,
		InventoryJSON: string(inventoryJSON),
		SkillsJSON:    string(skillsJSON),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

// ToJSONCharacter converts a database character to a JSON character
func (c *Character) ToJSONCharacter() (*characters.CharacterJSON, error) {
	// Convert stats JSON to struct
	var stats characters.Stats
	if c.StatsJSON != "" {
		if err := json.Unmarshal([]byte(c.StatsJSON), &stats); err != nil {
			return nil, err
		}
	}

	// Convert inventory JSON to slice
	var inventory []items.Item
	if c.InventoryJSON != "" {
		if err := json.Unmarshal([]byte(c.InventoryJSON), &inventory); err != nil {
			return nil, err
		}
	}

	// Convert skills JSON to slice
	var skills []characters.Skill
	if c.SkillsJSON != "" {
		if err := json.Unmarshal([]byte(c.SkillsJSON), &skills); err != nil {
			return nil, err
		}
	}

	return &characters.CharacterJSON{
		Schema:     "../schemas/character.schema.json",
		ID:         c.ID,
		Name:       c.Name,
		Race:       c.Race,
		Class:      c.Class,
		Stats:      stats,
		Health:     c.Health,
		Mana:       c.Mana,
		Experience: c.Experience,
		Level:      c.Level,
		IsVendor:   c.IsVendor,
		IsNpc:      c.IsNpc,
		Inventory:  inventory,
		Skills:     skills,
	}, nil
}