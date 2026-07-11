# RFC-0002 — Admission Boundary

**Status:** Locked
**Category:** Architecture Specification
**Version:** 1.0
**Last Updated:** 2026-07-11

## Abstract

This RFC defines the normative requirements for the Admission Boundary in IntentCore. The Admission layer is the governance gatekeeper responsible for determining if a validated `SemanticEnvelope` is authorized and trusted enough to enter the execution Lifecycle.

## Motivation

To ensure that only valid, authorized, and policy-compliant Intents induce state transitions, a strict governance boundary must exist before execution. Without this boundary, malicious or malformed intents could reach the State Machine and corrupt the Repository.

## Terminology

*   **Admission Pipeline:** The deterministic sequence of governance evaluations.
*   **Policy:** A bounded, decidable set of rules evaluating trust, capability, and authority.

## Architectural Context

Admission sits immediately after Normalization and immediately before the Lifecycle.
`Normalization → Admission → Lifecycle`

## Normative Specification

### 1. Formal Admission Pipeline
The Admission process MUST be executed as a strict, deterministic pipeline:
`Authenticate → Authorize → Policy Evaluation → Trust Evaluation → Accept/Reject`

### 2. Determinism and Decidability
All policies evaluated in the Admission layer MUST be strictly decidable.
*   Policies SHALL NOT contain unbounded loops or recursion.
*   Policies SHALL NOT execute arbitrary user-defined functions or code.
*   Admission evaluation MUST guarantee termination within a fixed computational bound.

### 3. Independence
The Admission layer SHALL NEVER mutate Intent state or repository state. Its sole output is a deterministic `Accept` or `Reject` decision, accompanied by audit evidence.

### 4. Rejection
If any stage in the Admission Pipeline fails, the Intent MUST be rejected immediately. The system MUST NOT attempt partial evaluation or fallback authorization mechanisms.

## State Model

The Admission Pipeline is a stateless decision gate yielding one of two terminal outcomes:
*   `AdmissionAccepted`
*   `AdmissionRejected`

## Interfaces

```go
type AdmissionPolicy interface {
    Evaluate(envelope SemanticEnvelope) (AdmissionDecision, error)
}
```

## Error Model

A failure in policy evaluation, or a missing identity, MUST trigger an `AdmissionError`. An `AdmissionError` explicitly prevents the intent from entering the Lifecycle.

## Security Considerations

The Admission Boundary is the primary governance defense. By restricting the policy language to a decidable DSL (YAML/JSON) and prohibiting Turing-complete execution, the system is protected against denial-of-service via infinite loops or logic bombs.

## Observability

Every admission decision MUST generate a verifiable decision record, including the policy identifier used, trust evidence, a timestamp, and verifier information.

## Compliance Requirements

| Requirement | RFC | Test |
| :--- | :--- | :--- |
| Reject recursive policies | RFC-0002 | policy_linter_test.go |
| Generate Audit Record on Reject | RFC-0002 | admission_audit_test.go |
| Guarantee immutable state during Admission | RFC-0002 | admission_isolation_test.go |

## Backward Compatibility

Changes to the Admission Pipeline sequence require a new RFC revision. Policy schema changes must be versioned.

## Rationale

Governance precedes execution. A purely declarative, decidable policy engine ensures that the system can mathematically prove whether an Intent was authorized without running the risk of halting problems or side effects.

## References

*   RFC-0000 — Architectural Principles
*   RFC-0001 — Semantic Envelope
