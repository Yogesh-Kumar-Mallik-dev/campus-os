# Contributing Guidelines — Campus OS

Thank you for your interest in contributing to **Campus OS**! We maintain rigorous standards for code quality, documentation precision, architectural integrity, and testing.

---

## 1. Core Development Invariants

1. **Strict Single-Change Policy:** Each Pull Request (PR) and commit must address exactly **one logical change**—one feature, one bug fix, or one refactor. Do not bundle multiple domains into a single commit.
2. **Persistence Boundary:** Never import Prisma or database drivers into `/backend` (Go). All database operations must route through canonical gRPC handles over Unix Domain Socket (`/var/run/campus-os/db.sock`) in `/db-layer` (TypeScript).
3. **100% Test Coverage for Domain Logic:** All business logic, state machines, and API handlers must include co-located tests.
4. **Unified Script Entrypoint:** All lifecycle workflows are accessible via `./script.sh` (Unix) and `.\script.ps1` (Windows PowerShell).

---

## 2. Getting Started & Local Verification

```bash
# 1. Install dependencies across all monorepo components
./script.sh deps

# 2. Generate Protobuf stubs for Go and TypeScript
./script.sh proto:gen

# 3. Run typechecking & static analysis
./script.sh check

# 4. Run automated test suites
./script.sh test

# 5. Start local development environment via Docker Compose
./script.sh dev
```

---

## 3. Pull Request & Commit Standards

1. **Signed Commits:** All commits must be cryptographically signed (`git commit -S -m "..."`).
2. **Conventional Commits:** Messages must follow the Conventional Commits format:
   * `feat(academic): add course semester hosting handler`
   * `fix(auth): correct AND-semantics evaluation for scoped assignments`
   * `test(hostel): add outpass approval transition unit tests`
   * `docs(standards): update documentation index`
3. **Linear History:** Rebase your branch on `main` before submitting your PR (`git pull --rebase origin main`).
4. **Context Tracking:** Update [`.agent/current_context.md`](file:///home/yogesh/campus_os/.agent/current_context.md) with every functional iteration.
