# 09 — Engineering Standards & AI Anti-Pattern Prevention System

**Derived from:** [Universal Engineering Standards](https://github.com/Yogesh-Kumar-Mallik-dev/engineering-standards)  
**Status:** Authoritative Repository Standard • Binding on all Human Engineers and AI Agents

---

## 1. Executive Summary & Purpose

Campus OS pairs autonomous AI coding agents with human engineers to construct an institutional ERP++ platform. Without rigid engineering guardrails, codebases suffer from code-doc drift, hallucinated dependencies, silent error swallowing, nested conditional spaghetti, and broken responsive layouts.

This document establishes the **authoritative engineering standards** for Campus OS, adapted directly from the audited universal engineering standards repository.

---

## 2. Elimination of the 10 AI Anti-Patterns

| Anti-Pattern | Description | Mandatory Remedy in Campus OS |
| :--- | :--- | :--- |
| **1. Hallucinated APIs & Imports** | Inventing uninstalled libraries or fantasy methods. | Only import packages declared in `go.mod` and `package.json`. Verify imports against local toolchains before committing. |
| **2. Lazy Placeholder Stubs** | Leaving `// TODO: implement later` or empty stubs. | Complete full, production-ready implementations with proper error handling and tests in the exact same task. |
| **3. Code-Doc Drift** | Updating code while leaving specs or READMEs outdated. | Code and documentation must update hand-in-hand. Any schema or behavioral change requires updating corresponding `/docs`. |
| **4. Unverified Tool Execution** | Assuming commands succeed without checking return codes. | Verify tool and command outputs synchronously. Inspect compiler and linter output before marking steps complete. |
| **5. Inverted Technical Depth** | Over-engineering boilerplate while trivializing domain logic. | Focus cognitive effort on domain rules, state machines, and scoped authorization; use minimal, idiomatic scaffolding. |
| **6. Silent Error Swallowing** | Using empty `catch` blocks or discarding Go errors (`_ = fn()`). | Every error must be explicitly handled, enriched with context (citing Block ID), and propagated or returned. |
| **7. Blind File Overwrites** | Obliterating existing implementations instead of surgical edits. | Use surgical edits. Preserve existing unrelated functions, comments, and signatures. |
| **8. Multi-Change Commit Sprawl** | Bundling multiple features or unrelated refactors into one commit. | **Strict Single-Change Policy**: Exactly one feature, one fix, or one refactoring per commit. Multi-task commits are rejected. |
| **9. Pyramid of Doom** | Deeply nested `if/else` ladders (nesting depth $\ge 3$). | Enforce guard clauses and early returns. Flatten all execution paths. |
| **10. Desktop-Only Fixed Layouts** | Building UI assuming 1920x1080 resolution with rigid `px` widths. | **Mobile-First & Viewport Resilience**: UI must adapt gracefully from 280px to 4K with fluid containers and zero root horizontal overflow. |

---

## 3. Modular Code Construction & Block Standards

All non-trivial functions, handlers, state transitions, and coordination logic must be structured into self-describing modular blocks with unique block headers:

```go
// BLOCK_ACAD_ASSIGN_TEACHER_001
// Purpose: Assigns a teacher to a course subject and section within a hosted semester.
// Inputs:  ctx context.Context, req *AssignTeacherRequest
// Outputs: *AssignTeacherResponse, error
// Errors:  ERR_TEACHER_NOT_FOUND, ERR_SECTION_NOT_FOUND, ERR_CONFLICTING_TIMETABLE
func (s *AcademicService) AssignTeacher(ctx context.Context, req *AssignTeacherRequest) (*AssignTeacherResponse, error) {
	// Guard clause 1: Validate input presence
	if req.TeacherId == "" {
		return nil, fmt.Errorf("BLOCK_ACAD_ASSIGN_TEACHER_001: teacher_id is required: %w", ErrInvalidInput)
	}

	// Guard clause 2: Check scope authorization
	if !s.authz.CanAssignTeacher(ctx, req.DepartmentId) {
		return nil, fmt.Errorf("BLOCK_ACAD_ASSIGN_TEACHER_001: unauthorized department scope: %w", ErrForbidden)
	}

	// Domain logic execution
	return s.store.AssignTeacher(ctx, req)
}
```

### Invariants:
1. **Traceable Errors:** Every returned or logged error message **must include the Block ID** (`BLOCK_<DOMAIN>_<ACTION>_<ID>`).
2. **Early Returns:** Preconditions, input validations, and tenancy checks are handled first via guard clauses; happy paths remain un-indented at the root function level.

---

## 4. Testing Rigor & Dependency Injection

1. **Mandatory Hand-in-Hand Tests:** Unit tests must be written or updated simultaneously with code modifications. No feature is complete without accompanying tests.
2. **Co-located Test Files:**
   * Go: `service.go` $\rightarrow$ `service_test.go` in the same directory.
   * TypeScript / Web: `handler.ts` $\rightarrow$ `handler.spec.ts` co-located with the source.
3. **100% Domain Logic Target:** Core business logic, state machines (Student Lifecycle, Outpass Approvals), and authorization evaluators must achieve comprehensive test coverage.
4. **Mockable Interfaces:** Persistence layers and external clients must be injected via mockable interfaces:
   ```go
   type AcademicStore interface {
       AssignTeacher(ctx context.Context, req *AssignTeacherRequest) (*AssignTeacherResponse, error)
   }
   ```

---

## 5. REST & gRPC API Standards

### 5.1 Collection Query Semantics (200 OK vs 404)
* **Empty Collections Return `200 OK`:** When querying collections (e.g. `GET /api/v1/departments/{id}/courses`), if the parent resource exists and the user has access, but no matching entities exist, the API **must return `200 OK` with `{"data": []}`**. Returning `404 Not Found` for an empty collection is strictly prohibited.
* **Single Resource Lookups Return `404 Not Found`:** If a specific ID is queried (`GET /api/v1/students/{id}`) and does not exist, return `404 Not Found`.

### 5.2 Error Envelopes (RFC 7807 Problem Details)
All HTTP error responses must adhere to the RFC 7807 specification:
```json
{
  "type": "https://campus-os.internal/errors/insufficient-permissions",
  "title": "Insufficient Permissions",
  "status": 403,
  "detail": "BLOCK_AUTH_EVAL_002: User lacks scoped assignment for Subject CS401 in Section B",
  "instance": "/api/v1/marks/cs401/section-b",
  "code": "ERR_FORBIDDEN_SCOPE",
  "invalidParams": []
}
```

### 5.3 IPC Error Isolation & UDS Path Concealment
* **Zero Leakage of Domain Sockets:** Internal socket paths (`/tmp/*.sock`, `/var/run/*.sock`) or raw gRPC dial failure strings must **never** be included in HTTP problem detail payloads or client UI status labels.
* **Transport-to-Problem Mapping:** All gRPC status codes must pass through `mapGRPCError`, translating transport faults (`codes.Unavailable`) to standard RFC 7807 problem envelopes with HTTP 503 (`SERVICE_UNAVAILABLE`).

---

## 6. Git Governance & Commit Discipline

1. **Strict Single-Change Policy:** One PR / commit = one logical change. Never combine a bug fix with a new feature or un-requested refactoring.
2. **Conventional Commits Format:**
   $$\text{<type>}(\text{<scope>}): \text{<concise expression>}$$
   * Supported types: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `chore`, `ci`.
   * Scopes: `academic`, `student`, `hostel`, `finance`, `auth`, `audit`, `db`, `api`, `ui`.
3. **Signed Commits:** Every commit must be cryptographically signed (`git commit -S`).
4. **Linear Git History:** Feature branches must rebase onto `main` before merging; merge commits on `main` are prohibited.

---

## 7. Mobile-First & Viewport Resilience (Web & Fyne Clients)

* **Viewport Support Range:** The web client (SvelteKit) and mobile/desktop clients (Fyne) must be fully functional across viewports from **280px to 4K**.
* **Zero Root Horizontal Overflow:** The root document (`html`, `body`) and native canvas must never produce an unintended horizontal scrollbar.
* **Fluid Spacing & Typography:** Use CSS `clamp()`, `rem`, and container-relative units rather than static pixel breakpoints.
* **Touch Targets:** Interactive elements (buttons, inputs, links) must meet the minimum touch target dimension of **44x44 CSS pixels** on mobile viewports.

---

## 8. Cross-Platform Lifecycle Orchestration

Developers and CI systems execute standard lifecycle tasks using the unified cross-platform orchestrator:
* **POSIX / Linux / macOS:** `./script.sh <command>`
* **Windows PowerShell:** `.\script.ps1 <command>`

Supported lifecycle targets: `dev`, `build`, `check`, `test`, `deps`, `proto:gen`, `envi`, `uenvi`, `flush`, `help`.

---

## 9. Fyne UI/UX Anti-Pattern Guardrails (2026 Modern Product Standards)

When building native client interfaces in Go Fyne, engineers must adhere to [`.agent/rules/ui_ux_anti_patterns.md`](file:///home/yogesh/campus_os/.agent/rules/ui_ux_anti_patterns.md):
1. **Never Replicate Legacy Desktop Apps:** Avoid Windows Forms-style beveled buttons, heavy gradients, grey panels, and dense toolbars. Interfaces must mirror modern 2026 products (Linear, Raycast, Notion, Vercel).
2. **Avoid Excessive Borders & Nested Cards:** Rely on spacing, background surfaces (`#10161C` canvas vs `#18202A` card), and typography hierarchy. Never create triple-nested cards (`Card -> Card -> Card`).
3. **No Raw Status Strings:** Replace `Status: ACTIVE` with scannable visual badges (`● Active`).
4. **Pure Vector Graphics:** Zero emojis, zero generic raster icons. Utilize pure vector Lucide SVGs with 24×24 geometry and uniform 2px stroke.
5. **Progressive Disclosure & Clear Hierarchy:** Do not cram all administrative data onto a single screen. Prioritize primary actions, use clear touch targets (min 44×44px), and provide explicit state feedback (loading, empty, success, error).
6. **Mobile-First Responsive Layouts & Breakpoints:** Enforce adaptive breakpoints (`<640` Mobile, `640..1024` Tablet, `>1024` Desktop). On compact/mobile viewports, navigation drawers must dock cleanly to the left edge with a semi-transparent backdrop overlay and ESC key listener rather than floating in the center.
7. **Accessible Modals & Overlays:** All modal dialogs (`ShowModal`) must support keyboard `Escape` dismissal and outside-click backdrop dismissal (enabled by default, optionally configurable).
8. **Responsive Card Grids & Text Wrapping:** Grid layouts must dynamically wrap (`AdaptiveGridLayout`, `FlowLayout`) and card labels must specify `TextWrapWord` to eliminate horizontal canvas clipping.

