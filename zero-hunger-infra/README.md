# Zero Hunger Infra

Local orchestration and end-to-end test assets for the Zero Hunger services.

## Running the stack

Compose runs the application services against a PostgreSQL host you manage; it does not create PostgreSQL containers or run migrations. Create databases `user_db`, `food_db`, `request_db`, and `claim_db`, then apply each service migration:

```bash
psql "$USER_DATABASE_URL" -f ../zero-hunger-user-service/migrations/001_init.sql
psql "$REQUEST_DATABASE_URL" -f ../zero-hunger-request-service/migrations/001_init.sql
psql "postgres://user:password@localhost:5434/food_db?sslmode=disable" -f ../zero-hunger-food-service/migrations/001_init.sql
psql "$CLAIM_DATABASE_URL" -f ../zero-hunger-claim-service/migrations/001_init.sql
```

Configure the external host databases in `.env`, then start:

```bash
cp .env.example .env
docker compose up --build
```

The service ports are:

| Service | HTTP | gRPC |
| --- | ---: | ---: |
| User | 8081 | 50051 |
| Food | 8082 | 50052 |
| Request | 8083 | 50053 |
| Claim | 8084 | — |

Container-to-container gRPC addresses are `user-service:50051`, `request-service:50053`, and `food-service:50052`; host clients use the localhost ports above.

Verify that each service logs its listener, then import `docs/Zero Hunger Local.postman_environment.json`, select it, and run `docs/Zero Hunger E2E.postman_collection.json` in order.

The services read configuration only from the process environment and do not load `.env` files themselves, so start them through Compose (or another tool that injects the environment).

## HTTPS edge + CORS (Caddy)

A `caddy` service terminates TLS and routes to the four services on one origin, and applies a single CORS policy for all of them. The services themselves need no CORS code.

- Local: `https://localhost:8443` (self-signed via Caddy's internal CA — accept the browser warning once).
- The app services still publish `8081`-`8084` for direct/Postman access; in production drop those `ports:` entries so only the proxy is public.
- Configure in `.env`:
  - `SITE_ADDRESS` — the public hostname (default `localhost`).
  - `CORS_ALLOWED_ORIGIN` — the exact origin of your browser frontend (default `http://localhost:3000`).

Routing is path-based; specific paths win. See `Caddyfile`. Swagger `/swagger/*` is not proxied.

For production, set `SITE_ADDRESS` to your domain, remove `tls internal` from the `Caddyfile`, and publish `80`/`443` instead of `8443`.

## Troubleshooting

- `connection refused`: check the database host/port, migration, and service startup order.
- gRPC `Unavailable`: confirm the target service is running on the expected port and that the gRPC address is set correctly (container names inside Compose, `localhost` for host clients).
- Claim migration errors: ensure the Claim database is separate and run `migrations/001_init.sql` once.
- Duplicate registration: change the email variables in the Postman environment or clean only the test rows in the external databases.

All REST responses use the envelope `{ "responseCode": "00", "responseMessage": "...", "responseData": ... }`.

See [END_TO_END_FLOW.md](docs/END_TO_END_FLOW.md) for ownership, lifecycle, and cancellation details.
