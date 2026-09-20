# Master Agent Instructions — Campus OS

Welcome, AI Agents & Engineers! When working within the **Campus OS** repository, you are strictly bound by the [Architecture Constitution](file:///home/yogesh/campus_os/docs/01-architecture-constitution.md), the [Product Bible v1.2](file:///home/yogesh/Downloads/Campus_OS_Product_Bible_v1.2.docx), and the Universal Engineering Standards.

---

## 1. Context Tracking Protocol
1. **At Task Start:** Always inspect [`.agent/current_context.md`](file:///home/yogesh/campus_os/.agent/current_context.md) to understand where the previous agent or developer left off.
2. **At Task End:** Before concluding your session, update [`.agent/current_context.md`](file:///home/yogesh/campus_os/.agent/current_context.md) recording:
   * Current working status
   * Completed steps in the current session
   * Decisions or deviations made
   * Immediate next steps for subsequent agents.

---

## 2. Architectural Boundaries & Invariants
1. **Persistence Boundary:**
   * The Go Logical Backend (`/backend`) **must never** import database drivers, Prisma, or execute raw SQL.
   * Persistence operations must route through canonical gRPC handles over Unix Domain Socket (`/var/run/campus-os/db.sock`) exposed by the TypeScript database layer (`/db-layer`).
2. **Domain Isolation:**
   * Each domain (`auth`, `academic`, `student`, `hostel`, `finance`, `document`, `audit`) strictly owns its PostgreSQL schema.
   * Domains communicate across boundaries via explicit synchronous contracts or asynchronous domain events.
3. **No Normal Hard Deletes:**
   * All entity deletion in standard workflows is soft deletion (`retired_at`, `retired_by_id`).
   * Hard deletion is exclusively reserved for Developer break-glass procedures.
4. **Scoped RBAC with AND Semantics:**
   * Permissions attach to Roles, never directly to users.
   * Assignments bind users to roles constrained by predefined scope dimensions (`Department`, `Course`, `Semester`, `Section`, `Subject`, `Lab`, `HostelBlockFloor`).

---

## 3. Mandatory Hand-in-Hand Testing
* **Co-located Tests:** Whenever a new feature, domain module, helper, or handler is added or modified, corresponding unit tests **must** be created or updated hand-in-hand in the exact same change.
* **Coverage Target:** Aim for 100% test coverage across core domain logic, state machine transitions, and authorization scope evaluators.
* **Mockable Interfaces:** All transport and external persistence dependencies must be injected via mockable interfaces.

---

## 4. Git Governance & Commit Discipline
* **Strict Single-Change Policy:** Execute only **one atomic change per task/commit** (one feature, one bugfix, or one refactoring). Never bundle unrelated modifications.
* **Conventional Commits Format:**
  `feat(academic): add course semester hosting handler`
  `fix(auth): correct AND-semantics evaluation for scoped assignments`
  `test(hostel): add outpass approval transition unit tests`
  `docs(standards): update documentation index`
* **Signed Commits:** Ensure commits remain cryptographically signed.

---

## 5. Modular Block Construction & Error Standards
* **Unique Block IDs:** Assign every major logical block (handler, state transition, validation pipeline, transaction coordinator) a unique ID header:
  ```go
  // BLOCK_ACAD_HOST_SEMESTER_001
  // Purpose: Hosts a course semester within a department and provisions sections.
  // Inputs:  deptID (string), courseID (string), semesterNum (int)
  // Outputs: HostedSemester, error
  // Errors:  ERR_DEPT_NOT_FOUND, ERR_SEMESTER_ALREADY_HOSTED
  ```
* **Traceable Errors:** All returned or logged errors must cite the unique block ID for instant observability and triage.
* **Flat Logic & Early Returns:** Eliminate nested conditional ladders (Pyramid of Doom). Validate inputs and preconditions with guard clauses and early exits.

---

## 6. API & Transport Standards
* **Explicit Versioning:** All HTTP endpoints follow `/api/v1/...`.
* **RFC 7807 Problem Details:** All error responses must use RFC 7807 format (`type`, `title`, `status`, `detail`, `instance`, `code`, `invalidParams`).
* **Query Empty Semantics:** Collection queries that find zero matching records must return `200 OK` with an empty array (`{"data": []}`), **never** a `404 Not Found`.

---

## 7. Mobile-First Responsive Design (Web & Clients)
* Ensure UI interfaces in SvelteKit and Fyne are fluid and resilient across viewports from **280px to 4K**.
* Never cause horizontal overflow on the root document viewport (`overflow-x: hidden` anti-pattern avoided by fluid containers).
* Touch targets must satisfy a minimum boundary of 44x44 CSS pixels.
