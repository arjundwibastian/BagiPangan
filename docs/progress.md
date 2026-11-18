# BagiPangan progress

## 2026-11-11 — Dedicated host ports after gateway collision

The user-reported gateway startup failed because `127.0.0.1:8081` was already allocated. The owning process/container has not been established; earlier containers can retain requested bindings, and another stack may also own the port. No existing service was stopped. Moved only local gateway host bindings to 18081–18084 and database 55433, keeping container ports and internal service addresses unchanged. Updated Vite's development proxies, their regression expectations and the local guide.

The recovery command recreates only `gateway` (`--no-deps --force-recreate`), preserving all five running app/database containers and their fixtures. Frontend and Compose checks are recorded in the audit. Runtime port availability and browser E2E still need verification; this change does not resolve the agent's independently observed permission restrictions.

## 2026-11-11 — Running containers exposed a local networking defect

User-supplied `compose ps` confirms all four app containers running and PostgreSQL healthy, but shows no host-port mappings. Source and parsed Compose configuration confirm the initial test file attached every container only to an internal network. Docker publishes ports through a default-gateway network, so this setup accepts `ports` yet does not publish them. This was a defect in the added test setup, separate from the agent's recorded permission restrictions; earlier parsing checks could not detect it. See [Docker maintainer explanation](https://github.com/moby/moby/discussions/53256).

Added a local HAProxy gateway on both internal and normal bridges; only it publishes the five loopback listeners. Apps and PostgreSQL remain internal-only. [Local guide](local-testing.md) includes a gateway-only start command that preserves the already-running database and its temporary fixtures. Compose parsing and topology checks are recorded in the audit; proxy runtime and browser E2E still require verification. The user output establishes running images, but a fresh build log has not been reviewed.

## 2026-11-11 — Follow-up after user started the local stack

The user reported running the local Compose command externally. Agent readiness requests to loopback ports 8081–8084 all failed with socket access forbidden; Docker inspection remains denied. Both opt-in PostgreSQL suites were now explicitly attempted against the named loopback fixture databases and **failed during fixture connection**, before schema creation or repository assertions. This does not establish a backend/database application defect or verify container startup.

Vite started on loopback port 3011, but the in-app browser timed out reaching it; an initial browser tab also failed to attach. The temporary Vite process was stopped. Browser E2E remains blocked by session connectivity, and no test accounts or claims were created. The earlier assurance that external stack startup alone would unblock testing was premature: local-network/browser connectivity is also required. External build/start results remain user-reported until their output can be reviewed.

## 2026-11-11 — Image-build and browser-E2E retry

Retried the standalone local stack without reading deployment credentials. Docker Desktop is now running, but the agent session cannot access `npipe:////./pipe/dockerDesktopLinuxEngine` (`permission denied`). The four-service `docker compose --env-file .env.example -f docker-compose.local.yml build` attempt also stopped before compilation because the default Buildx lock is inaccessible. Retrying with `BUILDX_CONFIG` in scratch stopped on a denied Buildx activity-file rename. `up -d --no-build` independently failed on engine access.

**Result:** image builds are blocked by environment permissions, not demonstrated Dockerfile/compiler failures. Database-backed browser E2E could not start; no browser-flow pass or failure was recorded. No application code, ignored settings, production data or Docker security settings were changed. The session needs authorized Docker access before either requested run can complete.

## 2026-11-11 — Repository audit and local improvements

Read both project notes and checked claims against source at `9a6cfca`; working tree was initially clean. [Audit report](audit.md) records exact locations, severity, reproduction steps, confirmed fixes, remaining consistency bugs and unverified operational risks. [Local testing](local-testing.md) describes the new standalone stack; no credentials or deployed data were read.

