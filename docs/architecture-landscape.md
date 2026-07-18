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

It establishes the normative architectural contracts, execution model, dependency rules, and system boundaries that every implementation SHALL preserve.

Unless explicitly stated otherwise, the terminology defined in RFC 2119 applies throughout this document.

Implementation details MAY evolve.

Architectural contracts MUST remain stable.

---

# 2. Architectural Identity

IntentCore is NOT:

- a message broker
- a workflow engine
- a transport protocol
- a service mesh
- a network framework

IntentCore SHALL operate exclusively as an Intent Coordination Kernel.

Its responsibilities are limited to:

- validation
- normalization
- admission governance
- deterministic lifecycle control
- authoritative state transitions
- repository consistency
- immutable history
- proof generation
- telemetry production

Everything else belongs to outer architectural layers.

---

# 3. Architectural Constitution (Normative Invariants)

Every implementation SHALL preserve the following invariants.

## 3.1 Intent Authority

Every authoritative state mutation MUST originate from a validated and admitted Intent.

Repository state MUST NOT be modified directly.

Only Lifecycle MAY request authoritative mutations.

---

## 3.2 Transport Independence

IntentCore MUST remain transport agnostic.

Transport implementations MAY evolve independently.

Examples include:

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

- lifecycle evaluation
- admission
- repository mutation
- policy evaluation
- business logic

Transport exists solely to deliver SemanticEnvelope objects into the kernel.

---

## 3.4 Single Source of Truth

Repository SHALL be the only authoritative state storage.

Every mutation MUST execute through Compare-And-Swap (CAS).

No alternative mutation path is permitted.

---

## 3.5 Immutable History

Every successful lifecycle transition MUST emit immutable historical evidence.

History SHALL be append-only.

Historical records MUST NOT be modified.

Historical records MUST NOT be deleted.

---

## 3.6 Strict Dependency Direction

Execution dependencies SHALL always move toward the kernel.

Cross-layer mutation is forbidden.

Outer layers MUST NOT bypass intermediate stages.

---

# 4. System Boundaries

## IntentCore

Responsible for:

- Validation
- Normalization
- Admission
- Lifecycle
- Repository
- History
- Proof
- Telemetry

---

## ABTP

Transport boundary only.

Responsible for:

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

Every external producer MUST communicate using SemanticEnvelope.

---

## Repository

Repository is the authoritative state boundary.

Repository guarantees:

- Compare-And-Swap (CAS)
- version consistency
- authoritative state storage
- snapshot support
- recovery support

Repository does NOT generate History, Proof, or Telemetry.

Those artifacts originate from Lifecycle.

---

# 5. Canonical Execution Pipeline

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
   ├──────────────► History
   ├──────────────► Proof
   ├──────────────► Telemetry
   │
   ▼
Repository
   ├──────────────► State Store
   │                    │
   │                    ├────────► State Cache
   │                    └────────► Snapshot Store
   │
   └──────────────► Ledger
                           │
                           ▼
                        Archive
```

Execution MUST remain strictly one-way.

No stage MAY bypass another stage.

---

# 6. Deterministic Lifecycle

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

Every transition MUST be:

- deterministic
- authorized
- atomic
- auditable

RFC-0004 defines the transition matrix.

---

# 7. Canonical Repository Layout

```text
IntentCore/
│
├── cmd/
├── contracts/
├── core/
├── admission/
├── lifecycle/
├── runtime/
├── state/
├── history/
├── proof/
├── telemetry/
├── transport/
├── internal/
└── docs/
```

---

# 8. RFC Mapping

| RFC | Responsibility |
|------|----------------|
| RFC-0001 | Semantic Envelope |
| RFC-0002 | Admission |
| RFC-0003 | Repository |
| RFC-0004 | Lifecycle |
| RFC-0005 | Event Bus Contract (Draft) |

---

# 9. Non-Goals

IntentCore SHALL NOT become:

- Message Broker
- Workflow Engine
- API Gateway
- Service Mesh
- Transport Stack
- Business Runtime

---

# 10. Implementation Philosophy

Specification-Driven Development

```
Architecture
    ↓
ADR
    ↓
RFC
    ↓
Contracts
    ↓
Implementation
    ↓
Testing
```

Implementation SHALL follow specifications.

Specifications SHALL NOT be derived from implementation.

---

# 11. Informative Roadmap

## Phase 1

Foundation

## Phase 2

Runnable Kernel

## Phase 3

Transformation

- Federation
- Distributed Coordination
- Semantic Routing
- Stability Foundation

## Phase 4

Knowledge Plane

- Intent Graph
- Zero-Trust
- Global Coordination

---

# 12. Summary

IntentCore is the architectural center of the system.

Lifecycle is the sole authority for state transitions.

Repository is the authoritative state store.

History is immutable.

Proof and Telemetry are emitted by Lifecycle.

ABTP remains outside the kernel boundary.

Every implementation SHALL preserve these architectural contracts.
