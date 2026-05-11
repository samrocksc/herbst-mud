package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/applog"
)

type entAppLogRepo struct {
	client *db.Client
}

func NewEntAppLogRepo(client *db.Client) AppLogRepo {
	return &entAppLogRepo{client: client}
}

func (r *entAppLogRepo) Create(ctx context.Context, level, message, service string, charID, roomID *int, templateID *string, metadata map[string]interface{}) (*db.AppLog, error) {
	builder := r.client.AppLog.Create().
		SetLevel(level).
		SetMessage(message).
		SetNillableService(&service)
	if charID != nil {
		builder = builder.SetNillableCharacterID(charID)
	}
	if roomID != nil {
		builder = builder.SetNillableRoomID(roomID)
	}
	if templateID != nil {
		builder = builder.SetNillableTemplateID(templateID)
	}
	if metadata != nil {
		builder = builder.SetMetadata(metadata)
	}
	return builder.Save(ctx)
}

func (r *entAppLogRepo) List(ctx context.Context, filter LogFilter) ([]*db.AppLog, int, error) {
	query := r.client.AppLog.Query()
	if filter.Level != "" {
		query = query.Where(applog.LevelEQ(filter.Level))
	}
	if filter.Service != "" {
		query = query.Where(applog.ServiceEQ(filter.Service))
	}
	if filter.CharacterID != nil {
		query = query.Where(applog.CharacterID(*filter.CharacterID))
	}
	if filter.RoomID != nil {
		query = query.Where(applog.RoomID(*filter.RoomID))
	}
	if filter.TemplateID != nil {
		query = query.Where(applog.TemplateID(*filter.TemplateID))
	}
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	results, err := query.All(ctx)
	return results, count, err
}

func (r *entAppLogRepo) ListServices(ctx context.Context) ([]string, error) {
	logs, err := r.client.AppLog.Query().
		Unique(true).
		Select(applog.FieldService).
		All(ctx)
	if err != nil {
		return nil, err
	}
	services := make([]string, 0, len(logs))
	for _, l := range logs {
		services = append(services, l.Service)
	}
	return services, nil
}