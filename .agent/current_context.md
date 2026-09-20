# Campus OS — Current Context Tracker

**Last Updated:** 2026-09-20  
**Current Phase:** Pre-Implementation Specification & Engineering Standards Scaffolding (Section 23 Compliance)  
**Active Branch:** `main`

---

## 1. Status Overview
* **Product Bible v1.2 Audited:** Foundational Product Bible read and analyzed.
* **Deferred / TBD Decisions Resolved:** All 11 open institutional and technical areas resolved via pair-programming session.
* **Section 23 Documentation Scaffolding Complete:** Initial specification suite (`docs/01` through `docs/08` and `docs/adrs/ADR-0001` through `ADR-0004`) created and committed.
* **Engineering Standards Audited & Adopted:** Universal engineering standards from `engineering-standards.git` incorporated into `.agent/`, root configs, and `docs/09-engineering-standards.md`.

---

## 2. Recent Actions Completed
1. Renamed git branch from `master` to `main`.
2. Cloned and audited `git@github.com:Yogesh-Kumar-Mallik-dev/engineering-standards.git`.
3. Created `.agent/` architecture with tailored `agents.md`, `current_context.md`, and 8 universal rules in `.agent/rules/`.
4. Established repository root configuration files (`.editorconfig`, `.gitattributes`, `.gitignore`, `CONTRIBUTING.md`, `SECURITY.md`, `script.sh`, `script.ps1`).
5. Synthesized `docs/09-engineering-standards.md` into the canonical documentation set.

---

## 3. Immediate Next Steps
1. Configure `buf.yaml` and initial Protobuf definitions in `/proto` for gRPC persistence contracts.
2. Initialize the Go logical backend module (`/backend`) with Go 1.22+ and package structure.
3. Initialize the TypeScript database layer (`/db-layer`) with Prisma and PostgreSQL domain schemas.
4. Set up Docker Compose topology with shared volume for Unix Domain Socket (`/var/run/campus-os/db.sock`).

---

## 4. Key Architectural Invariants
* **Go Backend:** Pure domain logic, workflows, authorization evaluation. Never imports SQL drivers or Prisma.
* **TypeScript DB Layer:** Implements canonical gRPC handles over UDS exposing domain operations (no generic CRUD).
* **Multi-Dimensional RBAC:** Strict AND semantics for scope dimensions (`Department`, `Course`, `Semester`, `Section`, `Subject`, `Lab`, `HostelBlockFloor`).
* **Strict Single-Change Commits:** One feature/fix per atomic commit, signed, conventional format.
* **Hand-in-Hand Tests:** Co-located tests accompany every domain logic change.
