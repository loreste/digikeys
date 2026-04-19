CREATE TABLE IF NOT EXISTS cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    citizen_id UUID NOT NULL REFERENCES citizens(id),
    card_number VARCHAR(30) NOT NULL UNIQUE,
    mrz_line1 VARCHAR(30) NOT NULL,
    mrz_line2 VARCHAR(30) NOT NULL,
    mrz_line3 VARCHAR(30) NOT NULL,
    embassy_id UUID NOT NULL REFERENCES embassies(id),
    issued_by UUID REFERENCES users(id),
    issued_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    print_batch_id VARCHAR(100),
    printed_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    previous_card_id UUID REFERENCES cards(id),
    renewal_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cards_citizen_id ON cards(citizen_id);
CREATE INDEX idx_cards_card_number ON cards(card_number);
CREATE INDEX idx_cards_embassy_id ON cards(embassy_id);
CREATE INDEX idx_cards_status ON cards(status);
CREATE INDEX idx_cards_print_batch_id ON cards(print_batch_id);

-- Sequence tracking for card numbers per embassy per year
CREATE TABLE IF NOT EXISTS card_sequences (
    embassy_id UUID NOT NULL REFERENCES embassies(id),
    year INT NOT NULL,
    last_sequence INT NOT NULL DEFAULT 0,
    PRIMARY KEY (embassy_id, year)
);
