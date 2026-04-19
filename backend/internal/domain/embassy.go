package domain

import "time"

type Embassy struct {
	ID          string    `json:"id" db:"id"`
	CountryCode string    `json:"countryCode" db:"country_code"`
	Name        string    `json:"name" db:"name"`
	City        string    `json:"city" db:"city"`
	Address     string    `json:"address,omitempty" db:"address"`
	Phone       string    `json:"phone,omitempty" db:"phone"`
	Email       string    `json:"email,omitempty" db:"email"`
	ConsulName  string    `json:"consulName,omitempty" db:"consul_name"`
	CardPrefix  string    `json:"cardPrefix" db:"card_prefix"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}
