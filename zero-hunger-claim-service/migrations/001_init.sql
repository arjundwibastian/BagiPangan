CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'claim_type') THEN
        CREATE TYPE claim_type AS ENUM ('waiting_for_pickup', 'picked_up', 'cancelled');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS food_claims (
    claim_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID NOT NULL,
    food_listing_id UUID NOT NULL,
    claim_code VARCHAR(10) NOT NULL UNIQUE,
    code_expires_at TIMESTAMPTZ NOT NULL,
    claimed_at TIMESTAMPTZ,
    claimed_quantity INTEGER NOT NULL CHECK (claimed_quantity > 0),
    status claim_type NOT NULL DEFAULT 'waiting_for_pickup',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_food_claims_request_id ON food_claims(request_id);
CREATE INDEX IF NOT EXISTS idx_food_claims_food_listing_id ON food_claims(food_listing_id);
