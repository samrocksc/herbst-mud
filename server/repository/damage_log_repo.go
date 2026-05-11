package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/damagelog"
)

type entDamageLogRepo struct {
	client *db.Client
}

func NewEntDamageLogRepo(client *db.Client) DamageLogRepo {
	return &entDamageLogRepo{client: client}
}

func (r *entDamageLogRepo) Create(ctx context.Context, attackerID, targetID, damage int) (*db.DamageLog, error) {
	return r.client.DamageLog.Create().
		SetAttackerID(attackerID).
		SetTargetID(targetID).
		SetDamage(damage).
		Save(ctx)
}

func (r *entDamageLogRepo) ListByCharacter(ctx context.Context, charID int, limit int) ([]*db.DamageLog, error) {
	return r.client.DamageLog.Query().
		Where(
			damagelog.Or(
				damagelog.AttackerID(charID),
				damagelog.TargetID(charID),
			),
		).
		Limit(limit).
		All(ctx)
}