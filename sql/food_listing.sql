CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE "food_category" AS ENUM (
  'cooked_food',
  'raw_food',
  'bakery'
);

CREATE TYPE "listing_status" AS ENUM (
  'available',
  'claimed',
  'expired'
);

CREATE TABLE "food_listings" (
  "food_listing_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID NOT NULL,

  "title" VARCHAR(255) NOT NULL,
  "description" TEXT,

  "latitude" DECIMAL(9,6) NOT NULL,
  "longitude" DECIMAL(9,6) NOT NULL,

  "food_type" food_category NOT NULL DEFAULT 'cooked_food',
  "quantity" INT NOT NULL CHECK (quantity >= 0),

  "available_from" TIMESTAMPTZ NOT NULL,
  "available_until" TIMESTAMPTZ NOT NULL,

  "status" listing_status NOT NULL DEFAULT 'available',
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT chk_food_listings_latitude
    CHECK (latitude BETWEEN -90 AND 90),

  CONSTRAINT chk_food_listings_longitude
    CHECK (longitude BETWEEN -180 AND 180),

  CONSTRAINT chk_food_listings_availability
    CHECK (available_until > available_from)
);

CREATE INDEX "idx_food_listings_user_id"
ON "food_listings" ("user_id");

CREATE INDEX "idx_food_listings_status"
ON "food_listings" ("status");

CREATE INDEX "idx_food_listings_latitude_longitude"
ON "food_listings" ("latitude", "longitude");

CREATE OR REPLACE FUNCTION update_food_listings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER food_listings_updated_at
BEFORE UPDATE ON "food_listings"
FOR EACH ROW
EXECUTE FUNCTION update_food_listings_updated_at();
