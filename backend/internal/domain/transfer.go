package domain

import "time"

type Transfer struct {
	ID              string     `json:"id" db:"id"`
	CitizenID       string     `json:"citizenId" db:"citizen_id"`
	BankAccountID   string     `json:"bankAccountId,omitempty" db:"bank_account_id"`
	Amount          int64      `json:"amount" db:"amount"`
	Currency        string     `json:"currency" db:"currency"`
	Type            string     `json:"type" db:"type"` // savings, fsb_contribution, remittance, withdrawal
	SourceProvider  string     `json:"sourceProvider,omitempty" db:"source_provider"`
	SourceReference string     `json:"sourceReference,omitempty" db:"source_reference"`
	Status          string     `json:"status" db:"status"`
	ExternalRef     string     `json:"externalRef,omitempty" db:"external_ref"`
	FailureReason   string     `json:"failureReason,omitempty" db:"failure_reason"`
	CompletedAt     *time.Time `json:"completedAt,omitempty" db:"completed_at"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
}

type BankAccount struct {
	ID            string     `json:"id" db:"id"`
	CitizenID     string     `json:"citizenId" db:"citizen_id"`
	BankCode      string     `json:"bankCode" db:"bank_code"`
	BankName      string     `json:"bankName" db:"bank_name"`
	AccountNumber string     `json:"accountNumber" db:"account_number"`
	IBAN          string     `json:"iban,omitempty" db:"iban"`
	Status        string     `json:"status" db:"status"`
	OpenedAt      *time.Time `json:"openedAt,omitempty" db:"opened_at"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
}

type TransferFilter struct {
	CitizenID string
	Type      string
	Status    string
	Page      int
	PageSize  int
}
