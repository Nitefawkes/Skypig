-- Drop trigger
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_tier;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_callsign;

-- Drop table
DROP TABLE IF EXISTS users;

-- Drop function (only if no other tables use it)
-- DROP FUNCTION IF EXISTS update_updated_at_column();
