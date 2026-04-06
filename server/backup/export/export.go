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

// Skills exports all skills to JSON
func Skills(ctx context.Context, client *db.Client, filePath string) (int, error) {
	skills, err := client.Skill.Query().All(ctx)
	if err != nil {
		return 0, err
	}
	return len(skills), writeJSON(filePath, skills)
}

// Talents exports all talents to JSON
func Talents(ctx context.Context, client *db.Client, filePath string) (int, error) {
	talents, err := client.Talent.Query().All(ctx)
	if err != nil {
		return 0, err
	}
	return len(talents), writeJSON(filePath, talents)
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

// CharacterSkills exports all character skills to JSON
func CharacterSkills(ctx context.Context, client *db.Client, filePath string) (int, error) {
	skills, err := client.CharacterSkill.Query().All(ctx)
	if err != nil {
		return 0, err
	}
	return len(skills), writeJSON(filePath, skills)
}

// CharacterTalents exports all character talents to JSON
func CharacterTalents(ctx context.Context, client *db.Client, filePath string) (int, error) {
	talents, err := client.CharacterTalent.Query().All(ctx)
	if err != nil {
		return 0, err
	}
	return len(talents), writeJSON(filePath, talents)
}

// AvailableTalents exports all available talents to JSON
func AvailableTalents(ctx context.Context, client *db.Client, filePath string) (int, error) {
	talents, err := client.AvailableTalent.Query().All(ctx)
	if err != nil {
		return 0, err
	}
	return len(talents), writeJSON(filePath, talents)
}