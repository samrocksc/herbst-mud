package repository

import (
	"context"

	"herbst-server/db"
)

func (r *entEquipmentRepo) Create(ctx context.Context, input CreateEquipmentInput) (*db.Equipment, error) {
	builder := r.client.Equipment.Create().
		SetName(input.Name).
		SetDescription(input.Description).
		SetSlot(input.Slot).
		SetLevel(input.Level).
		SetItemType(input.ItemType).
		SetArmorRating(input.ArmorRating).
		SetArmorType(input.ArmorType).
		SetDamageDiceCount(input.DamageDiceCount).
		SetDamageDiceSides(input.DamageDiceSides).
		SetDamageBonus(input.DamageBonus).
		SetDamageType(input.DamageType).
		SetWeaponType(input.WeaponType).
		SetIsTwoHanded(input.IsTwoHanded).
		SetStats(input.Stats).
		SetRarity(input.Rarity).
		SetSkillRequirement(input.SkillRequirement).
		SetSkillRequirementLevel(input.SkillRequirementLevel).
		SetWeight(input.Weight).
		SetIsEquipped(input.IsEquipped).
		SetIsImmovable(input.IsImmovable).
		SetColor(input.Color).
		SetIsVisible(input.IsVisible).
		SetEffectType(input.EffectType).
		SetEffectValue(input.EffectValue).
		SetEffectDuration(input.EffectDuration).
		SetHealing(input.Healing).
		SetIsContainer(input.IsContainer).
		SetContainerCapacity(input.ContainerCapacity).
		SetIsLocked(input.IsLocked).
		SetContainedItems(input.ContainedItems).
		SetRevealCondition(input.RevealCondition)
	if input.OwnerID != nil {
		builder = builder.SetOwnerId(*input.OwnerID)
	}
	if input.RoomID != nil {
		builder = builder.SetRoomID(*input.RoomID)
	}
	if input.EquipmentTemplateID != nil {
		builder = builder.SetNillableEquipmentTemplateID(input.EquipmentTemplateID)
	}
	if input.KeyItemID != nil {
		builder = builder.SetNillableKeyItemID(input.KeyItemID)
	}
	if input.ExpiresAt != nil {
		builder = builder.SetNillableExpiresAt(input.ExpiresAt)
	}
	return builder.Save(ctx)
}