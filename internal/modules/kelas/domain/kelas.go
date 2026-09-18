package domain

import (
	"context"
	"time"

	authDomain "siakad-pro/internal/modules/auth/domain"
	mkDomain "siakad-pro/internal/modules/matakuliah/domain"
	psDomain "siakad-pro/internal/modules/programstudi/domain"
)

type PesertaKelas struct {
	ID          string           `json:"id" gorm:"primaryKey;type:varchar(255)"`
	PengajuanID string           `json:"pengajuan_id" gorm:"type:varchar(255);uniqueIndex:idx_peserta_kelas;not null"`
	Pengajuan   *PengajuanKelas  `json:"pengajuan,omitempty" gorm:"foreignKey:PengajuanID"`
	MahasiswaID string           `json:"mahasiswa_id" gorm:"type:varchar(255);uniqueIndex:idx_peserta_kelas;not null"`
	Mahasiswa   *authDomain.User `json:"mahasiswa,omitempty" gorm:"foreignKey:MahasiswaID"`
	Status      string           `json:"status" gorm:"type:varchar(50);not null;default:'enrolled'"`
	CreatedAt   time.Time        `json:"created_at" gorm:"autoCreateTime"`
}

const (
	MinCapacity = 25
	MaxCapacity = 50
)

type Kelas struct {
	ID             string                 `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name           string                 `json:"name" gorm:"type:varchar(255);not null"`
	Capacity       int                    `json:"capacity" gorm:"not null"`
	Hari           string                 `json:"hari" gorm:"type:varchar(20)"`
	JamMulai       string                 `json:"jam_mulai" gorm:"type:varchar(10)"`
	JamSelesai     string                 `json:"jam_selesai" gorm:"type:varchar(10)"`
	ProgramStudiID string                 `json:"program_studi_id" gorm:"type:varchar(255);not null"`
	ProgramStudi   *psDomain.ProgramStudi `json:"program_studi,omitempty" gorm:"foreignKey:ProgramStudiID"`
	Pengajuan      []*PengajuanKelas      `json:"pengajuan,omitempty" gorm:"foreignKey:KelasID"`
	CreatedAt      time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
)

type CreateKelasRequest struct {
	Name           string `json:"name" validate:"required"`
	Capacity       int    `json:"capacity" validate:"required,min=25,max=50"`
	Hari           string `json:"hari" validate:"required"`
	JamMulai       string `json:"jam_mulai" validate:"required"`
	JamSelesai     string `json:"jam_selesai" validate:"required"`
	ProgramStudiID string `json:"program_studi_id" validate:"required"`
}

type KRSRequest struct {
	PengajuanID string `json:"pengajuan_id" validate:"required"`
}

type PengajuanKelas struct {
	ID           string               `json:"id" gorm:"primaryKey;type:varchar(255)"`
	PeriodeID    string               `json:"periode_id" gorm:"type:varchar(255)"`
	DosenID      string               `json:"dosen_id" gorm:"type:varchar(255);not null"`
	Dosen        *authDomain.User     `json:"dosen,omitempty" gorm:"foreignKey:DosenID"`
	KelasID      string               `json:"kelas_id" gorm:"type:varchar(255);not null"`
	Kelas        *Kelas               `json:"kelas,omitempty" gorm:"foreignKey:KelasID"`
	MataKuliahID string               `json:"mata_kuliah_id" gorm:"type:varchar(255);not null"`
	MataKuliah   *mkDomain.MataKuliah `json:"mata_kuliah,omitempty" gorm:"foreignKey:MataKuliahID"`
	Status       string               `json:"status" gorm:"type:varchar(50);not null;default:'pending'"`
	Code         string               `json:"code" gorm:"type:varchar(6);not null"`
	CreatedAt    time.Time            `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time            `json:"updated_at" gorm:"autoUpdateTime"`
}

type RequestKelasPayload struct {
	KelasID      string `json:"kelas_id" validate:"required"`
	MataKuliahID string `json:"mata_kuliah_id" validate:"required"`
}

