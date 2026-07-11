# Repository Blueprint

**Status:** Draft
**Phase:** Phase 0 — Specification & Blueprint
**Related RFCs:**
- RFC-0003 — State Repository
- RFC-0004 — Lifecycle Control

---

# Purpose

This blueprint defines the long-term storage strategy for IntentCore.

It extends RFC-0003 by describing implementation guidance for persistence, recovery, snapshot management, ledger maintenance, and scalability.

This document is informative and does not define protocol contracts.

---

# Scope

This blueprint covers:

- Repository architecture
- Snapshot strategy
- Recovery process
- Ledger compaction
- Archiving strategy
- Retention policy
- Repository scalability

It does NOT define:

- Lifecycle rules
- Admission policy
- Transport protocol
- Wire format

---

# Architectural Goals

- Single Source of Truth
- Immutable State History
- Event Sourcing
- Deterministic Recovery
- Horizontal Scalability
- Audit Completeness

---

# Repository Architecture

IntentCore Repository consists of:

- State Store
- Append-only Ledger
- Snapshot Store
- Archive Storage

```
Lifecycle
        │
        ▼
Repository
   ├── State
   ├── Ledger
   ├── Snapshot
   └── Archive
```

---

# Snapshot Strategy

Goals

- Fast restart
- Bounded recovery time
- Reduced replay cost

Possible policies

- Periodic snapshot
- Incremental snapshot
- Manual snapshot
- Automatic snapshot

Snapshot metadata

- Snapshot ID
- Repository Version
- Timestamp
- Intent Count
- Ledger Offset

---

# Recovery Strategy

Recovery pipeline

```
Snapshot
      │
      ▼
Replay Ledger
      │
      ▼
Rebuild State
```

Requirements

- Deterministic
- Replay-safe
- Crash-safe

---

# Ledger Management

Ledger remains append-only.

Maintenance includes:

- Compaction
- Archiving
- Retention
- Integrity verification

Ledger entries MUST NEVER be modified.

---

# Archiving Strategy

Old ledger segments MAY be archived.

Requirements:

- Immutable
- Verifiable
- Recoverable

Possible archive targets

- Object Storage
- Cold Storage
- Remote Repository

---

# Retention Policy

Repository SHOULD define:

- retention duration
- archive threshold
- snapshot interval
- cleanup policy

---

# Future Work

- Distributed Repository
- Multi-node Replication
- Consensus Layer
- Incremental Recovery
- Remote Snapshots

---

# References

RFC-0003
RFC-0004
