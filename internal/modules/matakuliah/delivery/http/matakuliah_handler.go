package http

import (
	"siakad-pro/internal/middleware"
	"siakad-pro/internal/modules/matakuliah/domain"
	"siakad-pro/internal/shared/apperrors"
	"siakad-pro/internal/shared/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type MataKuliahHandler struct {
	service domain.MataKuliahService
}

func NewMataKuliahHandler(service domain.MataKuliahService) *MataKuliahHandler {
	return &MataKuliahHandler{
		service: service,
	}
}

func (h *MataKuliahHandler) RegisterRoutes(router fiber.Router, jwtSecret string) {
	mkGroup := router.Group("/matakuliah")

	mkGroup.Use(middleware.Protected(jwtSecret))

	mkGroup.Get("/", h.GetMataKuliahList)
	mkGroup.Get("/mahasiswa", middleware.RequireRole("mahasiswa"), h.GetMataKuliahForMahasiswa)
	mkGroup.Post("/", middleware.RequireRole("admin"), h.CreateMataKuliah)
	mkGroup.Delete("/:id", middleware.RequireRole("admin"), h.DeleteMataKuliah)
	mkGroup.Post("/:id/lepas", middleware.RequireRole("admin"), h.LepasMataKuliah)

	mkGroup.Post("/requests", middleware.RequireRole("dosen"), h.RequestMataKuliah)
	mkGroup.Post("/requests/:id/accept-offer", middleware.RequireRole("dosen"), h.AcceptOffer)
	mkGroup.Post("/requests/:id/reject-offer", middleware.RequireRole("dosen"), h.RejectOffer)
	mkGroup.Get("/my-requests", middleware.RequireRole("dosen"), h.GetMyPengajuan)
	mkGroup.Get("/requests", middleware.RequireRole("admin"), h.GetAllPengajuan)
	mkGroup.Post("/requests/:id/approve", middleware.RequireRole("admin"), h.ApprovePengajuan)
	mkGroup.Post("/requests/:id/reject", middleware.RequireRole("admin"), h.RejectPengajuan)
}

func (h *MataKuliahHandler) CreateMataKuliah(c *fiber.Ctx) error {
	var req domain.CreateMataKuliahRequest

	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Invalid request payload", err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return apperrors.NewBadRequest("Data tidak valid", err.Error())
	}

	mk, err := h.service.CreateMataKuliah(c.UserContext(), req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Mata kuliah berhasil ditambahkan", mk)
}

func (h *MataKuliahHandler) GetMataKuliahList(c *fiber.Ctx) error {
	mkList, err := h.service.GetMataKuliahList(c.UserContext())
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil daftar mata kuliah", mkList)
}

func (h *MataKuliahHandler) GetMataKuliahForMahasiswa(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	mkList, err := h.service.GetMataKuliahForMahasiswa(c.UserContext(), userID)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil daftar mata kuliah", mkList)
}

func (h *MataKuliahHandler) DeleteMataKuliah(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID tidak valid")
	}

	if err := h.service.DeleteMataKuliah(c.UserContext(), id); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil menghapus mata kuliah", nil)
}

func (h *MataKuliahHandler) LepasMataKuliah(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID tidak valid")
	}

	if err := h.service.LepasMataKuliah(c.UserContext(), id); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil melepas mata kuliah dari dosen", nil)
}

func (h *MataKuliahHandler) RequestMataKuliah(c *fiber.Ctx) error {
	var req domain.RequestMataKuliahPayload
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Invalid request payload", err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return apperrors.NewBadRequest("Data tidak valid", err.Error())
	}

	dosenID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	pengajuan, err := h.service.RequestMataKuliah(c.UserContext(), dosenID, req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Berhasil mengajukan mata kuliah", pengajuan)
}

func (h *MataKuliahHandler) AcceptOffer(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID tidak valid")
	}

	dosenID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	if err := h.service.AcceptOffer(c.UserContext(), id, dosenID); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil menerima penawaran mata kuliah", nil)
}

func (h *MataKuliahHandler) RejectOffer(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID tidak valid")
	}

	dosenID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	if err := h.service.RejectOffer(c.UserContext(), id, dosenID); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil menolak penawaran mata kuliah", nil)
}

func (h *MataKuliahHandler) GetMyPengajuan(c *fiber.Ctx) error {
	dosenID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	list, err := h.service.GetMyPengajuan(c.UserContext(), dosenID)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil riwayat pengajuan", list)
}

func (h *MataKuliahHandler) GetAllPengajuan(c *fiber.Ctx) error {
	list, err := h.service.GetAllPengajuan(c.UserContext())
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil semua pengajuan", list)
}

func (h *MataKuliahHandler) ApprovePengajuan(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID Pengajuan wajib diisi")
	}

	if err := h.service.ApprovePengajuan(c.UserContext(), id); err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Pengajuan berhasil disetujui", nil)
}

func (h *MataKuliahHandler) RejectPengajuan(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID Pengajuan wajib diisi")
	}

	if err := h.service.RejectPengajuan(c.UserContext(), id); err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Pengajuan berhasil ditolak", nil)
}

