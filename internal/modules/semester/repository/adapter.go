package repository

import (
	"context"

	mkDomain "siakad-pro/internal/modules/matakuliah/domain"
	"siakad-pro/internal/modules/semester/domain"
)

type mataKuliahProviderAdapter struct {
	mkRepo mkDomain.MataKuliahCatalogRepository
}

func NewMataKuliahProviderAdapter(mkRepo mkDomain.MataKuliahCatalogRepository) domain.MataKuliahProvider {
	return &mataKuliahProviderAdapter{mkRepo: mkRepo}
}

func (a *mataKuliahProviderAdapter) GetMataKuliahByID(ctx context.Context, id string) (*domain.MataKuliahInfo, error) {
	mk, err := a.mkRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if mk == nil {
		return nil, nil
	}
	return mapToMataKuliahInfo(mk), nil
}

func (a *mataKuliahProviderAdapter) GetMataKuliahByIDs(ctx context.Context, ids []string) ([]*domain.MataKuliahInfo, error) {
	if len(ids) == 0 {
		return []*domain.MataKuliahInfo{}, nil
	}
	mks, err := a.mkRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	var res []*domain.MataKuliahInfo
	for _, mk := range mks {
		res = append(res, mapToMataKuliahInfo(mk))
	}
	return res, nil
}

func mapToMataKuliahInfo(mk *mkDomain.MataKuliah) *domain.MataKuliahInfo {
	return &domain.MataKuliahInfo{
		ID:             mk.ID,
		Name:           mk.Name,
		SKS:            mk.SKS,
		ProgramStudiID: mk.ProgramStudiID,
	}
}
