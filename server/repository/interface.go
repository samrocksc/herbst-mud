package repository

import (
	"context"
	"time"

	"herbst-server/db"
	"herbst-server/db/quest"
	"herbst-server/db/questprogress"
	"herbst-server/db/schema"
)

// CharacterRepo defines data access for characters.
type CharacterRepo interface {
	Get(ctx context.Context, id int) (*db.Character, error)
	GetByName(ctx context.Context, name string) (*db.Character, error)
	ListByUser(ctx context.Context, userID int) ([]*db.Character, error)
	ListByRoom(ctx context.Context, roomID int) ([]*db.Character, error)
	ListNPCsByRoom(ctx context.Context, roomID int) ([]*db.Character, error)
	ListAllNPCs(ctx context.Context) ([]*db.Character, error)
	ListAll(ctx context.Context) ([]*db.Character, error)
	ListAllByWorld(ctx context.Context, worldID string) ([]*db.Character, error)
	CountByUser(ctx context.Context, userID int) (int, error)
	Create(ctx context.Context, input CreateCharacterInput) (*db.Character, error)
	Update(ctx context.Context, id int, updates CharacterUpdates) (*db.Character, error)
	Delete(ctx context.Context, id int) error
	QueryAbilities(ctx context.Context, charID int) ([]*db.CharacterAbility, error)
	QueryAbilitiesWithDetails(ctx context.Context, charID int) ([]*db.CharacterAbility, error)
	QueryTags(ctx context.Context, charID int) ([]*db.CharacterTag, error)
	QueryFactions(ctx context.Context, charID int) ([]*db.CharacterFaction, error)
	QueryActiveEffects(ctx context.Context, charID int) ([]*db.ActiveEffect, error)
	QueryQuestProgress(ctx context.Context, charID int) ([]*db.QuestProgress, error)
}

// RoomRepo defines data access for rooms.
type RoomRepo interface {
	Get(ctx context.Context, id int) (*db.Room, error)
	List(ctx context.Context, worldID string) ([]*db.Room, error)
	GetRoot(ctx context.Context) ([]*db.Room, error)
	Create(ctx context.Context, input CreateRoomInput) (*db.Room, error)
	Update(ctx context.Context, id int, updates RoomUpdates) (*db.Room, error)
	Delete(ctx context.Context, id int) error
}

// QuestRepo defines data access for quest definitions.
type QuestRepo interface {
	Get(ctx context.Context, id int) (*db.Quest, error)
	List(ctx context.Context, worldID string) ([]*db.Quest, error)
	Create(ctx context.Context, input CreateQuestInput) (*db.Quest, error)
	Update(ctx context.Context, id int, updates QuestUpdates) (*db.Quest, error)
	Delete(ctx context.Context, id int) error
}

// QuestProgressRepo defines data access for quest progress.
type QuestProgressRepo interface {
	Get(ctx context.Context, id int) (*db.QuestProgress, error)
	GetWithRelations(ctx context.Context, id int) (*db.QuestProgress, error)
	ListByCharacter(ctx context.Context, charID int) ([]*db.QuestProgress, error)
	Create(ctx context.Context, input CreateQuestProgressInput) (*db.QuestProgress, error)
	Update(ctx context.Context, id int, updates QuestProgressUpdates) (*db.QuestProgress, error)
	Delete(ctx context.Context, id int) error
	CountActiveByCharacter(ctx context.Context, charID int, questID int) (int, error)
}

// EquipmentRepo defines data access for equipment instances.
type EquipmentRepo interface {
	Get(ctx context.Context, id int) (*db.Equipment, error)
	ListByOwner(ctx context.Context, ownerID int) ([]*db.Equipment, error)
	ListByRoom(ctx context.Context, roomID int) ([]*db.Equipment, error)
	Create(ctx context.Context, input CreateEquipmentInput) (*db.Equipment, error)
	Update(ctx context.Context, id int, updates EquipmentUpdates) (*db.Equipment, error)
	Delete(ctx context.Context, id int) error
	CountByTemplateID(ctx context.Context, templateID int) (int, error)
}

// NPCTemplateRepo defines data access for NPC templates.
type NPCTemplateRepo interface {
	Get(ctx context.Context, id string) (*db.NPCTemplate, error)
	List(ctx context.Context, worldID string) ([]*db.NPCTemplate, error)
	Create(ctx context.Context, input CreateNPCTemplateInput) (*db.NPCTemplate, error)
	Update(ctx context.Context, id string, updates NPCTemplateUpdates) (*db.NPCTemplate, error)
	Delete(ctx context.Context, id string) error
}

