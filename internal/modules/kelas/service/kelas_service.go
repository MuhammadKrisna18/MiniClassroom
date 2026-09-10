package service

import (
	"context"
	"log"
	"math/rand"
	"regexp"
	"time"

	"fmt"
	"github.com/google/uuid"
	authDomain "siakad-pro/internal/modules/auth/domain"
	"siakad-pro/internal/modules/kelas/domain"
	"siakad-pro/internal/shared/apperrors"
)

type kelasService struct {
	repo domain.KelasRepository
}

func NewKelasService(repo domain.KelasRepository) domain.KelasService {
	return &kelasService{repo: repo}
}

func (s *kelasService) Create(ctx context.Context, req domain.CreateKelasRequest) (*domain.Kelas, error) {

	matched, _ := regexp.MatchString(`^IF-[1-3]0[1-7]$`, req.Name)
	if !matched {
		return nil, apperrors.NewBadRequest("Format nama kelas tidak valid (contoh yang benar: IF-101 s/d IF-107, IF-201 s/d IF-207, IF-301 s/d IF-307)")
	}

	if req.Capacity < 25 || req.Capacity > 50 {
		return nil, apperrors.NewBadRequest("Kapasitas kelas harus antara 25 dan 50")
	}

	conflict, _ := s.repo.CheckScheduleConflict(ctx, req.Name, req.Hari, req.JamMulai)
	if conflict {
		return nil, apperrors.NewBadRequest("Kelas tersebut sudah terdaftar pada hari dan jam yang sama")
	}

	validHari := map[string]bool{"Senin": true, "Selasa": true, "Rabu": true, "Kamis": true, "Jumat": true}
	if !validHari[req.Hari] {
		return nil, apperrors.NewBadRequest("Hari harus antara Senin sampai Jumat")
	}

	if req.JamMulai == "" || req.JamSelesai == "" {
		return nil, apperrors.NewBadRequest("Jam mulai dan selesai harus diisi")
	}

	kelas := &domain.Kelas{
		ID:             uuid.New().String(),
		Name:           req.Name,
		Capacity:       req.Capacity,
		Hari:           req.Hari,
		JamMulai:       req.JamMulai,
		JamSelesai:     req.JamSelesai,
		ProgramStudiID: req.ProgramStudiID,
	}

	if err := s.repo.Create(ctx, kelas); err != nil {
		log.Println("DB INSERT ERROR:", err)
		return nil, apperrors.NewInternal("Gagal membuat kelas: " + err.Error())
	}

	return s.repo.GetByID(ctx, kelas.ID)
}

func (s *kelasService) GetAll(ctx context.Context) ([]*domain.Kelas, error) {
	return s.repo.GetAll(ctx)
}

func (s *kelasService) GetByID(ctx context.Context, id string) (*domain.Kelas, error) {
	kelas, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFound("Kelas tidak ditemukan")
	}
	return kelas, nil
}

func (s *kelasService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Kelas tidak ditemukan")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		log.Printf("Error deleting kelas: %v", err)
		return apperrors.NewInternal("Gagal menghapus kelas", err.Error())
	}
	return nil
}

