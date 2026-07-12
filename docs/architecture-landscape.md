# IntentCore Architecture Landscape & Specification

**Version:** 1.0.0  
**Status:** Active / Normative  
**Category:** Canonical Architecture Specification

> **One-line Definition**
>
> IntentCore is a transport-agnostic intent coordination kernel that enforces deterministic lifecycle control, authoritative state management, immutable system history, and proof-oriented governance for distributed autonomous systems.

---

# 1. Purpose

This document defines the canonical architecture of IntentCore.

It establishes the normative architectural contracts, system boundaries, dependency rules, and execution model that every implementation MUST follow.

Unless explicitly stated otherwise, the terminology defined in RFC 2119 applies throughout this specification.

This document is normative.

Implementation details MAY evolve over time, but the architectural contracts defined herein MUST remain stable.

---

# 2. Architectural Identity

IntentCore is NOT:

- a message broker
- a workflow engine
- a transport protocol
- a network framework

IntentCore SHALL operate exclusively as an **Intent Coordination Kernel**.

Its responsibilities are limited to:

- intent validation
- admission governance
- deterministic lifecycle control
- authoritative state transitions
- repository consistency
- immutable history
- proof generation
- telemetry production

Everything outside these responsibilities belongs to outer architectural layers.

---

# 3. Architectural Constitution (Normative Invariants)

The following invariants define the supreme architectural rules of IntentCore.

Every package, module, RFC, ADR, implementation, and future extension MUST comply with these rules.

## 3.1 Intent Authority

Every authoritative state mutation MUST originate from a fully validated and admitted Intent.

No component MAY mutate repository state directly.

---

## 3.2 Transport Independence

IntentCore MUST remain completely independent from transport implementations.

Transport technologies MAY evolve without requiring any modification to kernel contracts.

Examples include (but are not limited to):

- ABTP
- TCP
- QUIC
- gRPC
- RDMA

Transport evolution MUST NOT alter kernel semantics.

---

## 3.3 Stateless Transport

Transport implementations MUST remain stateless.

Transport MUST NOT perform:

- lifecycle decisions
- admission decisions
- repository mutations
- policy evaluation
- business logic

Transport exists solely to deliver SemanticEnvelope objects into the kernel.

---

## 3.4 Single Source of Truth

Repository SHALL be the only authoritative state storage.

Every state mutation MUST execute through Compare-And-Swap (CAS).

No alternative mutation path is permitted.

---

## 3.5 Immutable History

Every successful state transition MUST generate immutable evidence.

Historical records MUST be append-only.

Historical records MUST NOT be modified.

Historical records MUST NOT be deleted.

---

## 3.6 Strict Dependency Direction

Execution dependencies MUST always move toward the kernel.

Outer layers SHALL NOT bypass intermediate stages.

Cross-layer mutation is strictly forbidden.

---

# 4. System Boundaries

IntentCore is intentionally divided into explicit architectural boundaries.

## IntentCore

The coordination kernel.

Responsible for:

- validation
- normalization
- admission
- lifecycle
- repository
- history
- proof
- telemetry

---

## ABTP

The transport boundary.

Responsible only for:

- framing
- serialization
- checksum
- version negotiation
- protocol validation
- network communication

ABTP is NOT part of the kernel.

---

## SemanticEnvelope

SemanticEnvelope is the canonical wire contract.

Every external system MUST communicate with IntentCore using SemanticEnvelope.

Transport implementations MAY vary.

SemanticEnvelope MUST remain stable.

---

## Repository

Repository is the Single Source of Truth.

Only Lifecycle is permitted to request authoritative state mutations.

Repository guarantees:

- Compare-And-Swap
- version consistency
- snapshot support
- recovery support
- immutable persistence

---

# 5. Architectural Execution Pipeline

IntentCore enforces a strict one-way execution model.

```
External Systems
        │
        ▼
      ABTP
        │
        ▼
SemanticEnvelope
        │
        ▼
Validation
        │
        ▼
Normalization
        │
        ▼
Admission
        │
        ▼
Lifecycle
        │
        ▼
Repository
        │
        ▼
History
Proof
Telemetry
```

Pipeline execution MUST always remain unidirectional.

No stage MAY skip another stage.

No outer layer MAY directly mutate an inner layer.

---

# 6. Deterministic Lifecycle

IntentCore defines a deterministic lifecycle for every admitted Intent.

`StateUnknown` exists solely as an uninitialized sentinel.

Operational lifecycle begins at `Pending`.

```
Pending
    │
    ▼
Validated
    │
    ▼
Admitted
    │
    ▼
Scheduled
    │
    ▼
Executing
   ╱   ╲
  ▼     ▼
Completed Failed
          │
          ▼
     RolledBack
```

Every transition MUST satisfy:

- deterministic
- authorized
- atomic
- auditable

Transition rules are defined by RFC-0004.

---

# 7. Canonical Repository Layout

The repository SHALL expose architectural boundaries directly through its package layout.

```text
IntentCore/
│
├── cmd/
│   └── intentcored/
│
├── core/
├── lifecycle/
├── admission/
├── state/
├── history/
├── proof/
├── telemetry/
├── runtime/
│
├── transport/
│   ├── transport.go
│   ├── wire/
│   ├── abtp/
│   └── internal/
│
├── docs/
│   ├── adr/
│   ├── rfcs/
│   └── architecture-landscape.md
│
├── go.mod
└── README.md
```

This layout represents the canonical repository structure.

---

# 8. RFC Mapping

IntentCore architecture is governed by the following RFCs.

| RFC | Responsibility |
|------|----------------|
| RFC-0001 | Transport & Wire Contract |
| RFC-0002 | Admission Interface |
| RFC-0003 | Repository & State Topology |
| RFC-0004 | Lifecycle Control |

These RFCs define frozen architectural contracts.

Implementations MUST conform to them.

---

# 9. Non-Goals

IntentCore SHALL NOT become:

- a generic message broker
- a workflow orchestration engine
- an execution runtime
- a transport implementation
- a service mesh
- an API gateway

These responsibilities belong to external systems.

---

# 10. Implementation Philosophy

IntentCore follows Specification-Driven Development.

Architectural order SHALL always be:

```
Architecture

↓

ADR

↓

RFC

↓

Interfaces

↓

Implementation

↓

Testing
```

Implementation MUST follow specifications.

Specifications MUST NOT be derived from implementation.

---

# 11. Informative Appendix (Non-Normative)

The following information is provided for implementation planning only.

It does not define architectural contracts.

## Phase 1

- Core
- Lifecycle
- Repository
- Transport Boundary

## Phase 2

- Runtime Pipeline
- Telemetry
- Proof

## Phase 3

- Distributed Coordination
- Semantic Routing
- Agent Discovery
- Federation

## Phase 4

- Knowledge Plane
- Intent Graph
- Global Coordination
- Zero-Trust Infrastructure

---

# 12. Summary

IntentCore is the architectural center of the system.

ABTP is the transport boundary.

SemanticEnvelope is the canonical wire contract.

Lifecycle is the sole authority for state transitions.

Repository is the single source of truth.

History is immutable.

Architecture is specification-driven.

Every implementation SHALL preserve these architectural contracts.
