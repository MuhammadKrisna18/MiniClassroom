package repository

import (
	"context"
	authDomain "siakad-pro/internal/modules/auth/domain"
	mkDomain "siakad-pro/internal/modules/matakuliah/domain"
	periodeDomain "siakad-pro/internal/modules/periode/domain"
)

type userProviderAdapter struct {
	authRepo authDomain.AuthRepository
}

func NewUserProviderAdapter(authRepo authDomain.AuthRepository) mkDomain.UserProvider {
	return &userProviderAdapter{authRepo: authRepo}
}

func (a *userProviderAdapter) GetUserProdiID(ctx context.Context, userID string) (*string, error) {
	u, err := a.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return u.ProgramStudiID, nil
}

func (a *userProviderAdapter) GetDosenIDsByProdi(ctx context.Context, prodiID string) ([]string, error) {
	users, err := a.authRepo.GetUsersByRole(ctx, authDomain.RoleDosen)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, u := range users {
		if u.ProgramStudiID != nil && *u.ProgramStudiID == prodiID {
			ids = append(ids, u.ID)
		}
	}
	return ids, nil
}

type periodeProviderAdapter struct {
	periodeRepo periodeDomain.PeriodeRepository
}

func NewPeriodeProviderAdapter(periodeRepo periodeDomain.PeriodeRepository) mkDomain.PeriodeProvider {
	return &periodeProviderAdapter{periodeRepo: periodeRepo}
}

func (a *periodeProviderAdapter) GetActivePeriodeID(ctx context.Context) (string, error) {
	p, err := a.periodeRepo.GetActive(ctx)
	if err != nil {
		return "", err
	}
	return p.ID, nil
}
