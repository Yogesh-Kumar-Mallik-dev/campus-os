# 05 — Workflow & Approval Specification

**Governing Document:** [Campus_OS_Product_Bible_v1.2.docx](file:///home/yogesh/Downloads/Campus_OS_Product_Bible_v1.2.docx)  
**Status:** Canonical & Locked

---

## 1. Workflow Engine Principles

1. **Versioned Definitions:** Every workflow blueprint is versioned.
2. **In-Flight Immutability:** Requests already in progress remain bound to the workflow version with which they started; subsequent updates to a workflow definition apply solely to new requests.
3. **No Direct Overrides:** Automated rule-based outcomes cannot be overwritten directly. Exceptions must execute through a designated policy-exception workflow.
4. **Submitted Decision Immutability:** An approver may revise a draft decision freely, but once formally submitted, the decision is immutable.

---

## 2. Predefined Approval Policies

Administrators choose from standard, system-provided policy presets rather than constructing custom combinations:

```mermaid
flowchart TD
    Req["Workflow Request Submitted"]
    Req --> Pol{"Evaluated Approval Policy"}
    Pol -->|Single-Decisive| SD["Designated Decisive Approver Approves"]
    Pol -->|All-Must-Approve| AMA["100% Unanimous Approval Required"]
    Pol -->|Majority-Rules| MR["> 50% of Eligible Approvers Approve"]
    Pol -->|Threshold| TH["Configured Quorum (e.g. 75%) Approved"]
    SD --> Done["Request Approved"]
    AMA --> Done
    MR --> Done
    TH --> Done
```

### 2.1 Presets Catalog

1. **Single-Decisive:** Exactly one designated role holds final authority. Additional participants provide advisory input.
2. **All-Must-Approve (Unanimous):** Every participant in the approval chain must vote in favor. A single rejection terminates or returns the workflow.
3. **Majority-Rules (>50%):** Requires more than half of the eligible committee votes.
4. **Quorum / Threshold:** Requires a predefined quorum percentage (e.g. 75%) of assigned members.

---

## 3. Workflow Interaction Rules

* **Mandatory Rejection Comment:** Rejection requires a mandatory documented reason. Approval comments are optional.
* **Backward-Only Information Requests:** An approver may request clarifications only backwards along the existing approval path to previous participants.
* **One-Time Veto / Return:** An approver may return a request to an earlier stage exactly once while the final criteria remain unmet. Once the final approval criteria are satisfied, the request is closed and dissenters cannot reopen it via the veto mechanism.

---

## 4. Phase 1 Core Workflow Scenarios

### 4.1 Hostel Leave & Gate Outpass

```mermaid
sequenceDiagram
    autonumber
    actor Student
    actor Warden
    actor Guard
    Student->>Warden: Submit Outpass Request (Dates, Reason, Emergency Contact)
    Note over Warden: Evaluates curfew, past attendance & discipline
    alt Approved
        Warden->>Student: Digital Approval Issued (QR Token Generated)
        Student->>Guard: Present Digital Pass at Campus Gate
        Guard->>Guard: Scan QR & Verify Identity
        Guard->>System: Log Departure Timestamp
        Note over Student: Off-Campus Period
        Student->>Guard: Return to Campus Gate
        Guard->>System: Scan QR & Log Arrival Timestamp (Auto-Closes Pass)
    else Rejected
        Warden->>Student: Request Rejected (Mandatory Comment Recorded)
    end
```

### 4.2 Academic Marks & Assessment Finalization

$$\text{Teacher Submits Marks} \longrightarrow \text{HOD Departmental Review} \longrightarrow \text{Registrar / Dean Verification} \longrightarrow \text{Academic Records Locked}$$

### 4.3 Affiliating University (AKTU) Semester Result Ingestion & Conflict Resolution

```mermaid
sequenceDiagram
    autonumber
    actor Uploader as Student / Exam Cell Officer
    actor Engine as PDF Parser Engine
    actor COE as Exam Cell / COE
    
    Uploader->>Engine: Upload Result PDF (Gazette / OneView Scorecard)
    Engine->>Engine: Hash & Store PDF in document_schema.documents
    Engine->>Engine: Parse Metadata, Subject Rows, SGPA, & Result Status
    alt No Baseline or Exact Match
        Engine->>Engine: Generate New Revision & Lock (VERIFIED_LOCKED)
    else Discrepancy with Existing Verified Record
        Engine->>Engine: Transition Status to CONFLICT_RESOLUTION_NEEDED
        Engine->>COE: Route to Actionable Resolution Queue
        Note over COE: Side-by-Side Comparison of Extracted Marks & PDFs
        COE->>Engine: Sign Off (Accept Student / Accept Institution / Manual Override)
        Engine->>Engine: Commit New Revision with Audit Delta & Lock
    end
```

### 4.4 Master Data Governance Pipeline

$$\text{Registrar / Dean Academics Proposes} \longrightarrow \text{Director Formally Approves} \longrightarrow \text{Super Admin Atomically Commits}$$

---

## 5. Audit Records vs. Decision Records

| Dimension | Audit Record (`audit_schema`) | Decision Record (`audit_schema`) |
| :--- | :--- | :--- |
| **Primary Question** | *What happened?* | *Why was this decision reached?* |
| **Trigger** | Every state-modifying transactional mutation. | Every rule-based or human approval evaluation. |
| **Data Captured** | Actor ID, timestamp, entity changed, previous state, new state, client IP. | Evaluated policy version, input attributes, criteria satisfied/failed, voting breakdown, rationale. |
| **Immutability** | Cryptographically append-only; zero updates/deletions. | Cryptographically append-only; zero updates/deletions. |
