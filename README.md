# BagiPangan — share extra food with people nearby

BagiPangan (formerly ZeroHunger) connects people who have extra food with people who need it. Donors post what they have, recipients ask for what they need, and the system matches them through claims with a simple pickup code.

What you get:

- Accounts with roles. Register as a donor or recipient, log in, and get a JWT. Donor-only routes are enforced, so only donors can post listings.
- Food listings with quantity and location. Donors create listings with stock counts and pickup windows. Quantity is reserved on claim and released on cancel.
- Nearby search. Recipients create a request with a location and radius, then search for available listings around that request.
- Requests with a clear lifecycle. Every request moves `searching -> claimed -> completed`, with a clean path back to `searching` if a claim is cancelled.
- Claims with pickup codes. Claiming a listing generates a 6-digit code. The claim sits at `waiting_for_pickup` until the donor verifies the code, then it moves to `picked_up`.
- Email notifications on claim status via Resend.
- Microservices in Go with shared protobuf contracts, runnable together with Docker Compose, plus a Postman E2E collection to try the full flow.

## How the system works

![project flow](./project-flow.webp)

There are two sides running in parallel — food on top, requests at the bottom — and a claim joins them in the middle.

It starts with users. Someone registers as a donor, someone else as a recipient. A donor posts a food listing, so it becomes `available`. On the other side, a recipient posts a request, which starts out as `searching`.

When the recipient finds a suitable listing nearby, they create a claim. That's the link: the claim goes to `waiting_for_pickup` and the request moves from `searching` to `claimed`. The quantity on the listing is reserved at that point.

Pickup is confirmed with a code by the donor. Once verified, the claim becomes `picked_up` and the request becomes `completed`. That's the happy path.

If something falls through, there are two cancel paths shown as dashed lines in the diagram. A claim can be cancelled from `waiting_for_pickup` to `cancelled`, which triggers `release quantity` back to the food listing so it becomes available again. On the request side there's a `cancel flow` from `claimed` back to `searching`, so the recipient can try again.

## Architecture

```
zerohunger2/
  zero-hunger-contracts/       protobuf + generated Go code
  zero-hunger-user-service/    HTTP 8081 / gRPC 50051
  zero-hunger-food-service/    HTTP 8082 / gRPC 50052
  zero-hunger-request-service/ HTTP 8083 / gRPC 50053
  zero-hunger-claim-service/   HTTP 8084 / gRPC 50054
  zero-hunger-infra/           docker-compose.yml, .env.example, docs + Postman
```

Services talk to each other over gRPC with shared protobuf contracts. Details are in `zero-hunger-infra/docs/END_TO_END_FLOW.md`.

## How to run it

1. Install Go 1.22+, Docker, psql.

2. Create 4 databases:
`user_db`, `food_db`, `request_db`, `claim_db`

3. Apply migrations, one per DB:

```bash
psql "$USER_DATABASE_URL" -f zero-hunger-user-service/migrations/001_init.sql
psql "$REQUEST_DATABASE_URL" -f zero-hunger-request-service/migrations/001_init.sql
psql "postgres://user:password@localhost:5434/food_db?sslmode=disable" -f zero-hunger-food-service/migrations/001_init.sql
psql "$CLAIM_DATABASE_URL" -f zero-hunger-claim-service/migrations/001_init.sql
```

4. Configure env:

```bash
cp zero-hunger-infra/.env.example zero-hunger-infra/.env
```

Compose expects `*_HOST`, `*_USER`, `*_PASSWORD`, `JWT_SECRET`, `RESEND_API_KEY`. Defaults assume `host.docker.internal` and `sslmode=disable`.

5. Run:

```bash
cd zero-hunger-infra
docker compose up --build
```

Services will be on `http://localhost:8081` (user), `:8082` (food), `:8083` (request), `:8084` (claim). Import `zero-hunger-infra/docs/Zero Hunger Local.postman_environment.json` and run `Zero Hunger E2E.postman_collection.json` in order.

Or run from source, one terminal each:

```bash
cd ../zero-hunger-user-service && go run ./cmd/server
cd ../zero-hunger-request-service && go run ./cmd/server
cd ../zero-hunger-food-service && go run ./cmd
cd ../zero-hunger-claim-service && go run ./cmd
```

Local gRPC addresses default to `localhost:50051/50052/50053`, in Compose they become `user-service:50051` etc.
