package persistent

import (
    "context"

    "github.com/Masterminds/squirrel"
    "github.com/evrone/go-clean-template/internal/entity"
    "github.com/evrone/go-clean-template/pkg/postgres"
)

type UserRepo struct {
    db *postgres.Postgres
}

func NewUserRepo(db *postgres.Postgres) *UserRepo {
    return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user entity.User) error {
    query, args, err := r.db.Builder.Insert("users").
        Columns("email", "password").
        Values(user.Email, user.Password).
        Suffix("RETURNING id").
        ToSql()
    if err != nil {
        return err
    }
    return r.db.Pool.QueryRow(ctx, query, args...).Scan(&user.ID)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
    var user entity.User
    query, args, err := r.db.Builder.Select("id", "email", "password").From("users").Where(squirrel.Eq{"email": email}).ToSql()
    if err != nil {
        return user, err
    }
    err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Email, &user.Password)
    return user, err
}