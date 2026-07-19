# Platform Adapter Architecture

**Status:** Draft
**Category:** Blueprint
**Purpose:** Informative architectural reference
**Scope:** Platform adapters, host integration, SemanticEnvelope translation, and one-way pipeline boundaries

---

## 1. Purpose

This document defines the canonical architecture for platform adapters in IntentCore.

A platform adapter is an outer integration shell that translates host- or platform-specific events into `SemanticEnvelope` objects and forwards them into the IntentCore pipeline.

Platform adapters exist to connect external environments to IntentCore without redefining kernel semantics, bypassing pipeline stages, or mutating authoritative state directly.

This document is informative and architectural in nature. It does not change the meaning of RFC-0001 through RFC-0004.

---

## 2. Architectural Position

Platform adapters live outside the IntentCore kernel boundary.

They are not part of:

- Validation
- Normalization
- Admission
- Lifecycle
- Repository
- History
- Proof
- Telemetry

Instead, they are responsible for receiving external inputs from a host environment and translating those inputs into the canonical envelope used by the runtime pipeline.

---

## 3. Architectural Identity

IntentCore is the coordination kernel.

Platform adapters are the outer integration boundary.

The adapter layer MAY vary by platform, runtime, or deployment target, but the kernel contract MUST remain stable.

The adapter layer MUST NOT become a second kernel.

The adapter layer MUST NOT contain business logic, lifecycle logic, or repository logic.

---

## 4. Normative Design Principles

The platform adapter layer MUST follow these rules:

1. Platform adapters SHALL be external integration shells that translate host/platform-specific events into `SemanticEnvelope` objects and forward them into the IntentCore pipeline.
2. Platform adapters MUST NOT bypass validation, normalization, admission, or lifecycle stages.
3. Platform adapters MUST NOT bypass dispatcher control.
4. Platform adapters MUST NOT mutate repository state directly.
5. Platform adapters MUST NOT write to history, proof, or telemetry stores directly.
6. Platform adapters MUST remain implementation-specific and kernel-independent.

---

## 5. System Boundary

The platform adapter boundary sits between the host environment and the IntentCore runtime pipeline.

### Boundary Flow

```text
External Host / Platform
        │
        ▼
Platform Adapter
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
Dispatcher
        │
        ▼
Lifecycle
        │
        ▼
Repository
        │
        ├── History
        ├── Proof
        └── Telemetry
```
