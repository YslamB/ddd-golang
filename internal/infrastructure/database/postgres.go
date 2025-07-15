package database

import (
	"database/sql"
	"ddd/configs"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgresDB(config configs.DatabaseConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	return db, err
}

// func createTables(db *sql.DB) error {
// 	createUsersTable := `
//     CREATE TABLE IF NOT EXISTS users (
//         id UUID PRIMARY KEY,
//         name VARCHAR(100) NOT NULL,
//         email VARCHAR(255) NOT NULL UNIQUE,
//         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//     );
//     CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
//     CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
//     `

// 	_, err := db.Exec(createUsersTable)
// 	return err
// }
