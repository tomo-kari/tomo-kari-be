// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Filter struct {
	Field string      `json:"field"`
	Op    string      `json:"op"`
	Value interface{} `json:"value"`
}

type Pagination struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}

type OrderBy struct {
	Field string `json:"field"`
}

type GetManyRequestBody struct {
	Filters    []Filter    `json:"filters"`
	Pagination *Pagination `json:"pagination"`
	OrderBy    []OrderBy   `json:"orderBy"`
}

// Postgres -.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

// New -.
func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

// Close -.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func (p *Postgres) WithFilters(selectBuilder squirrel.SelectBuilder, filters []Filter) squirrel.SelectBuilder {
	for _, filter := range filters {
		switch filter.Op {
		case "eq":
			selectBuilder = selectBuilder.Where(filter.Field+" = ?", filter.Value)
		case "neq":
			selectBuilder = selectBuilder.Where(filter.Field+" != ?", filter.Value)
		case "gt":
			selectBuilder = selectBuilder.Where(filter.Field+" > ?", filter.Value)
		case "lt":
			selectBuilder = selectBuilder.Where(filter.Field+" < ?", filter.Value)
		case "gte":
			selectBuilder = selectBuilder.Where(filter.Field+" >= ?", filter.Value)
		case "lte":
			selectBuilder = selectBuilder.Where(filter.Field+" <= ?", filter.Value)
		//case "contains":
		//    selectBuilder = selectBuilder.Where(filter.Field+" @> ?", filter.Value)
		case "str-contains":
			selectBuilder = selectBuilder.Where(filter.Field+" like ?", "%"+filter.Value.(string)+"%")
		case "str-starts-with":
			selectBuilder = selectBuilder.Where(filter.Field+" like ?", filter.Value.(string)+"%")
		case "str-ends-with":
			selectBuilder = selectBuilder.Where(filter.Field+" like ?", "%"+filter.Value.(string))
		case "in":
			selectBuilder = selectBuilder.Where(filter.Field+" in ?", filter.Value)
		}
	}

	return selectBuilder
}

func (p *Postgres) WithSorting(selectBuilder squirrel.SelectBuilder, orderBy []OrderBy) squirrel.SelectBuilder {
	for _, order := range orderBy {
		selectBuilder = selectBuilder.OrderBy(order.Field)
	}

	return selectBuilder
}

func (p *Postgres) WithPagination(selectBuilder squirrel.SelectBuilder, pagination *Pagination) squirrel.SelectBuilder {
	if pagination == nil {
		pagination = &Pagination{
			Limit:  10,
			Offset: 0,
		}
	}
	if pagination.Offset > 0 {
		pagination.Offset -= 1
	}
	if pagination.Offset < 0 {
		pagination.Offset = 0
	}
	selectBuilder = selectBuilder.Limit(pagination.Limit).Offset(pagination.Offset)

	return selectBuilder
}