// AbilityRepo defines data access for abilities.
type AbilityRepo interface {
	Get(ctx context.Context, id int) (*db.Ability, error)
	List(ctx context.Context, worldID string) ([]*db.Ability, error)
	ListClassless(ctx context.Context, worldID string) ([]*db.Ability, error)
	ListByClass(ctx context.Context, worldID string, class string) ([]*db.Ability, error)
	Create(ctx context.Context, input CreateAbilityInput) (*db.Ability, error)
	Update(ctx context.Context, id int, updates AbilityUpdates) (*db.Ability, error)
	Delete(ctx context.Context, id int) error
}

// CharacterAbilityRepo defines data access for character-ability links.
type CharacterAbilityRepo interface {
	ListByCharacter(ctx context.Context, charID int) ([]*db.CharacterAbility, error)
	ListByCharacterWithDetails(ctx context.Context, charID int) ([]*db.CharacterAbility, error)
	ListByCharacterAndSlot(ctx context.Context, charID int, slot int) ([]*db.CharacterAbility, error)
	CountByCharacter(ctx context.Context, charID int) (int, error)
	ExistsByCharacterAndAbility(ctx context.Context, charID int, abilityID int) (bool, error)
	Create(ctx context.Context, charID int, abilityID int, slot int) (*db.CharacterAbility, error)
	Delete(ctx context.Context, id int) error
	DeleteByCharacterAndSlot(ctx context.Context, charID int, slot int) error
	DeleteByCharacterAndAbility(ctx context.Context, charID int, abilityID int) error
}

// TriggerRepo defines data access for triggers.
type TriggerRepo interface {
	Get(ctx context.Context, id int) (*db.Trigger, error)
	GetWithEdges(ctx context.Context, id int) (*db.Trigger, error)
	List(ctx context.Context) ([]*db.Trigger, error)
	ListWithEdges(ctx context.Context) ([]*db.Trigger, error)
	ListByRoom(ctx context.Context, roomID int) ([]*db.Trigger, error)
	ListByEquipment(ctx context.Context, equipmentID int) ([]*db.Trigger, error)
	ListByTriggerType(ctx context.Context, triggerType string) ([]*db.Trigger, error)
	ListByTargetType(ctx context.Context, targetType string) ([]*db.Trigger, error)
	Create(ctx context.Context, input CreateTriggerInput) (*db.Trigger, error)
	Update(ctx context.Context, id int, updates TriggerUpdates) (*db.Trigger, error)
	Delete(ctx context.Context, id int) error
}

// EffectRepo defines data access for effects.
type EffectRepo interface {
	Get(ctx context.Context, id int) (*db.Effect, error)
	GetWithHooks(ctx context.Context, id int) (*db.Effect, error)
	List(ctx context.Context) ([]*db.Effect, error)
	ListWithHooks(ctx context.Context) ([]*db.Effect, error)
	Create(ctx context.Context, input CreateEffectInput) (*db.Effect, error)
	Update(ctx context.Context, id int, updates EffectUpdates) (*db.Effect, error)
	Delete(ctx context.Context, id int) error
}

// ActiveEffectRepo defines data access for active effects on characters.
type ActiveEffectRepo interface {
	ListByCharacter(ctx context.Context, charID int) ([]*db.ActiveEffect, error)
	ListActiveByCharacter(ctx context.Context, charID int) ([]*db.ActiveEffect, error)
	GetActiveByCharacterAndEffect(ctx context.Context, charID, effectID int) (*db.ActiveEffect, error)
	GetWithEffect(ctx context.Context, id int) (*db.ActiveEffect, error)
	Create(ctx context.Context, input CreateActiveEffectInput) (*db.ActiveEffect, error)
	Update(ctx context.Context, id int, updates ActiveEffectUpdates) (*db.ActiveEffect, error)
	Delete(ctx context.Context, id int) error
	DeactivateExpired(ctx context.Context) ([]*db.ActiveEffect, error)
}

