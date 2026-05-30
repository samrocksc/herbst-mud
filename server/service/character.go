package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"herbst-server/constants"
	"herbst-server/db"
	"herbst-server/db/room"
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
)

type CreateCharacterInput struct {
	UserID   int
	Name     string
	Race     string
	Gender   string
	Description string
	WorldID  string
	Factions []string
}

func (s *characterService) CreateCharacter(ctx context.Context, input CreateCharacterInput) (*db.Character, error) {
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
	race := "human"
	if input.Race != "" {
		// Use default world for character creation (or the character's world)
		raceObj, err := s.repos.Race.GetByName(ctx, input.Race, input.WorldID)
		if err != nil || len(raceObj.RequirementTags) > 0 {
			return nil, ErrInvalidRace
		}
		race = raceObj.Name
	}
	gen := "he_him"
	if input.Gender != "" {
		// Use default world for character creation
		genderObj, err := s.repos.Gender.GetByWorld(ctx, input.Gender, input.WorldID)
		if err != nil {
			return nil, ErrInvalidGender
		}
		gen = genderObj.Name
	}
	// Class is derived from faction memberships in the "class" category.
	// Default to "survivor" if no class faction found (backwards compatibility).
	class := "survivor"
	classConfig := constants.GetClassConfig(class, "")
	baseStrength := constants.DefaultStats.Strength + classConfig.StatBonuses.Strength
	baseDexterity := constants.DefaultStats.Dexterity + classConfig.StatBonuses.Dexterity
	baseConstitution := constants.DefaultStats.Constitution + classConfig.StatBonuses.Constitution
	baseIntelligence := constants.DefaultStats.Intelligence + classConfig.StatBonuses.Intelligence
	baseWisdom := constants.DefaultStats.Wisdom + classConfig.StatBonuses.Wisdom
	startingRoomID := 1
	rootRooms, err := s.client.Room.Query().Where(room.IsRootRoom(true)).All(ctx)
	if err == nil && len(rootRooms) > 0 {
		startingRoomID = rootRooms[0].ID
	}
	char, err := s.repos.Character.Create(ctx, repository.CreateCharacterInput{
		Name:            input.Name,
		UserID:          input.UserID,
		RoomID:          startingRoomID,
		StartingRoomID:  startingRoomID,
		RespawnRoomID:   startingRoomID,
		WorldID:         input.WorldID,
		IsNPC:           false,
		Race:            race,
		Gender:          gen,
		Description:     input.Description,
		Class:           class,
		Specialty:       classConfig.Specialty,
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