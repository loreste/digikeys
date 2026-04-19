package ports

import (
	"context"
	"time"

	"github.com/digikeys/backend/internal/domain"
)

type CitizenRepository interface {
	Create(ctx context.Context, citizen *domain.Citizen) error
	GetByID(ctx context.Context, id string) (*domain.Citizen, error)
	GetByNationalID(ctx context.Context, nationalID string) (*domain.Citizen, error)
	GetByUniqueIdentifier(ctx context.Context, uid string) (*domain.Citizen, error)
	Search(ctx context.Context, filter domain.CitizenFilter) ([]*domain.Citizen, int, error)
	Update(ctx context.Context, citizen *domain.Citizen) error
	List(ctx context.Context, filter domain.CitizenFilter) ([]*domain.Citizen, int, error)
}

type CardRepository interface {
	Create(ctx context.Context, card *domain.Card) error
	GetByID(ctx context.Context, id string) (*domain.Card, error)
	GetByCardNumber(ctx context.Context, cardNumber string) (*domain.Card, error)
	GetByCitizenID(ctx context.Context, citizenID string) ([]*domain.Card, error)
	List(ctx context.Context, filter domain.CardFilter) ([]*domain.Card, int, error)
	Update(ctx context.Context, card *domain.Card) error
	UpdateStatus(ctx context.Context, id string, status string) error
	GetNextSequence(ctx context.Context, embassyID string, year int) (int, error)
}

type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *domain.Enrollment) error
	GetByID(ctx context.Context, id string) (*domain.Enrollment, error)
	ListByAgent(ctx context.Context, agentID string, page, pageSize int) ([]*domain.Enrollment, int, error)
	ListByEmbassy(ctx context.Context, embassyID string, page, pageSize int) ([]*domain.Enrollment, int, error)
	List(ctx context.Context, filter domain.EnrollmentFilter) ([]*domain.Enrollment, int, error)
	Update(ctx context.Context, enrollment *domain.Enrollment) error
	UpdateSyncStatus(ctx context.Context, id string, status string, syncedAt *time.Time) error
	UpdateReviewStatus(ctx context.Context, id, status, reviewedBy, notes string, reviewedAt *time.Time) error
}

type BiometricRepository interface {
	Create(ctx context.Context, bio *domain.Biometric) error
	GetByCitizenID(ctx context.Context, citizenID string) (*domain.Biometric, error)
	Update(ctx context.Context, bio *domain.Biometric) error
}

type EmbassyRepository interface {
	Create(ctx context.Context, embassy *domain.Embassy) error
	GetByID(ctx context.Context, id string) (*domain.Embassy, error)
	GetByCountryCode(ctx context.Context, code string) (*domain.Embassy, error)
	List(ctx context.Context, page, pageSize int) ([]*domain.Embassy, int, error)
	Update(ctx context.Context, embassy *domain.Embassy) error
}

type TransferRepository interface {
	Create(ctx context.Context, transfer *domain.Transfer) error
	GetByID(ctx context.Context, id string) (*domain.Transfer, error)
	GetByExternalRef(ctx context.Context, ref string) (*domain.Transfer, error)
	List(ctx context.Context, filter domain.TransferFilter) ([]*domain.Transfer, int, error)
	Update(ctx context.Context, transfer *domain.Transfer) error
}

type BankAccountRepository interface {
	Create(ctx context.Context, account *domain.BankAccount) error
	GetByID(ctx context.Context, id string) (*domain.BankAccount, error)
	GetByCitizenID(ctx context.Context, citizenID string) ([]*domain.BankAccount, error)
	Update(ctx context.Context, account *domain.BankAccount) error
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, role domain.UserRole, embassyID string, page, pageSize int) ([]*domain.User, int, error)
	Update(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, id string, passwordHash string) error
}

type CommunicationRepository interface {
	Create(ctx context.Context, comm *domain.Communication) error
	GetByID(ctx context.Context, id string) (*domain.Communication, error)
	ListByEmbassy(ctx context.Context, embassyID string, page, pageSize int) ([]*domain.Communication, int, error)
}

type AuditLogRepository interface {
	Create(ctx context.Context, action, entityType, entityID, userID, details string) error
	List(ctx context.Context, entityType, entityID string, page, pageSize int) ([]map[string]interface{}, int, error)
}

type StatisticsQuerier interface {
	TotalCitizens(ctx context.Context, embassyID string) (int, error)
	CardsByStatus(ctx context.Context, embassyID string) (map[string]int, error)
	EnrollmentsByStatus(ctx context.Context, embassyID string) (map[string]int, error)
	TransfersTotal(ctx context.Context, embassyID string) (int64, error)
	FSBTotal(ctx context.Context, embassyID string) (int64, error)
	CitizensByCountry(ctx context.Context) (map[string]int, error)
}
