package service

import (
	"log/slog"

	"herbst-server/db"
	"herbst-server/repository"
)

// Container holds all service instances.
type Container struct {
	Character          CharacterService
	XP                 XPAwardService
	SkillXP            SkillXPService
	AbilityEligibility AbilityEligibilityService
	Room               RoomService
	Quest              QuestService
	QuestProgress      QuestProgressService
	Combat             CombatService
	Equipment          EquipmentService
	NPC                NPCService
	Ability            AbilityService
	Effect             EffectService
	Dialog             DialogService
	Chat               ChatService
	Zone               *ZoneService
	Client             *db.Client
}

// NewContainer creates all services with their dependencies wired.
// Initially, existing services are wired with *db.Client directly.
// As services are migrated to use repositories, the wiring will switch
// to use repo interfaces.
func NewContainer(client *db.Client, repos *repository.Container, logger *slog.Logger) *Container {
	charSvc := NewCharacterService(client, repos)
	xpSvc := NewXPAwardService(client, logger)
	abilityEligSvc := NewAbilityEligibilityService(client)

	return &Container{
		Character:          charSvc,
		XP:                 xpSvc,
		AbilityEligibility: abilityEligSvc,
		Quest:             NewQuestService(repos.Quest, repos.QuestProgress),
		QuestProgress:     NewQuestProgressService(repos.QuestProgress, repos.Quest, repos.Character),
		Room:               NewRoomService(repos.Room, repos.Character, repos.Equipment, repos.NPCTemplate, repos.Tx, repos.Zone),
		Combat:             NewCombatService(repos.Character, repos.DamageLog, repos.NPCTemplate, logger),
		Ability:            NewAbilityService(repos.CharacterAbility, repos.Ability, repos.Character),
		Chat:               NewChatService(repos.Character, repos.ChannelSubscription, repos.OfflineTell, repos.Ignore),
		NPC:                NewNPCService(repos.NPCTemplate),
		Zone:               NewZoneService(repos.Zone, repos.NPCTemplate),
		Client:             client,
	}
}