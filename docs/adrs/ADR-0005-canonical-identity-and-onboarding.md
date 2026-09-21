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

#### 3.1 Genesis Bootstrapping: Super Admin (Chairperson) CLI Provisioning
The platform enforces a strict zero-trust root-of-trust bootstrap:
1. **Zero-Remote Creation Invariant:** The Super Admin account can **never** be created over HTTP, GraphQL, or external APIs. It must be initiated directly on the bare-metal host server via the backend CLI.
2. **Backend CLI Execution:**
   ```bash
   campus-backend bootstrap-superadmin \
     --name="Chairperson Full Name" \
     --email="chairperson@trust.edu" \
     --phone="+9198XXXXXXXX"
   ```
3. **Genesis Verification Protocol:**
   * Validates that `auth_schema.super_admin_seat` is vacant (`singleton` check).
   * Generates a 256-bit Ed25519 single-use cryptographic claim token bound to the Chairperson's registered phone number.
   * Renders an ANSI ASCII QR code directly to the secure terminal `stdout` and outputs an encrypted one-time claim token.
   * Records the genesis installation event in `audit_schema.audit_logs`.
4. **Chairperson Mobile SIM Claim:**
   * The Chairperson scans the terminal QR using the Campus OS native client.
   * SIM-presence and SMS-OTP verification are completed on the Chairperson's device.
   * Chairperson sets the master password, configures hardware/TOTP MFA, and crowns the active `super_admin_seat`.

#### 3.2 Executive Council Tier-1 Provisioning (Chairperson Authority)
The Chairperson (Super Admin) possesses exclusive institutional authority (`executive.qr.generate`) to generate sealed credential QR codes for the Top Institutional Executive Council:

| Executive Role | Institutional Mandate | Provisioned By | Verification & Security Invariant |
| :--- | :--- | :--- | :--- |
| **Director** | Academic & Campus Administrative Head | **Chairperson (Super Admin)** | Sealed QR + SIM-Bound Mobile OTP + Mandatory MFA |
| **Executive Director (ED)** | Institutional Operations, Finance & Infrastructure | **Chairperson (Super Admin)** | Sealed QR + SIM-Bound Mobile OTP + Mandatory MFA |
| **Dean** | Dean of Academic Affairs / Dean of Student Welfare | **Chairperson (Super Admin)** | Sealed QR + SIM-Bound Mobile OTP + Mandatory MFA |
| **Registrar** | Official Custodian of University Records, Admissions & Exams | **Chairperson (Super Admin)** | Sealed QR + SIM-Bound Mobile OTP + Mandatory MFA |

* **Zero-Knowledge Invariant:** The Chairperson never sees, sets, or handles the passwords of executive officers. The generated QR contains a single-use token that triggers private password creation and TOTP MFA setup during the executive's first scan.

#### 3.3 Cascading Institutional Authority Matrix
Credentials and onboarding QRs flow through a strict, multi-tiered hierarchy of authority:

| Tier | Target Cohort | Provisioning Authority | Permission Code | Mandate & Evidence Prerequisite |
| :--- | :--- | :--- | :--- | :--- |
| **Tier 0 (Genesis)** | **Super Admin (Chairperson)** | **Backend CLI (Local Server)** | `CLI_BARE_METAL_ONLY` | Physical bare-metal server access; singleton seat check. |
| **Tier 1 (Executives)** | **Director, Executive Director, Dean, Registrar** | **Chairperson (Super Admin)** | `executive.qr.generate` | Governing Trust / Board of Trustees resolution orders. |
| **Tier 2 (Students)** | **Students (UG, PG, Lateral Entry)** | **Registrar / Admissions Office** | `admissions.qr.generate_bulk` | Pre-verified, gazetted admission master list with phone & personal email. |
| **Tier 2 (Faculty)** | **Teaching Faculty, HODs, Lab Techs** | **Human Resources (HR) & Dean** | `hr.faculty_qr.generate_bulk` | Executed appointment letters and verified employee master files. |
| **Tier 2 (Operations)** | **Wardens, Estate Workers, Security Guards** | **Executive Director & CSO / Estate** | `staff.qr.generate_bulk` | Background verification and active operational deployment rosters. |

* **Zero Faculty Generation Invariant:** General teaching faculty and department staff have **zero authority** to generate student, peer, or executive QR credentials.
* **Immutable Batch Manifest:** Every QR generation batch creates an immutable audit log (`audit.qr_batch_manifest`) recording `batch_id`, operator `user_id`, cohort type, target count, cryptographic hash of the batch, generation timestamp, and workstation terminal ID.

#### 3.4 Zero-Faculty-Knowledge QR Delivery
* Credentials and account claim payloads are delivered directly to users via **sealed, tamper-evident physical QR slips** printed during official enrollment or issued on identity cards.
* **Cryptographic Claim Token:** The sealed QR encodes a signed single-use claim token:
  $$\text{ClaimToken} = \text{Sign}_{\text{Ed25519}}\left(\text{UserID}, \text{MobileHash}, \text{BatchID}, \text{Nonce}, \text{Expiry}\right)$$
  The token has an absolute expiry (e.g. 14 days from issue).

