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
func (u *UserRepo) Create(ctx context.Context, user entity.User) error {
	sql, args, err := u.Builder.
		Insert((&entity.User{}).Table()).
		Columns("email, phone, password, role, date_of_birth, description").
		Values(user.Email, user.Phone, user.Password, user.Role, user.DateOfBirth, user.Description).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Create - u.Builder: %w", err)
	}

	tag, err := u.Pool.Exec(ctx, sql, args...)
	fmt.Printf("%+v\n", tag.String())
	if err != nil {
		return fmt.Errorf("UserRepo - Create - u.Pool.Query: %w", err)
	}

	return nil
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	sql, args, err := u.Builder.
		Select("id, email, phone, role, password, date_of_birth").
		From((&entity.User{}).Table()).
		Where("email = $1", email).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetByEmail - u.Builder: %w", err)
	}

	rows, err := u.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetByEmail - u.Pool.Query: %w", err)
	}
	defer rows.Close()
	var user *entity.User
	for rows.Next() {
		user = &entity.User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Phone, &user.Role, &user.Password, &user.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("UserRepo - GetByEmail - rows.Scan: %w", err)
		}
	}
	return user, nil
}

func (u *UserRepo) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	sql, args, err := u.Builder.
		Select("id, email, phone, role, password, date_of_birth").
		From((&entity.User{}).Table()).
		Where("phone = $1", phone).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetByPhone - u.Builder: %w", err)
	}

	rows, err := u.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetByPhone - u.Pool.Query: %w", err)
	}
	defer rows.Close()

	var user *entity.User
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Email, &user.Phone, &user.Role, &user.Password, &user.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("UserRepo - GetByPhone - rows.Scan: %w", err)
		}
	}
	return user, nil
}