func (s *kelasService) RequestKelas(ctx context.Context, dosenID string, req domain.RequestKelasPayload) (*domain.PengajuanKelas, error) {
	log.Printf("[RequestKelas] dosenID=%s, kelasID=%s, mataKuliahID=%s", dosenID, req.KelasID, req.MataKuliahID)

	kelas, err := s.repo.GetByID(ctx, req.KelasID)
	if err != nil {
		log.Printf("[RequestKelas] ERROR: kelas not found: %v", err)
		return nil, apperrors.NewNotFound("Kelas tidak ditemukan")
	}
	log.Printf("[RequestKelas] Kelas found: %s, ProdiID: %s", kelas.Name, kelas.ProgramStudiID)

	hasMK, err := s.repo.IsMataKuliahValidForKelas(ctx, dosenID, req.MataKuliahID, kelas.ProgramStudiID)
	if err != nil {
		log.Printf("[RequestKelas] ERROR: IsMataKuliahValidForKelas failed: %v", err)
		return nil, apperrors.NewInternal("Gagal memvalidasi mata kuliah dosen", err.Error())
	}
	log.Printf("[RequestKelas] IsMataKuliahValidForKelas result: %v", hasMK)
	if !hasMK {
		return nil, apperrors.NewBadRequest("Mata kuliah yang dipilih tidak valid atau bukan berasal dari program studi kelas ini.")
	}

	activePengajuan, err := s.repo.GetActivePengajuanByKelasID(ctx, req.KelasID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengecek status kelas", err.Error())
	}

	for _, p := range activePengajuan {
		if p.Status == domain.StatusApproved {
			return nil, apperrors.NewBadRequest("Kelas ini sudah disetujui untuk dosen lain")
		}
		if p.DosenID == dosenID && p.Status == domain.StatusPending {
			return nil, apperrors.NewBadRequest("Anda sudah mengajukan kelas ini")
		}
	}

	dosenPengajuan, err := s.repo.GetPengajuanByDosenID(ctx, dosenID)
	if err == nil {
		for _, dp := range dosenPengajuan {
			if dp.Status != domain.StatusRejected && dp.Kelas != nil {
				if dp.Kelas.Hari == kelas.Hari {
					if kelas.JamMulai < dp.Kelas.JamSelesai && dp.Kelas.JamMulai < kelas.JamSelesai {
						return nil, apperrors.NewBadRequest("Jadwal kelas bentrok dengan kelas lain yang sudah Anda ajukan/ambil (" + dp.Kelas.Name + ")")
					}
				}
			}
		}
	}

	pengajuan := &domain.PengajuanKelas{
		ID:           uuid.NewString(),
		DosenID:      dosenID,
		KelasID:      req.KelasID,
		MataKuliahID: req.MataKuliahID,
		Status:       domain.StatusPending,
		Code:         generateRandomCode(),
	}

	if err := s.repo.CreatePengajuan(ctx, pengajuan); err != nil {
		log.Printf("[RequestKelas] ERROR: CreatePengajuan failed: %v", err)
		return nil, apperrors.NewInternal("Gagal mengajukan kelas", err.Error())
	}

	log.Printf("[RequestKelas] SUCCESS: pengajuan created with ID=%s", pengajuan.ID)
	pengajuan.Kelas = kelas
	return pengajuan, nil
}

func (s *kelasService) ApprovePengajuan(ctx context.Context, id string) error {
	p, err := s.repo.GetPengajuanByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Pengajuan tidak ditemukan")
	}

	if p.Status != domain.StatusPending {
		return apperrors.NewBadRequest("Hanya pengajuan berstatus pending yang dapat disetujui")
	}

	activePengajuan, err := s.repo.GetActivePengajuanByKelasID(ctx, p.KelasID)
	if err != nil {
		return apperrors.NewInternal("Gagal mengecek status kelas", err.Error())
	}

	for _, ap := range activePengajuan {
		if ap.Status == domain.StatusApproved && ap.ID != id {
			return apperrors.NewBadRequest("Kelas sudah disetujui untuk dosen lain")
		}
	}

	p.Status = domain.StatusApproved
	if err := s.repo.UpdatePengajuan(ctx, p); err != nil {
		return apperrors.NewInternal("Gagal menyetujui pengajuan", err.Error())
	}

	for _, ap := range activePengajuan {
		if ap.ID != id && ap.Status == domain.StatusPending {
			ap.Status = domain.StatusRejected
			s.repo.UpdatePengajuan(ctx, ap)
		}
	}

	return nil
}

func (s *kelasService) RejectPengajuan(ctx context.Context, id string) error {
	p, err := s.repo.GetPengajuanByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFound("Pengajuan tidak ditemukan")
	}

	if p.Status != domain.StatusPending {
		return apperrors.NewBadRequest("Hanya pengajuan berstatus pending yang dapat ditolak")
	}

	p.Status = domain.StatusRejected
	if err := s.repo.UpdatePengajuan(ctx, p); err != nil {
		return apperrors.NewInternal("Gagal menolak pengajuan", err.Error())
	}

	return nil
}

func (s *kelasService) GetMyPengajuan(ctx context.Context, dosenID string) ([]*domain.PengajuanKelas, error) {
	return s.repo.GetPengajuanByDosenID(ctx, dosenID)
}

