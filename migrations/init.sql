CREATE TABLE IF NOT EXISTS countries (
    id        VARCHAR(50)  PRIMARY KEY,
    name      VARCHAR(100) NOT NULL,
    flag      VARCHAR(10)  NOT NULL,
    saily_url VARCHAR(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS plans (
    id            VARCHAR(100)  PRIMARY KEY,
    country_id    VARCHAR(50)   NOT NULL,
    data_gb       INT           NOT NULL,
    validity_days INT           NOT NULL,
    price_eur     DECIMAL(10,2) NOT NULL,
    best_value    TINYINT(1)    NOT NULL DEFAULT 0,
    description   VARCHAR(255)  NOT NULL,
    FOREIGN KEY (country_id) REFERENCES countries(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
