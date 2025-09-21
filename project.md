
# Telegram SEA Ecosystem — Master Backlog (Markdown)
**Version:** v1.2 • **Date:** 2025-09-20 • **Owner:** TIKK  
**Columns:** Backlog → Ready → In‑Progress → Review → Done

---

## Legend
- **Labels:** `[BE]` backend (Go), `[FE-Web]` Next.js web, `[Android]` Kotlin, `[iOS]` Swift, `[AI]` ML/AI pipelines, `[Infra]` platform/DevOps, `[Sec]` security/compliance, `[Data]` analytics/warehouse, `[Prod]` product/PM, `[QA]` quality.
- **Format:** *Epic* → **Stories** (with checklists of **Subtasks**).  
- Use GitHub‑style checkboxes. Copy into GitHub/Notion/Obsidian as‑is.

---

# EPIC 1 — Foundations & Environments [Infra]
### Story 1.1 — Monorepo & CI/CD baseline [Infra]
- [ ] Init monorepo structure (apps/, services/, libs/, infra/) [Infra]
- [ ] Configure CI (lint, test, build matrices; Go/Node toolchains) [Infra]
- [ ] CODEOWNERS + branch protections + conventional commits [Infra]
- [ ] Pre‑commit hooks (lint-staged, gofmt, golangci-lint) [Infra]
- [ ] ArgoCD bootstrap (apps-of-apps, environment overlays) [Infra]

### Story 1.2 — Multi‑env setup (dev/stg/prod per region) [Infra]
- [ ] Namespaces & RBAC per env (TH‑BKK, SG‑1 seed) [Infra]
- [ ] Secrets via SOPS + KMS; sealed-secrets for bootstrap [Infra][Sec]
- [ ] Env config: feature flags, connector toggles [Infra]
- [ ] Cost/limits: resource quotas, HPA/VPA defaults [Infra]

### Story 1.3 — Observability baseline [Infra][Data]
- [ ] Deploy Grafana/Loki/Tempo + OpenTelemetry collectors [Infra]
- [ ] RED dashboards per service + SLO burn‑rates [Data]
- [ ] Synthetic probes from TH/SG POPs (HTTP/GRPC) [Infra]
- [ ] Alert routes & on‑call rotations (Pager) [Infra]

---

# EPIC 2 — Product IA & Design System [Prod]
### Story 2.1 — Top Tabs IA: Chat | Store | Social [Prod]
- [ ] Finalize navigation flows (web/android/ios) [Prod]
- [ ] Define empty/error/loading states per tab [Prod]
- [ ] Country‑level feature flags (TH default on) [Prod][Infra]

### Story 2.2 — Design system v1 [Prod]
- [ ] Tokens (typography/spacing/colors) + dark mode [Prod]
- [ ] Component parity (Web/ShadCN, Android/Compose, iOS/SwiftUI) [Prod]
- [ ] Accessibility & localization checklists [Prod]

### Story 2.3 — Onboarding Assistant PRD [Prod]
- [ ] User & merchant journeys (language, Buddhist calendar, payments bind, privacy) [Prod]
- [ ] Copy in TH/EN + tone of voice guidelines [Prod]
- [ ] Telemetry plan (activation events) [Prod][Data]

---

# EPIC 3 — Identity & Auth (OTP + Magic Link) [BE][Infra][Sec]
### Story 3.1 — Auth service skeleton (Go) [BE]
- [ ] gRPC + REST gateway (`StartAuth`, `Verify`) [BE]
- [ ] Audit logs + idempotency keys [BE]
- [ ] Token model (short access, rotating refresh) [BE]

### Story 3.2 — OTP & Magic Link providers [BE][Infra]
- [ ] SMS/Email providers with sandbox/live keys [Infra]
- [ ] Delivery callbacks & retries; JTI replay‑proof links [BE]
- [ ] Rate limits, Turnstile/reCAPTCHA, device/IP reputation [Sec][BE]

### Story 3.3 — Client SDKs (Web/Android/iOS) [FE-Web][Android][iOS]
- [ ] Secure storage (Keychain/Keystore/WebCrypto) [FE-Web][Android][iOS]
- [ ] Android SMS retriever/auto‑fill; iOS Universal Links [Android][iOS]
- [ ] Next.js server actions for auth exchange [FE-Web]

