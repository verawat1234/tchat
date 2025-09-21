# Tchat Monorepo

This repository hosts the foundational stacks for the Tchat project across backend (Go), web (Next.js), Android (Kotlin), and iOS (SwiftUI).

## Structure

- `apps/web` — Next.js application targeting the web client.
- `apps/mobile/android` — Kotlin Android application scaffold.
- `apps/mobile/ios` — SwiftUI iOS application scaffold.
- `backend/auth` — Go backend service skeleton (auth service baseline).
- `libs` — Shared libraries (future home for cross-cutting packages).

Each stack is intentionally minimal and ready for local bootstrapping, CI wiring, and further feature development.
