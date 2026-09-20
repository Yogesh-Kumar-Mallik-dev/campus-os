# 01 — Architecture Constitution

**Governing Document:** [Campus_OS_Product_Bible_v1.2.docx](file:///home/yogesh/Downloads/Campus_OS_Product_Bible_v1.2.docx)  
**Status:** Canonical & Locked

---

## 1. Core Principles

1. **Modular Monolith First:** Domain boundaries must be enforced within a modular monolith prior to any distributed service extraction.
2. **Strict Domain Data Ownership:** Each major domain owns its persistence schema. Domains never directly mutate records owned by another domain.
3. **Persistence Boundary:** The Go logical backend **never** accesses PostgreSQL or Prisma directly. It accesses persistence exclusively through canonical database handles implemented in TypeScript.
4. **Business-Level Database Handles:** Canonical database handles expose domain operations with business semantics, not unrestricted generic CRUD.
5. **No Normal Hard Deletes:** Deletion in ordinary application workflows is soft deletion / retirement. Hard deletion is restricted to the Developer break-glass path.
6. **Immutable Audit Records:** Audit logs answer *what happened* and cannot be altered, even by system administrators or domain owners.
7. **Separate Decision Records:** Decision records answer *how a policy or automated decision was reached* (recording rule versions, inputs, and evaluated criteria) separate from transactional audit traces.

---

## 2. Locked Decisions Register

| Decision | Status | Description |
| :--- | :--- | :--- |
| **Permissions assigned to roles, not users** | LOCKED | Users hold roles; roles hold permissions. No ad-hoc direct user permissions. |
| **Multiple roles per person** | LOCKED | A person may simultaneously hold multiple roles (e.g. Teacher + Warden). |
| **Permissions additive across active roles** | LOCKED | Authority aggregates across all currently active roles. |
| **Scoped assignments** | LOCKED | Roles are constrained by predefined scope dimensions. |
| **Multiple independent assignments per role** | LOCKED | A person can have distinct independent scopes for the same role. |
| **Scope combinations use AND semantics** | LOCKED | Dimensions within an assignment combine via AND (e.g. Subject=DAA AND Section=B). |
| **Manual delisting of assignments** | LOCKED | Assignments do not expire silently; they require an authorized delisting workflow. |
| **Workflow-governed authority changes** | LOCKED | Role assignment, delegation, and removal require authorized workflows. |
| **Sequential or parallel approval chains** | LOCKED | Supported workflows may define linear or concurrent approval steps. |
| **Predefined, non-combinable approval policies** | LOCKED | Approvals use standard system presets (Unanimous, Majority, Threshold, Decisive). |
| **Submitted decisions are immutable** | LOCKED | Once submitted, approval or rejection decisions cannot be retracted or modified. |
| **Rejection comments mandatory; approval optional** | LOCKED | Any rejection must provide a documented rationale. |
| **Information requests travel backward only** | LOCKED | Clarifications travel upstream along existing participants in the approval chain. |
| **One-time return/veto while criteria unmet** | LOCKED | Approvers can return a request to an earlier stage once before final criteria are satisfied. |
| **Versioned workflow definitions** | LOCKED | Active in-flight requests retain the workflow version with which they began. |
| **Super Admin exactly one active seat** | LOCKED | There is strictly one institutional Super Admin. |
| **Super Admin succession is atomic** | LOCKED | Outgoing holder is retired and incoming holder is activated in a single atomic transaction. |
| **Developer-only break-glass** | LOCKED | Emergency break-glass access is exclusively restricted to Developer authority. |
| **Developer explicitly terminates break-glass** | LOCKED | Break-glass sessions do not timeout automatically; they require explicit closure. |
| **Post-break-glass freeze & read-only access** | LOCKED | Records modified during break-glass are automatically frozen until incident review. |
| **Developer post-incident review mandatory** | LOCKED | Every emergency intervention requires a formal documented post-incident review. |
| **Master / Operational / Reference separation** | LOCKED | Data is formally categorized with corresponding governance rules. |
| **Canonical corrected value replaces old value** | LOCKED | Historical values remain preserved in audit history, while current value is canonical. |
| **PostgreSQL with domain schemas** | LOCKED | Dedicated schemas (`auth_schema`, `academic_schema`, etc.) enforce physical separation. |
| **Prisma + TypeScript database layer** | LOCKED | Persistence layer implemented in TypeScript wrapping Prisma. |
| **Logical backend accesses DB via canonical handles** | LOCKED | Go backend communicates via strictly typed gRPC handles over Unix Domain Socket. |
| **Automated decisions require exception workflow** | LOCKED | System-calculated rules cannot be directly overridden; exceptions use audited workflows. |

---

## 3. Architectural Guardrails

1. **Do not bypass the canonical database layer.**
2. **Do not let domains directly mutate another domain's records.**
3. **Do not expose unrestricted generic CRUD to the logical backend.**
4. **Do not grant permissions directly to users.**
5. **Do not treat a scope as a new permission.**
6. **Do not silently expire assignments.**
7. **Do not silently resolve unknown institutional data discrepancies.**
8. **Do not mutate submitted approval decisions.**
9. **Do not change workflow rules for in-progress requests.**
10. **Do not invent institutional authority hierarchy where it is unverified.**
11. **Do not expose hard deletion through ordinary application APIs.**
12. **Do not allow normal users or Super Admin to invoke Developer break-glass.**
13. **Do not allow AI agents to silently modify locked architectural decisions.**
