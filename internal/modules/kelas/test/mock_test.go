package test

import (
	"context"

	authDomain "siakad-pro/internal/modules/auth/domain"
	"siakad-pro/internal/modules/kelas/domain"
)

// --- Mock Implementations for Isolated Unit Testing ---

type mockKelasRepo struct {
	// Core Kelas
	CreateFn                func(ctx context.Context, kelas *domain.Kelas) error
	GetAllFn                func(ctx context.Context) ([]*domain.Kelas, error)
	GetByIDFn               func(ctx context.Context, id string) (*domain.Kelas, error)
	GetByNameFn             func(ctx context.Context, name string) (*domain.Kelas, error)
	CheckScheduleConflictFn func(ctx context.Context, name string, hari string, jamMulai string) (bool, error)
	DeleteFn                func(ctx context.Context, id string) error

	// Pengajuan Kelas
	CreatePengajuanFn               func(ctx context.Context, p *domain.PengajuanKelas) error
	GetPengajuanByIDFn              func(ctx context.Context, id string) (*domain.PengajuanKelas, error)
	GetPengajuanByDosenIDFn         func(ctx context.Context, dosenID string) ([]*domain.PengajuanKelas, error)
	GetActivePengajuanByKelasIDFn   func(ctx context.Context, kelasID string) ([]*domain.PengajuanKelas, error)
	GetAllPengajuanFn               func(ctx context.Context) ([]*domain.PengajuanKelas, error)
	UpdatePengajuanFn               func(ctx context.Context, p *domain.PengajuanKelas) error
	DeletePengajuanFn               func(ctx context.Context, id string) error
	LockPengajuanByIDFn             func(ctx context.Context, id string) (*domain.PengajuanKelas, error)
	GetApprovedPengajuanByProdiIDFn func(ctx context.Context, prodiID string) ([]*domain.PengajuanKelas, error)

	// Peserta Kelas
	CreatePesertaKelasFn             func(ctx context.Context, p *domain.PesertaKelas) error
	GetPesertaKelasByPengajuanIDFn   func(ctx context.Context, pengajuanID string) ([]*domain.PesertaKelas, error)
	GetPesertaKelasByMahasiswaIDFn   func(ctx context.Context, mahasiswaID string) ([]*domain.PesertaKelas, error)
	CountPesertaKelasFn              func(ctx context.Context, pengajuanID string) (int64, error)
	CheckPesertaMataKuliahConflictFn func(ctx context.Context, mahasiswaID string, mkID string) (bool, error)
	CheckPesertaScheduleConflictFn   func(ctx context.Context, mahasiswaID string, hari string, jamMulai string, jamSelesai string) (bool, error)

	// Pertemuan
	CreatePertemuanFn           func(ctx context.Context, p *domain.Pertemuan) error
	GetPertemuanByIDFn          func(ctx context.Context, id string) (*domain.Pertemuan, error)
	GetPertemuanByPengajuanIDFn func(ctx context.Context, pengajuanID string) ([]*domain.Pertemuan, error)
	UpdatePertemuanFn           func(ctx context.Context, p *domain.Pertemuan) error

	// Absensi
	CreateAbsensiFn           func(ctx context.Context, a *domain.Absensi) error
	GetAbsensiByPertemuanIDFn func(ctx context.Context, pertemuanID string) ([]*domain.Absensi, error)
	UpdateAbsensiBulkFn       func(ctx context.Context, pertemuanID string, data []domain.AbsensiUpdate) error

	// Transaction
	TransactionFn func(ctx context.Context, fn func(txRepo domain.KelasRepository) error) error
}

func (m *mockKelasRepo) Create(ctx context.Context, kelas *domain.Kelas) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, kelas)
	}
	return nil
}

func (m *mockKelasRepo) GetAll(ctx context.Context) ([]*domain.Kelas, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn(ctx)
	}
	return nil, nil
}

func (m *mockKelasRepo) GetByID(ctx context.Context, id string) (*domain.Kelas, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return &domain.Kelas{ID: id, Name: "IF-101", Capacity: 30}, nil
}

func (m *mockKelasRepo) GetByName(ctx context.Context, name string) (*domain.Kelas, error) {
	if m.GetByNameFn != nil {
		return m.GetByNameFn(ctx, name)
	}
	return nil, nil
}

