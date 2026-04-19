package ports

import (
	"context"
	"io"

	"github.com/digikeys/backend/internal/domain"
)

type StorageService interface {
	Upload(ctx context.Context, bucket, key string, reader io.Reader, contentType string) (url string, err error)
	Download(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, bucket, key string) error
	GetURL(ctx context.Context, bucket, key string) (string, error)
}

type NotificationService interface {
	SendSMS(ctx context.Context, phone, message string) error
	SendEmail(ctx context.Context, to, subject, body string) error
}

type BiometricMatcher interface {
	Match(ctx context.Context, template1, template2 []byte) (float64, error)
	ExtractTemplate(ctx context.Context, rawImage []byte) ([]byte, error)
}

type MRZGenerator interface {
	GenerateTD1(citizen *domain.Citizen, card *domain.Card, embassy *domain.Embassy) (line1, line2, line3 string, err error)
}

type PrintingService interface {
	QueueForPrinting(ctx context.Context, card *domain.Card, citizen *domain.Citizen) error
	GetBatchStatus(ctx context.Context, batchID string) (string, error)
}

type BankingService interface {
	OpenAccount(ctx context.Context, citizen *domain.Citizen) (*domain.BankAccount, error)
	InitiateTransfer(ctx context.Context, transfer *domain.Transfer) (externalRef string, err error)
	GetTransferStatus(ctx context.Context, externalRef string) (string, error)
}