// EffectHookRepo defines data access for effect hooks.
type EffectHookRepo interface {
	Get(ctx context.Context, id int) (*db.EffectHook, error)
	GetWithEdges(ctx context.Context, id int) (*db.EffectHook, error)
	List(ctx context.Context) ([]*db.EffectHook, error)
	ListWithEdges(ctx context.Context) ([]*db.EffectHook, error)
	ListByEvent(ctx context.Context, event string) ([]*db.EffectHook, error)
	ListByTemplateWithEdges(ctx context.Context, templateID string) ([]*db.EffectHook, error)
	CountByEffect(ctx context.Context, effectID int) (int, error)
	Create(ctx context.Context, input CreateEffectHookInput) (*db.EffectHook, error)
	Update(ctx context.Context, id int, updates EffectHookUpdates) (*db.EffectHook, error)
	Delete(ctx context.Context, id int) error
}

// DialogNodeRepo defines data access for dialog nodes.
type DialogNodeRepo interface {
	Get(ctx context.Context, id string) (*db.DialogNode, error)
	List(ctx context.Context, worldID string) ([]*db.DialogNode, error)
	ListByTemplate(ctx context.Context, templateID string) ([]*db.DialogNode, error)
	Create(ctx context.Context, input CreateDialogNodeInput) (*db.DialogNode, error)
	Update(ctx context.Context, id string, updates DialogNodeUpdates) (*db.DialogNode, error)
	Delete(ctx context.Context, id string) error
}

// UserRepo defines data access for users.
type UserRepo interface {
	Get(ctx context.Context, id int) (*db.User, error)
	GetByEmail(ctx context.Context, email string) (*db.User, error)
	List(ctx context.Context) ([]*db.User, error)
	Create(ctx context.Context, input CreateUserInput) (*db.User, error)
	Update(ctx context.Context, id int, updates UserUpdates) (*db.User, error)
	Delete(ctx context.Context, id int) error
}

// FactionRepo defines data access for factions.
type FactionRepo interface {
	Get(ctx context.Context, id int) (*db.Faction, error)
	GetWithEdges(ctx context.Context, id int) (*db.Faction, error)
	List(ctx context.Context, worldID string) ([]*db.Faction, error)
	Create(ctx context.Context, input CreateFactionInput) (*db.Faction, error)
	Update(ctx context.Context, id int, updates FactionUpdates) (*db.Faction, error)
	Delete(ctx context.Context, id int) error
}

// CharacterFactionRepo defines data access for character-faction links.
type CharacterFactionRepo interface {
	ListByCharacter(ctx context.Context, charID int) ([]*db.CharacterFaction, error)
	ListByFactionWithDetails(ctx context.Context, factionID int) ([]*db.CharacterFaction, error)
	Create(ctx context.Context, charID int, factionID int, reputation int) (*db.CharacterFaction, error)
	Delete(ctx context.Context, id int) error
}

// CharacterTagRepo defines data access for character tags.
type CharacterTagRepo interface {
	ListByCharacter(ctx context.Context, charID int) ([]*db.CharacterTag, error)
	Create(ctx context.Context, charID int, tag, source string) (*db.CharacterTag, error)
	Delete(ctx context.Context, id int) error
}

// TagRepo defines data access for tags.
type TagRepo interface {
	Get(ctx context.Context, id int) (*db.Tag, error)
	GetByName(ctx context.Context, name, worldID string) (*db.Tag, error)
	List(ctx context.Context, worldID string) ([]*db.Tag, error)
	Create(ctx context.Context, input CreateTagInput) (*db.Tag, error)
	Update(ctx context.Context, id int, updates TagUpdates) (*db.Tag, error)
	Delete(ctx context.Context, id int) error
}

// AchievementRepo defines data access for achievements.
type AchievementRepo interface {
	Get(ctx context.Context, id int) (*db.Achievement, error)
	List(ctx context.Context) ([]*db.Achievement, error)
	Create(ctx context.Context, input CreateAchievementInput) (*db.Achievement, error)
	Update(ctx context.Context, id int, updates AchievementUpdates) (*db.Achievement, error)
	Delete(ctx context.Context, id int) error
}

// DamageLogRepo defines data access for damage logs.
type DamageLogRepo interface {
	Create(ctx context.Context, attackerID, targetID, damage int) (*db.DamageLog, error)
	ListByCharacter(ctx context.Context, charID int, limit int) ([]*db.DamageLog, error)
}

