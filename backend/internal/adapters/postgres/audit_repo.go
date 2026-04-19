package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditLogRepo struct {
	pool *pgxpool.Pool
}

func NewAuditLogRepo(pool *pgxpool.Pool) *AuditLogRepo {
	return &AuditLogRepo{pool: pool}
}

func (r *AuditLogRepo) Create(ctx context.Context, action, entityType, entityID, userID, details string) error {
	query := `
		INSERT INTO audit_logs (id, action, entity_type, entity_id, user_id, details, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query,
		uuid.New().String(), action, entityType, entityID, userID, details, time.Now(),
	)
	return err
}

func (r *AuditLogRepo) List(ctx context.Context, entityType, entityID string, page, pageSize int) ([]map[string]interface{}, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if entityType != "" {
		conditions = append(conditions, fmt.Sprintf("entity_type = $%d", argIdx))
		args = append(args, entityType)
		argIdx++
	}
	if entityID != "" {
		conditions = append(conditions, fmt.Sprintf("entity_id = $%d", argIdx))
		args = append(args, entityID)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM audit_logs "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	query := fmt.Sprintf(
		"SELECT id, action, entity_type, entity_id, user_id, details, created_at FROM audit_logs %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		where, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var id, action, eType, eID, userID string
		var details *string
		var createdAt time.Time

		if err := rows.Scan(&id, &action, &eType, &eID, &userID, &details, &createdAt); err != nil {
			return nil, 0, err
		}

		entry := map[string]interface{}{
			"id":          id,
			"action":      action,
			"entityType":  eType,
			"entityId":    eID,
			"userId":      userID,
			"createdAt":   createdAt,
		}
		if details != nil {
			entry["details"] = *details
		}
		logs = append(logs, entry)
	}
	return logs, total, rows.Err()
}
