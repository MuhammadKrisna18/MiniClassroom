package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"siakad-pro/config"
	"siakad-pro/internal/modules/auth/domain"
	"siakad-pro/internal/shared/apperrors"
)

type authService struct {
	repo domain.AuthRepository
	cfg  *config.Config
}

func NewAuthService(repo domain.AuthRepository, cfg *config.Config) domain.AuthService {
	return &authService{repo: repo, cfg: cfg}
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, apperrors.NewUnauthorized("invalid email or password")
		}
		return nil, apperrors.NewInternal("login failed", err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, apperrors.NewUnauthorized("invalid email or password")
	}

	claims := jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, apperrors.NewInternal("failed to generate token", err.Error())
	}

	return &domain.LoginResponse{
		Token: t,
		Role:  user.Role,
	}, nil
}

// toProfileResponse maps a User entity to a UserProfileResponse DTO.
func toProfileResponse(user *domain.User) *domain.UserProfileResponse {
	return &domain.UserProfileResponse{
		ID:             user.ID,
		Name:           user.Name,
		Nickname:       user.Nickname,
		NID:            user.NID,
		NRP:            user.NRP,
		Email:          user.Email,
		Role:           user.Role,
		ProgramStudiID: user.ProgramStudiID,
		ProgramStudi:   user.ProgramStudi,
		PhotoURL:       user.PhotoURL,
		CreatedAt:      user.CreatedAt,
	}
}
