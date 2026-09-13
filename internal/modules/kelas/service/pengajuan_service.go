package service

import (
	"context"

	"github.com/google/uuid"
	"siakad-pro/internal/modules/kelas/domain"
	"siakad-pro/internal/shared/apperrors"
	"siakad-pro/pkg/utils"
)

func (s *kelasService) RequestKelas(ctx context.Context, dosenID string, req domain.RequestKelasPayload) (*domain.PengajuanKelas, error) {
	kelas, err := s.repo.GetByID(ctx, req.KelasID)
	if err != nil {
		return nil, apperrors.NewNotFound("Kelas tidak ditemukan")
	}

	hasMK, err := s.repo.IsMataKuliahValidForKelas(ctx, dosenID, req.MataKuliahID, kelas.ProgramStudiID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal memvalidasi mata kuliah dosen", err.Error())
	}
	if !hasMK {
		return nil, apperrors.NewBadRequest("Mata kuliah yang dipilih tidak valid atau bukan berasal dari program studi kelas ini.")
	}

	activePengajuan, err := s.repo.GetActivePengajuanByKelasID(ctx, req.KelasID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengecek status kelas", err.Error())
	}

	for _, p := range activePengajuan {
		if p.Status == domain.StatusApproved {
			return nil, apperrors.NewBadRequest("Kelas ini sudah disetujui untuk dosen lain")
		}
		if p.DosenID == dosenID && p.Status == domain.StatusPending {
			return nil, apperrors.NewBadRequest("Anda sudah mengajukan kelas ini")
		}
	}

	dosenPengajuan, err := s.repo.GetPengajuanByDosenID(ctx, dosenID)
	if err == nil {
		for _, dp := range dosenPengajuan {
			if dp.Status != domain.StatusRejected && dp.Kelas != nil {
				if dp.Kelas.Hari == kelas.Hari {
					if kelas.JamMulai < dp.Kelas.JamSelesai && dp.Kelas.JamMulai < kelas.JamSelesai {
						return nil, apperrors.NewBadRequest("Jadwal kelas bentrok dengan kelas lain yang sudah Anda ajukan/ambil (" + dp.Kelas.Name + ")")
					}
				}
			}
		}
	}

	periodeID, err := s.repo.GetActivePeriodeID(ctx)
	if err != nil {
		return nil, apperrors.NewBadRequest("Tidak ada periode akademik yang aktif")
	}

	pengajuan := &domain.PengajuanKelas{
		ID:           uuid.New().String(),
		PeriodeID:    periodeID,
		DosenID:      dosenID,
		KelasID:      req.KelasID,
		MataKuliahID: req.MataKuliahID,
		Status:       domain.StatusPending,
		Code:         utils.GenerateRandomString(6),
	}

	if err := s.repo.CreatePengajuan(ctx, pengajuan); err != nil {
		return nil, apperrors.NewInternal("Gagal mengajukan kelas", err.Error())
	}

	pengajuan.Kelas = kelas
	return pengajuan, nil
}

func (s *kelasService) ApprovePengajuan(ctx context.Context, id string) error {
	p, err := s.repo.GetPengajuanByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Pengajuan tidak ditemukan")
	}

	if p.Status != domain.StatusPending {
		return apperrors.NewBadRequest("Hanya pengajuan berstatus pending yang dapat disetujui")
	}

	activePengajuan, err := s.repo.GetActivePengajuanByKelasID(ctx, p.KelasID)
	if err != nil {
		return apperrors.NewInternal("Gagal mengecek status kelas", err.Error())
	}

	for _, ap := range activePengajuan {
		if ap.Status == domain.StatusApproved && ap.ID != id {
			return apperrors.NewBadRequest("Kelas sudah disetujui untuk dosen lain")
		}
	}

	p.Status = domain.StatusApproved
	if err := s.repo.UpdatePengajuan(ctx, p); err != nil {
		return apperrors.NewInternal("Gagal menyetujui pengajuan", err.Error())
	}

	for _, ap := range activePengajuan {
		if ap.ID != id && ap.Status == domain.StatusPending {
			ap.Status = domain.StatusRejected
			s.repo.UpdatePengajuan(ctx, ap)
		}
	}

	return nil
}

func (s *kelasService) RejectPengajuan(ctx context.Context, id string) error {
	p, err := s.repo.GetPengajuanByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Pengajuan tidak ditemukan")
	}

	if p.Status != domain.StatusPending {
		return apperrors.NewBadRequest("Hanya pengajuan berstatus pending yang dapat ditolak")
	}

	p.Status = domain.StatusRejected
	if err := s.repo.UpdatePengajuan(ctx, p); err != nil {
		return apperrors.NewInternal("Gagal menolak pengajuan", err.Error())
	}

	return nil
}

func (s *kelasService) GetMyPengajuan(ctx context.Context, dosenID string) ([]*domain.PengajuanKelas, error) {
	return s.repo.GetPengajuanByDosenID(ctx, dosenID)
}

func (s *kelasService) GetAllPengajuan(ctx context.Context) ([]*domain.PengajuanKelas, error) {
	return s.repo.GetAllPengajuan(ctx)
}
