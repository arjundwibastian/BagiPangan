CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE food_requests (
    request_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    latitude NUMERIC(9,6) NOT NULL,
    longitude NUMERIC(9,6) NOT NULL,
    radius_km NUMERIC(10,3) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'searching',
    claim_id UUID,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_food_requests_latitude CHECK (latitude BETWEEN -90 AND 90),
    CONSTRAINT chk_food_requests_longitude CHECK (longitude BETWEEN -180 AND 180),
    CONSTRAINT chk_food_requests_radius CHECK (radius_km > 0),
    CONSTRAINT chk_food_requests_status CHECK (status IN ('searching','claimed','completed','cancelled','expired'))
);

CREATE INDEX idx_food_requests_user_created ON food_requests(user_id, created_at DESC);
CREATE INDEX idx_food_requests_searching_expiry ON food_requests(status, expires_at);
