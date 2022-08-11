package usecase

import (
	"context"
	"fmt"

	"tomokari/internal/entity"
)

// UserUseCase -.
type UserUseCase struct {
	repo IUserRepo
}

// NewUserUseCase -.
func NewUserUseCase(r IUserRepo) *UserUseCase {
	return &UserUseCase{
		repo: r,
	}
}

// Register - create a new user.
func (uc *UserUseCase) Register(ctx context.Context, user entity.CreateUserRequestBody) error {
	err := uc.repo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("UserUseCase - Create - s.repo.Create: %w", err)
	}

	return nil
}

// Login -.
func (uc *TranslationUseCase) Login(ctx context.Context, t entity.Translation) (entity.Translation, error) {
	translation, err := uc.webAPI.Translate(t)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationUseCase - Translate - s.webAPI.Translate: %w", err)
	}

	err = uc.repo.Store(context.Background(), translation)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationUseCase - Translate - s.repo.Store: %w", err)
	}

	return translation, nil
}
