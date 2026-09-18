# Zero Hunger User Service

Owns authentication, user profile, and exclusive user roles: `admin`, `donor`, and `recipient`.

## Local development

```bash
go mod tidy
go test ./...
go run ./cmd/server
```

REST listens on `SERVER_PORT` and internal gRPC listens on `GRPC_PORT`.

Public registration accepts only `donor` and `recipient`. Admin users are provisioned separately.
