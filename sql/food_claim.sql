
CREATE TYPE "claim_type" AS ENUM (
  'waiting_for_pickup',
  'picked_up',
  'cancelled'
);

CREATE TABLE "food_claims" (
  "claim_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "request_id" UUID NOT NULL,
  "food_listing_id" UUID NOT NULL,
  "claim_code" VARCHAR(10) UNIQUE NOT NULL,
  "code_expires_at" TIMESTAMPTZ NOT NULL,
  "claimed_at" TIMESTAMPTZ,
  "claimed_quantity" INT NOT NULL CHECK (claimed_quantity > 0),
  "status" claim_type NOT NULL DEFAULT 'waiting_for_pickup',
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_food_claims_request_id"
ON "food_claims" ("request_id");

CREATE INDEX "idx_food_claims_food_listing_id"
ON "food_claims" ("food_listing_id");

CREATE INDEX "idx_food_claims_status"
ON "food_claims" ("status");

CREATE INDEX "idx_food_claims_claim_code"
ON "food_claims" ("claim_code");

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER food_claims_updated_at
BEFORE UPDATE ON "food_claims"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();


select * from food_claims;
