package domain

import (
	"context"
	"time"

	psDomain "siakad-pro/internal/modules/programstudi/domain"
)

const (
	RoleAdmin     = "admin"
	RoleDosen     = "dosen"
	RoleMahasiswa = "mahasiswa"
)

type User struct {
	ID             string                 `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name           string                 `json:"name" gorm:"type:varchar(255);not null"`
	Nickname       *string                `json:"nickname" gorm:"type:varchar(255)"`
	NID            *string                `json:"nid" gorm:"type:varchar(5);unique"`
	NRP            *string                `json:"nrp" gorm:"type:varchar(14);unique"`
	Email          string                 `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password       string                 `json:"-" gorm:"type:varchar(255);not null"`
	Role           string                 `json:"role" gorm:"type:varchar(50);not null;default:'user'"`
	ProgramStudiID *string                `json:"program_studi_id" gorm:"type:varchar(255)"`
	ProgramStudi   *psDomain.ProgramStudi `json:"program_studi,omitempty" gorm:"foreignKey:ProgramStudiID"`
	PhotoURL       *string                `json:"photo_url" gorm:"type:varchar(255)"`
	CreatedAt      time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

type EmailChangeRequest struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	UserID    string    `json:"user_id" gorm:"type:varchar(255);not null"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	NewEmail  string    `json:"new_email" gorm:"type:varchar(255);not null"`
	Status    string    `json:"status" gorm:"type:varchar(50);not null;default:'pending'"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterDosenRequest struct {
	Name           string `json:"name" validate:"required"`
	Username       string `json:"username" validate:"required"`
	Password       string `json:"password" validate:"required"`
	ProgramStudiID string `json:"program_studi_id" validate:"required"`
}

type UpdateProfileRequest struct {
	Name           string `json:"name" validate:"required"`
	Nickname       string `json:"nickname"`
	ProgramStudiID string `json:"program_studi_id"`
}

type EmailChangeRequestPayload struct {
	NewEmail string `json:"new_email" validate:"required,email"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type UserProfileResponse struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Nickname       *string                `json:"nickname,omitempty"`
	NID            *string                `json:"nid,omitempty"`
	NRP            *string                `json:"nrp,omitempty"`
	Email          string                 `json:"email"`
	Role           string                 `json:"role"`
	ProgramStudiID *string                `json:"program_studi_id,omitempty"`
	ProgramStudi   *psDomain.ProgramStudi `json:"program_studi,omitempty"`
	PendingEmail   *string                `json:"pending_email,omitempty"`
	PhotoURL       *string                `json:"photo_url,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
}

type AuthRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetUsersByRole(ctx context.Context, role string) ([]*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error

	CreateEmailChangeRequest(ctx context.Context, req *EmailChangeRequest) error
	GetPendingEmailRequestByUserID(ctx context.Context, userID string) (*EmailChangeRequest, error)
	GetAllPendingEmailRequests(ctx context.Context) ([]*EmailChangeRequest, error)
	GetEmailChangeRequestByID(ctx context.Context, id string) (*EmailChangeRequest, error)
	UpdateEmailChangeRequest(ctx context.Context, req *EmailChangeRequest) error
	
	Seed(ctx context.Context, prodis []*psDomain.ProgramStudi) error
}

type AuthService interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	GetProfile(ctx context.Context, id string) (*UserProfileResponse, error)
	RegisterDosen(ctx context.Context, req RegisterDosenRequest) (*UserProfileResponse, error)
	GetDosenList(ctx context.Context) ([]*UserProfileResponse, error)
	DeleteDosen(ctx context.Context, id string) error
	UpdateProfile(ctx context.Context, id string, req UpdateProfileRequest) (*UserProfileResponse, error)
	UpdateProfilePhoto(ctx context.Context, id string, photoURL string) (*UserProfileResponse, error)

	RequestEmailChange(ctx context.Context, userID string, req EmailChangeRequestPayload) error
	GetPendingEmailRequests(ctx context.Context) ([]*EmailChangeRequest, error)
	ReviewEmailRequest(ctx context.Context, requestID string, approve bool) error

	RegisterMahasiswa(ctx context.Context, req RegisterMahasiswaRequest) (*UserProfileResponse, error)
	GetMahasiswaList(ctx context.Context) ([]*UserProfileResponse, error)
}

type RegisterMahasiswaRequest struct {
	Name           string `json:"name" validate:"required"`
	NRP            string `json:"nrp" validate:"required"`
	Password       string `json:"password" validate:"required"`
	ProgramStudiID string `json:"program_studi_id" validate:"required"`
}
