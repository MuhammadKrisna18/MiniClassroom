package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"siakad-pro/internal/modules/matakuliah/domain"
	"siakad-pro/internal/shared/apperrors"
	"time"
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

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06d", rng.Intn(1000000))

	periodeID, err := s.repo.GetActivePeriodeID(ctx)
	if err != nil {
		return nil, apperrors.NewBadRequest("Tidak ada periode akademik yang aktif")
	}

	pengajuan := &domain.PengajuanMataKuliah{
		ID:           uuid.New().String(),
		PeriodeID:    periodeID,
		DosenID:      dosenID,
		MataKuliahID: req.MataKuliahID,
		Status:       "pending",
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

	if p.Status != "pending" {
		return apperrors.NewBadRequest("Pengajuan sudah tidak dalam status pending")
	}

	p.Status = "approved"
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

	if p.Status != "pending" {
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
