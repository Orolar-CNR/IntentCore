# RFC-0001 — Semantic Envelope

**Status:** Locked (Addendum Draft Merged)
**Category:** Architecture Specification
**Version:** 1.0
**Last Updated:** 2026-07-11

## Abstract

This RFC defines the canonical data structure, the `SemanticEnvelope`, which serves as the universal wire contract for IntentCore. It specifies the strict normative requirements for the envelope's validation, schema conformance, and normalization processes before it crosses the admission boundary into the kernel.

## Motivation

IntentCore requires a deterministic, language-agnostic, and versioned wire format that can transport semantic messages from external systems into the kernel without embedding business logic in the transport boundary. The wire contract must remain absolutely stable so that lifecycle, admission, and repository behaviors can evolve independently without breaking communication with remote agents.

## Terminology

*   **SemanticEnvelope:** The canonical transport object carried into IntentCore.
*   **Validation:** Structural and protocol checks performed before an envelope is admitted into the kernel pipeline.
*   **Normalization:** Deterministic transformation of an envelope into canonical form without altering its semantic intent.

## Architectural Context

The `SemanticEnvelope` acts as the interface between the Transport layer (which handles bytes) and the Validation/Admission layers (which handle structured intent). Transport MUST deliver a valid byte stream that decodes into a `SemanticEnvelope`.

## Normative Specification

### 1. Canonical Schema
The data structure of the `SemanticEnvelope` MUST strictly conform to the schema defined in this specification. IntentCore implementations MUST NOT redefine or diverge from this schema.

The `SemanticEnvelope` MUST contain the following fields:
*   `envelope_id`: A valid UUIDv4 string. REQUIRED. Used for idempotency.
*   `agent_identity`: A string identifying the source of the sender. REQUIRED.
*   `event_timestamp`: A string formatted as ISO 8601. REQUIRED. Used for logical ordering.
*   `telemetry_class`: A string. OPTIONAL.
*   `opaque_payload`: A JSON-encoded payload. REQUIRED.
*   `intent_id`: A stable, globally unique identity for the Intent. REQUIRED. Used for duplication / replay detection.
*   `expected_version`: The version of the target state the Intent was authored against. REQUIRED. Used for CAS conflict detection.

### 2. Validation
Upon receiving a `SemanticEnvelope`, the Validation layer MUST evaluate the structural integrity of the payload.
*   The Validation layer MUST act solely as a deterministic structural conformance check.
*   The Validation layer SHALL NOT inject any business logic.
*   Incoming envelopes with missing required fields or malformed types MUST be immediately rejected.
*   Schema version mismatches MUST result in immediate rejection before admission.

### 3. Normalization
Normalization is a canonicalization step.
*   The Normalization step MAY canonicalize field ordering and default values for optional metadata.
*   The Normalization step MUST NOT alter the payload meaning or rewrite intent semantics.

### 4. Kernel Identity and Mutation Fields
#### MUST
The SemanticEnvelope MUST include the following two kernel-owned fields, both present simultaneously. Neither substitutes for the other, as they guard against distinct failure modes:
*   `intent_id` — a stable, globally unique identity for the Intent. Used for duplication / replay detection.
*   `expected_version` — the version of the target state the Intent was authored against. Used for CAS conflict detection at Repository commit time.

`mutation_key` is explicitly not part of the SemanticEnvelope contract. Only `intent_id` and `expected_version` serve these roles.

The SemanticEnvelope MUST remain strictly decoupled from any event-plane or transport-specific envelope (e.g., the EventEnvelope defined in RFC-0005). If an Intent arrives wrapped inside an EventEnvelope (e.g., via bus delivery or DLQ replay), the transport/adapter boundary MUST unwrap or re-materialize it into a standalone SemanticEnvelope, carrying `intent_id` and `expected_version`, before it reaches Admission. Admission MUST NOT receive or inspect EventEnvelope fields directly.

#### MUST NOT
*   `idempotency_key` (an Event Bus / RFC-0005 concept) MUST NOT appear in, nor be conflated with, the SemanticEnvelope`'`s `intent_id`.

## State Model

The validation process follows a stateless model:
Received → Validated → Normalized → (Passed to Admission)

## Interfaces

```go
type EnvelopeValidator interface {
    Validate(envelope SemanticEnvelope) error
}
```

## Error Model

Validation failures MUST result in a definitive `ValidationError` containing the specific structural invariant that was violated. The envelope MUST NOT proceed to Admission.

## Security Considerations

The `SemanticEnvelope` is the primary attack surface from external actors. Strict structural validation prevents malformed data injection, oversized payload attacks, and replay attacks (via strict UUIDv4 checking). Transport and Validation layers SHALL NOT execute payload logic.

## Observability

Validation failures MUST log the rejection reason, specifically including version variances for system telemetry, without executing any further parsing.

## Compliance Requirements

| Requirement | RFC | Test |
| :--- | :--- | :--- |
| Enforce UUIDv4 on envelope_id | RFC-0001 | envelope_validation_test.go |
| Reject invalid ISO8601 | RFC-0001 | envelope_validation_test.go |
| Opaque Payload unaltered | RFC-0001 | envelope_normalization_test.go |

## Backward Compatibility

This RFC is frozen. Future changes to the `SemanticEnvelope` MUST preserve backward compatibility through a versioned envelope strategy. Breaking changes require a new RFC revision.

## Rationale

Separating the wire format (SemanticEnvelope) from the transport mechanism (eBPF/ABTP/TCP) ensures that the core kernel remains agnostic to network details, focusing entirely on intent coordination.

## References

*   RFC-0000 — Architectural Principles
