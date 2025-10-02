# Railway Deployment Implementation Plan

**Last updated:** 2025-10-02

---

## 🎯 Summary

- ✅ Eight backend services deployable via `./deploy-to-railway.sh`
- ✅ Railway configuration files consolidated in `backend/railway/services`
- ✅ Environment templates renamed with readable keys in `backend/railway/env`
- ✅ Database identifiers standardised to `<service>_service_db`
- ✅ Helper tooling updated (`deploy-to-railway.sh`, `setup-railway-env.sh`)

---

## 📋 Service Matrix

| # | Service | Railway Project | Port | Config Path | Database |
|---|---------|-----------------|------|-------------|----------|
| 1 | Auth API | `tchat-auth-service` | 8081 | `railway/services/auth-service.toml` | `auth_service_db` |
| 2 | Messaging API | `tchat-messaging-service` | 8082 | `railway/services/messaging-service.toml` | `messaging_service_db` |
| 3 | Payment API | `tchat-payment-service` | 8083 | `railway/services/payment-service.toml` | `payment_service_db` |
| 4 | Commerce API | `tchat-commerce-service` | 8084 | `railway/services/commerce-service.toml` | `commerce_service_db` |
| 5 | Notification API | `tchat-notification-service` | 8085 | `railway/services/notification-service.toml` | `notification_service_db` |
| 6 | Content API | `tchat-content-service` | 8086 | `railway/services/content-service.toml` | `content_service_db` |
| 7 | API Gateway | `tchat-gateway-service` | 8080 | `railway/services/gateway-service.toml` | — |
| 8 | Video API | `tchat-video-service` | 8091 | `railway/services/video-service.toml` | `video_service_db` |

---

## 🚀 Phase 1 — Prepare the Environment

```bash
npm install -g @railway/cli   # Railway CLI
railway login                 # Authenticate
railway whoami                # Verify session

cd backend
ls railway/services           # Expect eight *.toml files
ls railway/env                # Expect eight *.env.template files
```

---

## 🧭 Phase 2 — Deploy via Automation

```bash
./deploy-to-railway.sh
```

What the script does for each service:

1. Copies the matching TOML file to `railway.toml` temporarily.
2. Creates (or links) the target Railway project.
3. Attaches PostgreSQL and Redis plugins when flagged.
4. Seeds critical variables (`APP_ENVIRONMENT`, `SERVICE_PORT`, `DB_DATABASE`).
5. Executes `railway up --detach`.

Monitor the CLI output for warnings; the script will continue if a plugin already exists.

---

## 🔧 Phase 3 — Apply Environment Variables

```bash
# Apply templates for every project
./setup-railway-env.sh

# or target a single service
./setup-railway-env.sh payment-service
```

Each template uses descriptive keys and maps plugin variables automatically:

- Database aliases (`DB_*`, `DATABASE_*`, `<SERVICE>_DATABASE_NAME`)
- Service-specific port variables (`AUTH_SERVICE_PORT`, `PAYMENT_SERVICE_PORT`, ...)
- Redis namespaces per service (0–6)
- Placeholders for downstream service URLs to be filled with live Railway domains

If automatic linking fails, run `railway link` manually, then rerun the script for the affected service.

---

## 🧪 Phase 4 — Validate Deployments

1. Retrieve project URLs: `railway status --service <project>`
2. Update gateway variables with the discovered URLs.
3. Hit `/health` for every service, then `/health` on the gateway.
4. Execute smoke tests (for example `make test-e2e` if available).
5. Run migrations where applicable, e.g. `railway run --service tchat-auth-service -- make migrate`.

---

## 📈 Monitoring & Maintenance

- Use `railway logs --service <project>` to inspect runtime logs.
- Set alerting and metrics via Railway dashboard; variables `METRICS_ENABLED=true` are pre-populated.
- When adding a new service, create a TOML + env template and extend both helper scripts.
- Keep Dockerfiles patched; rebuild with `railway up` whenever dependencies change.

---

## ✅ Final Checklist

- [ ] All eight services deployed successfully.
- [ ] Environment variables synchronised from templates.
- [ ] Inter-service URLs configured in the gateway and dependent services.
- [ ] Schema migrations executed.
- [ ] Health endpoints returning HTTP 200.
- [ ] Monitoring/alerting configured (Railway, external tooling).

---

Once the checklist is complete the backend stack is fully operational on Railway with consistent naming across services, environments, and databases.
