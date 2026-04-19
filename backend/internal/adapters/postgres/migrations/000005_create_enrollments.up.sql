CREATE TABLE IF NOT EXISTS enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    citizen_id UUID NOT NULL REFERENCES citizens(id),
    embassy_id UUID NOT NULL REFERENCES embassies(id),
    agent_id UUID NOT NULL REFERENCES users(id),
    team_id VARCHAR(100),
    latitude DOUBLE PRECISION DEFAULT 0,
    longitude DOUBLE PRECISION DEFAULT 0,
    location_name VARCHAR(255),
    offline_data JSONB,
    sync_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    synced_at TIMESTAMPTZ,
    review_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    reviewed_by UUID REFERENCES users(id),
    review_notes TEXT,
    reviewed_at TIMESTAMPTZ,
    enrolled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_enrollments_citizen_id ON enrollments(citizen_id);
CREATE INDEX idx_enrollments_embassy_id ON enrollments(embassy_id);
CREATE INDEX idx_enrollments_agent_id ON enrollments(agent_id);
CREATE INDEX idx_enrollments_sync_status ON enrollments(sync_status);
CREATE INDEX idx_enrollments_review_status ON enrollments(review_status);
