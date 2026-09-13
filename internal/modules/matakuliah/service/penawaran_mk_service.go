package service

import (
	"context"
	"github.com/google/uuid"
	"siakad-pro/internal/modules/matakuliah/domain"
	"siakad-pro/internal/shared/apperrors"
	"siakad-pro/pkg/utils"
)

func (s *matakuliahService) RequestMataKuliah(ctx context.Context, dosenID string, req domain.RequestMataKuliahPayload) (*domain.PengajuanMataKuliah, error) {

	activeReqs, err := s.repo.GetActivePengajuanByMataKuliahID(ctx, req.MataKuliahID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengecek status mata kuliah", err.Error())
	}
	if len(activeReqs) > 0 {
		if activeReqs[0].DosenID == dosenID {
			return nil, apperrors.NewBadRequest("Anda sudah mengajukan mata kuliah ini")
		}
		return nil, apperrors.NewBadRequest("Mata kuliah ini sudah diajukan atau diambil oleh dosen lain")
	}

	code := utils.GenerateRandomNumberString(6)

	periodeID, err := s.repo.GetActivePeriodeID(ctx)
	if err != nil {
		return nil, apperrors.NewBadRequest("Tidak ada periode akademik yang aktif")
	}

	pengajuan := &domain.PengajuanMataKuliah{
		ID:           uuid.New().String(),
		PeriodeID:    periodeID,
		DosenID:      dosenID,
		MataKuliahID: req.MataKuliahID,
		Status:       domain.StatusPending,
		Code:         code,
	}

	if err := s.repo.CreatePengajuan(ctx, pengajuan); err != nil {
		return nil, apperrors.NewInternal("Gagal mengajukan mata kuliah", err.Error())
	}

	return pengajuan, nil
}

func (s *matakuliahService) ApprovePengajuan(ctx context.Context, id string) error {
	p, err := s.repo.GetPengajuanByID(ctx, id)
	if err != nil {
		return apperrors.NewInternal("Gagal mengambil data pengajuan", err.Error())
	}
	if p == nil {
		return apperrors.NewNotFound("Pengajuan tidak ditemukan")
	}

	if p.Status != domain.StatusPending {
		return apperrors.NewBadRequest("Pengajuan sudah tidak dalam status pending")
	}

	p.Status = domain.StatusApproved
	if err := s.repo.UpdatePengajuan(ctx, p); err != nil {
		return apperrors.NewInternal("Gagal menyetujui pengajuan", err.Error())
	}
	return nil
}

func (s *matakuliahService) RejectPengajuan(ctx context.Context, id string) error {
	p, err := s.repo.GetPengajuanByID(ctx, id)
	if err != nil {
		return apperrors.NewInternal("Gagal mengambil data pengajuan", err.Error())
	}
	if p == nil {
		return apperrors.NewNotFound("Pengajuan tidak ditemukan")
	}

	if p.Status != domain.StatusPending {
		return apperrors.NewBadRequest("Pengajuan sudah tidak dalam status pending")
	}

	p.Status = domain.StatusRejected
	if err := s.repo.UpdatePengajuan(ctx, p); err != nil {
		return apperrors.NewInternal("Gagal menolak pengajuan", err.Error())
	}
	return nil
}

func (s *matakuliahService) AcceptOffer(ctx context.Context, id string, dosenID string) error {
	p, err := s.repo.GetPengajuanByID(ctx, id)
	if err != nil {
		return apperrors.NewInternal("Gagal mengambil data penawaran", err.Error())
	}
	if p == nil {
		return apperrors.NewNotFound("Penawaran tidak ditemukan")
	}

	if p.DosenID != dosenID {
		return apperrors.NewBadRequest("Anda tidak berhak menerima penawaran ini")
	}

	if p.Status != domain.StatusOffered {
		return apperrors.NewBadRequest("Penawaran sudah tidak valid")
	}

	activeReqs, _ := s.repo.GetActivePengajuanByMataKuliahID(ctx, p.MataKuliahID)
	for _, req := range activeReqs {
		if req.Status == domain.StatusApproved {

			s.repo.DeletePengajuan(ctx, p.ID)
			return apperrors.NewBadRequest("Mata kuliah ini sudah diambil oleh dosen lain")
		}
	}

	p.Status = domain.StatusApproved
	if err := s.repo.UpdatePengajuan(ctx, p); err != nil {
		return apperrors.NewInternal("Gagal menyetujui penawaran", err.Error())
	}

	if allReqs, err := s.repo.GetAllPengajuan(ctx); err == nil {
		for _, req := range allReqs {
			if req.MataKuliahID == p.MataKuliahID && req.ID != p.ID && req.Status == domain.StatusOffered {
				s.repo.DeletePengajuan(ctx, req.ID)
			}
		}
	}

	return nil
}

func (s *matakuliahService) RejectOffer(ctx context.Context, id string, dosenID string) error {
	p, err := s.repo.GetPengajuanByID(ctx, id)
	if err != nil {
		return apperrors.NewInternal("Gagal mengambil data penawaran", err.Error())
	}
	if p == nil {
		return apperrors.NewNotFound("Penawaran tidak ditemukan")
	}

	if p.DosenID != dosenID {
		return apperrors.NewBadRequest("Anda tidak berhak menolak penawaran ini")
	}

	if p.Status != domain.StatusOffered {
		return apperrors.NewBadRequest("Penawaran sudah tidak valid")
	}

	if err := s.repo.DeletePengajuan(ctx, id); err != nil {
		return apperrors.NewInternal("Gagal menolak penawaran", err.Error())
	}
	return nil
}

func (s *matakuliahService) GetMyPengajuan(ctx context.Context, dosenID string) ([]*domain.PengajuanMataKuliah, error) {
	list, err := s.repo.GetPengajuanByDosenID(ctx, dosenID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil riwayat pengajuan", err.Error())
	}
	if list == nil {
		list = []*domain.PengajuanMataKuliah{}
	}
	return list, nil
}

func (s *matakuliahService) GetAllPengajuan(ctx context.Context) ([]*domain.PengajuanMataKuliah, error) {
	list, err := s.repo.GetAllPengajuan(ctx)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil daftar pengajuan", err.Error())
	}
	if list == nil {
		list = []*domain.PengajuanMataKuliah{}
	}
	return list, nil
}
