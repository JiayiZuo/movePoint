package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(connectionString string) error {
	if err := EnsureDatabase(connectionString); err != nil {
		return err
	}

	var err error
	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := RunMigrations(DB); err != nil {
		return err
	}

	log.Println("Database connection established and models migrated")
	return nil
}

func EnsureDatabase(connectionString string) error {
	cfg, err := mysqlDriver.ParseDSN(connectionString)
	if err != nil {
		return fmt.Errorf("failed to parse database dsn: %v", err)
	}

	dbName := cfg.DBName
	if dbName == "" {
		return fmt.Errorf("database name is required in DB_DSN")
	}

	cfg.DBName = ""
	serverDSN := cfg.FormatDSN()

	db, err := sql.Open("mysql", serverDSN)
	if err != nil {
		return fmt.Errorf("failed to open mysql server connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to mysql server: %v", err)
	}

	query := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		escapeIdentifier(dbName),
	)
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create database %q: %v", dbName, err)
	}

	log.Printf("Database %q is ready", dbName)
	return nil
}

func escapeIdentifier(value string) string {
	return strings.ReplaceAll(value, "`", "``")
}
