package domain

import "time"

type Communication struct {
	ID             string     `json:"id" db:"id"`
	EmbassyID      string     `json:"embassyId" db:"embassy_id"`
	SenderID       string     `json:"senderId" db:"sender_id"`
	Type           string     `json:"type" db:"type"`
	Channel        string     `json:"channel" db:"channel"`
	Subject        string     `json:"subject" db:"subject"`
	Body           string     `json:"body" db:"body"`
	TargetCountry  string     `json:"targetCountry,omitempty" db:"target_country"`
	RecipientCount int        `json:"recipientCount" db:"recipient_count"`
	SentAt         *time.Time `json:"sentAt,omitempty" db:"sent_at"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
}
