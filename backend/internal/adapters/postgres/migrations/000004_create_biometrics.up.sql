CREATE TABLE IF NOT EXISTS biometrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    citizen_id UUID NOT NULL UNIQUE REFERENCES citizens(id),
    fp_right_thumb BYTEA,
    fp_right_index BYTEA,
    fp_left_thumb BYTEA,
    fp_left_index BYTEA,
    fp_right_thumb_quality INT DEFAULT 0,
    fp_right_index_quality INT DEFAULT 0,
    fp_left_thumb_quality INT DEFAULT 0,
    fp_left_index_quality INT DEFAULT 0,
    photo_key TEXT,
    photo_hash VARCHAR(128),
    signature_key TEXT,
    encryption_key_id VARCHAR(100),
    captured_by UUID NOT NULL REFERENCES users(id),
    capture_device VARCHAR(100),
    capture_location VARCHAR(255),
    captured_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_biometrics_citizen_id ON biometrics(citizen_id);
