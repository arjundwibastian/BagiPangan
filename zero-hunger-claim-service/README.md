# Zero Hunger Claim Service

Owns `food_claims`, pickup code verification, claim lifecycle, and quantity reservation through Food Service.

## Local development

```bash
cp .env.example .env
go mod tidy
go test ./...
go run ./cmd/server
```

The HTTP API listens on `:8084`. The service connects to User Service (`localhost:50051`), Food Service (`localhost:50052`), and Request Service (`localhost:50053`) over gRPC. Configure the external Claim PostgreSQL database in `.env`, then apply `migrations/001_init.sql` before starting the service.

Claim verification is performed by the donor who owns the listing; the recipient owns the claim and may cancel it while it is waiting for pickup.
