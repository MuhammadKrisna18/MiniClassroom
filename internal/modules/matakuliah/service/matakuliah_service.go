package service

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"siakad-pro/internal/modules/matakuliah/domain"
	"siakad-pro/internal/shared/apperrors"
)

type matakuliahService struct {
	repo domain.MataKuliahRepository
}

func NewMataKuliahService(repo domain.MataKuliahRepository) domain.MataKuliahService {
	return &matakuliahService{
		repo: repo,
	}
}

func (s *matakuliahService) CreateMataKuliah(ctx context.Context, req domain.CreateMataKuliahRequest) (*domain.MataKuliah, error) {
	if req.ProgramStudiID == "" {
		return nil, apperrors.NewBadRequest("Program Studi wajib diisi")
	}

	existing, err := s.repo.GetByNameAndProdi(ctx, req.Name, req.ProgramStudiID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengecek mata kuliah", err.Error())
	}

	if existing != nil {
		return nil, &apperrors.AppError{Code: http.StatusConflict, Message: "Mata kuliah sudah terdaftar di Program Studi ini"}
	}

	newMk := &domain.MataKuliah{
		ID:             uuid.New().String(),
		Name:           req.Name,
		SKS:            req.SKS,
		ProgramStudiID: req.ProgramStudiID,
	}

	if err := s.repo.Create(ctx, newMk); err != nil {
		return nil, apperrors.NewInternal("Gagal menyimpan mata kuliah", err.Error())
	}

	dosenIDs, err := s.repo.GetDosenIDsByProdi(ctx, req.ProgramStudiID)
	if err == nil && len(dosenIDs) > 0 {
		for _, dID := range dosenIDs {
			pengajuan := &domain.PengajuanMataKuliah{
				ID:           uuid.New().String(),
				DosenID:      dID,
				MataKuliahID: newMk.ID,
				Status:       domain.StatusOffered,
				Code:         "",
			}
			s.repo.CreatePengajuan(ctx, pengajuan)
		}
	}

	return newMk, nil
}

func (s *matakuliahService) GetMataKuliahList(ctx context.Context) ([]*domain.MataKuliah, error) {
	mkList, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil daftar mata kuliah", err.Error())
	}

	if mkList == nil {
		mkList = []*domain.MataKuliah{}
	}

	return mkList, nil
}

func (s *matakuliahService) GetMataKuliahForMahasiswa(ctx context.Context, userID string) ([]*domain.MataKuliah, error) {
	prodiID, err := s.repo.GetUserProdiID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil data prodi mahasiswa", err.Error())
	}
	if prodiID == nil || *prodiID == "" {
		return nil, apperrors.NewBadRequest("Mahasiswa belum memiliki Program Studi")
	}

	mkList, err := s.repo.GetByProdi(ctx, *prodiID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil daftar mata kuliah", err.Error())
	}
	return mkList, nil
}

func (s *matakuliahService) DeleteMataKuliah(ctx context.Context, id string) error {

	activeReqs, err := s.repo.GetActivePengajuanByMataKuliahID(ctx, id)
	if err != nil {
		return apperrors.NewInternal("Gagal mengecek status mata kuliah", err.Error())
	}
	if len(activeReqs) > 0 {
		return apperrors.NewBadRequest("Mata kuliah tidak dapat dihapus karena sudah diajukan atau diambil oleh dosen")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return apperrors.NewInternal("Gagal menghapus mata kuliah", err.Error())
	}
	return nil
}

func (s *matakuliahService) LepasMataKuliah(ctx context.Context, mkID string) error {
	activeReqs, err := s.repo.GetActivePengajuanByMataKuliahID(ctx, mkID)
	if err != nil {
		return apperrors.NewInternal("Gagal mengecek status mata kuliah", err.Error())
	}

	for _, req := range activeReqs {
		if err := s.repo.DeletePengajuan(ctx, req.ID); err != nil {
			return apperrors.NewInternal("Gagal melepas mata kuliah", err.Error())
		}
	}

	return nil
}
