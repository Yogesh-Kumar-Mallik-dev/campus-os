# Campus OS — Current Context Tracker

**Last Updated:** 2026-09-21  
**Current Phase:** Native Client Component System & Institutional Identity Polish  
**Active Branch:** `main`

---

## 1. Status Overview
* **Section 23 Architectural Baseline Complete:** Formalized `docs/01` through `docs/09` and ADRs `ADR-0001` through `ADR-0006`.
* **Prisma 8 Multi-Schema Persistence Layer:** Implemented `db-layer/prisma/schema/*.prisma` with comprehensive models for auth, student, academic, and audit domains.
* **Go Logical Backend Core:** Implemented configuration, error envelopes (RFC 7807), collection query semantics, and UDS gRPC persistence client.
* **Go Fyne Native Client Modernized:** Focused strictly on Account Activation & Scholar Onboarding flow, recreating 10 shadcn UI primitives in pure Go, integrating Lucide vector icons, and embedding the official BBDIT institutional logo.
* **Continuous Integration & Quality:** All lifecycle scripts (`./script.sh check`, `./script.sh test`, `./script.sh build`) passing 100% across backend, client, and DB layer.

---

## 2. Recent Actions Completed
1. Renamed git branch from `master` to `main`.
2. Cloned and audited `engineering-standards.git`, scaffolding `.agent/` and root configs.
3. Created multi-module `go.work` linking `backend/` and `apps/client/`.
4. Upgraded to Prisma 8 (`8.1.0-dev.7`) and established multi-file domain schemas.
5. Implemented gRPC over Unix Domain Socket boundary (`/tmp/campus-db.sock`) with Protobuf generation.
6. Formalized Canonical Identity Conventions, Alt-Name modeling, and Genesis Bootstrapping in `ADR-0005`.
7. Formalized Affiliating University (AKTU) Result Ingestion & Revision Ledger in `ADR-0006`.
8. Implemented full-stack Onboarding & Authentication flow across Protobuf, TypeScript DB handlers, Go backend, and Fyne native client.
9. Standardized test naming to `*.test.ts` and `*_test.go` and added gRPC client test suites.
10. Added cross-platform PowerShell dev runners (`dev.ps1`, `dev.client.ps1`, `dev.backend.ps1`, `dev.db.ps1`, `dev.mock.ps1`).
11. Extracted and mapped institutional OKLCH color tokens from abandoned project's `+layout.css` (`#10161C`, `#18202A`, `#2E3844`, `#F45A51`, `#F2F6F8`).
12. Stripped extraneous unrequested navigation tabs from native client, focusing 100% on the 5-step Account Activation wizard.
13. Recreated 10 native shadcn-style UI primitives in Go Fyne with full unit test coverage:
    - Core: `NewShadcnCard`, `NewBadge`, `NewShadcnButton`, `NewShadcnInput`, `NewFormField`, `NewShadcnAlert`.
    - Feedback & Overlay: `NewSwitch`, `ShowToast`, `ShowAlertDialog`, `NewEmptyState`, `NewShadcnSeparator`.
14. Fixed responsive layout and text wrapping defects:
    - Enforced `TextWrapWord` across all labels in composite cards.
    - Resolved `container.NewCenter` 76px column squish in `NewEmptyState`, aligning to shadcn's full-width dropzone specification.
    - Removed redundant outer card nesting, establishing a clean canvas visual hierarchy.
15. Created pure vector Lucide icon suite in `icons.go` (`LucideScan`, `LucideQrCode`, `LucideSim`, `LucideSmartphone`, `LucideShieldCheck`, `LucideUserCheck`, `LucideFileCheck`, `LucideLock`, `LucideClock`, `LucideFingerprint`, `LucideBadgeCheck`, `LucideSparkles`, `LucideFlag`), removing all generic Fyne theme icons.
16. Imported 34 official BBDIT institutional assets from abandoned project into `apps/client/assets/` and embedded `bbdit-logo-transparent.png` into the native client top bar (`NewTopBar`).
17. Updated repository `README.md` and documentation set.

---

## 3. Immediate Next Steps
1. Gather user feedback on the updated client appearance and responsiveness via `./scripts/dev.client.sh`.
2. Expand Protobuf service definitions and persistence handlers for subsequent domain workflows when prioritized.
3. Maintain 100% test coverage and zero-broken-window policy across all commits.

---

## 4. Key Architectural Invariants
* **Go Backend:** Pure domain logic, workflows, authorization evaluation. Never imports SQL drivers or Prisma.
* **TypeScript DB Layer:** Implements canonical gRPC handles over UDS exposing domain operations (no generic CRUD).
* **Native Client UI:** Zero emojis, zero generic raster icons. Pure Lucide vector SVGs and embedded official BBDIT assets.
* **Multi-Dimensional RBAC:** Strict AND semantics for scope dimensions (`Department`, `Course`, `Semester`, `Section`, `Subject`, `Lab`, `HostelBlockFloor`).
* **Atomic Signed Commits:** Single-purpose, atomic, signed commits per task.
