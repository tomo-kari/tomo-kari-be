package repo

import (
	"context"
	"fmt"
	"time"
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
func (u *UserRepo) Create(ctx context.Context, user *entity.User) error {
	sql, args, err := u.Builder.
		Insert((&entity.User{}).Table()).
		Columns("email, phone, password, role, date_of_birth, description").
		Values(user.Email, user.Phone, user.Password, user.Role, user.DateOfBirth, user.Description).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Create - u.Builder: %w", err)
	}

	sql = fmt.Sprintf("%s RETURNING id", sql)

	//tag, err := u.Pool.Exec(ctx, sql, args...)
	err = u.Pool.QueryRow(ctx, sql, args...).Scan(&user.ID)
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

func (u *UserRepo) GetMany(ctx context.Context, filter postgres.GetManyRequestBody) ([]entity.Map, error) {
	builder := u.Builder.
		Select("id, email, phone, role, date_of_birth").
		From((&entity.User{}).Table())
	builder = u.WithFilters(builder, filter.Filters)
	builder = u.WithSorting(builder, filter.OrderBy)
	builder = u.WithPagination(builder, filter.Pagination)
	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetMany - u.Builder: %w", err)
	}
	rows, err := u.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetMany - u.Pool.Query: %w", err)
	}
	defer rows.Close()
	var users []entity.Map
	for rows.Next() {
		var id int64
		var email string
		var phone string
		var role string
		var dateOfBirth time.Time
		err := rows.Scan(&id, &email, &phone, &role, &dateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("UserRepo - GetMany - rows.Scan: %w", err)
		}
		user := entity.Map{
			"id":          id,
			"email":       email,
			"phone":       phone,
			"role":        role,
			"dateOfBirth": dateOfBirth,
		}
		users = append(users, user)
	}
	return users, nil
}
