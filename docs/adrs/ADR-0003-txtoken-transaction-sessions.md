# ADR-0003: Session-Based TxToken for Cross-Domain Transactions

* **Status:** Accepted / Locked
* **Date:** 2026-09-20
* **Deciders:** Developer, Yogesh

---

## Context

Campus OS permits atomic multi-operation and cross-domain transactions when strictly required, but forbids the Go logical backend from holding direct database connections or executing raw transaction syntax.

## Decision

We implement a **Session-Based TxToken protocol**:

1. Go calls `BeginTx(Context)` on the gRPC persistence service.
2. TypeScript DB layer opens an interactive Prisma transaction and registers a unique `TxToken` with an expiry lease (5 seconds default).
3. Go passes `x-tx-token: <token>` in the gRPC metadata for all operations participating in that atomic unit.
4. Go calls `CommitTx(token)` or `RollbackTx(token)`.
5. Automatic timeout rollback and socket disconnection listeners guarantee transaction reclamation without orphan locks.

## Consequences

* **Positive:** Preserves ACID guarantees for multi-table and cross-domain writes across the Go/TypeScript IPC boundary without opening raw database handles to Go.
* **Negative:** Requires stateful transaction session tracking in the TypeScript persistence service.
