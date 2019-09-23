package config

import (
	"database/sql"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Creates a *sql.DB instance from environment variables. The environment variables are:
// ${prefix}MYSQL_HOST, ${prefix}MYSQL_USER, ${prefix}MYSQL_PASSWORD, ${prefix}MYSQL_DATABASE.
func NewMySQL(prefix string) (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv(prefix + "MYSQL_HOST")
	cfg.User = os.Getenv(prefix + "MYSQL_USER")
	cfg.Passwd = os.Getenv(prefix + "MYSQL_PASSWORD")
	cfg.DBName = os.Getenv(prefix + "MYSQL_DATABASE")
	cfg.ClientFoundRows = true // Default to return matched rows instead of changed rows in an UPDATE query since that makes more sense
	cfg.ParseTime = true
	return sql.Open("mysql", cfg.FormatDSN())
}

// Creates a *sqlx.DB instance from environment variables. The environment variables are:
// ${prefix}MYSQL_HOST, ${prefix}MYSQL_USER, ${prefix}MYSQL_PASSWORD, ${prefix}MYSQL_DATABASE.
func NewMySQLx(prefix string) (*sqlx.DB, error) {
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv(prefix + "MYSQL_HOST")
	cfg.User = os.Getenv(prefix + "MYSQL_USER")
	cfg.Passwd = os.Getenv(prefix + "MYSQL_PASSWORD")
	cfg.DBName = os.Getenv(prefix + "MYSQL_DATABASE")
	cfg.ClientFoundRows = true // Default to return matched rows instead of changed rows in an UPDATE query since that makes more sense
	cfg.ParseTime = true
	return sqlx.Open("mysql", cfg.FormatDSN())
}