func (s *kelasService) GetAllPengajuan(ctx context.Context) ([]*domain.PengajuanKelas, error) {
	return s.repo.GetAllPengajuan(ctx)
}

func (s *kelasService) GetMahasiswaInKelas(ctx context.Context, pengajuanID string, dosenID string) ([]*authDomain.User, error) {
	p, err := s.repo.GetPengajuanByID(ctx, pengajuanID)
	if err != nil {
		return nil, apperrors.NewNotFound("Pengajuan tidak ditemukan")
	}

	if dosenID != "" && p.DosenID != dosenID {
		return nil, apperrors.NewUnauthorized("Anda tidak memiliki akses ke kelas ini", "")
	}

	if p.Kelas == nil {
		return nil, apperrors.NewInternal("Data kelas tidak valid", "")
	}

	pesertaList, err := s.repo.GetPesertaKelasByPengajuanID(ctx, p.ID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil data peserta kelas", err.Error())
	}
	
	users := make([]*authDomain.User, 0, len(pesertaList))
	for _, psrt := range pesertaList {
		if psrt.Mahasiswa != nil {
			users = append(users, psrt.Mahasiswa)
		}
	}
	
	return users, nil
}

func (s *kelasService) GetMyJadwal(ctx context.Context, userID string) ([]*domain.PengajuanKelas, error) {
	pesertaList, err := s.repo.GetPesertaKelasByMahasiswaID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil jadwal Anda", err.Error())
	}
	
	jadwal := make([]*domain.PengajuanKelas, 0, len(pesertaList))
	for _, psrt := range pesertaList {
		if psrt.Pengajuan != nil {
			jadwal = append(jadwal, psrt.Pengajuan)
		}
	}
	return jadwal, nil
}

func (s *kelasService) GetAvailableKelas(ctx context.Context, userID string) ([]*domain.PengajuanKelas, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil data user", err.Error())
	}
	if user == nil || user.ProgramStudiID == nil || *user.ProgramStudiID == "" {
		return nil, apperrors.NewBadRequest("Program Studi belum diatur")
	}
	return s.repo.GetApprovedPengajuanByProdiID(ctx, *user.ProgramStudiID)
}

func (s *kelasService) AmbilKelas(ctx context.Context, userID string, pengajuanID string) error {
	p, err := s.repo.GetPengajuanByID(ctx, pengajuanID)
	if err != nil {
		return apperrors.NewNotFound("Kelas tidak ditemukan")
	}
	
	if p.Status != domain.StatusApproved {
		return apperrors.NewBadRequest("Kelas belum disetujui")
	}
	
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return apperrors.NewInternal("Gagal mengambil data user", err.Error())
	}
	if user.ProgramStudiID == nil || *user.ProgramStudiID != p.Kelas.ProgramStudiID {
		return apperrors.NewForbidden("Kelas ini tidak tersedia untuk Program Studi Anda")
	}
	
	count, err := s.repo.CountPesertaKelas(ctx, pengajuanID)
	if err != nil {
		return apperrors.NewInternal("Gagal menghitung peserta", err.Error())
	}
	
	if count >= int64(p.Kelas.Capacity) {
		return apperrors.NewBadRequest("Kelas sudah penuh")
	}
	
	pesertaExist, _ := s.repo.GetPesertaKelasByMahasiswaID(ctx, userID)
	for _, psrt := range pesertaExist {
		if psrt.PengajuanID == pengajuanID {
			return apperrors.NewBadRequest("Anda sudah mengambil kelas ini")
		}
	}
	
	mkConflict, err := s.repo.CheckPesertaMataKuliahConflict(ctx, userID, p.MataKuliahID)
	if err != nil {
		return apperrors.NewInternal("Gagal memeriksa konflik mata kuliah", err.Error())
	}
	if mkConflict {
		return apperrors.NewBadRequest("Anda sudah mengambil kelas lain untuk Mata Kuliah ini")
	}
	
	schedConflict, err := s.repo.CheckPesertaScheduleConflict(ctx, userID, p.Kelas.Hari, p.Kelas.JamMulai, p.Kelas.JamSelesai)
	if err != nil {
		return apperrors.NewInternal("Gagal memeriksa konflik jadwal", err.Error())
	}
	if schedConflict {
		return apperrors.NewBadRequest("Jadwal kelas ini bentrok dengan kelas Anda yang lain")
	}
	
	peserta := &domain.PesertaKelas{
		ID:          uuid.NewString(),
		PengajuanID: pengajuanID,
		MahasiswaID: userID,
		Status:      "enrolled",
	}
	
	if err := s.repo.CreatePesertaKelas(ctx, peserta); err != nil {
		return apperrors.NewInternal("Gagal mendaftar kelas", err.Error())
	}
	
	return nil
}

