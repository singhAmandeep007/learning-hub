# Learning Hub API Contract

This package is the single source of truth for API contract definitions.

## What lives here

- OpenAPI spec: `openapi/openapi.yaml`
- Type generation scripts for consumer packages (`frontend`, `e2e`)

## Commands

```bash
npm install
npm run lint
npm run generate
```

## Generated outputs

- `../frontend/src/types/api-contract.generated.ts`
- `../e2e/helpers/types/api-contract.generated.ts`
- `../backend/contract/types.gen.go`
- `../backend/contract/openapi.yaml`

## Usage rules

1. Update only `openapi/openapi.yaml` manually.
2. Regenerate via `npm run generate` (or `make contract-types`).
3. Keep generated files committed so consumers build without local codegen.
