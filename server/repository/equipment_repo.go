package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/equipment"
	"herbst-server/db/room"
)

type entEquipmentRepo struct {
	client *db.Client
}

func NewEntEquipmentRepo(client *db.Client) EquipmentRepo {
	return &entEquipmentRepo{client: client}
}

func (r *entEquipmentRepo) Get(ctx context.Context, id int) (*db.Equipment, error) {
	return r.client.Equipment.Get(ctx, id)
}

func (r *entEquipmentRepo) ListByOwner(ctx context.Context, ownerID int) ([]*db.Equipment, error) {
	return r.client.Equipment.Query().
		Where(equipment.OwnerId(ownerID)).
		All(ctx)
}

func (r *entEquipmentRepo) ListByRoom(ctx context.Context, roomID int) ([]*db.Equipment, error) {
	return r.client.Equipment.Query().
		Where(equipment.HasRoomWith(room.ID(roomID))).
		All(ctx)
}

func (r *entEquipmentRepo) Update(ctx context.Context, id int, updates EquipmentUpdates) (*db.Equipment, error) {
	builder := r.client.Equipment.UpdateOneID(id)
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
	if updates.OwnerID != nil {
		builder = builder.SetOwnerId(*updates.OwnerID)
	}
	if updates.RoomID != nil {
		if *updates.RoomID == 0 && updates.ClearRoom {
			builder = builder.ClearRoom()
		} else {
			builder = builder.SetRoomID(*updates.RoomID)
		}
	}
	if updates.ClearRoom && updates.RoomID == nil {
		builder = builder.ClearRoom()
	}
	if updates.IsEquipped != nil {
		builder = builder.SetIsEquipped(*updates.IsEquipped)
	}
	if updates.IsImmovable != nil {
		builder = builder.SetIsImmovable(*updates.IsImmovable)
	}
	if updates.Color != nil {
		builder = builder.SetColor(*updates.Color)
	}
	if updates.IsVisible != nil {
		builder = builder.SetIsVisible(*updates.IsVisible)
	}
	if updates.ItemType != nil {
		builder = builder.SetItemType(*updates.ItemType)
	}
	if updates.Healing != nil {
		builder = builder.SetHealing(*updates.Healing)
	}
	if updates.Effect != nil {
		builder = builder.SetEffect(*updates.Effect)
	}
	if updates.ArmorRating != nil {
		builder = builder.SetArmorRating(*updates.ArmorRating)
	}
	if updates.ArmorType != nil {
		builder = builder.SetArmorType(*updates.ArmorType)
	}
	if updates.Stats != nil {
		builder = builder.SetStats(updates.Stats)
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
	if updates.ExpiresAt != nil {
		builder = builder.SetExpiresAt(*updates.ExpiresAt)
	}
	return builder.Save(ctx)
}

func (r *entEquipmentRepo) Delete(ctx context.Context, id int) error {
	return r.client.Equipment.DeleteOneID(id).Exec(ctx)
}

func (r *entEquipmentRepo) CountByTemplateID(ctx context.Context, templateID string) (int, error) {
	return r.client.Equipment.Query().
		Where(equipment.EquipmentTemplateIDEQ(templateID)).
		Count(ctx)
}