func (s *kelasService) MulaiPertemuan(ctx context.Context, pengajuanID string, judul string) (*domain.Pertemuan, error) {
	existing, err := s.repo.GetPertemuanByPengajuanID(ctx, pengajuanID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil data pertemuan", err.Error())
	}

	nomorNext := len(existing) + 1
	if nomorNext > domain.MaxPertemuan {
		return nil, apperrors.NewBadRequest("Semua 16 pertemuan sudah tercatat untuk kelas ini")
	}

	randBytes := make([]byte, 6)
	const charset = "0123456789"
	for i := range randBytes {
		randBytes[i] = charset[rand.Intn(len(charset))]
	}
	kodeAbsensi := string(randBytes)

	p := &domain.Pertemuan{
		ID:             uuid.NewString(),
		PengajuanID:    pengajuanID,
		Judul:          judul,
		Tanggal:        time.Now(),
		WaktuMulai:     time.Now(),
		Status:         domain.PertemuanStatusBerlangsung,
		KodeAbsensi:    kodeAbsensi,
		NomorPertemuan: nomorNext,
	}
	if err := s.repo.CreatePertemuan(ctx, p); err != nil {
		return nil, apperrors.NewInternal("Gagal memulai pertemuan", err.Error())
	}
	return p, nil
}

func (s *kelasService) AkhiriPertemuan(ctx context.Context, pertemuanID string) error {
	p, err := s.repo.GetPertemuanByID(ctx, pertemuanID)
	if err != nil {
		return apperrors.NewNotFound("Pertemuan tidak ditemukan")
	}
	if p.Status != domain.PertemuanStatusBerlangsung {
		return apperrors.NewBadRequest("Pertemuan sudah selesai atau tidak aktif")
	}
	now := time.Now()
	p.Status = domain.PertemuanStatusSelesai
	p.WaktuSelesai = &now
	if err := s.repo.UpdatePertemuan(ctx, p); err != nil {
		return apperrors.NewInternal("Gagal mengakhiri pertemuan", err.Error())
	}
	return nil
}

func (s *kelasService) GetPertemuanByPengajuan(ctx context.Context, pengajuanID string) ([]*domain.Pertemuan, error) {
	return s.repo.GetPertemuanByPengajuanID(ctx, pengajuanID)
}

func (s *kelasService) GetAbsensi(ctx context.Context, pertemuanID string) ([]*domain.Absensi, error) {
	return s.repo.GetAbsensiByPertemuanID(ctx, pertemuanID)
}

func (s *kelasService) SubmitAbsensi(ctx context.Context, pertemuanID string, data domain.BulkAbsensiRequest) error {
	existing, err := s.repo.GetAbsensiByPertemuanID(ctx, pertemuanID)
	if err != nil {
		return apperrors.NewInternal("Gagal mengecek absensi", err.Error())
	}

	if len(existing) == 0 {
		for _, v := range data.Data {
			a := &domain.Absensi{
				ID:              uuid.NewString(),
				PertemuanID:     pertemuanID,
				MahasiswaID:     v.MahasiswaID,
				StatusKehadiran: v.StatusKehadiran,
			}
			if err := s.repo.CreateAbsensi(ctx, a); err != nil {
				return apperrors.NewInternal("Gagal menyimpan absensi", err.Error())
			}
		}
		return nil
	}

	if err := s.repo.UpdateAbsensiBulk(ctx, pertemuanID, data.Data); err != nil {
		return apperrors.NewInternal("Gagal mengupdate absensi", err.Error())
	}
	return nil
}

