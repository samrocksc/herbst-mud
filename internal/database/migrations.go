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
	{
		Name: "004_create_rooms_table",
		SQL: `
CREATE TABLE IF NOT EXISTS rooms (
	id TEXT PRIMARY KEY,
	description TEXT NOT NULL,
	smells TEXT,
	exits_json TEXT,
	immovable_objects_json TEXT,
	movable_objects_json TEXT,
	npcs_json TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`,
	},
	{
		Name: "005_create_characters_table",
		SQL: `
CREATE TABLE IF NOT EXISTS characters (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	race TEXT NOT NULL,
	class TEXT NOT NULL,
	stats_json TEXT NOT NULL,
	health INTEGER NOT NULL,
	mana INTEGER NOT NULL,
	experience INTEGER NOT NULL,
	level INTEGER NOT NULL,
	is_vendor BOOLEAN NOT NULL,
	is_npc BOOLEAN NOT NULL,
	inventory_json TEXT,
	skills_json TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`,
	},
	{
		Name: "006_create_items_table",
		SQL: `
CREATE TABLE IF NOT EXISTS items (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	type TEXT NOT NULL,
	stats_json TEXT NOT NULL,
	is_magical BOOLEAN NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`,
	},
	{
		Name: "007_create_actions_table",
		SQL: `
CREATE TABLE IF NOT EXISTS actions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	type TEXT NOT NULL,
	description TEXT NOT NULL,
	min_level INTEGER NOT NULL,
	required_stats_json TEXT NOT NULL,
	required_skills_json TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`,
	},
	{
		Name: "008_create_global_state_characters_table",
		SQL: `
CREATE TABLE IF NOT EXISTS global_state_characters (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	character_id TEXT NOT NULL,
	current_room_id TEXT NOT NULL,
	health INTEGER NOT NULL,
	status TEXT NOT NULL,
	last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (character_id) REFERENCES characters(id),
	FOREIGN KEY (current_room_id) REFERENCES rooms(id)
);
`,
	},
	{
		Name: "009_create_global_state_rooms_table",
		SQL: `
CREATE TABLE IF NOT EXISTS global_state_rooms (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	room_id TEXT NOT NULL,
	player_count INTEGER NOT NULL DEFAULT 0,
	npc_state_json TEXT,
	item_state_json TEXT,
	last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (room_id) REFERENCES rooms(id)
);
`,
	},
	{
		Name: "010_add_username_password_to_users",
		SQL: `
ALTER TABLE users ADD COLUMN username TEXT;
ALTER TABLE users ADD COLUMN password TEXT;
UPDATE users SET username = 'user_' || id WHERE username IS NULL;
UPDATE users SET password = 'password' WHERE password IS NULL;
CREATE UNIQUE INDEX idx_users_username ON users(username);
`,
	},
}