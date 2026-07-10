RFC-0004: Lifecycle Control / State Machine

Status: Frozen
Authors: IntentCore Architecture Team
Created: 2026-07-10
Updated: 2026-07-10
Dependencies: RFC-0000, RFC-0001, RFC-0002, RFC-0003
Implements: Lifecycle state machine for intent coordination

1. Abstract

This RFC defines the canonical lifecycle state machine for IntentCore. It specifies the lifecycle states, valid transitions, transition authority, atomicity requirements, rollback model, retry semantics, and immutable transition history.

2. Motivation

IntentCore requires a deterministic and enforceable lifecycle model so that every intent can move through the system under explicit authority and auditable rules. Lifecycle control must be isolated from transport, admission, and repository implementation details while remaining the only authorized path for state mutation.

3. Scope

This RFC defines:

- lifecycle states
- transition rules
- transition authority
- atomic state mutation requirements
- rollback and retry semantics
- immutable transition history
- lifecycle-related invariants

This RFC does not define transport framing, admission policy internals, repository persistence mechanics, or proof construction semantics beyond lifecycle-facing history records.

4. Terminology

Lifecycle State
A canonical stage in the life of an intent.

Transition
A state change from one lifecycle state to another.

Authority
The component or role allowed to authorize or trigger a specific transition.

Atomic Transition
A state change that must occur all-or-nothing, without partial application.

History Record
An immutable audit record describing a completed transition.

5. Lifecycle States

The lifecycle consists of the following frozen states:

1. "StateUnknown"
2. "StatePending"
3. "StateValidated"
4. "StateAdmitted"
5. "StateScheduled"
6. "StateExecuting"
7. "StateCompleted"
8. "StateFailed"
9. "StateRolledBack"

"StateUnknown" is a non-operational sentinel and MUST NOT be treated as a valid runtime lifecycle position.

6. Architecture

Lifecycle is the sole authority for state mutation. It sits between admission and the state repository.

IntentCore
  → Admission
  → Lifecycle StateMachine
  → State Repository (CAS)
  → History

Lifecycle MUST be the only kernel-side component allowed to authorize and coordinate state transitions.

7. Transition Rules

The following transitions are allowed:

- "StatePending" → "StateValidated"
- "StateValidated" → "StateAdmitted"
- "StateAdmitted" → "StateScheduled"
- "StateScheduled" → "StateExecuting"
- "StateExecuting" → "StateCompleted"
- "StateExecuting" → "StateFailed"
- "StateFailed" → "StateRolledBack"

Forbidden transitions

- "StateCompleted" → any state
- "StateRolledBack" → any state
- "StateUnknown" → operational states without initialization
- any backward transition not explicitly listed above
- any transition that bypasses authority checks or repository version checks

8. Transition Authority

Each transition MUST be authorized by the correct authority.

Transition| Authority
"StatePending" → "StateValidated"| Validation
"StateValidated" → "StateAdmitted"| Admission
"StateAdmitted" → "StateScheduled"| Scheduler
"StateScheduled" → "StateExecuting"| Runtime
"StateExecuting" → "StateCompleted"| Runtime
"StateExecuting" → "StateFailed"| Runtime
"StateFailed" → "StateRolledBack"| Rollback

Requirements

- Lifecycle MUST enforce authority before state mutation.
- Unauthorized transitions MUST be rejected.
- Authority checks MUST be explicit and deterministic.
- Authority MUST NOT be inferred from transport details alone.
- Authority MUST be represented as a stable contract, not a stringly-typed guess.

9. Atomic Transition Requirements

Transitions MUST be atomic.

Requirements

- A transition MUST either complete fully or not happen at all.
- Partial state updates MUST NOT be visible.
- Transition version checks MUST occur before commit.
- The lifecycle engine SHOULD coordinate with repository CAS semantics.
- Exclusive ownership of an intent SHOULD be enforced during mutation.
- Retry-safe recovery MAY occur only after failure is fully resolved.

10. Retry and Rollback Semantics

Retry

Retry semantics MAY be supported for failed intents where policy permits it.

- Retry MUST NOT bypass authority or version rules.
- Retry SHOULD be policy-controlled.
- Retry MUST preserve auditability.

Rollback

Rollback is the compensating path from failure to terminal rollback.

- Rollback MUST be explicit.
- Rollback MUST be auditable.
- Rollback MUST be idempotent where possible.
- Rollback MUST NOT be treated as silent undo.
- Rollback MUST preserve the history of the failure that triggered it.

11. Transition History

Every transition MUST be recorded as an immutable history entry.

History Requirements

- Each history record MUST include intent identity.
- Each record MUST include previous and next state.
- Each record MUST include transition authority.
- Each record SHOULD include version before and after mutation.
- Each record SHOULD include timestamp metadata.
- History MUST be append-only from the perspective of consumers.
- History MUST NOT be externally mutable after publication.

12. Invariants

The following invariants MUST hold:

- Only Lifecycle MAY authorize state transitions.
- Only valid transitions MAY be committed.
- Completed and RolledBack states are terminal.
- Atomicity MUST prevent partial transition visibility.
- Every transition MUST be traceable through history.
- Repository version checks MUST be respected by lifecycle mutation.
- Transition authority MUST match the permitted transition table.

13. Compatibility

This RFC is frozen. Future lifecycle evolution MUST preserve the state machine, authority model, and atomic transition semantics.

Compatibility rules:

- New states MUST require a new RFC revision.
- New transitions MUST require a new RFC revision.
- Transition authority MAY be refined only if backward compatibility is maintained.
- The lifecycle engine implementation MAY evolve, but the contract MUST remain stable.

14. Security Considerations

Lifecycle control is a privileged boundary.

The system SHOULD protect against:

- unauthorized state mutation
- replayed transition requests
- stale version commits
- cross-layer mutation bypass
- concurrent transition races
- malicious rollback attempts

Lifecycle components MUST treat incoming requests as untrusted until authority and version rules are satisfied.

15. Reference Implementation Notes

A reference implementation MAY use:

- a transition table
- an authority checker
- an atomic transition helper
- an immutable history recorder

The implementation MUST preserve the frozen contract even if internal mechanics differ.

16. Future Work

Future revisions MAY define:

- richer retry policies
- distributed intent locking
- cross-shard transition coordination
- state-machine event sourcing
- more detailed failure categories
- formal verification of transition properties

These additions MUST preserve the frozen lifecycle contract defined by this RFC.
