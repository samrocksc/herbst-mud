package services

import (
	"context"
	"errors"
	"fmt"

	"herbst-server/constants"
	"herbst-server/db"
	"herbst-server/db/character"
	genderpkg "herbst-server/db/gender"
	racepkg "herbst-server/db/race"
	"herbst-server/db/room"
	"herbst-server/db/user"
	"herbst-server/dbinit"

	"golang.org/x/crypto/bcrypt"
)

// CharacterService handles character business logic.
type CharacterService struct {
	client *db.Client
}

// NewCharacterService creates a new CharacterService.
func NewCharacterService(client *db.Client) *CharacterService {
	return &CharacterService{client: client}
}

// ErrCharacterNotFound is returned when a character doesn't exist.
var ErrCharacterNotFound = errors.New("character not found")

// ErrCharacterNameTaken is returned when the character name is already in use.
var ErrCharacterNameTaken = errors.New("character name already exists")

// ErrTooManyCharacters is returned when the user already has 3 characters.
var ErrTooManyCharacters = errors.New("maximum of 3 characters per user reached")

// ErrInvalidRace is returned when the race is invalid or not playable.
var ErrInvalidRace = errors.New("invalid race")

// ErrInvalidGender is returned when the gender is invalid.
var ErrInvalidGender = errors.New("invalid gender")

// CreateCharacterInput contains the data needed to create a new character.
type CreateCharacterInput struct {
	UserID   int
	Name     string
	Password string
	Class    string
	Race     string
	Gender   string
}

// CreateCharacter creates a new player character and returns it.
// It validates name uniqueness, race, gender, and enforces max 3 chars per user.
// Password must already be hashed by the caller.
func (s *CharacterService) CreateCharacter(ctx context.Context, input CreateCharacterInput) (*db.Character, error) {
	// Validate character name: 1-23 chars, letters only
	if len(input.Name) < 1 || len(input.Name) > 23 {
		return nil, errors.New("character name must be 1-23 characters")
	}
	for _, ch := range input.Name {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
			return nil, errors.New("character name can only contain letters (a-z, A-Z)")
		}
	}

	// Check name uniqueness
	existingChar, err := s.client.Character.Query().
		Where(character.NameEQ(input.Name)).
		Only(ctx)
	if err == nil && existingChar != nil {
		return nil, ErrCharacterNameTaken
	}

	// Enforce max 3 characters per user
	userChars, err := s.client.Character.Query().
		Where(character.HasUserWith(user.IDEQ(input.UserID))).
		Count(ctx)
	if err == nil && userChars >= 3 {
		return nil, ErrTooManyCharacters
	}

	// Resolve race (default to human)
	race := "human"
	if input.Race != "" {
		raceObj, err := s.client.Race.Query().Where(racepkg.NameEQ(input.Race)).Only(ctx)
		if err != nil || !raceObj.IsPlayable {
			return nil, ErrInvalidRace
		}
		race = raceObj.Name
	}

	// Resolve gender (default to he_him)
	gen := "he_him"
	if input.Gender != "" {
		genderObj, err := s.client.Gender.Query().Where(genderpkg.NameEQ(input.Gender)).Only(ctx)
		if err != nil {
			return nil, ErrInvalidGender
		}
		gen = genderObj.Name
	}

	// Set class (default to survivor)
	class := "survivor"
	if input.Class != "" {
		class = input.Class
	}

	// Get class configuration
	classConfig := constants.GetClassConfig(class, "")

	// Calculate base stats with class bonuses
	baseStrength := constants.DefaultStats.Strength + classConfig.StatBonuses.Strength
	baseDexterity := constants.DefaultStats.Dexterity + classConfig.StatBonuses.Dexterity
	baseConstitution := constants.DefaultStats.Constitution + classConfig.StatBonuses.Constitution
	baseIntelligence := constants.DefaultStats.Intelligence + classConfig.StatBonuses.Intelligence
	baseWisdom := constants.DefaultStats.Wisdom + classConfig.StatBonuses.Wisdom

	// Get starting room (default to first starting room)
	startingRoomID := 1
	startingRooms, err := s.client.Room.Query().Where(room.IsStartingRoom(true)).All(ctx)
	if err == nil && len(startingRooms) > 0 {
		startingRoomID = startingRooms[0].ID
	}

	// Build character
	builder := s.client.Character.Create().
		SetName(input.Name).
		SetPassword(input.Password).
		SetUserID(input.UserID).
		SetIsNPC(false).
		SetHitpoints(100).
		SetMaxHitpoints(100).
		SetStamina(50).
		SetMaxStamina(50).
		SetMana(50).
		SetMaxMana(50).
		SetCurrentRoomId(startingRoomID).
		SetStartingRoomId(startingRoomID).
		SetRace(race).
		SetGender(gen).
		SetClass(class).
		SetSpecialty(classConfig.Specialty).
		SetStrength(baseStrength).
		SetDexterity(baseDexterity).
		SetConstitution(baseConstitution).
		SetIntelligence(baseIntelligence).
		SetWisdom(baseWisdom).
		SetLevel(1).
		SetXp(0)

	// Apply starting skills from class config
	for skill, level := range classConfig.StartingSkills {
		switch skill {
		case "blades":
			builder.SetSkillBlades(level)
		case "staves":
			builder.SetSkillStaves(level)
		case "knives":
			builder.SetSkillKnives(level)
		case "martial":
			builder.SetSkillMartial(level)
		case "brawling":
			builder.SetSkillBrawling(level)
		case "tech":
			builder.SetSkillTech(level)
		}
	}

	char, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Apply race stat modifiers from DB
	char, _ = dbinit.ApplyRaceToCharacter(ctx, s.client, char)

	// Auto-grant first_class tag on character creation
	if grantErr := s.GrantTag(ctx, char.ID, "first_class", "system"); grantErr != nil {
		// Log but don't fail character creation if tag grant fails
		fmt.Printf("Warning: failed to grant first_class tag to character %d: %v\n", char.ID, grantErr)
	}

	return char, nil
}

// DeleteCharacter deletes a character and its related data.
func (s *CharacterService) DeleteCharacter(ctx context.Context, charID int) error {
	char, err := s.client.Character.Get(ctx, charID)
	if err != nil {
		return ErrCharacterNotFound
	}
	return s.client.Character.DeleteOne(char).Exec(ctx)
}

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// GrantTag adds a tag to a character with the given source.
func (s *CharacterService) GrantTag(ctx context.Context, characterID int, tag, source string) error {
	_, err := s.client.CharacterTag.Create().
		SetTag(tag).
		SetSource(source).
		SetCharacterID(characterID).
		Save(ctx)
	return err
}
