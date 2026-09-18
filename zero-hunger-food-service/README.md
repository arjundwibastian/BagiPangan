# Zero Hunger Food Service

Food Service will own food listings, locations, PostGIS queries, quantity, pickup windows, and listing status.

This repository is intentionally separate from User Service. It will consume the User Service contract from `github.com/zero-hunger/contracts` and validate providers through gRPC.

## Local development

```bash
cp .env.example .env
go test ./...
go run ./cmd
```

Apply `migrations/001_init.sql` to the external Food PostgreSQL database first. REST listens on `:8082`, gRPC on `:50052`, and local gRPC dependencies are User `localhost:50051` and Request `localhost:50053`.
