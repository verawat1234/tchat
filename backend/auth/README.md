# Auth Service

Minimal Go HTTP scaffold for authentication flows. Exposes `/health` and `/ready` endpoints and sets the stage for OTP/Magic Link logic.

## Run locally

```bash
go run ./cmd/authsvc
```

Environment variables:

- `AUTH_HTTP_ADDR` — address to bind (default `:8080`).
- `AUTH_TOKEN_ISSUER` — token issuer URL (default `https://auth.tchat.dev`).
