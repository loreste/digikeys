package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/digikeys/backend/internal/domain"
)

type BankAccountRepo struct {
	pool *pgxpool.Pool
}

func NewBankAccountRepo(pool *pgxpool.Pool) *BankAccountRepo {
	return &BankAccountRepo{pool: pool}
}

func (r *BankAccountRepo) Create(ctx context.Context, a *domain.BankAccount) error {
	query := `
		INSERT INTO bank_accounts (
			id, citizen_id, bank_code, bank_name, account_number, iban,
			status, opened_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.pool.Exec(ctx, query,
		a.ID, a.CitizenID, a.BankCode, a.BankName, a.AccountNumber, a.IBAN,
		a.Status, a.OpenedAt, a.CreatedAt, a.UpdatedAt,
	)
	return err
}

func (r *BankAccountRepo) GetByID(ctx context.Context, id string) (*domain.BankAccount, error) {
	query := `SELECT ` + bankAccountColumns() + ` FROM bank_accounts WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	a, err := scanBankAccount(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return a, nil
}

func (r *BankAccountRepo) GetByCitizenID(ctx context.Context, citizenID string) ([]*domain.BankAccount, error) {
	query := `SELECT ` + bankAccountColumns() + ` FROM bank_accounts WHERE citizen_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, citizenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*domain.BankAccount
	for rows.Next() {
		a, err := scanBankAccountRows(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

func (r *BankAccountRepo) Update(ctx context.Context, a *domain.BankAccount) error {
	a.UpdatedAt = time.Now()
	query := `
		UPDATE bank_accounts SET
			bank_code = $2, bank_name = $3, account_number = $4,
			iban = $5, status = $6, opened_at = $7, updated_at = $8
		WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query,
		a.ID, a.BankCode, a.BankName, a.AccountNumber,
		a.IBAN, a.Status, a.OpenedAt, a.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func bankAccountColumns() string {
	return `id, citizen_id, bank_code, bank_name, account_number, iban,
		status, opened_at, created_at, updated_at`
}

func scanBankAccount(row pgx.Row) (*domain.BankAccount, error) {
	var a domain.BankAccount
	err := row.Scan(
		&a.ID, &a.CitizenID, &a.BankCode, &a.BankName, &a.AccountNumber, &a.IBAN,
		&a.Status, &a.OpenedAt, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func scanBankAccountRows(rows pgx.Rows) (*domain.BankAccount, error) {
	var a domain.BankAccount
	err := rows.Scan(
		&a.ID, &a.CitizenID, &a.BankCode, &a.BankName, &a.AccountNumber, &a.IBAN,
		&a.Status, &a.OpenedAt, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
