package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/digikeys/backend/internal/domain"
	"github.com/digikeys/backend/internal/ports"
)

type CardService struct {
	cardRepo       ports.CardRepository
	citizenRepo    ports.CitizenRepository
	embassyRepo    ports.EmbassyRepository
	enrollmentRepo ports.EnrollmentRepository
	mrzGen         ports.MRZGenerator
}

func NewCardService(
	cardRepo ports.CardRepository,
	citizenRepo ports.CitizenRepository,
	embassyRepo ports.EmbassyRepository,
	enrollmentRepo ports.EnrollmentRepository,
	mrzGen ports.MRZGenerator,
) *CardService {
	return &CardService{
		cardRepo:       cardRepo,
		citizenRepo:    citizenRepo,
		embassyRepo:    embassyRepo,
		enrollmentRepo: enrollmentRepo,
		mrzGen:         mrzGen,
	}
}

func (s *CardService) RequestCard(ctx context.Context, enrollmentID, issuedBy string) (*domain.Card, error) {
	enrollment, err := s.enrollmentRepo.GetByID(ctx, enrollmentID)
	if err != nil {
		return nil, err
	}

	if enrollment.ReviewStatus != "approved" {
		return nil, fmt.Errorf("%w: enrollment must be approved before card request", domain.ErrInvalidInput)
	}

	citizen, err := s.citizenRepo.GetByID(ctx, enrollment.CitizenID)
	if err != nil {
		return nil, err
	}

	embassy, err := s.embassyRepo.GetByID(ctx, enrollment.EmbassyID)
	if err != nil {
		return nil, err
	}

	// Generate card number: PREFIX + YEAR + SEQUENCE
	year := time.Now().Year()
	seq, err := s.cardRepo.GetNextSequence(ctx, embassy.ID, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get card sequence: %w", err)
	}

	cardNumber := fmt.Sprintf("%s%d%06d", embassy.CardPrefix, year, seq)

	card := &domain.Card{
		ID:         uuid.New().String(),
		CitizenID:  citizen.ID,
		CardNumber: cardNumber,
		EmbassyID:  embassy.ID,
		IssuedBy:   issuedBy,
		Status:     string(domain.CardStatusPending),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Generate MRZ
	line1, line2, line3, err := s.mrzGen.GenerateTD1(citizen, card, embassy)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidMRZ, err)
	}
	card.MRZLine1 = line1
	card.MRZLine2 = line2
	card.MRZLine3 = line3

	if err := s.cardRepo.Create(ctx, card); err != nil {
		return nil, err
	}

	return card, nil
}

func (s *CardService) ApproveCard(ctx context.Context, id, approvedBy string) error {
	card, err := s.cardRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if card.Status != string(domain.CardStatusPending) {
		return fmt.Errorf("%w: card must be pending to approve", domain.ErrInvalidInput)
	}

	card.Status = string(domain.CardStatusApproved)
	card.UpdatedAt = time.Now()
	return s.cardRepo.Update(ctx, card)
}

func (s *CardService) QueueForPrinting(ctx context.Context, id, batchID string) error {
	card, err := s.cardRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if card.Status != string(domain.CardStatusApproved) {
		return fmt.Errorf("%w: card must be approved to print", domain.ErrInvalidInput)
	}

	card.Status = string(domain.CardStatusPrinting)
	card.PrintBatchID = batchID
	card.UpdatedAt = time.Now()
	return s.cardRepo.Update(ctx, card)
}

func (s *CardService) MarkPrinted(ctx context.Context, id string) error {
	card, err := s.cardRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if card.Status != string(domain.CardStatusPrinting) {
		return fmt.Errorf("%w: card must be in printing to mark printed", domain.ErrInvalidInput)
	}

	now := time.Now()
	card.Status = string(domain.CardStatusPrinted)
	card.PrintedAt = &now
	card.UpdatedAt = now
	return s.cardRepo.Update(ctx, card)
}

