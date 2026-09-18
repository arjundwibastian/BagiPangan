# Zero Hunger Request Service

Owns `food_requests`, recipient request locations, and request lifecycle. It calls User Service for requester validation. Nearby matching is owned by Food Service and uses the request ID.

## Local development

```bash
go test ./...
go run ./cmd/server
```

Defaults are HTTP `:8083`, gRPC `:50053`, request PostgreSQL at `localhost:5435`, and a 24-hour request TTL. The service uses the same JWT secret as User Service.

`.env` is loaded automatically when running locally from the service directory. It is optional because the service also has defaults. To use explicit values:

```bash
cp .env.example .env
go run ./cmd/server
```

Existing shell environment variables take precedence over values in `.env`. The real `.env` is ignored by Git; only `.env.example` is committed.

Start the dependencies first with the User Service, Food Service, and PostgreSQL databases. The Postman files in `docs/` provide the local test flow:

1. Import `Zero Hunger Request Service.postman_environment.json` and select the environment.
2. Run `User Service - Setup recipient > Register Recipient`.
3. Run `User Service - Setup recipient > Login Recipient`.
4. Run `Request Service > Create Food Request`.
5. Run `Food Service - Nearby matching > Search Nearby Food For Request`.

The Postman create request calls Request Service over HTTP. Request Service calls User Service gRPC `GetUser` to validate that the JWT subject is a recipient. After creation, call Food Service `/api/v1/food-listings/nearby?request_id=...` to search nearby matches; Food Service uses Request Service gRPC `GetRequest` internally.

### VS Code debug

Open this repository as the VS Code workspace and use `Run and Debug > Debug Request Service`. If `.env` exists, the application loads it automatically; otherwise the built-in defaults are used. PostgreSQL, User Service, and Food Service must be running first. The task `Start Request Service dependencies` starts them through the infra Compose file.

Important log events include:

- `event=user_service_get_user_start` and `event=user_service_recipient_validated`
- `event=food_service_search_nearby_start` and `event=food_service_search_nearby_success`
- `event=request_created`, `event=request_claimed`, and `event=request_completed`

HTTP endpoints:

- `POST /api/v1/food-requests` — recipient creates a request.
- `GET /api/v1/food-requests/{request_id}` — owner/admin gets a request.
- `GET /api/v1/users/{user_id}/food-requests` — owner/admin lists requests.
- `PATCH /api/v1/food-requests/{request_id}/cancel` — owner cancels a searching request.
