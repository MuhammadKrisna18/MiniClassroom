package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *pgMataKuliahRepository) GetByID(ctx context.Context, id string) (*domain.MataKuliah, error) {
	var mk domain.MataKuliah
	err := r.db.WithContext(ctx).Preload("ProgramStudi").Where("id = ?", id).First(&mk).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &mk, nil
}

func (r *pgMataKuliahRepository) GetByIDs(ctx context.Context, ids []string) ([]*domain.MataKuliah, error) {
	var list []*domain.MataKuliah
	if len(ids) == 0 {
		return list, nil
	}
	err := r.db.WithContext(ctx).Preload("ProgramStudi").Where("id IN ?", ids).Find(&list).Error
	return list, err
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

func (r *pgMataKuliahRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.MataKuliah{}).Error
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

func (r *pgMataKuliahRepository) LockPengajuanByID(ctx context.Context, id string) (*domain.PengajuanMataKuliah, error) {
	var p domain.PengajuanMataKuliah
	err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("MataKuliah").
		Preload("Dosen").
		Where("id = ?", id).First(&p).Error
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

func (r *pgMataKuliahRepository) DeleteOtherOffersByMataKuliahID(ctx context.Context, mataKuliahID string, exceptPengajuanID string) error {
	return r.db.WithContext(ctx).
		Where("mata_kuliah_id = ? AND id != ? AND status = ?", mataKuliahID, exceptPengajuanID, domain.StatusOffered).
		Delete(&domain.PengajuanMataKuliah{}).Error
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

func (r *pgMataKuliahRepository) Transaction(ctx context.Context, fn func(txRepo domain.MataKuliahRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := NewPgMataKuliahRepository(tx)
		return fn(txRepo)
	})
}