func (s *kelasService) SubmitAbsensiMahasiswa(ctx context.Context, pertemuanID string, mahasiswaID string, kode string) error {
	p, err := s.repo.GetPertemuanByID(ctx, pertemuanID)
	if err != nil {
		return apperrors.NewNotFound("Pertemuan tidak ditemukan")
	}
	if p.Status != domain.PertemuanStatusBerlangsung {
		return apperrors.NewBadRequest("Pertemuan sudah tidak aktif")
	}
	if p.KodeAbsensi != kode {
		return apperrors.NewBadRequest("Kode absensi tidak valid")
	}

	existing, err := s.repo.GetAbsensiByPertemuanID(ctx, pertemuanID)
	if err != nil {
		return apperrors.NewInternal("Gagal memeriksa absensi", err.Error())
	}

	for _, a := range existing {
		if a.MahasiswaID == mahasiswaID {
			if a.StatusKehadiran == "hadir" {
				return nil
			}

			a.StatusKehadiran = "hadir"
			if err := s.repo.UpdateAbsensiBulk(ctx, pertemuanID, []domain.AbsensiUpdate{{MahasiswaID: mahasiswaID, StatusKehadiran: "hadir"}}); err != nil {
				return apperrors.NewInternal("Gagal update absensi", err.Error())
			}
			return nil
		}
	}

	a := &domain.Absensi{
		ID:              uuid.NewString(),
		PertemuanID:     pertemuanID,
		MahasiswaID:     mahasiswaID,
		StatusKehadiran: "hadir",
	}
	if err := s.repo.CreateAbsensi(ctx, a); err != nil {
		return apperrors.NewInternal("Gagal menyimpan absensi", err.Error())
	}

	return nil
}

func (s *kelasService) GetRekapKehadiran(ctx context.Context, pengajuanID string, dosenID string) (*domain.RekapKehadiranResponse, error) {

	p, err := s.repo.GetPengajuanByID(ctx, pengajuanID)
	if err != nil {
		return nil, apperrors.NewNotFound("Pengajuan tidak ditemukan")
	}
	if dosenID != "" && p.DosenID != dosenID {
		return nil, apperrors.NewForbidden("Anda tidak memiliki akses ke kelas ini")
	}

	pertemuanList, err := s.repo.GetPertemuanByPengajuanID(ctx, pengajuanID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil daftar pertemuan", err.Error())
	}

	students, err := s.GetMahasiswaInKelas(ctx, pengajuanID, dosenID)
	if err != nil {
		return nil, err
	}

	res := &domain.RekapKehadiranResponse{
		Pertemuan: make([]domain.PertemuanInfo, 0, len(pertemuanList)),
		Mahasiswa: make([]domain.MahasiswaRekap, 0, len(students)),
	}

	for _, p := range pertemuanList {
		res.Pertemuan = append(res.Pertemuan, domain.PertemuanInfo{
			ID:      p.ID,
			Judul:   p.Judul,
			Tanggal: p.Tanggal,
		})
	}

	for _, student := range students {
		nrp := "-"
		if student.NRP != nil {
			nrp = *student.NRP
		}
		rekap := domain.MahasiswaRekap{
			ID:        student.ID,
			NRP:       nrp,
			Name:      student.Name,
			Kehadiran: make(map[string]string),
		}
		res.Mahasiswa = append(res.Mahasiswa, rekap)
	}

	studentMap := make(map[string]*domain.MahasiswaRekap)
	for i := range res.Mahasiswa {
		studentMap[res.Mahasiswa[i].ID] = &res.Mahasiswa[i]
	}

	for _, p := range pertemuanList {
		absensiList, err := s.repo.GetAbsensiByPertemuanID(ctx, p.ID)
		if err != nil {
			continue
		}

		fmt.Println("Absensi count for meeting", p.ID, ":", len(absensiList))
		for _, a := range absensiList {
			fmt.Println("Found absensi for", a.MahasiswaID, "with status", a.StatusKehadiran)
			if m, ok := studentMap[a.MahasiswaID]; ok {
				m.Kehadiran[p.ID] = a.StatusKehadiran
				fmt.Println("Mapped to student", m.Name)
			} else {
				fmt.Println("Student not found in map!")
			}
		}

		for _, m := range res.Mahasiswa {
			if _, ok := m.Kehadiran[p.ID]; !ok {
				studentMap[m.ID].Kehadiran[p.ID] = domain.AbsensiAlpa
			}
		}
	}

	return res, nil
}

func generateRandomCode() string {
	b := make([]byte, 6)
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
