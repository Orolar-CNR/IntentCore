RFC-0002: Intent Admission Interface

Status: Frozen
Authors: IntentCore Architecture Team
Created: 2026-07-10
Updated: 2026-07-10
Dependencies: RFC-0000, RFC-0001
Implements: Admission boundary for intent entry into IntentCore

1. Abstract

This RFC defines the admission interface used by IntentCore to determine whether an incoming intent may enter the kernel pipeline. Admission is a policy boundary, not an execution boundary.

2. Motivation

IntentCore requires a deterministic and pluggable policy decision layer so that lifecycle control, repository mutation, and execution coordination remain separated from transport concerns. Admission must be explicit, testable, and stable.

3. Scope

This RFC defines:

- the "AdmissionPolicy" contract
- admission decision boundaries
- evaluation requirements
- rejection behavior
- policy invariants
- integration rules with validation and normalization

This RFC does not define lifecycle transitions, repository storage layout, proof construction, or execution runtime behavior.

4. Terminology

AdmissionPolicy
A pluggable policy interface used to determine whether an intent may proceed into the kernel.

Admission Decision
The result of policy evaluation: allow or reject.

Intent
A semantic request or directive that has passed transport validation and normalization and is being considered for kernel admission.

Policy Boundary
The architectural layer where access to the kernel is granted or denied.

5. Architecture

The admission flow is:

SemanticEnvelope
  → Validation
  → Normalization
  → AdmissionPolicy
  → Allow or Reject

Admission MUST occur after validation and normalization, and before lifecycle transition or repository mutation.

6. Contract

The canonical contract is an interface equivalent to:

AdmissionPolicy {
  Evaluate(ctx, intentID, state) -> (allowed, error)
}

Requirements

- AdmissionPolicy MUST be pluggable.
- AdmissionPolicy MUST NOT directly mutate repository state.
- AdmissionPolicy MUST be able to return allow/reject decisions.
- AdmissionPolicy MAY return structured errors for policy failures.
- AdmissionPolicy SHOULD be deterministic given the same inputs, configuration, and state snapshot.

7. Admission Flow

An incoming intent MUST follow this order:

1. Transport ingestion
2. Validation
3. Normalization
4. Admission evaluation
5. Lifecycle entry, if admitted

If admission is rejected:

- the intent MUST NOT enter lifecycle execution
- the intent MUST NOT mutate the repository
- the rejection SHOULD be emitted to telemetry
- the rejection SHOULD be traceable to an intent identity and reason

8. Integration Rules

Admission is dependent on validated and normalized intent input.

Admission MUST NOT:

- bypass validation
- bypass normalization
- mutate lifecycle state directly
- update the repository directly
- create proof records directly
- alter transport metadata as a side effect

Admission MAY:

- inspect intent identity
- inspect routing metadata
- inspect policy domain
- inspect trust state or resource feasibility where available through read-only context
- emit a reason for rejection

9. Invariants

The following invariants MUST hold:

- Admission MUST occur before lifecycle entry.
- Rejected intents MUST NOT transition state.
- Admission MUST be isolated from transport concerns.
- Admission MUST be isolated from repository mutation.
- Admission decisions MUST be logged or emitted for observability.
- Admission policy SHOULD remain pluggable across implementations.

10. Error Handling

Policy failures MUST be represented as structured errors.

The system SHOULD distinguish between:

- validation failure
- policy rejection
- internal evaluation failure
- missing context
- unsupported policy domain

A policy rejection is not necessarily an internal error. It is a valid decision outcome.

11. Compatibility

This RFC is frozen. Future admission evolution MUST preserve the interface-first contract and the separation between policy evaluation and state mutation.

Compatibility rules:

- New policy implementations MAY be added without breaking the contract.
- Additional policy metadata MAY be introduced if backward compatible.
- Breaking changes MUST require a new RFC revision.

12. Security Considerations

Admission is a security-sensitive boundary.

The system SHOULD protect against:

- unauthorized intent entry
- policy bypass
- malformed metadata
- spoofed identities
- privilege escalation through policy plugins

AdmissionPolicy implementations MUST be treated as privileged components.

13. Reference Implementation Notes

A reference implementation MAY use lightweight in-memory policy evaluation for prototype purposes, provided the interface contract is preserved and no repository mutation occurs during policy evaluation.

14. Future Work

Future revisions MAY define:

- policy composition
- trust scoring integration
- resource feasibility checks
- namespace-based policy isolation
- distributed policy evaluation

These additions MUST preserve the frozen admission boundary defined by this RFC.
