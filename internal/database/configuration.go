package database

import (
	"database/sql"
)

// Configuration represents the game configuration
type Configuration struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ConfigurationRepository provides methods for working with configuration
type ConfigurationRepository struct {
	db *DB
}

// NewConfigurationRepository creates a new configuration repository
func NewConfigurationRepository(db *DB) *ConfigurationRepository {
	return &ConfigurationRepository{db: db}
}

// Create creates a new configuration
func (r *ConfigurationRepository) Create(name string) (*Configuration, error) {
	stmt, err := r.db.Prepare("INSERT INTO configuration (name) VALUES (?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Configuration{
		ID:   int(id),
		Name: name,
	}, nil
}

// GetByName retrieves configuration by name
func (r *ConfigurationRepository) GetByName(name string) (*Configuration, error) {
	row := r.db.QueryRow("SELECT id, name FROM configuration WHERE name = ?", name)
	
	config := &Configuration{}
	err := row.Scan(&config.ID, &config.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return config, nil
}

// GetAll retrieves all configurations
func (r *ConfigurationRepository) GetAll() ([]*Configuration, error) {
	rows, err := r.db.Query("SELECT id, name FROM configuration")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*Configuration
	for rows.Next() {
		config := &Configuration{}
		err := rows.Scan(&config.ID, &config.Name)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

// Update updates a configuration
func (r *ConfigurationRepository) Update(config *Configuration) error {
	stmt, err := r.db.Prepare("UPDATE configuration SET name = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(config.Name, config.ID)
	return err
}

// Delete deletes a configuration by ID
func (r *ConfigurationRepository) Delete(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM configuration WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}