package export

import (
	"context"

	"herbst-server/db"
)

// Users exports all users to JSON
func Users(ctx context.Context, client *db.Client, filePath string) (int, error) {
	users, err := client.User.Query().All(ctx)
	if err != nil {
		return 0, err
	}
	return len(users), writeJSON(filePath, users)
}

// Rooms exports all rooms to JSON
func Rooms(ctx context.Context, client *db.Client, filePath string) (int, error) {
	rooms, err := client.Room.Query().All(ctx)
	if err != nil {
		return 0, err
	}
	return len(rooms), writeJSON(filePath, rooms)
}

// Abilities exports all abilities to JSON (file kept as skills.json for backward compatibility)
func Abilities(ctx context.Context, client *db.Client, filePath string) (int, error) {
	abilities, err := client.Ability.Query().All(ctx)
	if err != nil {
		return 0, err
	}
	return len(abilities), writeJSON(filePath, abilities)
}

// NPCTemplates exports all NPC templates to JSON
func NPCTemplates(ctx context.Context, client *db.Client, filePath string) (int, error) {
	templates, err := client.NPCTemplate.Query().All(ctx)
	if err != nil {
		return 0, err
	}
	return len(templates), writeJSON(filePath, templates)
}

// Equipment exports all equipment to JSON
func Equipment(ctx context.Context, client *db.Client, filePath string) (int, error) {
	equipment, err := client.Equipment.Query().All(ctx)
	if err != nil {
		return 0, err
	}
	return len(equipment), writeJSON(filePath, equipment)
}

// Characters exports all characters to JSON
func Characters(ctx context.Context, client *db.Client, filePath string) (int, error) {
	characters, err := client.Character.Query().All(ctx)
	if err != nil {
		return 0, err
	}
	return len(characters), writeJSON(filePath, characters)
}

// CharacterAbilities exports all character abilities to JSON
func CharacterAbilities(ctx context.Context, client *db.Client, filePath string) (int, error) {
	abilities, err := client.CharacterAbility.Query().All(ctx)
	if err != nil {
		return 0, err
	}
	return len(abilities), writeJSON(filePath, abilities)
}

