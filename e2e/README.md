# Learning Hub E2E (Playwright)

Dedicated end-to-end test module for validating full frontend + backend flows.

## What these tests validate

- API request payloads for user journeys (multipart form fields)
- API response status and JSON body expectations
- Frontend UI state updates after mutations

## Helper module structure

- `helpers/config.ts`: environment-aware API/product path helpers
- `helpers/api.ts`: reusable API waiters + multipart parsing + generic JSON response checks
- `helpers/resources.ts`: resource-specific helpers (form fill, API setup/cleanup, schema assertions)
- `helpers/index.ts`: single import surface for tests

Response validation uses `zod` schemas in helpers to keep assertions explicit and reusable.

## Prerequisites

- Docker + Docker Compose (recommended)
- or Node.js 22+ with app services running locally

## Run locally (recommended)

```bash
make e2e-docker
```

This brings up Firebase emulator, backend, frontend, and runs Playwright tests in a dedicated `e2e` container.

The Docker run uses a single standalone compose file: `docker-compose.e2e.yml`.

Playwright artifacts are written to `e2e/result/` on the host:
- `e2e/result/playwright-report`
- `e2e/result/test-results`

## Run locally without Docker

1. Start app stack from root:

```bash
make dev-local
```

2. In another terminal:

```bash
cd e2e
npm ci
npm run install:browsers
E2E_BASE_URL=http://localhost:3000 E2E_API_BASE_URL=http://localhost:8000 E2E_PRODUCT=ecomm npm test
```

## Visual regression

Visual specs are in `tests/resources.visual.spec.ts`.
The test runs a serial CRUD flow against the real backend and captures snapshots after:
- empty state
- create
- view details
- update
- delete

The suite clears backend state through APIs before/after each run so snapshots stay deterministic.

Run with Docker from repo root:

```bash
make e2e-docker-vrt
make e2e-docker-vrt-update
```

Run against locally running services from repo root:

```bash
make e2e-local-vrt
make e2e-local-vrt-update
```

Or run directly inside `e2e/`:

```bash
npm run test:visual
npm run test:visual:update
```

## CI usage

- GitHub Actions: `.github/workflows/e2e.yml`
- Jenkins declarative pipeline: `jenkins/Jenkinsfile.e2e`

Both run the exact same Docker Compose command for consistency.
