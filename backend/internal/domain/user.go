package domain

import "time"

type UserRole string

const (
	UserRoleSuperAdmin     UserRole = "super_admin"
	UserRoleEmbassyAdmin   UserRole = "embassy_admin"
	UserRoleEnrollmentAgent UserRole = "enrollment_agent"
	UserRolePrintOperator  UserRole = "print_operator"
	UserRoleBankAgent      UserRole = "bank_agent"
	UserRoleVerifier       UserRole = "verifier"
	UserRoleReadonly       UserRole = "readonly"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Phone        string    `json:"phone,omitempty" db:"phone"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         UserRole  `json:"role" db:"role"`
	EmbassyID    string    `json:"embassyId,omitempty" db:"embassy_id"`
	FirstName    string    `json:"firstName" db:"first_name"`
	LastName     string    `json:"lastName" db:"last_name"`
	Status       string    `json:"status" db:"status"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}
