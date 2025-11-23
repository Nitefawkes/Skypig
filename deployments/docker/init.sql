-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    callsign VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    qrz_user_id VARCHAR(100),
    tier VARCHAR(50) DEFAULT 'free' CHECK (tier IN ('free', 'operator', 'contester', 'partner')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User settings table
CREATE TABLE IF NOT EXISTS user_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    lotw_enabled BOOLEAN DEFAULT false,
    lotw_username VARCHAR(100),
    lotw_password_encrypted TEXT,
    propagation_alerts BOOLEAN DEFAULT true,
    preferred_bands TEXT[] DEFAULT ARRAY['20m', '40m'],
    grid_square VARCHAR(10),
    time_zone VARCHAR(100) DEFAULT 'UTC',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- QSOs table (main logbook)
CREATE TABLE IF NOT EXISTS qsos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    callsign VARCHAR(20) NOT NULL,
    frequency DECIMAL(10, 6),
    band VARCHAR(10),
    mode VARCHAR(20),
    rst_sent VARCHAR(10),
    rst_received VARCHAR(10),
    qso_date DATE NOT NULL,
    time_on TIMESTAMP WITH TIME ZONE NOT NULL,
    time_off TIMESTAMP WITH TIME ZONE,
    grid_square VARCHAR(10),
    country VARCHAR(100),
    state VARCHAR(50),
    county VARCHAR(100),
    comment TEXT,
    contest_id VARCHAR(100),
    propagation_mode VARCHAR(50),
    satellite_name VARCHAR(50),
    tx_power DECIMAL(10, 2),
    lotw_sent BOOLEAN DEFAULT false,
    lotw_confirmed BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Convert QSOs to hypertable for time-series optimization
SELECT create_hypertable('qsos', 'time_on', if_not_exists => TRUE);

-- Propagation data table
CREATE TABLE IF NOT EXISTS propagation_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    solar_flux DECIMAL(10, 2),
    sunspot_number INT,
    a_index INT,
    k_index INT,
    xray_flux VARCHAR(10),
    helium_line DECIMAL(10, 2),
    proton_flux INT,
    electron_flux INT,
    source VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Convert propagation data to hypertable
SELECT create_hypertable('propagation_data', 'timestamp', if_not_exists => TRUE);

-- SDR Receivers table (WebSDR directory)
CREATE TABLE IF NOT EXISTS sdr_receivers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    callsign VARCHAR(20),
    url TEXT NOT NULL UNIQUE,
    type VARCHAR(50) DEFAULT 'kiwisdr' CHECK (type IN ('kiwisdr', 'websdr', 'openwebrx', 'other')),
    location VARCHAR(255),
    grid_square VARCHAR(10),
    latitude DECIMAL(10, 6),
    longitude DECIMAL(10, 6),
    country VARCHAR(100),
    bands TEXT[] DEFAULT ARRAY[]::TEXT[],
    modes TEXT[] DEFAULT ARRAY[]::TEXT[],
    antenna_info TEXT,
    frequency_min DECIMAL(10, 6),
    frequency_max DECIMAL(10, 6),
    users_max INT,
    status VARCHAR(50) DEFAULT 'online' CHECK (status IN ('online', 'offline', 'unknown')),
    last_seen TIMESTAMP WITH TIME ZONE,
    description TEXT,
    avatar_url TEXT,
    is_public BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User favorite SDRs (many-to-many)
CREATE TABLE IF NOT EXISTS user_favorite_sdrs (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    sdr_id UUID REFERENCES sdr_receivers(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, sdr_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_qsos_user_id ON qsos(user_id);
CREATE INDEX IF NOT EXISTS idx_qsos_callsign ON qsos(callsign);
CREATE INDEX IF NOT EXISTS idx_qsos_band ON qsos(band);
CREATE INDEX IF NOT EXISTS idx_qsos_mode ON qsos(mode);
CREATE INDEX IF NOT EXISTS idx_qsos_qso_date ON qsos(qso_date);
CREATE INDEX IF NOT EXISTS idx_users_callsign ON users(callsign);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_sdr_receivers_type ON sdr_receivers(type);
CREATE INDEX IF NOT EXISTS idx_sdr_receivers_status ON sdr_receivers(status);
CREATE INDEX IF NOT EXISTS idx_sdr_receivers_country ON sdr_receivers(country);
CREATE INDEX IF NOT EXISTS idx_sdr_receivers_location ON sdr_receivers(location);
CREATE INDEX IF NOT EXISTS idx_sdr_receivers_grid_square ON sdr_receivers(grid_square);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_qsos_updated_at BEFORE UPDATE ON qsos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_settings_updated_at BEFORE UPDATE ON user_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sdr_receivers_updated_at BEFORE UPDATE ON sdr_receivers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Seed a test user (development only)
INSERT INTO users (callsign, email, name, tier)
VALUES ('W1AW', 'test@hamradio.cloud', 'Test User', 'operator')
ON CONFLICT (callsign) DO NOTHING;
