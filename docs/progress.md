# BagiPangan progress

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
