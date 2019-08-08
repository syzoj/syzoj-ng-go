package config

import (
	"database/sql"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func OpenMySQL(dbName string) (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv(dbName + "_MYSQL_HOST")
	cfg.User = os.Getenv(dbName + "_MYSQL_USER")
	cfg.Passwd = os.Getenv(dbName + "_MYSQL_PASSWORD")
	cfg.DBName = os.Getenv(dbName + "_MYSQL_DATABASE")
	cfg.ClientFoundRows = true // Default to return matched rows instead of changed rows in an UPDATE query since that makes more sense
	return sql.Open("mysql", cfg.FormatDSN())
}

func OpenMySQLGorm(dbName string) (*gorm.DB, error) {
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv(dbName + "_MYSQL_ADDR")
	cfg.User = os.Getenv(dbName + "_MYSQL_USER")
	cfg.Passwd = os.Getenv(dbName + "_MYSQL_PASSWORD")
	cfg.DBName = os.Getenv(dbName + "_MYSQL_DATABASE")
	cfg.ClientFoundRows = true // Default to return matched rows instead of changed rows in an UPDATE query since that makes more sense
	return gorm.Open("mysql", cfg.FormatDSN())
}
