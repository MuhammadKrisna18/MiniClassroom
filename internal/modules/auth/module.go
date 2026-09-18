package auth

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"siakad-pro/config"
	"siakad-pro/internal/modules/auth/delivery/http"
	"siakad-pro/internal/modules/auth/repository"
	"siakad-pro/internal/modules/auth/service"
	psDomain "siakad-pro/internal/modules/programstudi/domain"
)

type Module struct {
	Handler *http.AuthHandler
	cfg     *config.Config
}

func NewAuthModule(db *gorm.DB, cfg *config.Config) *Module {
	repo := repository.NewPgAuthRepository(db)
	svc := service.NewAuthService(repo, cfg)
	handler := http.NewAuthHandler(svc)

	// Fetch program studi via gorm to avoid circular dependency with programstudi module
	var prodis []*psDomain.ProgramStudi
	_ = db.Find(&prodis).Error
	
	if len(prodis) > 0 {
		_ = repo.Seed(context.Background(), prodis)
	}

	return &Module{
		Handler: handler,
		cfg:     cfg,
	}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	m.Handler.RegisterRoutes(router, m.cfg.JWTSecret)
}
