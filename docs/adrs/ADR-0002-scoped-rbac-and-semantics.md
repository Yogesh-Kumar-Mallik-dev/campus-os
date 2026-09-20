# ADR-0002: Multi-Dimensional Scoped RBAC with AND Semantics

* **Status:** Accepted / Locked
* **Date:** 2026-09-20
* **Deciders:** Developer, Yogesh

---

## Context

Permissions must not be granted directly to individuals. Furthermore, individuals often hold authority over specific segments (e.g. one section or one subject) without possessing department- or campus-wide authority.

## Decision

We implement a **Scoped RBAC** model:
`Person -> Role Assignment [Scope Dimensions] -> Role -> Permissions`.

* Scopes are predefined dimensions: `Department`, `Course`, `Semester`, `Section`, `Subject`, `Lab`, `HostelBlockFloor`.
* When an assignment contains multiple dimensions, they combine using strict **AND semantics**.
* A person may hold multiple distinct assignments for the same or different roles.
* Authority is additive across active assignments.
* Assignments are delisted manually via workflow; no silent automatic expiration.

## Consequences

* **Positive:** Eliminates permission explosion (no need for `TEACHER_DAA_SECTION_B` roles). Prevents privilege creep.
* **Negative:** Requires multi-attribute evaluation logic on every sensitive authorization check.
