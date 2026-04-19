package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/digikeys/backend/internal/domain"
)

type EmbassyRepo struct {
	pool *pgxpool.Pool
}

func NewEmbassyRepo(pool *pgxpool.Pool) *EmbassyRepo {
	return &EmbassyRepo{pool: pool}
}

func (r *EmbassyRepo) Create(ctx context.Context, e *domain.Embassy) error {
	query := `
		INSERT INTO embassies (
			id, country_code, name, city, address, phone, email,
			consul_name, card_prefix, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.pool.Exec(ctx, query,
		e.ID, e.CountryCode, e.Name, e.City, e.Address, e.Phone, e.Email,
		e.ConsulName, e.CardPrefix, e.Status, e.CreatedAt, e.UpdatedAt,
	)
	return err
}

func (r *EmbassyRepo) GetByID(ctx context.Context, id string) (*domain.Embassy, error) {
	query := `SELECT ` + embassyColumns() + ` FROM embassies WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	e, err := scanEmbassy(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return e, nil
}

func (r *EmbassyRepo) GetByCountryCode(ctx context.Context, code string) (*domain.Embassy, error) {
	query := `SELECT ` + embassyColumns() + ` FROM embassies WHERE country_code = $1`
	row := r.pool.QueryRow(ctx, query, code)
	e, err := scanEmbassy(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return e, nil
}

func (r *EmbassyRepo) List(ctx context.Context, page, pageSize int) ([]*domain.Embassy, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM embassies").Scan(&total); err != nil {
		return nil, 0, err
	}

	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	query := fmt.Sprintf(
		"SELECT %s FROM embassies ORDER BY name ASC LIMIT $1 OFFSET $2",
		embassyColumns(),
	)
	rows, err := r.pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var embassies []*domain.Embassy
	for rows.Next() {
		e, err := scanEmbassyRows(rows)
		if err != nil {
			return nil, 0, err
		}
		embassies = append(embassies, e)
	}
	return embassies, total, rows.Err()
}

func (r *EmbassyRepo) Update(ctx context.Context, e *domain.Embassy) error {
	e.UpdatedAt = time.Now()
	query := `
		UPDATE embassies SET
			country_code = $2, name = $3, city = $4, address = $5,
			phone = $6, email = $7, consul_name = $8, card_prefix = $9,
			status = $10, updated_at = $11
		WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query,
		e.ID, e.CountryCode, e.Name, e.City, e.Address,
		e.Phone, e.Email, e.ConsulName, e.CardPrefix,
		e.Status, e.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func embassyColumns() string {
	return `id, country_code, name, city, address, phone, email,
		consul_name, card_prefix, status, created_at, updated_at`
}

func scanEmbassy(row pgx.Row) (*domain.Embassy, error) {
	var e domain.Embassy
	err := row.Scan(
		&e.ID, &e.CountryCode, &e.Name, &e.City, &e.Address, &e.Phone, &e.Email,
		&e.ConsulName, &e.CardPrefix, &e.Status, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func scanEmbassyRows(rows pgx.Rows) (*domain.Embassy, error) {
	var e domain.Embassy
	err := rows.Scan(
		&e.ID, &e.CountryCode, &e.Name, &e.City, &e.Address, &e.Phone, &e.Email,
		&e.ConsulName, &e.CardPrefix, &e.Status, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
