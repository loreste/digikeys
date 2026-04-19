CREATE TABLE IF NOT EXISTS citizens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    maiden_name VARCHAR(100),
    date_of_birth DATE NOT NULL,
    place_of_birth VARCHAR(255) NOT NULL,
    gender VARCHAR(10) NOT NULL,
    nationality VARCHAR(50) NOT NULL DEFAULT 'Burkinabè',
    national_id VARCHAR(50),
    unique_identifier VARCHAR(50),
    passport_number VARCHAR(50),
    phone VARCHAR(50),
    email VARCHAR(255),
    country_of_residence VARCHAR(100) NOT NULL,
    city_of_residence VARCHAR(100) NOT NULL,
    address_abroad TEXT,
    province_of_origin VARCHAR(100),
    commune_of_origin VARCHAR(100),
    embassy_id UUID NOT NULL REFERENCES embassies(id),
    registered_by UUID NOT NULL REFERENCES users(id),
    photo_key TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_citizens_embassy_id ON citizens(embassy_id);
CREATE INDEX idx_citizens_national_id ON citizens(national_id);
CREATE INDEX idx_citizens_unique_identifier ON citizens(unique_identifier);
CREATE INDEX idx_citizens_country_of_residence ON citizens(country_of_residence);
CREATE INDEX idx_citizens_name ON citizens(last_name, first_name);
CREATE INDEX idx_citizens_status ON citizens(status);
