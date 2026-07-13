package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"herbst-server/constants"
	"herbst-server/db"
	"herbst-server/db/ability"
	"herbst-server/db/factioncategory"
	"herbst-server/db/room"
	"herbst-server/db/schema"
	"herbst-server/dbinit"
	"herbst-server/repository"

	"golang.org/x/crypto/bcrypt"
)

type characterService struct {
	repos   *repository.Container
	client  *db.Client
}

func NewCharacterService(client *db.Client, repos *repository.Container) CharacterService {
	return &characterService{repos: repos, client: client}
}

var (
	ErrCharacterNotFound = errors.New("character not found")
	ErrCharacterNameTaken = errors.New("character name already exists")
	ErrTooManyCharacters  = errors.New("maximum of 3 characters per user reached")
	ErrInvalidRace        = errors.New("invalid race")
	ErrInvalidGender      = errors.New("invalid gender")
	ErrWorldNotReady      = errors.New("world is not ready for character creation")
)

type CreateCharacterInput struct {
	UserID   int
	Name     string
	Race     string
	Gender   string
	Class    string
	Description string
	WorldID  string
	Factions []string
}

func (s *characterService) CreateCharacter(ctx context.Context, input CreateCharacterInput) (*db.Character, error) {
	if err := s.CheckWorldReady(ctx, input.WorldID); err != nil {
		return nil, err
	}
	if len(input.Name) < 1 || len(input.Name) > 23 {
		return nil, errors.New("character name must be 1-23 characters")
	}
	for _, ch := range input.Name {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
			return nil, errors.New("character name can only contain letters (a-z, A-Z)")
		}
	}
	existingChar, err := s.repos.Character.GetByName(ctx, input.Name)
	if err == nil && existingChar != nil {
		return nil, ErrCharacterNameTaken
	}
	if input.UserID > 0 {
		userChars, err := s.repos.Character.CountByUser(ctx, input.UserID)
		if err == nil && userChars >= 3 {
			return nil, ErrTooManyCharacters
		}
	}
	// Race is REQUIRED — validate it exists in the world
	if input.Race == "" {
		return nil, ErrInvalidRace
	}
	raceObj, err := s.repos.Race.GetByName(ctx, input.Race, input.WorldID)
	if err != nil || len(raceObj.RequirementTags) > 0 {
		return nil, ErrInvalidRace
	}
	race := raceObj.Name
	// Gender is REQUIRED — validate it exists in the world
	if input.Gender == "" {
		return nil, ErrInvalidGender
	}
	genderObj, err := s.repos.Gender.GetByWorld(ctx, input.Gender, input.WorldID)
	if err != nil {
		return nil, ErrInvalidGender
	}
	gen := genderObj.Name
	// Class defaults to "survivor" if not provided or invalid (backwards compatibility).
	// Validate against DB: a faction with name=class in a "class" category for this world.
	class := "survivor"
	specialty := "generalist"
	if input.Class != "" {
		if classFaction, err := s.getClassFactionByName(ctx, input.Class, input.WorldID); err == nil && classFaction != nil {
			class = input.Class
			if len(classFaction.Specialties) > 0 {
				specialty = classFaction.Specialties[0].ID
			}
		}
	}
	// Get stat bonuses from the DB faction for this class.
	var statBonuses schema.StatBonuses
	if classFaction, err := s.getClassFactionByName(ctx, class, input.WorldID); err == nil && classFaction != nil {
		statBonuses = classFaction.StatBonuses
	}
	baseStrength := constants.DefaultStats.Strength + statBonuses.Strength
	baseDexterity := constants.DefaultStats.Dexterity + statBonuses.Dexterity
	baseConstitution := constants.DefaultStats.Constitution + statBonuses.Constitution
	baseIntelligence := constants.DefaultStats.Intelligence + statBonuses.Intelligence
	baseWisdom := constants.DefaultStats.Wisdom + statBonuses.Wisdom
	baseCharisma := constants.DefaultStats.Charisma + statBonuses.Charisma
	worldIntID := 0
	if input.WorldID != "" {
		worldObj, err := s.repos.World.GetByName(ctx, input.WorldID)
		if err == nil && worldObj != nil {
			worldIntID = worldObj.ID
		}
	}
	startingRoomID := 1
	// Prefer a root room in the requested world. If none exists, fall back to
	// the lowest-id room in that world. Final fallback: room id 1 (dev world).
	if input.WorldID != "" {
		if rootRooms, rerr := s.client.Room.Query().
			Where(room.IsRootRoom(true), room.WorldID(input.WorldID)).
			All(ctx); rerr == nil && len(rootRooms) > 0 {
			startingRoomID = rootRooms[0].ID
		} else if worldRooms, werr := s.client.Room.Query().
			Where(room.WorldID(input.WorldID)).
			Limit(1).
			All(ctx); werr == nil && len(worldRooms) > 0 {
			startingRoomID = worldRooms[0].ID
		}
	} else if rootRooms, rerr := s.client.Room.Query().
		Where(room.IsRootRoom(true)).
		All(ctx); rerr == nil && len(rootRooms) > 0 {
		startingRoomID = rootRooms[0].ID
	}
	char, err := s.repos.Character.Create(ctx, repository.CreateCharacterInput{
		Name:            input.Name,
		UserID:          input.UserID,
		RoomID:          startingRoomID,
		StartingRoomID:  startingRoomID,
		RespawnRoomID:   startingRoomID,
		WorldID:         input.WorldID,
		WorldIntID:      worldIntID,
		IsNPC:           false,
		Race:            race,
		Gender:          gen,
		Description:     input.Description,
		Class:           class,
		Specialty:       specialty,
		HP:              100,
		MaxHP:           100,
		Stamina:         50,
		MaxStamina:      50,
		Mana:            50,
		MaxMana:         50,
		Level:           1,
		XP:              0,
		Strength:        baseStrength,
		Dexterity:       baseDexterity,
		Constitution:    baseConstitution,
		Intelligence:    baseIntelligence,
		Wisdom:          baseWisdom,
		Charisma:        baseCharisma,
	})
	if err != nil {
		return nil, err
	}
	char, _ = dbinit.ApplyRaceToCharacter(ctx, s.client, char)
	if grantErr := s.GrantTag(ctx, char.ID, "first_class", "system"); grantErr != nil {
		slog.Warn("failed to grant first_class tag", slog.Int("character_id", char.ID), slog.String("error", grantErr.Error()))
	}
	if syncErr := s.SyncRaceTags(ctx, char.ID, char.Race, char.CurrentWorld); syncErr != nil {
		slog.Warn("failed to sync race tags", slog.Int("character_id", char.ID), slog.String("error", syncErr.Error()))
	}
	// Equip default abilities for fresh characters (max 4 slots)
	s.equipDefaultAbilities(ctx, char.ID, class, input.WorldID)
	// Add initial faction memberships
	for _, factionStr := range input.Factions {
		factionID := 0
		fmt.Sscanf(factionStr, "%d", &factionID)
		if factionID > 0 {
			if _, cfErr := s.repos.CharacterFaction.Create(ctx, char.ID, factionID, 0); cfErr != nil {
				slog.Warn("failed to add faction to character", slog.Int("character_id", char.ID), slog.Int("faction_id", factionID), slog.String("error", cfErr.Error()))
			}
		}
	}
	return char, nil
}

