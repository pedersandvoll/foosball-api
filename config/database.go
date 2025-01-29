package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

type Config struct {
	Host      string
	Port      string
	User      string
	Password  string
	DBName    string
	SSLMode   string
	JWTSecret string
}

func NewConfig() *Config {
	return &Config{
		Host:      getEnv("DB_HOST", "localhost"),
		Port:      getEnv("DB_PORT", "5432"),
		User:      getEnv("DB_USER", "postgres"),
		Password:  getEnv("DB_PASSWORD", "password"),
		DBName:    getEnv("DB_NAME", "dbname"),
		SSLMode:   getEnv("DB_SSLMODE", "disable"),
		JWTSecret: getEnv("JWT_SECRET", "your-default-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func NewDatabase(config *Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return &Database{db}, nil
}
