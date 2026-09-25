package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"siakad-pro/internal/modules/programstudi/domain"
)

type pgProgramStudiRepository struct {
	db *gorm.DB
}

func NewPgProgramStudiRepository(db *gorm.DB) domain.ProgramStudiRepository {
	return &pgProgramStudiRepository{db: db}
}

func (r *pgProgramStudiRepository) Create(ctx context.Context, ps *domain.ProgramStudi) error {
	return r.db.WithContext(ctx).Create(ps).Error
}

func (r *pgProgramStudiRepository) GetByName(ctx context.Context, name string) (*domain.ProgramStudi, error) {
	var ps domain.ProgramStudi
	err := r.db.WithContext(ctx).Where("LOWER(name) = LOWER(?)", name).First(&ps).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ps, nil
}

func (r *pgProgramStudiRepository) GetAll(ctx context.Context) ([]*domain.ProgramStudi, error) {
	var psList []*domain.ProgramStudi
	err := r.db.WithContext(ctx).Order("name asc").Find(&psList).Error
	if err != nil {
		return nil, err
	}
	return psList, nil
}

func (r *pgProgramStudiRepository) Seed(ctx context.Context) error {
	prodis := []domain.ProgramStudi{
		{Name: "Teknik Informatika", Code: "TI"},
		{Name: "Rekayasa Perangkat Lunak", Code: "RPL"},
		{Name: "Rekayasa Kecerdasan Artificial", Code: "RKA"},
	}

	for _, p := range prodis {
		var existing domain.ProgramStudi
		err := r.db.WithContext(ctx).Where("code = ?", p.Code).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				p.ID = uuid.New().String()
				if err := r.db.WithContext(ctx).Create(&p).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}
