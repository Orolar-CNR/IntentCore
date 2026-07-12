# IntentCore Architecture Landscape

> This document provides the architectural landscape of IntentCore.
>
> Document 1 defines the stable architectural foundation.
> Document 2 defines the long-term architectural vision.
>
> Together they describe the evolution path from a specification-driven
> coordination kernel into a distributed intent coordination infrastructure.

---

# Architecture Philosophy

IntentCore is **not a message broker**, **not a workflow engine**, and **not a transport protocol**.

IntentCore is an **Intent Coordination Kernel**.

Its responsibility is to coordinate, validate, govern, and maintain the lifecycle of Intent while remaining independent from transport implementations.

The architecture follows one fundamental principle:

> Every state mutation must originate from a validated Intent.

Nothing inside the system is allowed to modify state directly.

---

# Two-Layer Architecture

The project is intentionally described using two complementary documents.

## Layer 1 — Architectural Foundation

(Document 1)

This layer establishes the immutable architectural contracts.

Its purpose is to answer:

- What is Intent?
- What is Transport?
- What is Admission?
- What is Lifecycle?
- What is State?
- What is Governance?
- How are these components related?

The outcome of this layer is a stable architecture suitable for RFCs and production implementation.

Core characteristics include:

- Intent as the execution primitive
- ABTP as transport boundary
- SemanticEnvelope as canonical wire format
- Deterministic lifecycle
- Single Source of Truth
- Immutable history
- Strict dependency direction
- Separation of concerns

This layer defines **how the system is built**.

---

## Layer 2 — Architectural Vision

(Document 2)

The second document extends the foundation toward future distributed systems.

It answers a different question:

> What should IntentCore eventually become?

Instead of redefining the architecture, it expands its capabilities.

Future capabilities include:

- Semantic Routing
- Distributed Coordination
- Multi-Agent Governance
- Trust Infrastructure
- Knowledge Plane
- Intent Discovery
- Adaptive Coordination
- Global Observability

This layer defines **where the architecture is heading**.

---

# Relationship Between Both Documents

The two documents are not alternatives.

They represent different abstraction levels.

```
Document 1
    │
    │ establishes
    ▼

Architecture Foundation
    │
    │ enables
    ▼

Document 2
    │
    │ expands
    ▼

Distributed Intent Infrastructure
```

Document 1 provides architectural stability.

Document 2 provides architectural evolution.

---

# Core Architectural Components

## Intent

Intent is the fundamental semantic primitive.

It expresses desired outcomes instead of execution procedures.

Intent is the only object allowed to initiate state transitions.

---

## Transport Boundary (ABTP)

ABTP is responsible only for transporting data.

Responsibilities include:

- framing
- serialization
- checksum
- protocol validation
- version negotiation

ABTP never performs:

- lifecycle decisions
- policy evaluation
- state mutation
- governance logic

Transport remains completely independent from the coordination kernel.

---

## SemanticEnvelope

SemanticEnvelope is the canonical wire contract.

Every external system communicates with IntentCore using SemanticEnvelope.

It guarantees:

- interoperability
- deterministic parsing
- stable wire compatibility

---

## Admission

Admission evaluates whether an Intent is allowed to enter the system.

Responsibilities include:

- schema validation
- identity verification
- authorization
- policy enforcement
- trust evaluation

Admission does not mutate system state.

---

## Lifecycle

Lifecycle is the only authority allowed to perform state transitions.

Every transition must be:

- deterministic
- atomic
- auditable

Transition rules are defined by RFC-0004.

---

## Repository

Repository acts as the Single Source of Truth.

It guarantees:

- Compare-And-Swap
- Version control
- Snapshot
- Recovery
- Immutable state history

---

## History / Proof / Telemetry

Every transition generates evidence.

The system maintains:

- audit history
- proof records
- telemetry events
- distributed tracing

Observability is considered a first-class architectural component.

---

# Dependency Direction

