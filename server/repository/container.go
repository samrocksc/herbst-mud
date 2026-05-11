package repository

import (
	"herbst-server/db"
)

// Container holds all repository instances.
type Container struct {
	Character         CharacterRepo
	Room              RoomRepo
	Quest             QuestRepo
	QuestProgress     QuestProgressRepo
	Equipment         EquipmentRepo
	NPCTemplate       NPCTemplateRepo
	Ability           AbilityRepo
	CharacterAbility  CharacterAbilityRepo
	Effect            EffectRepo
	ActiveEffect      ActiveEffectRepo
	EffectHook        EffectHookRepo
	DialogNode        DialogNodeRepo
	User              UserRepo
	Faction           FactionRepo
	CharacterFaction  CharacterFactionRepo
	CharacterTag      CharacterTagRepo
	Tag               TagRepo
	Achievement       AchievementRepo
	DamageLog         DamageLogRepo
	GameConfig        GameConfigRepo
	Competency        CompetencyRepo
	Race              RaceRepo
	Gender            GenderRepo
	EquipmentTemplate EquipmentTemplateRepo
	ItemInstance      ItemInstanceRepo
	AppLog            AppLogRepo
	Tx                TransactionRunner
}

// NewContainer creates all ent-backed repositories.
func NewContainer(client *db.Client) *Container {
	return &Container{
		Character:         NewEntCharacterRepo(client),
		Room:              NewEntRoomRepo(client),
		Quest:             NewEntQuestRepo(client),
		QuestProgress:     NewEntQuestProgressRepo(client),
		Equipment:         NewEntEquipmentRepo(client),
		NPCTemplate:       NewEntNPCTemplateRepo(client),
		Ability:           NewEntAbilityRepo(client),
		CharacterAbility:  NewEntCharacterAbilityRepo(client),
		Effect:            NewEntEffectRepo(client),
		ActiveEffect:      NewEntActiveEffectRepo(client),
		EffectHook:        NewEntEffectHookRepo(client),
		DialogNode:        NewEntDialogNodeRepo(client),
		User:              NewEntUserRepo(client),
		Faction:           NewEntFactionRepo(client),
		CharacterFaction:  NewEntCharacterFactionRepo(client),
		CharacterTag:      NewEntCharacterTagRepo(client),
		Tag:               NewEntTagRepo(client),
		Achievement:       NewEntAchievementRepo(client),
		DamageLog:         NewEntDamageLogRepo(client),
		GameConfig:        NewEntGameConfigRepo(client),
		Competency:        NewEntCompetencyRepo(client),
		Race:              NewEntRaceRepo(client),
		Gender:            NewEntGenderRepo(client),
		EquipmentTemplate: NewEntEquipmentTemplateRepo(client),
		ItemInstance:      NewEntItemInstanceRepo(client),
		AppLog:            NewEntAppLogRepo(client),
		Tx:                NewEntTransactionRunner(client),
	}
}