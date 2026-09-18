# Zero Hunger REST API

All endpoints use the envelope below:

```json
{
  "rc": "00",
  "message": "success",
  "data": {}
}
```

Error codes: `40` invalid input, `41` unauthorized, `42` not found, `43` conflict, `44` forbidden, `50` internal error.

## User Service (`:8081`)

| Method | Path | Auth | Purpose |
|---|---|---|---|
| POST | `/api/v1/auth/register` | No | Register donor or recipient |
| POST | `/api/v1/auth/login` | No | Login and issue JWT |
| POST | `/api/v1/auth/refresh` | No | Issue a new access token |
| POST | `/api/v1/auth/logout` | No | Revoke refresh token |
| GET | `/api/v1/users/me` | Bearer | Get current profile |
| PATCH | `/api/v1/users/me` | Bearer | Update name and phone |
| GET | `/api/v1/users/{user_id}` | Internal/Bearer | Get user summary |

Register request:

```json
{
  "name": "Budi",
  "email": "budi@example.com",
  "password": "password123",
  "phone": "08123456789",
  "role": "donor"
}
```

## Food Service (`:8082`)

| Method | Path | Auth | Purpose |
|---|---|---|---|
| POST | `/api/v1/food-listings` | Donor | Create food listing |
| GET | `/api/v1/food-listings` | Bearer | List active listings |
| GET | `/api/v1/food-listings/{food_listing_id}` | Bearer | Get listing detail |
| PATCH | `/api/v1/food-listings/{food_listing_id}` | Donor owner | Update listing |
| DELETE | `/api/v1/food-listings/{food_listing_id}` | Donor owner | Cancel listing |
| GET | `/api/v1/food-listings/nearby` | Recipient | Search by PostGIS radius |

Create request:

```json
{
  "title": "Nasi Kotak",
  "description": "Nasi ayam dan sayur",
  "latitude": -6.2,
  "longitude": 106.816666,
  "food_type": "cooked_food",
  "quantity": 20,
  "available_from": "2026-09-07T10:00:00Z",
  "available_until": "2026-09-07T14:00:00Z"
}
```

Nearby query:

```text
GET /api/v1/food-listings/nearby?latitude=-6.2&longitude=106.816666&radius_km=5
```

## Request Service (`:8083`)

| Method | Path | Auth | Purpose |
|---|---|---|---|
| POST | `/api/v1/food-requests` | Recipient | Create a nearby food request |
| GET | `/api/v1/food-requests/{request_id}` | Owner/Admin | Get request detail |
| GET | `/api/v1/users/{user_id}/food-requests` | Owner/Admin | List user requests |
| PATCH | `/api/v1/food-requests/{request_id}/cancel` | Owner | Cancel request |

Create request:

```json
{
  "latitude": -6.2,
  "longitude": 106.816666,
  "radius_km": 5
}
```

## Claim Service (`:8084`)

| Method | Path | Auth | Purpose |
|---|---|---|---|
| POST | `/api/v1/claims` | Recipient | Claim available food |
| GET | `/api/v1/claims/{claim_id}` | Owner/Donor | Get claim detail |
| POST | `/api/v1/claims/{claim_id}/verify-pickup` | Donor | Verify pickup code |
| POST | `/api/v1/claims/{claim_id}/cancel` | Recipient/Donor | Cancel claim |

Create request:

```json
{
  "request_id": "uuid",
  "food_listing_id": "uuid",
  "claimed_quantity": 2
}
```

Pickup request:

```json
{
  "claim_code": "123456"
}
```

## Review Service (`:8085`)

| Method | Path | Auth | Purpose |
|---|---|---|---|
| POST | `/api/v1/ratings` | Recipient | Rate completed claim |
| GET | `/api/v1/claims/{claim_id}/rating` | Bearer | Get claim rating |
| GET | `/api/v1/users/{user_id}/rating-summary` | Bearer | Get rating summary |
| POST | `/api/v1/food-reports` | Recipient | Report a claim |
| GET | `/api/v1/food-reports/{report_id}` | Owner/Admin | Get report |
| PATCH | `/api/v1/food-reports/{report_id}/status` | Admin | Update report status |

Rating request:

```json
{
  "claim_id": "uuid",
  "rating": 5,
  "comment": "Makanan sesuai dan masih baik."
}
```

Report request:

```json
{
  "claim_id": "uuid",
  "reason": "food_quality",
  "description": "Makanan tidak sesuai kondisi."
}
```
