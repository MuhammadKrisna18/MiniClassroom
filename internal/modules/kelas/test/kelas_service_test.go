package test

import (
	"context"
	"errors"
	"strings"
	"testing"

	authDomain "siakad-pro/internal/modules/auth/domain"
	"siakad-pro/internal/modules/kelas/domain"
	"siakad-pro/internal/modules/kelas/service"
	"siakad-pro/internal/shared/apperrors"
)

// ==============================================================================
// 1. UNIT TEST: KelasCoreService.Create using EP and BVA
// ==============================================================================

func TestKelasService_Create_EP_BVA(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		reqModifier func(req *domain.CreateKelasRequest)
		mockSetup   func(m *mockKelasRepo)
		wantErr     bool
		errContains string
	}{
		// --------------------------------------------------------------------------
		// BVA & EP: Kapasitas Kelas (MinCapacity = 25, MaxCapacity = 50)
		// --------------------------------------------------------------------------
		{
			name: "BVA - Capacity 24 [Just below min 25 - Invalid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Capacity = 24
			},
			wantErr:     true,
			errContains: "Kapasitas kelas harus antara 25 dan 50",
		},
		{
			name: "BVA - Capacity 25 [Min boundary - Valid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Capacity = 25
			},
			wantErr: false,
		},
		{
			name: "BVA - Capacity 26 [Just above min 25 - Valid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Capacity = 26
			},
			wantErr: false,
		},
		{
			name: "BVA - Capacity 49 [Just below max 50 - Valid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Capacity = 49
			},
			wantErr: false,
		},
		{
			name: "BVA - Capacity 50 [Max boundary - Valid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Capacity = 50
			},
			wantErr: false,
		},
		{
			name: "BVA - Capacity 51 [Just above max 50 - Invalid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Capacity = 51
			},
			wantErr:     true,
			errContains: "Kapasitas kelas harus antara 25 dan 50",
		},
		{
			name: "EP - Capacity -5 [Negative number partition - Invalid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Capacity = -5
			},
			wantErr:     true,
			errContains: "Kapasitas kelas harus antara 25 dan 50",
		},
		{
			name: "EP - Capacity 0 [Zero partition - Invalid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Capacity = 0
			},
			wantErr:     true,
			errContains: "Kapasitas kelas harus antara 25 dan 50",
		},
		{
			name: "EP - Capacity 35 [Nominal valid partition - Valid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Capacity = 35
			},
			wantErr: false,
		},

		// --------------------------------------------------------------------------
		// BVA & EP: Format Nama Kelas (Regex: ^[A-Z]{2,4}-[1-9]0[1-9]$)
		// --------------------------------------------------------------------------
		{
			name: "BVA - Name 'A-101' [Prodi length 1 char, just below min 2 - Invalid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Name = "A-101"
			},
			wantErr:     true,
			errContains: "Format nama kelas tidak valid",
		},
		{
			name: "BVA - Name 'IF-101' [Prodi length 2 chars, min boundary - Valid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Name = "IF-101"
			},
			wantErr: false,
		},
		{
			name: "BVA - Name 'RPL-201' [Prodi length 3 chars, nominal - Valid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Name = "RPL-201"
			},
			wantErr: false,
		},
		{
			name: "BVA - Name 'INFT-301' [Prodi length 4 chars, max boundary - Valid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Name = "INFT-301"
			},
			wantErr: false,
		},
		{
			name: "BVA - Name 'INFOT-101' [Prodi length 5 chars, just above max 4 - Invalid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Name = "INFOT-101"
			},
			wantErr:     true,
			errContains: "Format nama kelas tidak valid",
		},
		{
			name: "EP - Name 'if-101' [Lowercase partition - Invalid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Name = "if-101"
			},
			wantErr:     true,
			errContains: "Format nama kelas tidak valid",
		},
		{
			name: "BVA - Name 'IF-001' [Semester digit 0, just below min 1 - Invalid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Name = "IF-001"
			},
			wantErr:     true,
			errContains: "Format nama kelas tidak valid",
		},
		{
			name: "BVA - Name 'IF-101' [Semester digit 1, min boundary - Valid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Name = "IF-101"
			},
			wantErr: false,
		},
		{
			name: "BVA - Name 'IF-901' [Semester digit 9, max boundary - Valid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Name = "IF-901"
			},
			wantErr: false,
		},
		{
			name: "EP - Name 'IF-111' [Middle digit non-zero, invalid partition]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Name = "IF-111"
			},
			wantErr:     true,
			errContains: "Format nama kelas tidak valid",
		},
		{
			name: "BVA - Name 'IF-100' [Class digit 0, just below min 1 - Invalid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Name = "IF-100"
			},
			wantErr:     true,
			errContains: "Format nama kelas tidak valid",
		},
		{
			name: "BVA - Name 'IF-101' [Class digit 1, min boundary - Valid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Name = "IF-101"
			},
			wantErr: false,
		},
		{
			name: "BVA - Name 'IF-109' [Class digit 9, max boundary - Valid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Name = "IF-109"
			},
			wantErr: false,
		},

		// --------------------------------------------------------------------------
		// EP: Validasi Hari Kuliah (Senin s/d Jumat)
		// --------------------------------------------------------------------------
		{
			name: "EP - Hari 'Senin' [Valid partition]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Hari = "Senin"
			},
			wantErr: false,
		},
		{
			name: "EP - Hari 'Rabu' [Valid partition]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Hari = "Rabu"
			},
			wantErr: false,
		},
		{
			name: "EP - Hari 'Jumat' [Valid partition]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Hari = "Jumat"
			},
			wantErr: false,
		},
		{
			name: "EP - Hari 'Sabtu' [Weekend partition - Invalid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Hari = "Sabtu"
			},
			wantErr:     true,
			errContains: "Hari harus antara Senin sampai Jumat",
		},
		{
			name: "EP - Hari 'Minggu' [Weekend partition - Invalid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Hari = "Minggu"
			},
			wantErr:     true,
			errContains: "Hari harus antara Senin sampai Jumat",
		},
		{
			name: "EP - Hari '' [Empty string partition - Invalid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Hari = ""
			},
			wantErr:     true,
			errContains: "Hari harus antara Senin sampai Jumat",
		},
		{
			name: "EP - Hari 'senin' [Lowercase partition - Invalid]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.Hari = "senin"
			},
			wantErr:     true,
			errContains: "Hari harus antara Senin sampai Jumat",
		},

		// --------------------------------------------------------------------------
		// EP: Waktu Jam Mulai & Jam Selesai
		// --------------------------------------------------------------------------
		{
			name: "EP - JamMulai empty [Invalid partition]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.JamMulai = ""
			},
			wantErr:     true,
			errContains: "Jam mulai dan selesai harus diisi",
		},
		{
			name: "EP - JamSelesai empty [Invalid partition]",
			reqModifier: func(req *domain.CreateKelasRequest) {
				req.JamSelesai = ""
			},
			wantErr:     true,
			errContains: "Jam mulai dan selesai harus diisi",
		},

		// --------------------------------------------------------------------------
		// EP: Konflik Jadwal di Repository
		// --------------------------------------------------------------------------
		{
			name: "EP - Schedule Conflict Detected [Conflict partition - Invalid]",
			mockSetup: func(m *mockKelasRepo) {
				m.CheckScheduleConflictFn = func(ctx context.Context, name string, hari string, jamMulai string) (bool, error) {
					return true, nil
				}
			},
			wantErr:     true,
			errContains: "Kelas tersebut sudah terdaftar pada hari dan jam yang sama",
		},

		// --------------------------------------------------------------------------
		// EP: Kegagalan Database Repository
		// --------------------------------------------------------------------------
		{
			name: "EP - Database Create Error [Database failure partition]",
			mockSetup: func(m *mockKelasRepo) {
				m.CreateFn = func(ctx context.Context, kelas *domain.Kelas) error {
					return errors.New("db connection failure")
				}
			},
			wantErr:     true,
			errContains: "Gagal membuat kelas: db connection failure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockKelasRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			// Instantiate service with segregated mock
			svc := service.NewKelasService(mockRepo, nil, nil, nil)

			req := defaultValidKelasRequest()
			if tt.reqModifier != nil {
				tt.reqModifier(&req)
			}

			result, err := svc.Create(ctx, req)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, but got nil", tt.errContains)
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain expected substring %q", err.Error(), tt.errContains)
				}
				// Verify error type is AppError
				if _, ok := err.(*apperrors.AppError); !ok {
					t.Errorf("expected error of type *apperrors.AppError, got %T", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result == nil {
					t.Fatalf("expected non-nil kelas result, got nil")
				}
			}
		})
	}
}

