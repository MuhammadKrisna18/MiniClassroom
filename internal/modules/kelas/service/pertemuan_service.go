package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"siakad-pro/internal/modules/kelas/domain"
	"siakad-pro/internal/shared/apperrors"
	"siakad-pro/pkg/utils"
)

func (s *kelasService) MulaiPertemuan(ctx context.Context, pengajuanID string, judul string) (*domain.Pertemuan, error) {
	existing, err := s.repo.GetPertemuanByPengajuanID(ctx, pengajuanID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil data pertemuan", err.Error())
	}

	nomorNext := len(existing) + 1
	if nomorNext > domain.MaxPertemuan {
		return nil, apperrors.NewBadRequest("Semua 16 pertemuan sudah tercatat untuk kelas ini")
	}

	kodeAbsensi := utils.GenerateRandomNumberString(6)

	p := &domain.Pertemuan{
		ID:             uuid.NewString(),
		PengajuanID:    pengajuanID,
		Judul:          judul,
		Tanggal:        time.Now(),
		WaktuMulai:     time.Now(),
		Status:         domain.PertemuanStatusBerlangsung,
		KodeAbsensi:    kodeAbsensi,
		NomorPertemuan: nomorNext,
	}
	if err := s.repo.CreatePertemuan(ctx, p); err != nil {
		return nil, apperrors.NewInternal("Gagal memulai pertemuan", err.Error())
	}
	return p, nil
}

func (s *kelasService) AkhiriPertemuan(ctx context.Context, pertemuanID string) error {
	p, err := s.repo.GetPertemuanByID(ctx, pertemuanID)
	if err != nil {
		return apperrors.NewNotFound("Pertemuan tidak ditemukan")
	}
	if p.Status != domain.PertemuanStatusBerlangsung {
		return apperrors.NewBadRequest("Pertemuan sudah selesai atau tidak aktif")
	}
	now := time.Now()
	p.Status = domain.PertemuanStatusSelesai
	p.WaktuSelesai = &now
	if err := s.repo.UpdatePertemuan(ctx, p); err != nil {
		return apperrors.NewInternal("Gagal mengakhiri pertemuan", err.Error())
	}
	return nil
}

func (s *kelasService) GetPertemuanByPengajuan(ctx context.Context, pengajuanID string) ([]*domain.Pertemuan, error) {
	return s.repo.GetPertemuanByPengajuanID(ctx, pengajuanID)
}