// GameConfigRepo defines data access for game configuration.
type GameConfigRepo interface {
	Get(ctx context.Context, key string) (*db.GameConfig, error)
	GetOrCreate(ctx context.Context, key, defaultValue string) (*db.GameConfig, error)
	Set(ctx context.Context, key, value string) (*db.GameConfig, error)
	List(ctx context.Context) ([]*db.GameConfig, error)
	Delete(ctx context.Context, key string) error
}

// CompetencyRepo defines data access for competency categories and thresholds.
type CompetencyRepo interface {
	GetCategory(ctx context.Context, id string) (*db.CompetencyCategory, error)
	GetCategoryWithThresholds(ctx context.Context, id string) (*db.CompetencyCategory, error)
	ListCategories(ctx context.Context) ([]*db.CompetencyCategory, error)
	CreateCategory(ctx context.Context, input CreateCompetencyInput) (*db.CompetencyCategory, error)
	UpdateCategory(ctx context.Context, id string, updates CompetencyCategoryUpdates) (*db.CompetencyCategory, error)
	DeleteCategory(ctx context.Context, id string) error
	CountCompetenciesByCategory(ctx context.Context, categoryID string) (int, error)
	GetCharacterCompetency(ctx context.Context, charID int, categoryID string) (*db.CharacterCompetency, error)
	UpsertCharacterCompetency(ctx context.Context, charID int, categoryID string, xp, level int) (*db.CharacterCompetency, error)
}

// RaceRepo defines data access for races.
type RaceRepo interface {
	Get(ctx context.Context, id int) (*db.Race, error)
	GetWithTags(ctx context.Context, id int) (*db.Race, error)
	GetByName(ctx context.Context, name, worldID string) (*db.Race, error)
	List(ctx context.Context, worldID string) ([]*db.Race, error)
	ListWithTags(ctx context.Context, worldID string) ([]*db.Race, error)
	ListPlayable(ctx context.Context, worldID string) ([]*Race, error)
	Create(ctx context.Context, input CreateRaceInput) (*db.Race, error)
	Update(ctx context.Context, id int, updates RaceUpdates) (*db.Race, error)
	Delete(ctx context.Context, id int) error
	CountCharactersByRaceName(ctx context.Context, raceName, worldID string) (int, error)
}

// GenderRepo defines data access for genders.
type GenderRepo interface {
	Get(ctx context.Context, id int) (*db.Gender, error)
	GetByName(ctx context.Context, name string) (*db.Gender, error)
	GetByWorld(ctx context.Context, name, worldID string) (*db.Gender, error)
	List(ctx context.Context, worldID string) ([]*db.Gender, error)
	Create(ctx context.Context, input CreateGenderInput) (*db.Gender, error)
	Update(ctx context.Context, id int, updates GenderUpdates) (*db.Gender, error)
	Delete(ctx context.Context, id int) error
}

// WorldRepo defines data access for worlds.
type WorldRepo interface {
	Get(ctx context.Context, id int) (*db.World, error)
	GetByName(ctx context.Context, name string) (*db.World, error)
	List(ctx context.Context) ([]*db.World, error)
	GetActive(ctx context.Context) ([]*db.World, error)
	Create(ctx context.Context, input CreateWorldInput) (*db.World, error)
	Update(ctx context.Context, id int, updates WorldUpdates) (*db.World, error)
	Delete(ctx context.Context, id int) error
}

// EquipmentTemplateRepo defines data access for equipment templates.
type EquipmentTemplateRepo interface {
	Get(ctx context.Context, id int) (*db.EquipmentTemplate, error)
	GetBySlug(ctx context.Context, slug string, worldID string) (*db.EquipmentTemplate, error)
	List(ctx context.Context, worldID string) ([]*db.EquipmentTemplate, error)
	Create(ctx context.Context, input CreateEquipmentTemplateInput) (*db.EquipmentTemplate, error)
	Update(ctx context.Context, id int, updates EquipmentTemplateUpdates) (*db.EquipmentTemplate, error)
	Delete(ctx context.Context, id int) error
}

// ItemInstanceRepo defines data access for NPC instances and item instances.
type ItemInstanceRepo interface {
	ListNPCsByRoom(ctx context.Context, roomID int) ([]*db.Character, error)
}

// AppLogRepo defines data access for application logs.
type AppLogRepo interface {
	Create(ctx context.Context, level, message, service string, charID, roomID *int, templateID *string, metadata map[string]interface{}) (*db.AppLog, error)
	List(ctx context.Context, filter LogFilter) ([]*db.AppLog, int, error)
	ListServices(ctx context.Context) ([]string, error)
}

