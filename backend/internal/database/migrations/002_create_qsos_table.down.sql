-- Drop triggers
DROP TRIGGER IF EXISTS decrement_qso_count_on_delete ON qsos;
DROP TRIGGER IF EXISTS increment_qso_count_on_insert ON qsos;
DROP TRIGGER IF EXISTS update_qsos_updated_at ON qsos;

-- Drop functions
DROP FUNCTION IF EXISTS decrement_user_qso_count();
DROP FUNCTION IF EXISTS increment_user_qso_count();

-- Drop indexes
DROP INDEX IF EXISTS idx_qsos_user_date;
DROP INDEX IF EXISTS idx_qsos_user_band_mode;
DROP INDEX IF EXISTS idx_qsos_user_time;
DROP INDEX IF EXISTS idx_qsos_gridsquare;
DROP INDEX IF EXISTS idx_qsos_country;
DROP INDEX IF EXISTS idx_qsos_mode;
DROP INDEX IF EXISTS idx_qsos_band;
DROP INDEX IF EXISTS idx_qsos_time_on;
DROP INDEX IF EXISTS idx_qsos_callsign;
DROP INDEX IF EXISTS idx_qsos_user_id;

-- Drop table
DROP TABLE IF EXISTS qsos;
