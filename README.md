# Campus OS

**Canonical Institutional ERP++ Platform for Single-Institution Operation**  
*Official Digital Operating System for Babu Banarasi Das Institute of Technology (BBDIT)*

Campus OS is a high-assurance, digitally sovereign campus operating system engineered to eliminate administrative paper friction, enforce cryptographic auditability, preserve institutional governance, and streamline life for scholars, faculty, and leadership.

---

## 1. High-Level Technology Stack

| Layer | Technology | Role |
| :--- | :--- | :--- |
| **Presentation (Desktop/Mobile)** | **Go Fyne v2** | Cross-platform native client for Linux, macOS, Windows, Android, and iOS featuring recreated shadcn components and Lucide vector icons. |
| **Presentation (Web)** | **SvelteKit** | Progressive web application client for administrative dashboards and self-service portals. |
| **Logical Backend** | **Go (Modular Monolith)** | Pure domain business logic, workflow state machines, scoped RBAC evaluation, and orchestration. |
| **Persistence Boundary** | **TypeScript + Prisma 8** | Canonical database handle microservice exposing strictly-typed business operations. |
| **Inter-Process Comm (IPC)** | **gRPC over Unix Domain Socket** | High-performance, zero-network-exposure typed boundary between Go and TypeScript (`/tmp/campus-db.sock`). |
| **System of Record** | **PostgreSQL (Multi-Schema)** | Relational database divided into `auth`, `academic`, `student`, `hostel`, `finance`, and `audit` schemas. |
| **Observability** | **OpenTelemetry + Prometheus** | Structured JSON logging with trace context, collection metrics, and audit event streams. |

---

## 2. Native Client & Design System (`apps/client`)

The native cross-platform application is implemented in Go Fyne, adhering to institutional design standards and zero-emoji, pure vector aesthetics.

### Recreated Native shadcn-Style Primitives
All UI components mirror the component architecture of `shadcn-svelte` natively in Go:
* **`NewShadcnCard`**: Compound card widget with structured header, pill badge, body content, and footer actions.
* **`NewShadcnButton`**: Modern button with variants (`ButtonDefault`, `ButtonSecondary`, `ButtonOutline`, `ButtonDestructive`, `ButtonGhost`), responsive sizing, and vector SVG icon support.
* **`NewBadge`**: Visual status tag with variant coloring (`BadgeDefault`, `BadgeSecondary`, `BadgeSuccess`, `BadgeWarning`, `BadgeDestructive`, `BadgeOutline`) and pill/rounded geometry.
* **`NewShadcnInput` & `NewFormField`**: Clean entry fields with label typography, placeholder tokens, and validation error messages.
* **`NewShadcnAlert`**: Themed notification and warning banners matching institutional severity tokens.
* **`NewSwitch`**: Accessible animated toggle switch widget for feature activation (e.g. biometric authentication).
* **`ShowToast`**: Non-blocking auto-dismissing floating notification banner overlay.
* **`ShowAlertDialog`**: Modal confirmation and audit-flagging dialog for high-stakes actions.
* **`NewEmptyState`**: Full-width dropzone placeholder with centered media, title, and descriptive helper text.
* **`NewShadcnSeparator`**: Subtle 1px structural dividing lines utilizing theme border tokens.

### Official Institutional Assets & Lucide Icons
* **Official BBDIT Institutional Logo**: Embedded directly via `//go:embed` from `bbdit-logo-transparent.png`, encapsulated within a high-contrast pill container in `NewTopBar`.
* **Lucide Vector Icon Pack**: Pure vector SVGs with 24×24 geometry, 2px stroke width, and rounded terminals (`LucideScan`, `LucideQrCode`, `LucideSim`, `LucideSmartphone`, `LucideShieldCheck`, `LucideUserCheck`, `LucideFileCheck`, `LucideKeyRound`, `LucideLock`, `LucideClock`, `LucideFingerprint`, `LucideBadgeCheck`, `LucideSparkles`, `LucideFlag`).
* **Design Tokens**: Strict color mapping from `+layout.css`:
  * Dark Canvas: `oklch(0.16 0.018 242)` $\rightarrow$ `#10161C`
  * Dark Card: `oklch(0.205 0.02 240)` $\rightarrow$ `#18202A`
  * Dark Border: `oklch(0.72 0.018 225 / 16%)` $\rightarrow$ `#2E3844`
  * Dark Primary: `oklch(0.67 0.19 27)` $\rightarrow$ `#F45A51`
  * Text Contrast: `oklch(0.97 0.004 225)` $\rightarrow$ `#F2F6F8`

---

## 3. Account Activation & Onboarding Flow

The client delivers a 5-step zero-knowledge onboarding flow:
1. **Optical QR Capture**: Scans sealed admission QR voucher or accepts fallback token verification.
2. **Hardware SIM Telephony Binding**: Validates carrier SIM presence and processes SMS OTP challenge.
3. **Admission Dossier Verification**: Displays academic vs. legal full name, lateral entry semester, and allows discrepancy reporting to the Registrar.
4. **Credential Provisioning**: Sets account access password under the 72-hour orientation grace period policy and configures biometric keyrings.
5. **Digital Campus Pass Generation**: Issues active student gate credentials with NFC/QR checkpoint authorization.