---

# EPIC 4 — Messaging Core & Presence [BE]
### Story 4.1 — Timelines in Scylla & DAL [BE]
- [ ] Schema (partition by peer_id; seqno clustering) [BE]
- [ ] Seqno allocator & idempotent write path [BE]
- [ ] Unit tests for hot path & error taxonomy [BE]

### Story 4.2 — Outbox/ACK + Delta Sync [BE][FE-Web][Android][iOS]
- [ ] Local outbox & resend with backoff [Android][iOS][FE-Web]
- [ ] Delta window negotiation & conflict rules [BE]
- [ ] Typing indicators & read cursors [BE]

### Story 4.3 — Push Notifications [BE][Android][iOS][FE-Web]
- [ ] Token registration & topic mapping [BE]
- [ ] FCM/APNs/WebPush with collapse keys [BE]
- [ ] Silent push background sync (mobile/web) [Android][iOS][FE-Web]

---

# EPIC 5 — Media Pipeline [BE][Infra]
### Story 5.1 — Resumable Uploads & CDN [BE]
- [ ] Chunked upload API + pre‑hash dedupe [BE]
- [ ] Signed URLs via CDN; range requests [Infra]
- [ ] AVIF/WebP, Opus, placeholder/blurhash [BE]

### Story 5.2 — Transcode Workers & Ladders [Infra][BE]
- [ ] HLS/DASH packaging; low‑bitrate ladders for 3G [Infra]
- [ ] Thumbnail sprites; error handling & retries [BE]
- [ ] Observability (queue depth, transcode time) [Infra]

---

# EPIC 6 — Payments, Wallets & Ledger [BE][Sec][Data]
### Story 6.1 — Ledger v1 (double‑entry) [BE][Data]
- [ ] Append‑only entries; invariants dashboard [BE][Data]
- [ ] Daily close & reconciliation jobs [BE]
- [ ] Payout files & audit exports [BE]

### Story 6.2 — PromptPay (TH) Connector [BE][Sec]
- [ ] Static/dynamic QR; signature verification [BE]
- [ ] Webhook verification + idempotency [BE][Sec]
- [ ] Refunds & dispute skeletons [BE]

### Story 6.3 — Wallet service & limits [BE][Sec]
- [ ] KYC tiers & velocity checks [BE]
- [ ] Freeze/unfreeze & sanctions hooks [Sec]
- [ ] Events for analytics & fraud models [Data][AI]

### Story 6.4 — Payment Intents API [BE][FE-Web]
- [ ] Create/confirm/capture + error taxonomy [BE]
- [ ] 3DS/OTP hooks & client secrets [BE][FE-Web]
- [ ] Webhooks → Orders service handoff [BE]

---

# EPIC 7 — Merchants, Orders & Storefront [BE][FE-Web]
### Story 7.1 — Merchant onboarding [FE-Web][BE]
- [ ] Business profile, verification badge [FE-Web]
- [ ] Bank account binding & settlements [BE]
- [ ] Role/permissions (owner, staff) [BE]

### Story 7.2 — Catalog & inventory [FE-Web][BE]
- [ ] CSV/Google Sheets import [FE-Web]
- [ ] Variants, stock adjustments, price rules [BE]
- [ ] Product media & SEO fields [FE-Web]

### Story 7.3 — Chat‑to‑cart flow [FE-Web][Android][iOS][BE]
- [ ] Product sheet in chat + inline cart [FE-Web]
- [ ] Order create & payment handoff [BE]
- [ ] Receipt thread & status updates [Android][iOS]

### Story 7.4 — Logistics integration [BE]
- [ ] Ninja Van/J&T/Kerry adapters [BE]
- [ ] Tracking webhooks + label PDFs [BE]
- [ ] SLA/exception dashboards [Data]

---

# EPIC 8 — Mini‑Apps Platform [BE][FE-Web][Sec]
### Story 8.1 — Runtime & sandbox [BE][Sec]
- [ ] JS isolates; CPU/mem/time quotas [BE]
- [ ] Capability tokens; network egress policy [Sec]
- [ ] Billing events for paid apps [BE]