func (m *mockKelasRepo) CheckScheduleConflict(ctx context.Context, name string, hari string, jamMulai string) (bool, error) {
	if m.CheckScheduleConflictFn != nil {
		return m.CheckScheduleConflictFn(ctx, name, hari, jamMulai)
	}
	return false, nil
}

func (m *mockKelasRepo) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *mockKelasRepo) CreatePengajuan(ctx context.Context, p *domain.PengajuanKelas) error {
	if m.CreatePengajuanFn != nil {
		return m.CreatePengajuanFn(ctx, p)
	}
	return nil
}

func (m *mockKelasRepo) GetPengajuanByID(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
	if m.GetPengajuanByIDFn != nil {
		return m.GetPengajuanByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockKelasRepo) GetPengajuanByDosenID(ctx context.Context, dosenID string) ([]*domain.PengajuanKelas, error) {
	if m.GetPengajuanByDosenIDFn != nil {
		return m.GetPengajuanByDosenIDFn(ctx, dosenID)
	}
	return nil, nil
}

func (m *mockKelasRepo) GetActivePengajuanByKelasID(ctx context.Context, kelasID string) ([]*domain.PengajuanKelas, error) {
	if m.GetActivePengajuanByKelasIDFn != nil {
		return m.GetActivePengajuanByKelasIDFn(ctx, kelasID)
	}
	return nil, nil
}

func (m *mockKelasRepo) GetAllPengajuan(ctx context.Context) ([]*domain.PengajuanKelas, error) {
	if m.GetAllPengajuanFn != nil {
		return m.GetAllPengajuanFn(ctx)
	}
	return nil, nil
}

func (m *mockKelasRepo) UpdatePengajuan(ctx context.Context, p *domain.PengajuanKelas) error {
	if m.UpdatePengajuanFn != nil {
		return m.UpdatePengajuanFn(ctx, p)
	}
	return nil
}

func (m *mockKelasRepo) DeletePengajuan(ctx context.Context, id string) error {
	if m.DeletePengajuanFn != nil {
		return m.DeletePengajuanFn(ctx, id)
	}
	return nil
}

func (m *mockKelasRepo) LockPengajuanByID(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
	if m.LockPengajuanByIDFn != nil {
		return m.LockPengajuanByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockKelasRepo) GetApprovedPengajuanByProdiID(ctx context.Context, prodiID string) ([]*domain.PengajuanKelas, error) {
	if m.GetApprovedPengajuanByProdiIDFn != nil {
		return m.GetApprovedPengajuanByProdiIDFn(ctx, prodiID)
	}
	return nil, nil
}

func (m *mockKelasRepo) CreatePesertaKelas(ctx context.Context, p *domain.PesertaKelas) error {
	if m.CreatePesertaKelasFn != nil {
		return m.CreatePesertaKelasFn(ctx, p)
	}
	return nil
}

func (m *mockKelasRepo) GetPesertaKelasByPengajuanID(ctx context.Context, pengajuanID string) ([]*domain.PesertaKelas, error) {
	if m.GetPesertaKelasByPengajuanIDFn != nil {
		return m.GetPesertaKelasByPengajuanIDFn(ctx, pengajuanID)
	}
	return nil, nil
}

func (m *mockKelasRepo) GetPesertaKelasByMahasiswaID(ctx context.Context, mahasiswaID string) ([]*domain.PesertaKelas, error) {
	if m.GetPesertaKelasByMahasiswaIDFn != nil {
		return m.GetPesertaKelasByMahasiswaIDFn(ctx, mahasiswaID)
	}
	return nil, nil
}

func (m *mockKelasRepo) CountPesertaKelas(ctx context.Context, pengajuanID string) (int64, error) {
	if m.CountPesertaKelasFn != nil {
		return m.CountPesertaKelasFn(ctx, pengajuanID)
	}
	return 0, nil
}

func (m *mockKelasRepo) CheckPesertaMataKuliahConflict(ctx context.Context, mahasiswaID string, mkID string) (bool, error) {
	if m.CheckPesertaMataKuliahConflictFn != nil {
		return m.CheckPesertaMataKuliahConflictFn(ctx, mahasiswaID, mkID)
	}
	return false, nil
}

func (m *mockKelasRepo) CheckPesertaScheduleConflict(ctx context.Context, mahasiswaID string, hari string, jamMulai string, jamSelesai string) (bool, error) {
	if m.CheckPesertaScheduleConflictFn != nil {
		return m.CheckPesertaScheduleConflictFn(ctx, mahasiswaID, hari, jamMulai, jamSelesai)
	}
	return false, nil
}

func (m *mockKelasRepo) CreatePertemuan(ctx context.Context, p *domain.Pertemuan) error {
	if m.CreatePertemuanFn != nil {
		return m.CreatePertemuanFn(ctx, p)
	}
	return nil
}

func (m *mockKelasRepo) GetPertemuanByID(ctx context.Context, id string) (*domain.Pertemuan, error) {
	if m.GetPertemuanByIDFn != nil {
		return m.GetPertemuanByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockKelasRepo) GetPertemuanByPengajuanID(ctx context.Context, pengajuanID string) ([]*domain.Pertemuan, error) {
	if m.GetPertemuanByPengajuanIDFn != nil {
		return m.GetPertemuanByPengajuanIDFn(ctx, pengajuanID)
	}
	return nil, nil
}

func (m *mockKelasRepo) UpdatePertemuan(ctx context.Context, p *domain.Pertemuan) error {
	if m.UpdatePertemuanFn != nil {
		return m.UpdatePertemuanFn(ctx, p)
	}
	return nil
}

func (m *mockKelasRepo) CreateAbsensi(ctx context.Context, a *domain.Absensi) error {
	if m.CreateAbsensiFn != nil {
		return m.CreateAbsensiFn(ctx, a)
	}
	return nil
}

func (m *mockKelasRepo) GetAbsensiByPertemuanID(ctx context.Context, pertemuanID string) ([]*domain.Absensi, error) {
	if m.GetAbsensiByPertemuanIDFn != nil {
		return m.GetAbsensiByPertemuanIDFn(ctx, pertemuanID)
	}
	return nil, nil
}

func (m *mockKelasRepo) UpdateAbsensiBulk(ctx context.Context, pertemuanID string, data []domain.AbsensiUpdate) error {
	if m.UpdateAbsensiBulkFn != nil {
		return m.UpdateAbsensiBulkFn(ctx, pertemuanID, data)
	}
	return nil
}

func (m *mockKelasRepo) Transaction(ctx context.Context, fn func(txRepo domain.KelasRepository) error) error {
	if m.TransactionFn != nil {
		return m.TransactionFn(ctx, fn)
	}
	return fn(m)
}

// --- Mock External Providers ---

type mockUserProvider struct {
	GetUserByIDFn func(ctx context.Context, userID string) (*authDomain.User, error)
}

func (m *mockUserProvider) GetUserByID(ctx context.Context, userID string) (*authDomain.User, error) {
	if m.GetUserByIDFn != nil {
		return m.GetUserByIDFn(ctx, userID)
	}
	prodiID := "prodi-ti-1"
	return &authDomain.User{ID: userID, ProgramStudiID: &prodiID}, nil
}

type mockPeriodeProvider struct {
	GetActivePeriodeIDFn func(ctx context.Context) (string, error)
}

func (m *mockPeriodeProvider) GetActivePeriodeID(ctx context.Context) (string, error) {
	if m.GetActivePeriodeIDFn != nil {
		return m.GetActivePeriodeIDFn(ctx)
	}
	return "periode-2024-genap", nil
}

type mockMataKuliahProvider struct {
	IsMataKuliahValidForKelasFn func(ctx context.Context, dosenID string, mkID string, prodiID string) (bool, error)
}

func (m *mockMataKuliahProvider) IsMataKuliahValidForKelas(ctx context.Context, dosenID string, mkID string, prodiID string) (bool, error) {
	if m.IsMataKuliahValidForKelasFn != nil {
		return m.IsMataKuliahValidForKelasFn(ctx, dosenID, mkID, prodiID)
	}
	return true, nil
}

// helper valid request
func defaultValidKelasRequest() domain.CreateKelasRequest {
	return domain.CreateKelasRequest{
		Name:           "IF-101",
		Capacity:       30,
		Hari:           "Senin",
		JamMulai:       "08:00",
		JamSelesai:     "10:30",
		ProgramStudiID: "prodi-ti-1",
	}
}
