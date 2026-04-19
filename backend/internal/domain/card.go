package domain

import "time"

type CardStatus string

const (
	CardStatusPending   CardStatus = "pending"
	CardStatusApproved  CardStatus = "approved"
	CardStatusPrinting  CardStatus = "printing"
	CardStatusPrinted   CardStatus = "printed"
	CardStatusDelivered CardStatus = "delivered"
	CardStatusActive    CardStatus = "active"
	CardStatusSuspended CardStatus = "suspended"
	CardStatusRevoked   CardStatus = "revoked"
	CardStatusExpired   CardStatus = "expired"
)

type Card struct {
	ID             string     `json:"id" db:"id"`
	CitizenID      string     `json:"citizenId" db:"citizen_id"`
	CardNumber     string     `json:"cardNumber" db:"card_number"`
	MRZLine1       string     `json:"mrzLine1" db:"mrz_line1"`
	MRZLine2       string     `json:"mrzLine2" db:"mrz_line2"`
	MRZLine3       string     `json:"mrzLine3" db:"mrz_line3"`
	EmbassyID      string     `json:"embassyId" db:"embassy_id"`
	IssuedBy       string     `json:"issuedBy,omitempty" db:"issued_by"`
	IssuedAt       *time.Time `json:"issuedAt,omitempty" db:"issued_at"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty" db:"expires_at"`
	Status         string     `json:"status" db:"status"`
	PrintBatchID   string     `json:"printBatchId,omitempty" db:"print_batch_id"`
	PrintedAt      *time.Time `json:"printedAt,omitempty" db:"printed_at"`
	DeliveredAt    *time.Time `json:"deliveredAt,omitempty" db:"delivered_at"`
	PreviousCardID string     `json:"previousCardId,omitempty" db:"previous_card_id"`
	RenewalReason  string     `json:"renewalReason,omitempty" db:"renewal_reason"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
}

type CardFilter struct {
	EmbassyID string
	Status    string
	CitizenID string
	Page      int
	PageSize  int
}
