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

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (
			id, email, phone, password_hash, role, embassy_id,
			first_name, last_name, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.pool.Exec(ctx, query,
		u.ID, u.Email, u.Phone, u.PasswordHash, u.Role, nullIfEmpty(u.EmbassyID),
		u.FirstName, u.LastName, u.Status, u.CreatedAt, u.UpdatedAt,
	)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT ` + userColumns() + ` FROM users WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	u, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT ` + userColumns() + ` FROM users WHERE email = $1`
	row := r.pool.QueryRow(ctx, query, email)
	u, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) List(ctx context.Context, role domain.UserRole, embassyID string, page, pageSize int) ([]*domain.User, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if role != "" {
		conditions = append(conditions, fmt.Sprintf("role = $%d", argIdx))
		args = append(args, string(role))
		argIdx++
	}
	if embassyID != "" {
		conditions = append(conditions, fmt.Sprintf("embassy_id = $%d", argIdx))
		args = append(args, embassyID)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	query := fmt.Sprintf(
		"SELECT %s FROM users %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		userColumns(), where, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u, err := scanUserRows(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	u.UpdatedAt = time.Now()
	query := `
		UPDATE users SET
			email = $2, phone = $3, role = $4, embassy_id = $5,
			first_name = $6, last_name = $7, status = $8, updated_at = $9
		WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query,
		u.ID, u.Email, u.Phone, u.Role, nullIfEmpty(u.EmbassyID),
		u.FirstName, u.LastName, u.Status, u.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	query := `UPDATE users SET password_hash = $2, updated_at = $3 WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id, passwordHash, time.Now())
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func userColumns() string {
	return `id, email, phone, password_hash, role, embassy_id,
		first_name, last_name, status, created_at, updated_at`
}

func scanUser(row pgx.Row) (*domain.User, error) {
	var u domain.User
	var embassyID *string
	err := row.Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.Role, &embassyID,
		&u.FirstName, &u.LastName, &u.Status, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if embassyID != nil {
		u.EmbassyID = *embassyID
	}
	return &u, nil
}

func scanUserRows(rows pgx.Rows) (*domain.User, error) {
	var u domain.User
	var embassyID *string
	err := rows.Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.Role, &embassyID,
		&u.FirstName, &u.LastName, &u.Status, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if embassyID != nil {
		u.EmbassyID = *embassyID
	}
	return &u, nil
}
