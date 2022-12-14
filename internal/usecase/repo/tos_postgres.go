package repo

import (
	"context"
	"fmt"
	"tomokari/internal/entity"
	"tomokari/pkg/postgres"
)

// TOSRepo -.
type TOSRepo struct {
	*postgres.Postgres
}

// NewTOSRepo -.
func NewTOSRepo(pg *postgres.Postgres) *TOSRepo {
	return &TOSRepo{pg}
}

func (tos TOSRepo) GetByID(ctx context.Context, id int64) (*entity.TermsOfService, error) {
	sql, args, err := tos.Builder.
		Select("id, content").
		From((&entity.TermsOfService{}).Table()).
		Where("id = $1", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TOSRepo - GetByID - tos.Builder: %w", err)
	}

	rows, err := tos.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetByID - u.Pool.Query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var t entity.TermsOfService
		err := rows.Scan(&t.ID, &t.Content)
		if err != nil {
			return nil, fmt.Errorf("TOSRepo - GetByID - rows.Scan: %w", err)
		}
		return &t, nil
	}
	return nil, fmt.Errorf("TOSRepo - GetByID - not found")
}
