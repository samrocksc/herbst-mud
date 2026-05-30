package dbinit

import (
	"context"
	"encoding/json"
	"log"

	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/gender"
	"herbst-server/db/race"
)

// InitRaces seeds the race table for the specified world if empty.
// Default world_id is "1" for backward compatibility.
func InitRaces(client *db.Client, worldID string) error {
	if worldID == "" {
		worldID = "1"
	}

	ctx := context.Background()

	// Check if races exist for this world
	if n, _ := client.Race.Query().Where(race.WorldID(worldID)).Count(ctx); n > 0 {
		log.Printf("Races already seeded for world %s, skipping", worldID)
		return nil
	}

	races := []struct {
		name           string
		displayName    string
		description    string
		statModifiers  map[string]int
		skillGrants    []string
		equipmentSlots []string
	}{
		{
			name:        "human",
			displayName: "Human",
			description: "A regular human survivor. Balanced and adaptable.",
			statModifiers: map[string]int{
				"strength":     0,
				"dexterity":    0,
				"constitution": 0,
				"intelligence": 0,
				"wisdom":       0,
			},
			skillGrants:    []string{},
			equipmentSlots: []string{"head", "neck", "chest", "back", "hands", "legs", "feet", "finger_left", "finger_right", "main_hand", "off_hand"},
		},
		{
			name:        "turtle",
			displayName: "Turtle",
			description: "A mutant turtle with a hard shell and slow but sturdy nature. +2 CON.",
			statModifiers: map[string]int{
				"strength":     0,
				"dexterity":    -1,
				"constitution": 2,
				"intelligence": 0,
				"wisdom":       0,
			},
			skillGrants:    []string{"shell_defense"},
			equipmentSlots: []string{"head", "neck", "chest", "back", "hands", "legs", "feet", "finger_left", "finger_right", "main_hand", "off_hand"},
		},
		{
			name:        "mutant",
			displayName: "Mutant",
			description: "A strange mutant with unusual abilities. +2 INT, +1 STR.",
			statModifiers: map[string]int{
				"strength":     1,
				"dexterity":    0,
				"constitution": 0,
				"intelligence": 2,
				"wisdom":       0,
			},
			skillGrants:    []string{"mutant_armor"},
			equipmentSlots: []string{"head", "neck", "chest", "back", "hands", "legs", "feet", "finger_left", "finger_right", "main_hand", "off_hand", "tail"},
		},
	}

	for _, r := range races {
		statJSON, _ := json.Marshal(r.statModifiers)
		skillJSON, _ := json.Marshal(r.skillGrants)
		_, err := client.Race.Create().
			SetWorldID(worldID).
			SetName(r.name).
			SetDisplayName(r.displayName).
			SetDescription(r.description).
			SetStatModifiers(string(statJSON)).
			SetSkillGrants(string(skillJSON)).
			SetEquipmentSlots(r.equipmentSlots).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to seed race %s for world %s: %v", r.name, worldID, err)
		}
	}

	log.Printf("Races seeded successfully for world %s", worldID)
	return nil
}

// InitGenders seeds the gender table for the specified world if empty.
// Default world_id is "1" for backward compatibility.
func InitGenders(client *db.Client, worldID string) error {
	if worldID == "" {
		worldID = "1"
	}

	ctx := context.Background()

	// Check if genders exist for this world
	if n, _ := client.Gender.Query().Where(gender.WorldID(worldID)).Count(ctx); n > 0 {
		log.Printf("Genders already seeded for world %s, skipping", worldID)
		return nil
	}

	genders := []struct {
		name              string
		displayName       string
		subjectPronoun    string
		objectPronoun     string
		possessivePronoun string
	}{
		{name: "he_him", displayName: "He/Him", subjectPronoun: "he", objectPronoun: "him", possessivePronoun: "his"},
		{name: "she_her", displayName: "She/Her", subjectPronoun: "she", objectPronoun: "her", possessivePronoun: "hers"},
		{name: "they_them", displayName: "They/Them", subjectPronoun: "they", objectPronoun: "them", possessivePronoun: "theirs"},
	}

	for _, g := range genders {
		_, err := client.Gender.Create().
			SetWorldID(worldID).
			SetName(g.name).
			SetDisplayName(g.displayName).
			SetSubjectPronoun(g.subjectPronoun).
			SetObjectPronoun(g.objectPronoun).
			SetPossessivePronoun(g.possessivePronoun).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to seed gender %s for world %s: %v", g.name, worldID, err)
		}
	}

	log.Printf("Genders seeded successfully for world %s", worldID)
	return nil
}

// GetPlayableRaces returns all races for the specified world where requirement_tags is empty.
func GetPlayableRaces(ctx context.Context, client *db.Client, worldID string) ([]*db.Race, error) {
	if worldID == "" {
		worldID = "1"
	}
	all, err := client.Race.Query().Where(race.WorldID(worldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*db.Race, 0, len(all))
	for _, r := range all {
		if len(r.RequirementTags) == 0 {
			result = append(result, r)
		}
	}
	return result, nil
}

// GetAllGenders returns all genders for the specified world.
func GetAllGenders(ctx context.Context, client *db.Client, worldID string) ([]*db.Gender, error) {
	if worldID == "" {
		worldID = "1"
	}
	return client.Gender.Query().Where(gender.WorldID(worldID)).All(ctx)
}

// ApplyRaceToCharacter reads the race from DB and applies stat modifiers.
// It mutates the character in-place and returns the updated character.
// Race lookup is scoped to the character's world using current_world field.
func ApplyRaceToCharacter(ctx context.Context, client *db.Client, char *db.Character) (*db.Character, error) {
	// Get the character's current world
	charObj, err := client.Character.Query().Where(character.ID(char.ID)).Only(ctx)
	if err != nil {
		return char, err
	}

	// Get world_id from character's current_world field, default to "1"
	worldID := charObj.CurrentWorld
	if worldID == "" {
		worldID = "1"
	}

	raceObj, err := client.Race.Query().Where(race.Name(char.Race), race.WorldID(worldID)).Only(ctx)
	if err != nil {
		return char, err
	}

	var statMods map[string]int
	if raceObj.StatModifiers != "" {
		json.Unmarshal([]byte(raceObj.StatModifiers), &statMods)
	}

	updater := client.Character.UpdateOneID(char.ID)
	if v, ok := statMods["strength"]; ok {
		updater.SetStrength(char.Strength + v)
	}
	if v, ok := statMods["dexterity"]; ok {
		updater.SetDexterity(char.Dexterity + v)
	}
	if v, ok := statMods["constitution"]; ok {
		updater.SetConstitution(char.Constitution + v)
	}
	if v, ok := statMods["intelligence"]; ok {
		updater.SetIntelligence(char.Intelligence + v)
	}
	if v, ok := statMods["wisdom"]; ok {
		updater.SetWisdom(char.Wisdom + v)
	}

	_, err = updater.Save(ctx)
	return char, err
}

// ApplyGenderToCharacter sets the gender field on a character.
// Gender lookup is scoped to the specified world.
func ApplyGenderToCharacter(ctx context.Context, client *db.Client, charID int, genderName string, worldID string) error {
	if worldID == "" {
		worldID = "1"
	}
	char, err := client.Character.Get(ctx, charID)
	if err != nil {
		return err
	}

	// Verify the gender exists in this world
	_, err = client.Gender.Query().Where(gender.Name(genderName), gender.WorldID(worldID)).Only(ctx)
	if err != nil {
		return err
	}

	_, err = client.Character.UpdateOneID(char.ID).SetGender(genderName).Save(ctx)
	return err
}
