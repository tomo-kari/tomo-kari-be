package repo

import (
	"context"
	"fmt"
	"tomokari/internal/entity"

	"tomokari/pkg/postgres"
)

//const _defaultEntityCap = 64

// UserRepo -.
type UserRepo struct {
	*postgres.Postgres
}

// NewUserRepo -.
func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

// Create -.
func (r *UserRepo) Create(ctx context.Context, user entity.CreateUserRequestBody) error {
	sql, args, err := r.Builder.
		Insert("user").
		Columns("email, phone, password, date_of_birth, terms_of_service_id").
		Values(user.Email, user.Phone, user.Password, user.DateOfBirth, user.TermsOfServiceId).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Create - r.Builder: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	fmt.Printf("%+v\n", tag.String())
	if err != nil {
		return fmt.Errorf("UserRepo - Create - r.Pool.Query: %w", err)
	}

	return nil
}
