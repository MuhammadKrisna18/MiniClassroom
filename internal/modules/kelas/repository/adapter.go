package repository

import (
	"context"
	"errors"

	authDomain "siakad-pro/internal/modules/auth/domain"
	"siakad-pro/internal/modules/kelas/domain"
	mkDomain "siakad-pro/internal/modules/matakuliah/domain"
	periodeDomain "siakad-pro/internal/modules/periode/domain"
)

type userProviderAdapter struct {
	authRepo authDomain.AuthRepository
}

func NewUserProviderAdapter(authRepo authDomain.AuthRepository) domain.UserProvider {
	return &userProviderAdapter{authRepo: authRepo}
}

func (a *userProviderAdapter) GetUserByID(ctx context.Context, userID string) (*authDomain.User, error) {
	return a.authRepo.GetByID(ctx, userID)
}

type periodeProviderAdapter struct {
	periodeRepo periodeDomain.PeriodeRepository
}

func NewPeriodeProviderAdapter(periodeRepo periodeDomain.PeriodeRepository) domain.PeriodeProvider {
	return &periodeProviderAdapter{periodeRepo: periodeRepo}
}

func (a *periodeProviderAdapter) GetActivePeriodeID(ctx context.Context) (string, error) {
	p, err := a.periodeRepo.GetActive(ctx)
	if err != nil {
		return "", err
	}
	if p == nil {
		return "", errors.New("tidak ada periode aktif")
	}
	return p.ID, nil
}

type mataKuliahProviderAdapter struct {
	mkRepo mkDomain.MataKuliahRepository
}

func NewMataKuliahProviderAdapter(mkRepo mkDomain.MataKuliahRepository) domain.MataKuliahProvider {
	return &mataKuliahProviderAdapter{mkRepo: mkRepo}
}

func (a *mataKuliahProviderAdapter) IsMataKuliahValidForKelas(ctx context.Context, dosenID string, mkID string, prodiID string) (bool, error) {
	return a.mkRepo.IsMataKuliahValidForKelas(ctx, dosenID, mkID, prodiID)
}
