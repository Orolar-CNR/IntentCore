RFC-0001: Transport & Wire Protocol

Status: Frozen
Authors: IntentCore Architecture Team
Created: 2026-07-10
Updated: 2026-07-10
Dependencies: RFC-0000
Implements: IntentCore transport boundary via ABTP

1. Abstract

This RFC defines the canonical transport and wire contract for IntentCore. It specifies the "SemanticEnvelope" structure, validation requirements, normalization requirements, and wire-level invariants used to carry intent-bearing messages into the kernel.

2. Motivation

IntentCore requires a deterministic, language-agnostic, and versioned wire format that can transport semantic messages from external systems into the kernel without embedding business logic in the transport boundary. The wire contract must remain stable so that lifecycle, admission, and repository behavior can evolve independently.

3. Scope

This RFC defines:

- the canonical envelope structure
- required and optional fields
- validation requirements
- normalization requirements
- wire-level invariants
- compatibility rules

This RFC does not define admission policy, lifecycle transitions, repository internals, proof models, or telemetry semantics beyond transport-level metadata propagation.

4. Terminology

SemanticEnvelope
The canonical transport object carried by ABTP into IntentCore.

Wire Contract
The immutable field set, ordering, and semantics used for transport interoperability.

Validation
Structural and protocol checks performed before an envelope is admitted into the kernel pipeline.

Normalization
Deterministic transformation of an envelope into canonical form without altering its semantic intent.

5. Architecture

The transport flow is:

External Systems
  → ABTP
  → SemanticEnvelope
  → IntentCore Validation
  → IntentCore Normalization
  → IntentCore Admission

ABTP is a transport protocol only. It MUST NOT contain business logic, lifecycle logic, or state mutation logic.

6. Canonical Wire Contract

The canonical envelope is "SemanticEnvelope". The following fields are part of the frozen wire contract.

Field| Type| Requirement| Purpose
"envelope_id"| UUIDv4| Required| Idempotency and duplicate detection
"agent_identity"| String| Required| Source identity of the sender
"event_timestamp"| ISO 8601| Required| Logical ordering and traceability
"telemetry_class"| String| Optional| High-level observability classification
"opaque_payload"| JSON| Required| Payload forwarded without transport-side interpretation

Wire Contract Rules

- Field names and semantics MUST remain stable.
- The transport layer MUST treat "opaque_payload" as opaque.
- The transport layer MUST NOT interpret payload business meaning.
- "envelope_id" MUST be used as the idempotency key when applicable.
- "event_timestamp" MUST be preserved through the transport boundary.

7. Validation Pipeline

Validation is the first kernel-side step after transport ingestion.

7.1 Structural Validation

The system MUST verify that:

- all required fields are present
- field types conform to the contract
- envelope identity is structurally valid
- payload encoding is parseable

7.2 Protocol Validation

The system SHOULD verify transport-level compatibility such as:

- protocol version compatibility
- checksum or integrity markers, if present
- future framing metadata compatibility

7.3 Rejection Behavior

If validation fails:

- the envelope MUST NOT enter admission
- the failure SHOULD be recorded as telemetry
- the envelope MAY be rejected with a structured error
- the kernel MUST preserve the rejection reason for observability

8. Normalization Requirements

Normalization MUST be deterministic.

The system MAY:

- canonicalize field ordering for internal representation
- assign default values to optional metadata fields
- normalize timestamps into a canonical representation
- normalize routing metadata into an internal form

The system MUST NOT:

- alter payload meaning
- rewrite intent semantics
- mutate the semantic meaning of the original envelope

Normalization is a canonicalization step, not a transformation of intent.

9. Invariants

The following invariants MUST hold:

- Validation MUST NOT mutate the incoming envelope.
- Normalization MUST be deterministic for the same input.
- Transport MUST NOT access repository state directly.
- Transport MUST NOT decide admission.
- Transport MUST NOT change lifecycle state.
- The same "envelope_id" MUST map to the same logical message identity.
- A rejected envelope MUST remain rejected regardless of transport retries unless explicitly re-submitted.
- IntentCore MUST remain independent from any specific transport implementation.

10. Compatibility

This RFC is frozen. Future transport changes MUST preserve backward compatibility through a versioned envelope strategy.

Compatibility rules:

- Version changes MUST be explicit.
- Older envelopes SHOULD remain readable where possible.
- Breaking changes MUST require a new RFC revision.
- Transport evolution MAY occur at the edges, but the frozen core contract MUST remain stable.

11. Security Considerations

The transport layer MUST assume that incoming envelopes are untrusted.

The system SHOULD guard against:

- message injection
- malformed envelope attacks
- duplicate replay
- oversized payload abuse
- malformed timestamp attacks

The transport layer MUST NOT execute payload logic, tool calls, or policy decisions.

12. Reference Implementation Notes

The current reference implementation uses the transport boundary as a thin protocol layer. Implementations MAY use any transport substrate, provided they preserve this RFC’s frozen contract.

13. Future Work

Future versions MAY define:

- binary framing
- compression strategy
- session negotiation
- cross-shard routing metadata
- backpressure signaling

Those concerns MUST NOT change the frozen semantics of this RFC.
