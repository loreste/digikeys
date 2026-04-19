package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/digikeys/backend/internal/domain"
)

type CardRepo struct {
	pool *pgxpool.Pool
}

func NewCardRepo(pool *pgxpool.Pool) *CardRepo {
	return &CardRepo{pool: pool}
}

func (r *CardRepo) Create(ctx context.Context, c *domain.Card) error {
	query := `
		INSERT INTO cards (
			id, citizen_id, card_number, mrz_line1, mrz_line2, mrz_line3,
			embassy_id, issued_by, issued_at, expires_at, status,
			print_batch_id, printed_at, delivered_at,
			previous_card_id, renewal_reason, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14,
			$15, $16, $17, $18
		)`
	_, err := r.pool.Exec(ctx, query,
		c.ID, c.CitizenID, c.CardNumber, c.MRZLine1, c.MRZLine2, c.MRZLine3,
		c.EmbassyID, nullIfEmpty(c.IssuedBy), c.IssuedAt, c.ExpiresAt, c.Status,
		nullIfEmpty(c.PrintBatchID), c.PrintedAt, c.DeliveredAt,
		nullIfEmpty(c.PreviousCardID), c.RenewalReason, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *CardRepo) GetByID(ctx context.Context, id string) (*domain.Card, error) {
	query := `SELECT ` + cardColumns() + ` FROM cards WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	c, err := scanCard(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

func (r *CardRepo) GetByCardNumber(ctx context.Context, cardNumber string) (*domain.Card, error) {
	query := `SELECT ` + cardColumns() + ` FROM cards WHERE card_number = $1`
	row := r.pool.QueryRow(ctx, query, cardNumber)
	c, err := scanCard(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

func (r *CardRepo) GetByCitizenID(ctx context.Context, citizenID string) ([]*domain.Card, error) {
	query := `SELECT ` + cardColumns() + ` FROM cards WHERE citizen_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, citizenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*domain.Card
	for rows.Next() {
		c, err := scanCardRows(rows)
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, rows.Err()
}

func (r *CardRepo) List(ctx context.Context, filter domain.CardFilter) ([]*domain.Card, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.EmbassyID != "" {
		conditions = append(conditions, fmt.Sprintf("embassy_id = $%d", argIdx))
		args = append(args, filter.EmbassyID)
		argIdx++
	}
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.CitizenID != "" {
		conditions = append(conditions, fmt.Sprintf("citizen_id = $%d", argIdx))
		args = append(args, filter.CitizenID)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM cards "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	offset := (page - 1) * pageSize

	query := fmt.Sprintf(
		"SELECT %s FROM cards %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		cardColumns(), where, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var cards []*domain.Card
	for rows.Next() {
		c, err := scanCardRows(rows)
		if err != nil {
			return nil, 0, err
		}
		cards = append(cards, c)
	}
	return cards, total, rows.Err()
}

func (r *CardRepo) Update(ctx context.Context, c *domain.Card) error {
	c.UpdatedAt = time.Now()
	query := `
		UPDATE cards SET
			citizen_id = $2, card_number = $3, mrz_line1 = $4, mrz_line2 = $5,
			mrz_line3 = $6, embassy_id = $7, issued_by = $8, issued_at = $9,
			expires_at = $10, status = $11, print_batch_id = $12, printed_at = $13,
			delivered_at = $14, previous_card_id = $15, renewal_reason = $16,
			updated_at = $17
		WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query,
		c.ID, c.CitizenID, c.CardNumber, c.MRZLine1, c.MRZLine2,
		c.MRZLine3, c.EmbassyID, nullIfEmpty(c.IssuedBy), c.IssuedAt,
		c.ExpiresAt, c.Status, nullIfEmpty(c.PrintBatchID), c.PrintedAt,
		c.DeliveredAt, nullIfEmpty(c.PreviousCardID), c.RenewalReason,
		c.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CardRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE cards SET status = $2, updated_at = $3 WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id, status, time.Now())
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CardRepo) GetNextSequence(ctx context.Context, embassyID string, year int) (int, error) {
	query := `
		INSERT INTO card_sequences (embassy_id, year, last_sequence)
		VALUES ($1, $2, 1)
		ON CONFLICT (embassy_id, year)
		DO UPDATE SET last_sequence = card_sequences.last_sequence + 1
		RETURNING last_sequence`
	var seq int
	err := r.pool.QueryRow(ctx, query, embassyID, year).Scan(&seq)
	return seq, err
}

func cardColumns() string {
	return `id, citizen_id, card_number, mrz_line1, mrz_line2, mrz_line3,
		embassy_id, issued_by, issued_at, expires_at, status,
		print_batch_id, printed_at, delivered_at,
		previous_card_id, renewal_reason, created_at, updated_at`
}

func scanCard(row pgx.Row) (*domain.Card, error) {
	var c domain.Card
	var issuedBy, printBatchID, previousCardID *string
	err := row.Scan(
		&c.ID, &c.CitizenID, &c.CardNumber, &c.MRZLine1, &c.MRZLine2, &c.MRZLine3,
		&c.EmbassyID, &issuedBy, &c.IssuedAt, &c.ExpiresAt, &c.Status,
		&printBatchID, &c.PrintedAt, &c.DeliveredAt,
		&previousCardID, &c.RenewalReason, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if issuedBy != nil {
		c.IssuedBy = *issuedBy
	}
	if printBatchID != nil {
		c.PrintBatchID = *printBatchID
	}
	if previousCardID != nil {
		c.PreviousCardID = *previousCardID
	}
	return &c, nil
}

func scanCardRows(rows pgx.Rows) (*domain.Card, error) {
	var c domain.Card
	var issuedBy, printBatchID, previousCardID *string
	err := rows.Scan(
		&c.ID, &c.CitizenID, &c.CardNumber, &c.MRZLine1, &c.MRZLine2, &c.MRZLine3,
		&c.EmbassyID, &issuedBy, &c.IssuedAt, &c.ExpiresAt, &c.Status,
		&printBatchID, &c.PrintedAt, &c.DeliveredAt,
		&previousCardID, &c.RenewalReason, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if issuedBy != nil {
		c.IssuedBy = *issuedBy
	}
	if printBatchID != nil {
		c.PrintBatchID = *printBatchID
	}
	if previousCardID != nil {
		c.PreviousCardID = *previousCardID
	}
	return &c, nil
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
