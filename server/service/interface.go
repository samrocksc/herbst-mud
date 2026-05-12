package service

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/schema"
)

// CharacterService handles character creation, deletion, tags, and race sync.
type CharacterService interface {
	CreateCharacter(ctx context.Context, input CreateCharacterInput) (*db.Character, error)
	DeleteCharacter(ctx context.Context, charID int) error
	GrantTag(ctx context.Context, charID int, tag, source string) error
	SyncRaceTags(ctx context.Context, charID int, raceName string) error
	QueryCharacterByName(ctx context.Context, name string) (*db.Character, error)
}

// XPAwardService handles XP awards, death penalties, and competency tracking.
type XPAwardService interface {
	AwardXP(ctx context.Context, characterID, xpGained int) (newXP, newLevel int, leveledUp bool, err error)
	ApplyDeathPenalty(ctx context.Context, characterID, penaltyPercent int) (xpLost, newXP int, err error)
	AwardCompetencyXP(ctx context.Context, characterID int, categoryID string, rawXP int) error
	SeedCompetencyCategories(ctx context.Context) error
	GetCharacterXP(ctx context.Context, characterID int) (xp, level int, err error)
	QueryCharacter(ctx context.Context, id int) (*db.Character, error)
}

// AbilityEligibilityService checks whether characters meet ability prerequisites.
type AbilityEligibilityService interface {
	CheckEligibility(ctx context.Context, charID int, sk *db.Ability, activeFactionIDs map[int]bool, tagSet map[string]bool) AbilityEligibility
	CheckEligibilityForCharacter(ctx context.Context, charID int) (map[int]AbilityEligibility, error)
	GetEligiblePassiveAbilitiesForEvent(ctx context.Context, charID int, procEvent string) ([]*db.Ability, error)
	AbilitiesForCharacterWithEligibility(ctx context.Context, charID int) ([]AbilityWithEligibility, error)
}

// RoomService handles room CRUD, exit management, and bidirectional exits.
type RoomService interface {
	CreateRoom(ctx context.Context, input CreateRoomInput) (*db.Room, error)
	GetRoom(ctx context.Context, id int) (*db.Room, error)
	ListRooms(ctx context.Context) ([]*db.Room, error)
	UpdateRoom(ctx context.Context, id int, input UpdateRoomInput) (*db.Room, error)
	DeleteRoom(ctx context.Context, id int) error
	CleanupOrphanExits(ctx context.Context) (int, error)
	CreateBidirectionalExit(ctx context.Context, sourceID int, dir string, targetID int) (*BidirectionalExitResult, error)
	DeleteBidirectionalExit(ctx context.Context, sourceID int, dir string) error
}

// QuestService handles quest definition CRUD.
type QuestService interface {
	CreateQuest(ctx context.Context, input CreateQuestInput) (*db.Quest, error)
	GetQuest(ctx context.Context, id int) (*db.Quest, error)
	ListQuests(ctx context.Context) ([]*db.Quest, error)
	UpdateQuest(ctx context.Context, id int, input UpdateQuestInput) (*db.Quest, error)
	DeleteQuest(ctx context.Context, id int) error
}

// QuestProgressService handles quest acceptance, advancement, and abandonment.
type QuestProgressService interface {
	Accept(ctx context.Context, charID int, questID int) (*QuestProgressView, error)
	Advance(ctx context.Context, charID int, questID int, objectiveKey string, count int) (*QuestProgressView, error)
	Abandon(ctx context.Context, charID int, questID int) error
	CheckAll(ctx context.Context, charID int, objectiveType, targetID string) ([]QuestProgressView, error)
	ListByCharacter(ctx context.Context, charID int) ([]QuestProgressView, error)
	ValidateAcceptance(ctx context.Context, charID int, questID int) error
}

