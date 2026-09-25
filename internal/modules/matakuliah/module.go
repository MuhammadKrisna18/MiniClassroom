package matakuliah

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"siakad-pro/config"
	authRepo "siakad-pro/internal/modules/auth/repository"
	"siakad-pro/internal/modules/matakuliah/delivery/http"
	"siakad-pro/internal/modules/matakuliah/repository"
	"siakad-pro/internal/modules/matakuliah/service"
	periodeRepo "siakad-pro/internal/modules/periode/repository"
)

type Module struct {
	Handler *http.MataKuliahHandler
	cfg     *config.Config
}

func NewMataKuliahModule(db *gorm.DB, cfg *config.Config) *Module {
	repo := repository.NewPgMataKuliahRepository(db)
	userProvider := repository.NewUserProviderAdapter(authRepo.NewPgAuthRepository(db))
	periodeProvider := repository.NewPeriodeProviderAdapter(periodeRepo.NewPostgresPeriodeRepository(db))
	svc := service.NewMataKuliahService(repo, userProvider, periodeProvider)
	handler := http.NewMataKuliahHandler(svc)

	return &Module{
		Handler: handler,
		cfg:     cfg,
	}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	m.Handler.RegisterRoutes(router, m.cfg.JWTSecret)
}
