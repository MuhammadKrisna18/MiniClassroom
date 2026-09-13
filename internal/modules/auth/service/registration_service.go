package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"siakad-pro/internal/modules/auth/domain"
	"siakad-pro/internal/shared/apperrors"
	"siakad-pro/pkg/utils"
)

func (s *authService) RegisterDosen(ctx context.Context, req domain.RegisterDosenRequest) (*domain.UserProfileResponse, error) {
	email := req.Username + "@DosenGO.id"

	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	nidStr := utils.GenerateRandomNumberString(5)

	newUser := &domain.User{
		ID:             uuid.New().String(),
		Name:           req.Name,
		Email:          email,
		NID:            &nidStr,
		Password:       string(hashedPassword),
		Role:           domain.RoleDosen,
		ProgramStudiID: &req.ProgramStudiID,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, errors.New("failed to create dosen account")
	}

	return toProfileResponse(newUser), nil
}

func (s *authService) GetDosenList(ctx context.Context) ([]*domain.UserProfileResponse, error) {
	users, err := s.repo.GetUsersByRole(ctx, domain.RoleDosen)
	if err != nil {
		return nil, errors.New("failed to fetch dosen list")
	}

	var res []*domain.UserProfileResponse
	for _, u := range users {
		res = append(res, toProfileResponse(u))
	}
	return res, nil
}

func (s *authService) DeleteDosen(ctx context.Context, id string) error {

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Akun dosen tidak ditemukan", err.Error())
	}
	if user.Role != "dosen" {
		return &apperrors.AppError{Code: 400, Message: "Akun ini bukan dosen"}
	}
	return s.repo.DeleteUser(ctx, id)
}

func (s *authService) RegisterMahasiswa(ctx context.Context, req domain.RegisterMahasiswaRequest) (*domain.UserProfileResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengenkripsi password", err.Error())
	}

	email := req.NRP + "@student.its.golang"

	existingUser, err := s.repo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, apperrors.NewConflict("NRP ini sudah terdaftar untuk mahasiswa lain", "")
	}

	user := &domain.User{
		ID:             uuid.New().String(),
		Name:           req.Name,
		NRP:            &req.NRP,
		Email:          email,
		Password:       string(hashedPassword),
		Role:           domain.RoleMahasiswa,
		ProgramStudiID: &req.ProgramStudiID,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, apperrors.NewInternal("Gagal mendaftarkan mahasiswa", err.Error())
	}

	return s.GetProfile(ctx, user.ID)
}

func (s *authService) GetMahasiswaList(ctx context.Context) ([]*domain.UserProfileResponse, error) {
	users, err := s.repo.GetUsersByRole(ctx, domain.RoleMahasiswa)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil data mahasiswa", err.Error())
	}

	var response []*domain.UserProfileResponse
	for _, user := range users {
		response = append(response, toProfileResponse(user))
	}
	return response, nil
}
