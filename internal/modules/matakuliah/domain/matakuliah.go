package domain

import (
	"context"
	"time"

	authDomain "siakad-pro/internal/modules/auth/domain"
	psDomain "siakad-pro/internal/modules/programstudi/domain"
)

type MataKuliah struct {
	ID             string                 `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name           string                 `json:"name" gorm:"type:varchar(255);uniqueIndex:idx_name_prodi;not null"`
	SKS            int                    `json:"sks" gorm:"not null"`
	ProgramStudiID string                 `json:"program_studi_id" gorm:"type:varchar(255);uniqueIndex:idx_name_prodi;not null"`
	ProgramStudi   *psDomain.ProgramStudi `json:"program_studi,omitempty" gorm:"foreignKey:ProgramStudiID"`
	Pengajuan      []*PengajuanMataKuliah `json:"pengajuan,omitempty" gorm:"foreignKey:MataKuliahID"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

const (
	StatusPending  = "pending"
	StatusOffered  = "offered"
	StatusApproved = "approved"
	StatusRejected = "rejected"
)

type CreateMataKuliahRequest struct {
	Name           string `json:"name" validate:"required"`
	SKS            int    `json:"sks" validate:"required,min=1"`
	ProgramStudiID string `json:"program_studi_id" validate:"required"`
}

// UserProvider defines the port for fetching user information needed by matakuliah
type UserProvider interface {
	GetUserProdiID(ctx context.Context, userID string) (*string, error)
	GetDosenIDsByProdi(ctx context.Context, prodiID string) ([]string, error)
}

// PeriodeProvider defines the port for fetching academic period information needed by matakuliah
type PeriodeProvider interface {
	GetActivePeriodeID(ctx context.Context) (string, error)
}

// MataKuliahCatalogRepository handles mata kuliah catalog persistence
type MataKuliahCatalogRepository interface {
	Create(ctx context.Context, mk *MataKuliah) error
	GetByNameAndProdi(ctx context.Context, name string, prodiID string) (*MataKuliah, error)
	GetByID(ctx context.Context, id string) (*MataKuliah, error)
	GetByIDs(ctx context.Context, ids []string) ([]*MataKuliah, error)
	GetAll(ctx context.Context) ([]*MataKuliah, error)
	GetByProdi(ctx context.Context, prodiID string) ([]*MataKuliah, error)
	Delete(ctx context.Context, id string) error
}

// MataKuliahPenawaranRepository handles lecturer applications and course offering workflows
type MataKuliahPenawaranRepository interface {
	CreatePengajuan(ctx context.Context, p *PengajuanMataKuliah) error
	GetPengajuanByID(ctx context.Context, id string) (*PengajuanMataKuliah, error)
	LockPengajuanByID(ctx context.Context, id string) (*PengajuanMataKuliah, error)
	GetPengajuanByDosenID(ctx context.Context, dosenID string) ([]*PengajuanMataKuliah, error)
	GetActivePengajuanByMataKuliahID(ctx context.Context, mkID string) ([]*PengajuanMataKuliah, error)
	GetAllPengajuan(ctx context.Context) ([]*PengajuanMataKuliah, error)
	IsMataKuliahValidForKelas(ctx context.Context, dosenID string, mkID string, prodiID string) (bool, error)
	UpdatePengajuan(ctx context.Context, p *PengajuanMataKuliah) error
	DeletePengajuan(ctx context.Context, id string) error
	DeleteOtherOffersByMataKuliahID(ctx context.Context, mataKuliahID string, exceptPengajuanID string) error
	Transaction(ctx context.Context, fn func(txRepo MataKuliahRepository) error) error
}

// MataKuliahRepository is the composite interface
type MataKuliahRepository interface {
	MataKuliahCatalogRepository
	MataKuliahPenawaranRepository
}

// MataKuliahCatalogService handles curriculum and catalog operations
type MataKuliahCatalogService interface {
	CreateMataKuliah(ctx context.Context, req CreateMataKuliahRequest) (*MataKuliah, error)
	GetMataKuliahList(ctx context.Context) ([]*MataKuliah, error)
	GetMataKuliahForMahasiswa(ctx context.Context, userID string) ([]*MataKuliah, error)
	DeleteMataKuliah(ctx context.Context, id string) error
	LepasMataKuliah(ctx context.Context, mkID string) error
}

// MataKuliahPenawaranService handles lecturer requests, offers, and approvals
type MataKuliahPenawaranService interface {
	RequestMataKuliah(ctx context.Context, dosenID string, req RequestMataKuliahPayload) (*PengajuanMataKuliah, error)
	ApprovePengajuan(ctx context.Context, id string) error
	RejectPengajuan(ctx context.Context, id string) error
	AcceptOffer(ctx context.Context, id string, dosenID string) error
	RejectOffer(ctx context.Context, id string, dosenID string) error
	GetMyPengajuan(ctx context.Context, dosenID string) ([]*PengajuanMataKuliah, error)
	GetAllPengajuan(ctx context.Context) ([]*PengajuanMataKuliah, error)
}

// MataKuliahService is the composite interface
type MataKuliahService interface {
	MataKuliahCatalogService
	MataKuliahPenawaranService
}

type PengajuanMataKuliah struct {
	ID           string           `json:"id" gorm:"primaryKey;type:varchar(255)"`
	PeriodeID    string           `json:"periode_id" gorm:"type:varchar(255)"`
	DosenID      string           `json:"dosen_id" gorm:"type:varchar(255);not null"`
	Dosen        *authDomain.User `json:"dosen,omitempty" gorm:"foreignKey:DosenID"`
	MataKuliahID string           `json:"mata_kuliah_id" gorm:"type:varchar(255);not null"`
	MataKuliah   *MataKuliah      `json:"mata_kuliah,omitempty" gorm:"foreignKey:MataKuliahID"`
	Status       string           `json:"status" gorm:"type:varchar(50);not null;default:'pending'"`
	Code         string           `json:"code" gorm:"type:varchar(6);not null"`
	CreatedAt    time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

type RequestMataKuliahPayload struct {
	MataKuliahID string `json:"mata_kuliah_id" validate:"required"`
}

type ApproveRejectPengajuanPayload struct {
}