IntentCore enforces a strict one-way dependency model.

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
History / Proof / Telemetry
```

No component may bypass another layer.

No outer layer may mutate inner state.

---

# IntentCore Roadmap — Phase 0: Specification & Blueprint

## วัตถุประสงค์

Phase นี้มุ่งเน้นการสร้าง **สัญญาสถาปัตยกรรม (Architectural Contracts)** และ Blueprint ที่มั่นคง ก่อนเริ่มพัฒนาโค้ดจริง โดยยึดหลัก **Specification First** เพื่อให้การพัฒนาในระยะต่อไปเป็นไปอย่างสอดคล้องและตรวจสอบได้

## Stage 1 — Core Architectural Contracts

**เป้าหมาย**
กำหนดขอบเขตของระบบให้สมบูรณ์และไม่คลุมเครือ

**งานที่ต้องดำเนินการ**
- `RFC-0000` — Architectural Principles (กฎสูงสุดของทั้งชุด RFC)
- `RFC-0001` — SemanticEnvelope
- `RFC-0002` — Admission Boundary
- `RFC-0003` — State Repository
- `RFC-0004` — Lifecycle Control

**ผลลัพธ์**
- Stable Core Contracts
- Frozen Public Interfaces
- Canonical Architecture
- Dependency Direction ถูกกำหนดอย่างถาวร

## Stage 2 — Repository Blueprint

**เป้าหมาย**
กำหนดพฤติกรรมของระบบจัดเก็บข้อมูลระยะยาว

**งานที่ต้องดำเนินการ**
- Snapshot Strategy
- Incremental Snapshot
- Recovery Procedure
- Ledger Compaction
- Archive Strategy
- Retention Policy

**ผลลัพธ์**
Repository สามารถ:
- Recover ได้
- Scale ได้
- Audit ได้
- รองรับ Append-only Ledger ระยะยาว

## Stage 3 — Transport Blueprint (ABTP)

**เป้าหมาย**
กำหนดขอบเขตของ Transport อย่างสมบูรณ์

**งานที่ต้องดำเนินการ**
- Transport Pipeline
- Decoder
- Encoder
- Version Negotiation
- Frame Validation
- Checksum
- Packet Drop Rules
- Malformed Envelope Handling
- Transport Error Model

**สิ่งที่ "ไม่อยู่" ใน ABTP**
- Policy
- Lifecycle
- State Mutation
- Repository
- Governance
- Semantic Evaluation

**ผลลัพธ์**
ABTP กลายเป็น Pure Transport Boundary อย่างแท้จริง

## Stage 4 — Governance Blueprint

**เป้าหมาย**
สร้างระบบกำกับดูแลที่พิสูจน์ได้

**งานที่ต้องดำเนินการ**
- Policy Language
- Logic Linter
- Decidability Checking
- Authorization
- Trust Evaluation
- Credential Validation

**ผลลัพธ์**
Admission Layer สามารถรับประกันว่า:
- Policy ไม่มี Infinite Loop
- Policy วิเคราะห์จบได้
- Governance ตรวจสอบย้อนหลังได้

---

# สิ่งที่ยังไม่ทำใน Phase นี้

เพื่อรักษาขอบเขตของโครงการ สิ่งต่อไปนี้จะยังไม่อยู่ใน Phase 0:
- Network Optimization
- eBPF
- XDP
- AF_XDP
- DPDK
- RDMA
- NUMA Scheduling
- Zero-copy Optimization

**เหตุผลคือ:**
สิ่งเหล่านี้เป็นเพียง **Implementation Technology** ไม่ใช่ **Architectural Contract** การนำเข้ามาพิจารณาในระยะนี้อาจทำให้เกิดการเบี่ยงเบนทางสถาปัตยกรรม (Architectural Drift) ได้ จึงต้องผลักการตัดสินใจเหล่านี้ไปกระทำหลังจากที่ RFC, Wire Protocol, และ Transport Contract ถูกล็อกอย่างสมบูรณ์แล้วเท่านั้น

---

# หลังจาก Phase 0

เมื่อ Blueprint สมบูรณ์แล้ว จึงเข้าสู่ **Phase 1: Implementation**

```text
Core Package
    ↓
Lifecycle
    ↓
Repository
    ↓
Transport
    ↓
Telemetry
    ↓
Testing
```

จากนั้นจึงเข้าสู่ **Phase 2: Distributed Expansion** เช่น:
- Semantic Routing
- Agent Discovery
- Federation
- Intent Graph
- Distributed Scheduling
- Knowledge Plane

---

# Architectural Principles (RFC-0000)

RFC-0000 จะไม่อธิบาย API หรือโปรโตคอล แต่จะอธิบาย "กฎสูงสุด" (Supreme Rule) ของทั้งชุด RFC ที่ทุกฉบับต้องปฏิบัติตามอย่างเคร่งครัด เช่น:
- Intent First
- Specification First
- One-Way Dependency
- Separation of Concerns
- Transport Independence
- Deterministic Lifecycle
- Single Source of Truth
- Immutable History
- Governance Before Execution
- Observability by Design

ข้อดีคือ เมื่อโครงการขยายไปเป็น RFC-0010, RFC-0020 หรือมากกว่านั้น ทุกเอกสารจะยังยึดหลักการเดียวกัน และช่วยป้องกันการออกแบบที่เบี่ยงเบนจากสถาปัตยกรรมหลักโดยไม่ตั้งใจ

# Summary

Document 1 defines the architecture.

Document 2 defines the direction.

Together they establish the complete evolution path of IntentCore—from a specification-driven coordination kernel to a distributed infrastructure for intent-aware autonomous systems.
