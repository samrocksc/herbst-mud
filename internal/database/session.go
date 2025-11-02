package database

import (
	"database/sql"
	"time"
)

// SessionRepository provides methods for working with sessions
type SessionRepository struct {
	db *DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create creates a new session
func (r *SessionRepository) Create(session *Session) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO sessions (id, user_id, character_id, room_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		session.ID,
		session.UserID,
		session.CharacterID,
		session.RoomID,
		session.CreatedAt,
		session.UpdatedAt,
	)
	return err
}

// GetByID retrieves a session by its ID
func (r *SessionRepository) GetByID(id string) (*Session, error) {
	row := r.db.QueryRow(`
		SELECT id, user_id, character_id, room_id, created_at, updated_at
		FROM sessions
		WHERE id = ?
	`, id)

	session := &Session{}
	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.CharacterID,
		&session.RoomID,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return session, nil
}

// Update updates a session
func (r *SessionRepository) Update(session *Session) error {
	stmt, err := r.db.Prepare(`
		UPDATE sessions
		SET user_id = ?, character_id = ?, room_id = ?, updated_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		session.UserID,
		session.CharacterID,
		session.RoomID,
		time.Now(),
		session.ID,
	)
	return err
}

// Delete deletes a session by ID
func (r *SessionRepository) Delete(id string) error {
	stmt, err := r.db.Prepare("DELETE FROM sessions WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// GetAll retrieves all sessions
func (r *SessionRepository) GetAll() ([]*Session, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, character_id, room_id, created_at, updated_at
		FROM sessions
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*Session
	for rows.Next() {
		session := &Session{}
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.CharacterID,
			&session.RoomID,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}