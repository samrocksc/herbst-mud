package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/user"
)

type entUserRepo struct {
	client *db.Client
}

func NewEntUserRepo(client *db.Client) UserRepo {
	return &entUserRepo{client: client}
}

func (r *entUserRepo) Get(ctx context.Context, id int) (*db.User, error) {
	return r.client.User.Get(ctx, id)
}

func (r *entUserRepo) GetByEmail(ctx context.Context, email string) (*db.User, error) {
	return r.client.User.Query().
		Where(user.Email(email)).
		Only(ctx)
}

func (r *entUserRepo) List(ctx context.Context) ([]*db.User, error) {
	return r.client.User.Query().All(ctx)
}

func (r *entUserRepo) Create(ctx context.Context, input CreateUserInput) (*db.User, error) {
	return r.client.User.Create().
		SetEmail(input.Email).
		SetPassword(input.Password).
		SetIsAdmin(input.IsAdmin).
		Save(ctx)
}

func (r *entUserRepo) Update(ctx context.Context, id int, updates UserUpdates) (*db.User, error) {
	builder := r.client.User.UpdateOneID(id)
	if updates.Email != nil {
		builder = builder.SetEmail(*updates.Email)
	}
	if updates.Password != nil {
		builder = builder.SetPassword(*updates.Password)
	}
	if updates.IsAdmin != nil {
		builder = builder.SetIsAdmin(*updates.IsAdmin)
	}
	return builder.Save(ctx)
}

func (r *entUserRepo) Delete(ctx context.Context, id int) error {
	return r.client.User.DeleteOneID(id).Exec(ctx)
}