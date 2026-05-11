package repository

import (
	"context"

	"herbst-server/db"
)

type entEquipmentTemplateRepo struct {
	client *db.Client
}

func NewEntEquipmentTemplateRepo(client *db.Client) EquipmentTemplateRepo {
	return &entEquipmentTemplateRepo{client: client}
}

func (r *entEquipmentTemplateRepo) Get(ctx context.Context, id string) (*db.EquipmentTemplate, error) {
	return r.client.EquipmentTemplate.Get(ctx, id)
}

func (r *entEquipmentTemplateRepo) List(ctx context.Context) ([]*db.EquipmentTemplate, error) {
	return r.client.EquipmentTemplate.Query().All(ctx)
}

func (r *entEquipmentTemplateRepo) Create(ctx context.Context, input CreateEquipmentTemplateInput) (*db.EquipmentTemplate, error) {
	builder := r.client.EquipmentTemplate.Create().
		SetID(input.ID).
		SetName(input.Name).
		SetDescription(input.Description).
		SetSlot(input.Slot).
		SetLevel(input.Level).
		SetWeight(input.Weight).
		SetItemType(input.ItemType).
		SetStats(input.Stats).
		SetIsVisible(input.IsVisible).
		SetIsImmovable(input.IsImmovable).
		SetEffectType(input.EffectType).
		SetEffectValue(input.EffectValue).
		SetEffectDuration(input.EffectDuration).
		SetIsContainer(input.IsContainer).
		SetContainerCapacity(input.ContainerCapacity).
		SetIsLocked(input.IsLocked).
		SetRevealCondition(input.RevealCondition).
		SetArmorRating(input.ArmorRating).
		SetArmorType(input.ArmorType).
		SetRarity(input.Rarity).
		SetSkillRequirement(input.SkillRequirement).
		SetSkillRequirementLevel(input.SkillRequirementLevel).
		SetDamageDiceCount(input.DamageDiceCount).
		SetDamageDiceSides(input.DamageDiceSides).
		SetDamageBonus(input.DamageBonus).
		SetDamageType(input.DamageType).
		SetWeaponType(input.WeaponType).
		SetIsTwoHanded(input.IsTwoHanded)
	if input.Color != "" {
		builder = builder.SetColor(input.Color)
	}
	if input.KeyItemID != "" {
		builder = builder.SetKeyItemID(input.KeyItemID)
	}
	return builder.Save(ctx)
}

func (r *entEquipmentTemplateRepo) Update(ctx context.Context, id string, updates EquipmentTemplateUpdates) (*db.EquipmentTemplate, error) {
	builder := r.client.EquipmentTemplate.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.Description != nil {
		builder = builder.SetDescription(*updates.Description)
	}
	if updates.Slot != nil {
		builder = builder.SetSlot(*updates.Slot)
	}
	if updates.Level != nil {
		builder = builder.SetLevel(*updates.Level)
	}
	if updates.Weight != nil {
		builder = builder.SetWeight(*updates.Weight)
	}
	if updates.ItemType != nil {
		builder = builder.SetItemType(*updates.ItemType)
	}
	if updates.Stats != nil {
		builder = builder.SetStats(updates.Stats)
	}
	if updates.Color != nil {
		builder = builder.SetColor(*updates.Color)
	}
	if updates.IsVisible != nil {
		builder = builder.SetIsVisible(*updates.IsVisible)
	}
	if updates.IsImmovable != nil {
		builder = builder.SetIsImmovable(*updates.IsImmovable)
	}
	if updates.EffectType != nil {
		builder = builder.SetEffectType(*updates.EffectType)
	}
	if updates.EffectValue != nil {
		builder = builder.SetEffectValue(*updates.EffectValue)
	}
	if updates.EffectDuration != nil {
		builder = builder.SetEffectDuration(*updates.EffectDuration)
	}
	if updates.IsContainer != nil {
		builder = builder.SetIsContainer(*updates.IsContainer)
	}
	if updates.ContainerCapacity != nil {
		builder = builder.SetContainerCapacity(*updates.ContainerCapacity)
	}
	if updates.IsLocked != nil {
		builder = builder.SetIsLocked(*updates.IsLocked)
	}
	if updates.KeyItemID != nil {
		builder = builder.SetKeyItemID(*updates.KeyItemID)
	}
	if updates.RevealCondition != nil {
		builder = builder.SetRevealCondition(*updates.RevealCondition)
	}
	if updates.ArmorRating != nil {
		builder = builder.SetArmorRating(*updates.ArmorRating)
	}
	if updates.ArmorType != nil {
		builder = builder.SetArmorType(*updates.ArmorType)
	}
	if updates.Rarity != nil {
		builder = builder.SetRarity(*updates.Rarity)
	}
	if updates.SkillRequirement != nil {
		builder = builder.SetSkillRequirement(*updates.SkillRequirement)
	}
	if updates.SkillRequirementLevel != nil {
		builder = builder.SetSkillRequirementLevel(*updates.SkillRequirementLevel)
	}
	if updates.DamageDiceCount != nil {
		builder = builder.SetDamageDiceCount(*updates.DamageDiceCount)
	}
	if updates.DamageDiceSides != nil {
		builder = builder.SetDamageDiceSides(*updates.DamageDiceSides)
	}
	if updates.DamageBonus != nil {
		builder = builder.SetDamageBonus(*updates.DamageBonus)
	}
	if updates.DamageType != nil {
		builder = builder.SetDamageType(*updates.DamageType)
	}
	if updates.WeaponType != nil {
		builder = builder.SetWeaponType(*updates.WeaponType)
	}
	if updates.IsTwoHanded != nil {
		builder = builder.SetIsTwoHanded(*updates.IsTwoHanded)
	}
	return builder.Save(ctx)
}

func (r *entEquipmentTemplateRepo) Delete(ctx context.Context, id string) error {
	return r.client.EquipmentTemplate.DeleteOneID(id).Exec(ctx)
}