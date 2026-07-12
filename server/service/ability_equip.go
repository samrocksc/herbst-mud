package service

import (
	"context"
	"encoding/json"
	"fmt"

	"herbst-server/db"
)

func (s *abilityService) EquipAbility(ctx context.Context, charID int, abilityID int, slot int) error {
	if slot < 1 || slot > 5 {
		return ErrSlotOutOfRange
	}
	_, err := s.abilityRepo.Get(ctx, abilityID)
	if err != nil {
		return ErrAbilityNotFound
	}
	existing, _ := s.charAbilityRepo.ListByCharacter(ctx, charID)
	slotsUsed := 0
	for _, ca := range existing {
		if ca.Slot != slot {
			slotsUsed++
		}
	}
	if slotsUsed >= 5 {
		return ErrMaxAbilities
	}
	if err := s.validateSkillRequirements(ctx, charID, abilityID); err != nil {
		return err
	}
	_ = s.charAbilityRepo.DeleteByCharacterAndSlot(ctx, charID, slot)
	_, err = s.charAbilityRepo.Create(ctx, charID, abilityID, slot)
	return err
}

func (s *abilityService) UnequipAbility(ctx context.Context, charID int, slot int) error {
	if slot < 1 || slot > 5 {
		return ErrSlotOutOfRange
	}
	abilities, err := s.charAbilityRepo.ListByCharacterAndSlot(ctx, charID, slot)
	if err != nil || len(abilities) == 0 {
		return ErrNotEquipped
	}
	for _, ca := range abilities {
		_ = s.charAbilityRepo.Delete(ctx, ca.ID)
	}
	return nil
}

func (s *abilityService) SwapAbilities(ctx context.Context, charID int, slot1, slot2 int) (*SwapResult, error) {
	if slot1 < 1 || slot1 > 5 || slot2 < 1 || slot2 > 5 {
		return nil, ErrSlotOutOfRange
	}
	ability1ID, _ := s.getAbilityIDInSlot(ctx, charID, slot1)
	ability2ID, _ := s.getAbilityIDInSlot(ctx, charID, slot2)
	_ = s.charAbilityRepo.DeleteByCharacterAndSlot(ctx, charID, slot1)
	_ = s.charAbilityRepo.DeleteByCharacterAndSlot(ctx, charID, slot2)
	if ability1ID > 0 {
		s.charAbilityRepo.Create(ctx, charID, ability1ID, slot2)
	}
	if ability2ID > 0 {
		s.charAbilityRepo.Create(ctx, charID, ability2ID, slot1)
	}
	return &SwapResult{}, nil
}

func (s *abilityService) getAbilityIDInSlot(ctx context.Context, charID int, slot int) (int, string) {
	cas, err := s.charAbilityRepo.ListByCharacterAndSlot(ctx, charID, slot)
	if err != nil || len(cas) == 0 {
		return 0, ""
	}
	ca := cas[0]
	if ca.Edges.Ability != nil {
		return ca.Edges.Ability.ID, ca.Edges.Ability.Name
	}
	return 0, ""
}

func (s *abilityService) validateSkillRequirements(ctx context.Context, charID int, abilityID int) error {
	abilityObj, err := s.abilityRepo.Get(ctx, abilityID)
	if err != nil {
		return ErrAbilityNotFound
	}
	if abilityObj.Requirements == "" {
		return nil
	}
	requirements := map[string]int{}
	if err := json.Unmarshal([]byte(abilityObj.Requirements), &requirements); err != nil {
		return nil
	}
	_, err = s.charRepo.Get(ctx, charID)
	if err != nil {
		return ErrCharNotFound
	}
	skillLevels, err := s.charRepo.GetSkillLevels(ctx, charID)
	if err != nil {
		return nil
	}
	for skillName, requiredLevel := range requirements {
		if skillLevels[skillName] < requiredLevel {
			return fmt.Errorf("%w: requires %s level %d", ErrSkillRequirements, skillName, requiredLevel)
		}
	}
	return nil
}

// GetAbilitiesWithDetails returns all character abilities with ability+effects eager-loaded.
func (s *abilityService) GetAbilitiesWithDetails(ctx context.Context, charID int) ([]*db.CharacterAbility, error) {
	return s.charAbilityRepo.ListByCharacterWithDetails(ctx, charID)
}

