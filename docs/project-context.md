# BagiPangan project context

Snapshot of the project discussions through 9 November 2026. Code references were checked against commit `8ba474e`.

**11 November source audit:** this historical snapshot has been checked against `9a6cfca` and the subsequent local changes. See [findings and verification](audit.md) and [isolated test setup](local-testing.md). Development Vite now forces local proxies; a separate disposable Compose stack creates test databases and migrations. CI now includes frontend and opt-in PostgreSQL regressions. My Listings uses donor-owned discovery. Pickup retries can repair request completion; new code expiry is capped by food availability; cleanup has a separate bounded context. Stock release preserves donor cancellation, and fully reserved listings can expire. The original deployment Compose probes and 10-minute updater cadence remain unchanged. Historical deployment, provider, email and browser-testing reports below have not been independently reverified against production.

**Current local ports:** after observed internal-network publication failure and an 8081 allocation conflict, the isolated stack uses a gateway on host HTTP ports **18081–18084** and PostgreSQL **55433**. Vite development proxies match those host ports. Container/service-to-service ports are unchanged. Older port examples below are historical; use [local-testing.md](local-testing.md) for current commands.

## Project and deployment

BagiPangan began as a Hacktiv8 final project. Donors list available food, recipients search nearby and reserve portions, and donors confirm collection using a pickup code. The project continued after the course as a backend portfolio project. The [README](../README.md) describes the four Go services, their gRPC contracts, PostgreSQL storage, and the React and TypeScript frontend.