### Story 8.2 — Signed init payload & GraphQL bridge [BE][FE-Web]
- [ ] HMAC signed payload (user/session/device) [BE]
- [ ] GraphQL edge for webviews/postMessage [FE-Web]
- [ ] Sample mini‑app template (Next.js) [FE-Web]

### Story 8.3 — Dev portal [FE-Web]
- [ ] App submission & webhook tester [FE-Web]
- [ ] Analytics (MAU, errors, perf) [FE-Web][Data]
- [ ] Payments SDK docs & examples [FE-Web]

---

# EPIC 9 — AI Lead Engine & RAG [AI][BE][Data]
### Story 9.1 — Event contracts & Kafka topics [BE][AI]
- [ ] Define `login`, `chatIntent`, `view`, `liveEngage` events [BE]
- [ ] PII minimization & schemas in proto [AI][Sec]
- [ ] Producer libraries in SDKs [FE-Web][Android][iOS]

### Story 9.2 — ai‑orchestrator (Go) [AI][BE]
- [ ] Model router + circuit breakers [AI][BE]
- [ ] Feature flags per country; shadow mode [AI]
- [ ] Telemetry (latency, failure, QoS) [Data]

### Story 9.3 — Feature store & scoring [Data][AI]
- [ ] ClickHouse features & LightGBM v1 [Data][AI]
- [ ] Thresholds & AB buckets [AI][Prod]
- [ ] Lead write path + dedupe [BE]

### Story 9.4 — RAG over catalogs/FAQs [AI][BE]
- [ ] Ingest pipeline; embeddings (pgvector/Weaviate) [AI]
- [ ] Grounding templates & eval set [AI]
- [ ] Merchant bot notify + feedback loop [BE]

---

# EPIC 10 — Mobile Apps (Android/iOS) [Android][iOS]
### Story 10.1 — App skeletons & navigation
- [ ] Android single‑activity, 3 fragments (Chat/Store/Social) [Android]
- [ ] iOS SwiftUI TabView + background tasks [iOS]
- [ ] Shared design primitives & theming [Android][iOS]

### Story 10.2 — Auth screens & deep links
- [ ] OTP/magic link screens & retries [Android][iOS]
- [ ] Autofill (Android), Universal Links (iOS) [Android][iOS]
- [ ] Secure session storage [Android][iOS]

### Story 10.3 — Chat UI v1
- [ ] Virtualized list & bubble system [Android][iOS]
- [ ] Attachments picker + upload state [Android][iOS]
- [ ] Local outbox/resend; error toasts [Android][iOS]

### Story 10.4 — Store & Social v1
- [ ] Product feed + cart sheet [Android][iOS]
- [ ] Social feed + follows [Android][iOS]
- [ ] Empty/error states + pull‑to‑refresh [Android][iOS]

---

# EPIC 11 — Web (Next.js) [FE-Web]
### Story 11.1 — PWA & offline shell
- [ ] SW caching for chat shell & background sync [FE-Web]
- [ ] Install prompts & icon set [FE-Web]
- [ ] Error boundaries & retry logic [FE-Web]

### Story 11.2 — GraphQL edge (BFF)
- [ ] Server actions; auth context propagation [FE-Web]
- [ ] Tracing to OpenTelemetry [FE-Web]
- [ ] Rate limits & caching policies [FE-Web]

### Story 11.3 — Merchant console
- [ ] Catalog, orders, payouts, analytics [FE-Web]
- [ ] RBAC with roles/permissions [FE-Web][BE]
- [ ] Bulk import & validation flows [FE-Web]

---

# EPIC 12 — Search & Discovery [BE][Infra]
### Story 12.1 — OpenSearch cluster & ILM
- [ ] Index templates; ILM policies; snapshots [Infra]
- [ ] Service auth & tenancy model [Infra]
- [ ] Cost guardrails & SLOs [Infra]

### Story 12.2 — Indexers
- [ ] Messages/merchants/products pipelines [BE]
- [ ] Backfill jobs + progress tracking [BE]
- [ ] Search API & query hints [BE]

---

