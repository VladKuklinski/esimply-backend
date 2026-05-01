package mysql

import (
	"database/sql"

	"esimply/internal/domain"
)

type countryRepo struct {
	db *sql.DB
}

func NewCountryRepository(db *sql.DB) domain.CountryRepository {
	return &countryRepo{db: db}
}

func (r *countryRepo) GetAll() ([]domain.Country, error) {
	rows, err := r.db.Query("SELECT id, name, flag, saily_url FROM countries ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	countries := make([]domain.Country, 0)
	for rows.Next() {
		var c domain.Country
		if err := rows.Scan(&c.ID, &c.Name, &c.Flag, &c.SailyURL); err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}
	return countries, rows.Err()
}

func (r *countryRepo) GetPlansByCountryID(id string) ([]domain.Plan, error) {
	var exists bool
	if err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM countries WHERE id = ?)", id).Scan(&exists); err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrNotFound
	}

	rows, err := r.db.Query(
		"SELECT id, data_gb, validity_days, price_eur, best_value, description FROM plans WHERE country_id = ? ORDER BY data_gb",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := make([]domain.Plan, 0)
	for rows.Next() {
		var p domain.Plan
		if err := rows.Scan(&p.ID, &p.DataGB, &p.ValidityDays, &p.PriceEUR, &p.BestValue, &p.Description); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, rows.Err()
}
