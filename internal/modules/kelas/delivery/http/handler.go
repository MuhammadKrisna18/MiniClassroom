package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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

	kelas, err := h.service.Create(c.Context(), req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Kelas berhasil ditambahkan", kelas)
}

func (h *KelasHandler) GetAll(c *fiber.Ctx) error {
	kelases, err := h.service.GetAll(c.Context())
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

	kelas, err := h.service.GetByID(c.Context(), id)
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

	if err := h.service.Delete(c.Context(), id); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Kelas berhasil dihapus", nil)
}

func (h *KelasHandler) RequestKelas(c *fiber.Ctx) error {
	dosenID := c.Locals("userID").(string)

	var req domain.RequestKelasPayload
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Format payload tidak valid")
	}

	if err := h.validate.Struct(req); err != nil {
		return apperrors.NewBadRequest("Validasi gagal: kelas_id diperlukan")
	}

	pengajuan, err := h.service.RequestKelas(c.Context(), dosenID, req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Berhasil mengajukan kelas", pengajuan)
}

func (h *KelasHandler) GetMyPengajuan(c *fiber.Ctx) error {
	dosenID := c.Locals("userID").(string)

	list, err := h.service.GetMyPengajuan(c.Context(), dosenID)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil riwayat pengajuan kelas", list)
}

func (h *KelasHandler) GetAllPengajuan(c *fiber.Ctx) error {
	list, err := h.service.GetAllPengajuan(c.Context())
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

	if err := h.service.ApprovePengajuan(c.Context(), id); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil menyetujui kelas", nil)
}

func (h *KelasHandler) RejectPengajuan(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID pengajuan diperlukan")
	}

	if err := h.service.RejectPengajuan(c.Context(), id); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil menolak kelas", nil)
}

func (h *KelasHandler) GetMahasiswaInKelas(c *fiber.Ctx) error {
	dosenID, ok := c.Locals("userID").(string)
	if !ok {
		return apperrors.NewUnauthorized("Unauthorized")
	}

	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID pengajuan diperlukan")
	}

	list, err := h.service.GetMahasiswaInKelas(c.Context(), id, dosenID)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil data mahasiswa di kelas ini", list)
}

func (h *KelasHandler) GetMyJadwal(c *fiber.Ctx) error {
	mahasiswaID, ok := c.Locals("userID").(string)
	if !ok {
		return apperrors.NewUnauthorized("Unauthorized")
	}

	list, err := h.service.GetMyJadwal(c.Context(), mahasiswaID)
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

	p, err := h.service.MulaiPertemuan(c.Context(), payload.PengajuanID, payload.Judul)
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
	if err := h.service.AkhiriPertemuan(c.Context(), id); err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Berhasil mengakhiri pertemuan", nil)
}

func (h *KelasHandler) GetPertemuanByPengajuan(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID Pengajuan diperlukan")
	}
	list, err := h.service.GetPertemuanByPengajuan(c.Context(), id)
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
	list, err := h.service.GetAbsensi(c.Context(), id)
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

	if err := h.service.SubmitAbsensi(c.Context(), id, req); err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Berhasil mengupdate absensi", nil)
}

func (h *KelasHandler) SubmitAbsensiMahasiswa(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID Pertemuan diperlukan")
	}

	mahasiswaID, ok := c.Locals("userID").(string)
	if !ok {
		return apperrors.NewUnauthorized("Unauthorized")
	}

	var payload struct {
		Kode string `json:"kode" validate:"required"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return apperrors.NewBadRequest("Format request tidak valid")
	}

	if err := h.service.SubmitAbsensiMahasiswa(c.Context(), id, mahasiswaID, payload.Kode); err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Kehadiran Berhasil Dicatat", nil)
}

func (h *KelasHandler) GetAvailableKelas(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	jadwal, err := h.service.GetAvailableKelas(c.Context(), userID)
	if err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "Data bursa kelas", jadwal)
}

func (h *KelasHandler) AmbilKelas(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	
	var req domain.KRSRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Format payload tidak valid")
	}

	if err := h.validate.Struct(req); err != nil {
		return apperrors.NewBadRequest("Validasi gagal")
	}

	if err := h.service.AmbilKelas(c.Context(), userID, req.PengajuanID); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil kelas", nil)
}

func (h *KelasHandler) GetRekapKehadiran(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID Pengajuan diperlukan")
	}

	dosenID, ok := c.Locals("userID").(string)
	if !ok {
		return apperrors.NewUnauthorized("Unauthorized")
	}

	rekap, err := h.service.GetRekapKehadiran(c.Context(), id, dosenID)
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

	rekap, err := h.service.GetRekapKehadiran(c.Context(), id, "")
	if err != nil {
		return err
	}

	type AdminRekapResponse struct {
		*domain.RekapKehadiranResponse
		TotalPertemuan int                      `json:"total_pertemuan"`
		Summary        []map[string]interface{} `json:"summary"`
	}

	totalPertemuan := len(rekap.Pertemuan)
	var summary []map[string]interface{}

	for _, m := range rekap.Mahasiswa {
		hadirCount := 0
		for _, status := range m.Kehadiran {
			if status == domain.AbsensiHadir {
				hadirCount++
			}
		}

		summary = append(summary, map[string]interface{}{
			"id":              m.ID,
			"nrp":             m.NRP,
			"name":            m.Name,
			"total_hadir":     hadirCount,
			"total_pertemuan": totalPertemuan,
		})
	}

	res := AdminRekapResponse{
		RekapKehadiranResponse: rekap,
		TotalPertemuan:         totalPertemuan,
		Summary:                summary,
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil rekap kehadiran", res)
}