// UnlockPassiveAbility equips a passive ability in the next available slot.
func (s *abilityService) UnlockPassiveAbility(ctx context.Context, charID int, abilityID int) (*db.CharacterAbility, error) {
	abilityObj, err := s.abilityRepo.Get(ctx, abilityID)
	if err != nil {
		return nil, ErrAbilityNotFound
	}
	if abilityObj.AbilityClass != "passive" {
		return nil, ErrNotPassive
	}
	exists, _ := s.charAbilityRepo.ExistsByCharacterAndAbility(ctx, charID, abilityID)
	if exists {
		return nil, ErrAlreadyEquipped
	}
	charAbilities, _ := s.charAbilityRepo.ListByCharacter(ctx, charID)
	usedSlots := map[int]bool{}
	for _, ca := range charAbilities {
		usedSlots[ca.Slot] = true
	}
	slot := 0
	for s := 1; s <= 5; s++ {
		if !usedSlots[s] {
			slot = s
			break
		}
	}
	if slot == 0 {
		return nil, ErrNoAvailableSlots
	}
	return s.charAbilityRepo.Create(ctx, charID, abilityID, slot)
}

// RemovePassiveAbility removes an ability from a character by ability ID.
func (s *abilityService) RemovePassiveAbility(ctx context.Context, charID int, abilityID int) error {
	return s.charAbilityRepo.DeleteByCharacterAndAbility(ctx, charID, abilityID)
}

// EquipClasslessSkill equips an ability in a specific slot (simpler validation than EquipAbility).
func (s *abilityService) EquipClasslessSkill(ctx context.Context, charID int, skillID int, slot int) error {
	if slot < 1 || slot > 5 {
		return ErrSlotOutOfRange
	}
	_, err := s.abilityRepo.Get(ctx, skillID)
	if err != nil {
		return ErrAbilityNotFound
	}
	_ = s.charAbilityRepo.DeleteByCharacterAndSlot(ctx, charID, slot)
	_, err = s.charAbilityRepo.Create(ctx, charID, skillID, slot)
	return err
}

// SwapClasslessSkills swaps abilities between two slots.
func (s *abilityService) SwapClasslessSkills(ctx context.Context, charID int, slot1, slot2 int) error {
	if slot1 < 1 || slot1 > 5 || slot2 < 1 || slot2 > 5 {
		return ErrSlotOutOfRange
	}
	ability1ID, _ := s.getAbilityIDInSlot(ctx, charID, slot1)
	ability2ID, _ := s.getAbilityIDInSlot(ctx, charID, slot2)
	_ = s.charAbilityRepo.DeleteByCharacterAndSlot(ctx, charID, slot1)
	_ = s.charAbilityRepo.DeleteByCharacterAndSlot(ctx, charID, slot2)
	if ability1ID > 0 {
		s.charAbilityRepo.Create(ctx, charID, ability1ID, slot2)
	}
	if ability2ID > 0 {
		s.charAbilityRepo.Create(ctx, charID, ability2ID, slot1)
	}
	return nil
}

// FormatAbilitySlot formats a character ability with its details for API responses.
func FormatAbilitySlot(ca *db.CharacterAbility) map[string]interface{} {
	entry := map[string]interface{}{
		"slot":       ca.Slot,
	}
	if ca.Edges.Ability != nil {
		ab := ca.Edges.Ability
		entry["ability_id"] = ab.ID
		entry["name"] = ab.Name
		entry["description"] = ab.Description
		entry["cooldown"] = ab.Cooldown
		entry["manaCost"] = ab.ManaCost
		entry["staminaCost"] = ab.StaminaCost
		effects := make([]map[string]interface{}, 0)
		for _, e := range ab.Edges.Effects {
			effects = append(effects, map[string]interface{}{
				"effectType":    e.EffectType,
				"damageSubtype": e.DamageSubtype,
				"target":        e.Target,
				"value":         e.Value,
				"duration":      e.Duration,
				"scalingStat":   e.ScalingStat,
				"scalingRatio":  e.ScalingRatio,
				"sortOrder":     e.SortOrder,
			})
		}
		entry["effects"] = effects
	}
	return entry
}

// CalcSkillBonus returns a string bonus for a skill level.
func CalcSkillBonus(skill int) string {
	switch {
	case skill >= 91:
		return "+75%"
	case skill >= 76:
		return "+50%"
	case skill >= 51:
		return "+25%"
	case skill >= 26:
		return "+10%"
	default:
		return "+0%"
	}
}