- Baseline: all five Go modules passed tests/vet; frontend built; lint failed on the three documented state-in-effect errors.
- Batch 1: local-only dev proxies, HTTP-status-based session expiry, useful gateway errors, owned listing discovery, lint fixes, frontend regression tests and CI checks. Frontend tests/lint/build pass.
- Batch 2: owner/admin-only user profiles and ownership checks on request-based nearby searches. Denied-access tests pass. Profile tests fail on an isolated snapshot of the original source.
- Batch 3: bounded cleanup despite caller cancellation, pickup expiry limited by food availability, same-code pickup/completion retries, and repeated-cancellation protection tests. Cleanup and expiry tests fail on the original snapshot and pass with fixes; claim/request tests/vet pass.
- Batch 4: release preserves cancelled listings and rejects overflow; fully reserved food expires; edits use a server snapshot check; cancellation condition is rechecked in SQL; query cancellation propagates. Added isolated PostgreSQL concurrency/expiry suites and CI configuration. Suites compile and skip locally because no database is running; their runtime behavior is **not yet verified**.
- Batch 5: nearby UI discards consumed requests/stale results and validates whole quantities; forbidden searches remain visible errors; empty nearby search returns `[]`; finite-number validation and zero-coordinate input tests; Postman negative test now requires 403. Food/request tests/vet and frontend checks pass.

Both Compose files parse with example settings (deployment file uses a non-secret tunnel placeholder). Docker engine, psql and local race support are unavailable; image builds, Caddy runtime validation, PostgreSQL execution, browser E2E and GitHub CI execution were not performed. The audit did not deploy, push, rewrite the architecture or alter existing ignored environment files. High-priority open work: durable reservation identity/recovery and expired-claim reclamation; see report for failure scenarios.

## 2026-11-09 — Conversation context consolidated