// CombatService handles damage, healing, stamina, mana, and NPC combat.
type CombatService interface {
	ApplyDamage(ctx context.Context, targetID int, damage int) (*CombatResult, error)
	LogDamage(ctx context.Context, attackerID, targetID, damage int)
	GetCombatStatus(ctx context.Context, charID int) (*CombatStatusResult, error)
	HealCharacter(ctx context.Context, charID int, amount int) (*HealResult, error)
	AdjustStamina(ctx context.Context, charID int, amount int) (*StatResult, error)
	AdjustMana(ctx context.Context, charID int, amount int) (*StatResult, error)
	HealNPCsInRoom(ctx context.Context, roomID int, amount int) (int, error)
	PassiveHealNPCsInRoom(ctx context.Context, roomID int) (*PassiveHealResult, error)
}

// EquipmentService handles equip/unequip and slot management.
type EquipmentService interface {
	EquipItem(ctx context.Context, charID int, itemID int, slot string) (*db.Equipment, error)
	UnequipItem(ctx context.Context, charID int, slot string) (*db.Equipment, error)
	GetEquipment(ctx context.Context, charID int) (map[string]*db.Equipment, error)
	GetInventory(ctx context.Context, charID int) ([]*db.Equipment, error)
	UseConsumable(ctx context.Context, charID int, itemID int) error
}

// NPCService handles NPC template and instance CRUD.
type NPCService interface {
	GetTemplate(ctx context.Context, id string) (*db.NPCTemplate, error)
	ListTemplates(ctx context.Context) ([]*db.NPCTemplate, error)
	CreateTemplate(ctx context.Context, input CreateNPCTemplateInput) (*db.NPCTemplate, error)
	UpdateTemplate(ctx context.Context, id string, input UpdateNPCTemplateInput) (*db.NPCTemplate, error)
	DeleteTemplate(ctx context.Context, id string) error
}

// AbilityService handles ability CRUD and slot management.
type AbilityService interface {
	GetAbility(ctx context.Context, id int) (*db.Ability, error)
	ListAbilities(ctx context.Context) ([]*db.Ability, error)
	ListClasslessAbilities(ctx context.Context) ([]*db.Ability, error)
	ListPassiveAbilities(ctx context.Context) ([]*db.Ability, error)
	CreateAbility(ctx context.Context, input CreateAbilityInput) (*db.Ability, error)
	UpdateAbility(ctx context.Context, id int, input UpdateAbilityInput) (*db.Ability, error)
	DeleteAbility(ctx context.Context, id int) error
	EquipAbility(ctx context.Context, charID int, abilityID int, slot int) error
	UnequipAbility(ctx context.Context, charID int, slot int) error
	SwapAbilities(ctx context.Context, charID int, slot1, slot2 int) (*SwapResult, error)
	GetAbilitiesWithDetails(ctx context.Context, charID int) ([]*db.CharacterAbility, error)
	UnlockPassiveAbility(ctx context.Context, charID int, abilityID int) (*db.CharacterAbility, error)
	RemovePassiveAbility(ctx context.Context, charID int, abilityID int) error
	EquipClasslessSkill(ctx context.Context, charID int, skillID int, slot int) error
	SwapClasslessSkills(ctx context.Context, charID int, slot1, slot2 int) error
}

// EffectService handles effect CRUD and active effect management.
type EffectService interface {
	CreateEffect(ctx context.Context, input CreateEffectInput) (*db.Effect, error)
	GetEffect(ctx context.Context, id int) (*db.Effect, error)
	ListEffects(ctx context.Context) ([]*db.Effect, error)
	UpdateEffect(ctx context.Context, id int, input UpdateEffectInput) (*db.Effect, error)
	DeleteEffect(ctx context.Context, id int) error
	ApplyEffect(ctx context.Context, charID int, effectID int, appliedBy int) (*db.ActiveEffect, error)
	RemoveActiveEffect(ctx context.Context, id int) error
}

// DialogService handles dialog node CRUD.
type DialogService interface {
	GetNode(ctx context.Context, id string) (*db.DialogNode, error)
	ListNodes(ctx context.Context) ([]*db.DialogNode, error)
	ListNodesByTemplate(ctx context.Context, templateID string) ([]*db.DialogNode, error)
	CreateNode(ctx context.Context, input CreateDialogNodeInput) (*db.DialogNode, error)
	UpdateNode(ctx context.Context, id string, input UpdateDialogNodeInput) (*db.DialogNode, error)
	DeleteNode(ctx context.Context, id string) error
}

