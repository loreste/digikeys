package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/digikeys/backend/internal/domain"
)

type TransferRepo struct {
	pool *pgxpool.Pool
}

func NewTransferRepo(pool *pgxpool.Pool) *TransferRepo {
	return &TransferRepo{pool: pool}
}

func (r *TransferRepo) Create(ctx context.Context, t *domain.Transfer) error {
	query := `
		INSERT INTO transfers (
			id, citizen_id, bank_account_id, amount, currency, type,
			source_provider, source_reference, status, external_ref,
			failure_reason, completed_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.pool.Exec(ctx, query,
		t.ID, t.CitizenID, nullIfEmpty(t.BankAccountID), t.Amount, t.Currency, t.Type,
		t.SourceProvider, t.SourceReference, t.Status, t.ExternalRef,
		t.FailureReason, t.CompletedAt, t.CreatedAt,
	)
	return err
}

func (r *TransferRepo) GetByID(ctx context.Context, id string) (*domain.Transfer, error) {
	query := `SELECT ` + transferColumns() + ` FROM transfers WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	t, err := scanTransfer(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return t, nil
}

func (r *TransferRepo) GetByExternalRef(ctx context.Context, ref string) (*domain.Transfer, error) {
	query := `SELECT ` + transferColumns() + ` FROM transfers WHERE external_ref = $1`
	row := r.pool.QueryRow(ctx, query, ref)
	t, err := scanTransfer(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return t, nil
}

func (r *TransferRepo) List(ctx context.Context, filter domain.TransferFilter) ([]*domain.Transfer, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.CitizenID != "" {
		conditions = append(conditions, fmt.Sprintf("citizen_id = $%d", argIdx))
		args = append(args, filter.CitizenID)
		argIdx++
	}
	if filter.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, filter.Type)
		argIdx++
	}
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM transfers "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	offset := (page - 1) * pageSize

	query := fmt.Sprintf(
		"SELECT %s FROM transfers %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		transferColumns(), where, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transfers []*domain.Transfer
	for rows.Next() {
		t, err := scanTransferRows(rows)
		if err != nil {
			return nil, 0, err
		}
		transfers = append(transfers, t)
	}
	return transfers, total, rows.Err()
}

func (r *TransferRepo) Update(ctx context.Context, t *domain.Transfer) error {
	query := `
		UPDATE transfers SET
			citizen_id = $2, bank_account_id = $3, amount = $4, currency = $5,
			type = $6, source_provider = $7, source_reference = $8,
			status = $9, external_ref = $10, failure_reason = $11,
			completed_at = $12
		WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query,
		t.ID, t.CitizenID, nullIfEmpty(t.BankAccountID), t.Amount, t.Currency,
		t.Type, t.SourceProvider, t.SourceReference,
		t.Status, t.ExternalRef, t.FailureReason,
		t.CompletedAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func transferColumns() string {
	return `id, citizen_id, bank_account_id, amount, currency, type,
		source_provider, source_reference, status, external_ref,
		failure_reason, completed_at, created_at`
}

func scanTransfer(row pgx.Row) (*domain.Transfer, error) {
	var t domain.Transfer
	var bankAccountID *string
	err := row.Scan(
		&t.ID, &t.CitizenID, &bankAccountID, &t.Amount, &t.Currency, &t.Type,
		&t.SourceProvider, &t.SourceReference, &t.Status, &t.ExternalRef,
		&t.FailureReason, &t.CompletedAt, &t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if bankAccountID != nil {
		t.BankAccountID = *bankAccountID
	}
	return &t, nil
}

func scanTransferRows(rows pgx.Rows) (*domain.Transfer, error) {
	var t domain.Transfer
	var bankAccountID *string
	err := rows.Scan(
		&t.ID, &t.CitizenID, &bankAccountID, &t.Amount, &t.Currency, &t.Type,
		&t.SourceProvider, &t.SourceReference, &t.Status, &t.ExternalRef,
		&t.FailureReason, &t.CompletedAt, &t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if bankAccountID != nil {
		t.BankAccountID = *bankAccountID
	}
	return &t, nil
}
