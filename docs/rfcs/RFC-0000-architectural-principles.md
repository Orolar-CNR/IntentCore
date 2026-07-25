# RFC-0000: Architectural Principles

## 1. Introduction

This document establishes the **supreme architectural rules** of IntentCore. All subsequent RFCs, blueprints, and implementations MUST strictly adhere to the principles defined herein. Any deviation from these principles constitutes an architectural violation and MUST NOT be admitted into the codebase.

## 2. Supreme Rules

### 2.1 Specification First
IntentCore is a specification-driven system.
* Implementations MUST NOT proceed without formal RFC approval.
* All external contracts (e.g., SemanticEnvelope, Lifecycle Matrix) MUST be treated as frozen unless a formal RFC revision is adopted.

### 2.2 Intent First
* State mutations MUST originate exclusively from a formally validated Intent.
* Systems or components MUST NOT bypass the Intent lifecycle to modify the canonical state.

### 2.3 Strict One-Way Dependency
The system MUST enforce a strictly unidirectional pipeline for data flow and dependency management.
* The pipeline MUST flow sequentially: External Systems -> ABTP -> SemanticEnvelope -> Validation -> Normalization -> Admission -> Lifecycle -> Repository (CAS), while History / Proof / Telemetry are emitted by Lifecycle.
* Components in one layer MUST NOT import or depend on packages or structures from layers deeper in the pipeline (e.g., `contracts` MUST NOT depend on `core` implementation details).
* Cyclic dependencies are strictly prohibited.

### 2.4 Separation of Concerns
* Transport (ABTP) MUST handle data delivery and serialization only. It MUST NOT evaluate semantics, apply governance, or mutate state.
* Admission MUST evaluate policies and govern entry. It MUST NOT mutate state.
* Lifecycle MUST coordinate state transitions. It MUST NOT perform transport functions or mutate state outside the Repository CAS constraint.
* Repository MUST provide deterministic CAS and append-only immutability. It MUST NOT perform governance or lifecycle coordination.

### 2.5 Deterministic Single Source of Truth
* The Repository MUST act as the absolute and sole Single Source of Truth.
* All state mutations MUST be performed using Compare-And-Swap (CAS) against a monotonically increasing version.
* The historical ledger MUST be append-only and completely immutable. Past records SHALL NOT be rewritten.
