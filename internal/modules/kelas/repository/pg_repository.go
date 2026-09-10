package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	authDomain "siakad-pro/internal/modules/auth/domain"
	"siakad-pro/internal/modules/kelas/domain"
)

type pgKelasRepository struct {
	db *gorm.DB
}

func NewPgKelasRepository(db *gorm.DB) domain.KelasRepository {
	return &pgKelasRepository{db: db}
}

func (r *pgKelasRepository) Create(ctx context.Context, kelas *domain.Kelas) error {
	return r.db.WithContext(ctx).Create(kelas).Error
}

func (r *pgKelasRepository) GetAll(ctx context.Context) ([]*domain.Kelas, error) {
	var kelases []*domain.Kelas
	err := r.db.WithContext(ctx).
		Preload("ProgramStudi").
		Preload("Pengajuan", "status IN ?", []string{domain.StatusPending, domain.StatusApproved}).
		Preload("Pengajuan.Dosen").
		Order("created_at desc").Find(&kelases).Error
	return kelases, err
}

func (r *pgKelasRepository) GetByID(ctx context.Context, id string) (*domain.Kelas, error) {
	var kelas domain.Kelas
	err := r.db.WithContext(ctx).
		Preload("ProgramStudi").
		Preload("Pengajuan", "status IN ?", []string{domain.StatusPending, domain.StatusApproved}).
		Preload("Pengajuan.Dosen").
		Preload("Pengajuan.MataKuliah").
		Where("id = ?", id).First(&kelas).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kelas not found")
		}
		return nil, err
	}
	return &kelas, nil
}

func (r *pgKelasRepository) GetByName(ctx context.Context, name string) (*domain.Kelas, error) {
	var kelas domain.Kelas
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&kelas).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kelas not found")
		}
		return nil, err
	}
	return &kelas, nil
}

func (r *pgKelasRepository) CheckScheduleConflict(ctx context.Context, name string, hari string, jamMulai string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Kelas{}).
		Where("name = ? AND hari = ? AND jam_mulai = ?", name, hari, jamMulai).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *pgKelasRepository) HasApprovedMataKuliah(ctx context.Context, dosenID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("pengajuan_mata_kuliahs").
		Where("dosen_id = ? AND status = ?", dosenID, domain.StatusApproved).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *pgKelasRepository) IsMataKuliahValidForKelas(ctx context.Context, dosenID string, mkID string, prodiID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("pengajuan_mata_kuliahs").
		Joins("JOIN mata_kuliahs mk ON mk.id = pengajuan_mata_kuliahs.mata_kuliah_id").
		Where("pengajuan_mata_kuliahs.dosen_id = ? AND pengajuan_mata_kuliahs.mata_kuliah_id = ? AND pengajuan_mata_kuliahs.status = ? AND mk.program_studi_id = ?", dosenID, mkID, domain.StatusApproved, prodiID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *pgKelasRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Kelas{}).Error
}

func (r *pgKelasRepository) CreatePengajuan(ctx context.Context, p *domain.PengajuanKelas) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *pgKelasRepository) GetPengajuanByID(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
	var p domain.PengajuanKelas
	err := r.db.WithContext(ctx).Preload("Kelas").Preload("Dosen").Preload("MataKuliah").Where("id = ?", id).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengajuan not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgKelasRepository) GetPengajuanByDosenID(ctx context.Context, dosenID string) ([]*domain.PengajuanKelas, error) {
	var list []*domain.PengajuanKelas
	err := r.db.WithContext(ctx).Preload("Kelas").Preload("Kelas.ProgramStudi").Preload("MataKuliah").Where("dosen_id = ?", dosenID).Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *pgKelasRepository) GetActivePengajuanByKelasID(ctx context.Context, kelasID string) ([]*domain.PengajuanKelas, error) {
	var list []*domain.PengajuanKelas
	err := r.db.WithContext(ctx).Where("kelas_id = ? AND status IN ?", kelasID, []string{domain.StatusPending, domain.StatusApproved}).Find(&list).Error
	return list, err
}

func (r *pgKelasRepository) GetAllPengajuan(ctx context.Context) ([]*domain.PengajuanKelas, error) {
	var list []*domain.PengajuanKelas
	err := r.db.WithContext(ctx).Preload("Kelas").Preload("Kelas.ProgramStudi").Preload("Dosen").Preload("MataKuliah").Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *pgKelasRepository) UpdatePengajuan(ctx context.Context, p *domain.PengajuanKelas) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *pgKelasRepository) DeletePengajuan(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.PengajuanKelas{}).Error
}

func (r *pgKelasRepository) GetMahasiswaByProgramStudiID(ctx context.Context, prodiID string) ([]*authDomain.User, error) {
	var users []*authDomain.User
	err := r.db.WithContext(ctx).Table("users").Where("role = ? AND program_studi_id = ?", "mahasiswa", prodiID).Find(&users).Error
	return users, err
}

