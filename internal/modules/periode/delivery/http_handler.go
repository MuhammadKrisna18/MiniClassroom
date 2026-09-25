package delivery

import (
	"github.com/gofiber/fiber/v2"
	"siakad-pro/internal/middleware"
	"siakad-pro/internal/modules/periode/domain"
	"siakad-pro/internal/shared/apperrors"
	"siakad-pro/internal/shared/response"
)

type PeriodeHandler struct {
	service domain.PeriodeService
}

func NewPeriodeHandler(service domain.PeriodeService) *PeriodeHandler {
	return &PeriodeHandler{service: service}
}

func (h *PeriodeHandler) RegisterRoutes(router fiber.Router, jwtSecret string) {
	periodeRoute := router.Group("/periode")

	// Public / All Roles (untuk melihat periode aktif)
	periodeRoute.Get("/active", middleware.Protected(jwtSecret), h.GetActive)

	// Admin Only
	adminRoute := periodeRoute.Group("", middleware.Protected(jwtSecret), middleware.RequireRole("admin"))
	adminRoute.Post("/", h.Create)
	adminRoute.Get("/", h.GetAll)
	adminRoute.Delete("/:id", h.Delete)
	adminRoute.Post("/:id/active", h.SetActive)
}

func (h *PeriodeHandler) Create(c *fiber.Ctx) error {
	var req domain.CreatePeriodeRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Invalid request body", err.Error())
	}

	p, err := h.service.Create(c.UserContext(), req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Periode akademik berhasil dibuat", p)
}

func (h *PeriodeHandler) GetAll(c *fiber.Ctx) error {
	periodes, err := h.service.GetAll(c.UserContext())
	if err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Daftar periode akademik", periodes)
}

func (h *PeriodeHandler) GetActive(c *fiber.Ctx) error {
	p, err := h.service.GetActive(c.UserContext())
	if err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Periode akademik aktif saat ini", p)
}

func (h *PeriodeHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID periode diperlukan")
	}

	if err := h.service.Delete(c.UserContext(), id); err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Periode berhasil dihapus", nil)
}

func (h *PeriodeHandler) SetActive(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID periode diperlukan")
	}

	if err := h.service.SetActive(c.UserContext(), id); err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Periode berhasil diaktifkan", nil)
}

