package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/digikeys/backend/internal/domain"
)

// ── Mock Repositories ─────────────────────────────────────────────────

type mockCardRepo struct {
	cards    map[string]*domain.Card
	sequence int
}

func newMockCardRepo() *mockCardRepo {
	return &mockCardRepo{cards: make(map[string]*domain.Card), sequence: 0}
}

func (m *mockCardRepo) Create(ctx context.Context, card *domain.Card) error {
	m.cards[card.ID] = card
	return nil
}

func (m *mockCardRepo) GetByID(ctx context.Context, id string) (*domain.Card, error) {
	c, ok := m.cards[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return c, nil
}

func (m *mockCardRepo) GetByCardNumber(ctx context.Context, num string) (*domain.Card, error) {
	for _, c := range m.cards {
		if c.CardNumber == num {
			return c, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockCardRepo) GetByCitizenID(ctx context.Context, citizenID string) ([]*domain.Card, error) {
	return nil, nil
}

func (m *mockCardRepo) List(ctx context.Context, filter domain.CardFilter) ([]*domain.Card, int, error) {
	return nil, 0, nil
}

func (m *mockCardRepo) Update(ctx context.Context, card *domain.Card) error {
	m.cards[card.ID] = card
	return nil
}

func (m *mockCardRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	c, ok := m.cards[id]
	if !ok {
		return domain.ErrNotFound
	}
	c.Status = status
	return nil
}

func (m *mockCardRepo) GetNextSequence(ctx context.Context, embassyID string, year int) (int, error) {
	m.sequence++
	return m.sequence, nil
}

type mockCitizenRepo struct {
	citizens map[string]*domain.Citizen
}

func newMockCitizenRepo() *mockCitizenRepo {
	return &mockCitizenRepo{citizens: make(map[string]*domain.Citizen)}
}

func (m *mockCitizenRepo) Create(ctx context.Context, c *domain.Citizen) error {
	m.citizens[c.ID] = c
	return nil
}

func (m *mockCitizenRepo) GetByID(ctx context.Context, id string) (*domain.Citizen, error) {
	c, ok := m.citizens[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return c, nil
}

func (m *mockCitizenRepo) GetByNationalID(ctx context.Context, nid string) (*domain.Citizen, error) {
	return nil, domain.ErrNotFound
}

func (m *mockCitizenRepo) GetByUniqueIdentifier(ctx context.Context, uid string) (*domain.Citizen, error) {
	return nil, domain.ErrNotFound
}

func (m *mockCitizenRepo) Search(ctx context.Context, f domain.CitizenFilter) ([]*domain.Citizen, int, error) {
	return nil, 0, nil
}

func (m *mockCitizenRepo) Update(ctx context.Context, c *domain.Citizen) error {
	m.citizens[c.ID] = c
	return nil
}

func (m *mockCitizenRepo) List(ctx context.Context, f domain.CitizenFilter) ([]*domain.Citizen, int, error) {
	return nil, 0, nil
}

type mockEmbassyRepo struct {
	embassies map[string]*domain.Embassy
}

func newMockEmbassyRepo() *mockEmbassyRepo {
	return &mockEmbassyRepo{embassies: make(map[string]*domain.Embassy)}
}

func (m *mockEmbassyRepo) Create(ctx context.Context, e *domain.Embassy) error {
	m.embassies[e.ID] = e
	return nil
}

func (m *mockEmbassyRepo) GetByID(ctx context.Context, id string) (*domain.Embassy, error) {
	e, ok := m.embassies[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return e, nil
}

func (m *mockEmbassyRepo) GetByCountryCode(ctx context.Context, code string) (*domain.Embassy, error) {
	return nil, domain.ErrNotFound
}

func (m *mockEmbassyRepo) List(ctx context.Context, page, pageSize int) ([]*domain.Embassy, int, error) {
	return nil, 0, nil
}

func (m *mockEmbassyRepo) Update(ctx context.Context, e *domain.Embassy) error {
	m.embassies[e.ID] = e
	return nil
}

type mockEnrollmentRepo struct {
	enrollments map[string]*domain.Enrollment
}

func newMockEnrollmentRepo() *mockEnrollmentRepo {
	return &mockEnrollmentRepo{enrollments: make(map[string]*domain.Enrollment)}
}

func (m *mockEnrollmentRepo) Create(ctx context.Context, e *domain.Enrollment) error {
	m.enrollments[e.ID] = e
	return nil
}

func (m *mockEnrollmentRepo) GetByID(ctx context.Context, id string) (*domain.Enrollment, error) {
	e, ok := m.enrollments[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return e, nil
}

func (m *mockEnrollmentRepo) ListByAgent(ctx context.Context, agentID string, page, pageSize int) ([]*domain.Enrollment, int, error) {
	return nil, 0, nil
}

func (m *mockEnrollmentRepo) ListByEmbassy(ctx context.Context, embassyID string, page, pageSize int) ([]*domain.Enrollment, int, error) {
	return nil, 0, nil
}

func (m *mockEnrollmentRepo) List(ctx context.Context, filter domain.EnrollmentFilter) ([]*domain.Enrollment, int, error) {
	return nil, 0, nil
}

func (m *mockEnrollmentRepo) Update(ctx context.Context, e *domain.Enrollment) error {
	m.enrollments[e.ID] = e
	return nil
}

func (m *mockEnrollmentRepo) UpdateSyncStatus(ctx context.Context, id string, status string, syncedAt *time.Time) error {
	return nil
}

func (m *mockEnrollmentRepo) UpdateReviewStatus(ctx context.Context, id, status, reviewedBy, notes string, reviewedAt *time.Time) error {
	return nil
}

type mockMRZGenerator struct{}

func (m *mockMRZGenerator) GenerateTD1(citizen *domain.Citizen, card *domain.Card, embassy *domain.Embassy) (string, string, string, error) {
	return strings.Repeat("A", 30), strings.Repeat("B", 30), strings.Repeat("C", 30), nil
}

// ── Test Helpers ──────────────────────────────────────────────────────

func setupCardServiceTest() (*CardService, *mockCardRepo, *mockCitizenRepo, *mockEmbassyRepo, *mockEnrollmentRepo) {
	cardRepo := newMockCardRepo()
	citizenRepo := newMockCitizenRepo()
	embassyRepo := newMockEmbassyRepo()
	enrollmentRepo := newMockEnrollmentRepo()
	mrzGen := &mockMRZGenerator{}

	svc := NewCardService(cardRepo, citizenRepo, embassyRepo, enrollmentRepo, mrzGen)

	// Seed test data
	citizenRepo.citizens["cit-001"] = &domain.Citizen{
		ID:          "cit-001",
		FirstName:   "Amadou",
		LastName:    "Ouedraogo",
		DateOfBirth: time.Date(1990, 3, 20, 0, 0, 0, 0, time.UTC),
		Gender:      "M",
	}
	embassyRepo.embassies["emb-001"] = &domain.Embassy{
		ID:         "emb-001",
		CardPrefix: "BFCI",
		Name:       "Ambassade Abidjan",
	}
	enrollmentRepo.enrollments["enr-001"] = &domain.Enrollment{
		ID:           "enr-001",
		CitizenID:    "cit-001",
		EmbassyID:    "emb-001",
		ReviewStatus: "approved",
	}

	return svc, cardRepo, citizenRepo, embassyRepo, enrollmentRepo
}

// ── Tests ─────────────────────────────────────────────────────────────

func TestRequestCardSuccess(t *testing.T) {
	svc, cardRepo, _, _, _ := setupCardServiceTest()
	ctx := context.Background()

	card, err := svc.RequestCard(ctx, "enr-001", "agent-001")
	if err != nil {
		t.Fatalf("RequestCard failed: %v", err)
	}

	if card.ID == "" {
		t.Error("expected card ID to be set")
	}
	if card.CitizenID != "cit-001" {
		t.Errorf("expected CitizenID=cit-001, got %s", card.CitizenID)
	}
	if card.EmbassyID != "emb-001" {
		t.Errorf("expected EmbassyID=emb-001, got %s", card.EmbassyID)
	}
	if card.Status != string(domain.CardStatusPending) {
		t.Errorf("expected Status=pending, got %s", card.Status)
	}
	if card.IssuedBy != "agent-001" {
		t.Errorf("expected IssuedBy=agent-001, got %s", card.IssuedBy)
	}

	// Card number should start with embassy prefix
	if !strings.HasPrefix(card.CardNumber, "BFCI") {
		t.Errorf("expected card number to start with BFCI, got %s", card.CardNumber)
	}

	// Should contain year
	year := fmt.Sprintf("%d", time.Now().Year())
	if !strings.Contains(card.CardNumber, year) {
		t.Errorf("expected card number to contain year %s, got %s", year, card.CardNumber)
	}

	// MRZ should be set
	if card.MRZLine1 == "" || card.MRZLine2 == "" || card.MRZLine3 == "" {
		t.Error("expected MRZ lines to be set")
	}

	// Stored in repo
	if _, err := cardRepo.GetByID(ctx, card.ID); err != nil {
		t.Errorf("card not found in repo: %v", err)
	}
}

func TestRequestCardEnrollmentNotApproved(t *testing.T) {
	svc, _, _, _, enrollmentRepo := setupCardServiceTest()
	ctx := context.Background()

	// Change enrollment status
	enrollmentRepo.enrollments["enr-001"].ReviewStatus = "pending"

	_, err := svc.RequestCard(ctx, "enr-001", "agent-001")
	if err == nil {
		t.Fatal("expected error for non-approved enrollment")
	}
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got: %v", err)
	}
}

func TestCardStatusTransitions(t *testing.T) {
	svc, cardRepo, _, _, _ := setupCardServiceTest()
	ctx := context.Background()

	card, _ := svc.RequestCard(ctx, "enr-001", "agent-001")

	// pending -> approved
	if err := svc.ApproveCard(ctx, card.ID, "admin-001"); err != nil {
		t.Fatalf("ApproveCard failed: %v", err)
	}
	updated, _ := cardRepo.GetByID(ctx, card.ID)
	if updated.Status != string(domain.CardStatusApproved) {
		t.Errorf("expected approved, got %s", updated.Status)
	}

	// approved -> printing
	if err := svc.QueueForPrinting(ctx, card.ID, "batch-001"); err != nil {
		t.Fatalf("QueueForPrinting failed: %v", err)
	}
	updated, _ = cardRepo.GetByID(ctx, card.ID)
	if updated.Status != string(domain.CardStatusPrinting) {
		t.Errorf("expected printing, got %s", updated.Status)
	}
	if updated.PrintBatchID != "batch-001" {
		t.Errorf("expected PrintBatchID=batch-001, got %s", updated.PrintBatchID)
	}

	// printing -> printed
	if err := svc.MarkPrinted(ctx, card.ID); err != nil {
		t.Fatalf("MarkPrinted failed: %v", err)
	}
	updated, _ = cardRepo.GetByID(ctx, card.ID)
	if updated.Status != string(domain.CardStatusPrinted) {
		t.Errorf("expected printed, got %s", updated.Status)
	}
	if updated.PrintedAt == nil {
		t.Error("expected PrintedAt to be set")
	}

	// printed -> delivered
	if err := svc.MarkDelivered(ctx, card.ID); err != nil {
		t.Fatalf("MarkDelivered failed: %v", err)
	}
	updated, _ = cardRepo.GetByID(ctx, card.ID)
	if updated.Status != string(domain.CardStatusDelivered) {
		t.Errorf("expected delivered, got %s", updated.Status)
	}
	if updated.DeliveredAt == nil {
		t.Error("expected DeliveredAt to be set")
	}

	// delivered -> active
	if err := svc.Activate(ctx, card.ID); err != nil {
		t.Fatalf("Activate failed: %v", err)
	}
	updated, _ = cardRepo.GetByID(ctx, card.ID)
	if updated.Status != string(domain.CardStatusActive) {
		t.Errorf("expected active, got %s", updated.Status)
	}
	if updated.IssuedAt == nil {
		t.Error("expected IssuedAt to be set")
	}
	if updated.ExpiresAt == nil {
		t.Error("expected ExpiresAt to be set")
	}

	// active -> suspended
	if err := svc.Suspend(ctx, card.ID); err != nil {
		t.Fatalf("Suspend failed: %v", err)
	}
	updated, _ = cardRepo.GetByID(ctx, card.ID)
	if updated.Status != string(domain.CardStatusSuspended) {
		t.Errorf("expected suspended, got %s", updated.Status)
	}
}

func TestInvalidStatusTransitions(t *testing.T) {
	svc, _, _, _, _ := setupCardServiceTest()
	ctx := context.Background()

	card, _ := svc.RequestCard(ctx, "enr-001", "agent-001")

	// pending -> cannot print directly
	err := svc.QueueForPrinting(ctx, card.ID, "batch-001")
	if err == nil {
		t.Fatal("expected error: cannot print a pending card")
	}
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got: %v", err)
	}

	// pending -> cannot deliver
	err = svc.MarkDelivered(ctx, card.ID)
	if err == nil {
		t.Fatal("expected error: cannot deliver a pending card")
	}

	// pending -> cannot activate
	err = svc.Activate(ctx, card.ID)
	if err == nil {
		t.Fatal("expected error: cannot activate a pending card")
	}

	// pending -> cannot suspend
	err = svc.Suspend(ctx, card.ID)
	if err == nil {
		t.Fatal("expected error: cannot suspend a pending card")
	}
}

func TestRevokeCard(t *testing.T) {
	svc, cardRepo, _, _, _ := setupCardServiceTest()
	ctx := context.Background()

	card, _ := svc.RequestCard(ctx, "enr-001", "agent-001")
	_ = svc.ApproveCard(ctx, card.ID, "admin")

	err := svc.Revoke(ctx, card.ID)
	if err != nil {
		t.Fatalf("Revoke failed: %v", err)
	}
	updated, _ := cardRepo.GetByID(ctx, card.ID)
	if updated.Status != string(domain.CardStatusRevoked) {
		t.Errorf("expected revoked, got %s", updated.Status)
	}

	// Cannot revoke again
	err = svc.Revoke(ctx, card.ID)
	if err == nil {
		t.Fatal("expected error revoking already-revoked card")
	}
}

func TestCardNumberFormat(t *testing.T) {
	svc, _, _, _, _ := setupCardServiceTest()
	ctx := context.Background()

	card, err := svc.RequestCard(ctx, "enr-001", "agent-001")
	if err != nil {
		t.Fatalf("RequestCard failed: %v", err)
	}

	// Format: PREFIX + YEAR + 6-digit sequence
	// e.g., "BFCI2026000001"
	prefix := "BFCI"
	year := fmt.Sprintf("%d", time.Now().Year())
	expectedPrefix := prefix + year

	if !strings.HasPrefix(card.CardNumber, expectedPrefix) {
		t.Errorf("card number %q should start with %q", card.CardNumber, expectedPrefix)
	}

	// Total length: 4 (prefix) + 4 (year) + 6 (sequence) = 14
	if len(card.CardNumber) != 14 {
		t.Errorf("card number length = %d, want 14 for %q", len(card.CardNumber), card.CardNumber)
	}
}
