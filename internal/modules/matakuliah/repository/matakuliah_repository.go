package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"siakad-pro/internal/modules/matakuliah/domain"
)

type pgMataKuliahRepository struct {
	db *gorm.DB
}

func NewPgMataKuliahRepository(db *gorm.DB) domain.MataKuliahRepository {
	return &pgMataKuliahRepository{db: db}
}

func (r *pgMataKuliahRepository) Create(ctx context.Context, mk *domain.MataKuliah) error {
	return r.db.WithContext(ctx).Create(mk).Error
}

func (r *pgMataKuliahRepository) GetByNameAndProdi(ctx context.Context, name string, prodiID string) (*domain.MataKuliah, error) {
	var mk domain.MataKuliah
	err := r.db.WithContext(ctx).Where("LOWER(name) = LOWER(?) AND program_studi_id = ?", name, prodiID).First(&mk).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &mk, nil
}

func (r *pgMataKuliahRepository) GetAll(ctx context.Context) ([]*domain.MataKuliah, error) {
	var mkList []*domain.MataKuliah
	err := r.db.WithContext(ctx).Preload("ProgramStudi").Preload("Pengajuan").Preload("Pengajuan.Dosen").Order("created_at desc").Find(&mkList).Error
	if err != nil {
		return nil, err
	}
	return mkList, nil
}

func (r *pgMataKuliahRepository) GetByProdi(ctx context.Context, prodiID string) ([]*domain.MataKuliah, error) {
	var mkList []*domain.MataKuliah
	err := r.db.WithContext(ctx).Preload("ProgramStudi").Preload("Pengajuan").Preload("Pengajuan.Dosen").Where("program_studi_id = ?", prodiID).Order("created_at desc").Find(&mkList).Error
	if err != nil {
		return nil, err
	}
	return mkList, nil
}

func (r *pgMataKuliahRepository) GetUserProdiID(ctx context.Context, userID string) (*string, error) {
	var prodiID *string
	err := r.db.WithContext(ctx).Table("users").Select("program_studi_id").Where("id = ?", userID).Scan(&prodiID).Error
	return prodiID, err
}

func (r *pgMataKuliahRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.MataKuliah{}).Error
}

func (r *pgMataKuliahRepository) GetActivePeriodeID(ctx context.Context) (string, error) {
	var id string
	err := r.db.WithContext(ctx).Table("periode_akademiks").Select("id").Where("is_active = ?", true).Scan(&id).Error
	if err != nil {
		return "", err
	}
	if id == "" {
		return "", errors.New("tidak ada periode aktif")
	}
	return id, nil
}

func (r *pgMataKuliahRepository) CreatePengajuan(ctx context.Context, p *domain.PengajuanMataKuliah) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *pgMataKuliahRepository) GetPengajuanByID(ctx context.Context, id string) (*domain.PengajuanMataKuliah, error) {
	var p domain.PengajuanMataKuliah
	err := r.db.WithContext(ctx).Preload("MataKuliah").Preload("Dosen").Where("id = ?", id).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgMataKuliahRepository) GetPengajuanByDosenID(ctx context.Context, dosenID string) ([]*domain.PengajuanMataKuliah, error) {
	var list []*domain.PengajuanMataKuliah
	err := r.db.WithContext(ctx).Preload("MataKuliah").Preload("Dosen").Where("dosen_id = ?", dosenID).Order("created_at desc").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *pgMataKuliahRepository) GetActivePengajuanByMataKuliahID(ctx context.Context, mkID string) ([]*domain.PengajuanMataKuliah, error) {
	var list []*domain.PengajuanMataKuliah
	err := r.db.WithContext(ctx).Preload("Dosen").Where("mata_kuliah_id = ? AND status IN ('pending', 'approved')", mkID).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *pgMataKuliahRepository) GetAllPengajuan(ctx context.Context) ([]*domain.PengajuanMataKuliah, error) {
	var list []*domain.PengajuanMataKuliah
	err := r.db.WithContext(ctx).Preload("MataKuliah").Preload("Dosen").Order("created_at desc").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *pgMataKuliahRepository) UpdatePengajuan(ctx context.Context, p *domain.PengajuanMataKuliah) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *pgMataKuliahRepository) DeletePengajuan(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.PengajuanMataKuliah{}).Error
}

func (r *pgMataKuliahRepository) GetDosenIDsByProdi(ctx context.Context, prodiID string) ([]string, error) {
	var ids []string
	err := r.db.WithContext(ctx).Table("users").
		Where("role = ? AND program_studi_id = ?", "dosen", prodiID).
		Pluck("id", &ids).Error
	return ids, err
}

func (r *pgMataKuliahRepository) IsMataKuliahValidForKelas(ctx context.Context, dosenID string, mkID string, prodiID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("pengajuan_mata_kuliahs").
		Joins("JOIN mata_kuliahs mk ON mk.id = pengajuan_mata_kuliahs.mata_kuliah_id").
		Where("pengajuan_mata_kuliahs.dosen_id = ? AND pengajuan_mata_kuliahs.mata_kuliah_id = ? AND pengajuan_mata_kuliahs.status = ? AND mk.program_studi_id = ?", dosenID, mkID, domain.StatusApproved, prodiID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
