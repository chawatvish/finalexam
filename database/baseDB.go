package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	return sql.Open("postgres", os.Getenv("DATABASE_URL"))
}
