package service

import (
	"context"

	"github.com/google/uuid"
	"siakad-pro/internal/modules/kelas/domain"
	"siakad-pro/internal/shared/apperrors"
)

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

		for _, a := range absensiList {
			if m, ok := studentMap[a.MahasiswaID]; ok {
				m.Kehadiran[p.ID] = a.StatusKehadiran
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