func (s *CardService) MarkDelivered(ctx context.Context, id string) error {
	card, err := s.cardRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if card.Status != string(domain.CardStatusPrinted) {
		return fmt.Errorf("%w: card must be printed to deliver", domain.ErrInvalidInput)
	}

	now := time.Now()
	card.Status = string(domain.CardStatusDelivered)
	card.DeliveredAt = &now
	card.UpdatedAt = now
	return s.cardRepo.Update(ctx, card)
}

func (s *CardService) Activate(ctx context.Context, id string) error {
	card, err := s.cardRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if card.Status != string(domain.CardStatusDelivered) {
		return fmt.Errorf("%w: card must be delivered to activate", domain.ErrInvalidInput)
	}

	now := time.Now()
	expiry := now.Add(5 * 365 * 24 * time.Hour) // 5-year validity
	card.Status = string(domain.CardStatusActive)
	card.IssuedAt = &now
	card.ExpiresAt = &expiry
	card.UpdatedAt = now
	return s.cardRepo.Update(ctx, card)
}

func (s *CardService) Suspend(ctx context.Context, id string) error {
	card, err := s.cardRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if card.Status != string(domain.CardStatusActive) {
		return fmt.Errorf("%w: card must be active to suspend", domain.ErrInvalidInput)
	}

	card.Status = string(domain.CardStatusSuspended)
	card.UpdatedAt = time.Now()
	return s.cardRepo.Update(ctx, card)
}

func (s *CardService) Revoke(ctx context.Context, id string) error {
	card, err := s.cardRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if card.Status == string(domain.CardStatusRevoked) {
		return fmt.Errorf("%w: card is already revoked", domain.ErrInvalidInput)
	}

	card.Status = string(domain.CardStatusRevoked)
	card.UpdatedAt = time.Now()
	return s.cardRepo.Update(ctx, card)
}

func (s *CardService) Renew(ctx context.Context, oldCardID, reason, issuedBy string) (*domain.Card, error) {
	oldCard, err := s.cardRepo.GetByID(ctx, oldCardID)
	if err != nil {
		return nil, err
	}

	citizen, err := s.citizenRepo.GetByID(ctx, oldCard.CitizenID)
	if err != nil {
		return nil, err
	}

	embassy, err := s.embassyRepo.GetByID(ctx, oldCard.EmbassyID)
	if err != nil {
		return nil, err
	}

	// Revoke old card
	oldCard.Status = string(domain.CardStatusRevoked)
	oldCard.UpdatedAt = time.Now()
	if err := s.cardRepo.Update(ctx, oldCard); err != nil {
		return nil, err
	}

	// Generate new card number
	year := time.Now().Year()
	seq, err := s.cardRepo.GetNextSequence(ctx, embassy.ID, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get card sequence: %w", err)
	}

	cardNumber := fmt.Sprintf("%s%d%06d", embassy.CardPrefix, year, seq)

	newCard := &domain.Card{
		ID:             uuid.New().String(),
		CitizenID:      citizen.ID,
		CardNumber:     cardNumber,
		EmbassyID:      embassy.ID,
		IssuedBy:       issuedBy,
		Status:         string(domain.CardStatusPending),
		PreviousCardID: oldCardID,
		RenewalReason:  reason,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	line1, line2, line3, err := s.mrzGen.GenerateTD1(citizen, newCard, embassy)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidMRZ, err)
	}
	newCard.MRZLine1 = line1
	newCard.MRZLine2 = line2
	newCard.MRZLine3 = line3

	if err := s.cardRepo.Create(ctx, newCard); err != nil {
		return nil, err
	}

	return newCard, nil
}

func (s *CardService) GetByID(ctx context.Context, id string) (*domain.Card, error) {
	return s.cardRepo.GetByID(ctx, id)
}

func (s *CardService) GetByCardNumber(ctx context.Context, cardNumber string) (*domain.Card, error) {
	return s.cardRepo.GetByCardNumber(ctx, cardNumber)
}

func (s *CardService) List(ctx context.Context, filter domain.CardFilter) ([]*domain.Card, int, error) {
	return s.cardRepo.List(ctx, filter)
}
