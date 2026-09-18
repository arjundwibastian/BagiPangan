# Zero Hunger Infra

Local orchestration and end-to-end test assets for the Zero Hunger services.

## Local source workflow (recommended)

This workflow runs all application services from source and keeps PostgreSQL on the developer-managed host. Create databases `user_db`, `food_db`, `request_db`, and `claim_db`, then apply each service migration:

```bash
psql "$USER_DATABASE_URL" -f ../zero-hunger-user-service/migrations/001_init.sql
psql "$REQUEST_DATABASE_URL" -f ../zero-hunger-request-service/migrations/001_init.sql
psql "postgres://user:password@localhost:5434/food_db?sslmode=disable" -f ../zero-hunger-food-service/migrations/001_init.sql
psql "$CLAIM_DATABASE_URL" -f ../zero-hunger-claim-service/migrations/001_init.sql
```

Copy `.env.example` to `.env` in each service directory and replace the database credentials/hosts with the developer's external database values. The local service ports are:

| Service | HTTP | gRPC |
| --- | ---: | ---: |
| User | 8081 | 50051 |
| Food | 8082 | 50052 |
| Request | 8083 | 50053 |
| Claim | 8084 | — |

Start in this order, each command in its own terminal:

```bash
cd ../zero-hunger-user-service && go run ./cmd/server
cd ../zero-hunger-request-service && go run ./cmd/server
cd ../zero-hunger-food-service && go run ./cmd
cd ../zero-hunger-claim-service && go run ./cmd
```

Verify that each process logs its HTTP listener and that gRPC clients connect without repeated errors. Then import `docs/Zero Hunger Local.postman_environment.json`, select it, and run `docs/Zero Hunger E2E.postman_collection.json` in order.

### VS Code

Open this infra directory in VS Code and use the committed multi-target launch configuration in `.vscode/launch.json`. It starts all four Go services from source; `.env` files are loaded from each service directory. The external database remains outside VS Code and must already be reachable.

## Docker Compose alternative

Compose runs published service images only. It does not create PostgreSQL containers or run migrations. Configure the external host databases in `.env` (Compose uses `host.docker.internal` by default for database connections), then run:

```bash
cp .env.example .env
docker compose pull
docker compose up
```

Container-to-container gRPC addresses are `user-service:50051`, `request-service:50053`, and `food-service:50052`; host clients still use the localhost ports in the table above.

## Troubleshooting

- `connection refused`: check the database host/port, migration, and service startup order.
- gRPC `Unavailable`: confirm the target service is running on the expected port and that the `.env` address uses `localhost` for source mode.
- Claim migration errors: ensure the Claim database is separate and run `migrations/001_init.sql` once.
- Duplicate registration: change the email variables in the Postman environment or clean only the test rows in the external databases.

All REST responses use `{ "rc": "00", "message": "...", "data": ... }` for User/Request and `{ "responseCode": "00", "responseMessage": "...", "responseData": ... }` for Food/Claim.

See [END_TO_END_FLOW.md](docs/END_TO_END_FLOW.md) for ownership, lifecycle, and cancellation details.
