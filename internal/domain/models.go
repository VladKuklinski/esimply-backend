package domain

import "errors"

var ErrNotFound = errors.New("not found")

type Country struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Flag     string `json:"flag"`
	SailyURL string `json:"saily_url"`
}

type Plan struct {
	ID           string  `json:"id"`
	DataGB       int     `json:"data_gb"`
	ValidityDays int     `json:"validity_days"`
	PriceEUR     float64 `json:"price_eur"`
	BestValue    bool    `json:"best_value"`
	Description  string  `json:"description"`
}

type CountryRepository interface {
	GetAll() ([]Country, error)
	GetPlansByCountryID(id string) ([]Plan, error)
}

type CountryUsecase interface {
	GetAllCountries() ([]Country, error)
	GetPlansByCountryID(id string) ([]Plan, error)
}
