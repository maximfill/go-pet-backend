package postgres

import "database/sql"

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id),
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT false,
			image_url TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`)
	return err
}
