# Blueprint: Snapshots and Archiving Strategy

**Related RFCs:** RFC-0003 (State Repository)

## 1. Introduction

IntentCore's State Repository utilizes an event-sourced, append-only ledger to ensure absolute immutability and precise history (RFC-0003). While this guarantees an auditable trace of every intent lifecycle transition, it poses long-term scalability challenges in terms of storage size and recovery time.

This document outlines the strategic blueprint for maintaining performance and unbounded durability through a structured snapshot and archiving strategy.

## 2. Snapshot Strategy

To avoid replaying the entire history of the system upon recovery or initialization, the repository MUST support periodic snapshots.

### 2.1 Incremental and Full Snapshots
*   **Full Snapshots:** Captures the complete active state of all non-terminal intents at a specific ledger offset.
*   **Incremental Snapshots:** Captures only the state mutations that have occurred since the last snapshot.

### 2.2 Triggering Mechanisms
Snapshots should be triggered based on configurable thresholds:
*   **Time-based:** Every $N$ hours.
*   **Volume-based:** After $M$ ledger entries are appended.
*   **Manual/Governance:** On-demand triggering via administrative action.

### 2.3 Atomicity
Snapshot creation MUST NOT block ongoing CAS operations. Implementations should utilize concurrent structures (e.g., MVCC, Copy-On-Write) to generate a consistent snapshot of the state at a precise ledger offset without halting system throughput.

## 3. Archiving and Ledger Compaction

As the append-only ledger grows indefinitely, older entries must be archived to "cold" storage to keep the active ledger ("hot" storage) fast and manageable.

### 3.1 Cold Storage Transition
*   Ledger entries preceding the oldest required snapshot can be safely migrated to cold storage.
*   Cold storage can be an object store (e.g., S3), specialized time-series databases, or compressed archival files.

### 3.2 Immutability in Archive
Archived records MUST retain their immutability guarantees. Cryptographic hashes (e.g., Merkle trees) should be used to verify the integrity of cold storage records against the live system's historical root.

## 4. Recovery Procedure

The recovery process should follow these steps:
1.  **Locate Latest Snapshot:** Retrieve the most recent verified snapshot.
2.  **Restore State:** Load the snapshot into the fast-access Cache.
3.  **Replay Ledger:** Replay the append-only ledger starting from the snapshot's exact offset to the tip of the ledger.
4.  **Verify Integrity:** Ensure the resulting state matches the expected hash/checksum.

This approach ensures the system can recover quickly while preserving the single source of truth defined in RFC-0003.
