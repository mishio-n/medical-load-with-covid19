package shared

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")
	DATABASE := os.Getenv("DATABASE")

	return sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@("+DB_HOST+":"+DB_PORT+")/"+DATABASE)
}
