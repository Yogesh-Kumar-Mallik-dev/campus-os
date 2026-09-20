# 08 — AI Development Constitution & Agent Guardrails

**Governing Document:** [Campus_OS_Product_Bible_v1.2.docx](file:///home/yogesh/Downloads/Campus_OS_Product_Bible_v1.2.docx)  
**Status:** Mandatory & Non-Negotiable

---

## 1. Prime Directives for AI Coding Assistants

All AI agents assisting in the development of Campus OS must strictly adhere to the following rules:

1. **The Product Bible is a Governing Law, Not a Suggestion:** Agents must never treat specifications as loose suggestions or override locked architectural choices for convenience.
2. **Never Invent Institutional Structures:** Agents must not assume or fabricate reporting chains, approval rules, or institutional policies marked as TBD or unverified.
3. **Preserve the Persistence Boundary:**
   * Never import Prisma or database drivers into the Go logical backend.
   * Never write raw SQL queries from domain business logic.
   * All persistence operations must route through canonical TypeScript gRPC handles.
4. **No Direct User-to-Permission Grants:** Always enforce `Person -> Scoped Role Assignment -> Role -> Permissions`.
5. **No Normal Hard Deletes:** Never generate generic `DELETE FROM` endpoints or methods in normal operational APIs.
6. **No Silent Resolutions:** Ingestion errors, data conflicts, or unmapped values must be captured as discrepancies with audit provenance, never silently swallowed.
7. **Approval Decision Immutability:** Never introduce mechanisms to edit, overwrite, or mutate a submitted approval decision.
8. **Require Formal ADRs for Architectural Changes:** Any proposed modification to the locked architecture or technology stack requires human approval and a formal Architecture Decision Record (ADR).
9. **Explainable and Auditable Code:** All generated code must include clear tests, structured error handling, and correlation IDs for auditability.
