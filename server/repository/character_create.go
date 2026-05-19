package repository

import (
	"context"

	"herbst-server/db"
)

func (r *entCharacterRepo) Create(ctx context.Context, input CreateCharacterInput) (*db.Character, error) {
	builder := r.client.Character.Create().
		SetName(input.Name).
		SetCurrentRoomId(input.RoomID).
		SetIsAdmin(input.IsAdmin).
		SetIsNPC(input.IsNPC).
		SetHitpoints(input.HP).
		SetMaxHitpoints(input.MaxHP).
		SetStamina(input.Stamina).
		SetMaxStamina(input.MaxStamina).
		SetMana(input.Mana).
		SetMaxMana(input.MaxMana).
		SetLevel(input.Level).
		SetXp(input.XP).
		SetStrength(input.Strength).
		SetDexterity(input.Dexterity).
		SetConstitution(input.Constitution).
		SetIntelligence(input.Intelligence).
		SetWisdom(input.Wisdom).
		SetRace(input.Race).
		SetGender(input.Gender).
		SetClass(input.Class).
		SetUserID(input.UserID).
		SetNillableNpcTemplateID(&input.NPCTemplateID)
	if input.StartingRoomID > 0 {
		builder.SetStartingRoomId(input.StartingRoomID)
	}
	if input.RespawnRoomID > 0 {
		builder.SetRespawnRoomId(input.RespawnRoomID)
	}
	if input.WorldID != "" {
		builder.SetCurrentWorld(input.WorldID)
	}
	if input.Specialty != "" {
		builder.SetSpecialty(input.Specialty)
	}
	if input.Description != "" {
		builder.SetDescription(input.Description)
	}
	if input.SkillBlades > 0 {
		builder.SetSkillBlades(input.SkillBlades)
	}
	if input.SkillStaves > 0 {
		builder.SetSkillStaves(input.SkillStaves)
	}
	if input.SkillKnives > 0 {
		builder.SetSkillKnives(input.SkillKnives)
	}
	if input.SkillMartial > 0 {
		builder.SetSkillMartial(input.SkillMartial)
	}
	if input.SkillBrawling > 0 {
		builder.SetSkillBrawling(input.SkillBrawling)
	}
	if input.SkillTech > 0 {
		builder.SetSkillTech(input.SkillTech)
	}
	if input.SkillLightArmor > 0 {
		builder.SetSkillLightArmor(input.SkillLightArmor)
	}
	if input.SkillClothArmor > 0 {
		builder.SetSkillClothArmor(input.SkillClothArmor)
	}
	if input.SkillHeavyArmor > 0 {
		builder.SetSkillHeavyArmor(input.SkillHeavyArmor)
	}
	return builder.Save(ctx)
}

func (r *entCharacterRepo) Update(ctx context.Context, id int, updates CharacterUpdates) (*db.Character, error) {
	builder := r.client.Character.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.CurrentRoomID != nil {
		builder = builder.SetCurrentRoomId(*updates.CurrentRoomID)
	}
	if updates.StartingRoomID != nil {
		builder = builder.SetStartingRoomId(*updates.StartingRoomID)
	}
	if updates.RespawnRoomID != nil {
		builder = builder.SetRespawnRoomId(*updates.RespawnRoomID)
	}
	if updates.Hitpoints != nil {
		builder = builder.SetHitpoints(*updates.Hitpoints)
	}
	if updates.MaxHitpoints != nil {
		builder = builder.SetMaxHitpoints(*updates.MaxHitpoints)
	}
	if updates.Stamina != nil {
		builder = builder.SetStamina(*updates.Stamina)
	}
	if updates.MaxStamina != nil {
		builder = builder.SetMaxStamina(*updates.MaxStamina)
	}
	if updates.Mana != nil {
		builder = builder.SetMana(*updates.Mana)
	}
	if updates.MaxMana != nil {
		builder = builder.SetMaxMana(*updates.MaxMana)
	}
	if updates.Level != nil {
		builder = builder.SetLevel(*updates.Level)
	}
	if updates.Xp != nil {
		builder = builder.SetXp(*updates.Xp)
	}
	if updates.IsNPC != nil {
		builder = builder.SetIsNPC(*updates.IsNPC)
	}
	if updates.IsImmortal != nil {
		builder = builder.SetIsImmortal(*updates.IsImmortal)
	}
	if updates.IsAdmin != nil {
		builder = builder.SetIsAdmin(*updates.IsAdmin)
	}
	if updates.IsTest != nil {
		builder = builder.SetIsTest(*updates.IsTest)
	}
	if updates.Race != nil {
		builder = builder.SetRace(*updates.Race)
	}
	if updates.Gender != nil {
		builder = builder.SetGender(*updates.Gender)
	}
	if updates.Class != nil {
		builder = builder.SetClass(*updates.Class)
	}
	if updates.Specialty != nil {
		builder = builder.SetSpecialty(*updates.Specialty)
	}
	if updates.Description != nil {
		builder = builder.SetDescription(*updates.Description)
	}
	if updates.LastSeenAt != nil {
		builder = builder.SetLastSeenAt(*updates.LastSeenAt)
	}
	if updates.Strength != nil {
		builder = builder.SetStrength(*updates.Strength)
	}
	if updates.Dexterity != nil {
		builder = builder.SetDexterity(*updates.Dexterity)
	}
	if updates.Constitution != nil {
		builder = builder.SetConstitution(*updates.Constitution)
	}
	if updates.Intelligence != nil {
		builder = builder.SetIntelligence(*updates.Intelligence)
	}
	if updates.Wisdom != nil {
		builder = builder.SetWisdom(*updates.Wisdom)
	}
	if updates.SkillBlades != nil {
		builder = builder.SetSkillBlades(*updates.SkillBlades)
	}
	if updates.SkillStaves != nil {
		builder = builder.SetSkillStaves(*updates.SkillStaves)
	}
	if updates.SkillKnives != nil {
		builder = builder.SetSkillKnives(*updates.SkillKnives)
	}
	if updates.SkillMartial != nil {
		builder = builder.SetSkillMartial(*updates.SkillMartial)
	}
	if updates.SkillBrawling != nil {
		builder = builder.SetSkillBrawling(*updates.SkillBrawling)
	}
	if updates.SkillTech != nil {
		builder = builder.SetSkillTech(*updates.SkillTech)
	}
	if updates.SkillLightArmor != nil {
		builder = builder.SetSkillLightArmor(*updates.SkillLightArmor)
	}
	if updates.SkillClothArmor != nil {
		builder = builder.SetSkillClothArmor(*updates.SkillClothArmor)
	}
	if updates.SkillHeavyArmor != nil {
		builder = builder.SetSkillHeavyArmor(*updates.SkillHeavyArmor)
	}
	if updates.DiedAt != nil {
		builder = builder.SetDiedAt(*updates.DiedAt)
	}
	return builder.Save(ctx)
}