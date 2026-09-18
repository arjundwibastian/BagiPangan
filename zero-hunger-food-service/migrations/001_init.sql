CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'food_category') THEN
        CREATE TYPE food_category AS ENUM ('cooked_food', 'raw_food', 'bakery');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'listing_status') THEN
        CREATE TYPE listing_status AS ENUM ('available', 'claimed', 'expired');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS food_listings (
    food_listing_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    latitude NUMERIC(9,6) NOT NULL CHECK (latitude BETWEEN -90 AND 90),
    longitude NUMERIC(9,6) NOT NULL CHECK (longitude BETWEEN -180 AND 180),
    food_type food_category NOT NULL DEFAULT 'cooked_food',
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    available_from TIMESTAMPTZ NOT NULL,
    available_until TIMESTAMPTZ NOT NULL,
    status listing_status NOT NULL DEFAULT 'available',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_food_listing_window CHECK (available_until > available_from)
);

CREATE INDEX IF NOT EXISTS idx_food_listings_active_window
    ON food_listings(status, available_from, available_until);
