package application

import (
	"context"

	"github.com/digikeys/backend/internal/domain"
	"github.com/digikeys/backend/internal/ports"
)

type VerificationService struct {
	cardRepo    ports.CardRepository
	citizenRepo ports.CitizenRepository
}

func NewVerificationService(cardRepo ports.CardRepository, citizenRepo ports.CitizenRepository) *VerificationService {
	return &VerificationService{
		cardRepo:    cardRepo,
		citizenRepo: citizenRepo,
	}
}

type VerificationResult struct {
	Valid     bool   `json:"valid"`
	Status   string `json:"status"`
	CardNumber string `json:"cardNumber"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	EmbassyID  string `json:"embassyId,omitempty"`
	ExpiresAt  string `json:"expiresAt,omitempty"`
	Message    string `json:"message,omitempty"`
}

func (s *VerificationService) VerifyCard(ctx context.Context, cardNumber string) (*VerificationResult, error) {
	card, err := s.cardRepo.GetByCardNumber(ctx, cardNumber)
	if err != nil {
		return &VerificationResult{
			Valid:   false,
			Message: "Carte non trouvée",
		}, nil
	}

	result := &VerificationResult{
		CardNumber: card.CardNumber,
		Status:     card.Status,
		EmbassyID:  card.EmbassyID,
	}

	switch domain.CardStatus(card.Status) {
	case domain.CardStatusActive:
		result.Valid = true
		result.Message = "Carte valide"
	case domain.CardStatusExpired:
		result.Valid = false
		result.Message = "Carte expirée"
	case domain.CardStatusSuspended:
		result.Valid = false
		result.Message = "Carte suspendue"
	case domain.CardStatusRevoked:
		result.Valid = false
		result.Message = "Carte révoquée"
	default:
		result.Valid = false
		result.Message = "Carte non active"
	}

	if card.ExpiresAt != nil {
		result.ExpiresAt = card.ExpiresAt.Format("2006-01-02")
	}

	// Include basic citizen info (no biometrics)
	citizen, err := s.citizenRepo.GetByID(ctx, card.CitizenID)
	if err == nil {
		result.FirstName = citizen.FirstName
		result.LastName = citizen.LastName
	}

	return result, nil
}
