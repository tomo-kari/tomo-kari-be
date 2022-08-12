// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"tomokari/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Translation -.
	Translation interface {
		Translate(context.Context, entity.Translation) (entity.Translation, error)
		History(context.Context) ([]entity.Translation, error)
	}

	// TranslationRepo -.
	TranslationRepo interface {
		Store(context.Context, entity.Translation) error
		GetHistory(context.Context) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}
)

type (
	User interface {
		Register(ctx context.Context, user entity.CreateUserRequestBody) (*entity.AuthUserResponse, Status, error)
		Login(ctx context.Context, user entity.LoginUserRequestBody) (*entity.AuthUserResponse, Status, error)
	}

	IUserRepo interface {
		Create(ctx context.Context, user entity.User) error
		GetByEmail(ctx context.Context, email string) (*entity.User, error)
		GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	}
)

type (
	ITOSRepo interface {
		GetByID(ctx context.Context, id int64) (*entity.TermsOfService, error)
	}
)
