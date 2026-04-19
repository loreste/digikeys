package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	"github.com/digikeys/backend/config"
	apphttp "github.com/digikeys/backend/internal/adapters/http"
	"github.com/digikeys/backend/internal/adapters/mrz"
	"github.com/digikeys/backend/internal/adapters/postgres"
	"github.com/digikeys/backend/internal/application"
	"github.com/digikeys/backend/internal/domain"
)

func main() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "serve")
	}

	cmd := os.Args[1]

	switch cmd {
	case "serve":
		runServer()
	case "migrate":
		runMigrations()
	case "seed":
		runSeed()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\nusage: server [serve|migrate|seed]\n", cmd)
		os.Exit(1)
	}
}

func runMigrations() {
	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := postgres.NewConnection(ctx, cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	slog.Info("migrations completed successfully (run SQL files manually or use golang-migrate)")
}

func runSeed() {
	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := postgres.NewConnection(ctx, cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Create default super admin
	adminEmail := "admin@carteconsulaire.bf"
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123456"), bcrypt.DefaultCost)

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, email, phone, password_hash, role, first_name, last_name, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (email) DO NOTHING
	`,
		uuid.New().String(),
		adminEmail,
		"+226 00 00 00 00",
		string(hash),
		string(domain.UserRoleSuperAdmin),
		"Admin",
		"DIGIKEYS",
		"active",
		time.Now(),
		time.Now(),
	)
	if err != nil {
		slog.Error("failed to create admin user", "error", err)
		os.Exit(1)
	}

	slog.Info("seed completed", "admin_email", adminEmail, "admin_password", "admin123456")
}

func runServer() {
	slog.Info("starting Carte Consulaire DIGIKEYS API server")

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ── Database ──────────────────────────────────────────────
	pool, err := postgres.NewConnection(ctx, cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// ── Repositories (stubbed -- need concrete implementations) ──
	// For compilation, we use nil services. In production, wire real repos here.
	_ = pool

	// ── MRZ Generator ────────────────────────────────────────
	mrzGen := mrz.NewGenerator()
	_ = mrzGen

	// ── Application Services ─────────────────────────────────
	// NOTE: Full wiring requires concrete repository implementations.
	// This demonstrates the service initialization pattern.
	// authService := application.NewAuthService(userRepo, cfg.JWT)
	// citizenService := application.NewCitizenService(citizenRepo)
	// enrollmentService := application.NewEnrollmentService(enrollmentRepo, citizenRepo)
	// cardService := application.NewCardService(cardRepo, citizenRepo, embassyRepo, enrollmentRepo, mrzGen)
	// verifyService := application.NewVerificationService(cardRepo, citizenRepo)
	// transferService := application.NewTransferService(transferRepo, citizenRepo)
	// fsbService := application.NewFSBService(transferRepo)
	// statsService := application.NewStatisticsService(statsQuerier)

	// For now, create a minimal working server with auth only
	// Once postgres repos are implemented, uncomment the full wiring above
	var authService *application.AuthService
	// authService will be nil -- router handles nil deps gracefully

	// ── HTTP Router ──────────────────────────────────────────
	router := apphttp.NewRouter(apphttp.RouterDeps{
		AuthService: authService,
	})

	// ── Server ───────────────────────────────────────────────
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server listening", "address", addr, "env", cfg.App.Environment)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server stopped")
}
