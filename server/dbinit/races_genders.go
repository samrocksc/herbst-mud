package dbinit

import (
	"context"
	"encoding/json"
	"log"

	"herbst-server/db"
	"herbst-server/db/race"
)

// InitRaces seeds the race table if empty.
func InitRaces(client *db.Client) error {
	ctx := context.Background()

	if n, _ := client.Race.Query().Count(ctx); n > 0 {
		log.Println("Races already seeded, skipping")
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
				"constitution":  0,
				"intelligence":  0,
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
				"constitution":  2,
				"intelligence":  0,
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
				"constitution":  0,
				"intelligence":  2,
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
			SetName(r.name).
			SetDisplayName(r.displayName).
			SetDescription(r.description).
			SetStatModifiers(string(statJSON)).
			SetSkillGrants(string(skillJSON)).
			SetEquipmentSlots(r.equipmentSlots).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to seed race %s: %v", r.name, err)
		}
	}

	log.Println("Races seeded successfully")
	return nil
}

// InitGenders seeds the gender table if empty.
func InitGenders(client *db.Client) error {
	ctx := context.Background()

	if n, _ := client.Gender.Query().Count(ctx); n > 0 {
		log.Println("Genders already seeded, skipping")
		return nil
	}

	genders := []struct {
		name             string
		displayName      string
		subjectPronoun   string
		objectPronoun    string
		possessivePronoun string
	}{
		{name: "he_him", displayName: "He/Him", subjectPronoun: "he", objectPronoun: "him", possessivePronoun: "his"},
		{name: "she_her", displayName: "She/Her", subjectPronoun: "she", objectPronoun: "her", possessivePronoun: "hers"},
		{name: "they_them", displayName: "They/Them", subjectPronoun: "they", objectPronoun: "them", possessivePronoun: "theirs"},
	}

	for _, g := range genders {
		_, err := client.Gender.Create().
			SetName(g.name).
			SetDisplayName(g.displayName).
			SetSubjectPronoun(g.subjectPronoun).
			SetObjectPronoun(g.objectPronoun).
			SetPossessivePronoun(g.possessivePronoun).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to seed gender %s: %v", g.name, err)
		}
	}

	log.Println("Genders seeded successfully")
	return nil
}

// GetPlayableRaces returns all races where requirement_tags is empty.
func GetPlayableRaces(ctx context.Context, client *db.Client) ([]*db.Race, error) {
	all, err := client.Race.Query().All(ctx)
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

// GetAllGenders returns all genders.
func GetAllGenders(ctx context.Context, client *db.Client) ([]*db.Gender, error) {
	return client.Gender.Query().All(ctx)
}

// ApplyRaceToCharacter reads the race from DB and applies stat modifiers.
// It mutates the character in-place and returns the updated character.
func ApplyRaceToCharacter(ctx context.Context, client *db.Client, char *db.Character) (*db.Character, error) {
	raceObj, err := client.Race.Query().Where(race.NameEQ(char.Race)).Only(ctx)
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
func ApplyGenderToCharacter(ctx context.Context, client *db.Client, charID int, genderName string) error {
	char, err := client.Character.Get(ctx, charID)
	if err != nil {
		return err
	}
	_, err = client.Character.UpdateOneID(char.ID).SetGender(genderName).Save(ctx)
	return err
}