package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"siakad-pro/internal/modules/periode/domain"
)

type postgresPeriodeRepository struct {
	db *gorm.DB
}

func NewPostgresPeriodeRepository(db *gorm.DB) domain.PeriodeRepository {
	return &postgresPeriodeRepository{db: db}
}

func (r *postgresPeriodeRepository) Create(ctx context.Context, p *domain.PeriodeAkademik) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *postgresPeriodeRepository) GetAll(ctx context.Context) ([]*domain.PeriodeAkademik, error) {
	var periodes []*domain.PeriodeAkademik
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&periodes).Error
	return periodes, err
}

func (r *postgresPeriodeRepository) GetByID(ctx context.Context, id string) (*domain.PeriodeAkademik, error) {
	var p domain.PeriodeAkademik
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("periode not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *postgresPeriodeRepository) GetActive(ctx context.Context) (*domain.PeriodeAkademik, error) {
	var p domain.PeriodeAkademik
	err := r.db.WithContext(ctx).Where("is_active = ?", true).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No active period is okay, handle upstream
		}
		return nil, err
	}
	return &p, nil
}

func (r *postgresPeriodeRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.PeriodeAkademik{}, "id = ?", id).Error
}

func (r *postgresPeriodeRepository) SetActive(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&domain.PeriodeAkademik{}).Where("id = ?", id).Update("is_active", true).Error
}

func (r *postgresPeriodeRepository) DeactivateAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Model(&domain.PeriodeAkademik{}).Where("1=1").Update("is_active", false).Error
}
