package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func RunMigrations(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS countries (
			id        VARCHAR(50)  PRIMARY KEY,
			name      VARCHAR(100) NOT NULL,
			flag      VARCHAR(10)  NOT NULL,
			saily_url VARCHAR(255) NOT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS plans (
			id            VARCHAR(100)  PRIMARY KEY,
			country_id    VARCHAR(50)   NOT NULL,
			data_gb       INT           NOT NULL,
			validity_days INT           NOT NULL,
			price_eur     DECIMAL(10,2) NOT NULL,
			best_value    TINYINT(1)    NOT NULL DEFAULT 0,
			description   VARCHAR(255)  NOT NULL,
			FOREIGN KEY (country_id) REFERENCES countries(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

func SeedIfEmpty(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM countries").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	type country struct {
		id, name, flag, sailyURL string
	}
	type plan struct {
		suffix, description  string
		dataGB, validityDays int
		priceEUR             float64
		bestValue            bool
	}

	countries := []country{
		{"france", "France", "🇫🇷", "https://saily.com/esim-france/"},
		{"italy", "Italy", "🇮🇹", "https://saily.com/esim-italy/"},
		{"spain", "Spain", "🇪🇸", "https://saily.com/esim-spain/"},
		{"germany", "Germany", "🇩🇪", "https://saily.com/esim-germany/"},
		{"greece", "Greece", "🇬🇷", "https://saily.com/esim-greece/"},
		{"portugal", "Portugal", "🇵🇹", "https://saily.com/esim-portugal/"},
		{"netherlands", "Netherlands", "🇳🇱", "https://saily.com/esim-netherlands/"},
		{"poland", "Poland", "🇵🇱", "https://saily.com/esim-poland/"},
		{"croatia", "Croatia", "🇭🇷", "https://saily.com/esim-croatia/"},
		{"czech-republic", "Czech Republic", "🇨🇿", "https://saily.com/esim-czech-republic/"},
		{"austria", "Austria", "🇦🇹", "https://saily.com/esim-austria/"},
		{"switzerland", "Switzerland", "🇨🇭", "https://saily.com/esim-switzerland/"},
		{"hungary", "Hungary", "🇭🇺", "https://saily.com/esim-hungary/"},
		{"sweden", "Sweden", "🇸🇪", "https://saily.com/esim-sweden/"},
		{"norway", "Norway", "🇳🇴", "https://saily.com/esim-norway/"},
	}
	plans := []plan{
		{"1gb", "Perfect for quick trips", 1, 7, 2.99, false},
		{"5gb", "Great for most travelers", 5, 14, 7.99, true},
		{"10gb", "Long stays and heavy use", 10, 30, 13.99, false},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, c := range countries {
		if _, err := tx.Exec(
			"INSERT INTO countries (id, name, flag, saily_url) VALUES (?, ?, ?, ?)",
			c.id, c.name, c.flag, c.sailyURL,
		); err != nil {
			return err
		}
		for _, p := range plans {
			if _, err := tx.Exec(
				"INSERT INTO plans (id, country_id, data_gb, validity_days, price_eur, best_value, description) VALUES (?, ?, ?, ?, ?, ?, ?)",
				c.id+"-"+p.suffix, c.id, p.dataGB, p.validityDays, p.priceEUR, p.bestValue, p.description,
			); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