// TransactionRunner wraps database transaction support.
type TransactionRunner interface {
	WithTx(ctx context.Context, fn func(tx *db.Tx) error) error
}

// --- Input/Update types ---

type CreateCharacterInput struct {
	Name             string
	UserID           int
	RoomID           int
	StartingRoomID   int
	RespawnRoomID    int
	WorldID          string
	IsAdmin          bool
	IsNPC            bool
	NPCTemplateID    string
	Strength         int
	Dexterity        int
	Constitution     int
	Intelligence     int
	Wisdom           int
	Charisma         int
	Race             string
	Gender           string
	Description      string
	Class            string
	Specialty        string
	HP               int
	MaxHP            int
	Stamina          int
	MaxStamina       int
	Mana             int
	MaxMana          int
	Level            int
	XP               int
	SkillBlades      int
	SkillStaves      int
	SkillKnives      int
	SkillMartial     int
	SkillBrawling    int
	SkillTech        int
	SkillLightArmor  int
	SkillClothArmor  int
	SkillHeavyArmor  int
}

type CharacterUpdates struct {
	Name            *string
	CurrentRoomID   *int
	StartingRoomID  *int
	RespawnRoomID   *int
	Hitpoints       *int
	MaxHitpoints    *int
	Stamina         *int
	MaxStamina      *int
	Mana            *int
	MaxMana         *int
	Level           *int
	Xp              *int
	IsNPC           *bool
	IsImmortal      *bool
	IsAdmin         *bool
	IsTest          *bool
	Race            *string
	Gender          *string
	Class           *string
	Specialty       *string
	Description     *string
	LastSeenAt      *time.Time
	Strength        *int
	Dexterity       *int
	Constitution    *int
	Intelligence    *int
	Wisdom          *int
	SkillBlades     *int
	SkillStaves     *int
	SkillKnives     *int
	SkillMartial    *int
	SkillBrawling   *int
	SkillTech       *int
	SkillLightArmor *int
	SkillClothArmor *int
	SkillHeavyArmor *int
	DiedAt          *time.Time
}

type CreateRoomInput struct {
	Name           string
	Description    string
	IsStartingRoom bool
	IsRootRoom     bool
	Exits          map[string]int
	Atmosphere     string
	PosX           int
	PosY           int
	PosZ           int
	WorldID        string
}

type RoomUpdates struct {
	Name           *string
	Description    *string
	IsStartingRoom *bool
	IsRootRoom     *bool
	Exits          *map[string]int
	Atmosphere     *string
	PosX           *int
	PosY           *int
	PosZ           *int
	WorldID        *string
}

type CreateQuestInput struct {
	Name                 string
	Description          string
	PrerequisiteQuestIDs []string
	Objectives           []schema.QuestObjective
	Rewards              schema.QuestRewards
	RepeatMode           quest.RepeatMode
	CooldownHours        int
	IsActive             bool
	WorldID              string
}

type QuestUpdates struct {
	Name                 *string
	Description          *string
	PrerequisiteQuestIDs *[]string
	Objectives           *[]schema.QuestObjective
	Rewards             *schema.QuestRewards
	RepeatMode          *quest.RepeatMode
	CooldownHours       *int
	IsActive            *bool
	WorldID             *string
}

type CreateQuestProgressInput struct {
	CharacterID      int
	QuestID           int
	Status           questprogress.Status
	StartedAt        time.Time
	CurrentStep      int
	ObjectiveCounts  map[string]int
}

type QuestProgressUpdates struct {
	Status          *questprogress.Status
	CurrentStep     *int
	ObjectiveCounts *map[string]int
	CompletedAt     *time.Time
}

type CreateEquipmentInput struct {
	Name                  string
	Description           string
	Slot                  string
	OwnerID               *int
	RoomID                *int
	Level                 int
	ItemType              string
	ArmorRating           int
	ArmorType             string
	DamageDiceCount       int
	DamageDiceSides       int
	DamageBonus           int
	DamageType            string
	WeaponType             string
	IsTwoHanded           bool
	Stats                 map[string]int
	Rarity                string
	SkillRequirement      string
	SkillRequirementLevel int
	Weight                 int
	IsEquipped             bool
	IsImmovable           bool
	Color                  string
	IsVisible             bool
	EffectType            string
	EffectValue           int
	EffectDuration        int
	EquipmentTemplateID   *int
	Healing               int
	IsContainer           bool
	ContainerCapacity     int
	IsLocked              bool
	KeyItemID             *string
	ContainedItems        string
	RevealCondition       string
	ExpiresAt             *time.Time
}

