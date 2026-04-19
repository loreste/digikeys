package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	Storage   StorageConfig
	Biometric BiometricConfig
	Banking   BankingConfig
	SMS       SMSConfig
	Email     EmailConfig
	App       AppConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  int // minutes
	RefreshTokenTTL int // days
}

type StorageConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type BiometricConfig struct {
	EncryptionKeyPath string
}

type BankingConfig struct {
	BaseURL   string
	APIKey    string
	APISecret string
	WebhookSecret string
}

type SMSConfig struct {
	Provider string
	APIKey   string
	SenderID string
}

type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	FromAddr string
	FromName string
}

type AppConfig struct {
	BaseURL     string
	Environment string
	Country     string // "BF" or "CD" - determines which country this deployment serves
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode)
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "carteconsulaire"),
			Password: getEnv("DB_PASSWORD", "carteconsulaire_dev_password"),
			DBName:   getEnv("DB_NAME", "carteconsulaire"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "carteconsulaire-dev-secret-change-in-production"),
			AccessTokenTTL:  getEnvInt("JWT_ACCESS_TTL", 15),
			RefreshTokenTTL: getEnvInt("JWT_REFRESH_TTL", 7),
		},
		Storage: StorageConfig{
			Endpoint:  getEnv("STORAGE_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("STORAGE_ACCESS_KEY", "carteconsulaire"),
			SecretKey: getEnv("STORAGE_SECRET_KEY", "carteconsulaire_dev_password"),
			Bucket:    getEnv("STORAGE_BUCKET", "carteconsulaire"),
			UseSSL:    getEnvBool("STORAGE_USE_SSL", false),
		},
		Biometric: BiometricConfig{
			EncryptionKeyPath: getEnv("BIOMETRIC_ENCRYPTION_KEY_PATH", ""),
		},
		Banking: BankingConfig{
			BaseURL:       getEnv("BANKING_BASE_URL", ""),
			APIKey:        getEnv("BANKING_API_KEY", ""),
			APISecret:     getEnv("BANKING_API_SECRET", ""),
			WebhookSecret: getEnv("BANKING_WEBHOOK_SECRET", ""),
		},
		SMS: SMSConfig{
			Provider: getEnv("SMS_PROVIDER", ""),
			APIKey:   getEnv("SMS_API_KEY", ""),
			SenderID: getEnv("SMS_SENDER_ID", "DIGIKEYS"),
		},
		Email: EmailConfig{
			SMTPHost: getEnv("EMAIL_SMTP_HOST", "localhost"),
			SMTPPort: getEnvInt("EMAIL_SMTP_PORT", 587),
			Username: getEnv("EMAIL_USERNAME", ""),
			Password: getEnv("EMAIL_PASSWORD", ""),
			FromAddr: getEnv("EMAIL_FROM_ADDR", "noreply@carteconsulaire.bf"),
			FromName: getEnv("EMAIL_FROM_NAME", "Carte Consulaire DIGIKEYS"),
		},
		App: AppConfig{
			BaseURL:     getEnv("APP_BASE_URL", "http://localhost:8080"),
			Environment: getEnv("APP_ENV", "development"),
			Country:     getEnv("APP_COUNTRY", "BF"), // BF or CD
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
