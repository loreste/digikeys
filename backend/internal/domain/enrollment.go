package domain

import "time"

type Enrollment struct {
	ID           string                 `json:"id" db:"id"`
	CitizenID    string                 `json:"citizenId" db:"citizen_id"`
	EmbassyID    string                 `json:"embassyId" db:"embassy_id"`
	AgentID      string                 `json:"agentId" db:"agent_id"`
	TeamID       string                 `json:"teamId,omitempty" db:"team_id"`
	Latitude     float64                `json:"latitude" db:"latitude"`
	Longitude    float64                `json:"longitude" db:"longitude"`
	LocationName string                 `json:"locationName,omitempty" db:"location_name"`
	OfflineData  map[string]interface{} `json:"offlineData,omitempty" db:"offline_data"`
	SyncStatus   string                 `json:"syncStatus" db:"sync_status"`
	SyncedAt     *time.Time             `json:"syncedAt,omitempty" db:"synced_at"`
	ReviewStatus string                 `json:"reviewStatus" db:"review_status"`
	ReviewedBy   string                 `json:"reviewedBy,omitempty" db:"reviewed_by"`
	ReviewNotes  string                 `json:"reviewNotes,omitempty" db:"review_notes"`
	ReviewedAt   *time.Time             `json:"reviewedAt,omitempty" db:"reviewed_at"`
	EnrolledAt   time.Time              `json:"enrolledAt" db:"enrolled_at"`
	CreatedAt    time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time              `json:"updatedAt" db:"updated_at"`
}

type EnrollmentFilter struct {
	EmbassyID    string
	AgentID      string
	SyncStatus   string
	ReviewStatus string
	Page         int
	PageSize     int
}