func (r *pgKelasRepository) GetApprovedPengajuanByProdiID(ctx context.Context, prodiID string) ([]*domain.PengajuanKelas, error) {
	var list []*domain.PengajuanKelas
	err := r.db.WithContext(ctx).
		Joins("JOIN kelas ON kelas.id = pengajuan_kelas.kelas_id").
		Preload("Kelas").
		Preload("Dosen").
		Preload("MataKuliah").
		Where("kelas.program_studi_id = ? AND pengajuan_kelas.status = ?", prodiID, domain.StatusApproved).
		Order("pengajuan_kelas.created_at desc").
		Find(&list).Error
	return list, err
}

func (r *pgKelasRepository) GetUserByID(ctx context.Context, userID string) (*authDomain.User, error) {
	var user authDomain.User
	err := r.db.WithContext(ctx).Table("users").Where("id = ?", userID).First(&user).Error
	return &user, err
}

func (r *pgKelasRepository) CreatePertemuan(ctx context.Context, p *domain.Pertemuan) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *pgKelasRepository) GetPertemuanByID(ctx context.Context, id string) (*domain.Pertemuan, error) {
	var p domain.Pertemuan
	err := r.db.WithContext(ctx).Preload("Pengajuan").Where("id = ?", id).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pertemuan not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgKelasRepository) GetPertemuanByPengajuanID(ctx context.Context, pengajuanID string) ([]*domain.Pertemuan, error) {
	var list []*domain.Pertemuan
	err := r.db.WithContext(ctx).Where("pengajuan_id = ?", pengajuanID).Order("nomor_pertemuan asc").Find(&list).Error
	return list, err
}

func (r *pgKelasRepository) UpdatePertemuan(ctx context.Context, p *domain.Pertemuan) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *pgKelasRepository) CreateAbsensi(ctx context.Context, a *domain.Absensi) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *pgKelasRepository) GetAbsensiByPertemuanID(ctx context.Context, pertemuanID string) ([]*domain.Absensi, error) {
	var list []*domain.Absensi
	err := r.db.WithContext(ctx).Preload("Mahasiswa").Where("pertemuan_id = ?", pertemuanID).Find(&list).Error
	return list, err
}

func (r *pgKelasRepository) UpdateAbsensiBulk(ctx context.Context, pertemuanID string, data []domain.AbsensiUpdate) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, update := range data {
			err := tx.Model(&domain.Absensi{}).
				Where("pertemuan_id = ? AND mahasiswa_id = ?", pertemuanID, update.MahasiswaID).
				Update("status_kehadiran", update.StatusKehadiran).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *pgKelasRepository) CreatePesertaKelas(ctx context.Context, p *domain.PesertaKelas) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *pgKelasRepository) GetPesertaKelasByPengajuanID(ctx context.Context, pengajuanID string) ([]*domain.PesertaKelas, error) {
	var list []*domain.PesertaKelas
	err := r.db.WithContext(ctx).Preload("Mahasiswa").Where("pengajuan_id = ?", pengajuanID).Find(&list).Error
	return list, err
}

func (r *pgKelasRepository) GetPesertaKelasByMahasiswaID(ctx context.Context, mahasiswaID string) ([]*domain.PesertaKelas, error) {
	var list []*domain.PesertaKelas
	err := r.db.WithContext(ctx).
		Preload("Pengajuan").
		Preload("Pengajuan.Kelas").
		Preload("Pengajuan.Dosen").
		Preload("Pengajuan.MataKuliah").
		Where("mahasiswa_id = ?", mahasiswaID).Find(&list).Error
	return list, err
}

func (r *pgKelasRepository) CountPesertaKelas(ctx context.Context, pengajuanID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.PesertaKelas{}).Where("pengajuan_id = ?", pengajuanID).Count(&count).Error
	return count, err
}

func (r *pgKelasRepository) CheckPesertaMataKuliahConflict(ctx context.Context, mahasiswaID string, mkID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("peserta_kelas").
		Joins("JOIN pengajuan_kelas pk ON pk.id = peserta_kelas.pengajuan_id").
		Where("peserta_kelas.mahasiswa_id = ? AND pk.mata_kuliah_id = ?", mahasiswaID, mkID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *pgKelasRepository) CheckPesertaScheduleConflict(ctx context.Context, mahasiswaID string, hari string, jamMulai string, jamSelesai string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("peserta_kelas").
		Joins("JOIN pengajuan_kelas pk ON pk.id = peserta_kelas.pengajuan_id").
		Joins("JOIN kelas k ON k.id = pk.kelas_id").
		Where("peserta_kelas.mahasiswa_id = ? AND k.hari = ?", mahasiswaID, hari).
		Where("((k.jam_mulai <= ? AND k.jam_selesai > ?) OR (k.jam_mulai < ? AND k.jam_selesai >= ?) OR (? <= k.jam_mulai AND ? >= k.jam_selesai))", 
			jamSelesai, jamMulai, jamSelesai, jamMulai, jamMulai, jamSelesai).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
