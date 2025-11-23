-- Seed test data for development

-- Insert test user (callsign: W1AW)
INSERT INTO users (callsign, email, name, grid_square, qrz_verified, tier, qso_limit, qso_count)
VALUES ('W1AW', 'w1aw@example.com', 'Test User', 'FN31pr', true, 'operator', 20000, 0)
ON CONFLICT (callsign) DO NOTHING;

-- Note: This migration should only be applied in development environments
-- In production, use proper user registration flow
