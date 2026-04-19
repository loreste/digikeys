-- Multi-country support: Burkina Faso (BF) and DRC (CD)
-- The platform can issue consular cards for both countries' diasporas

ALTER TABLE citizens ADD COLUMN IF NOT EXISTS origin_country VARCHAR(5) DEFAULT 'BF';
ALTER TABLE cards ADD COLUMN IF NOT EXISTS origin_country VARCHAR(5) DEFAULT 'BF';
ALTER TABLE embassies ADD COLUMN IF NOT EXISTS origin_country VARCHAR(5) DEFAULT 'BF';

-- Insert DRC as origin country
INSERT INTO countries (code, name, name_en, region, currency, timezone, active) VALUES
    ('BF', 'Burkina Faso', 'Burkina Faso', 'Afrique de l''Ouest', 'XOF', 'Africa/Ouagadougou', TRUE),
    ('CD', 'République Démocratique du Congo', 'Democratic Republic of the Congo', 'Afrique Centrale', 'CDF', 'Africa/Kinshasa', TRUE)
ON CONFLICT (code) DO NOTHING;

-- DRC-specific card prefix patterns
-- BF cards: CC-{country}-{year}-{seq} (e.g., CC-FR-2026-000001)
-- CD cards: CD-{country}-{year}-{seq} (e.g., CD-FR-2026-000001)

-- Sample DRC embassies/consulates
INSERT INTO embassies (id, country_code, name, city, origin_country, card_prefix, status) VALUES
    (gen_random_uuid(), 'FR', 'Ambassade de la RDC en France', 'Paris', 'CD', 'CD-FR', 'active'),
    (gen_random_uuid(), 'BE', 'Ambassade de la RDC en Belgique', 'Bruxelles', 'CD', 'CD-BE', 'active'),
    (gen_random_uuid(), 'US', 'Ambassade de la RDC aux États-Unis', 'Washington', 'CD', 'CD-US', 'active'),
    (gen_random_uuid(), 'ZA', 'Ambassade de la RDC en Afrique du Sud', 'Pretoria', 'CD', 'CD-ZA', 'active'),
    (gen_random_uuid(), 'GB', 'Ambassade de la RDC au Royaume-Uni', 'Londres', 'CD', 'CD-GB', 'active'),
    (gen_random_uuid(), 'CA', 'Ambassade de la RDC au Canada', 'Ottawa', 'CD', 'CD-CA', 'active'),
    (gen_random_uuid(), 'CN', 'Ambassade de la RDC en Chine', 'Pékin', 'CD', 'CD-CN', 'active'),
    (gen_random_uuid(), 'AO', 'Ambassade de la RDC en Angola', 'Luanda', 'CD', 'CD-AO', 'active'),
    (gen_random_uuid(), 'CG', 'Ambassade de la RDC au Congo-Brazzaville', 'Brazzaville', 'CD', 'CD-CG', 'active'),
    (gen_random_uuid(), 'TZ', 'Ambassade de la RDC en Tanzanie', 'Dar es Salaam', 'CD', 'CD-TZ', 'active'),
    (gen_random_uuid(), 'KE', 'Ambassade de la RDC au Kenya', 'Nairobi', 'CD', 'CD-KE', 'active'),
    (gen_random_uuid(), 'IN', 'Ambassade de la RDC en Inde', 'New Delhi', 'CD', 'CD-IN', 'active')
ON CONFLICT DO NOTHING;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_citizens_origin ON citizens(origin_country);
CREATE INDEX IF NOT EXISTS idx_cards_origin ON cards(origin_country);
CREATE INDEX IF NOT EXISTS idx_embassies_origin ON embassies(origin_country);
