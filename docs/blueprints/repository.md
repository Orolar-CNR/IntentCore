# Repository Blueprint

**Phase:** Phase 0 — Specification & Blueprint
**Related RFCs:** RFC-0003 (State Repository), RFC-0004 (Lifecycle Control)

## Purpose
This blueprint defines the long-term storage strategy for IntentCore. It provides implementation guidance on persistence, recovery, snapshot management, and how the append-only ledger operates under the hood.

## Repository Architecture
The Repository is expected to be built around an Event Sourcing model, utilizing an append-only ledger for all state mutations.

```text
Lifecycle
        │ (CAS mutations)
        ▼
Repository
   ├── State Cache (In-Memory / Fast Access)
   ├── Ledger (Append-only Immutable Log)
   ├── Snapshot Store (Periodic checkpoints)
   └── Archive (Cold storage for old ledgers)
```

## Snapshot and Recovery Strategy
Because an append-only ledger grows indefinitely, reading the entire ledger from genesis to rebuild state is impractical.
The system is expected to perform periodic snapshots:
1.  **Snapshotting:** The current aggregate state is serialized to a Snapshot Store at a specific Ledger Offset.
2.  **Recovery:** On startup, the system loads the latest snapshot, then replays only the ledger events that occurred *after* the snapshot's offset.

## Compaction and Archiving
To manage disk space, older ledger segments that have been fully captured in a stable snapshot can be archived to cold storage (e.g., AWS S3, local tape backups). The active repository node only needs to retain the recent ledger segments necessary for immediate recovery and current CAS validation.

## Storage Independence
The repository interface is designed to be agnostic to the underlying database engine. Implementations could utilize PostgreSQL, specialized event stores (like EventStoreDB), or highly optimized in-memory stores backed by WAL (Write-Ahead Logging), as long as they satisfy the CAS and immutability contracts of RFC-0003.