// KelasRepository defines all data access methods for the kelas module.
type KelasRepository interface {
	Create(ctx context.Context, kelas *Kelas) error
	GetAll(ctx context.Context) ([]*Kelas, error)
	GetByID(ctx context.Context, id string) (*Kelas, error)
	GetByName(ctx context.Context, name string) (*Kelas, error)
	CheckScheduleConflict(ctx context.Context, name string, hari string, jamMulai string) (bool, error)
	Delete(ctx context.Context, id string) error
	HasApprovedMataKuliah(ctx context.Context, dosenID string) (bool, error)
	IsMataKuliahValidForKelas(ctx context.Context, dosenID string, mkID string, prodiID string) (bool, error)
	GetActivePeriodeID(ctx context.Context) (string, error)

	CreatePengajuan(ctx context.Context, p *PengajuanKelas) error
	GetPengajuanByID(ctx context.Context, id string) (*PengajuanKelas, error)
	GetPengajuanByDosenID(ctx context.Context, dosenID string) ([]*PengajuanKelas, error)
	GetActivePengajuanByKelasID(ctx context.Context, kelasID string) ([]*PengajuanKelas, error)
	GetAllPengajuan(ctx context.Context) ([]*PengajuanKelas, error)
	UpdatePengajuan(ctx context.Context, p *PengajuanKelas) error
	DeletePengajuan(ctx context.Context, id string) error
	LockPengajuanByID(ctx context.Context, id string) (*PengajuanKelas, error)

	Transaction(ctx context.Context, fn func(txRepo KelasRepository) error) error

	GetMahasiswaByProgramStudiID(ctx context.Context, prodiID string) ([]*authDomain.User, error)
	GetApprovedPengajuanByProdiID(ctx context.Context, prodiID string) ([]*PengajuanKelas, error)
	GetUserByID(ctx context.Context, userID string) (*authDomain.User, error)

	CreatePesertaKelas(ctx context.Context, p *PesertaKelas) error
	GetPesertaKelasByPengajuanID(ctx context.Context, pengajuanID string) ([]*PesertaKelas, error)
	GetPesertaKelasByMahasiswaID(ctx context.Context, mahasiswaID string) ([]*PesertaKelas, error)
	CountPesertaKelas(ctx context.Context, pengajuanID string) (int64, error)
	CheckPesertaMataKuliahConflict(ctx context.Context, mahasiswaID string, mkID string) (bool, error)
	CheckPesertaScheduleConflict(ctx context.Context, mahasiswaID string, hari string, jamMulai string, jamSelesai string) (bool, error)

	CreatePertemuan(ctx context.Context, p *Pertemuan) error
	GetPertemuanByID(ctx context.Context, id string) (*Pertemuan, error)
	GetPertemuanByPengajuanID(ctx context.Context, pengajuanID string) ([]*Pertemuan, error)
	UpdatePertemuan(ctx context.Context, p *Pertemuan) error

	CreateAbsensi(ctx context.Context, a *Absensi) error
	GetAbsensiByPertemuanID(ctx context.Context, pertemuanID string) ([]*Absensi, error)
	UpdateAbsensiBulk(ctx context.Context, pertemuanID string, data []AbsensiUpdate) error
}

// KelasService defines all business logic methods for the kelas module.
type KelasService interface {
	Create(ctx context.Context, req CreateKelasRequest) (*Kelas, error)
	GetAll(ctx context.Context) ([]*Kelas, error)
	GetByID(ctx context.Context, id string) (*Kelas, error)
	Delete(ctx context.Context, id string) error

	RequestKelas(ctx context.Context, dosenID string, req RequestKelasPayload) (*PengajuanKelas, error)
	ApprovePengajuan(ctx context.Context, id string) error
	RejectPengajuan(ctx context.Context, id string) error
	GetMyPengajuan(ctx context.Context, dosenID string) ([]*PengajuanKelas, error)
	GetAllPengajuan(ctx context.Context) ([]*PengajuanKelas, error)
	GetMahasiswaInKelas(ctx context.Context, pengajuanID string, dosenID string) ([]*authDomain.User, error)
	GetMyJadwal(ctx context.Context, userID string) ([]*PengajuanKelas, error)
	GetAvailableKelas(ctx context.Context, userID string) ([]*PengajuanKelas, error)
	AmbilKelas(ctx context.Context, userID string, pengajuanID string) error

	MulaiPertemuan(ctx context.Context, pengajuanID string, judul string) (*Pertemuan, error)
	AkhiriPertemuan(ctx context.Context, pertemuanID string) error
	GetPertemuanByPengajuan(ctx context.Context, pengajuanID string) ([]*Pertemuan, error)
	GetAbsensi(ctx context.Context, pertemuanID string) ([]*Absensi, error)
	SubmitAbsensi(ctx context.Context, pertemuanID string, data BulkAbsensiRequest) error
	SubmitAbsensiMahasiswa(ctx context.Context, pertemuanID string, mahasiswaID string, kode string) error
	GetRekapKehadiran(ctx context.Context, pengajuanID string, dosenID string) (*RekapKehadiranResponse, error)
	GetRekapKehadiranAdmin(ctx context.Context, pengajuanID string, dosenID string) (*AdminRekapResponse, error)
}