The developer described running the backend with Docker Compose on an old MacBook with a broken display, Xubuntu, Wi-Fi, and no battery. Cloudflare Pages serves the [frontend](https://bagipangan.my.id). API traffic reaches the MacBook through Cloudflare Tunnel and Caddy. A power or Wi-Fi outage can therefore leave the homepage visible while API calls fail.

An earlier external check of `GET /api/v1/users/me` and `GET /api/v1/claims` returned the expected HTTP 401 responses without credentials. This established reachability at the time of that check. Each service also has an internal `/readyz` database check. The [Caddy configuration](../zero-hunger-infra/Caddyfile) currently has no public route to those checks.

## Portfolio availability and database hosting

The goal is for a reviewer opening the CV link at any hour to be able to try the donation, claim and pickup flow. Turning off the backend overnight would interrupt that flow even if the homepage remained available.

In October 2026 the developer reported using about 80 of Neon's 100 monthly CU-hours after roughly two weeks, with the backend running except for one night. At a constant 0.25 CU, 14 days of continuous database activity would consume 84 CU-hours. The estimate fits the report, but the live compute size, activity history and quota reset date were not verified. Neon can suspend after five idle minutes and wake on the next database connection. [Neon pricing](https://neon.com/pricing) · [Compute lifecycle](https://neon.com/docs/introduction/compute-lifecycle)

Two sources of activity were identified in the repository:

- All four services in [Docker Compose](../zero-hunger-infra/docker-compose.yml) call `/readyz` every 10 seconds. Each call pings PostgreSQL. If deployed as written, these checks explain how the database could stay awake with no visitors. The existing `/healthz` checks only the service process. Switching frequent checks to that endpoint would also change what Compose waits for when starting services with `depends_on: service_healthy`.
- The [Food Service listing updater](../zero-hunger-food-service/cmd/server/main.go) runs at startup and every 10 minutes to update expired and claimed listings. It can wake the database even after the health checks change. Any scheduling change should preserve listing behavior.

Hosting options discussed in October 2026 remain under review. [Supabase Free](https://supabase.com/pricing) had no equivalent monthly CU-hour cap, a 500 MB database limit per project and 5 GB of monthly egress. A project with low activity over seven days could be paused and require a [dashboard resume](https://supabase.com/docs/guides/platform/free-project-pausing). [Aiven Free](https://aiven.io/docs/products/postgresql/concepts/pg-free-tier) had no fixed runtime limit, 1 GB of disk and a 20-connection ceiling shared across the application's services. Inactive services could be powered off. Self-hosting would require checking available RAM and disk, maintaining PostgreSQL, and testing restoration from an off-host backup. No migration was selected.

**Host identity to confirm.** The earlier deployment account identifies the Xubuntu MacBook. In the hosting discussion the developer called the running backend host a VPS. Confirm whether these refer to the same machine or a later deployment before applying host-specific instructions.

**Database provider to confirm.** During the 26 October local Docker troubleshooting, inspected containers pointed to a Supabase database through its Session pooler. This establishes that test configuration at that time, not a confirmed production migration. Preserve the Neon usage discussion above; determine which provider each live/test environment actually uses before changing credentials or hosting.

## Frontend direction and architecture decisions

The user knows HTML, CSS and JavaScript and is learning TypeScript. For employability, they explicitly prioritized job relevance over ease of learning. The proposed order was React with TypeScript, then Next.js, with Vue deferred unless a target employer uses it. Next.js builds on React. This was learning guidance, not authorization or a confirmed decision to migrate the application. The [current frontend package](../frontend/package.json) uses React, TypeScript, Vite and Leaflet; it does not currently use Next.js or React Router.

Keep `frontend/` beside the service folders in the same repository, with its own dependencies, build and deployment. Repository layout does not require the applications to deploy together. The design exploration used a community food board, warm cream, forest green and peach accents. A softer neo-brutalist mockup added dark borders, small solid offset shadows and bold headings. These were concepts rather than a complete design specification; food photos would require API support beyond the listing model inspected during that discussion.

Microservices were a bootcamp requirement. The architecture discussion treated them as useful for independent ownership, deployment and scaling when those benefits justify distributed failures, networking and operational complexity. A modular monolith was suggested for a hypothetical small-team version, not as a rewrite decision for this project. Extraction was explained as moving a business capability, data ownership and callers incrementally, with recovery for cross-service failures. A monolith and extracted services can coexist.

## Local development and environment boundaries

The [Compose stack](../zero-hunger-infra/docker-compose.yml) expects externally managed PostgreSQL. It does not create databases or apply migrations, and service startup does not run GORM `AutoMigrate`. Apply each service's `migrations/001_init.sql` to a new test database manually. The donor pickup change did not require a schema migration.

An ignored `zero-hunger-infra/.env.local` was proposed to separate test settings from the deployed `.env`. Invoke it explicitly with `docker compose --env-file .env.local ...`; the October inspection found `.env` and no `.env.local`, so check which file exists and is being used. Compose maps these variables into the containers; per-service environment files are not loaded automatically. Environment-only changes require container recreation, while code changes require image rebuilding.

For PostgreSQL installed on the Windows host, containers use `host.docker.internal` and the actual database port. A database in a container instead uses its service name on a shared network. Existing gRPC addresses such as `food-service:50052` must retain Compose service names. Four databases on one PostgreSQL instance were recommended to match the existing migrations. One shared test database with service-owned tables was discussed as possible, but migration object collisions and schema/search-path configuration need checking first. Services should retain API-based data access across ownership boundaries.

For local browser testing, start only `user-service`, `request-service`, `food-service` and `claim-service`; Caddy, Tunnel and Pages are not needed. With `.env.local`, the command from `zero-hunger-infra` is:

```powershell
docker compose --env-file .env.local config --quiet
docker compose --env-file .env.local up -d --build user-service request-service food-service claim-service
docker compose --env-file .env.local ps
```

Configuration validation checks interpolation, not database connectivity or migrations. The existing required tunnel variable still needs a non-empty local placeholder even when the tunnel is not started. Keep real keys, credentials, tokens and certificates out of Git.

The local frontend plan is browser `localhost:3000` -> Vite proxy -> service HTTP ports `8081` (users/auth), `8082` (food), `8083` (requests) and `8084` (claims). The user-specific `/users/{id}/food-requests` proxy must precede the general users rule and target Request Service. Use empty `VITE_API_BASE_URL` in `frontend/.env.development.local`, then restart Vite. **Current-state difference:** the checked-in [Vite configuration](../frontend/vite.config.ts) still proxies to `https://app.bagipangan.my.id`; running the frontend locally alone does not imply a local backend or test database.

### Observed Supabase connection failures

- User Service failed DNS resolution for the direct Supabase hostname. A host DNS check found an IPv6 record but no IPv4 record, consistent with an IPv4/IPv6 compatibility issue in that Docker setup. The Session pooler was recommended for IPv4 connectivity. Copy its complete connection parameters, including its different username, from Supabase.
- Food Service later reached the pooler but failed with `x509: certificate signed by unknown authority`. Both food and claim containers had `DB_SSLMODE=verify-full` with an empty `DB_SSLROOTCERT`. Claim Service was still `Created`, waiting on dependencies, rather than independently demonstrating the same error.
- The proposed fix was to place the provider CA at `zero-hunger-infra/certs/supabase-ca.crt`, keep `FOOD_DB_SSLMODE` and `CLAIM_DB_SSLMODE` at `verify-full`, set both corresponding `*_DB_SSLROOTCERT` values to `/certs/supabase-ca.crt`, and recreate the containers. Compose already mounts that directory read-only. Disabling TLS was only discussed for a local database, not the remote provider.

## GitHub CI understanding

The conversation explained that GitHub Actions workflows in `.github/workflows/*.yml` determine which events run which checks. Even a one-line change can trigger CI; change size does not determine whether it runs. Branch and path filters can limit runs. A `pull_request` workflow normally runs when a PR is opened, reopened or updated with new commits. Configuring both push and PR events for the same branch can produce two runs for an update. The example discussed used pushes to `main` and PRs targeting `main`: a feature-branch push before a PR would not trigger that example, opening or updating the PR would, and merging would trigger the push workflow. This was an explanation, not verification or modification of this repository's workflows. [GitHub workflow triggers](https://docs.github.com/en/actions/how-tos/write-workflows/choose-when-workflows-run/trigger-a-workflow)

Separately, this conversation inspected the [actual CI workflow](../.github/workflows/ci.yml). It runs on pushes to `main` and on `pull_request` without a base-branch filter. It runs Go vet/tests across the five modules, validates Compose and Caddy, and builds images. It does not deploy the VPS or run frontend checks. A PR failure occurred because `.env.example` leaves `CLOUDFLARE_TUNNEL_TOKEN` empty while Compose uses `${CLOUDFLARE_TUNNEL_TOKEN:?set CLOUDFLARE_TUNNEL_TOKEN}`. Parsing the full configuration requires a non-empty value even for validation/build. The suggested build-job `env` placeholder is now present in the checked-in workflow. It is not a real credential; this CI job does not need a production tunnel secret. A successful GitHub rerun was not reported here.

## Release order

Frontend and backend changes share a commit/PR but deploy separately. For this backward-compatible pickup change, the proposed order was: push the tested feature branch -> PR checks -> deploy the tested food/claim backend commit and check health -> merge into the configured Pages production branch -> test the public flow. Checking out the exact tested commit on the backend host temporarily in detached HEAD was suggested, followed by returning that checkout to `main` after merge. Confirm host path, clean working tree and branch settings first; `~/apps/bagipangan` was the README example, not a verified path on every host.

At one inspection, `dev` and its local `origin/dev` tracking ref matched `001b6dd` (`change: pickup flow`), while `main` was one commit behind. This is historical state, not a branch instruction. Pages production settings discussed were root `frontend`, build `npm run build`, output `dist`, and `VITE_API_BASE_URL=https://app.bagipangan.my.id`. Production and preview deployments can have separate environment values; a preview connected to the live API still changes live data. Production database configuration and certificates should stay on the host. No production deployment was performed by the assistant in this conversation.

## Learning while using AI

The developer raised a tradeoff: implementing improvements while learning takes too long, but delegating the full implementation can leave them without understanding it. The suggested workflow is to choose one learning objective per improvement, attempt the core reasoning or implementation first, and ask AI for specific hints when stuck. AI can handle familiar boilerplate, fixtures, mechanical refactors and documentation drafts. Design unfamiliar behavior together and keep generated changes small enough to explain, debug and modify.

For a stock-reservation improvement, the developer would own the quantity invariant and learn how the transaction enforces it; AI could draft supporting setup and concurrent tests. Verification should include explaining failure cases and making a small follow-up change independently. Security, migration and deployment consequences still need review. This is proposed learning guidance, not a requirement to turn every implementation task into an exercise or an instruction authorizing code changes.

## Claim lifecycle

The [HTTP handler](../zero-hunger-claim-service/internal/handler/http.go) takes the signed-in user's ID, parses the request and listing IDs, and passes the requested quantity to the [claim use case](../zero-hunger-claim-service/internal/usecase/claim_usecase.go).

| Action | Order of operations |
| --- | --- |
| Create | Check recipient, request ownership and searching status. Check listing availability. Ask Food Service to reserve stock. Insert the claim as `waiting_for_pickup`. Ask Request Service to mark the request `claimed`. Start the email notification asynchronously. |
| Verify | The donor selects their listing and enters the recipient's six-digit code. `POST /api/v1/claims/verify-pickup` accepts `food_listing_id` and `claim_code`. The backend resolves the claim ID, checks ownership, waiting status, code and expiry, marks the claim `picked_up`, then marks the request `completed`. The older route with a claim ID in the URL also exists. |
| Cancel | Check recipient ownership. Mark the claim `cancelled`, release the reserved quantity, then return the request to `searching`. |

### Donor listing discovery and quantity meaning

The donor verification UI originally required manually entering a claim ID. Using a claim ID was valid for identifying a particular recipient's reservation; the missing piece was discovery. The [current verification page](../frontend/src/pages/donor/VerifyPickupPage.tsx) loads `GET /api/v1/food-listings/mine` and provides one code field per owned listing. It includes listing statuses `available` and `claimed`, including fully reserved food with zero unreserved stock, and hides `expired` listings. The server matches listing plus code to the individual claim; it does not verify every claim on that listing.

Listing `quantity` is remaining unreserved stock, not a preserved original total. A claim reduces it, cancellation releases it, and pickup does not deduct it a second time. The recipient-facing active-listings query requires available status, a current availability window and quantity greater than zero, so it was unsuitable for discovering all donor pickups. Individual donor access to a known claim already existed through ownership checks; `GET /claims` discovers claims through the signed-in user's requests and remains recipient-oriented. A donor incoming-claims list was discussed as a later option, not implemented as part of this flow.

## Email delivery and mailbox placement

The [email adapter](../zero-hunger-claim-service/internal/repository/http/email_services.go) is named `SendGridEmailService` but actually uses the Resend Go SDK. It sends claim/cancellation notifications to `user.GetEmail()` using `Zero Hunger <noreply@bagipangan.my.id>`. Pickup verification itself does not send a notification in the inspected flow. Both send methods currently discard Resend's returned email ID.

Three outcomes must be distinguished: claim API success, Resend request acceptance, and recipient mail-server delivery. Notification sends run in a goroutine; failures are logged without failing the already successful claim operation. There is no durable retry for failed notifications. Logging the provider message ID and checking its dashboard event history were recommended, not implemented here.

The local Docker logs showed a successful claim creation followed by `API key is invalid`; many HTTP 200 entries were `/readyz` health probes. The user confirmed the local key was old. They separately reported a valid deployed key, Resend HTTP 200, correct destination, `Delivered` status, and a message visible in Gmail All Mail but absent from Inbox, Spam and Trash. Those facts establish delivery for that message and are consistent with archiving or a mailbox rule; the specific rule or cause was not confirmed. Do not carry the local invalid-key diagnosis over to the deployed message.

Resend also reported no DMARC record. A Cloudflare DNS TXT record at `_dmarc` with `v=DMARC1; p=none;` was proposed as an initial non-enforcing policy, with SPF/DKIM verification and later review of authentication. It was not established as the cause of missing Inbox placement, and no DNS change was verified. A `p=none` record alone does not guarantee Inbox placement. See [Resend DMARC guidance](https://resend.com/docs/dashboard/domains/dmarc).

## Link-preview proposal

To show an image when sharing the public URL, add Open Graph tags directly to [frontend/index.html](../frontend/index.html), with absolute public HTTPS URLs, and provide a 1200 x 630 PNG/JPG such as `frontend/public/og-image.png`. Suggested fields were title, type `website`, site name, canonical page URL, description, image URL, image dimensions and image alt text. Static HTML allows preview crawlers to read metadata without running React. Preview caches and platform display choices can delay or alter results; per-listing previews would need additional work. [Open Graph protocol](https://ogp.me/)

This remains a proposal: at commit `8ba474e`, `index.html` still has title `frontend`, no Open Graph tags, and the public directory has no `og-image.png`. No backend change is required for a site-wide preview.

## Concurrency and failure handling

The [food repository](../zero-hunger-food-service/internal/repository/db/food_listing_repository.go) locks the listing row and rechecks availability before deducting stock. If two recipients try to reserve the last portion, the second reservation sees the first one's committed change. Verification and cancellation similarly lock the same [claim row](../zero-hunger-claim-service/internal/repository/db/claim_repository.go) and require `waiting_for_pickup`. This makes only one of those competing transitions eligible to succeed.

The food transaction commits before the claim is inserted. If the insert returns an error, the claim use case already calls `ReleaseFoodQuantity`. If marking the request claimed fails, the use case attempts to delete the claim and release stock. Several failure cases remain relevant:

- A crash after reservation can leave stock reduced without a saved claim. A failed release during cleanup is logged.
- Cancellation can commit before stock is returned. Pickup can commit before the request is completed. A retry can encounter a claim that has already changed status.
- Food Service accepts a claim ID on reserve and release, but its repository updates stock by listing ID and quantity. Replaying the same release directly can add stock twice. A second normal cancel request is guarded by the claim status check.
- A donor editing an absolute quantity based on an earlier reading can overwrite an intervening stock change.

The discussion prioritized testing the existing cleanup path for the junior portfolio. Tracking reservations by claim ID and persisting unfinished operations for retry remain possible later improvements.

## Why the My Claims page changed

The earlier page kept claim IDs in browser storage and called `GET /claims/{id}` for each one using `Promise.all`. Three stored IDs appeared as three UUID-named requests in DevTools. The calls ran concurrently, but rendering waited for the slowest one. Reported durations ranged from hundreds of milliseconds to about 2.5 seconds. The cause of that latency was not measured.

The [current page](../frontend/src/pages/recipient/MyClaimsPage.tsx) calls [getClaimList](../frontend/src/api/claims.ts) once per load. The authenticated `GET /api/v1/claims` endpoint finds the user's requests and retrieves their claims together. This also lets a recipient load their claims from another browser. The earlier loader remains in a comment for comparison.

For reading the React code, follow `getClaimList() → setClaims(results) → claims.map(...)`. The TypeScript return type `Promise<FoodClaim[]>` describes an array of any length. `[FoodClaim]` would describe a tuple with exactly one element. The frontend production build passed when the change was made. No post-change browser timing was recorded in the discussion.

## Test evidence and next exercise

The developer reported that the manual scenarios behaved as expected, deliberate invalid actions were rejected, and stock quantities matched. To repeat the check, use separate donor and recipient sessions with test data. Start with five portions, claim two, and check that three remain. A wrong pickup code should leave the claim waiting. A correct code should complete pickup with three portions still available. On a fresh equivalent claim, cancellation should restore five portions, and a repeated cancellation should leave five.

The [Go use-case tests](../zero-hunger-claim-service/internal/usecase/claim_usecase_test.go) cover successful creation, cancellation, and several pickup-code cases using mocks. From the repository root in PowerShell, run:

```powershell
Set-Location .\zero-hunger-claim-service
go test ./internal/usecase -v -count=1
```

The [Postman collection](../zero-hunger-infra/docs/Zero%20Hunger%20E2E.postman_collection.json) exercises the API across services. Its pickup step uses the claim-ID route. Its recipient-create-listing test currently accepts either 403 or 500, so a passing result there can hide a server error.

The creation flow has already been traced. The next testing exercise is to make the mock claim insert fail after a successful reservation and check that cleanup releases the correct quantity and returns the original error. A separate, smaller Go exercise is to extract the existing input guard and test missing IDs and non-positive quantities. PostgreSQL concurrency tests would provide separate evidence for row-lock behavior.