---

## 4. Repository Structure

```
campus_os/
├── apps/
│   └── client/                 # Go Fyne native desktop & mobile client
│       ├── assets/             # Official BBDIT image assets & photography
│       └── internal/ui/        # Recreated shadcn primitives, theme, and onboarding wizard
├── backend/                    # Go logical backend (Modular Monolith)
│   ├── cmd/server/             # Backend server & CLI bootstrap entrypoints
│   ├── internal/config/        # Environment and runtime configuration
│   ├── internal/transport/     # HTTP REST & gRPC Unix Domain Socket clients
│   └── pkg/proto/              # Generated Go Protobuf stubs
├── db-layer/                   # TypeScript Prisma 8 persistence boundary
│   ├── prisma/schema/          # Multi-file domain schemas (auth, student, academic, etc.)
│   └── src/handles/            # Domain operation handlers over gRPC UDS
├── docs/                       # Canonical architecture constitution and ADRs
│   ├── 01-architecture-constitution.md
│   ├── 02-product-requirements.md
│   ├── 03-domain-catalog.md
│   ├── 04-authorization-specification.md
│   ├── 05-workflow-specification.md
│   ├── 06-data-governance-specification.md
│   ├── 07-api-and-ipc-contracts.md
│   ├── 08-ai-agent-instructions.md
│   ├── 09-engineering-standards.md
│   └── adrs/                   # Architecture Decision Records (ADR-0001 through ADR-0006)
├── packages/                   # Shared TypeScript packages & proto contracts
├── proto/                      # Canonical Protobuf schema definitions
└── scripts/                    # Unified lifecycle runners (Bash & PowerShell)
```

---

## 5. Development & Lifecycle Commands

The repository provides standardized scripts across Linux/macOS (`script.sh`) and Windows (`script.ps1`):

### Core Workflows
```bash
# Verify static analysis, Go vet, TypeScript typechecking, and Protobuf linting
./script.sh check

# Execute comprehensive test suites across Go backend, native client, and TypeScript DB layer
./script.sh test

# Build production binaries (bin/campus-backend and bin/campus-client)
./script.sh build

# Regenerate Protobuf contracts across Go and TypeScript
./script.sh proto_gen
```

### Dev Runners
```bash
# Launch Go Fyne Native Client
./scripts/dev.client.sh         # On Linux/macOS
./scripts/dev.client.ps1        # On Windows (PowerShell)

# Launch Go Backend
./scripts/dev.backend.sh        # On Linux/macOS
./scripts/dev.backend.ps1       # On Windows (PowerShell)

# Launch TypeScript DB Layer
./scripts/dev.db.sh             # On Linux/macOS
./scripts/dev.db.ps1            # On Windows (PowerShell)
```

---

## 6. Architecture Reference Index

* [Architecture Constitution](docs/01-architecture-constitution.md) — Persistence boundary, zero-exposure database access, and architectural invariants.
* [Product Requirements](docs/02-product-requirements.md) — Product vision, institutional scope, and operational domains.
* [Domain Catalog](docs/03-domain-catalog.md) — Domain schemas, table ownership, and time-series attendance models.
* [Authorization Specification](docs/04-authorization-specification.md) — Multi-dimensional RBAC with strict `AND` semantics, dual-identity tri-name handling, and delegation.
* [Workflow Specification](docs/05-workflow-specification.md) — Workflow state machines and audit vs. decision records.
* [Data Governance Specification](docs/06-data-governance-specification.md) — Master/Operational/Reference data classification and migration pipelines.
* [API & IPC Contracts Specification](docs/07-api-and-ipc-contracts.md) — Go $\leftrightarrow$ TypeScript gRPC boundary over Unix Domain Socket with TxToken transaction sessions.
* [AI Development Constitution](docs/08-ai-agent-instructions.md) — Guardrails and invariants for AI pair programming.
* [Engineering Standards](docs/09-engineering-standards.md) — Block construction (`BLOCK_<DOMAIN>_<ACTION>_<ID>`), testing rigor, RFC 7807 error envelopes, and responsive layout standards.
* [Architecture Decision Records (ADRs)](docs/adrs/) — Chronological record of locked architectural decisions:
  * [ADR-0001: Go-TypeScript gRPC over UDS Boundary](docs/adrs/ADR-0001-go-ts-grpc-uds-boundary.md)
  * [ADR-0002: Scoped Multi-Dimensional RBAC](docs/adrs/ADR-0002-scoped-rbac-and-semantics.md)
  * [ADR-0003: TxToken Interactive Transaction Sessions](docs/adrs/ADR-0003-txtoken-transaction-sessions.md)
  * [ADR-0004: Academic Course Semester Model](docs/adrs/ADR-0004-academic-course-semester-model.md)
  * [ADR-0005: Canonical Identity Conventions & Onboarding Architecture](docs/adrs/ADR-0005-canonical-identity-and-onboarding.md)
  * [ADR-0006: Affiliating University (AKTU) Result Ingestion & Revision Ledger](docs/adrs/ADR-0006-affiliating-university-result-ingestion-and-marks-ledger.md)