# EPIC 13 — Data & Analytics [Data]
### Story 13.1 — Event taxonomy
- [ ] Tracking plan & schemas [Data][Prod]
- [ ] Privacy budgets & retention [Data][Sec]
- [ ] QA of event completeness [Data]

### Story 13.2 — Dashboards
- [ ] Activation/lead funnels [Data]
- [ ] Payment success & failure trees [Data]
- [ ] SLO & error budget burn [Data]

---

# EPIC 14 — Security & Compliance [Sec][BE][Infra]
### Story 14.1 — Threat model v1
- [ ] STRIDE per service; prioritized mitigations [Sec]
- [ ] Secrets posture & rotation policy [Sec]
- [ ] Break‑glass & incident drills [Sec][Infra]

### Story 14.2 — PDPA/AML controls
- [ ] Data map & DSR process [Sec]
- [ ] Retention schedules by country [Sec]
- [ ] Sanctions screening hooks [Sec][BE]

---

# EPIC 15 — QA & Reliability [QA][Infra]
### Story 15.1 — Deterministic data seeding
- [ ] Fixtures for users/merchants/products/wallets [QA][BE]
- [ ] Reset scripts & synthetic datasets [QA]

### Story 15.2 — Playwright E2E (web)
- [ ] Signup → onboarding → chat send [QA][FE-Web]
- [ ] Chat‑to‑cart → payment → order tracking [QA][FE-Web]
- [ ] Visual regression & trace retention [QA]

### Story 15.3 — Load & chaos
- [ ] k6 fan‑out baseline (5k msg/sec) [QA][Infra]
- [ ] Broker/DB/CDN fault injection [QA][Infra]
- [ ] SLO validation & reports [QA]

---

# EPIC 16 — Localization & Content [Prod]
### Story 16.1 — L10N kits (TH/EN first)
- [ ] ICU messages; Buddhist calendar [Prod]
- [ ] Number/date formats; glossary [Prod]
- [ ] RTL readiness (future) [Prod]

---

# EPIC 17 — Partnerships & Legal [Prod][Sec]
### Story 17.1 — Payments MOUs (TH)
- [ ] PromptPay/banks agreements signed [Prod]
- [ ] Webhook whitelists & security review [Sec]
- [ ] Legal terms & pricing compliance [Prod]

### Story 17.2 — Telco/SMS Sender IDs
- [ ] Register sender IDs & verify throughput [Prod]
- [ ] SLA docs & escalation contacts [Prod]

---

# EPIC 18 — Growth & Community [Prod][BE]
### Story 18.1 — Campus soft launch
- [ ] Ambassador program & referral mechanics [Prod]
- [ ] LINE migration toolkit [Prod][BE]
- [ ] Activation incentives & tracking [Prod][Data]

### Story 18.2 — Merchant referral program
- [ ] Tier rules & fraud protection [Prod][Sec]
- [ ] Tracking & payouts [BE]
- [ ] Dashboard & comms templates [FE-Web]

---

# EPIC 19 — Support & Ops [Infra][Prod]
### Story 19.1 — Runbooks & on‑call
- [ ] Incident sev levels & playbooks [Infra]
- [ ] Rollback & data‑recovery steps [Infra]
- [ ] Post‑mortem template & cadence [Prod]

### Story 19.2 — Help Center v1
- [ ] Articles for login, payments, orders, privacy [Prod]
- [ ] Search & feedback widgets [FE-Web]
- [ ] Localization TH/EN [Prod]

---

## Initial Sprint Suggestions (TH Launch)
- Sprint 1: AuthSvc core, OTP/Magic Link, Android/iOS auth screens, Web magic link
- Sprint 2: Messaging core schema + outbox, Push gateway, Observability baseline
- Sprint 3: Media uploads, PromptPay connector, Merchant onboarding
- Sprint 4: Chat‑to‑cart, Orders + Webhooks, Logistics adapters (TH)
- Sprint 5: AI event contracts + ai‑orchestrator (shadow), Feature store v1
- Sprint 6: PWA offline shell, Playwright E2E happy path, Threat model v1

---

> **Notes:** Ensure every PR links to its metric/trace panel; ship behind country flags; enforce idempotency on all payment/order/media paths.
