# Security Policy & Vulnerability Disclosure — Campus OS

## 1. Supported Versions

| Version | Supported |
| :--- | :--- |
| Current `main` | Yes |
| Prior Major Releases | Critical Security Patches Only |

---

## 2. Reporting a Vulnerability

If you discover a security vulnerability in Campus OS, please **do NOT open a public issue**. Instead, follow these steps:

1. **Email Disclosure:** Send an email with full reproduction steps to `security@campus-os.internal` (or report via GitHub Private Vulnerability Reporting).
2. **Include Details:**
   - Detailed description of the vulnerability.
   - Proof of concept (PoC) code or script.
   - Affected domains, endpoints, or gRPC RPCs.
   - Potential impact and threat assessment.
3. **Response SLA:**
   - Acknowledgement within **24 hours**.
   - Triage and mitigation timeline within **72 hours**.
   - Coordinated security patch delivery.

---

## 3. Defense-in-Depth Security Invariants

1. **Zero Direct User Permissions:** All access evaluates through `Person -> Scoped Role Assignment -> Role -> Permissions`.
2. **Multi-Dimensional Scope Enforcement:** Permissions are filtered via strict `AND` evaluation across supported dimensions (`Department`, `Course`, `Semester`, `Section`, `Subject`, `Lab`, `HostelBlockFloor`).
3. **Zero Direct DB Access from Domain Logic:** The Go logical backend communicates exclusively via canonical gRPC handles over Unix Domain Sockets; no SQL injection vector can bridge directly to domain logic.
4. **Developer Break-Glass Protocol:** Emergency access requires Developer cryptographic keys, emits instant immutable audit logs, and triggers an automated data freeze upon termination.
