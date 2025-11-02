package database

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/sam/makeathing/internal/items"
)

// Item represents an item in the database
type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	StatsJSON   string    `json:"stats_json"`
	IsMagical   bool      `json:"is_magical"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ItemRepository provides methods for working with items
type ItemRepository struct {
	db *DB
}

// NewItemRepository creates a new item repository
func NewItemRepository(db *DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// Create creates a new item
func (r *ItemRepository) Create(item *Item) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO items (id, name, description, type, stats_json, is_magical, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		item.ID,
		item.Name,
		item.Description,
		item.Type,
		item.StatsJSON,
		item.IsMagical,
		item.CreatedAt,
		item.UpdatedAt,
	)
	return err
}

// GetByID retrieves an item by its ID
func (r *ItemRepository) GetByID(id string) (*Item, error) {
	row := r.db.QueryRow(`
		SELECT id, name, description, type, stats_json, is_magical, created_at, updated_at
		FROM items
		WHERE id = ?
	`, id)

	item := &Item{}
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Type,
		&item.StatsJSON,
		&item.IsMagical,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return item, nil
}

// Update updates an item
func (r *ItemRepository) Update(item *Item) error {
	stmt, err := r.db.Prepare(`
		UPDATE items
		SET name = ?, description = ?, type = ?, stats_json = ?, is_magical = ?, updated_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		item.Name,
		item.Description,
		item.Type,
		item.StatsJSON,
		item.IsMagical,
		time.Now(),
		item.ID,
	)
	return err
}

// Delete deletes an item by ID
func (r *ItemRepository) Delete(id string) error {
	stmt, err := r.db.Prepare("DELETE FROM items WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// GetAll retrieves all items
func (r *ItemRepository) GetAll() ([]*Item, error) {
	rows, err := r.db.Query(`
		SELECT id, name, description, type, stats_json, is_magical, created_at, updated_at
		FROM items
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		item := &Item{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Type,
			&item.StatsJSON,
			&item.IsMagical,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// Helper functions to convert between JSON item and database item

// ItemFromJSONItem converts a JSON item to a database item
func ItemFromJSONItem(jsonItem *items.ItemJSON) (*Item, error) {
	// Convert stats to JSON
	statsJSON, err := json.Marshal(jsonItem.Stats)
	if err != nil {
		return nil, err
	}

	return &Item{
		ID:          jsonItem.ID,
		Name:        jsonItem.Name,
		Description: jsonItem.Description,
		Type:        jsonItem.Type,
		StatsJSON:   string(statsJSON),
		IsMagical:   jsonItem.IsMagical,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// ToJSONItem converts a database item to a JSON item
func (i *Item) ToJSONItem() (*items.ItemJSON, error) {
	// Convert stats JSON to struct
	var stats items.ItemStats
	if i.StatsJSON != "" {
		if err := json.Unmarshal([]byte(i.StatsJSON), &stats); err != nil {
			return nil, err
		}
	}

	return &items.ItemJSON{
		Schema:      "../schemas/item.schema.json",
		ID:          i.ID,
		Name:        i.Name,
		Description: i.Description,
		Type:        i.Type,
		Stats:       stats,
		IsMagical:   i.IsMagical,
	}, nil
}