type EquipmentUpdates struct {
	Name                      *string
	Description               *string
	Slot                      *string
	Level                     *int
	Weight                    *int
	OwnerID                   *int
	RoomID                    *int
	IsEquipped                *bool
	IsImmovable               *bool
	Color                     *string
	IsVisible                 *bool
	ItemType                  *string
	Hitpoints                 *int
	MaxHitpoints              *int
	Healing                   *int
	Effect                    *string
	ClearRoom                 bool
	ArmorRating               *int
	ArmorType                 *string
	Stats                     map[string]int
	Rarity                    *string
	SkillRequirement          *string
	SkillRequirementLevel     *int
	DamageDiceCount           *int
	DamageDiceSides           *int
	DamageBonus               *int
	DamageType                *string
	WeaponType                *string
	IsTwoHanded               *bool
	ExpiresAt                 *time.Time
}

type CreateNPCTemplateInput struct {
	ID               string
	Slug             string
	Name             string
	Description      string
	RaceID           int
	Disposition      string
	Level            int
	XPValue          int
	Skills           map[string]int
	TradesWith       []string
	Greeting         string
	RespawnRooms     []string
	RespawnCooldown  *int
	WorldID          string
}

type NPCTemplateUpdates struct {
	Name             *string
	Slug             *string
	Description      *string
	RaceID           *int
	Disposition      *string
	Level            *int
	XPValue          *int
	Skills           *map[string]int
	TradesWith       *[]string
	Greeting         *string
	RespawnRooms     *[]string
	RespawnCooldown  *int
	WorldID          *string
}

type CreateAbilityInput struct {
	Name             string
	Description      string
	AbilityType     string
	AbilityClass    string
	Cost             int
	Cooldown         int
	ManaCost         int
	StaminaCost      int
	HPCost           int
	Requirements     string
	RequiredTag      string
	ProcChance       float64
	ProcEvent        string
	CooldownSeconds  int
	Slug             string
	FactionID        *int
	WorldID          string
}

type AbilityUpdates struct {
	Name            *string
	Description     *string
	AbilityType    *string
	AbilityClass   *string
	Cost            *int
	Cooldown        *int
	ManaCost        *int
	StaminaCost     *int
	HPCost          *int
	Requirements    *string
	RequiredTag     *string
	ProcChance      *float64
	ProcEvent       *string
	CooldownSeconds *int
	Slug            *string
	FactionID       *int
	WorldID         *string
}

type CreateEffectInput struct {
	Name         string
	Description  string
	EffectType  string
	Parameters  map[string]interface{}
	StackMode   string
	StackLimit  int
	IsPermanent bool
	DurationSecs int
	Messages    map[string]string
}

type EffectUpdates struct {
	Name         *string
	Description  *string
	EffectType  *string
	Parameters  *map[string]interface{}
	StackMode   *string
	StackLimit  *int
	IsPermanent *bool
	DurationSecs *int
	Messages    *map[string]string
}

type CreateActiveEffectInput struct {
	CharacterID int
	EffectID    int
	AppliedByID int
	StackCount  int
	ExpiresAt   *time.Time
}

type ActiveEffectUpdates struct {
	StackCount *int
	IsActive   *bool
	ExpiresAt  *time.Time
	StartedAt  *time.Time
}

type CreateEffectHookInput struct {
	Name          string
	Event         string
	Target        string
	Condition     string
	Enabled       bool
	EffectID      int
	NPCTemplateID *string
}

type EffectHookUpdates struct {
	Name      *string
	Event     *string
	Target    *string
	Condition *string
	Enabled   *bool
	EffectID  *int
}

type CreateDialogNodeInput struct {
	ID             string
	NPCTemplateID  string
	NPCText        string
	Responses      []schema.DialogResponse
	IsEntry        bool
	EntryCondition string
	OnEnterEffects []int
	WorldID        string
}

type DialogNodeUpdates struct {
	NPCText        *string
	Responses      *[]schema.DialogResponse
	IsEntry        *bool
	EntryCondition *string
	OnEnterEffects *[]int
	NPCTemplateID  *string
	WorldID        *string
}

type CreateUserInput struct {
	Email         string
	Password      string
	IsAdmin       bool
	AllowedWorlds string
}

