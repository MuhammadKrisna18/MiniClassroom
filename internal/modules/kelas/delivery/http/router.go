package http

import (
	"github.com/gofiber/fiber/v2"
	"siakad-pro/internal/middleware"
)

func RegisterRoutes(router fiber.Router, handler *KelasHandler, jwtSecret string) {
	kelas := router.Group("/kelas")

	kelas.Use(middleware.Protected(jwtSecret))
	kelas.Get("/", handler.GetAll)

	kelas.Post("/request", middleware.RequireRole("dosen"), handler.RequestKelas)
	kelas.Get("/pengajuan/me", middleware.RequireRole("dosen"), handler.GetMyPengajuan)
	kelas.Get("/pengajuan/:id/mahasiswa", middleware.RequireRole("dosen"), handler.GetMahasiswaInKelas)

	kelas.Get("/mahasiswa/my-jadwal", middleware.RequireRole("mahasiswa"), handler.GetMyJadwal)
	kelas.Get("/mahasiswa/krs/available", middleware.RequireRole("mahasiswa"), handler.GetAvailableKelas)
	kelas.Post("/mahasiswa/krs", middleware.RequireRole("mahasiswa"), handler.AmbilKelas)

	kelas.Get("/pengajuan", middleware.RequireRole("admin"), handler.GetAllPengajuan)
	kelas.Post("/pengajuan/:id/approve", middleware.RequireRole("admin"), handler.ApprovePengajuan)
	kelas.Post("/pengajuan/:id/reject", middleware.RequireRole("admin"), handler.RejectPengajuan)

	kelas.Post("/", middleware.RequireRole("admin"), handler.Create)
	kelas.Delete("/:id", middleware.RequireRole("admin"), handler.Delete)
	kelas.Get("/:id", handler.GetByID)

	pertemuan := router.Group("/pertemuan")
	pertemuan.Use(middleware.Protected(jwtSecret))
	pertemuan.Post("/", middleware.RequireRole("dosen"), handler.MulaiPertemuan)
	pertemuan.Put("/:id/end", middleware.RequireRole("dosen"), handler.AkhiriPertemuan)
	pertemuan.Get("/pengajuan/:id", handler.GetPertemuanByPengajuan)
	pertemuan.Get("/:id/absensi", handler.GetAbsensi)
	pertemuan.Put("/:id/absensi", middleware.RequireRole("dosen"), handler.SubmitAbsensi)
	pertemuan.Post("/:id/attend", middleware.RequireRole("mahasiswa"), handler.SubmitAbsensiMahasiswa)
	pertemuan.Get("/rekap/:id", middleware.RequireRole("dosen"), handler.GetRekapKehadiran)
	pertemuan.Get("/admin/rekap/:id", middleware.RequireRole("admin"), handler.GetRekapKehadiranAdmin)
}
