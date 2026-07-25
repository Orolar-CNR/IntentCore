# Canonical Glossary

This document defines the canonical terminology used throughout the IntentCore architecture, RFCs, Blueprints, and implementation. Standardizing terms ensures clear communication and prevents semantic drift.

## A
*   **ABTP (AetherBus Transport Protocol):** The strictly defined transport boundary. Responsible solely for data delivery, framing, and serialization. It does *not* evaluate semantics or mutate state.
*   **Admission:** The architectural layer responsible for governing entry into the system by evaluating deterministic policies and validating identities.
*   **Agent:** An external entity (system, user, or service) that submits an Intent to the system.

## C
*   **CAS (Compare-And-Swap):** A concurrency control mechanism used exclusively by the State Repository to ensure atomic state mutations based on version matching.

## I
*   **Intent:** The fundamental semantic primitive representing a desired outcome rather than an explicit execution procedure. It is the sole entity allowed to initiate state transitions.
*   **IntentCore:** The coordination kernel. The central component responsible for lifecycle management, state consistency, and policy enforcement.

## L
*   **Lifecycle:** The definitive state machine governing the allowed transitions of an Intent. It coordinates mutations but does not store state itself.

## O
*   **One-Way Dependency:** The strict architectural rule dictating that data and execution flow in a single, irreversible pipeline (e.g., Transport -> Validation -> Admission -> Lifecycle -> Repository).

## R
*   **Repository (State Repository):** The isolated, canonical single source of truth for the system's authoritative state, functioning as an append-only ledger.
*   **RFC (Request for Comments):** The normative, mandatory contracts that define the unchangeable rules of the system (using MUST/SHALL language).

## S
*   **SemanticEnvelope:** The canonical wire format used for communication between external systems and IntentCore.
*   **Specification-First:** The guiding principle that architectural contracts and RFCs must be formalized and approved before implementation begins.
