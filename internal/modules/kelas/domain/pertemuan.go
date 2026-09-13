package domain

import (
	"time"
)

const (
	PertemuanStatusBerlangsung = "berlangsung"
	PertemuanStatusSelesai     = "selesai"

	MaxPertemuan = 16
)

type Pertemuan struct {
	ID             string          `json:"id" gorm:"primaryKey;type:varchar(255)"`
	PengajuanID    string          `json:"pengajuan_id" gorm:"type:varchar(255);not null"`
	Pengajuan      *PengajuanKelas `json:"pengajuan,omitempty" gorm:"foreignKey:PengajuanID"`
	Judul          string          `json:"judul" gorm:"type:varchar(255);not null"`
	Tanggal        time.Time       `json:"tanggal" gorm:"not null"`
	WaktuMulai     time.Time       `json:"waktu_mulai" gorm:"not null"`
	WaktuSelesai   *time.Time      `json:"waktu_selesai"`
	Status         string          `json:"status" gorm:"type:varchar(50);not null;default:'berlangsung'"`
	KodeAbsensi    string          `json:"kode_absensi" gorm:"type:varchar(10)"`
	NomorPertemuan int             `json:"nomor_pertemuan" gorm:"not null;default:1"`
	Absensi        []*Absensi      `json:"absensi,omitempty" gorm:"foreignKey:PertemuanID"`
	CreatedAt      time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}
