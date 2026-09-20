# ADR-0005: Canonical Username Formats, Onboarding Workflows & Authentication Architecture

* **Status:** Accepted / Locked
* **Date:** 2026-09-21
* **Deciders:** Developer, Yogesh

---

## Context
In a collegiate institutional operating system (ERP++), account creation cannot allow arbitrary public self-registration. Furthermore, usernames and institutional email addresses must follow deterministic, collision-resistant naming conventions that instantly convey cohort, program, and admission status.

## Decision

### 1. Student Canonical Username Format (Option A)
$$\mathbf{\left[\text{first}\right]} \boldsymbol{.} \mathbf{\left[\text{course}\right]} \boldsymbol{.} \mathbf{\left[\text{YYYY}\right]} \left[\boldsymbol{.} \mathbf{\text{l}}\right] \left[\boldsymbol{.} \mathbf{\text{increment}}\right]$$

* **`first`:** Sanitized lowercase first name (alphanumeric only).
* **`course`:** Canonical lowercase course/program code (e.g. `cse`, `ece`, `me`, `btech`).
* **`YYYY`:** 4-digit entry year (e.g. `2024`, `2026`).
* **`l`:** Appended with dot delimiter **strictly if** `admissionType == "LATERAL_ENTRY"`.
* **`increment`:** Appended with dot delimiter upon collision (`.1`, `.2`...).
* **Institutional Email:** `<username>@campus.edu`.

**Examples:**
* Regular Student: `yogesh.cse.2024` $\rightarrow$ Collision: `yogesh.cse.2024.1`
* Lateral Entry: `yogesh.cse.2024.l` $\rightarrow$ Collision: `yogesh.cse.2024.l.1`

### 2. Faculty & Staff Canonical Username Format
$$\mathbf{\left[\text{first}\right]} \boldsymbol{.} \mathbf{\left[\text{last}\right]} \boldsymbol{.} \mathbf{\left[\text{YYYY}\right]} \left[\boldsymbol{.} \mathbf{\text{increment}}\right]$$

* **`first`:** Sanitized lowercase first name.
* **`last`:** Sanitized lowercase last name.
* **`YYYY`:** 4-digit appointment/joining year.
* **`increment`:** Appended with dot delimiter upon collision (`.1`, `.2`...).
* **Institutional Email:** `<username>@campus.edu`.

**Examples:**
* Faculty: `amit.sharma.2026` $\rightarrow$ Collision: `amit.sharma.2026.1`

### 3. Onboarding & Account Claiming Protocols
1. **Admissions Pre-Ingestion:** The Admissions / Registrar office ingests verified student records. Usernames and emails are pre-calculated and reserved.
2. **Student Account Claim:** Students claim accounts via `PRN + Date of Birth + Mobile OTP`.
   * Lateral entry students must review and acknowledge their prior diploma details, Semester 3 starting state, and course exemptions.
3. **Faculty Account Invitation:** HR issues a 72-hour cryptographic invitation token (`/activate?token=...`). Faculty sets credentials and configures mandatory TOTP Multi-Factor Authentication (MFA).
4. **Operational Staff Provisioning:** Estate and Gate Guards are onboarded via supervisor-provisioned Phone Number + PIN credentials.

### 4. Authentication & Session Security
* **Universal Login Identifiers:** Login accepts `username` (e.g. `yogesh.cse.2024.l` or `amit.sharma.2026`), institutional email, or PRN.
* **Password Hashing:** Argon2id (memory=64MB, iterations=3, parallelism=4).
* **Session Lifecycle:** 15-minute EdDSA/Ed25519 JWT Access Token + 7-day HttpOnly cookie Refresh Token with automatic Family Token rotation and breach invalidation.
* **Brute-Force Defense:** Rate limiting enforcing temporary lockout after 5 consecutive failed attempts per user/IP.

## Consequences
* **Positive:** Completely eliminates ambiguous identities; institutional emails and usernames are predictable, human-readable, and deterministic; lateral entry status is transparent; security guardrails adhere to RFC 7807 and zero-trust standards.
* **Negative:** Requires atomic sequence/collision checks during batch student admissions ingestion.
