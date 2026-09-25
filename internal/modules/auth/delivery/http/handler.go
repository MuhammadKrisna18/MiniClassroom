package http

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"siakad-pro/internal/middleware"
	"siakad-pro/internal/modules/auth/domain"
	"siakad-pro/internal/shared/apperrors"
	"siakad-pro/internal/shared/response"
)

type AuthHandler struct {
	service domain.AuthService
}

func NewAuthHandler(s domain.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) RegisterRoutes(router fiber.Router, jwtSecret string) {
	authGroup := router.Group("/auth")
	authGroup.Post("/login", h.Login)
	authGroup.Get("/me", middleware.Protected(jwtSecret), h.Me)
	authGroup.Put("/me", middleware.Protected(jwtSecret), h.UpdateProfile)
	authGroup.Post("/profile/photo", middleware.Protected(jwtSecret), h.UploadPhoto)

	authGroup.Post("/email-request", middleware.Protected(jwtSecret), h.RequestEmailChange)
	authGroup.Get("/email-request", middleware.Protected(jwtSecret), middleware.RequireRole("admin"), h.GetEmailRequests)
	authGroup.Put("/email-request/:id", middleware.Protected(jwtSecret), middleware.RequireRole("admin"), h.ReviewEmailRequest)

	authGroup.Post("/dosen", middleware.Protected(jwtSecret), middleware.RequireRole("admin"), h.RegisterDosen)
	authGroup.Get("/dosen", middleware.Protected(jwtSecret), middleware.RequireRole("admin"), h.GetDosenList)
	authGroup.Delete("/dosen/:id", middleware.Protected(jwtSecret), middleware.RequireRole("admin"), h.DeleteDosen)

	authGroup.Post("/mahasiswa", middleware.Protected(jwtSecret), middleware.RequireRole("admin"), h.RegisterMahasiswa)
	authGroup.Get("/mahasiswa", middleware.Protected(jwtSecret), middleware.RequireRole("admin"), h.GetMahasiswaList)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("invalid request body", err.Error())
	}

	if req.Email == "" || req.Password == "" {
		return apperrors.NewBadRequest("email and password are required")
	}

	res, err := h.service.Login(c.UserContext(), req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "login successful", res)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	profile, err := h.service.GetProfile(c.UserContext(), userID)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "success", profile)
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	var req domain.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("invalid request body", err.Error())
	}

	if req.Name == "" {
		return apperrors.NewBadRequest("Name is required")
	}

	res, err := h.service.UpdateProfile(c.UserContext(), userID, req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Profile updated successfully", res)
}

func (h *AuthHandler) RequestEmailChange(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	var req domain.EmailChangeRequestPayload
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("invalid request body", err.Error())
	}

	if req.NewEmail == "" {
		return apperrors.NewBadRequest("New email is required")
	}

	if err := h.service.RequestEmailChange(c.UserContext(), userID, req); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Permintaan ganti email berhasil diajukan", nil)
}

func (h *AuthHandler) GetEmailRequests(c *fiber.Ctx) error {
	list, err := h.service.GetPendingEmailRequests(c.UserContext())
	if err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "success", list)
}

func (h *AuthHandler) ReviewEmailRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID tidak valid")
	}

	var req struct {
		Approve bool `json:"approve"`
	}
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("invalid request body", err.Error())
	}

	if err := h.service.ReviewEmailRequest(c.UserContext(), id, req.Approve); err != nil {
		return err
	}

	statusStr := "ditolak"
	if req.Approve {
		statusStr = "disetujui"
	}
	return response.Success(c, fiber.StatusOK, "Permintaan ganti email "+statusStr, nil)
}

func (h *AuthHandler) RegisterDosen(c *fiber.Ctx) error {
	var req domain.RegisterDosenRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("invalid request body", err.Error())
	}

	if req.Name == "" || req.Username == "" || req.Password == "" || req.ProgramStudiID == "" {
		return apperrors.NewBadRequest("name, username, password, and program studi are required")
	}

	res, err := h.service.RegisterDosen(c.UserContext(), req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "dosen account created successfully", res)
}

func (h *AuthHandler) GetDosenList(c *fiber.Ctx) error {
	list, err := h.service.GetDosenList(c.UserContext())
	if err != nil {
		return err
	}
	return response.Success(c, fiber.StatusOK, "success", list)
}

func (h *AuthHandler) DeleteDosen(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperrors.NewBadRequest("ID tidak valid")
	}

	if err := h.service.DeleteDosen(c.UserContext(), id); err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil menghapus akun dosen", nil)
}

func (h *AuthHandler) UploadPhoto(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	file, err := c.FormFile("photo")
	if err != nil {
		return apperrors.NewBadRequest("Foto tidak ditemukan", err.Error())
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return apperrors.NewBadRequest("Format file tidak didukung. Gunakan JPG atau PNG.")
	}

	if file.Size > 2*1024*1024 {
		return apperrors.NewBadRequest("Ukuran file maksimal 2MB")
	}

	uploadDir := "./uploads/profiles"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return apperrors.NewInternal("Gagal membuat direktori upload", err.Error())
	}

	fileName := fmt.Sprintf("%s_%d%s", userID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, fileName)

	if err := c.SaveFile(file, filePath); err != nil {
		return apperrors.NewInternal("Gagal menyimpan file", err.Error())
	}

	photoURL := fmt.Sprintf("/uploads/profiles/%s", fileName)
	res, err := h.service.UpdateProfilePhoto(c.UserContext(), userID, photoURL)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Foto profil berhasil diperbarui", res)
}

func (h *AuthHandler) RegisterMahasiswa(c *fiber.Ctx) error {
	var req domain.RegisterMahasiswaRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequest("Format data tidak valid", err.Error())
	}

	if req.Name == "" || req.NRP == "" || req.Password == "" || req.ProgramStudiID == "" {
		return apperrors.NewBadRequest("Semua field harus diisi")
	}

	res, err := h.service.RegisterMahasiswa(c.UserContext(), req)
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusCreated, "Berhasil mendaftarkan mahasiswa", res)
}

func (h *AuthHandler) GetMahasiswaList(c *fiber.Ctx) error {
	res, err := h.service.GetMahasiswaList(c.UserContext())
	if err != nil {
		return err
	}

	return response.Success(c, fiber.StatusOK, "Berhasil mengambil data mahasiswa", res)
}
