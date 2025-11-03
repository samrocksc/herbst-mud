package database

import (
	"database/sql"
	"time"

	"github.com/sam/makeathing/internal/users"
)

// User represents a user in the database
type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	CharacterID string    `json:"character_id"`
	RoomID      string    `json:"room_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserFromJSONUser converts a JSON user to a database user
func UserFromJSONUser(jsonUser *users.UserJSON) (*User, error) {
	return &User{
		ID:          jsonUser.ID,
		Username:    jsonUser.Username,
		Password:    jsonUser.Password,
		CharacterID: jsonUser.CharacterID,
		RoomID:      jsonUser.RoomID,
		CreatedAt:   jsonUser.CreatedAt,
		UpdatedAt:   jsonUser.UpdatedAt,
	}, nil
}

// ToJSONUser converts a database user to a JSON user
func (u *User) ToJSONUser() (*users.UserJSON, error) {
	return &users.UserJSON{
		Schema:      "", // This would typically be set when saving to JSON
		ID:          u.ID,
		Username:    u.Username,
		Password:    u.Password,
		CharacterID: u.CharacterID,
		RoomID:      u.RoomID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}, nil
}

// UserRepository provides methods for working with users
type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *User) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO users (username, password, character_id, room_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		user.Username,
		user.Password,
		user.CharacterID,
		user.RoomID,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

// GetByID retrieves a user by its ID
func (r *UserRepository) GetByID(id int) (*User, error) {
	row := r.db.QueryRow(`
		SELECT id, username, password, character_id, room_id, created_at, updated_at
		FROM users
		WHERE id = ?
	`, id)

	user := &User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CharacterID,
		&user.RoomID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*User, error) {
	row := r.db.QueryRow(`
		SELECT id, username, password, character_id, room_id, created_at, updated_at
		FROM users
		WHERE username = ?
	`, username)

	user := &User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CharacterID,
		&user.RoomID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// GetByCharacterID retrieves a user by character ID
func (r *UserRepository) GetByCharacterID(characterID string) (*User, error) {
	row := r.db.QueryRow(`
		SELECT id, username, password, character_id, room_id, created_at, updated_at
		FROM users
		WHERE character_id = ?
	`, characterID)

	user := &User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CharacterID,
		&user.RoomID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *User) error {
	stmt, err := r.db.Prepare(`
		UPDATE users
		SET username = ?, password = ?, character_id = ?, room_id = ?, updated_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		user.Username,
		user.Password,
		user.CharacterID,
		user.RoomID,
		time.Now(),
		user.ID,
	)
	return err
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// GetAll retrieves all users
func (r *UserRepository) GetAll() ([]*User, error) {
	rows, err := r.db.Query(`
		SELECT id, username, password, character_id, room_id, created_at, updated_at
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.CharacterID,
			&user.RoomID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}