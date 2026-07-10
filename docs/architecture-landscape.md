# IntentCore Architecture Landscape

IntentCore is a specification-driven intent coordination kernel, with ABTP (AetherBus Transport Protocol) as its transport boundary, designed to enforce deterministic lifecycle control, admission governance, and state consistency for distributed autonomous systems.

## Architecture at a Glance

```text
External Systems
        │
        ▼
+----------------------+
|        ABTP          |
|  Transport Boundary  |
+----------------------+
        │
        ▼
+----------------------+
|  SemanticEnvelope    |
|   (Wire Contract)    |
+----------------------+
        │
        ▼
+----------------------+
|     IntentCore       |
| Coordination Kernel  |
+----------------------+
        │
        ├── Validation
        ├── Normalization
        ├── Admission
        ├── Lifecycle
        ├── State Repository
        ├── History
        ├── Proof
        └── Telemetry
```

## 1. Architecture Boundaries

IntentCore is no longer a message broker in the traditional sense. It is the coordination kernel that governs intent lifecycle, state mutation, authority, and proof-oriented coordination.

| Component | Responsibility | Architectural identity |
| --- | --- | --- |
| Repository / Project: IntentCore | Core kernel for lifecycle, state, admission, and coordination | Kernel |
| Transport Protocol: ABTP | Low-level transport that carries `SemanticEnvelope` into the kernel | Transport boundary |
| Wire Format: `SemanticEnvelope` | Canonical envelope format and metadata contract carried by ABTP | Wire format |
| RFC | Frozen or approved implementation contract | Locked standard |
| Architecture Family: IntentCore Architecture | Full architectural envelope governing structure, flow, and development rules | System architecture |

## 2. Core Contracts

The architectural contracts below are frozen or approved and form the stable foundation of the system.

The canonical RFC mapping is:

| RFC | Scope | Kernel responsibility |
| --- | --- | --- |
| RFC-0001 | Transport / Wire Protocol | Carries `SemanticEnvelope` into the kernel through ABTP |
| RFC-0002 | Admission | Defines the admission interface and decision boundary |
| RFC-0003 | State Repository | Defines the single source of truth and repository mutation primitives |
| RFC-0004 | Lifecycle | Defines lifecycle states, transitions, authority, and history |

Runtime command and execution-level contracts are not core RFCs in this kernel layer. Those concerns belong above IntentCore in runtime or execution components.

### ADR-0001 — IntentCore Internal Architecture

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

`core/` provides the shared contracts, identity types, common abstractions, and type-safe foundations used across the entire kernel.

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

These components are implementation work only. They do not introduce new architectural contracts.

## 4. Internal Data Flow

The system obeys a one-way pipeline only:

```text
External Systems
    ↓
ABTP
    ↓
SemanticEnvelope
    ↓
IntentCore
    ├── Validation
    ├── Normalization
    ├── Admission
    ├── Lifecycle
    ├── State Repository
    ├── History
    ├── Proof
    └── Telemetry
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

The rebrand from AetherBus-Tachyon to IntentCore is therefore not only a rename. It is a change in architectural identity:

- IntentCore is now the repository and architecture name.
- ABTP remains the transport protocol name.
- `SemanticEnvelope` is the wire format carried by ABTP.
- RFC documents are the locked contracts for implementation behavior.
- Legacy references to broker-centric framing are historical only.

This naming model follows separation of concerns and supports the specification-driven architecture model used across ADRs, RFCs, and package boundaries.

## 7. Target Package Structure

The long-term repository shape should make the separation visible in the filesystem:

```text
IntentCore/
│
├── core/
├── lifecycle/
├── admission/
├── state/
├── proof/
├── history/
├── telemetry/
├── runtime/
├── transport/
│   └── aetherbus/
│
├── docs/
│   ├── adr/
│   └── rfc/
│
└── README.md
```

In this structure, ABTP is an implementation of the transport layer for IntentCore. It is not the architectural center of the system.

## 8. Current State Summary

The project is now structurally stable.

- Core Contracts: complete and frozen
- Repository API: stable
- Core: complete
- Lifecycle: complete
- State Repository: in progress
- Architecture contracts: frozen
- Development model: specification-driven
- Next phase: finish `state/`, then build `runtime/pipeline`

## 9. One-line Definition

IntentCore is a specification-driven intent coordination kernel that uses ABTP as its transport boundary to provide deterministic lifecycle control, admission governance, repository-backed state consistency, and proof-oriented coordination for distributed autonomous systems.