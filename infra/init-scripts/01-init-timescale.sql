-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Create initial schema
CREATE SCHEMA IF NOT EXISTS hamradio;

-- Set search path
SET search_path TO hamradio, public;

-- Log successful initialization
DO $$
BEGIN
  RAISE NOTICE 'TimescaleDB extension enabled successfully';
END $$;
