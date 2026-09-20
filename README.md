# Campus OS

Canonical Institutional ERP++ Platform for Single-Institution Operation.

Campus OS is a digital operating system engineered to make campus operations transparent, eliminate paper friction, preserve accountability, and streamline institutional life for students, staff, and leadership.

---

## 1. Documentation Index (Section 23 Architecture Set)

Before code generation, the system architecture and domain models must be internally consistent. Refer to the canonical repository specifications:

* [Architecture Constitution](file:///home/yogesh/campus_os/docs/01-architecture-constitution.md) — Core principles, persistence boundaries, guardrails, and locked decisions.
* [Product Requirements](file:///home/yogesh/campus_os/docs/02-product-requirements.md) — Product vision, institutional boundaries, and core functional scopes.
* [Domain Catalog](file:///home/yogesh/campus_os/docs/03-domain-catalog.md) — Domain schemas, table ownership, and persistence boundaries.
* [Authorization Specification](file:///home/yogesh/campus_os/docs/04-authorization-specification.md) — Scoped RBAC, AND semantics, privacy classification, and role delegation.
* [Workflow & Approval Specification](file:///home/yogesh/campus_os/docs/05-workflow-specification.md) — Workflow state machines, predefined approval policies, audit vs. decision records.
* [Data Governance Specification](file:///home/yogesh/campus_os/docs/06-data-governance-specification.md) — Master/Operational/Reference data classification, migration, and discrepancies.
* [API & IPC Contracts Specification](file:///home/yogesh/campus_os/docs/07-api-and-ipc-contracts.md) — Go $\leftrightarrow$ TypeScript gRPC boundary over Unix Domain Socket, Protobuf schemas, and TxToken transactions.
* [AI Development Constitution](file:///home/yogesh/campus_os/docs/08-ai-agent-instructions.md) — Mandatory governing rules and guardrails for AI coding assistants.
* [Engineering Standards & AI Anti-Pattern Prevention](file:///home/yogesh/campus_os/docs/09-engineering-standards.md) — Authoritative engineering standards, block construction (`BLOCK_<DOMAIN>_<ACTION>_<ID>`), testing rigor, RFC 7807, and mobile-first rules.
* [Architecture Decision Records (ADRs)](file:///home/yogesh/campus_os/docs/adrs) — Chronological index of locked architectural decisions.

---

## 2. High-Level Technology Stack

| Layer | Technology | Role |
| :--- | :--- | :--- |
| **Presentation (Web)** | SvelteKit | Responsive web client for students, faculty, and administrators. |
| **Presentation (Desktop/Mobile)** | Fyne (Go) | Cross-platform client for desktop (Linux/macOS/Windows) and mobile (Android/iOS). |
| **Logical Backend** | Go (Modular Monolith) | Domain logic, workflow engine, authorization evaluator, orchestration. |
| **Persistence Boundary** | TypeScript + Prisma | Canonical database handles exposing business-level operations. |
| **Inter-Process Comm (IPC)** | gRPC over Unix Domain Socket | High-performance, strictly typed boundary between Go and TypeScript. |
| **System of Record** | PostgreSQL | Multi-schema database (`auth`, `academic`, `student`, `hostel`, `finance`, `audit`). |
| **Observability** | OpenTelemetry + Prometheus + Grafana + Loki | End-to-end tracing, metrics, dashboards, and structured JSON log collection. |
| **Deployment** | Docker Compose on Ubuntu | Containerized deployment with shared volume for Unix Domain Socket. |
