package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/gender"
)

type entGenderRepo struct {
	client *db.Client
}

func NewEntGenderRepo(client *db.Client) GenderRepo {
	return &entGenderRepo{client: client}
}

func (r *entGenderRepo) Get(ctx context.Context, id int) (*db.Gender, error) {
	return r.client.Gender.Get(ctx, id)
}

func (r *entGenderRepo) GetByName(ctx context.Context, name string) (*db.Gender, error) {
	return r.client.Gender.Query().
		Where(gender.Name(name)).
		Only(ctx)
}

func (r *entGenderRepo) GetByWorld(ctx context.Context, name, worldID string) (*db.Gender, error) {
	return r.client.Gender.Query().
		Where(gender.Name(name), gender.WorldID(worldID)).
		Only(ctx)
}

func (r *entGenderRepo) List(ctx context.Context, worldID string) ([]*db.Gender, error) {
	return r.client.Gender.Query().
		Where(gender.WorldID(worldID)).
		All(ctx)
}

func (r *entGenderRepo) Create(ctx context.Context, input CreateGenderInput) (*db.Gender, error) {
	created, err := r.client.Gender.Create().
		SetName(input.Name).
		SetDisplayName(input.DisplayName).
		SetSubjectPronoun(input.SubjectPronoun).
		SetObjectPronoun(input.ObjectPronoun).
		SetPossessivePronoun(input.PossessivePronoun).
		SetWorldID(input.WorldID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.client.Gender.Query().
		Where(gender.IDEQ(created.ID)).
		Only(ctx)
}

func (r *entGenderRepo) Update(ctx context.Context, id int, updates GenderUpdates) (*db.Gender, error) {
	builder := r.client.Gender.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.DisplayName != nil {
		builder = builder.SetDisplayName(*updates.DisplayName)
	}
	if updates.SubjectPronoun != nil {
		builder = builder.SetSubjectPronoun(*updates.SubjectPronoun)
	}
	if updates.ObjectPronoun != nil {
		builder = builder.SetObjectPronoun(*updates.ObjectPronoun)
	}
	if updates.PossessivePronoun != nil {
		builder = builder.SetPossessivePronoun(*updates.PossessivePronoun)
	}
	if updates.WorldID != nil {
		builder = builder.SetWorldID(*updates.WorldID)
	}
	_, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.client.Gender.Query().
		Where(gender.IDEQ(id)).
		Only(ctx)
}

func (r *entGenderRepo) Delete(ctx context.Context, id int) error {
	return r.client.Gender.DeleteOneID(id).Exec(ctx)
}