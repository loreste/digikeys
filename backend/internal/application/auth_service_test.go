package application

import (
	"context"
	"errors"
	"testing"

	"github.com/digikeys/backend/config"
	"github.com/digikeys/backend/internal/domain"
)

// ── Mock User Repository ──────────────────────────────────────────────

type mockUserRepo struct {
	users map[string]*domain.User // keyed by email
	byID  map[string]*domain.User // keyed by id
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[string]*domain.User),
		byID:  make(map[string]*domain.User),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	m.users[user.Email] = user
	m.byID[user.ID] = user
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	u, ok := m.byID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return u, nil
}

func (m *mockUserRepo) List(ctx context.Context, role domain.UserRole, embassyID string, page, pageSize int) ([]*domain.User, int, error) {
	return nil, 0, nil
}

func (m *mockUserRepo) Update(ctx context.Context, user *domain.User) error {
	m.users[user.Email] = user
	m.byID[user.ID] = user
	return nil
}

func (m *mockUserRepo) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	u, ok := m.byID[id]
	if !ok {
		return domain.ErrNotFound
	}
	u.PasswordHash = passwordHash
	return nil
}

// ── Test Helpers ──────────────────────────────────────────────────────

func testJWTConfig() config.JWTConfig {
	return config.JWTConfig{
		Secret:          "test-secret-key-for-unit-tests-only",
		AccessTokenTTL:  15,
		RefreshTokenTTL: 7,
	}
}

func testUser() *domain.User {
	return &domain.User{
		ID:        "user-001",
		Email:     "agent@embassy.bf",
		FirstName: "Amadou",
		LastName:  "Ouedraogo",
		Role:      domain.UserRoleEnrollmentAgent,
		EmbassyID: "embassy-001",
		Status:    "active",
	}
}

// ── Tests ─────────────────────────────────────────────────────────────

func TestRegisterSuccess(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, testJWTConfig())
	ctx := context.Background()

	user := testUser()
	err := svc.Register(ctx, user, "strongpassword123")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if user.PasswordHash == "" {
		t.Fatal("expected password hash to be set")
	}
	if user.PasswordHash == "strongpassword123" {
		t.Fatal("password hash should not be the raw password")
	}
	if user.Status != "active" {
		t.Errorf("expected status=active, got %s", user.Status)
	}

	// Verify stored in repo
	stored, err := repo.GetByEmail(ctx, "agent@embassy.bf")
	if err != nil {
		t.Fatalf("user not found in repo: %v", err)
	}
	if stored.ID != "user-001" {
		t.Errorf("expected stored user ID=user-001, got %s", stored.ID)
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, testJWTConfig())
	ctx := context.Background()

	user := testUser()
	_ = svc.Register(ctx, user, "password123")

	user2 := &domain.User{
		ID:    "user-002",
		Email: "agent@embassy.bf",
		Role:  domain.UserRoleEnrollmentAgent,
	}
	err := svc.Register(ctx, user2, "password456")
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Errorf("expected ErrAlreadyExists, got: %v", err)
	}
}

func TestLoginSuccess(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, testJWTConfig())
	ctx := context.Background()

	user := testUser()
	_ = svc.Register(ctx, user, "correctpassword")

	tokens, loggedUser, err := svc.Login(ctx, "agent@embassy.bf", "correctpassword")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if tokens == nil {
		t.Fatal("expected token pair, got nil")
	}
	if tokens.AccessToken == "" {
		t.Fatal("expected non-empty access token")
	}
	if tokens.RefreshToken == "" {
		t.Fatal("expected non-empty refresh token")
	}
	if tokens.ExpiresIn != 15*60 {
		t.Errorf("expected ExpiresIn=%d, got %d", 15*60, tokens.ExpiresIn)
	}
	if loggedUser.Email != "agent@embassy.bf" {
		t.Errorf("expected email=agent@embassy.bf, got %s", loggedUser.Email)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, testJWTConfig())
	ctx := context.Background()

	user := testUser()
	_ = svc.Register(ctx, user, "correctpassword")

	_, _, err := svc.Login(ctx, "agent@embassy.bf", "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got: %v", err)
	}
}

func TestLoginNonexistentUser(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, testJWTConfig())
	ctx := context.Background()

	_, _, err := svc.Login(ctx, "nobody@example.com", "password")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got: %v", err)
	}
}

func TestLoginInactiveUser(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, testJWTConfig())
	ctx := context.Background()

	user := testUser()
	_ = svc.Register(ctx, user, "password123")
	user.Status = "suspended"
	_ = repo.Update(ctx, user)

	_, _, err := svc.Login(ctx, "agent@embassy.bf", "password123")
	if err == nil {
		t.Fatal("expected error for inactive user")
	}
	if !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got: %v", err)
	}
}

func TestValidateToken(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, testJWTConfig())
	ctx := context.Background()

	user := testUser()
	_ = svc.Register(ctx, user, "password123")

	tokens, _, _ := svc.Login(ctx, "agent@embassy.bf", "password123")

	claims, err := svc.ValidateToken(tokens.AccessToken)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if claims.Subject != "user-001" {
		t.Errorf("expected Subject=user-001, got %s", claims.Subject)
	}
	if claims.Role != domain.UserRoleEnrollmentAgent {
		t.Errorf("expected Role=enrollment_agent, got %s", claims.Role)
	}
	if claims.EmbassyID != "embassy-001" {
		t.Errorf("expected EmbassyID=embassy-001, got %s", claims.EmbassyID)
	}
	if claims.Issuer != "carteconsulaire" {
		t.Errorf("expected Issuer=carteconsulaire, got %s", claims.Issuer)
	}
}

func TestValidateTokenInvalid(t *testing.T) {
	svc := NewAuthService(newMockUserRepo(), testJWTConfig())

	_, err := svc.ValidateToken("invalid.token.string")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestValidateTokenWrongSecret(t *testing.T) {
	repo := newMockUserRepo()
	svc1 := NewAuthService(repo, testJWTConfig())
	ctx := context.Background()

	user := testUser()
	_ = svc1.Register(ctx, user, "password123")
	tokens, _, _ := svc1.Login(ctx, "agent@embassy.bf", "password123")

	otherCfg := testJWTConfig()
	otherCfg.Secret = "different-secret-key"
	svc2 := NewAuthService(repo, otherCfg)

	_, err := svc2.ValidateToken(tokens.AccessToken)
	if err == nil {
		t.Fatal("expected error validating token with different secret")
	}
}

func TestRefreshToken(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, testJWTConfig())
	ctx := context.Background()

	user := testUser()
	_ = svc.Register(ctx, user, "password123")
	tokens, _, _ := svc.Login(ctx, "agent@embassy.bf", "password123")

	newTokens, err := svc.RefreshToken(ctx, tokens.RefreshToken)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}
	if newTokens.AccessToken == "" {
		t.Fatal("expected non-empty new access token")
	}
	if newTokens.RefreshToken == "" {
		t.Fatal("expected non-empty new refresh token")
	}
	// Tokens are valid JWT strings
	if len(newTokens.AccessToken) < 50 {
		t.Error("access token seems too short to be a valid JWT")
	}
}

func TestRefreshTokenInvalid(t *testing.T) {
	svc := NewAuthService(newMockUserRepo(), testJWTConfig())
	ctx := context.Background()

	_, err := svc.RefreshToken(ctx, "garbage-token")
	if err == nil {
		t.Fatal("expected error for invalid refresh token")
	}
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got: %v", err)
	}
}
