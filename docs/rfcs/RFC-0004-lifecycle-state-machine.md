# RFC-0004 — Lifecycle Control

**Status:** Locked
**Category:** Architecture Specification
**Version:** 1.0
**Last Updated:** 2026-07-11

## Abstract

This RFC defines the normative state machine and transition rules for the Lifecycle component. The Lifecycle component is the sole kernel-side authority for authorizing and coordinating state transitions for admitted Intents.

## Motivation

To maintain deterministic behavior, a single authority must govern how an Intent progresses from admission to completion or failure. If multiple components can trigger arbitrary state changes, the system loses its auditability and reasoning guarantees.

## Terminology

*   **Lifecycle Engine:** The core state machine coordinator.
*   **State Transition:** The movement of an Intent from one recognized state to another.
*   **Terminal State:** A state from which no further transitions are allowed.

## Architectural Context

The Lifecycle engine receives accepted Intents from Admission, computes the required state transition, and attempts to commit the new state to the Repository via CAS.
`Admission → Lifecycle → Repository`

## Normative Specification

### 1. Sole Authority
The Lifecycle component SHALL be the sole authority for state transitions. No other component (Transport, Validation, Admission) MAY authorize a state mutation.

### 2. Deterministic State Machine
The Lifecycle MUST operate as a formal, deterministic state machine.
*   All allowed transitions MUST be explicitly defined in a Transition Matrix.
*   Any transition not explicitly allowed MUST be rejected.

### 3. Allowed Transition Matrix
The following transitions are ALLOWED:
*   `Pending` → `Validated`
*   `Validated` → `Admitted`
*   `Admitted` → `Scheduled`
*   `Scheduled` → `Executing`
*   `Executing` → `Completed` (Terminal)
*   `Executing` → `Failed`
*   `Failed` → `RolledBack` (Terminal)

### 4. Atomic Transitions
Transitions MUST be atomic.
*   A transition MUST either complete fully or not happen at all.
*   Partial state updates SHALL NOT be visible to any consumer.
*   The Lifecycle engine MUST coordinate with repository CAS semantics to ensure atomicity.

### 5. Invariants
*   `Completed` and `RolledBack` states MUST be terminal.
*   Transitions MUST NOT bypass authority checks.

## State Model

The Formal State Model is a Directed Acyclic Graph (DAG) with cycles only permitted in strictly defined retry policies (which are out of scope for the base state machine).
```
Pending → Validated → Admitted → Scheduled → Executing → Completed
                                                  ↓
                                                Failed → RolledBack
```

## Interfaces

```go
type Lifecycle interface {
    Transition(ctx context.Context, intent core.IntentID, targetState State) error
}
```

## Error Model

*   `ErrInvalidTransition`: Returned when attempting a transition not in the Allowed Matrix.
*   `ErrTerminalState`: Returned when attempting to mutate a completed or rolled-back Intent.

## Security Considerations

The deterministic state machine prevents attackers from forcing an Intent into an arbitrary state (e.g., jumping from `Pending` directly to `Completed` without passing `Admission`).

## Observability

Every state transition attempted by the Lifecycle MUST be recorded, including both successful commits and rejected transitions (due to CAS failures or matrix violations).

## Compliance Requirements

| Requirement | RFC | Test |
| :--- | :--- | :--- |
| Enforce Allowed Transition Matrix | RFC-0004 | lifecycle_matrix_test.go |
| Reject transitions from Terminal states | RFC-0004 | lifecycle_terminal_test.go |
| Atomicity via Repository CAS | RFC-0004 | lifecycle_atomic_test.go |

## Backward Compatibility

The state machine is frozen. Adding new states or transitions requires a new RFC revision. Implementations MAY evolve internally, but the external transition matrix contract MUST remain stable.

## Rationale

A strict Lifecycle state machine guarantees that the system's operational flow is mathematically verifiable and easy to reason about, which is critical for a high-assurance coordination kernel.

## References

*   RFC-0000 — Architectural Principles
*   RFC-0003 — State Repository
