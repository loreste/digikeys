package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Clear any env vars that might interfere
	envVars := []string{
		"SERVER_HOST", "SERVER_PORT",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE",
		"JWT_SECRET", "JWT_ACCESS_TTL", "JWT_REFRESH_TTL",
		"STORAGE_ENDPOINT", "STORAGE_ACCESS_KEY", "STORAGE_SECRET_KEY", "STORAGE_BUCKET", "STORAGE_USE_SSL",
		"APP_BASE_URL", "APP_ENV", "APP_COUNTRY",
		"SMS_SENDER_ID",
	}
	for _, k := range envVars {
		os.Unsetenv(k)
	}

	cfg := Load()

	// Server defaults
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected Server.Host=0.0.0.0, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected Server.Port=8080, got %d", cfg.Server.Port)
	}

	// Database defaults
	if cfg.Database.Host != "localhost" {
		t.Errorf("expected DB.Host=localhost, got %s", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("expected DB.Port=5432, got %d", cfg.Database.Port)
	}
	if cfg.Database.User != "carteconsulaire" {
		t.Errorf("expected DB.User=carteconsulaire, got %s", cfg.Database.User)
	}
	if cfg.Database.SSLMode != "disable" {
		t.Errorf("expected DB.SSLMode=disable, got %s", cfg.Database.SSLMode)
	}

	// JWT defaults
	if cfg.JWT.AccessTokenTTL != 15 {
		t.Errorf("expected JWT.AccessTokenTTL=15, got %d", cfg.JWT.AccessTokenTTL)
	}
	if cfg.JWT.RefreshTokenTTL != 7 {
		t.Errorf("expected JWT.RefreshTokenTTL=7, got %d", cfg.JWT.RefreshTokenTTL)
	}

	// App defaults
	if cfg.App.Environment != "development" {
		t.Errorf("expected App.Environment=development, got %s", cfg.App.Environment)
	}
	if cfg.App.Country != "BF" {
		t.Errorf("expected App.Country=BF, got %s", cfg.App.Country)
	}
	if cfg.App.BaseURL != "http://localhost:8080" {
		t.Errorf("expected App.BaseURL=http://localhost:8080, got %s", cfg.App.BaseURL)
	}

	// SMS defaults
	if cfg.SMS.SenderID != "DIGIKEYS" {
		t.Errorf("expected SMS.SenderID=DIGIKEYS, got %s", cfg.SMS.SenderID)
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	t.Setenv("SERVER_HOST", "127.0.0.1")
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("DB_HOST", "db.prod.example.com")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USER", "produser")
	t.Setenv("DB_PASSWORD", "s3cret")
	t.Setenv("DB_NAME", "proddb")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("JWT_SECRET", "my-prod-secret")
	t.Setenv("JWT_ACCESS_TTL", "30")
	t.Setenv("JWT_REFRESH_TTL", "14")
	t.Setenv("APP_ENV", "production")
	t.Setenv("APP_COUNTRY", "CD")
	t.Setenv("STORAGE_USE_SSL", "true")

	cfg := Load()

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("expected Server.Host=127.0.0.1, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("expected Server.Port=9090, got %d", cfg.Server.Port)
	}
	if cfg.Database.Host != "db.prod.example.com" {
		t.Errorf("expected DB.Host=db.prod.example.com, got %s", cfg.Database.Host)
	}
	if cfg.Database.Port != 5433 {
		t.Errorf("expected DB.Port=5433, got %d", cfg.Database.Port)
	}
	if cfg.Database.User != "produser" {
		t.Errorf("expected DB.User=produser, got %s", cfg.Database.User)
	}
	if cfg.Database.Password != "s3cret" {
		t.Errorf("expected DB.Password=s3cret, got %s", cfg.Database.Password)
	}
	if cfg.Database.SSLMode != "require" {
		t.Errorf("expected DB.SSLMode=require, got %s", cfg.Database.SSLMode)
	}
	if cfg.JWT.Secret != "my-prod-secret" {
		t.Errorf("expected JWT.Secret=my-prod-secret, got %s", cfg.JWT.Secret)
	}
	if cfg.JWT.AccessTokenTTL != 30 {
		t.Errorf("expected JWT.AccessTokenTTL=30, got %d", cfg.JWT.AccessTokenTTL)
	}
	if cfg.JWT.RefreshTokenTTL != 14 {
		t.Errorf("expected JWT.RefreshTokenTTL=14, got %d", cfg.JWT.RefreshTokenTTL)
	}
	if cfg.App.Environment != "production" {
		t.Errorf("expected App.Environment=production, got %s", cfg.App.Environment)
	}
	if cfg.App.Country != "CD" {
		t.Errorf("expected App.Country=CD, got %s", cfg.App.Country)
	}
	if !cfg.Storage.UseSSL {
		t.Error("expected Storage.UseSSL=true")
	}
}

func TestLoadInvalidIntFallsBackToDefault(t *testing.T) {
	t.Setenv("SERVER_PORT", "not-a-number")
	t.Setenv("DB_PORT", "xyz")

	cfg := Load()

	if cfg.Server.Port != 8080 {
		t.Errorf("expected Server.Port fallback=8080 on invalid int, got %d", cfg.Server.Port)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("expected DB.Port fallback=5432 on invalid int, got %d", cfg.Database.Port)
	}
}

func TestLoadInvalidBoolFallsBackToDefault(t *testing.T) {
	t.Setenv("STORAGE_USE_SSL", "not-a-bool")

	cfg := Load()

	if cfg.Storage.UseSSL != false {
		t.Error("expected Storage.UseSSL fallback=false on invalid bool")
	}
}

func TestDatabaseDSN(t *testing.T) {
	db := DatabaseConfig{
		Host:     "dbhost",
		Port:     5432,
		User:     "user",
		Password: "pass",
		DBName:   "mydb",
		SSLMode:  "disable",
	}
	expected := "postgres://user:pass@dbhost:5432/mydb?sslmode=disable"
	if got := db.DSN(); got != expected {
		t.Errorf("DSN() = %s, want %s", got, expected)
	}
}
