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
6. Initialized `pnpm-workspace.yaml` and root `package.json` with `@bufbuild/buf`, `typescript`, and approved package builds.
7. Initialized multi-module `go.work` linking `backend/` and `apps/client/` (Fyne).
8. Upgraded Prisma and `@prisma/client` to Prisma 8 (`8.1.0-dev.7`).
9. Implemented actual core scripts in `scripts/` (`scripts/*.sh` and `scripts/*.ps1`) for `dev`, `build`, `check`, `test`, `deps`, `proto_gen`, `envi`, `uenvi`, and `flush_db` per engineering standards.
10. Added Go gRPC & Protobuf dependencies (`google.golang.org/grpc v1.84.0`).
11. Generated Protobuf stubs into `backend/pkg/proto/` and `packages/ts-proto/src/` with `buf generate`.
12. Verified all lifecycle commands (`./script.sh check`, `./script.sh test`, `./script.sh build`, `./script.sh deps`) pass cleanly.
13. Converted single `schema.prisma` into multi-file directory target `db-layer/prisma/schema/*.prisma` with dedicated domain models: `base.prisma`, `auth.prisma`, `academic.prisma`, `student.prisma`, `hostel.prisma`, `finance.prisma`, `document.prisma`, and `audit.prisma`.
14. Validated and generated typed Prisma 8 Client from multi-file schemas (`prisma generate --schema=prisma/schema`).

---

## 3. Immediate Next Steps
1. Expand Protobuf service contracts in `/proto/campus/v1/` for academic, auth, and student domain operations.
2. Implement gRPC server handlers and Prisma wrappers in `db-layer/src/handles/`.
3. Implement `TxToken` interactive session manager in `db-layer/src/tx/`.
4. Implement Go gRPC client stubs and domain service interfaces in `backend/internal/domains/`.

---

## 4. Key Architectural Invariants
* **Go Backend:** Pure domain logic, workflows, authorization evaluation. Never imports SQL drivers or Prisma.
* **TypeScript DB Layer:** Implements canonical gRPC handles over UDS exposing domain operations (no generic CRUD).
* **Multi-Dimensional RBAC:** Strict AND semantics for scope dimensions (`Department`, `Course`, `Semester`, `Section`, `Subject`, `Lab`, `HostelBlockFloor`).
* **Strict Single-Change Commits:** One feature/fix per atomic commit, signed, conventional format.
* **Hand-in-Hand Tests:** Co-located tests accompany every domain logic change.
