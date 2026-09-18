package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"siakad-pro/internal/modules/semester/domain"
)

type pgSemesterRepository struct {
	db *gorm.DB
}

func NewPgSemesterRepository(db *gorm.DB) domain.SemesterRepository {
	return &pgSemesterRepository{db: db}
}

func (r *pgSemesterRepository) Create(ctx context.Context, s *domain.Semester) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *pgSemesterRepository) GetAll(ctx context.Context) ([]*domain.Semester, error) {
	var semesters []*domain.Semester
	err := r.db.WithContext(ctx).Preload("MataKuliah").Preload("MataKuliah.MataKuliah").Preload("MataKuliah.MataKuliah.ProgramStudi").Preload("SKSProdi").Preload("SKSProdi.ProgramStudi").Order("nomor ASC").Find(&semesters).Error
	return semesters, err
}

func (r *pgSemesterRepository) GetByID(ctx context.Context, id string) (*domain.Semester, error) {
	var s domain.Semester
	err := r.db.WithContext(ctx).Preload("MataKuliah").Preload("MataKuliah.MataKuliah").Preload("MataKuliah.MataKuliah.ProgramStudi").Preload("SKSProdi").Preload("SKSProdi.ProgramStudi").First(&s, "id = ?", id).Error
	return &s, err
}

func (r *pgSemesterRepository) GetByNomor(ctx context.Context, nomor int) (*domain.Semester, error) {
	var s domain.Semester
	err := r.db.WithContext(ctx).First(&s, "nomor = ?", nomor).Error
	return &s, err
}

func (r *pgSemesterRepository) Update(ctx context.Context, s *domain.Semester) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *pgSemesterRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("semester_id = ?", id).Delete(&domain.SemesterMataKuliah{}).Error; err != nil {
			return err
		}

		return tx.Delete(&domain.Semester{}, "id = ?", id).Error
	})
}

func (r *pgSemesterRepository) AssignMataKuliah(ctx context.Context, sm *domain.SemesterMataKuliah) error {
	if sm.ID == "" {
		sm.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(sm).Error
}

func (r *pgSemesterRepository) UnassignMataKuliah(ctx context.Context, semesterID string, mkID string) error {
	return r.db.WithContext(ctx).Where("semester_id = ? AND mata_kuliah_id = ?", semesterID, mkID).Delete(&domain.SemesterMataKuliah{}).Error
}

func (r *pgSemesterRepository) GetSemesterMataKuliah(ctx context.Context, semesterID string) ([]*domain.SemesterMataKuliah, error) {
	var items []*domain.SemesterMataKuliah
	err := r.db.WithContext(ctx).Preload("MataKuliah").Preload("MataKuliah.ProgramStudi").Where("semester_id = ?", semesterID).Find(&items).Error
	return items, err
}

func (r *pgSemesterRepository) GetTotalSKS(ctx context.Context, semesterID string, prodiID string) (int, error) {
	var total int
	row := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(mk.sks), 0)
		FROM semester_mata_kuliahs smk
		JOIN mata_kuliahs mk ON mk.id = smk.mata_kuliah_id
		WHERE smk.semester_id = ? AND mk.program_studi_id = ?
	`, semesterID, prodiID).Row()
	err := row.Scan(&total)
	return total, err
}

func (r *pgSemesterRepository) GetMataKuliahProdiID(ctx context.Context, mkID string) (string, error) {
	var prodiID string
	err := r.db.WithContext(ctx).Raw("SELECT program_studi_id FROM mata_kuliahs WHERE id = ?", mkID).Row().Scan(&prodiID)
	return prodiID, err
}

func (r *pgSemesterRepository) SetSKSProdi(ctx context.Context, semesterID string, sksProdis []*domain.SemesterSKSProdi) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing entries for this semester
		if err := tx.Where("semester_id = ?", semesterID).Delete(&domain.SemesterSKSProdi{}).Error; err != nil {
			return err
		}
		
		if len(sksProdis) > 0 {
			if err := tx.Create(&sksProdis).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
}

func (r *pgSemesterRepository) GetSKSProdi(ctx context.Context, semesterID string) ([]*domain.SemesterSKSProdi, error) {
	var items []*domain.SemesterSKSProdi
	err := r.db.WithContext(ctx).Preload("ProgramStudi").Where("semester_id = ?", semesterID).Find(&items).Error
	return items, err
}


