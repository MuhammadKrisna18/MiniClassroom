package http

import (
	"github.com/gofiber/fiber/v2"
	"siakad-pro/internal/modules/semester/domain"
	"siakad-pro/internal/shared/apperrors"
	"siakad-pro/internal/shared/response"
)

type SemesterHandler struct {
	service domain.SemesterService
}

func NewSemesterHandler(svc domain.SemesterService) *SemesterHandler {
	return &SemesterHandler{service: svc}
}

func (h *SemesterHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateSemesterRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Invalid request body", err.Error())
	}

	sem, err := h.service.Create(c.UserContext(), req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Semester berhasil dibuat", sem)
}

func (h *SemesterHandler) GetAll(c *fiber.Ctx) error {
	semesters, err := h.service.GetAll(c.UserContext())
	if err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Daftar semester", semesters)
}

func (h *SemesterHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID semester diperlukan")
	}

	sem, err := h.service.GetByID(c.UserContext(), id)
	if err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Detail semester", sem)
}

func (h *SemesterHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID semester diperlukan")
	}

	var req domain.UpdateSemesterRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Invalid request body", err.Error())
	}

	sem, err := h.service.Update(c.UserContext(), id, req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Semester berhasil diupdate", sem)
}

func (h *SemesterHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID semester diperlukan")
	}

	if err := h.service.Delete(c.UserContext(), id); err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Semester berhasil dihapus", nil)
}

func (h *SemesterHandler) AssignMataKuliah(c *fiber.Ctx) error {
	semesterID := c.Params("id")
	if semesterID == "" {
		return apperrors.NewBadRequest("ID semester diperlukan")
	}

	var req domain.AssignMataKuliahRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Invalid request body", err.Error())
	}

	sm, err := h.service.AssignMataKuliah(c.UserContext(), semesterID, req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Mata kuliah berhasil ditambahkan ke semester", sm)
}

func (h *SemesterHandler) UnassignMataKuliah(c *fiber.Ctx) error {
	semesterID := c.Params("id")
	mkID := c.Params("mkId")
	if semesterID == "" || mkID == "" {
		return apperrors.NewBadRequest("ID semester dan ID mata kuliah diperlukan")
	}

	if err := h.service.UnassignMataKuliah(c.UserContext(), semesterID, mkID); err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Mata kuliah berhasil dihapus dari semester", nil)
}

func (h *SemesterHandler) SetSKSProdi(c *fiber.Ctx) error {
	semesterID := c.Params("id")
	if semesterID == "" {
		return apperrors.NewBadRequest("ID semester diperlukan")
	}

	var req domain.SetSemesterSKSProdiRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Invalid request body", err.Error())
	}

	sksProdis, err := h.service.SetSKSProdi(c.UserContext(), semesterID, req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "SKS Prodi berhasil diupdate", sksProdis)
}

