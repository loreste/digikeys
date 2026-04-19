package domain

import "time"

type Biometric struct {
	ID                    string    `json:"id" db:"id"`
	CitizenID             string    `json:"citizenId" db:"citizen_id"`
	FPRightThumb          []byte    `json:"-" db:"fp_right_thumb"`
	FPRightIndex          []byte    `json:"-" db:"fp_right_index"`
	FPLeftThumb           []byte    `json:"-" db:"fp_left_thumb"`
	FPLeftIndex           []byte    `json:"-" db:"fp_left_index"`
	FPRightThumbQuality   int       `json:"fpRightThumbQuality" db:"fp_right_thumb_quality"`
	FPRightIndexQuality   int       `json:"fpRightIndexQuality" db:"fp_right_index_quality"`
	FPLeftThumbQuality    int       `json:"fpLeftThumbQuality" db:"fp_left_thumb_quality"`
	FPLeftIndexQuality    int       `json:"fpLeftIndexQuality" db:"fp_left_index_quality"`
	PhotoKey              string    `json:"photoKey,omitempty" db:"photo_key"`
	PhotoHash             string    `json:"photoHash,omitempty" db:"photo_hash"`
	SignatureKey          string    `json:"signatureKey,omitempty" db:"signature_key"`
	EncryptionKeyID       string    `json:"encryptionKeyId,omitempty" db:"encryption_key_id"`
	CapturedBy            string    `json:"capturedBy" db:"captured_by"`
	CaptureDevice         string    `json:"captureDevice,omitempty" db:"capture_device"`
	CaptureLocation       string    `json:"captureLocation,omitempty" db:"capture_location"`
	CapturedAt            time.Time `json:"capturedAt" db:"captured_at"`
	CreatedAt             time.Time `json:"createdAt" db:"created_at"`
}
