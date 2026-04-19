package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/digikeys/backend/internal/domain"
)

type BiometricRepo struct {
	pool *pgxpool.Pool
}

func NewBiometricRepo(pool *pgxpool.Pool) *BiometricRepo {
	return &BiometricRepo{pool: pool}
}

func (r *BiometricRepo) Create(ctx context.Context, b *domain.Biometric) error {
	query := `
		INSERT INTO biometrics (
			id, citizen_id,
			fp_right_thumb, fp_right_index, fp_left_thumb, fp_left_index,
			fp_right_thumb_quality, fp_right_index_quality,
			fp_left_thumb_quality, fp_left_index_quality,
			photo_key, photo_hash, signature_key, encryption_key_id,
			captured_by, capture_device, capture_location,
			captured_at, created_at
		) VALUES (
			$1, $2,
			$3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13, $14,
			$15, $16, $17,
			$18, $19
		)`
	_, err := r.pool.Exec(ctx, query,
		b.ID, b.CitizenID,
		b.FPRightThumb, b.FPRightIndex, b.FPLeftThumb, b.FPLeftIndex,
		b.FPRightThumbQuality, b.FPRightIndexQuality,
		b.FPLeftThumbQuality, b.FPLeftIndexQuality,
		b.PhotoKey, b.PhotoHash, b.SignatureKey, b.EncryptionKeyID,
		b.CapturedBy, b.CaptureDevice, b.CaptureLocation,
		b.CapturedAt, b.CreatedAt,
	)
	return err
}

func (r *BiometricRepo) GetByCitizenID(ctx context.Context, citizenID string) (*domain.Biometric, error) {
	query := `SELECT ` + biometricColumns() + ` FROM biometrics WHERE citizen_id = $1`
	row := r.pool.QueryRow(ctx, query, citizenID)
	b, err := scanBiometric(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return b, nil
}

func (r *BiometricRepo) Update(ctx context.Context, b *domain.Biometric) error {
	query := `
		UPDATE biometrics SET
			fp_right_thumb = $2, fp_right_index = $3,
			fp_left_thumb = $4, fp_left_index = $5,
			fp_right_thumb_quality = $6, fp_right_index_quality = $7,
			fp_left_thumb_quality = $8, fp_left_index_quality = $9,
			photo_key = $10, photo_hash = $11, signature_key = $12,
			encryption_key_id = $13, captured_by = $14,
			capture_device = $15, capture_location = $16,
			captured_at = $17
		WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query,
		b.ID,
		b.FPRightThumb, b.FPRightIndex,
		b.FPLeftThumb, b.FPLeftIndex,
		b.FPRightThumbQuality, b.FPRightIndexQuality,
		b.FPLeftThumbQuality, b.FPLeftIndexQuality,
		b.PhotoKey, b.PhotoHash, b.SignatureKey,
		b.EncryptionKeyID, b.CapturedBy,
		b.CaptureDevice, b.CaptureLocation,
		b.CapturedAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func biometricColumns() string {
	return `id, citizen_id,
		fp_right_thumb, fp_right_index, fp_left_thumb, fp_left_index,
		fp_right_thumb_quality, fp_right_index_quality,
		fp_left_thumb_quality, fp_left_index_quality,
		photo_key, photo_hash, signature_key, encryption_key_id,
		captured_by, capture_device, capture_location,
		captured_at, created_at`
}

func scanBiometric(row pgx.Row) (*domain.Biometric, error) {
	var b domain.Biometric
	err := row.Scan(
		&b.ID, &b.CitizenID,
		&b.FPRightThumb, &b.FPRightIndex, &b.FPLeftThumb, &b.FPLeftIndex,
		&b.FPRightThumbQuality, &b.FPRightIndexQuality,
		&b.FPLeftThumbQuality, &b.FPLeftIndexQuality,
		&b.PhotoKey, &b.PhotoHash, &b.SignatureKey, &b.EncryptionKeyID,
		&b.CapturedBy, &b.CaptureDevice, &b.CaptureLocation,
		&b.CapturedAt, &b.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
