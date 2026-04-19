CREATE TABLE IF NOT EXISTS transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    citizen_id UUID NOT NULL REFERENCES citizens(id),
    bank_account_id UUID REFERENCES bank_accounts(id),
    amount BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'XOF',
    type VARCHAR(30) NOT NULL,
    source_provider VARCHAR(50),
    source_reference VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    external_ref VARCHAR(100),
    failure_reason TEXT,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transfers_citizen_id ON transfers(citizen_id);
CREATE INDEX idx_transfers_type ON transfers(type);
CREATE INDEX idx_transfers_status ON transfers(status);
CREATE INDEX idx_transfers_external_ref ON transfers(external_ref);
