package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"siakad-pro/internal/modules/semester/domain"
	"siakad-pro/internal/shared/apperrors"
)

type semesterService struct {
	repo       domain.SemesterRepository
	mkProvider domain.MataKuliahProvider
}

func NewSemesterService(repo domain.SemesterRepository, mkProvider domain.MataKuliahProvider) domain.SemesterService {
	return &semesterService{
		repo:       repo,
		mkProvider: mkProvider,
	}
}

func (s *semesterService) Create(ctx context.Context, req domain.CreateSemesterRequest) (*domain.Semester, error) {
	if req.Nomor < domain.MinSemester || req.Nomor > domain.MaxSemester {
		return nil, apperrors.NewBadRequest(fmt.Sprintf("Nomor semester harus antara %d dan %d", domain.MinSemester, domain.MaxSemester))
	}

	if req.MinSKS >= req.MaxSKS {
		return nil, apperrors.NewBadRequest("Minimum SKS harus lebih kecil dari Maksimum SKS")
	}

	existing, _ := s.repo.GetByNomor(ctx, req.Nomor)
	if existing != nil {
		return nil, apperrors.NewConflict(fmt.Sprintf("Semester %d sudah ada", req.Nomor))
	}

	semester := &domain.Semester{
		Nomor:  req.Nomor,
		MinSKS: req.MinSKS,
		MaxSKS: req.MaxSKS,
	}

	if err := s.repo.Create(ctx, semester); err != nil {
		return nil, apperrors.NewInternal("Gagal membuat semester: " + err.Error())
	}

	return semester, nil
}

func (s *semesterService) GetAll(ctx context.Context) ([]*domain.Semester, error) {
	return s.repo.GetAll(ctx)
}

func (s *semesterService) GetByID(ctx context.Context, id string) (*domain.Semester, error) {
	sem, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFound("Semester tidak ditemukan")
	}
	return sem, nil
}

func (s *semesterService) Update(ctx context.Context, id string, req domain.UpdateSemesterRequest) (*domain.Semester, error) {
	sem, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFound("Semester tidak ditemukan")
	}

	if req.MinSKS != nil {
		sem.MinSKS = *req.MinSKS
	}
	if req.MaxSKS != nil {
		sem.MaxSKS = *req.MaxSKS
	}

	if sem.MinSKS >= sem.MaxSKS {
		return nil, apperrors.NewBadRequest("Minimum SKS harus lebih kecil dari Maksimum SKS")
	}

	if err := s.repo.Update(ctx, sem); err != nil {
		return nil, apperrors.NewInternal("Gagal mengupdate semester")
	}

	return sem, nil
}

func (s *semesterService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Semester tidak ditemukan")
	}
	return s.repo.Delete(ctx, id)
}



func (s *semesterService) AssignMataKuliah(ctx context.Context, semesterID string, req domain.AssignMataKuliahRequest) (*domain.SemesterMataKuliah, error) {
	sem, err := s.repo.GetByID(ctx, semesterID)
	if err != nil {
		return nil, apperrors.NewNotFound("Semester tidak ditemukan")
	}

	for _, smk := range sem.MataKuliah {
		if smk.MataKuliahID == req.MataKuliahID {
			return nil, apperrors.NewConflict("Mata kuliah sudah ada di semester ini")
		}
	}

	mkInfo, err := s.mkProvider.GetMataKuliahByID(ctx, req.MataKuliahID)
	if err != nil || mkInfo == nil {
		return nil, apperrors.NewBadRequest("Mata kuliah tidak valid atau prodi tidak ditemukan")
	}

	var mkIDs []string
	for _, smk := range sem.MataKuliah {
		mkIDs = append(mkIDs, smk.MataKuliahID)
	}

	totalSKS := 0
	if len(mkIDs) > 0 {
		assignedMKs, err := s.mkProvider.GetMataKuliahByIDs(ctx, mkIDs)
		if err != nil {
			return nil, apperrors.NewInternal("Gagal mengambil data mata kuliah semester")
		}
		for _, m := range assignedMKs {
			if m.ProgramStudiID == mkInfo.ProgramStudiID {
				totalSKS += m.SKS
			}
		}
	}

	// Find the specific SKS limit for this prodi in this semester
	maxSKS := sem.MaxSKS
	for _, p := range sem.SKSProdi {
		if p.ProgramStudiID == mkInfo.ProgramStudiID {
			maxSKS = p.MaxSKS
			break
		}
	}

	if totalSKS+mkInfo.SKS > maxSKS {
		return nil, apperrors.NewBadRequest(fmt.Sprintf("Total SKS prodi ini sudah mencapai batas maksimum (%d SKS)", maxSKS))
	}

	kategori := req.Kategori
	if kategori == "" {
		kategori = domain.KategoriWajib
	}

	sm := &domain.SemesterMataKuliah{
		SemesterID:   semesterID,
		MataKuliahID: req.MataKuliahID,
		Kategori:     kategori,
	}

	if err := s.repo.AssignMataKuliah(ctx, sm); err != nil {
		return nil, apperrors.NewInternal("Gagal menambahkan mata kuliah ke semester")
	}

	return sm, nil
}

func (s *semesterService) UnassignMataKuliah(ctx context.Context, semesterID string, mkID string) error {
	return s.repo.UnassignMataKuliah(ctx, semesterID, mkID)
}

func (s *semesterService) SetSKSProdi(ctx context.Context, semesterID string, req domain.SetSemesterSKSProdiRequest) ([]*domain.SemesterSKSProdi, error) {
	_, err := s.repo.GetByID(ctx, semesterID)
	if err != nil {
		return nil, apperrors.NewNotFound("Semester tidak ditemukan")
	}

	var sksProdis []*domain.SemesterSKSProdi
	for _, cfg := range req.Configs {
		if cfg.MinSKS >= cfg.MaxSKS {
			return nil, apperrors.NewBadRequest("Minimum SKS harus lebih kecil dari Maksimum SKS")
		}
		sksProdis = append(sksProdis, &domain.SemesterSKSProdi{
			ID:             uuid.New().String(),
			SemesterID:     semesterID,
			ProgramStudiID: cfg.ProgramStudiID,
			MinSKS:         cfg.MinSKS,
			MaxSKS:         cfg.MaxSKS,
		})
	}

	if err := s.repo.SetSKSProdi(ctx, semesterID, sksProdis); err != nil {
		return nil, apperrors.NewInternal("Gagal menyimpan SKS Prodi")
	}

	return s.repo.GetSKSProdi(ctx, semesterID)
}

func (s *semesterService) GetSKSProdi(ctx context.Context, semesterID string) ([]*domain.SemesterSKSProdi, error) {
	return s.repo.GetSKSProdi(ctx, semesterID)
}
