package database

// Migration represents a single database migration
type Migration struct {
	Name string
	SQL  string
}

// migrations is the list of all migrations in order
var migrations = []Migration{
	{
		Name: "001_create_configuration_table",
		SQL: `
CREATE TABLE IF NOT EXISTS configuration (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);
`,
	},
	{
		Name: "002_create_sessions_table",
		SQL: `
CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY,
	user_id INTEGER NOT NULL,
	character_id TEXT NOT NULL,
	room_id TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`,
	},
	{
		Name: "003_create_users_table",
		SQL: `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	character_id TEXT NOT NULL,
	room_id TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`,
	},
}