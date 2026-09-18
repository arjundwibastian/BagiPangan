CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE "report_status" AS ENUM (
  'waiting_for_review',
  'reviewed'
);

CREATE TABLE "ratings" (
  "rating_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID NOT NULL,
  "claim_id" UUID NOT NULL,
  "rating" INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
  "comment" TEXT,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_ratings_user_claim
    UNIQUE ("user_id", "claim_id")
);

CREATE INDEX "idx_ratings_user_id"
ON "ratings" ("user_id");

CREATE INDEX "idx_ratings_claim_id"
ON "ratings" ("claim_id");


CREATE TABLE "food_reports" (
  "report_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID NOT NULL,
  "claim_id" UUID NOT NULL,
  "reason" VARCHAR(255) NOT NULL,
  "description" TEXT,
  "status" report_status NOT NULL DEFAULT 'waiting_for_review',
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_food_reports_user_id"
ON "food_reports" ("user_id");

CREATE INDEX "idx_food_reports_claim_id"
ON "food_reports" ("claim_id");

CREATE INDEX "idx_food_reports_status"
ON "food_reports" ("status");