// ==============================================================================
// 2. UNIT TEST: PesertaKelasService.AmbilKelas (KRS) using EP and BVA
// ==============================================================================

func TestKelasService_AmbilKelas_EP_BVA(t *testing.T) {
	ctx := context.Background()

	defaultProdiID := "prodi-ti-1"
	otherProdiID := "prodi-dkv-2"

	tests := []struct {
		name        string
		userID      string
		pengajuanID string
		mockSetup   func(repo *mockKelasRepo, userProv *mockUserProvider)
		wantErr     bool
		errContains string
	}{
		// --------------------------------------------------------------------------
		// BVA: Batas Kuota Kapasitas Kelas (count >= Capacity)
		// Kelas Capacity = 30
		// --------------------------------------------------------------------------
		{
			name: "BVA - Current count 29, Capacity 30 [Just below capacity boundary - Valid]",
			mockSetup: func(repo *mockKelasRepo, userProv *mockUserProvider) {
				repo.LockPengajuanByIDFn = func(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
					return &domain.PengajuanKelas{
						ID:     id,
						Status: domain.StatusApproved,
						Kelas:  &domain.Kelas{ID: "kelas-1", Capacity: 30, ProgramStudiID: defaultProdiID},
					}, nil
				}
				repo.CountPesertaKelasFn = func(ctx context.Context, pengajuanID string) (int64, error) {
					return 29, nil // 1 seat available
				}
			},
			wantErr: false,
		},
		{
			name: "BVA - Current count 30, Capacity 30 [Exact capacity boundary, Class Full - Invalid]",
			mockSetup: func(repo *mockKelasRepo, userProv *mockUserProvider) {
				repo.LockPengajuanByIDFn = func(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
					return &domain.PengajuanKelas{
						ID:     id,
						Status: domain.StatusApproved,
						Kelas:  &domain.Kelas{ID: "kelas-1", Capacity: 30, ProgramStudiID: defaultProdiID},
					}, nil
				}
				repo.CountPesertaKelasFn = func(ctx context.Context, pengajuanID string) (int64, error) {
					return 30, nil // Exact full
				}
			},
			wantErr:     true,
			errContains: "Kelas sudah penuh",
		},
		{
			name: "BVA - Current count 31, Capacity 30 [Over capacity boundary - Invalid]",
			mockSetup: func(repo *mockKelasRepo, userProv *mockUserProvider) {
				repo.LockPengajuanByIDFn = func(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
					return &domain.PengajuanKelas{
						ID:     id,
						Status: domain.StatusApproved,
						Kelas:  &domain.Kelas{ID: "kelas-1", Capacity: 30, ProgramStudiID: defaultProdiID},
					}, nil
				}
				repo.CountPesertaKelasFn = func(ctx context.Context, pengajuanID string) (int64, error) {
					return 31, nil
				}
			},
			wantErr:     true,
			errContains: "Kelas sudah penuh",
		},

		// --------------------------------------------------------------------------
		// EP: Status Pengajuan Kelas
		// --------------------------------------------------------------------------
		{
			name: "EP - Class status 'pending' [Not yet approved partition - Invalid]",
			mockSetup: func(repo *mockKelasRepo, userProv *mockUserProvider) {
				repo.LockPengajuanByIDFn = func(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
					return &domain.PengajuanKelas{
						ID:     id,
						Status: domain.StatusPending,
						Kelas:  &domain.Kelas{ID: "kelas-1", Capacity: 30, ProgramStudiID: defaultProdiID},
					}, nil
				}
			},
			wantErr:     true,
			errContains: "Kelas belum disetujui",
		},
		{
			name: "EP - Class status 'rejected' [Rejected partition - Invalid]",
			mockSetup: func(repo *mockKelasRepo, userProv *mockUserProvider) {
				repo.LockPengajuanByIDFn = func(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
					return &domain.PengajuanKelas{
						ID:     id,
						Status: domain.StatusRejected,
						Kelas:  &domain.Kelas{ID: "kelas-1", Capacity: 30, ProgramStudiID: defaultProdiID},
					}, nil
				}
			},
			wantErr:     true,
			errContains: "Kelas belum disetujui",
		},

		// --------------------------------------------------------------------------
		// EP: Kesesuaian Program Studi
		// --------------------------------------------------------------------------
		{
			name: "EP - Student Prodi differs from Class Prodi [Cross-prodi partition - Invalid]",
			mockSetup: func(repo *mockKelasRepo, userProv *mockUserProvider) {
				repo.LockPengajuanByIDFn = func(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
					return &domain.PengajuanKelas{
						ID:     id,
						Status: domain.StatusApproved,
						Kelas:  &domain.Kelas{ID: "kelas-1", Capacity: 30, ProgramStudiID: defaultProdiID},
					}, nil
				}
				userProv.GetUserByIDFn = func(ctx context.Context, userID string) (*authDomain.User, error) {
					return &authDomain.User{ID: userID, ProgramStudiID: &otherProdiID}, nil
				}
			},
			wantErr:     true,
			errContains: "Kelas ini tidak tersedia untuk Program Studi Anda",
		},
		{
			name: "EP - Student has no Prodi [Unassigned prodi partition - Invalid]",
			mockSetup: func(repo *mockKelasRepo, userProv *mockUserProvider) {
				repo.LockPengajuanByIDFn = func(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
					return &domain.PengajuanKelas{
						ID:     id,
						Status: domain.StatusApproved,
						Kelas:  &domain.Kelas{ID: "kelas-1", Capacity: 30, ProgramStudiID: defaultProdiID},
					}, nil
				}
				userProv.GetUserByIDFn = func(ctx context.Context, userID string) (*authDomain.User, error) {
					return &authDomain.User{ID: userID, ProgramStudiID: nil}, nil
				}
			},
			wantErr:     true,
			errContains: "Kelas ini tidak tersedia untuk Program Studi Anda",
		},

		// --------------------------------------------------------------------------
		// EP: Duplikasi dan Bentrok (Conflict Partitions)
		// --------------------------------------------------------------------------
		{
			name: "EP - Student already enrolled in this class [Duplicate enrollment partition - Invalid]",
			pengajuanID: "pengajuan-1",
			mockSetup: func(repo *mockKelasRepo, userProv *mockUserProvider) {
				repo.LockPengajuanByIDFn = func(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
					return &domain.PengajuanKelas{
						ID:     id,
						Status: domain.StatusApproved,
						Kelas:  &domain.Kelas{ID: "kelas-1", Capacity: 30, ProgramStudiID: defaultProdiID},
					}, nil
				}
				repo.CountPesertaKelasFn = func(ctx context.Context, pengajuanID string) (int64, error) {
					return 10, nil
				}
				repo.GetPesertaKelasByMahasiswaIDFn = func(ctx context.Context, mahasiswaID string) ([]*domain.PesertaKelas, error) {
					return []*domain.PesertaKelas{
						{PengajuanID: "pengajuan-1", MahasiswaID: mahasiswaID},
					}, nil
				}
			},
			wantErr:     true,
			errContains: "Anda sudah mengambil kelas ini",
		},
		{
			name: "EP - Course conflict: Same Mata Kuliah already taken [Duplicate MK partition - Invalid]",
			pengajuanID: "pengajuan-1",
			mockSetup: func(repo *mockKelasRepo, userProv *mockUserProvider) {
				repo.LockPengajuanByIDFn = func(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
					return &domain.PengajuanKelas{
						ID:           id,
						MataKuliahID: "mk-algo-1",
						Status:       domain.StatusApproved,
						Kelas:        &domain.Kelas{ID: "kelas-1", Capacity: 30, ProgramStudiID: defaultProdiID},
					}, nil
				}
				repo.CountPesertaKelasFn = func(ctx context.Context, pengajuanID string) (int64, error) {
					return 10, nil
				}
				repo.CheckPesertaMataKuliahConflictFn = func(ctx context.Context, mahasiswaID string, mkID string) (bool, error) {
					return true, nil // MK conflict exists
				}
			},
			wantErr:     true,
			errContains: "Anda sudah mengambil kelas lain untuk Mata Kuliah ini",
		},
		{
			name: "EP - Schedule conflict: Overlapping time with another class [Schedule conflict partition - Invalid]",
			pengajuanID: "pengajuan-1",
			mockSetup: func(repo *mockKelasRepo, userProv *mockUserProvider) {
				repo.LockPengajuanByIDFn = func(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
					return &domain.PengajuanKelas{
						ID:           id,
						MataKuliahID: "mk-algo-1",
						Status:       domain.StatusApproved,
						Kelas: &domain.Kelas{
							ID: "kelas-1", Capacity: 30, ProgramStudiID: defaultProdiID,
							Hari: "Senin", JamMulai: "08:00", JamSelesai: "10:30",
						},
					}, nil
				}
				repo.CountPesertaKelasFn = func(ctx context.Context, pengajuanID string) (int64, error) {
					return 10, nil
				}
				repo.CheckPesertaScheduleConflictFn = func(ctx context.Context, mahasiswaID string, hari string, jamMulai string, jamSelesai string) (bool, error) {
					return true, nil // Schedule overlaps
				}
			},
			wantErr:     true,
			errContains: "Jadwal kelas ini bentrok dengan kelas Anda yang lain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockKelasRepo{}
			mockUser := &mockUserProvider{}

			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo, mockUser)
			}

			svc := service.NewKelasService(mockRepo, mockUser, nil, nil)

			userID := tt.userID
			if userID == "" {
				userID = "mhs-1"
			}
			pengajuanID := tt.pengajuanID
			if pengajuanID == "" {
				pengajuanID = "pengajuan-1"
			}

			err := svc.AmbilKelas(ctx, userID, pengajuanID)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, but got nil", tt.errContains)
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain expected substring %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// ==============================================================================
// 3. UNIT TEST: Core Read & Delete Operations
// ==============================================================================

func TestKelasService_GetByID_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("GetByID - Not Found Partition", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetByIDFn: func(ctx context.Context, id string) (*domain.Kelas, error) {
				return nil, errors.New("kelas not found")
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		kelas, err := svc.GetByID(ctx, "non-existent-id")
		if err == nil {
			t.Fatalf("expected not found error, got nil")
		}
		if kelas != nil {
			t.Fatalf("expected nil kelas, got %v", kelas)
		}
	})

	t.Run("GetByID - Success Partition", func(t *testing.T) {
		expectedKelas := &domain.Kelas{ID: "kelas-123", Name: "IF-101", Capacity: 30}
		mockRepo := &mockKelasRepo{
			GetByIDFn: func(ctx context.Context, id string) (*domain.Kelas, error) {
				return expectedKelas, nil
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		kelas, err := svc.GetByID(ctx, "kelas-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if kelas.ID != expectedKelas.ID || kelas.Name != expectedKelas.Name {
			t.Errorf("expected %+v, got %+v", expectedKelas, kelas)
		}
	})

	t.Run("Delete - Not Found Partition", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetByIDFn: func(ctx context.Context, id string) (*domain.Kelas, error) {
				return nil, errors.New("kelas not found")
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		err := svc.Delete(ctx, "non-existent-id")
		if err == nil {
			t.Fatalf("expected error when deleting non-existent kelas, got nil")
		}
	})

	t.Run("Delete - Success Partition", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetByIDFn: func(ctx context.Context, id string) (*domain.Kelas, error) {
				return &domain.Kelas{ID: "kelas-123"}, nil
			},
			DeleteFn: func(ctx context.Context, id string) error {
				return nil
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		err := svc.Delete(ctx, "kelas-123")
		if err != nil {
			t.Fatalf("unexpected error deleting kelas: %v", err)
		}
	})
}

// ==============================================================================
// 4. UNIT TEST: PengajuanKelasService (RequestKelas, Approve, Reject)
// ==============================================================================

func TestKelasService_PengajuanKelas_EP(t *testing.T) {
	ctx := context.Background()

	t.Run("RequestKelas - Kelas Not Found", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetByIDFn: func(ctx context.Context, id string) (*domain.Kelas, error) {
				return nil, errors.New("not found")
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		_, err := svc.RequestKelas(ctx, "dosen-1", domain.RequestKelasPayload{KelasID: "k-none", MataKuliahID: "mk-1"})
		if err == nil || !strings.Contains(err.Error(), "Kelas tidak ditemukan") {
			t.Fatalf("expected Kelas tidak ditemukan, got: %v", err)
		}
	})

	t.Run("RequestKelas - Invalid Mata Kuliah for Dosen/Prodi", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetByIDFn: func(ctx context.Context, id string) (*domain.Kelas, error) {
				return &domain.Kelas{ID: id, ProgramStudiID: "prodi-1"}, nil
			},
		}
		mockMK := &mockMataKuliahProvider{
			IsMataKuliahValidForKelasFn: func(ctx context.Context, dosenID, mkID, prodiID string) (bool, error) {
				return false, nil
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, mockMK)

		_, err := svc.RequestKelas(ctx, "dosen-1", domain.RequestKelasPayload{KelasID: "k-1", MataKuliahID: "mk-1"})
		if err == nil || !strings.Contains(err.Error(), "Mata kuliah yang dipilih tidak valid") {
			t.Fatalf("expected Mata kuliah yang dipilih tidak valid error, got: %v", err)
		}
	})

	t.Run("RequestKelas - Class Already Approved for Another Dosen", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetByIDFn: func(ctx context.Context, id string) (*domain.Kelas, error) {
				return &domain.Kelas{ID: id, ProgramStudiID: "prodi-1"}, nil
			},
			GetActivePengajuanByKelasIDFn: func(ctx context.Context, kelasID string) ([]*domain.PengajuanKelas, error) {
				return []*domain.PengajuanKelas{
					{ID: "p-other", DosenID: "dosen-other", Status: domain.StatusApproved},
				}, nil
			},
		}
		mockMK := &mockMataKuliahProvider{}
		svc := service.NewKelasService(mockRepo, nil, nil, mockMK)

		_, err := svc.RequestKelas(ctx, "dosen-1", domain.RequestKelasPayload{KelasID: "k-1", MataKuliahID: "mk-1"})
		if err == nil || !strings.Contains(err.Error(), "Kelas ini sudah disetujui untuk dosen lain") {
			t.Fatalf("expected Kelas ini sudah disetujui untuk dosen lain, got: %v", err)
		}
	})

	t.Run("RequestKelas - Dosen Already Has Pending Proposal for This Class", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetByIDFn: func(ctx context.Context, id string) (*domain.Kelas, error) {
				return &domain.Kelas{ID: id, ProgramStudiID: "prodi-1"}, nil
			},
			GetActivePengajuanByKelasIDFn: func(ctx context.Context, kelasID string) ([]*domain.PengajuanKelas, error) {
				return []*domain.PengajuanKelas{
					{ID: "p-1", DosenID: "dosen-1", Status: domain.StatusPending},
				}, nil
			},
		}
		mockMK := &mockMataKuliahProvider{}
		svc := service.NewKelasService(mockRepo, nil, nil, mockMK)

		_, err := svc.RequestKelas(ctx, "dosen-1", domain.RequestKelasPayload{KelasID: "k-1", MataKuliahID: "mk-1"})
		if err == nil || !strings.Contains(err.Error(), "Anda sudah mengajukan kelas ini") {
			t.Fatalf("expected Anda sudah mengajukan kelas ini, got: %v", err)
		}
	})

	t.Run("RequestKelas - No Active Academic Period", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetByIDFn: func(ctx context.Context, id string) (*domain.Kelas, error) {
				return &domain.Kelas{ID: id, ProgramStudiID: "prodi-1"}, nil
			},
		}
		mockMK := &mockMataKuliahProvider{}
		mockPeriode := &mockPeriodeProvider{
			GetActivePeriodeIDFn: func(ctx context.Context) (string, error) {
				return "", errors.New("tidak ada periode aktif")
			},
		}
		svc := service.NewKelasService(mockRepo, nil, mockPeriode, mockMK)

		_, err := svc.RequestKelas(ctx, "dosen-1", domain.RequestKelasPayload{KelasID: "k-1", MataKuliahID: "mk-1"})
		if err == nil || !strings.Contains(err.Error(), "Tidak ada periode akademik yang aktif") {
			t.Fatalf("expected Tidak ada periode akademik yang aktif, got: %v", err)
		}
	})

	t.Run("RequestKelas - Success Partition", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetByIDFn: func(ctx context.Context, id string) (*domain.Kelas, error) {
				return &domain.Kelas{ID: id, ProgramStudiID: "prodi-1"}, nil
			},
			CreatePengajuanFn: func(ctx context.Context, p *domain.PengajuanKelas) error {
				return nil
			},
		}
		mockMK := &mockMataKuliahProvider{}
		mockPeriode := &mockPeriodeProvider{}
		svc := service.NewKelasService(mockRepo, nil, mockPeriode, mockMK)

		result, err := svc.RequestKelas(ctx, "dosen-1", domain.RequestKelasPayload{KelasID: "k-1", MataKuliahID: "mk-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Status != domain.StatusPending {
			t.Errorf("expected pending status, got %s", result.Status)
		}
	})

	t.Run("ApprovePengajuan - Non-pending status rejected", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetPengajuanByIDFn: func(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
				return &domain.PengajuanKelas{ID: id, Status: domain.StatusApproved}, nil
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		err := svc.ApprovePengajuan(ctx, "p-1")
		if err == nil || !strings.Contains(err.Error(), "Hanya pengajuan berstatus pending") {
			t.Fatalf("expected error for approving non-pending proposal, got: %v", err)
		}
	})

	t.Run("ApprovePengajuan - Success automatically rejects other pending proposals", func(t *testing.T) {
		rivalProposal := &domain.PengajuanKelas{ID: "p-rival", KelasID: "k-1", Status: domain.StatusPending}
		targetProposal := &domain.PengajuanKelas{ID: "p-target", KelasID: "k-1", Status: domain.StatusPending}

		mockRepo := &mockKelasRepo{
			GetPengajuanByIDFn: func(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
				return targetProposal, nil
			},
			GetActivePengajuanByKelasIDFn: func(ctx context.Context, kelasID string) ([]*domain.PengajuanKelas, error) {
				return []*domain.PengajuanKelas{targetProposal, rivalProposal}, nil
			},
			UpdatePengajuanFn: func(ctx context.Context, p *domain.PengajuanKelas) error {
				return nil
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		err := svc.ApprovePengajuan(ctx, "p-target")
		if err != nil {
			t.Fatalf("unexpected error approving proposal: %v", err)
		}
		if targetProposal.Status != domain.StatusApproved {
			t.Errorf("expected target proposal to be approved, got %s", targetProposal.Status)
		}
		if rivalProposal.Status != domain.StatusRejected {
			t.Errorf("expected rival proposal to be rejected, got %s", rivalProposal.Status)
		}
	})

	t.Run("RejectPengajuan - Success", func(t *testing.T) {
		targetProposal := &domain.PengajuanKelas{ID: "p-target", Status: domain.StatusPending}

		mockRepo := &mockKelasRepo{
			GetPengajuanByIDFn: func(ctx context.Context, id string) (*domain.PengajuanKelas, error) {
				return targetProposal, nil
			},
			UpdatePengajuanFn: func(ctx context.Context, p *domain.PengajuanKelas) error {
				return nil
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		err := svc.RejectPengajuan(ctx, "p-target")
		if err != nil {
			t.Fatalf("unexpected error rejecting proposal: %v", err)
		}
		if targetProposal.Status != domain.StatusRejected {
			t.Errorf("expected proposal to be rejected, got %s", targetProposal.Status)
		}
	})
}

// ==============================================================================
// 5. UNIT TEST: PertemuanService using BVA and EP
// ==============================================================================

func TestKelasService_Pertemuan_EP_BVA(t *testing.T) {
	ctx := context.Background()

	t.Run("BVA - MulaiPertemuan with 15 existing meetings [Next is 16, Max boundary - Valid]", func(t *testing.T) {
		existing := make([]*domain.Pertemuan, 15)
		mockRepo := &mockKelasRepo{
			GetPertemuanByPengajuanIDFn: func(ctx context.Context, pengajuanID string) ([]*domain.Pertemuan, error) {
				return existing, nil
			},
			CreatePertemuanFn: func(ctx context.Context, p *domain.Pertemuan) error {
				return nil
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		p, err := svc.MulaiPertemuan(ctx, "p-1", "Pertemuan 16")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.NomorPertemuan != 16 {
			t.Errorf("expected meeting number 16, got %d", p.NomorPertemuan)
		}
	})

	t.Run("BVA - MulaiPertemuan with 16 existing meetings [Next is 17, Over max boundary - Invalid]", func(t *testing.T) {
		existing := make([]*domain.Pertemuan, 16)
		mockRepo := &mockKelasRepo{
			GetPertemuanByPengajuanIDFn: func(ctx context.Context, pengajuanID string) ([]*domain.Pertemuan, error) {
				return existing, nil
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		_, err := svc.MulaiPertemuan(ctx, "p-1", "Pertemuan 17")
		if err == nil || !strings.Contains(err.Error(), "Semua 16 pertemuan sudah tercatat") {
			t.Fatalf("expected error for exceeding 16 meetings, got: %v", err)
		}
	})

	t.Run("EP - AkhiriPertemuan: Meeting Not Found", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetPertemuanByIDFn: func(ctx context.Context, id string) (*domain.Pertemuan, error) {
				return nil, errors.New("not found")
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		err := svc.AkhiriPertemuan(ctx, "pt-none")
		if err == nil || !strings.Contains(err.Error(), "Pertemuan tidak ditemukan") {
			t.Fatalf("expected Pertemuan tidak ditemukan, got: %v", err)
		}
	})

	t.Run("EP - AkhiriPertemuan: Meeting Already Finished", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetPertemuanByIDFn: func(ctx context.Context, id string) (*domain.Pertemuan, error) {
				return &domain.Pertemuan{ID: id, Status: domain.PertemuanStatusSelesai}, nil
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		err := svc.AkhiriPertemuan(ctx, "pt-1")
		if err == nil || !strings.Contains(err.Error(), "Pertemuan sudah selesai atau tidak aktif") {
			t.Fatalf("expected error for already finished meeting, got: %v", err)
		}
	})

	t.Run("EP - AkhiriPertemuan: Success", func(t *testing.T) {
		meeting := &domain.Pertemuan{ID: "pt-1", Status: domain.PertemuanStatusBerlangsung}
		mockRepo := &mockKelasRepo{
			GetPertemuanByIDFn: func(ctx context.Context, id string) (*domain.Pertemuan, error) {
				return meeting, nil
			},
			UpdatePertemuanFn: func(ctx context.Context, p *domain.Pertemuan) error {
				return nil
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		err := svc.AkhiriPertemuan(ctx, "pt-1")
		if err != nil {
			t.Fatalf("unexpected error ending meeting: %v", err)
		}
		if meeting.Status != domain.PertemuanStatusSelesai {
			t.Errorf("expected meeting status selesai, got %s", meeting.Status)
		}
	})
}

// ==============================================================================
// 6. UNIT TEST: AbsensiService using EP
// ==============================================================================

func TestKelasService_Absensi_EP(t *testing.T) {
	ctx := context.Background()

	t.Run("SubmitAbsensiMahasiswa - Pertemuan Not Found", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetPertemuanByIDFn: func(ctx context.Context, id string) (*domain.Pertemuan, error) {
				return nil, errors.New("not found")
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		err := svc.SubmitAbsensiMahasiswa(ctx, "pt-none", "mhs-1", "123456")
		if err == nil || !strings.Contains(err.Error(), "Pertemuan tidak ditemukan") {
			t.Fatalf("expected Pertemuan tidak ditemukan, got: %v", err)
		}
	})

	t.Run("SubmitAbsensiMahasiswa - Pertemuan Inactive", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetPertemuanByIDFn: func(ctx context.Context, id string) (*domain.Pertemuan, error) {
				return &domain.Pertemuan{ID: id, Status: domain.PertemuanStatusSelesai}, nil
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		err := svc.SubmitAbsensiMahasiswa(ctx, "pt-1", "mhs-1", "123456")
		if err == nil || !strings.Contains(err.Error(), "Pertemuan sudah tidak aktif") {
			t.Fatalf("expected Pertemuan sudah tidak aktif, got: %v", err)
		}
	})

	t.Run("SubmitAbsensiMahasiswa - Invalid Attendance Code", func(t *testing.T) {
		mockRepo := &mockKelasRepo{
			GetPertemuanByIDFn: func(ctx context.Context, id string) (*domain.Pertemuan, error) {
				return &domain.Pertemuan{
					ID:          id,
					Status:      domain.PertemuanStatusBerlangsung,
					KodeAbsensi: "654321",
				}, nil
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		err := svc.SubmitAbsensiMahasiswa(ctx, "pt-1", "mhs-1", "WRONG!")
		if err == nil || !strings.Contains(err.Error(), "Kode absensi tidak valid") {
			t.Fatalf("expected Kode absensi tidak valid, got: %v", err)
		}
	})

	t.Run("SubmitAbsensiMahasiswa - Success (New Attendance Record)", func(t *testing.T) {
		created := false
		mockRepo := &mockKelasRepo{
			GetPertemuanByIDFn: func(ctx context.Context, id string) (*domain.Pertemuan, error) {
				return &domain.Pertemuan{
					ID:          id,
					Status:      domain.PertemuanStatusBerlangsung,
					KodeAbsensi: "123456",
				}, nil
			},
			GetAbsensiByPertemuanIDFn: func(ctx context.Context, pertemuanID string) ([]*domain.Absensi, error) {
				return nil, nil // no prior record
			},
			CreateAbsensiFn: func(ctx context.Context, a *domain.Absensi) error {
				created = true
				if a.StatusKehadiran != "hadir" {
					t.Errorf("expected hadir status, got %s", a.StatusKehadiran)
				}
				return nil
			},
		}
		svc := service.NewKelasService(mockRepo, nil, nil, nil)

		err := svc.SubmitAbsensiMahasiswa(ctx, "pt-1", "mhs-1", "123456")
		if err != nil {
			t.Fatalf("unexpected error submitting attendance: %v", err)
		}
		if !created {
			t.Errorf("expected CreateAbsensi to be called")
		}
	})
}
