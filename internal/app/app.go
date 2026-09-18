package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"siakad-pro/config"
	_ "siakad-pro/docs"
	"siakad-pro/internal/middleware"
	"siakad-pro/internal/modules/auth"
	"siakad-pro/internal/modules/kelas"
	"siakad-pro/internal/modules/matakuliah"
	"siakad-pro/internal/modules/periode"
	"siakad-pro/internal/modules/programstudi"
	"siakad-pro/internal/modules/semester"
	"siakad-pro/internal/shared/apperrors"
	"siakad-pro/internal/shared/cache"
	"siakad-pro/internal/shared/database"
	"siakad-pro/internal/shared/response"
)

type App struct {
	cfg   *config.Config
	db    *gorm.DB
	redis *redis.Client
	fiber *fiber.App
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Start() error {
	var err error

	a.db, err = database.NewPostgresConnection(a.cfg)
	if err != nil {
		return fmt.Errorf("failed to init postgres: %w", err)
	}

	a.redis, err = cache.NewRedisClient(a.cfg)
	if err != nil {
		return fmt.Errorf("failed to init redis: %w", err)
	}

	a.fiber = fiber.New(fiber.Config{
		AppName:      "SIAKAD Pro API v1.0",
		ErrorHandler: globalErrorHandler,
	})

	a.fiber.Use(recover.New())
	a.fiber.Use(cors.New())

	a.fiber.Use(middleware.PerformanceLogger())

	a.fiber.Static("/uploads", "./uploads")

	a.fiber.Get("/", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, "System is healthy", fiber.Map{
			"env":   a.cfg.Env,
			"redis": "connected",
			"db":    "connected",
		})
	})

	a.fiber.Get("/swagger/*", swagger.HandlerDefault)

	api := a.fiber.Group("/api/v1")

	psModule := programstudi.NewProgramStudiModule(a.db, a.cfg)
	psModule.RegisterRoutes(api)

	authModule := auth.NewAuthModule(a.db, a.cfg)
	authModule.RegisterRoutes(api)

	mkModule := matakuliah.NewMataKuliahModule(a.db, a.cfg)
	mkModule.RegisterRoutes(api)

	kelasModule := kelas.NewKelasModule(a.db, a.cfg)
	kelasModule.RegisterRoutes(api)

	periodeModule := periode.NewPeriodeModule(a.db, a.cfg)
	periodeModule.RegisterRoutes(api)

	semesterModule := semester.NewSemesterModule(a.db, a.cfg)
	semesterModule.RegisterRoutes(api)

	serverErrors := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf(":%d", a.cfg.Port)
		log.Printf("Server is running on %s\n", addr)
		serverErrors <- a.fiber.Listen(addr)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Printf("Start shutdown, received signal: %v\n", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.fiber.ShutdownWithContext(ctx); err != nil {
			log.Printf("Could not gracefully shutdown server: %v\n", err)
		}

		if a.db != nil {
			log.Println("Closing PostgreSQL connection pool...")
			sqlDB, err := a.db.DB()
			if err == nil && sqlDB != nil {
				sqlDB.Close()
			}
		}

		if a.redis != nil {
			log.Println("Closing Redis client...")
			_ = a.redis.Close()
		}

		log.Println("Server gracefully stopped")
	}

	return nil
}

func globalErrorHandler(c *fiber.Ctx, err error) error {
	if e, ok := err.(*apperrors.AppError); ok {
		return response.Error(c, e.Code, e.Message, e.Details)
	}

	if e, ok := err.(*fiber.Error); ok {
		return response.Error(c, e.Code, e.Message, nil)
	}

	log.Printf("[GLOBAL ERROR] %v\n", err)
	return response.Error(c, fiber.StatusInternalServerError, "An unexpected error occurred", err.Error())
}
