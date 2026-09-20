# ADR-0001: Go Logical Backend & TypeScript DB Boundary via gRPC over UDS

* **Status:** Accepted / Locked
* **Date:** 2026-09-20
* **Deciders:** Developer, Yogesh

---

## Context
The Campus OS Product Bible locks the system architecture to a Go Logical Backend for domain logic and workflows, and a separate TypeScript + Prisma persistence layer. Direct database access from Go is forbidden to enforce business-level persistence boundaries.

## Decision
We establish **gRPC with Protocol Buffers over a local Unix Domain Socket (`/var/run/campus-os/db.sock`)** as the official IPC boundary.
* Contract toolchain: Single monorepo `/proto` directory managed via `buf`.
* Deployment: Go and TypeScript run in isolated containers inside Docker Compose, sharing a mounted socket volume.

## Consequences
* **Positive:** High throughput, low latency (avoids TCP loopback overhead), strong compile-time type safety across Go and TypeScript.
* **Negative:** Requires running two runtime containers (Go and Node/Bun) in deployment; contract modifications require code generation step.
