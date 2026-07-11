RFC-0003: State Topology / State Repository

Status: Locked
Authors: IntentCore Architecture Team
Created: 2026-07-10
Updated: 2026-07-10
Dependencies: RFC-0000, RFC-0001, RFC-0002, RFC-0004
Implements: State repository contract for IntentCore

1. Abstract

This RFC defines the state topology of IntentCore. It specifies the repository as the single source of truth for intent state, including the mutation primitive, state versioning, commitment ledger, snapshot model, recovery model, and repository invariants.

2. Motivation

IntentCore requires a stable and auditable state layer that can support deterministic lifecycle control, concurrent access, crash recovery, and future distributed coordination. The state repository must be isolated from transport and admission concerns while remaining the authoritative source of truth for lifecycle state.

3. Scope

This RFC defines:

- the repository as the canonical state boundary
- the "IntentState" model
- repository mutation semantics
- optimistic concurrency via version checks
- commitment ledger semantics
- snapshot and recovery requirements
- repository health and versioning semantics
- state visibility and ownership rules

This RFC does not define transport framing, admission policy evaluation, lifecycle transition rules, or proof generation behavior beyond repository-facing persistence of state facts.

4. Terminology

State Repository
The system component responsible for storing and retrieving authoritative intent state.

IntentState
The current canonical state of a single intent at a point in time.

Commitment Ledger
The append-oriented history of lifecycle-relevant state facts, commitments, or transitions.

Snapshot
A point-in-time representation of repository state used for recovery or checkpointing.

State Version
A monotonically increasing version used to support optimistic concurrency control.

5. Architecture

The repository sits below the lifecycle layer and above persistence-specific implementation details.

IntentCore
  → Lifecycle
  → State Repository
  → Commitment Ledger
  → Snapshot / Recovery

The repository MUST be the only kernel-side authority for durable state mutation.

6. Contract

The canonical repository contract is equivalent to the following capabilities:

Get(intentID) -> IntentState
Put(intentState) -> error
CompareAndSwap(intentID, expectedVersion, newState) -> error
Persist() -> error
Snapshot() -> Snapshot
Recover(snapshot) -> error
Health() -> RepositoryHealth
Version() -> StateVersion

Requirements

- The repository MUST be the single source of truth for authoritative state.
- The repository MUST support optimistic concurrency through version checks.
- The repository MUST support copy-in / copy-out behavior to prevent external mutation.
- The repository MUST expose read and mutation operations through a stable interface.
- The repository MUST be isolated from transport and admission layers.
- The repository MUST NOT permit direct mutation of internal storage from outside the repository boundary.

7. IntentState Model

"IntentState" represents the current authoritative state of an intent.

Recommended fields:

- "IntentID"
- "State"
- "Version"
- "UpdatedAt"
- "Metadata"

Requirements

- "IntentID" MUST uniquely identify a single intent instance.
- "State" MUST be a valid lifecycle state.
- "Version" MUST increase on successful state mutation.
- "UpdatedAt" SHOULD reflect the last repository-applied update time.
- "Metadata" MAY store non-authoritative context needed by adjacent layers.

8. Mutation Semantics

The repository MUST treat "CompareAndSwap" as the strict mutation primitive for authoritative state updates.

Requirements

- "CompareAndSwap" MUST succeed only when the expected version matches the current version.
- "CompareAndSwap" MUST fail with a version mismatch when the state is stale.
- "CompareAndSwap" MUST be the primary mechanism used by lifecycle transitions.
- "Put" MAY be used for initialization or controlled insertion, but MUST NOT bypass repository invariants.
- "Put" MUST preserve repository consistency.

9. Commitment Ledger

The commitment ledger records authoritative state facts over time.

Requirements

- The ledger SHOULD be append-oriented.
- The ledger MUST preserve historical intent state facts.
- The ledger MUST support audit and recovery use cases.
- The ledger MUST NOT be directly mutated by transport or admission layers.
- The ledger SHOULD remain logically separate from the current-state record.

10. Snapshot Model

The repository MUST support snapshotting for checkpoint and recovery purposes.

Snapshot SHOULD include:

- current authoritative intent state
- repository version
- active ledger-relevant state
- recovery metadata
- timestamp or checkpoint metadata

Requirements

- Snapshot creation MUST be deterministic for a given repository state.
- Snapshot data MUST be internally consistent.
- Snapshot content MUST be sufficient to restore repository state without reconstructing from transport traffic.
- Snapshot generation SHOULD preserve intent identity and version information.

11. Recovery Model

Recovery restores the repository from a snapshot and associated durable state.

Recovery Requirements

- Recovery MUST restore authoritative repository state.
- Recovery MUST preserve state identity and version metadata.
- Recovery SHOULD rebuild any derived or transient structures required by the repository implementation.
- Recovery MUST not invent state that is not present in the source snapshot or durable ledger.
- Recovery MUST be safe to repeat if needed.

12. Concurrency Model

The repository MUST support concurrent access safely.

Requirements

- Reads MAY occur concurrently.
- Writes MUST be synchronized.
- Concurrent mutation of the same intent MUST NOT corrupt state.
- Copy-in / copy-out behavior SHOULD prevent accidental mutation leakage.
- Version checks MUST be used to detect stale writers.

13. State Visibility

The repository defines visibility boundaries for each component.

Component| Visibility
Current Intent State| Repository + Lifecycle
Commitment Ledger| Repository + History/Audit
Snapshot| Repository + Recovery
Health| Repository + Observability
Version| Repository contract consumers

14. Invariants

The following invariants MUST hold:

- Every stored "IntentState" MUST have a valid identity.
- Every authoritative state update MUST increase or preserve consistency rules.
- Version mismatches MUST be detectable.
- No external component MAY mutate internal repository storage directly.
- Snapshot and recovery MUST preserve repository identity semantics.
- The repository MUST remain the sole authority for durable state mutation.

15. Compatibility

This RFC is frozen. Future state-layer evolution MUST preserve the repository as the single source of truth.

Compatibility rules:

- New repository implementations MAY be added without breaking the contract.
- Persistence backends MAY change as long as repository semantics are preserved.
- Breaking changes MUST require a new RFC revision.
- Future sharding MUST preserve identity, version, and mutation semantics.

16. Security Considerations

The repository MUST assume that upstream inputs may be malformed or adversarial.

The system SHOULD guard against:

- stale writes
- state replay
- unauthorized mutation
- concurrent update races
- corrupted snapshots
- inconsistent recovery inputs

Repository access MUST remain isolated from untrusted transport inputs.

17. Reference Implementation Notes

A prototype repository MAY use an in-memory map plus synchronization primitives, provided that:

- the interface contract is preserved
- CAS semantics remain authoritative
- copy-in / copy-out protection is enforced
- context-aware operations are supported where applicable

18. Future Work

Future revisions MAY define:

- distributed repository backends
- append-only durable ledger implementations
- incremental snapshots
- online migration
- cross-shard state ownership transfer
- repository replication protocols

These additions MUST preserve the frozen topology and authority model defined by this RFC.
