package periode

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"siakad-pro/config"
	"siakad-pro/internal/modules/periode/delivery"
	"siakad-pro/internal/modules/periode/repository"
	"siakad-pro/internal/modules/periode/service"
)

type PeriodeModule struct {
	Handler *delivery.PeriodeHandler
	cfg     *config.Config
}

func NewPeriodeModule(db *gorm.DB, cfg *config.Config) *PeriodeModule {
	repo := repository.NewPostgresPeriodeRepository(db)
	svc := service.NewPeriodeService(repo)
	handler := delivery.NewPeriodeHandler(svc)

	return &PeriodeModule{
		Handler: handler,
		cfg:     cfg,
	}
}

func (m *PeriodeModule) RegisterRoutes(router fiber.Router) {
	m.Handler.RegisterRoutes(router, m.cfg.JWTSecret)
}
