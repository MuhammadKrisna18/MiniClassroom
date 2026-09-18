package http

import (
	"github.com/gofiber/fiber/v2"
	"siakad-pro/internal/middleware"
)

func RegisterRoutes(router fiber.Router, handler *SemesterHandler, jwtSecret string) {
	sem := router.Group("/semester")

	sem.Use(middleware.Protected(jwtSecret))
	sem.Get("/", handler.GetAll)

	sem.Get("/:id", handler.GetByID)

	sem.Post("/", middleware.RequireRole("admin"), handler.Create)
	sem.Put("/:id", middleware.RequireRole("admin"), handler.Update)
	sem.Delete("/:id", middleware.RequireRole("admin"), handler.Delete)
	sem.Post("/:id/matakuliah", middleware.RequireRole("admin"), handler.AssignMataKuliah)
	sem.Delete("/:id/matakuliah/:mkId", middleware.RequireRole("admin"), handler.UnassignMataKuliah)
	sem.Put("/:id/sks-prodi", middleware.RequireRole("admin"), handler.SetSKSProdi)

}
