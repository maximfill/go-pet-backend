package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func New(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	var pingErr error

	for i := 0; i < 10; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			return db, nil
		}
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("database not ready after retries: %w", pingErr)
}
