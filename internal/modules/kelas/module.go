package kelas

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"siakad-pro/config"
	authRepo "siakad-pro/internal/modules/auth/repository"
	"siakad-pro/internal/modules/kelas/delivery/http"
	"siakad-pro/internal/modules/kelas/repository"
	"siakad-pro/internal/modules/kelas/service"
	mkRepo "siakad-pro/internal/modules/matakuliah/repository"
	periodeRepo "siakad-pro/internal/modules/periode/repository"
)

type KelasModule struct {
	handler *http.KelasHandler
	cfg     *config.Config
}

func NewKelasModule(db *gorm.DB, cfg *config.Config) *KelasModule {
	repo := repository.NewPgKelasRepository(db)
	userProvider := repository.NewUserProviderAdapter(authRepo.NewPgAuthRepository(db))
	periodeProvider := repository.NewPeriodeProviderAdapter(periodeRepo.NewPostgresPeriodeRepository(db))
	mkProvider := repository.NewMataKuliahProviderAdapter(mkRepo.NewPgMataKuliahRepository(db))

	svc := service.NewKelasService(repo, userProvider, periodeProvider, mkProvider)
	handler := http.NewKelasHandler(svc)

	return &KelasModule{
		handler: handler,
		cfg:     cfg,
	}
}

func (m *KelasModule) RegisterRoutes(router fiber.Router) {
	http.RegisterRoutes(router, m.handler, m.cfg.JWTSecret)
}
