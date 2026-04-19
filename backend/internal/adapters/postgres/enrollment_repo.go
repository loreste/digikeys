package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/digikeys/backend/internal/domain"
)

type EnrollmentRepo struct {
	pool *pgxpool.Pool
}

func NewEnrollmentRepo(pool *pgxpool.Pool) *EnrollmentRepo {
	return &EnrollmentRepo{pool: pool}
}

func (r *EnrollmentRepo) Create(ctx context.Context, e *domain.Enrollment) error {
	offlineJSON, err := json.Marshal(e.OfflineData)
	if err != nil {
		return fmt.Errorf("marshal offline data: %w", err)
	}

	query := `
		INSERT INTO enrollments (
			id, citizen_id, embassy_id, agent_id, team_id,
			latitude, longitude, location_name, offline_data,
			sync_status, synced_at, review_status, reviewed_by,
			review_notes, reviewed_at, enrolled_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12, $13,
			$14, $15, $16, $17, $18
		)`
	_, err = r.pool.Exec(ctx, query,
		e.ID, e.CitizenID, e.EmbassyID, e.AgentID, e.TeamID,
		e.Latitude, e.Longitude, e.LocationName, offlineJSON,
		e.SyncStatus, e.SyncedAt, e.ReviewStatus, nullIfEmpty(e.ReviewedBy),
		e.ReviewNotes, e.ReviewedAt, e.EnrolledAt, e.CreatedAt, e.UpdatedAt,
	)
	return err
}

func (r *EnrollmentRepo) GetByID(ctx context.Context, id string) (*domain.Enrollment, error) {
	query := `SELECT ` + enrollmentColumns() + ` FROM enrollments WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	e, err := scanEnrollment(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return e, nil
}

func (r *EnrollmentRepo) ListByAgent(ctx context.Context, agentID string, page, pageSize int) ([]*domain.Enrollment, int, error) {
	return r.List(ctx, domain.EnrollmentFilter{
		AgentID:  agentID,
		Page:     page,
		PageSize: pageSize,
	})
}

func (r *EnrollmentRepo) ListByEmbassy(ctx context.Context, embassyID string, page, pageSize int) ([]*domain.Enrollment, int, error) {
	return r.List(ctx, domain.EnrollmentFilter{
		EmbassyID: embassyID,
		Page:      page,
		PageSize:  pageSize,
	})
}

func (r *EnrollmentRepo) List(ctx context.Context, filter domain.EnrollmentFilter) ([]*domain.Enrollment, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.EmbassyID != "" {
		conditions = append(conditions, fmt.Sprintf("embassy_id = $%d", argIdx))
		args = append(args, filter.EmbassyID)
		argIdx++
	}
	if filter.AgentID != "" {
		conditions = append(conditions, fmt.Sprintf("agent_id = $%d", argIdx))
		args = append(args, filter.AgentID)
		argIdx++
	}
	if filter.SyncStatus != "" {
		conditions = append(conditions, fmt.Sprintf("sync_status = $%d", argIdx))
		args = append(args, filter.SyncStatus)
		argIdx++
	}
	if filter.ReviewStatus != "" {
		conditions = append(conditions, fmt.Sprintf("review_status = $%d", argIdx))
		args = append(args, filter.ReviewStatus)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM enrollments "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	offset := (page - 1) * pageSize

	query := fmt.Sprintf(
		"SELECT %s FROM enrollments %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		enrollmentColumns(), where, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var enrollments []*domain.Enrollment
	for rows.Next() {
		e, err := scanEnrollmentRows(rows)
		if err != nil {
			return nil, 0, err
		}
		enrollments = append(enrollments, e)
	}
	return enrollments, total, rows.Err()
}

func (r *EnrollmentRepo) Update(ctx context.Context, e *domain.Enrollment) error {
	e.UpdatedAt = time.Now()
	offlineJSON, err := json.Marshal(e.OfflineData)
	if err != nil {
		return fmt.Errorf("marshal offline data: %w", err)
	}

	query := `
		UPDATE enrollments SET
			citizen_id = $2, embassy_id = $3, agent_id = $4, team_id = $5,
			latitude = $6, longitude = $7, location_name = $8, offline_data = $9,
			sync_status = $10, synced_at = $11, review_status = $12, reviewed_by = $13,
			review_notes = $14, reviewed_at = $15, updated_at = $16
		WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query,
		e.ID, e.CitizenID, e.EmbassyID, e.AgentID, e.TeamID,
		e.Latitude, e.Longitude, e.LocationName, offlineJSON,
		e.SyncStatus, e.SyncedAt, e.ReviewStatus, nullIfEmpty(e.ReviewedBy),
		e.ReviewNotes, e.ReviewedAt, e.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *EnrollmentRepo) UpdateSyncStatus(ctx context.Context, id string, status string, syncedAt *time.Time) error {
	query := `UPDATE enrollments SET sync_status = $2, synced_at = $3, updated_at = $4 WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id, status, syncedAt, time.Now())
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *EnrollmentRepo) UpdateReviewStatus(ctx context.Context, id, status, reviewedBy, notes string, reviewedAt *time.Time) error {
	query := `
		UPDATE enrollments SET
			review_status = $2, reviewed_by = $3, review_notes = $4,
			reviewed_at = $5, updated_at = $6
		WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id, status, nullIfEmpty(reviewedBy), notes, reviewedAt, time.Now())
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func enrollmentColumns() string {
	return `id, citizen_id, embassy_id, agent_id, team_id,
		latitude, longitude, location_name, offline_data,
		sync_status, synced_at, review_status, reviewed_by,
		review_notes, reviewed_at, enrolled_at, created_at, updated_at`
}

func scanEnrollment(row pgx.Row) (*domain.Enrollment, error) {
	var e domain.Enrollment
	var offlineJSON []byte
	var reviewedBy *string
	err := row.Scan(
		&e.ID, &e.CitizenID, &e.EmbassyID, &e.AgentID, &e.TeamID,
		&e.Latitude, &e.Longitude, &e.LocationName, &offlineJSON,
		&e.SyncStatus, &e.SyncedAt, &e.ReviewStatus, &reviewedBy,
		&e.ReviewNotes, &e.ReviewedAt, &e.EnrolledAt, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if reviewedBy != nil {
		e.ReviewedBy = *reviewedBy
	}
	if offlineJSON != nil {
		_ = json.Unmarshal(offlineJSON, &e.OfflineData)
	}
	return &e, nil
}

func scanEnrollmentRows(rows pgx.Rows) (*domain.Enrollment, error) {
	var e domain.Enrollment
	var offlineJSON []byte
	var reviewedBy *string
	err := rows.Scan(
		&e.ID, &e.CitizenID, &e.EmbassyID, &e.AgentID, &e.TeamID,
		&e.Latitude, &e.Longitude, &e.LocationName, &offlineJSON,
		&e.SyncStatus, &e.SyncedAt, &e.ReviewStatus, &reviewedBy,
		&e.ReviewNotes, &e.ReviewedAt, &e.EnrolledAt, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if reviewedBy != nil {
		e.ReviewedBy = *reviewedBy
	}
	if offlineJSON != nil {
		_ = json.Unmarshal(offlineJSON, &e.OfflineData)
	}
	return &e, nil
}
