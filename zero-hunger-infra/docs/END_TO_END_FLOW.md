# Zero Hunger end-to-end flow

```mermaid
sequenceDiagram
    participant R as Recipient
    participant U as User Service
    participant Q as Request Service
    participant F as Food Service
    participant C as Claim Service
    participant D as Donor

    R->>U: Register/Login
    U-->>R: JWT recipient
    D->>U: Register/Login
    U-->>D: JWT donor
    D->>F: Create food listing
    F->>U: Validate donor via gRPC
    F-->>D: Food listing created
    R->>Q: Create food request
    Q->>U: Validate recipient via gRPC
    Q-->>R: Request SEARCHING
    R->>F: Search nearby by request_id
    F->>Q: GetRequest via gRPC
    F-->>R: Nearby food listings
    R->>C: Create claim
    C->>U: Validate recipient
    C->>Q: GetRequest
    C->>F: Get listing and reserve quantity
    C->>Q: MarkRequestClaimed
    C-->>R: Claim WAITING_FOR_PICKUP
    D->>C: Verify pickup code
    C->>Q: MarkRequestCompleted
    C-->>D: Claim PICKED_UP
```

Cancellation is a compensating flow:

```mermaid
flowchart LR
    A[Claim WAITING_FOR_PICKUP] --> B[Cancel Claim]
    B --> C[Release food quantity]
    C --> D[Mark request SEARCHING]
    D --> E[Request can search again]
```

## Ownership

- User Service: authentication, JWT, roles, and user validation.
- Food Service: listings, quantity, nearby matching, reserve/release.
- Request Service: food requests and request status.
- Claim Service: claim lifecycle and pickup verification.
- Contracts: protobuf APIs used between services.

The happy path expects the recipient to create and inspect the claim. Pickup verification is intentionally performed with the donor token because only the donor who owns the listing may confirm the code. The request status progresses `searching` → `claimed` → `completed`; claim status progresses `waiting_for_pickup` → `picked_up`.