func (s *characterService) CheckWorldReady(ctx context.Context, worldID string) error {
	roomCount, err := s.repos.Room.CountByWorld(ctx, worldID)
	if err != nil || roomCount < 1 {
		return ErrWorldNotReady
	}
	// Require at least one root room — players need a known entry point.
	rootRoomCount, err := s.client.Room.Query().
		Where(room.IsRootRoom(true), room.WorldID(worldID)).
		Count(ctx)
	if err != nil || rootRoomCount < 1 {
		return ErrWorldNotReady
	}
	raceCount, err := s.repos.Race.CountByWorld(ctx, worldID)
	if err != nil || raceCount < 1 {
		return ErrWorldNotReady
	}
	genderCount, err := s.repos.Gender.CountByWorld(ctx, worldID)
	if err != nil || genderCount < 1 {
		return ErrWorldNotReady
	}
	return nil
}

func (s *characterService) DeleteCharacter(ctx context.Context, charID int) error {
	return s.repos.Character.Delete(ctx, charID)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *characterService) GrantTag(ctx context.Context, characterID int, tag, source string) error {
	_, err := s.repos.CharacterTag.Create(ctx, characterID, tag, source)
	return err
}

func (s *characterService) SyncRaceTags(ctx context.Context, characterID int, raceName, worldID string) error {
	raceObj, err := s.repos.Race.GetByName(ctx, raceName, worldID)
	if err != nil {
		return err
	}
	// Delete existing race-source tags for this character
	tags, err := s.repos.Character.QueryTags(ctx, characterID)
	if err != nil {
		return err
	}
	for _, t := range tags {
		if t.Source == "race" {
			_ = s.repos.CharacterTag.Delete(ctx, t.ID)
		}
	}
	for _, t := range raceObj.Edges.Tags {
		if grantErr := s.GrantTag(ctx, characterID, t.Name, "race"); grantErr != nil {
			return grantErr
		}
	}
	return nil
}

func (s *characterService) QueryCharacterByName(ctx context.Context, name string) (*db.Character, error) {
	return s.repos.Character.GetByName(ctx, name)
}

// getClassFactionByName queries the DB for a faction with the given name in a
// "class" faction category for the given world. Returns nil if not found.
func (s *characterService) getClassFactionByName(ctx context.Context, name, worldID string) (*db.Faction, error) {
	cat, err := s.client.FactionCategory.Query().
		Where(
			factioncategory.Name("class"),
			factioncategory.WorldID(worldID),
		).
		WithFactions().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	for _, f := range cat.Edges.Factions {
		if f.Name == name {
			return f, nil
		}
	}
	return nil, fmt.Errorf("class faction %q not found in world %q", name, worldID)
}

// equipDefaultAbilities queries abilities for the character's class and creates
// character_abilities rows (max 4, starting at slot 1). Falls back to classless
// abilities if none found for the given class.
func (s *characterService) equipDefaultAbilities(ctx context.Context, charID int, class, worldID string) {
	// Pick the first N active abilities in the world. There's no per-character-class
	// ability table (ability_class is active/passive/classless, not warrior/chef/etc).
	// Default to active abilities so the player can use slots 1-4 on a fresh char.
	abilities, err := s.client.Ability.Query().
		Where(ability.WorldID(worldID), ability.AbilityClassEQ("active")).
		Limit(4).
		All(ctx)
	if err != nil || len(abilities) == 0 {
		slog.Warn("no default abilities found for character", slog.Int("character_id", charID), slog.String("class", class), slog.String("world_id", worldID))
		return
	}
	for slot, ab := range abilities {
		slotNum := slot + 1
		_, err := s.repos.CharacterAbility.Create(ctx, charID, ab.ID, slotNum)
		if err != nil {
			slog.Warn("failed to equip default ability", slog.Int("character_id", charID), slog.Int("ability_id", ab.ID), slog.String("error", err.Error()))
		}
	}
}