Earlier discussions clarified CI triggers and proposed a way to balance AI-assisted delivery with learning. The [project context](project-context.md#github-ci-understanding) records the CI example and [learning guidance](project-context.md#learning-while-using-ai). No repository workflow inspection, CI change or implementation was performed in those discussions. The developer has not explicitly adopted the proposed learning workflow.

**Hosting recommendation to review.** An earlier assistant recommendation was to keep Neon unless cost or measured performance justified moving, and prioritize integration tests, reliable claims and performance evidence. The availability review below records a concrete reported CU-hour concern and raises availability ahead of the cleanup exercise. These are context-dependent recommendations, not an agreed migration or a settled priority order. Review them using live usage and host evidence; preserve the existing open cleanup test.

## 2026-11-09 — Portfolio availability review

The developer reported that the main application flows ran smoothly in their own testing. Independent user testing remains difficult to arrange. An agent review suggested additional error handling, which the developer considered low urgency. An overlapping last-portion claim test can be scripted with concurrent requests. That test was discussed but no result was reported.

The hosting investigation identified recurring database activity in the local configuration. The [project context](project-context.md#portfolio-availability-and-database-hosting) records the relevant files, the reported usage and the hosting alternatives. The discussion produced a diagnosis to verify. No probe change, deployment or migration was reported.

**Next operational check.** Confirm which host runs the live backend and whether its containers use the repository's 10-second database checks. Check the production branch's Active/Idle state and billing reset date in Neon. If the live settings match, review startup dependencies before switching frequent checks to `/healthz`. Then observe Neon across multiple 10-minute listing-update cycles during a quiet period and compare CU-hour consumption before choosing a host.

**Priorities to review.** The existing engineering exercise is to test stock release after a failed claim insert. This conversation raised availability as a more immediate concern for CV visitors. The cleanup test remains open. Once availability is settled, a five-minute donor-to-pickup walkthrough was proposed as a portfolio milestone before deciding whether to put the MVP into maintenance.

## 2026-10-01 to 2026-10-03 — Frontend concept and learning direction

Proposed community food-card and simple-list designs, then generated a softer neo-brutalist image concept using cream, forest green, dark borders and small offset shadows. Recommended a sibling `frontend/` folder in the same repository. For job-focused learning, proposed React/TypeScript before Next.js, with Vue deferred unless a target role needs it. The implemented stack remains React/TypeScript/Vite; no Next.js migration was agreed. Microservices were clarified as a bootcamp requirement with learning value, rather than a universal default or a reason to rewrite the existing project as a monolith.

## 2026-10-26 — Donor verification flow implemented and tested

Replaced manual claim-ID entry with owned-listing cards and one six-digit pickup-code field per listing. Added donor-only listing discovery and listing-plus-code verification while keeping the older claim-ID route. Fully reserved listings remain eligible for display; backend authorization, expiry and claim-state checks still apply. See [donor discovery and quantity semantics](project-context.md#donor-listing-discovery-and-quantity-meaning) for the API distinction.

At implementation time, the frontend build, food/claim service test suites and changed-file lint passed. A headless browser check with mocked APIs passed owned-listing filtering, fully reserved visibility, absence of claim-ID input, short/incorrect code handling, successful leading-zero code submission, input reset, mobile width and empty state. These browser checks did not prove real email delivery or a production rollout. Full frontend lint reported three pre-existing state-in-effect errors in My Listings, My Claims and My Requests; their current status has not been rechecked in this documentation update.

The user later reported that local testing worked well and asked how to release both applications. The backend-first release sequence is recorded in [project context](project-context.md#release-order). At the release-advice inspection, `dev` contained commit `001b6dd`; a subsequent successful merge or deployed commit was not verified in this conversation. Do not treat that historical branch state as current.

## 2026-10-26 — Local connection, email and CI troubleshooting

- **Database setup:** confirmed manual migrations and externally managed PostgreSQL. Proposed ignored `.env.local` plus local Vite proxies. The inspected environment actually used `.env` and Supabase; the checked-in Vite configuration still targets the public API as of `8ba474e`. Local frontend execution alone is not evidence of environment isolation.
- **Supabase connectivity:** diagnosed the user container's direct-hostname DNS failure and suggested the IPv4-compatible Session pooler. The request container was initially only `Created`, not successfully connected. Food later failed TLS CA verification with `verify-full` and an empty root-certificate setting; Claim had not started because of dependencies. Provided CA-mount/environment instructions. The user subsequently reported successful local code testing, but individual post-fix container health results were not captured here.
- **Email:** local logs identified an invalid Resend key, which the user confirmed was old. The deployed-message investigation instead ended with user-reported `Delivered` status and presence in Gmail All Mail. Specific mailbox filters remain unverified. Provider-ID logging, durable retries and the proposed DMARC record remain follow-up work, not completed fixes.
- **GitHub CI:** diagnosed required `CLOUDFLARE_TUNNEL_TOKEN` interpolation while validating/building Compose. Suggested a non-secret build-job placeholder. Current workflow inspection confirms that placeholder is now present; no successful CI rerun was reported in this conversation. Production credentials are not required for this validation/build-only job.
- **Link preview:** supplied Open Graph HTML and an image-file convention. Current source still lacks the tags and preview image, so implementation/deployment remains pending.

## 2026-11-09 — This conversation added to shared documentation

Preserved the earlier CI explanation, claim consistency notes, availability investigation, My Claims work and open cleanup exercise. Added source-checked operational details without applying application, workflow, database, DNS or deployment changes in this update. No build or runtime tests were rerun for this documentation-only change.

**Conflicts and decisions for review:**

- The historical Xubuntu MacBook versus "VPS" host description remains unresolved; verify the actual production host before using host-specific deployment instructions.
- Local Supabase evidence and the separate Neon availability discussion do not establish one universal/current database provider or an agreed migration. Record live/test provider assignments before changing them.
- Next.js was a proposed learning progression, while the current frontend remains Vite. The softer neo-brutalist mockup was exploratory, not a complete approved redesign.
- `.env.local` and local Vite proxies were recommended, but inspected file/runtime states differed. Confirm the intended test environment before exercising data-changing flows.
- The generic CI example mentioned PRs targeting `main`; the actual workflow's `pull_request` event has no base-branch filter. The earlier note about no CI inspection applies to those other discussions; this conversation did inspect the workflow and its placeholder is present now.
- DMARC absence was a provider insight, not a demonstrated cause of archiving. DNS changes and stricter DMARC enforcement were not confirmed.
