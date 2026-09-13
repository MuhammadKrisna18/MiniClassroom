package service

import (
	"context"

	"github.com/google/uuid"
	"siakad-pro/internal/modules/auth/domain"
	"siakad-pro/internal/shared/apperrors"
)

func (s *authService) GetProfile(ctx context.Context, id string) (*domain.UserProfileResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	profile := toProfileResponse(user)

	pendingReq, err := s.repo.GetPendingEmailRequestByUserID(ctx, id)
	if err == nil && pendingReq != nil {
		profile.PendingEmail = &pendingReq.NewEmail
	}

	return profile, nil
}

func (s *authService) UpdateProfile(ctx context.Context, id string, req domain.UpdateProfileRequest) (*domain.UserProfileResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFound("Akun tidak ditemukan", err.Error())
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Nickname != "" {
		user.Nickname = &req.Nickname
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, apperrors.NewInternal("Gagal mengupdate profil", err.Error())
	}

	return s.GetProfile(ctx, id)
}

func (s *authService) UpdateProfilePhoto(ctx context.Context, id string, photoURL string) (*domain.UserProfileResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFound("Akun tidak ditemukan", err.Error())
	}

	user.PhotoURL = &photoURL

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, apperrors.NewInternal("Gagal mengupdate foto profil", err.Error())
	}

	return s.GetProfile(ctx, id)
}

func (s *authService) RequestEmailChange(ctx context.Context, userID string, req domain.EmailChangeRequestPayload) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.NewNotFound("Akun tidak ditemukan", err.Error())
	}

	if user.Email == req.NewEmail {
		return apperrors.NewBadRequest("Email baru tidak boleh sama dengan email saat ini")
	}

	existingUser, err := s.repo.GetByEmail(ctx, req.NewEmail)
	if err == nil && existingUser.ID != user.ID {
		return apperrors.NewBadRequest("Email sudah digunakan oleh pengguna lain")
	}

	existingReq, err := s.repo.GetPendingEmailRequestByUserID(ctx, userID)
	if err == nil && existingReq != nil {

		existingReq.NewEmail = req.NewEmail
		return s.repo.UpdateEmailChangeRequest(ctx, existingReq)
	}

	newReq := &domain.EmailChangeRequest{
		ID:       uuid.New().String(),
		UserID:   userID,
		NewEmail: req.NewEmail,
		Status:   "pending",
	}

	return s.repo.CreateEmailChangeRequest(ctx, newReq)
}

func (s *authService) GetPendingEmailRequests(ctx context.Context) ([]*domain.EmailChangeRequest, error) {
	return s.repo.GetAllPendingEmailRequests(ctx)
}

func (s *authService) ReviewEmailRequest(ctx context.Context, requestID string, approve bool) error {
	req, err := s.repo.GetEmailChangeRequestByID(ctx, requestID)
	if err != nil {
		return apperrors.NewNotFound("Request tidak ditemukan", err.Error())
	}
	if req.Status != "pending" {
		return apperrors.NewBadRequest("Request sudah diproses")
	}

	if approve {
		req.Status = "approved"

		user, err := s.repo.GetByID(ctx, req.UserID)
		if err != nil {
			return apperrors.NewNotFound("User tidak ditemukan", err.Error())
		}

		user.Email = req.NewEmail
		if err := s.repo.Update(ctx, user); err != nil {
			return apperrors.NewInternal("Gagal mengupdate email user", err.Error())
		}
	} else {
		req.Status = "rejected"
	}

	return s.repo.UpdateEmailChangeRequest(ctx, req)
}
