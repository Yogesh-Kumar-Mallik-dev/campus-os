# ADR-0007: Super Admin Provisioning, OS Keyring Sessions, and Base Dashboard Shell

* **Status:** Accepted / Locked
* **Date:** 2026-09-21
* **Deciders:** Yogesh, Developer
* **Related Documentation:** [01 — Architecture Constitution](file:///home/yogesh/campus_os/docs/01-architecture-constitution.md), [10 — Super Admin & Dashboard Architecture](file:///home/yogesh/campus_os/docs/10-super-admin-and-dashboard-architecture.md)

---

## Context

Campus OS requires an absolute single-tenant institutional governance structure. Following successful implementation of student onboarding and frontend-backend synchronization, the system requires:
1. A secure mechanism to bootstrap and activate the institutional Chairperson (`SUPER_ADMIN`).
2. Persistent local session management in the native client that survives app restarts securely.
3. A responsive, multi-role dashboard foundation adhering to Design System 2026 (pure Go, native Fyne widgets, high-contrast accessible styling).

---

## Decisions

### 1. Genesis Super Admin Lifecycle
* **Provisioning:** Strictly initiated via server bare-metal CLI (`campus-server bootstrap-superadmin`). No web/UI bypass.
* **Docket Output:** Dual output — terminal ASCII QR matrix plus an exportable high-resolution image docket (`genesis_docket.png` / `.svg`) on disk for physical sealing.
* **Telephony Binding:** Strict hardware SIM match enforced against the registered phone number (`+919876543210`).
* **Credentials:** Standard unified password setup without mandatory immediate TOTP 2FA.

### 2. Smart Client Routing & OS Keyring Session
* **Smart Entrypoint:** At launch, queries the native OS Keyring. If a valid session exists, enters the Dashboard directly. Otherwise, presents the Institutional Login screen with an "Activate with Sealed QR Docket" button.
* **Session Storage:** Native Linux OS Keyring via Secret Service D-Bus API (`org.freedesktop.secrets`).
* **Dual-Revocation:** Logging out purges in-memory state, wipes the OS keyring entry, and dispatches a revocation request to the backend.

### 3. Base Dashboard Shell Layout
* **Sidebar Geometry:** Fixed `220px` desktop sidebar with responsive slide-out drawer on compact/mobile screens.
* **TopBar Layout:**
  * Left: Hamburger toggle (compact/mobile) + BBDIT Logo.
  * Center: Current Page Name (`EXECUTIVE OVERVIEW`).
  * Right: Theme toggle (Dark / Light) + Backend live status ping + Profile pill (`SUPER ADMIN`) + Logout.
* **Chairperson Navigation Items:**
  1. Executive Overview
  2. Academic Structure
  3. Governance & Approvals
  4. Audit Ledger
  5. Executive Credentialing (Dean, Registrar, Director, Executive Director access tokens)
* **Primary KPI Cards:**
  1. Total Scholars (`1,420 Enrolled • 98.4% Active`)
  2. Academic Departments (`8 Active • 34 Hosted Semesters`)
  3. Pending Presidential Approvals (`3 Requiring Chairperson Signature`)
  4. Total Staff & Faculty (`142 Active Officers & Educators`)

---

## Consequences

* **Positive:**
  - High security: Genesis creation cannot be triggered over an exposed network endpoint.
  - Native desktop experience: Keyring integration avoids storing plaintext tokens on disk.
  - Clear structural hierarchy: Super Admin has direct executive oversight and credentialing powers from Day 1.
* **Negative:**
  - Linux headless setups without D-Bus Secret Service daemon require a graceful fallback to encrypted local config.
