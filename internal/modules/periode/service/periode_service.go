package service

import (
	"context"

	"github.com/google/uuid"
	"siakad-pro/internal/modules/periode/domain"
	"siakad-pro/internal/shared/apperrors"
)

type periodeService struct {
	repo domain.PeriodeRepository
}

func NewPeriodeService(repo domain.PeriodeRepository) domain.PeriodeService {
	return &periodeService{repo: repo}
}

func (s *periodeService) Create(ctx context.Context, req domain.CreatePeriodeRequest) (*domain.PeriodeAkademik, error) {
	p := &domain.PeriodeAkademik{
		ID:    uuid.New().String(),
		Tahun: req.Tahun,
		Jenis: req.Jenis,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, apperrors.NewInternal("Gagal membuat periode akademik: " + err.Error())
	}

	return p, nil
}

func (s *periodeService) GetAll(ctx context.Context) ([]*domain.PeriodeAkademik, error) {
	return s.repo.GetAll(ctx)
}

func (s *periodeService) GetActive(ctx context.Context) (*domain.PeriodeAkademik, error) {
	p, err := s.repo.GetActive(ctx)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil periode aktif")
	}
	if p == nil {
		return nil, apperrors.NewNotFound("Belum ada periode akademik yang aktif")
	}
	return p, nil
}

func (s *periodeService) Delete(ctx context.Context, id string) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Periode tidak ditemukan")
	}
	if p.IsActive {
		return apperrors.NewBadRequest("Tidak bisa menghapus periode yang sedang aktif")
	}
	return s.repo.Delete(ctx, id)
}

func (s *periodeService) SetActive(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Periode tidak ditemukan")
	}

	if err := s.repo.SetActive(ctx, id); err != nil {
		return apperrors.NewInternal("Gagal mengaktifkan periode")
	}

	return nil
}
