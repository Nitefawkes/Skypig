-- Create QSOs table (ADIF-compliant schema)
CREATE TABLE IF NOT EXISTS qsos (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Contact information
    callsign VARCHAR(20) NOT NULL,
    operator_call VARCHAR(20),
    station_callsign VARCHAR(20),

    -- Time information
    qso_date DATE NOT NULL,
    time_on TIMESTAMP WITH TIME ZONE NOT NULL,
    time_off TIMESTAMP WITH TIME ZONE,

    -- Frequency & Band
    band VARCHAR(10),
    band_rx VARCHAR(10),
    freq DECIMAL(10, 6),
    freq_rx DECIMAL(10, 6),

    -- Mode information
    mode VARCHAR(20),
    submode VARCHAR(20),

    -- Signal reports
    rst_sent VARCHAR(10),
    rst_rcvd VARCHAR(10),

    -- Location information
    name VARCHAR(255),
    qth VARCHAR(255),
    gridsquare VARCHAR(10),
    country VARCHAR(100),
    dxcc INTEGER,
    state VARCHAR(50),
    county VARCHAR(100),

    -- Comments & Notes
    comment TEXT,
    notes TEXT,

    -- Power
    tx_pwr DECIMAL(8, 2),
    rx_pwr DECIMAL(8, 2),

    -- Propagation
    prop_mode VARCHAR(20),
    sat_name VARCHAR(50),
    sat_mode VARCHAR(20),

    -- Contest information
    contest_id VARCHAR(100),
    stx INTEGER,
    srx INTEGER,

    -- QSL information
    lotw_qsl_sent CHAR(1) DEFAULT 'N',
    lotw_qsl_rcvd CHAR(1) DEFAULT 'N',
    lotw_qslrdate TIMESTAMP WITH TIME ZONE,
    eqsl_qsl_sent CHAR(1) DEFAULT 'N',
    eqsl_qsl_rcvd CHAR(1) DEFAULT 'N',

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Convert to TimescaleDB hypertable for time-series optimization
SELECT create_hypertable('qsos', 'time_on', if_not_exists => TRUE);

-- Create indexes for common queries
CREATE INDEX idx_qsos_user_id ON qsos(user_id);
CREATE INDEX idx_qsos_callsign ON qsos(callsign);
CREATE INDEX idx_qsos_time_on ON qsos(time_on DESC);
CREATE INDEX idx_qsos_band ON qsos(band);
CREATE INDEX idx_qsos_mode ON qsos(mode);
CREATE INDEX idx_qsos_country ON qsos(country);
CREATE INDEX idx_qsos_gridsquare ON qsos(gridsquare);
CREATE INDEX idx_qsos_user_time ON qsos(user_id, time_on DESC);

-- Composite indexes for common filter combinations
CREATE INDEX idx_qsos_user_band_mode ON qsos(user_id, band, mode);
CREATE INDEX idx_qsos_user_date ON qsos(user_id, qso_date DESC);

-- Create trigger for updated_at
CREATE TRIGGER update_qsos_updated_at
    BEFORE UPDATE ON qsos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create function to increment user QSO count
CREATE OR REPLACE FUNCTION increment_user_qso_count()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users
    SET qso_count = qso_count + 1
    WHERE id = NEW.user_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to auto-increment QSO count
CREATE TRIGGER increment_qso_count_on_insert
    AFTER INSERT ON qsos
    FOR EACH ROW
    EXECUTE FUNCTION increment_user_qso_count();

-- Create function to decrement user QSO count
CREATE OR REPLACE FUNCTION decrement_user_qso_count()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users
    SET qso_count = qso_count - 1
    WHERE id = OLD.user_id AND qso_count > 0;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to auto-decrement QSO count
CREATE TRIGGER decrement_qso_count_on_delete
    AFTER DELETE ON qsos
    FOR EACH ROW
    EXECUTE FUNCTION decrement_user_qso_count();