// ChatService handles messaging: say/yell/shout/tell/whisper/emote, channel chat, and ignore system.
type ChatService interface {
	SendSay(ctx context.Context, charID, roomID int, message string) (*MessageResult, error)
	SendYell(ctx context.Context, charID, roomID int, message string) (*MessageResult, error)
	SendShout(ctx context.Context, charID int, message string) (*MessageResult, error)
	SendTell(ctx context.Context, fromID, toID int, message string) (*MessageResult, error)
	SendWhisper(ctx context.Context, fromID, toID int, message string) (*MessageResult, error)
	SendEmote(ctx context.Context, charID int, action string) (*MessageResult, error)
	SendChannel(ctx context.Context, channel, message string, charID int) (*MessageResult, error)
	GetChannels(charID int) ([]ChannelState, error)
	SetChannelEnabled(ctx context.Context, charID int, channel string, enabled bool) error
	SetChannelColor(ctx context.Context, charID int, channel string, color string) error
	IgnorePlayer(ctx context.Context, charID, ignoredID int) error
	UnignorePlayer(ctx context.Context, charID, ignoredID int) error
	GetIgnoredPlayers(ctx context.Context, charID int) ([]int, error)
	QueueOfflineTell(ctx context.Context, fromID int, recipientName string, message string) error
	DeliverQueuedTells(ctx context.Context, charID int) ([]QueuedTell, error)
}

// --- View/Result types (service-layer DTOs) ---

// CombatResult is returned by combat operations.
type CombatResult struct {
	ID       int
	HP       int
	MaxHP    int
	Defeated bool
	Immortal bool
	Message  string
}

// CombatStatusResult is returned by GetCombatStatus.
type CombatStatusResult struct {
	ID    int
	HP    int
	MaxHP int
	IsNPC bool
}

// HealResult is returned by heal operations.
type HealResult struct {
	ID    int
	HP    int
	MaxHP int
}

// StatResult is returned by stamina/mana adjustments.
type StatResult struct {
	ID      int
	Current int
	Max     int
}

// PassiveHealResult is returned by NPC passive heal operations.
type PassiveHealResult struct {
	Healed int
}

// SwapResult is returned by ability swap operations.
type SwapResult struct {
	Slot1 *db.CharacterAbility
	Slot2 *db.CharacterAbility
}

// AbilityEligibility describes whether a character can use an ability.
type AbilityEligibility struct {
	Eligible bool   `json:"eligible"`
	Reason   string `json:"reason,omitempty"`
}

// AbilityWithEligibility pairs an ability with its eligibility.
type AbilityWithEligibility struct {
	Ability    *db.Ability
	Eligibility AbilityEligibility
}

// QuestProgressView is the DTO returned by quest progress operations.
type QuestProgressView struct {
	ID               int                       `json:"id"`
	CharacterID      int                       `json:"character_id"`
	QuestID          int                       `json:"quest_id"`
	QuestName        string                    `json:"quest_name"`
	QuestDescription string                    `json:"quest_description,omitempty"`
	Status           string                    `json:"status"`
	CurrentStep      int                       `json:"current_step"`
	ObjectiveCounts  map[string]int            `json:"objective_counts"`
	StartedAt        string                    `json:"started_at"`
	CompletedAt      *string                   `json:"completed_at,omitempty"`
	RewardsApplied   map[string]interface{}    `json:"rewards_applied,omitempty"`
}

// --- Input types for services ---

// These are defined in existing files (character.go, xp.go, ability_eligibility.go)
// or will be defined in new service implementation files.
// The service-layer input types are separate from repository input types.

// CreateRoomInput for RoomService.
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
}

// UpdateRoomInput for RoomService.
type UpdateRoomInput struct {
	Name           *string
	Description    *string
	IsStartingRoom *bool
	IsRootRoom     *bool
	Exits          *map[string]int
	Atmosphere     *string
	PosX           *int
	PosY           *int
	PosZ           *int
	Version        *int
}

