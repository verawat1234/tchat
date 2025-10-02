# Railway Deployment Guide

This guide explains how to deploy the eight Tchat backend services into the single Railway project **Tchat** (project ID `0a1f3508-2150-4d0c-8ae9-878f74a607a0`). All manifests live in `railway/services`, and matching environment templates live in `railway/env`.

## Service Inventory

| Service           | Railway Service Name | Port | Manifest Path                              | Default Database |
|-------------------|----------------------|------|--------------------------------------------|------------------|
| Auth API          | `auth-service`       | 8081 | `railway/services/auth-service.toml`       | `auth_service_db`|
| Messaging API     | `messaging-service`  | 8082 | `railway/services/messaging-service.toml`  | `messaging_service_db`|
| Payment API       | `payment-service`    | 8083 | `railway/services/payment-service.toml`    | `payment_service_db`|
| Commerce API      | `commerce-service`   | 8084 | `railway/services/commerce-service.toml`   | `commerce_service_db`|
| Notification API  | `notification-service`| 8085 | `railway/services/notification-service.toml`| `notification_service_db`|
| Content API       | `content-service`    | 8086 | `railway/services/content-service.toml`    | `content_service_db`|
| API Gateway       | `gateway-service`    | 8080 | `railway/services/gateway-service.toml`    | —                |
| Video API         | `video-service`      | 8091 | `railway/services/video-service.toml`      | `video_service_db`|

## Prerequisites

1. Install the Railway CLI v4.5 or newer (`npm install -g @railway/cli`).
2. Authenticate: `railway login`, confirm with `railway whoami`.
3. From the repository root run `cd backend`.
4. Ensure the repository is clean (`git status` succeeds); Railway builds from git snapshots.

## Step 1 — Deploy Services

Automated deployment is handled by `deploy-to-railway.sh`:

```bash
cd backend
./deploy-to-railway.sh
```

The script will:

- Link the CLI to the `Tchat` project and `production` environment.
- Provision shared Postgres/Redis plugins if they do not already exist.
- Ensure each Railway service (`auth-service`, `messaging-service`, …) exists.
- Copy the corresponding manifest under `railway/services` to `railway.toml` and trigger `railway up --service <name> --detach` for every backend.

If you need to redeploy a single service, copy the relevant manifest manually and run:

```bash
cp railway/services/auth-service.toml railway.toml
railway up --service auth-service --environment production --detach
rm railway.toml
```

> **Tip:** use `railway logs --service <name> -d` to follow the latest build and deployment output.

## Step 2 — Apply Environment Variables

Environment templates live in `railway/env/*.env.template`. Sync them after the services finish building:

```bash
./setup-railway-env.sh           # apply all templates
./setup-railway-env.sh gateway-service  # or target one service
```

The helper script links to the `Tchat` project, verifies the service exists, and calls

```bash
railway variables --service <name> --environment production --set KEY=value
```

for every entry in the template. Update the placeholders (JWT secrets, Stripe keys, inter-service URLs, etc.) in Railway once public endpoints are available.

## Step 3 — Database Preparation

All services share the same Postgres plugin. Each template sets a unique `DB_DATABASE`/`DB_NAME` value (for example `auth_service_db`). Create these databases before running migrations:

```bash
railway run --service Postgres --environment production -- \
  psql -c "CREATE DATABASE auth_service_db;"
```

Repeat for each database name. Redis namespaces are separated via `REDIS_DATABASE` per template.

## Step 4 — Post-Deployment Checklist

- [ ] Record service URLs via `railway status --json` (each service exposes a `publicUrl`).
- [ ] Update downstream `*_SERVICE_URL` variables in the gateway and dependent services.
- [ ] Run migrations/tests per service, e.g. `railway run --service auth-service -- make migrate`.
- [ ] Curl `/health` for every backend plus the gateway.
- [ ] Configure secrets for external integrations (Twilio, Stripe, Kafka, ScyllaDB, etc.).

## Troubleshooting

| Issue | Resolution |
|-------|------------|
| `railway up` fails immediately | Re-run with `railway up --service <name>` and inspect build logs via `railway logs --service <name> -b` |
| Container starts then stops | Check service logs (`railway logs --service <name>`) for runtime errors, verify environment variables |
| Database authentication errors | Ensure the Postgres plugin is present and templates were applied; recreate databases if needed |
| Gateway upstream errors | Confirm each downstream service URL/environment variable is correct and health endpoints respond with 200 |

## Maintenance Notes

- Keep manifests (`railway/services/*.toml`) aligned with Dockerfile locations when new services are added.
- Extend `deploy-to-railway.sh` and `setup-railway-env.sh` when introducing additional services.
- Use `railway status --json` to script health checks or capture deployment metadata.
- Remove unused services (for example prototypes) from the Railway dashboard to avoid drift and unnecessary spend.

Following these steps results in a fully deployed backend stack inside the single Railway project with readable service, environment, and database naming.
