package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	authDomain "siakad-pro/internal/modules/auth/domain"
	"siakad-pro/internal/modules/kelas/domain"
	"siakad-pro/internal/shared/apperrors"
)

type kelasService struct {
	repo            domain.KelasRepository
	userProvider    domain.UserProvider
	periodeProvider domain.PeriodeProvider
	mkProvider      domain.MataKuliahProvider
}

func NewKelasService(
	repo domain.KelasRepository,
	userProvider domain.UserProvider,
	periodeProvider domain.PeriodeProvider,
	mkProvider domain.MataKuliahProvider,
) domain.KelasService {
	return &kelasService{
		repo:            repo,
		userProvider:    userProvider,
		periodeProvider: periodeProvider,
		mkProvider:      mkProvider,
	}
}

var validHari = map[string]bool{"Senin": true, "Selasa": true, "Rabu": true, "Kamis": true, "Jumat": true}

func (s *kelasService) Create(ctx context.Context, req domain.CreateKelasRequest) (*domain.Kelas, error) {

	matched, _ := regexp.MatchString(`^[A-Z]{2,4}-[1-9]0[1-9]$`, req.Name)
	if !matched {
		return nil, apperrors.NewBadRequest("Format nama kelas tidak valid (contoh yang benar: [KODE]-101 s/d [KODE]-409, misal IF-101, RPL-201)")
	}

	if req.Capacity < domain.MinCapacity || req.Capacity > domain.MaxCapacity {
		return nil, apperrors.NewBadRequest(fmt.Sprintf("Kapasitas kelas harus antara %d dan %d", domain.MinCapacity, domain.MaxCapacity))
	}

	conflict, _ := s.repo.CheckScheduleConflict(ctx, req.Name, req.Hari, req.JamMulai)
	if conflict {
		return nil, apperrors.NewBadRequest("Kelas tersebut sudah terdaftar pada hari dan jam yang sama")
	}

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
		return apperrors.NewInternal("Gagal menghapus kelas", err.Error())
	}
	return nil
}

// --- KRS & Mahasiswa ---

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
	user, err := s.userProvider.GetUserByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternal("Gagal mengambil data user", err.Error())
	}
	if user == nil || user.ProgramStudiID == nil || *user.ProgramStudiID == "" {
		return nil, apperrors.NewBadRequest("Program Studi belum diatur")
	}
	return s.repo.GetApprovedPengajuanByProdiID(ctx, *user.ProgramStudiID)
}

func (s *kelasService) AmbilKelas(ctx context.Context, userID string, pengajuanID string) error {
	return s.repo.Transaction(ctx, func(txRepo domain.KelasRepository) error {
		p, err := txRepo.LockPengajuanByID(ctx, pengajuanID)
		if err != nil {
			return apperrors.NewNotFound("Kelas tidak ditemukan")
		}

		if p.Status != domain.StatusApproved {
			return apperrors.NewBadRequest("Kelas belum disetujui")
		}

		user, err := s.userProvider.GetUserByID(ctx, userID)
		if err != nil {
			return apperrors.NewInternal("Gagal mengambil data user", err.Error())
		}
		if user.ProgramStudiID == nil || *user.ProgramStudiID != p.Kelas.ProgramStudiID {
			// MKUB / DEPT bypass: if the class is DEPT or MKUB, it might be open to all, but for now we follow original logic
			return apperrors.NewForbidden("Kelas ini tidak tersedia untuk Program Studi Anda")
		}

		count, err := txRepo.CountPesertaKelas(ctx, pengajuanID)
		if err != nil {
			return apperrors.NewInternal("Gagal menghitung peserta", err.Error())
		}

		if count >= int64(p.Kelas.Capacity) {
			return apperrors.NewBadRequest("Kelas sudah penuh")
		}

		pesertaExist, _ := txRepo.GetPesertaKelasByMahasiswaID(ctx, userID)
		for _, psrt := range pesertaExist {
			if psrt.PengajuanID == pengajuanID {
				return apperrors.NewBadRequest("Anda sudah mengambil kelas ini")
			}
		}

		mkConflict, err := txRepo.CheckPesertaMataKuliahConflict(ctx, userID, p.MataKuliahID)
		if err != nil {
			return apperrors.NewInternal("Gagal memeriksa konflik mata kuliah", err.Error())
		}
		if mkConflict {
			return apperrors.NewBadRequest("Anda sudah mengambil kelas lain untuk Mata Kuliah ini")
		}

		schedConflict, err := txRepo.CheckPesertaScheduleConflict(ctx, userID, p.Kelas.Hari, p.Kelas.JamMulai, p.Kelas.JamSelesai)
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

		if err := txRepo.CreatePesertaKelas(ctx, peserta); err != nil {
			return apperrors.NewInternal("Gagal mendaftar kelas", err.Error())
		}

		return nil
	})
}
