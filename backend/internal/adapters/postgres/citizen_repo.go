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

type CitizenRepo struct {
	pool *pgxpool.Pool
}

func NewCitizenRepo(pool *pgxpool.Pool) *CitizenRepo {
	return &CitizenRepo{pool: pool}
}

func (r *CitizenRepo) Create(ctx context.Context, c *domain.Citizen) error {
	query := `
		INSERT INTO citizens (
			id, first_name, last_name, maiden_name, date_of_birth, place_of_birth,
			gender, nationality, national_id, unique_identifier, passport_number,
			phone, email, country_of_residence, city_of_residence, address_abroad,
			province_of_origin, commune_of_origin, embassy_id, registered_by,
			photo_key, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16,
			$17, $18, $19, $20,
			$21, $22, $23, $24
		)`
	_, err := r.pool.Exec(ctx, query,
		c.ID, c.FirstName, c.LastName, c.MaidenName, c.DateOfBirth, c.PlaceOfBirth,
		c.Gender, c.Nationality, c.NationalID, c.UniqueIdentifier, c.PassportNumber,
		c.Phone, c.Email, c.CountryOfResidence, c.CityOfResidence, c.AddressAbroad,
		c.ProvinceOfOrigin, c.CommuneOfOrigin, c.EmbassyID, c.RegisteredBy,
		c.PhotoKey, c.Status, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *CitizenRepo) GetByID(ctx context.Context, id string) (*domain.Citizen, error) {
	return r.getByField(ctx, "id", id)
}

func (r *CitizenRepo) GetByNationalID(ctx context.Context, nationalID string) (*domain.Citizen, error) {
	return r.getByField(ctx, "national_id", nationalID)
}

func (r *CitizenRepo) GetByUniqueIdentifier(ctx context.Context, uid string) (*domain.Citizen, error) {
	return r.getByField(ctx, "unique_identifier", uid)
}

func (r *CitizenRepo) getByField(ctx context.Context, field, value string) (*domain.Citizen, error) {
	query := `SELECT ` + citizenColumns() + ` FROM citizens WHERE ` + field + ` = $1`
	row := r.pool.QueryRow(ctx, query, value)
	c, err := scanCitizen(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

func (r *CitizenRepo) Search(ctx context.Context, filter domain.CitizenFilter) ([]*domain.Citizen, int, error) {
	return r.listInternal(ctx, filter, true)
}

func (r *CitizenRepo) List(ctx context.Context, filter domain.CitizenFilter) ([]*domain.Citizen, int, error) {
	return r.listInternal(ctx, filter, false)
}

func (r *CitizenRepo) listInternal(ctx context.Context, filter domain.CitizenFilter, searchMode bool) ([]*domain.Citizen, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.EmbassyID != "" {
		conditions = append(conditions, fmt.Sprintf("embassy_id = $%d", argIdx))
		args = append(args, filter.EmbassyID)
		argIdx++
	}
	if filter.CountryOfResidence != "" {
		conditions = append(conditions, fmt.Sprintf("country_of_residence = $%d", argIdx))
		args = append(args, filter.CountryOfResidence)
		argIdx++
	}
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.Query != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(first_name ILIKE $%d OR last_name ILIKE $%d OR national_id ILIKE $%d OR unique_identifier ILIKE $%d OR passport_number ILIKE $%d)",
			argIdx, argIdx, argIdx, argIdx, argIdx,
		))
		args = append(args, "%"+filter.Query+"%")
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count
	countQuery := "SELECT COUNT(*) FROM citizens " + where
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	offset := (page - 1) * pageSize

	query := fmt.Sprintf(
		"SELECT %s FROM citizens %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		citizenColumns(), where, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var citizens []*domain.Citizen
	for rows.Next() {
		c, err := scanCitizenRows(rows)
		if err != nil {
			return nil, 0, err
		}
		citizens = append(citizens, c)
	}
	return citizens, total, rows.Err()
}

func (r *CitizenRepo) Update(ctx context.Context, c *domain.Citizen) error {
	c.UpdatedAt = time.Now()
	query := `
		UPDATE citizens SET
			first_name = $2, last_name = $3, maiden_name = $4, date_of_birth = $5,
			place_of_birth = $6, gender = $7, nationality = $8, national_id = $9,
			unique_identifier = $10, passport_number = $11, phone = $12, email = $13,
			country_of_residence = $14, city_of_residence = $15, address_abroad = $16,
			province_of_origin = $17, commune_of_origin = $18, photo_key = $19,
			status = $20, updated_at = $21
		WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query,
		c.ID, c.FirstName, c.LastName, c.MaidenName, c.DateOfBirth,
		c.PlaceOfBirth, c.Gender, c.Nationality, c.NationalID,
		c.UniqueIdentifier, c.PassportNumber, c.Phone, c.Email,
		c.CountryOfResidence, c.CityOfResidence, c.AddressAbroad,
		c.ProvinceOfOrigin, c.CommuneOfOrigin, c.PhotoKey,
		c.Status, c.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func citizenColumns() string {
	return `id, first_name, last_name, maiden_name, date_of_birth, place_of_birth,
		gender, nationality, national_id, unique_identifier, passport_number,
		phone, email, country_of_residence, city_of_residence, address_abroad,
		province_of_origin, commune_of_origin, embassy_id, registered_by,
		photo_key, status, created_at, updated_at`
}

func scanCitizen(row pgx.Row) (*domain.Citizen, error) {
	var c domain.Citizen
	err := row.Scan(
		&c.ID, &c.FirstName, &c.LastName, &c.MaidenName, &c.DateOfBirth, &c.PlaceOfBirth,
		&c.Gender, &c.Nationality, &c.NationalID, &c.UniqueIdentifier, &c.PassportNumber,
		&c.Phone, &c.Email, &c.CountryOfResidence, &c.CityOfResidence, &c.AddressAbroad,
		&c.ProvinceOfOrigin, &c.CommuneOfOrigin, &c.EmbassyID, &c.RegisteredBy,
		&c.PhotoKey, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func scanCitizenRows(rows pgx.Rows) (*domain.Citizen, error) {
	var c domain.Citizen
	err := rows.Scan(
		&c.ID, &c.FirstName, &c.LastName, &c.MaidenName, &c.DateOfBirth, &c.PlaceOfBirth,
		&c.Gender, &c.Nationality, &c.NationalID, &c.UniqueIdentifier, &c.PassportNumber,
		&c.Phone, &c.Email, &c.CountryOfResidence, &c.CityOfResidence, &c.AddressAbroad,
		&c.ProvinceOfOrigin, &c.CommuneOfOrigin, &c.EmbassyID, &c.RegisteredBy,
		&c.PhotoKey, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func normalizePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
