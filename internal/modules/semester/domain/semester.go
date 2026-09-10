package domain

import (
	"context"
	"time"

	mkDomain "siakad-pro/internal/modules/matakuliah/domain"
)

type Semester struct {
	ID         string                `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Nomor      int                   `json:"nomor" gorm:"uniqueIndex;not null"`
	MinSKS     int                   `json:"min_sks" gorm:"not null;default:18"`
	MaxSKS     int                   `json:"max_sks" gorm:"not null;default:24"`
	MataKuliah []*SemesterMataKuliah `json:"mata_kuliah,omitempty" gorm:"foreignKey:SemesterID"`
	CreatedAt  time.Time             `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time             `json:"updated_at" gorm:"autoUpdateTime"`
}

type SemesterMataKuliah struct {
	ID           string               `json:"id" gorm:"primaryKey;type:varchar(255)"`
	SemesterID   string               `json:"semester_id" gorm:"type:varchar(255);uniqueIndex:idx_sem_mk;not null"`
	Semester     *Semester            `json:"semester,omitempty" gorm:"foreignKey:SemesterID"`
	MataKuliahID string               `json:"mata_kuliah_id" gorm:"type:varchar(255);uniqueIndex:idx_sem_mk;not null"`
	MataKuliah   *mkDomain.MataKuliah `json:"mata_kuliah,omitempty" gorm:"foreignKey:MataKuliahID"`
	Kategori     string               `json:"kategori" gorm:"type:varchar(50);not null;default:'wajib'"`
	CreatedAt    time.Time            `json:"created_at" gorm:"autoCreateTime"`
}

const (
	KategoriWajib   = "wajib"
	KategoriPilihan = "pilihan"

	MaxPertemuan = 16
	MinSemester  = 1
	MaxSemester  = 8
)

type CreateSemesterRequest struct {
	Nomor  int `json:"nomor" validate:"required,min=1,max=8"`
	MinSKS int `json:"min_sks" validate:"required,min=1"`
	MaxSKS int `json:"max_sks" validate:"required,min=1"`
}

type UpdateSemesterRequest struct {
	MinSKS *int `json:"min_sks"`
	MaxSKS *int `json:"max_sks"`
}

type AssignMataKuliahRequest struct {
	MataKuliahID string `json:"mata_kuliah_id" validate:"required"`
	Kategori     string `json:"kategori" validate:"required"`
}


type SemesterRepository interface {
	Create(ctx context.Context, s *Semester) error
	GetAll(ctx context.Context) ([]*Semester, error)
	GetByID(ctx context.Context, id string) (*Semester, error)
	GetByNomor(ctx context.Context, nomor int) (*Semester, error)
	Update(ctx context.Context, s *Semester) error
	Delete(ctx context.Context, id string) error

	AssignMataKuliah(ctx context.Context, sm *SemesterMataKuliah) error
	UnassignMataKuliah(ctx context.Context, semesterID string, mkID string) error
	GetSemesterMataKuliah(ctx context.Context, semesterID string) ([]*SemesterMataKuliah, error)
	GetTotalSKS(ctx context.Context, semesterID string) (int, error)


}

type SemesterService interface {
	Create(ctx context.Context, req CreateSemesterRequest) (*Semester, error)
	GetAll(ctx context.Context) ([]*Semester, error)
	GetByID(ctx context.Context, id string) (*Semester, error)
	Update(ctx context.Context, id string, req UpdateSemesterRequest) (*Semester, error)
	Delete(ctx context.Context, id string) error

	AssignMataKuliah(ctx context.Context, semesterID string, req AssignMataKuliahRequest) (*SemesterMataKuliah, error)
	UnassignMataKuliah(ctx context.Context, semesterID string, mkID string) error


}
