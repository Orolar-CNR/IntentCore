# IntentCore Architecture Landscape

IntentCore is a specification-driven intent coordination kernel, with AetherBus as its transport protocol, designed to enforce deterministic lifecycle control, admission governance, and state consistency for distributed autonomous systems.

## 1. Architecture Boundaries

IntentCore is no longer a message broker in the traditional sense. It is the coordination kernel that governs intent lifecycle, state mutation, authority, and proof-oriented coordination.

| Component | Responsibility |
| --- | --- |
| Repository / Project: IntentCore | Core kernel for lifecycle, state, admission, and coordination |
| Transport / Wire Protocol: AetherBus | Low-level transport that carries `SemanticEnvelope` into the kernel |
| Messaging / Wire Format: AetherBus Protocol | Canonical frame/envelope format and metadata contract |
| Architecture Family: IntentCore Architecture | Full architectural envelope governing structure, flow, and development rules |

## 2. Core Contracts

The architectural contracts below are frozen or approved and form the stable foundation of the system.

### ADR-0001 — Broker Internal Architecture

**Status:** Approved

Defines internal development rules:

- Interface-first design
- Single-direction pipeline
- Strict module ownership
- No layer crossing
- Event-driven flow
- Specification-driven implementation

### RFC-0001 — Transport & Wire Protocol

**Status:** Frozen

Defines `SemanticEnvelope`, validation, normalization, and wire-level constraints.

### RFC-0002 — Intent Admission Interface

**Status:** Frozen

Defines `AdmissionPolicy` and the admission decision boundary before execution.

### RFC-0003 — State Topology / State Repository

**Status:** Frozen

Defines the repository as the single source of truth, including:

- `CompareAndSwap`
- `CommitmentLedger`
- `Snapshot`
- `Recovery`
- `StateVersion`

### RFC-0004 — Lifecycle Control / State Machine

**Status:** Frozen

Defines:

- the 8-state lifecycle
- transition table
- authority enforcement
- atomic transitions
- immutable history

## 3. Implementation Status

### Phase A — Core Contracts

**Status:** Complete

`core/` is the shared dictionary and type-safety layer of the system.

Key elements:

- `IDGenerator`
- `IntentID`, `TraceID`, `TransitionID`, `NodeID`
- `CommitmentState`
- `StateVersion`
- `TransitionRequest`, `TransitionResult`, `TransitionRecord`
- `BrokerError` and error codes
- `BrokerContext`
- `Clock`, `Event`, and runtime abstractions

### Phase B — Lifecycle Module

**Status:** Complete

`lifecycle/` is the central authority for state mutation.

Key elements:

- `StateMachine`
- transition rules
- authority enforcement
- atomic transition handling
- immutable transition history

### Phase C — State Repository

**Status:** In progress

`state/` is the memory and truth layer of the system.

Implemented:

- `StateRepository` interface
- thread-safe `InMemoryRepository`
- `CompareAndSwap` as the only mutation primitive
- copy-in / copy-out protection
- context-aware operations
- clock injection

Pending:

- `ledger.go`
- `version.go`
- `snapshot.go`
- `recovery.go`
- `health.go`

## 4. Internal Data Flow

The system obeys a one-way pipeline only:

```text
SemanticEnvelope
  → Transport
  → Validation
  → Normalization
  → Admission
  → StateMachine
  → Repository (CAS)
  → History
  → Proof
  → Telemetry
```

No layer is allowed to mutate a lower or unrelated layer directly.

## 5. Architectural Principles

IntentCore is built on the following principles:

- Simple contract in the middle
- Evolution at the edges
- No business logic in transport
- State mutation only through the lifecycle engine
- Repository as the single source of truth
- Immutable history for auditability
- Frozen contracts for stability
- Specification-driven development as the default development model

## 6. Historical Context

The system originally lived under the AetherBus name, where the transport and message-routing idea first took shape. As the architecture matured, the project was re-centered around the actual responsibility of the kernel: intent coordination, lifecycle control, and state governance.

For this reason:

- IntentCore is now the repository and architecture name.
- AetherBus remains the transport / wire protocol name.
- Legacy references to broker-centric framing are historical only.

## 7. Current State Summary

The project is now structurally stable.

- Core: complete
- Lifecycle: complete
- State Repository: in progress
- Architecture contracts: frozen
- Development model: specification-driven
- Next phase: finish `state/`, then build `runtime/pipeline`

## 8. One-line Definition

IntentCore is a frozen-contract intent coordination kernel with AetherBus transport, deterministic lifecycle control, strict admission governance, and repository-backed state consistency for distributed autonomous systems.
