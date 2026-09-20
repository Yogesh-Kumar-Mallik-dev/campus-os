# 02 — Product Requirements

**Governing Document:** [Campus_OS_Product_Bible_v1.2.docx](file:///home/yogesh/Downloads/Campus_OS_Product_Bible_v1.2.docx)  
**Status:** Canonical

---

## 1. Product Vision & Operating Scope

Campus OS is an institutional operating system designed to serve as the system of record and execution engine for a single collegiate institution.

### 1.1 Guiding Tenets
* **Transparency:** Users clearly see the live status of any request, responsible authority, current stage, and required next action within their authorized scope.
* **Paperwork Elimination:** Manual forms and repeated data collection are converted into typed digital workflows with contextual pre-filling.
* **Frictionless Operation:** Administrative overhead is minimized through rule-based automation, while preserving mandatory institutional checks.
* **Accountability & Evidence:** Every consequential action, authority delegation, and status transition leaves an immutable, reconstructable audit trail.
* **Single-Tenant Digital-First:** Tailored exclusively to the specific institution's physical and administrative realities, supporting digital workflows while providing graceful paper fallback mechanisms where mandatory.

---

## 2. Institutional Personas & Responsibilities

| Role | Operational Scope |
| :--- | :--- |
| **Student** | Requests leaves/outpasses, views personal attendance and finalized marks, submits grievances, tracks fee receipts. |
| **Teacher** | Manages assigned class sessions, submits internal assessment marks, marks attendance, requests leave, reports academic incidents. |
| **Lab Technician** | Manages lab slots, inventory, equipment maintenance, and assists teachers with practical sessions. |
| **Head of Department (HOD)** | Hosts semesters for departmental courses, configures sections, allocates timetables and teacher assignments, reviews department leaves and marks. |
| **Dean Academics** | Formulates curriculum structures, approves academic calendars, reviews departmental marks submission, resolves cross-departmental academic matters. |
| **Registrar** | Institutional custodian of student admissions, enrollment states, degree conferrals, and official institutional master data. |
| **Director / Apex** | Institutional executive authority; approves institutional master data changes, apex escalations, and major policy changes. |
| **Warden** | Final institutional authority for hostel allocations, discipline, curfew management, and student leave/outpass approvals. |
| **Guard** | Physical checkpoint authority at campus gates; scans/verifies Warden approvals and records real-time entry and exit timestamps. |
| **Caterer (Mess / Canteen)** | Manages meal plans, dietary logs, and service status under Estate supervision. |
| **Librarian** | Catalog management, book circulations, fine collections, and library clearance certificates. |
| **Management / Moderator** | Institutional supervisory body; handles grievance oversight, cross-departmental supervision, and administrative checks. |
| **Super Admin** | Exactly one active seat; atomic succession; configures institutional workflows, assigns institutional roles, commits master data changes. |
| **Developer** | System governance, capability-based administrative maintenance, emergency break-glass operator. |

---

## 3. Core Functional Scopes (Phase 1 Priority)

### 3.1 Academic Delivery & Administration
* Curriculum modeling: Courses $\rightarrow$ Semesters $\rightarrow$ Subjects and Labs.
* Departmental semester hosting and operational configuration (Sections, Timetables, Slots).
* Teacher-to-Subject/Section assignment workflow.
* Internal assessments and marks submission with multi-tier verification.

### 3.2 Student Lifecycle Management
* State progression: `APPLICANT` $\rightarrow$ `ENROLLED` $\rightarrow$ `ACTIVE` $\rightarrow$ `SUSPENDED_OR_ON_LEAVE` $\rightarrow$ `ALUMNI_GRADUATED`.
* Terminal state exits: `TERMINATED_DISCONTINUED` and `EXPELLED`.
* Student profile, departmental enrollment, and historical transcript maintenance.

### 3.3 Campus Life & Hostel Administration
* Hostel block, room, and bed allocation workflows.
* Outpass and Leave workflow: Student application $\rightarrow$ Warden digital approval $\rightarrow$ Guard gate verification and logging.
* Student grievance submission and supervisory review by Moderator.

### 3.4 Institutional Governance & Master Data
* Controlled creation and editing of institutional structures (Departments, Courses, Academic Years).
* Single Super Admin atomic handover workflow.
* Developer break-glass emergency response system.
