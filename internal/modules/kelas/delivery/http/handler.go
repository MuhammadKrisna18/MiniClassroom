package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"siakad-pro/internal/middleware"
	"siakad-pro/internal/modules/kelas/domain"
	"siakad-pro/internal/shared/apperrors"
	"siakad-pro/internal/shared/response"
)

type KelasHandler struct {
	service  domain.KelasService
	validate *validator.Validate
}

func NewKelasHandler(service domain.KelasService) *KelasHandler {
	return &KelasHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *KelasHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateKelasRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Format payload tidak valid")
	}

	if err := h.validate.Struct(req); err != nil {
		return apperrors.NewBadRequest("Validasi gagal: periksa kembali data yang dimasukkan")
	}

	kelas, err := h.service.Create(c.UserContext(), req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Kelas berhasil ditambahkan", kelas)
}

func (h *KelasHandler) GetAll(c *fiber.Ctx) error {
	kelases, err := h.service.GetAll(c.UserContext())
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil data kelas", kelases)
}

func (h *KelasHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID kelas diperlukan")
	}

	kelas, err := h.service.GetByID(c.UserContext(), id)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil detail kelas", kelas)
}

func (h *KelasHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID kelas diperlukan")
	}

	if err := h.service.Delete(c.UserContext(), id); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Kelas berhasil dihapus", nil)
}

func (h *KelasHandler) RequestKelas(c *fiber.Ctx) error {
	dosenID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	var req domain.RequestKelasPayload
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Format payload tidak valid")
	}

	if err := h.validate.Struct(req); err != nil {
		return apperrors.NewBadRequest("Validasi gagal: kelas_id diperlukan")
	}

	pengajuan, err := h.service.RequestKelas(c.UserContext(), dosenID, req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Berhasil mengajukan kelas", pengajuan)
}

func (h *KelasHandler) GetMyPengajuan(c *fiber.Ctx) error {
	dosenID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	list, err := h.service.GetMyPengajuan(c.UserContext(), dosenID)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil riwayat pengajuan kelas", list)
}

func (h *KelasHandler) GetAllPengajuan(c *fiber.Ctx) error {
	list, err := h.service.GetAllPengajuan(c.UserContext())
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil semua pengajuan kelas", list)
}

func (h *KelasHandler) ApprovePengajuan(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID pengajuan diperlukan")
	}

	if err := h.service.ApprovePengajuan(c.UserContext(), id); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil menyetujui kelas", nil)
}

func (h *KelasHandler) RejectPengajuan(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID pengajuan diperlukan")
	}

	if err := h.service.RejectPengajuan(c.UserContext(), id); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil menolak kelas", nil)
}

func (h *KelasHandler) GetMahasiswaInKelas(c *fiber.Ctx) error {
	dosenID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID pengajuan diperlukan")
	}

	list, err := h.service.GetMahasiswaInKelas(c.UserContext(), id, dosenID)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil data mahasiswa di kelas ini", list)
}

func (h *KelasHandler) GetMyJadwal(c *fiber.Ctx) error {
	mahasiswaID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	list, err := h.service.GetMyJadwal(c.UserContext(), mahasiswaID)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil jadwal kelas", list)
}

func (h *KelasHandler) MulaiPertemuan(c *fiber.Ctx) error {
	var payload struct {
		PengajuanID string `json:"pengajuan_id" validate:"required"`
		Judul       string `json:"judul" validate:"required"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return apperrors.NewBadRequest("Format request tidak valid")
	}

	p, err := h.service.MulaiPertemuan(c.UserContext(), payload.PengajuanID, payload.Judul)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Berhasil memulai pertemuan", p)
}

func (h *KelasHandler) AkhiriPertemuan(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID Pertemuan diperlukan")
	}
	if err := h.service.AkhiriPertemuan(c.UserContext(), id); err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Berhasil mengakhiri pertemuan", nil)
}

func (h *KelasHandler) GetPertemuanByPengajuan(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID Pengajuan diperlukan")
	}
	list, err := h.service.GetPertemuanByPengajuan(c.UserContext(), id)
	if err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Berhasil mengambil daftar pertemuan", list)
}

func (h *KelasHandler) GetAbsensi(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID Pertemuan diperlukan")
	}
	list, err := h.service.GetAbsensi(c.UserContext(), id)
	if err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Berhasil mengambil absensi", list)
}

func (h *KelasHandler) SubmitAbsensi(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID Pertemuan diperlukan")
	}

	var req domain.BulkAbsensiRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Format request tidak valid")
	}

	if err := h.service.SubmitAbsensi(c.UserContext(), id, req); err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Berhasil mengupdate absensi", nil)
}

func (h *KelasHandler) SubmitAbsensiMahasiswa(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID Pertemuan diperlukan")
	}

	mahasiswaID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	var payload struct {
		Kode string `json:"kode" validate:"required"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return apperrors.NewBadRequest("Format request tidak valid")
	}

	if err := h.service.SubmitAbsensiMahasiswa(c.UserContext(), id, mahasiswaID, payload.Kode); err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Kehadiran Berhasil Dicatat", nil)
}

func (h *KelasHandler) GetAvailableKelas(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}
	jadwal, err := h.service.GetAvailableKelas(c.UserContext(), userID)
	if err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Data bursa kelas", jadwal)
}

func (h *KelasHandler) AmbilKelas(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}
	
	var req domain.KRSRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Format payload tidak valid")
	}

	if err := h.validate.Struct(req); err != nil {
		return apperrors.NewBadRequest("Validasi gagal")
	}

	if err := h.service.AmbilKelas(c.UserContext(), userID, req.PengajuanID); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil kelas", nil)
}

func (h *KelasHandler) GetRekapKehadiran(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID Pengajuan diperlukan")
	}

	dosenID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	rekap, err := h.service.GetRekapKehadiran(c.UserContext(), id, dosenID)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil rekap kehadiran", rekap)
}

func (h *KelasHandler) GetRekapKehadiranAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID Pengajuan diperlukan")
	}

	res, err := h.service.GetRekapKehadiranAdmin(c.UserContext(), id, "")
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil rekap kehadiran", res)
}

