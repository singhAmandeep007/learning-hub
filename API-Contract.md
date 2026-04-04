# Learning Hub API Contract

The API contract source of truth is now maintained in the dedicated package:

- `api-contract/openapi/openapi.yaml`

Generated outputs:

- `frontend/src/types/api-contract.generated.ts`
- `e2e/helpers/types/api-contract.generated.ts`
- `backend/contract/types.gen.go`
- `backend/contract/openapi.yaml`

Use these commands from `v1` root:

```bash
make contract-lint
make contract-types
```

Notes:

- Update only `api-contract/openapi/openapi.yaml` manually.
- Do not edit generated files by hand.
- Backend handlers and models should continue matching the OpenAPI schema.
