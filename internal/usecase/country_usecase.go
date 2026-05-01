package usecase

import "esimply/internal/domain"

type countryUsecase struct {
	repo domain.CountryRepository
}

func NewCountryUsecase(repo domain.CountryRepository) domain.CountryUsecase {
	return &countryUsecase{repo: repo}
}

func (u *countryUsecase) GetAllCountries() ([]domain.Country, error) {
	return u.repo.GetAll()
}

func (u *countryUsecase) GetPlansByCountryID(id string) ([]domain.Plan, error) {
	return u.repo.GetPlansByCountryID(id)
}
