package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/trigger"
)

type entTriggerRepo struct {
	client *db.Client
}

func NewEntTriggerRepo(client *db.Client) TriggerRepo {
	return &entTriggerRepo{client: client}
}

func (r *entTriggerRepo) Get(ctx context.Context, id int) (*db.Trigger, error) {
	return r.client.Trigger.Get(ctx, id)
}

func (r *entTriggerRepo) GetWithEdges(ctx context.Context, id int) (*db.Trigger, error) {
	return r.client.Trigger.Query().
		Where(trigger.IDEQ(id)).
		WithEffect().
		WithRecipe().
		WithDialogNode().
		Only(ctx)
}

func (r *entTriggerRepo) List(ctx context.Context) ([]*db.Trigger, error) {
	return r.client.Trigger.Query().All(ctx)
}

func (r *entTriggerRepo) ListWithEdges(ctx context.Context) ([]*db.Trigger, error) {
	return r.client.Trigger.Query().
		Order(db.Asc(trigger.FieldName)).
		WithEffect().
		WithRecipe().
		WithDialogNode().
		All(ctx)
}

func (r *entTriggerRepo) ListByRoom(ctx context.Context, roomID int) ([]*db.Trigger, error) {
	return r.client.Trigger.Query().
		Where(trigger.RoomIDEQ(roomID)).
		WithEffect().
		WithRecipe().
		WithDialogNode().
		All(ctx)
}

func (r *entTriggerRepo) ListByEquipment(ctx context.Context, equipmentID int) ([]*db.Trigger, error) {
	return r.client.Trigger.Query().
		Where(trigger.EquipmentIDEQ(equipmentID)).
		WithEffect().
		WithRecipe().
		WithDialogNode().
		All(ctx)
}

func (r *entTriggerRepo) ListByTriggerType(ctx context.Context, triggerType string) ([]*db.Trigger, error) {
	return r.client.Trigger.Query().
		Where(trigger.TriggerTypeEQ(triggerType)).
		WithEffect().
		WithRecipe().
		WithDialogNode().
		All(ctx)
}

func (r *entTriggerRepo) ListByTargetType(ctx context.Context, targetType string) ([]*db.Trigger, error) {
	return r.client.Trigger.Query().
		Where(trigger.TargetTypeEQ(targetType)).
		WithEffect().
		WithRecipe().
		WithDialogNode().
		All(ctx)
}

func (r *entTriggerRepo) Create(ctx context.Context, input CreateTriggerInput) (*db.Trigger, error) {
	builder := r.client.Trigger.Create().
		SetName(input.Name).
		SetWorldID(input.WorldID).
		SetTriggerType(input.TriggerType).
		SetTargetType(input.TargetType).
		SetTargetID(input.TargetID).
		SetEnabled(input.Enabled)
	if input.RoomID != nil {
		builder = builder.SetRoomID(*input.RoomID)
	}
	if input.EquipmentID != nil {
		builder = builder.SetEquipmentID(*input.EquipmentID)
	}
	if input.Condition != "" {
		builder = builder.SetCondition(input.Condition)
	}
	return builder.Save(ctx)
}

func (r *entTriggerRepo) Update(ctx context.Context, id int, updates TriggerUpdates) (*db.Trigger, error) {
	builder := r.client.Trigger.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.WorldID != nil {
		builder = builder.SetWorldID(*updates.WorldID)
	}
	if updates.TriggerType != nil {
		builder = builder.SetTriggerType(*updates.TriggerType)
	}
	if updates.TargetType != nil {
		builder = builder.SetTargetType(*updates.TargetType)
	}
	if updates.TargetID != nil {
		builder = builder.SetTargetID(*updates.TargetID)
	}
	if updates.RoomID != nil {
		if *updates.RoomID == 0 {
			builder = builder.ClearRoomID()
		} else {
			builder = builder.SetRoomID(*updates.RoomID)
		}
	}
	if updates.EquipmentID != nil {
		if *updates.EquipmentID == 0 {
			builder = builder.ClearEquipmentID()
		} else {
			builder = builder.SetEquipmentID(*updates.EquipmentID)
		}
	}
	if updates.Condition != nil {
		builder = builder.SetCondition(*updates.Condition)
	}
	if updates.Enabled != nil {
		builder = builder.SetEnabled(*updates.Enabled)
	}
	return builder.Save(ctx)
}

func (r *entTriggerRepo) Delete(ctx context.Context, id int) error {
	return r.client.Trigger.DeleteOneID(id).Exec(ctx)
}