type UserUpdates struct {
	Email         *string
	Password      *string
	IsAdmin       *bool
	AllowedWorlds *string
}

type CreateFactionInput struct {
	Name        string
	DisplayName string
	Description string
	MemberTags  []string
}

type FactionUpdates struct {
	Name        *string
	DisplayName *string
	Description *string
	MemberTags  []string
}

type CreateTagInput struct {
	Name  string
	Color string
}

type TagUpdates struct {
	Name  *string
	Color *string
}

type CreateAchievementInput struct {
	Name        string
	Description string
	Icon        string
	XPReward    int
	Criteria    string
}

type AchievementUpdates struct {
	Name        *string
	Description *string
	Icon        *string
	XPReward    *int
	Criteria    *string
}

type CreateCompetencyInput struct {
	ID           string
	Name         string
	XPMultiplier float64
}

type CompetencyCategoryUpdates struct {
	Name         *string
	XPMultiplier *float64
}

type CreateRaceInput struct {
	Name           string
	DisplayName    string
	Description    string
	StatModifiers  *string
	RequirementTags []string
	Color          string
	EquipmentSlots []string
	TagIDs         []int
}

type RaceUpdates struct {
	Name           *string
	DisplayName    *string
	Description    *string
	StatModifiers  *string
	RequirementTags []string
	Color          *string
	EquipmentSlots []string
	ClearTags      bool
	AddTagIDs      []int
}

type CreateEquipmentTemplateInput struct {
	Slug                  string
	Name                  string
	Description           string
	Slot                  string
	Level                 int
	Weight                int
	ItemType              string
	Stats                 map[string]int
	Color                 string
	IsVisible             bool
	IsImmovable           bool
	EffectType            string
	EffectValue           int
	EffectDuration        int
	IsContainer           bool
	ContainerCapacity     int
	IsLocked              bool
	KeyItemID             string
	RevealCondition       string
	ArmorRating           int
	ArmorType             string
	Rarity                string
	SkillRequirement      string
	SkillRequirementLevel int
	DamageDiceCount       int
	DamageDiceSides       int
	DamageBonus           int
	DamageType            string
	WeaponType            string
	IsTwoHanded           bool
	WorldID               string
}

type EquipmentTemplateUpdates struct {
	Name                    *string
	Description             *string
	Slot                    *string
	Level                   *int
	Weight                  *int
	ItemType                *string
	Stats                   map[string]int
	Color                   *string
	IsVisible               *bool
	IsImmovable             *bool
	EffectType              *string
	EffectValue             *int
	EffectDuration          *int
	IsContainer             *bool
	ContainerCapacity       *int
	IsLocked                *bool
	KeyItemID               *string
	RevealCondition         *string
	ArmorRating             *int
	ArmorType               *string
	Rarity                  *string
	SkillRequirement        *string
	SkillRequirementLevel   *int
	DamageDiceCount         *int
	DamageDiceSides         *int
	DamageBonus             *int
	DamageType              *string
	WeaponType              *string
	IsTwoHanded             *bool
	WorldID                 *string
}

type LogFilter struct {
	Level       string
	Service     string
	CharacterID *int
	RoomID      *int
	TemplateID  *string
	Limit       int
	Offset      int
}

type CreateWorldInput struct {
	Name        string
	Title       string
	Description string
	Active      bool
}

type WorldUpdates struct {
	Name        *string
	Title       *string
	Description *string
	Active      *bool
}

// CreateTriggerInput defines the input for creating a trigger.
type CreateTriggerInput struct {
	Name        string
	WorldID     string
	TriggerType string
	TargetType  string
	TargetID    int
	RoomID      *int
	EquipmentID *int
	Condition   string
	Enabled     bool
}

// TriggerUpdates defines the updates for a trigger.
type TriggerUpdates struct {
	Name        *string
	WorldID     *string
	TriggerType *string
	TargetType  *string
	TargetID    *int
	RoomID      *int
	EquipmentID *int
	Condition   *string
	Enabled     *bool
}

type CreateGenderInput struct {
	Name             string
	DisplayName      string
	SubjectPronoun   string
	ObjectPronoun    string
	PossessivePronoun string
	WorldID          string
}

type GenderUpdates struct {
	Name             *string
	DisplayName      *string
	SubjectPronoun   *string
	ObjectPronoun    *string
	PossessivePronoun *string
	WorldID          *string
}