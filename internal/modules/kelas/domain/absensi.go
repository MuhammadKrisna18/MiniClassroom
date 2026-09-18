package domain

import (
	"time"

	authDomain "siakad-pro/internal/modules/auth/domain"
)

const (
	AbsensiHadir = "hadir"
	AbsensiIzin  = "izin"
	AbsensiSakit = "sakit"
	AbsensiAlpa  = "alpa"
)

type Absensi struct {
	ID              string           `json:"id" gorm:"primaryKey;type:varchar(255)"`
	PertemuanID     string           `json:"pertemuan_id" gorm:"type:varchar(255);not null"`
	Pertemuan       *Pertemuan       `json:"pertemuan,omitempty" gorm:"foreignKey:PertemuanID"`
	MahasiswaID     string           `json:"mahasiswa_id" gorm:"type:varchar(255);not null"`
	Mahasiswa       *authDomain.User `json:"mahasiswa,omitempty" gorm:"foreignKey:MahasiswaID"`
	StatusKehadiran string           `json:"status_kehadiran" gorm:"type:varchar(20);not null;default:'alpa'"`
	CreatedAt       time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

type BulkAbsensiRequest struct {
	Data []AbsensiUpdate `json:"data" validate:"required,min=1"`
}

type AbsensiUpdate struct {
	MahasiswaID     string `json:"mahasiswa_id" validate:"required"`
	StatusKehadiran string `json:"status_kehadiran" validate:"required,oneof=hadir izin sakit alpa"`
}

type RekapKehadiranResponse struct {
	Pertemuan []PertemuanInfo  `json:"pertemuan"`
	Mahasiswa []MahasiswaRekap `json:"mahasiswa"`
}

type PertemuanInfo struct {
	ID      string    `json:"id"`
	Judul   string    `json:"judul"`
	Tanggal time.Time `json:"tanggal"`
}

type MahasiswaRekap struct {
	ID        string            `json:"id"`
	NRP       string            `json:"nrp"`
	Name      string            `json:"name"`
	Kehadiran map[string]string `json:"kehadiran"`
}

type AdminRekapResponse struct {
	*RekapKehadiranResponse
	TotalPertemuan int                      `json:"total_pertemuan"`
	Summary        []map[string]interface{} `json:"summary"`
}
