# RFC-0003 — State Repository

**Status:** Locked
**Category:** Architecture Specification
**Version:** 1.0
**Last Updated:** 2026-07-11

## Abstract

This RFC specifies the normative requirements for the IntentCore State Repository. The Repository is the isolated, canonical single source of truth for the system, responsible for persisting intent states, managing versions, and maintaining an immutable history.

## Motivation

Distributed coordination requires an authoritative source of state. If multiple components maintain their own state or mutate state directly without concurrency control, data corruption is inevitable. The Repository centralizes all state mutations behind a strict Compare-And-Swap (CAS) interface to guarantee consistency.

## Terminology

*   **Repository:** The canonical single source of truth.
*   **Compare-And-Swap (CAS):** An atomic operation ensuring state is only updated if the expected version matches the current version.
*   **Ledger:** The append-only log of historical state transitions.

## Architectural Context

The Repository is the bottom-most layer of the core pipeline. It is only accessible by the Lifecycle component for mutations.
`Lifecycle → Repository`
Note: History is generated as a result of Lifecycle transitions, not as a direct output of the Repository itself.

## Normative Specification

### 1. Single Source of Truth
The Repository SHALL be the sole, canonical source of truth for the entire system. No other component MAY cache authoritative state or perform direct mutations.

### 2. State Mutation Restrictions
All state mutations within IntentCore MUST occur exclusively via Compare-And-Swap (CAS) in the state repository.
*   Direct overwrites without version checking ARE STRICTLY PROHIBITED.
*   If a CAS operation fails due to a version mismatch, the mutation MUST be rejected, and the Lifecycle MUST handle the conflict (e.g., retry or fail).

### 3. Business Logic Isolation
The Repository SHALL NEVER contain business logic, transition logic, or policy evaluation. It acts purely as a dumb, highly consistent storage mechanism.

### 4. Immutable History
Every successful CAS operation MUST result in an immutable entry appended to the History Ledger. Past states or historical records SHALL NOT be rewritten or altered under any circumstances.

### 5. Data Governance
The repository requires formal snapshot and archiving policies to prevent memory bottlenecks from its append-only ledger structure.

## State Model

State Versioning Model:
`State_v(N) + Transition -> State_v(N+1)`
If `CAS(Expected: N, Next: N+1)` succeeds, commit. Else, abort.

## Interfaces

```go
type StateRepository interface {
    LoadIntent(ctx context.Context, id core.IntentID) (*IntentState, error)
    CompareAndSwap(ctx context.Context, current Version, next IntentState) error
    Snapshot(ctx context.Context) (*Snapshot, error)
    Recover(ctx context.Context, snapshot Snapshot) error
}
```

## Error Model

*   `ErrVersionConflict`: Emitted when CAS fails.
*   `ErrNotFound`: Emitted when an IntentID does not exist.

## Security Considerations

By forcing all mutations through CAS, the system prevents race conditions and concurrent modification attacks. The isolation of the Repository ensures that a compromised transport or admission layer cannot directly overwrite database records.

## Observability

Every successful CAS MUST emit an event to the History component for auditability.

## Compliance Requirements

| Requirement | RFC | Test |
| :--- | :--- | :--- |
| Enforce CAS on mutation | RFC-0003 | repository_cas_test.go |
| Reject business logic in Repo | RFC-0003 | architecture_linter_test.go |
| Append-only history | RFC-0003 | repository_history_test.go |

## Backward Compatibility

The CAS contract and append-only ledger requirements are frozen. Storage backends (e.g., PostgreSQL, memory, Redis) may change, but they MUST implement the exact `StateRepository` contract.

## Rationale

The Repository must remain incredibly simple to guarantee correctness. Moving all business logic into the Lifecycle and Admission layers ensures the storage tier focuses solely on I/O and consistency.

## References

*   RFC-0000 — Architectural Principles
*   RFC-0004 — Lifecycle Control
