# Campus OS — Current Context Tracker

**Last Updated:** 2026-09-22  
**Current Phase:** Production-Grade Responsive Architecture & Dev Environment Hardening  
**Active Branch:** `main`

---

## 1. Status Overview
* **Section 23 Architectural Baseline Complete:** Formalized `docs/01` through `docs/10` and ADRs `ADR-0001` through `ADR-0007`.
* **Prisma 8 Multi-Schema Persistence Layer:** Implemented `db-layer/prisma/schema/*.prisma` with comprehensive models for auth, student, academic, and audit domains.
* **Go Logical Backend Core:** Implemented configuration, error envelopes (RFC 7807), collection query semantics, UDS gRPC persistence client, and sanitized gRPC transport error mapping (`mapGRPCError`).
* **Go Fyne Native Client Modernized:** Pure vector design system (shadcn primitives, Lucide icons, BBDIT branding) with full mobile-first responsive layout engine, accessible modal overlays (`ShowModal`), left-docked mobile navigation drawer, and universal theme switcher.
* **Continuous Integration & Quality:** All lifecycle scripts (`./scripts/check.sh`, `./scripts/test.sh`, `./scripts/build.sh`) passing 100% across backend, client, and DB layer.

---

## 2. Recent Actions Completed
1. **Responsive Mobile-First Layout Engine (`apps/client/internal/ui/`):**
   - Implemented dynamic adaptive breakpoint evaluation (`<640` Mobile, `640..1024` Tablet, `>1024` Desktop).
   - Designed and integrated `AdaptiveGridLayout` and `FlowLayout` for dynamic multi-column wrapping without clipping.
   - De-cluttered Institutional TopBar: removed redundant student role badge, converted theme button to an interactive toggle switch, docked action buttons top-right, and added symmetric horizontal padding.
   - Docked mobile navigation drawer to left screen edge with dark semi-transparent backdrop overlay and ESC key dismissal.
   - Upgraded all modals (`ShowModal`, camera permission, camera scanner, overflow actions) with keyboard `Escape` dismissal, outside-click backdrop dismissal, and responsive max-width clamping.
   - Scaled BBDIT institutional logo (`150x28`) and added symmetric gutters across topbar, mobile drawer, sidebar, and workspace to eliminate element clipping.
2. **IPC Boundary Error Sanitization & RFC 7807 Protection (`backend/internal/transport/http/`):**
   - Implemented `mapGRPCError` to intercept low-level gRPC transport connection faults and `codes.Unavailable`.
   - Transformed Unix socket dial errors into clean HTTP 503 (`SERVICE_UNAVAILABLE`: "Database Persistence Offline") responses, strictly preventing Unix domain socket paths (`/tmp/*.sock`) from leaking to clients or end users.
   - Added unit test `TestAuthHTTPHandler_UnavailablePersistence`.
   - Added client-side detection (`IsPersistenceUnavailable`) and error message sanitization (`sanitizeUIError`).
3. **Local Dev Lifecycle & Port Conflict Hardening (`scripts/`):**
   - Added automated pre-flight port reclamation in `scripts/dev.sh` to cleanly terminate any lingering zombie processes occupying `:8080`.
   - Added synchronous PID crash detection in `scripts/dev.sh` (`kill -0 "$BACKEND_PID"`) preventing false-positive readiness reporting.
   - Replaced subshell execution with `exec go run` in `scripts/dev.backend.sh` and `scripts/dev.client.sh` for direct POSIX signal propagation.
   - Ensured socket file unlinking and port release in dev `cleanup()` trap.
4. **Documentation Alignment:**
   - Updated `README.md`, `docs/07-api-and-ipc-contracts.md`, `docs/09-engineering-standards.md`, and `docs/10-super-admin-and-dashboard-architecture.md` to reflect all responsive architecture, accessible modal standards, and IPC error containment guarantees.
5. **Fully Adaptive Responsive Layout Engine (`apps/client/internal/layout/`):**
   - Implemented production-grade, constraint-based responsive layout engine for Go + Fyne (`viewport.go`, `breakpoint.go`, `constraints.go`, `spacing.go`, `container.go`, `grid.go`, `flow.go`, `row.go`, `column.go`, `stack.go`, `spacer.go`, `responsive.go`, `shell.go`, `engine.go`, `utilities.go`).
   - Pure geometry and constraint-driven architecture with zero device-class checks (`SizeClass`, `HeightClass`, `Orientation`, `AspectScale`, `Columns`, `FluidWidth`).
   - Fluid max-width container with automatic horizontal centering on ultrawide displays and adaptive padding scaling.
   - Mathematical dynamic grid calculating columns from available width without hardcoded device breakpoints.
   - Flex Row & Column layouts with proportional weight distribution (`FlexItem`, `Spacer`).
   - State-preserving `Responsive` switcher and `ResponsiveAppShell` (Desktop fixed sidebar, Tablet header drawer, Mobile bottom navigation).
   - Created standalone interactive demonstration app (`apps/client/cmd/demo/main.go`), `--demo` flag, and launcher script (`scripts/demo.sh`).
   - Authored comprehensive developer documentation in `docs/11-responsive-layout-engine.md`.
   - 100% passing unit test suite in `apps/client/internal/layout/layout_test.go` covering all boundaries, resizes, grids, constraints, and awkward aspect-ratio matrices.

---

## 3. Immediate Next Steps
1. Gather user verification on full-stack dev startup (`./scripts/dev.sh`) and responsive client layouts across desktop/mobile resolutions.
2. Expand Protobuf service definitions and persistence handlers for subsequent domain workflows (Academics, Outpass, Marks Ledger) as prioritized.
3. Maintain 100% test coverage and zero-broken-window policy across all commits.

---

## 4. Key Architectural Invariants
* **Go Backend:** Pure domain logic, workflows, authorization evaluation. Never imports SQL drivers or Prisma.
* **TypeScript DB Layer:** Implements canonical gRPC handles over UDS exposing domain operations (no generic CRUD).
* **IPC Isolation:** Internal domain socket paths and transport dial traces never leak through HTTP problem details or UI messages.
* **Native Client UI:** Zero emojis, zero generic raster icons. Pure Lucide vector SVGs, embedded official BBDIT assets, accessible modal overlays (`ShowModal`), and mobile-first adaptive layouts.
* **Multi-Dimensional RBAC:** Strict AND semantics for scope dimensions (`Department`, `Course`, `Semester`, `Section`, `Subject`, `Lab`, `HostelBlockFloor`).
* **Atomic Signed Commits:** Single-purpose, atomic, cryptographically signed commits (`git commit -S`) per task.
