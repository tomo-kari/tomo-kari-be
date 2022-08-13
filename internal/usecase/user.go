package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"tomokari/internal/constants"
	"tomokari/internal/entity"
	"tomokari/internal/utils"
	"tomokari/pkg/postgres"
)

// UserUseCase -.
type UserUseCase struct {
	userRepo IUserRepo
	tosRepo  ITOSRepo
}

// NewUserUseCase -.
func NewUserUseCase(u IUserRepo, tosRepo ITOSRepo) *UserUseCase {
	return &UserUseCase{
		userRepo: u,
		tosRepo:  tosRepo,
	}
}

func createAuthResponse(user *entity.User) *entity.AuthUserResponse {
	tokenData := entity.UserTokenData{
		ID:    user.ID,
		Email: user.Email,
		Phone: user.Phone,
		Role:  string(user.Role),
	}
	accessToken, _ := utils.GenerateToken(tokenData, constants.AccessToken)
	refreshToken, _ := utils.GenerateToken(tokenData, constants.RefreshToken)
	return &entity.AuthUserResponse{
		ID: user.ID,
		BasicInfo: entity.BasicInfo{
			Email: user.Email,
			Phone: user.Phone,
		},
		DateOfBirth: user.DateOfBirth,
		Token: entity.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
}

// Register - create a new user.
func (uc *UserUseCase) Register(ctx context.Context, registerInfo entity.CreateUserRequestBody) (*entity.AuthUserResponse, Status, error) {
	tos, err := uc.tosRepo.GetByID(ctx, registerInfo.TermsOfServiceId)
	if tos == nil {
		return nil, http.StatusBadRequest, fmt.Errorf(constants.TermsOfServiceNotAcceptedErrorMessage)
	}

	existedUser, err := uc.userRepo.GetByEmail(ctx, registerInfo.Email)
	if existedUser != nil {
		return nil, http.StatusBadRequest, fmt.Errorf(constants.DuplicatedEmailErrorMessage)
	}

	existedUser, err = uc.userRepo.GetByPhone(ctx, registerInfo.Phone)
	if existedUser != nil {
		return nil, http.StatusBadRequest, fmt.Errorf(constants.DuplicatedPhoneErrorMessage)
	}

	var user entity.User
	user.Email = registerInfo.Email
	user.Phone = registerInfo.Phone
	user.Password = registerInfo.Password
	user.DateOfBirth, err = time.Parse(time.RFC3339, registerInfo.DateOfBirth)
	user.HashPassword()
	user.Role = entity.USER
	err = uc.userRepo.Create(ctx, &user)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	respData := createAuthResponse(&user)
	return respData, http.StatusOK, nil
}

// Login -.
func (uc *UserUseCase) Login(ctx context.Context, loginInfo entity.LoginUserRequestBody) (*entity.AuthUserResponse, Status, error) {
	user, _ := uc.userRepo.GetByEmail(ctx, loginInfo.Email)
	if user == nil {
		return nil, http.StatusBadRequest, fmt.Errorf(constants.IncorrectUserCredentialsErrorMessage)
	}
	if !user.CheckPasswordHash(loginInfo.Password) {
		return nil, http.StatusBadRequest, fmt.Errorf(constants.IncorrectUserCredentialsErrorMessage)
	}
	respData := createAuthResponse(user)
	return respData, http.StatusOK, nil
}

func (uc *UserUseCase) GetCandidates(ctx context.Context, filter postgres.GetManyRequestBody) ([]entity.Map, error) {
	return uc.userRepo.GetMany(ctx, filter)
}