// CreateQuestInput for QuestService.
type CreateQuestInput struct {
	Name                 string
	Description          string
	PrerequisiteQuestIDs []string
	Objectives           []schema.QuestObjective
	Rewards              schema.QuestRewards
	RepeatMode           string
	CooldownHours        int
	IsActive             bool
}

// UpdateQuestInput for QuestService.
type UpdateQuestInput struct {
	Name                 *string
	Description          *string
	PrerequisiteQuestIDs *[]string
	Objectives           *[]schema.QuestObjective
	Rewards             *schema.QuestRewards
	RepeatMode          *string
	CooldownHours       *int
	IsActive            *bool
}

// CreateNPCTemplateInput for NPCService.
type CreateNPCTemplateInput struct {
	ID              string
	Name            string
	Description     string
	Race            string
	Disposition     string
	Level           int
	XPValue         int
	Skills          map[string]int
	TradesWith      []string
	Greeting        string
	RespawnRooms    []string
	RespawnCooldown *int
}

// UpdateNPCTemplateInput for NPCService.
type UpdateNPCTemplateInput struct {
	Name            *string
	Description     *string
	Disposition     *string
	Level           *int
	XPValue         *int
	Skills          *map[string]int
	TradesWith      *[]string
	Greeting        *string
	RespawnRooms    *[]string
	RespawnCooldown *int
}

// CreateAbilityInput for AbilityService.
type CreateAbilityInput struct {
	Name             string
	Description      string
	AbilityType      string
	AbilityClass     string
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
}

// UpdateAbilityInput for AbilityService.
type UpdateAbilityInput struct {
	Name             *string
	Description      *string
	AbilityType      *string
	AbilityClass     *string
	Cost             *int
	Cooldown         *int
	ManaCost         *int
	StaminaCost      *int
	HPCost           *int
	Requirements     *string
	RequiredTag      *string
	ProcChance       *float64
	ProcEvent        *string
	CooldownSeconds  *int
	Slug             *string
	FactionID        *int
}

// CreateEffectInput for EffectService.
type CreateEffectInput struct {
	Name          string
	Description   string
	EffectType    string
	Parameters   map[string]interface{}
	StackMode    string
	StackLimit   int
	IsPermanent  bool
	DurationSecs int
	Messages     map[string]string
}

// UpdateEffectInput for EffectService.
type UpdateEffectInput struct {
	Name          *string
	Description   *string
	EffectType    *string
	Parameters   *map[string]interface{}
	StackMode    *string
	StackLimit   *int
	IsPermanent  *bool
	DurationSecs *int
	Messages     *map[string]string
}

// CreateDialogNodeInput for DialogService.
type CreateDialogNodeInput struct {
	ID             string
	NPCTemplateID  string
	NPCText        string
	Responses      []schema.DialogResponse
	IsEntry        bool
	EntryCondition string
	OnEnterEffects []int
}

// UpdateDialogNodeInput for DialogService.
type UpdateDialogNodeInput struct {
	NPCText        *string
	Responses      *[]schema.DialogResponse
	IsEntry        *bool
	EntryCondition *string
	OnEnterEffects *[]int
}

// MessageResult is returned by chat operations.
type MessageResult struct {
	FromCharacterID   int      `json:"from_character_id"`
	FromCharacterName string   `json:"from_character_name"`
	ToCharacterIDs    []int    `json:"to_character_ids"`
	Channel           string   `json:"channel"`
	Message           string   `json:"message"`
	Type              string   `json:"type"` // say, yell, shout, tell, whisper, emote, channel, system
	RoomID            *int     `json:"room_id,omitempty"`
	DisplayMessage    string   `json:"display_message,omitempty"` // Message to display to the sender
}

// ChannelState describes a chat channel subscription.
type ChannelState struct {
	Name    string `json:"name"`
	Color   string `json:"color"`
	Enabled bool   `json:"enabled"`
}

// QueuedTell is an offline message waiting for delivery.
type QueuedTell struct {
	ID            int    `json:"id"`
	FromID        int    `json:"from_id"`
	FromName      string `json:"from_name"`
	RecipientName string `json:"recipient_name"`
	Message       string `json:"message"`
	QueuedAt      string `json:"queued_at"`
}