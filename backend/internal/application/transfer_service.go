package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/digikeys/backend/internal/domain"
	"github.com/digikeys/backend/internal/ports"
)

type TransferService struct {
	transferRepo ports.TransferRepository
	citizenRepo  ports.CitizenRepository
}

func NewTransferService(transferRepo ports.TransferRepository, citizenRepo ports.CitizenRepository) *TransferService {
	return &TransferService{
		transferRepo: transferRepo,
		citizenRepo:  citizenRepo,
	}
}

func (s *TransferService) InitiateTransfer(ctx context.Context, transfer *domain.Transfer) error {
	// Verify citizen exists
	_, err := s.citizenRepo.GetByID(ctx, transfer.CitizenID)
	if err != nil {
		return fmt.Errorf("%w: citizen not found", domain.ErrNotFound)
	}

	transfer.ID = uuid.New().String()
	transfer.Status = "pending"
	transfer.CreatedAt = time.Now()

	return s.transferRepo.Create(ctx, transfer)
}

func (s *TransferService) ProcessWebhook(ctx context.Context, externalRef, status, failureReason string) error {
	transfer, err := s.transferRepo.GetByExternalRef(ctx, externalRef)
	if err != nil {
		return err
	}

	transfer.Status = status
	if failureReason != "" {
		transfer.FailureReason = failureReason
	}

	if status == "completed" {
		now := time.Now()
		transfer.CompletedAt = &now
	}

	return s.transferRepo.Update(ctx, transfer)
}

func (s *TransferService) GetTransferHistory(ctx context.Context, filter domain.TransferFilter) ([]*domain.Transfer, int, error) {
	return s.transferRepo.List(ctx, filter)
}
