package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/digikeys/backend/internal/domain"
)

type CommunicationRepo struct {
	pool *pgxpool.Pool
}

func NewCommunicationRepo(pool *pgxpool.Pool) *CommunicationRepo {
	return &CommunicationRepo{pool: pool}
}

func (r *CommunicationRepo) Create(ctx context.Context, c *domain.Communication) error {
	query := `
		INSERT INTO communications (
			id, embassy_id, sender_id, type, channel, subject, body,
			target_country, recipient_count, sent_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.pool.Exec(ctx, query,
		c.ID, c.EmbassyID, c.SenderID, c.Type, c.Channel, c.Subject, c.Body,
		c.TargetCountry, c.RecipientCount, c.SentAt, c.CreatedAt,
	)
	return err
}

func (r *CommunicationRepo) GetByID(ctx context.Context, id string) (*domain.Communication, error) {
	query := `SELECT ` + communicationColumns() + ` FROM communications WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	c, err := scanCommunication(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

func (r *CommunicationRepo) ListByEmbassy(ctx context.Context, embassyID string, page, pageSize int) ([]*domain.Communication, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM communications WHERE embassy_id = $1", embassyID).Scan(&total); err != nil {
		return nil, 0, err
	}

	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	query := fmt.Sprintf(
		"SELECT %s FROM communications WHERE embassy_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		communicationColumns(),
	)
	rows, err := r.pool.Query(ctx, query, embassyID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comms []*domain.Communication
	for rows.Next() {
		c, err := scanCommunicationRows(rows)
		if err != nil {
			return nil, 0, err
		}
		comms = append(comms, c)
	}
	return comms, total, rows.Err()
}

func communicationColumns() string {
	return `id, embassy_id, sender_id, type, channel, subject, body,
		target_country, recipient_count, sent_at, created_at`
}

func scanCommunication(row pgx.Row) (*domain.Communication, error) {
	var c domain.Communication
	err := row.Scan(
		&c.ID, &c.EmbassyID, &c.SenderID, &c.Type, &c.Channel, &c.Subject, &c.Body,
		&c.TargetCountry, &c.RecipientCount, &c.SentAt, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func scanCommunicationRows(rows pgx.Rows) (*domain.Communication, error) {
	var c domain.Communication
	err := rows.Scan(
		&c.ID, &c.EmbassyID, &c.SenderID, &c.Type, &c.Channel, &c.Subject, &c.Body,
		&c.TargetCountry, &c.RecipientCount, &c.SentAt, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
