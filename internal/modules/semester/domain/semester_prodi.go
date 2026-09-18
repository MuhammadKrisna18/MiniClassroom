package domain

import (
	"time"

	psDomain "siakad-pro/internal/modules/programstudi/domain"
)

type SemesterSKSProdi struct {
	ID             string                 `json:"id" gorm:"primaryKey;type:varchar(255)"`
	SemesterID     string                 `json:"semester_id" gorm:"type:varchar(255);uniqueIndex:idx_sem_prodi;not null"`
	ProgramStudiID string                 `json:"program_studi_id" gorm:"type:varchar(255);uniqueIndex:idx_sem_prodi;not null"`
	ProgramStudi   *psDomain.ProgramStudi `json:"program_studi,omitempty" gorm:"foreignKey:ProgramStudiID"`
	MinSKS         int                    `json:"min_sks" gorm:"not null"`
	MaxSKS         int                    `json:"max_sks" gorm:"not null"`
	CreatedAt      time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

type SKSProdiConfig struct {
	ProgramStudiID string `json:"program_studi_id" validate:"required"`
	MinSKS         int    `json:"min_sks" validate:"required,min=1"`
	MaxSKS         int    `json:"max_sks" validate:"required,min=1"`
}

type SetSemesterSKSProdiRequest struct {
	Configs []SKSProdiConfig `json:"configs" validate:"required,dive"`
}
