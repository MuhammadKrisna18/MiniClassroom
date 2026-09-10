package domain

import (
	"context"
	"time"
)

type PeriodeAkademik struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Tahun     string    `json:"tahun" gorm:"type:varchar(20);not null"` // e.g. "2024/2025"
	Jenis     string    `json:"jenis" gorm:"type:varchar(10);not null"` // "ganjil", "genap", "pendek"
	IsActive  bool      `json:"is_active" gorm:"not null;default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreatePeriodeRequest struct {
	Tahun string `json:"tahun" validate:"required"`
	Jenis string `json:"jenis" validate:"required,oneof=ganjil genap pendek"`
}

type PeriodeRepository interface {
	Create(ctx context.Context, p *PeriodeAkademik) error
	GetAll(ctx context.Context) ([]*PeriodeAkademik, error)
	GetByID(ctx context.Context, id string) (*PeriodeAkademik, error)
	GetActive(ctx context.Context) (*PeriodeAkademik, error)
	Delete(ctx context.Context, id string) error
	SetActive(ctx context.Context, id string) error
	DeactivateAll(ctx context.Context) error
}

type PeriodeService interface {
	Create(ctx context.Context, req CreatePeriodeRequest) (*PeriodeAkademik, error)
	GetAll(ctx context.Context) ([]*PeriodeAkademik, error)
	GetActive(ctx context.Context) (*PeriodeAkademik, error)
	Delete(ctx context.Context, id string) error
	SetActive(ctx context.Context, id string) error
}
