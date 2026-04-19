package banking

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/digikeys/backend/internal/domain"
)

// BankAdapter implements ports.BankingService.
// This is a placeholder implementation that logs operations and returns mock data.
// Real bank API integration (Coris Bank, Rawbank, etc.) will replace this later.
type BankAdapter struct {
	baseURL   string
	apiKey    string
	apiSecret string
}

// NewBankAdapter creates a new banking adapter.
func NewBankAdapter(baseURL, apiKey, apiSecret string) *BankAdapter {
	return &BankAdapter{
		baseURL:   baseURL,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

// OpenAccount opens a new bank account for a citizen and returns the account details.
// Placeholder: generates a mock account number and logs the request.
func (b *BankAdapter) OpenAccount(ctx context.Context, citizen *domain.Citizen) (*domain.BankAccount, error) {
	if citizen == nil {
		return nil, fmt.Errorf("citizen is required")
	}

	slog.Info("opening bank account (placeholder)",
		"citizenId", citizen.ID,
		"name", citizen.FirstName+" "+citizen.LastName,
		"country", citizen.CountryOfResidence,
	)

	// Generate a mock account number.
	mockAccountNum := fmt.Sprintf("CC%s%d", citizen.ID[:8], time.Now().UnixNano()%100000)

	now := time.Now()
	account := &domain.BankAccount{
		ID:            uuid.New().String(),
		CitizenID:     citizen.ID,
		BankCode:      "CORIS",
		BankName:      "Coris Bank International",
		AccountNumber: mockAccountNum,
		IBAN:          fmt.Sprintf("BF00CORIS%s", mockAccountNum),
		Status:        "active",
		OpenedAt:      &now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	slog.Info("bank account opened (placeholder)",
		"citizenId", citizen.ID,
		"accountNumber", account.AccountNumber,
	)

	return account, nil
}

// InitiateTransfer initiates a money transfer and returns an external reference.
// Placeholder: generates a mock reference and logs the request.
func (b *BankAdapter) InitiateTransfer(ctx context.Context, transfer *domain.Transfer) (string, error) {
	if transfer == nil {
		return "", fmt.Errorf("transfer is required")
	}

	slog.Info("initiating transfer (placeholder)",
		"transferId", transfer.ID,
		"citizenId", transfer.CitizenID,
		"amount", transfer.Amount,
		"currency", transfer.Currency,
		"type", transfer.Type,
	)

	externalRef := fmt.Sprintf("EXT-%s-%d", transfer.ID[:8], time.Now().UnixNano()%100000)

	slog.Info("transfer initiated (placeholder)",
		"transferId", transfer.ID,
		"externalRef", externalRef,
	)

	return externalRef, nil
}

// GetTransferStatus returns the status of a transfer by its external reference.
// Placeholder: always returns "completed".
func (b *BankAdapter) GetTransferStatus(ctx context.Context, externalRef string) (string, error) {
	if externalRef == "" {
		return "", fmt.Errorf("external reference is required")
	}

	slog.Info("checking transfer status (placeholder)",
		"externalRef", externalRef,
	)

	// Placeholder: all transfers are "completed".
	return "completed", nil
}
