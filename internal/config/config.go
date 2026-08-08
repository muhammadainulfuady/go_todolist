package config

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
	JWT    JWTConfig
	SMTP   SMTPConfig
	OTP    OTPConfig
}

type ServerConfig struct {
	Port    string
	BaseURL string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
}

type OTPConfig struct {
	ExpiresIn time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET wajib diisi di file .env")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:    envOrDefault("APP_PORT", "8080"),
			BaseURL: envOrDefault("BASE_URL", "http://localhost:8080"),
		},
		DB: DBConfig{
			Host:     envOrDefault("DB_HOST", "127.0.0.1"),
			Port:     envOrDefault("DB_PORT", "3306"),
			User:     envOrDefault("DB_USER", "root"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     envOrDefault("DB_NAME", "go_todolist"),
		},
		SMTP: SMTPConfig{
			Host: envOrDefault("SMTP_HOST", "smtp.gmail.com"),
			Port: envOrDefault("SMTP_PORT", "587"),
			User: os.Getenv("SMTP_USER"),
			Pass: os.Getenv("SMTP_PASS"),
		},
	}

	if cfg.DB.Host == "" || cfg.DB.Name == "" {
		return nil, fmt.Errorf("DB_HOST dan DB_NAME wajib diisi")
	}
	if cfg.SMTP.User == "" || cfg.SMTP.Pass == "" {
		return nil, fmt.Errorf("SMTP_USER dan SMTP_PASS wajib diisi")
	}

	jwtExp, err := time.ParseDuration(envOrDefault("JWT_EXPIRES_IN", "24h"))
	if err != nil {
		return nil, fmt.Errorf("JWT_EXPIRES_IN tidak valid: %w", err)
	}
	cfg.JWT.ExpiresIn = jwtExp
	cfg.JWT.Secret = jwtSecret

	otpExp, err := time.ParseDuration(envOrDefault("OTP_EXPIRES_IN", "5m"))
	if err != nil {
		return nil, fmt.Errorf("OTP_EXPIRES_IN tidak valid: %w", err)
	}
	cfg.OTP.ExpiresIn = otpExp

	return cfg, nil
}

func (c *Config) ConnectDB() (*sql.DB, error) {
	baseCfg := mysql.Config{
		User:                 c.DB.User,
		Passwd:               c.DB.Password,
		Net:                  "tcp",
		Addr:                 net.JoinHostPort(c.DB.Host, c.DB.Port),
		ParseTime:            true,
		MultiStatements:      true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  time.Local,
		AllowNativePasswords: true,
	}

	adminCfg := baseCfg
	adminDB, err := sql.Open("mysql", adminCfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi admin: %w", err)
	}
	defer adminDB.Close()

	createStmt := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		c.DB.Name,
	)
	if _, err := adminDB.Exec(createStmt); err != nil {
		return nil, fmt.Errorf("gagal membuat database %s: %w", c.DB.Name, err)
	}

	baseCfg.DBName = c.DB.Name
	db, err := sql.Open("mysql", baseCfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("gagal ping database: %w", err)
	}

	return db, nil
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
