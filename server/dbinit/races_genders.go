package dbinit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"herbst-server/db"
	"herbst-server/db/gameconfig"
	"herbst-server/db/race"
	"herbst-server/db/room"
)

// InitRaces seeds the race table if empty.
func InitRaces(client *db.Client) error {
	ctx := context.Background()

	if n, _ := client.Race.Query().Count(ctx); n > 0 {
		log.Println("Races already seeded, skipping")
		return nil
	}

	races := []struct {
		name          string
		displayName   string
		description   string
		statModifiers map[string]int
		skillGrants   []string
		isPlayable    bool
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
			skillGrants: []string{},
			isPlayable:  true,
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
			skillGrants: []string{"shell_defense"},
			isPlayable:  true,
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
			skillGrants: []string{"mutant_armor"},
			isPlayable:  true,
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
			SetIsPlayable(r.isPlayable).
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

// GetFountainRoomID returns the configured fountain room ID from game_config,
// falling back to a room named "The Fountain" if not set.
func GetFountainRoomID(ctx context.Context, client *db.Client) (int, error) {
	cfg, err := client.GameConfig.Query().Where(gameconfig.KeyEQ("fountain_room_id")).Only(ctx)
	if err == nil && cfg != nil {
		var id int
		if _, scanErr := fmt.Sscanf(cfg.Value, "%d", &id); scanErr == nil {
			return id, nil
		}
	}
	// Fallback: look up by room name
	fountain, err := client.Room.Query().Where(room.NameEQ("The Fountain")).Only(ctx)
	if err != nil {
		return 0, err
	}
	return fountain.ID, nil
}

// SetFountainRoomID stores the fountain room ID in game_config.
func SetFountainRoomID(ctx context.Context, client *db.Client, roomID int) error {
	ctx_ := context.Background()
	key := "fountain_room_id"

	existing, err := client.GameConfig.Query().Where(gameconfig.KeyEQ(key)).Only(ctx_)
	if err == nil && existing != nil {
		_, err = client.GameConfig.UpdateOne(existing).SetValue(fmt.Sprintf("%d", roomID)).Save(ctx_)
		return err
	}

	_, err = client.GameConfig.Create().SetKey(key).SetValue(fmt.Sprintf("%d", roomID)).Save(ctx_)
	return err
}

// GetPlayableRaces returns all races where is_playable = true.
func GetPlayableRaces(ctx context.Context, client *db.Client) ([]*db.Race, error) {
	return client.Race.Query().Where(race.IsPlayable(true)).All(ctx)
}

// GetAllGenders returns all genders.
func GetAllGenders(ctx context.Context, client *db.Client) ([]*db.Gender, error) {
	return client.Gender.Query().All(ctx)
}

// ApplyRaceToCharacter reads the race from DB and applies stat modifiers + skill grants.
func ApplyRaceToCharacter(ctx context.Context, client *db.Client, charID int, raceName string) error {
	raceObj, err := client.Race.Query().Where(race.NameEQ(raceName)).Only(ctx)
	if err != nil {
		return err
	}

	char, err := client.Character.Get(ctx, charID)
	if err != nil {
		return err
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
	return err
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