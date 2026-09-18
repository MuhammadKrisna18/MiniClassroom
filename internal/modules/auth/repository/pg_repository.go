package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"siakad-pro/internal/modules/auth/domain"
	psDomain "siakad-pro/internal/modules/programstudi/domain"
	"siakad-pro/pkg/utils"
)

type pgAuthRepository struct {
	db *gorm.DB
}

func NewPgAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &pgAuthRepository{db: db}
}

func (r *pgAuthRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&u)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &u, nil
}

func (r *pgAuthRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	result := r.db.WithContext(ctx).Preload("ProgramStudi").Where("id = ?", id).First(&u)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &u, nil
}

func (r *pgAuthRepository) GetUsersByRole(ctx context.Context, role string) ([]*domain.User, error) {
	var users []*domain.User
	result := r.db.WithContext(ctx).Preload("ProgramStudi").Where("role = ?", role).Order("created_at desc").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *pgAuthRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *pgAuthRepository) DeleteUser(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", id).Delete(&domain.EmailChangeRequest{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&domain.User{}).Error
	})
}

func (r *pgAuthRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *pgAuthRepository) CreateEmailChangeRequest(ctx context.Context, req *domain.EmailChangeRequest) error {
	return r.db.WithContext(ctx).Create(req).Error
}

func (r *pgAuthRepository) GetPendingEmailRequestByUserID(ctx context.Context, userID string) (*domain.EmailChangeRequest, error) {
	var req domain.EmailChangeRequest
	err := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userID, "pending").First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *pgAuthRepository) GetAllPendingEmailRequests(ctx context.Context) ([]*domain.EmailChangeRequest, error) {
	var reqs []*domain.EmailChangeRequest
	err := r.db.WithContext(ctx).Preload("User").Where("status = ?", "pending").Find(&reqs).Error
	return reqs, err
}

func (r *pgAuthRepository) GetEmailChangeRequestByID(ctx context.Context, id string) (*domain.EmailChangeRequest, error) {
	var req domain.EmailChangeRequest
	err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *pgAuthRepository) UpdateEmailChangeRequest(ctx context.Context, req *domain.EmailChangeRequest) error {
	return r.db.WithContext(ctx).Save(req).Error
}

func (r *pgAuthRepository) Seed(ctx context.Context, prodis []*psDomain.ProgramStudi) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	for _, prodi := range prodis {
		prodiCodeLower := strings.ToLower(prodi.Code)

		// 1. Dosen
		dosenEmail := fmt.Sprintf("dosen.%s@dosengo.id", prodiCodeLower)
		_, err := r.GetByEmail(ctx, dosenEmail)
		if err != nil && err.Error() == "user not found" {
			nidStr := utils.GenerateRandomNumberString(5)
			newDosen := &domain.User{
				ID:             uuid.New().String(),
				Name:           fmt.Sprintf("Dosen %s", prodi.Code),
				Email:          dosenEmail,
				NID:            &nidStr,
				Password:       string(hashedPassword),
				Role:           domain.RoleDosen,
				ProgramStudiID: &prodi.ID,
			}
			if err := r.Create(ctx, newDosen); err != nil {
				fmt.Printf("Error creating Dosen %s: %v\n", prodi.Code, err)
			}
		}

		// 2. Mahasiswa
		mhsEmail := fmt.Sprintf("mhs.%s@student.its.golang", prodiCodeLower)
		_, err = r.GetByEmail(ctx, mhsEmail)
		if err != nil && err.Error() == "user not found" {
			nrpStr := fmt.Sprintf("50252%s", utils.GenerateRandomNumberString(5))
			newMhs := &domain.User{
				ID:             uuid.New().String(),
				Name:           fmt.Sprintf("Mahasiswa %s", prodi.Code),
				Email:          mhsEmail,
				NRP:            &nrpStr,
				Password:       string(hashedPassword),
				Role:           domain.RoleMahasiswa,
				ProgramStudiID: &prodi.ID,
			}
			if err := r.Create(ctx, newMhs); err != nil {
				fmt.Printf("Error creating Mahasiswa %s: %v\n", prodi.Code, err)
			}
		}
	}

	return nil
}
