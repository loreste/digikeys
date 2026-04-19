package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/digikeys/backend/internal/domain"
	"github.com/digikeys/backend/internal/ports"
)

// FSBService manages the Fonds de Solidarité Burkinabè contributions
// that are collected upon card issuance.
type FSBService struct {
	transferRepo ports.TransferRepository
}

func NewFSBService(transferRepo ports.TransferRepository) *FSBService {
	return &FSBService{transferRepo: transferRepo}
}

func (s *FSBService) RecordContribution(ctx context.Context, citizenID string, amount int64, currency string) error {
	transfer := &domain.Transfer{
		ID:        uuid.New().String(),
		CitizenID: citizenID,
		Amount:    amount,
		Currency:  currency,
		Type:      "fsb_contribution",
		Status:    "completed",
		CreatedAt: time.Now(),
	}

	now := time.Now()
	transfer.CompletedAt = &now

	return s.transferRepo.Create(ctx, transfer)
}

type FSBSummary struct {
	TotalContributions int64  `json:"totalContributions"`
	TotalAmount        int64  `json:"totalAmount"`
	Currency           string `json:"currency"`
}

func (s *FSBService) GetSummary(ctx context.Context) (*FSBSummary, error) {
	transfers, total, err := s.transferRepo.List(ctx, domain.TransferFilter{
		Type:     "fsb_contribution",
		Status:   "completed",
		Page:     1,
		PageSize: 1,
	})
	_ = transfers

	var totalAmount int64
	// Sum is done via query; for simplicity, return count-based summary
	return &FSBSummary{
		TotalContributions: int64(total),
		TotalAmount:        totalAmount,
		Currency:           "XOF",
	}, err
}

type FSBReport struct {
	TotalContributions int64            `json:"totalContributions"`
	TotalAmount        int64            `json:"totalAmount"`
	Currency           string           `json:"currency"`
	ByCountry          map[string]int64 `json:"byCountry,omitempty"`
}

func (s *FSBService) GetReport(ctx context.Context) (*FSBReport, error) {
	summary, err := s.GetSummary(ctx)
	if err != nil {
		return nil, err
	}

	return &FSBReport{
		TotalContributions: summary.TotalContributions,
		TotalAmount:        summary.TotalAmount,
		Currency:           summary.Currency,
		ByCountry:          make(map[string]int64),
	}, nil
}