#### 3.5 SIM-Presence Hardware Binding & Mobile OTP Verification
Claiming an account requires dual verification combining physical QR possession with active SIM hardware binding:
1. **Device SIM Verification:** When scanned via the Campus OS native mobile client, the application probes device telephony services to confirm that the physical/eSIM currently inserted in the device matches the pre-registered mobile number recorded during admissions.
2. **SIM-Bound SMS OTP:** An automated one-time passkey (OTP) is dispatched over SMS to the registered mobile number.
3. **Hardware Anti-Theft Lock:** The onboarding flow verifies that the OTP is received and processed directly on the device hosting that active SIM card (via platform SMS verification APIs).
4. **Rejection of Remote Claims:** If an unauthorized party steals or photographs a physical QR slip, they cannot claim the account on their own device because the physical SIM is absent (`ERR_SIM_ABSENT_OR_MISMATCH`).
5. **Workstation / Laboratory Fallback:** In institutional computer laboratories or on desktop clients where mobile telephony hardware is absent, onboarding requires a Proctored Dual-Factor mode (SMS OTP sent to student's phone + physical proctor / admission officer verification code).

#### 3.6 First-Scan Mandatory Password Reset & 3-Day Grace Period
* **Mandatory Initial Update:** Scanning the QR and completing SIM verification immediately places the user in an isolated onboarding state. The platform requires the user to set a personal password before any institutional resources can be accessed.
* **Days 1–3 (Orientation Grace Period):** To eliminate cognitive fatigue and boarding friction during campus orientation days, users may set a simplified password (minimum 6 characters, no complex character rules required).
* **Post Day 3 (Hard Compliance Enforcement):** After 72 hours from first claim, any user with a simplified password is automatically blocked with an unskippable credential upgrade modal enforcing enterprise complexity (minimum 10 characters, upper, lower, numeric, and special characters).
* **Lateral Entry Review:** Lateral entry students review and acknowledge their prior polytechnic/diploma credits, Semester 3 placement, and course exemptions before password commitment.

#### 3.7 Post-Onboarding Mobile Number Update Workflow
Users are completely free to change their registered mobile phone number after initial onboarding:
1. **Self-Service Dual-OTP Verification:**
   * User logs into their authenticated profile and requests a phone number update.
   * Step-up authentication is enforced (current password or biometric verification).
   * **Dual Verification:** An authorization OTP is sent to the *current* registered mobile number, and a verification OTP is sent to the *new* mobile number. Both must be entered to finalize the update.
2. **Lost / Damaged SIM Exception Procedure:**
   * If the user has lost access to their old SIM card, self-service dual-OTP cannot complete.
   * The user must visit the Registrar’s Office (for students) or HR (for staff) in person.
   * The designated officer performs biometric / physical government ID verification and executes an audited break-glass contact update (`auth.mobile_override_update`).

### 4. Authentication & Session Security
* **Universal Login Identifiers:** Login accepts canonical username (e.g. `yogesh.cse.2024.l` or `amit.sharma.2026`), institutional email (`<username>@campus.edu`), or University PRN.
* **Password Hashing:** Argon2id (memory=64MB, iterations=3, parallelism=4).
* **Session Lifecycle:** 15-minute EdDSA/Ed25519 JWT Access Token paired with a 7-day HttpOnly cookie Refresh Token with automatic Family Token rotation and breach invalidation.
* **Brute-Force Defense:** Exponential backoff and account lock after 5 consecutive failed attempts per identity or IP address.

### 5. Dual-Identity Architecture: Academic vs. Legal Alt-Names & Secondary School Discrepancies

#### 5.1 The Institutional Reality in Indian Education
In the Indian educational ecosystem, secondary school examination boards (CBSE, ICSE, NIOS, State Boards) notoriously record student and parent names with severe discrepancies compared to statutory government identification (Aadhaar, Passport, PAN, Birth Certificates):
* **Truncated / Mononym Discrepancy:** A student whose full legal name is `Yogesh Kumar Mallik` may be registered solely as `Yogesh` on their Class X matriculation certificate.
* **Initials vs. Expansion:** A student named `K. V. Ramesh` on board rolls may hold `Kandukuri Venkata Ramesh` on their Passport and Aadhaar.
* **Parental Discrepancy:** Parents' names on marksheets frequently omit surnames, truncate initials (e.g. `R K Mallik` vs `Rajesh Kumar Mallik`), or list a mother's maiden name without update.
* **The Single-Field Trap:** Forcing a single `full_name` field guarantees administrative failure:
  * Using the legal name on grade cards causes state universities and exam boards to reject transcripts for failing to match the Class X prerequisite record.
  * Using the academic name on placement, scholarship, or visa letters causes national DBT/PFMS scholarship portals, banks, and corporate background checkers (EPFO/TCS/Infosys KYC) to fail identity verification.

#### 5.2 Tri-Name Data Modeling
Campus OS explicitly segregates identity into three distinct first-class constructs across `auth_schema` and `student_schema`:

1. **`academic_name` (Matriculation Record Name):**
   * Verbatim character sequence as printed on the student's Class X Secondary School Certificate and university registration board rolls (e.g. `Yogesh`).
   * **Governing Rule:** Sole authority for all academic, examination, and graduation artifacts.
2. **`legal_full_name` (Statutory National Identity):**
   * Statutory legal identity recorded on Government Photo IDs (Aadhaar, Passport, PAN, Voter ID) (e.g. `Yogesh Kumar Mallik`).
   * **Governing Rule:** Sole authority for financial, statutory, legal, immigration, and placement KYC records.
3. **`preferred_name` (Display / Community Name):**
   * Everyday name chosen by the user for UI greeting, student directory, forums, and campus messaging.
4. **Canonical Username Derivation Rule:**
   * In $\mathbf{\left[\text{first}\right]} \boldsymbol{.} \mathbf{\left[\text{course}\right]} \boldsymbol{.} \mathbf{\left[\text{YYYY}\right]}$, the $\mathbf{\left[\text{first}\right]}$ component is extracted from the primary token of `academic_name` (e.g. `yogesh`), ensuring immediate congruence with examination roll sheets and attendance registers.

#### 5.3 Parents & Legal Guardian Discrepancy Modeling
Both sets of parents and legal guardians are modeled with decoupled academic vs. legal identities in `student_profiles`:

| Persona | Academic Record (`academic_name`) | Statutory Record (`legal_name`) | Additional Attributes |
| :--- | :--- | :--- | :--- |
| **Father** | `father_academic_name` (as per Class X marksheet) | `father_legal_name` (as per Aadhaar/Passport) | Phone, Email, Occupation |
| **Mother** | `mother_academic_name` (as per Class X marksheet) | `mother_legal_name` (as per Aadhaar/Passport) | Phone, Email, Occupation |
| **Legal Guardian** | `guardian_academic_name` (as per school admission) | `guardian_legal_name` (as per guardianship deed) | Relationship, Phone, Email, Address |

#### 5.4 Document Generation Routing Matrix
The automated document generation engine in `document_schema` enforces strict field-level routing based on the regulatory destination of the document:

| Document Category | Documents Generated | Consumed Name Fields |
| :--- | :--- | :--- |
| **Academic & Examination** | Semester Grade Cards, Consolidated Transcripts, Degree Parchments, Provisional Certificates, Migration Certificates, Exam Hall Tickets / Admit Cards | `academic_name`<br>`father_academic_name`<br>`mother_academic_name` |
| **Statutory & Financial** | Direct Benefit Transfer (DBT) / National Scholarship (PFMS) verification slips, Official Fee Invoices, Bank Education Loan Letters | `legal_full_name`<br>`father_legal_name`<br>`mother_legal_name` |
| **Immigration & Placement** | Passport / Visa Bonafide Letters, Embassy Verification Certificates, Campus Placement / Corporate KYC Dossiers | `legal_full_name`<br>`academic_name` (as alias reference) |
| **Legal & Disciplinary** | Disciplinary Committee Summon Notices, Court Affidavits, Institutional Show-Cause Notices | `legal_full_name`<br>`guardian_legal_name` |

#### 5.5 Institutional Dual-Verification Certificate (The "One and the Same" Resolver)
When a student with discrepant records applies for higher studies abroad, passports, or corporate background checks, Campus OS provides a cryptographically verifiable **Dual-Verification Certificate** template (`DOC_CERT_DUAL_IDENTITY`):

$$\begin{aligned}
\text{“This is to certify that } &\mathbf{\text{[legal\_full\_name]}}\text{, holder of Govt ID/Passport }\mathbf{\text{[GovtID]}}\text{, is a bona fide student} \\
\text{of this Institute enrolled in } &\mathbf{\text{[course]}}\text{ under PRN }\mathbf{\text{[prn]}}\text{. It is further certified that their name is officially} \\
\text{recorded as } &\mathbf{\text{[academic\_name]}}\text{ in their Secondary School Examination (Class X) Certificate} \\
\text{and University Examination Ledgers, and that both names refer to one and the same individual.”}
\end{aligned}$$

* **Evidentiary Attachment:** The database maintains `has_name_discrepancy = true` and links to `name_discrepancy_affidavit_doc_id` whenever a notarized affidavit or gazette notification is submitted.

## Consequences
* **Positive:** Completely eliminates ambiguous identities; institutional emails and usernames are predictable, human-readable, and deterministic; lateral entry status is transparent; security guardrails adhere to RFC 7807 and zero-trust standards; physical QR slips cannot be abused without possession of the registered SIM card; bulk issuance authority is strictly audited and limited to authorized administrative offices; solves the pervasive Indian academic vs. legal name discrepancy permanently across all university transcripts, degree certificates, scholarships, and placements.
* **Negative:** Requires dual data capture during student onboarding and admissions ingestion; requires document templates to explicitly declare their regulatory identity routing.
