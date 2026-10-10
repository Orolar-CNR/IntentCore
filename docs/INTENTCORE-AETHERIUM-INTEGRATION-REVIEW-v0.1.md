# IntentCore → Aetherium-Manifest Integration Plan Review & Feedback Report
**Document Version:** v0.1
**Status:** Architectural Review & Recommendations
**Target Proposal:** `docs/INTENTCORE-AETHERIUM-INTEGRATION-PLAN-v0.1.md`
**Reference Roadmap:** `docs/blueprints/phase-3-roadmap.md`

---

## 1. Executive Summary

This report provides a comprehensive architectural and technical evaluation of the **IntentCore → Aetherium-Manifest Integration Proposal (v0.1)**.

The evaluation confirms that the core architectural boundary proposed in `v0.1`—positioning IntentCore strictly as the admission, authority, and coordination kernel while delegating visual interpretation and rendering to Aetherium—is **sound, compliant with canonical RFCs (RFC-0000 through RFC-0004), and essential for maintaining strict layer separation**.

However, this review highlights key gaps and recommended refinements across five dimensions:
1. **Architectural Boundary & Invariants**
2. **Phase Alignment & Development Sequence** (Conflict between Roadmap Phase 3.1 and Proposal prerequisites)
3. **IntentCore Kernel Items (IC-001 through IC-007)**
4. **Security, Trust Boundary & Cryptographic Authentication**
5. **Integration Contract, Schema & Error Dynamics**

---

## 2. Evaluation Matrix & Key Findings

| Domain | Proposal Provision | Assessment | Recommendation |
| :--- | :--- | :--- | :--- |
| **Architectural Boundary** | No Canvas2D/WebGPU/SDF/Particle concepts inside IntentCore kernel | **EXCELLENT** | Enforce strict payload opaqueness; prevent leaking UI/Manifestation logic into IntentCore. |
| **Development Sequence** | Requires Zero-Trust cryptographic verification (IC-002) prior to Phase-0/1 integration | **MISALIGNED** | Roadmap `phase-3-roadmap.md` places Zero-Trust Primitives in **Phase 3.3**, whereas Phase 3.1 focuses on **Stability Foundation**. Decouple interface readiness from full cryptographic zero-trust rollouts. |
| **Payload Typing (IC-001)** | Introduces application-level payload metadata conventions | **GOOD WITH WARNING** | Preserve `OpaquePayload []byte` in kernel state. Metadata (`content_type`, `payload_version`) should exist in envelope attributes or header specs, avoiding kernel dependencies on external schemas. |
| **Signature Verification (IC-002)** | Notes current `len(Signatures) == 0` is placeholder check | **CRITICAL GAP** | Specify public key resolution, signature algorithms (e.g., Ed25519/ECDSA), and canonical serialization (canonical JSON/Protobuf) explicitly. |
| **Validation Enforcements (IC-003)** | Recommends UUID, timestamp, schema validations in `ValidateEnvelope()` | **VERY GOOD** | Implement deterministic validation rules in `admission/policy.go` rather than mutating kernel data structures. |
| **Trace & ID Mapping (IC-004 & IC-005)** | Maps `TraceID` and `IntentID`/`EnvelopeID` cleanly | **EXCELLENT** | Maintain clear domain identity separation: `IntentCore.TraceID -> Aetherium.trace_id`, `IntentCore.EnvelopeID -> Aetherium.external_message_id`. |
| **Rejection Feedback (IC-006)** | Machine-readable rejection codes | **VERY GOOD** | Align error/rejection codes with `contracts.RejectionCode` (RFC-0002). |
| **Lifecycle Mapping (IC-007)** | Differentiates IntentCore lifecycle state vs Aetherium manifestation state | **CRITICAL NECESSITY** | Formally define mapping state machine so `Pending` or `Admitted` state in IntentCore is not confused with active visual rendering. |

---

## 3. Detailed Review & Technical Recommendations

### 3.1 Alignment with Phase 3 Roadmap (`phase-3-roadmap.md`)

* **Observation:** The integration proposal states that IC-002 (Cryptographic Signature Verification) must be completed before production Aetherium integration.
* **Roadmap Analysis:** `docs/blueprints/phase-3-roadmap.md` mandates that **Phase 3.1** is strictly limited to **Stability Foundation** (snapshots, archiving, memory/CPU regression checks, and load simulation), while **Zero-Trust Primitives** (identity verification, capability authorization, signed event verification) are scheduled under **Phase 3.3**.
* **Recommendation:**
  1. Define a Pluggable `SignatureVerifier` Interface during Phase 3.1 without introducing complex distributed PKI or breaking Phase 3.1 focus.
  2. Implement mock or basic cryptographic verification adapters for early Phase 0/1 testing, postponing full distributed key resolution to Phase 3.3.

### 3.2 Detailed Analysis of IntentCore Kernel Proposals (IC-001 to IC-007)

#### IC-001: Application Payload Typing
* **Current Code:** `contracts.SemanticEnvelope` contains `OpaquePayload []byte`.
* **Feedback:** Keep `OpaquePayload` as raw `[]byte` within the kernel storage and memory layer. If typing is needed, introduce `ContentType string` and `PayloadVersion string` as optional metadata fields on `SemanticEnvelope` without giving `IntentCore` code any knowledge of Aetherium-specific JSON structures.

#### IC-002: Signature Verification
* **Current Code:** `admission/policy.go` only checks `len(env.Signatures) == 0`.
* **Feedback:** This is indeed a simulation check. A formal specification must define:
  1. Canonical byte representation of `SemanticEnvelope` prior to hashing/signing.
  2. Supported signature algorithms (e.g. Ed25519).
  3. Rejection code mapping (`CodeInvalidSignature`).

#### IC-003: Payload Schema Validation
* **Current Code:** `runtime/validation.go` unmarshals JSON but has commented-out structural checks.
* **Feedback:** `ValidateEnvelope` should perform deterministic structural checks (valid UUIDv4, non-zero timestamp, valid UTF-8 if text-based). Application-level schema validation should be delegated to specialized `AdmissionPolicy` implementations.

#### IC-004 & IC-005: Trace Propagation & Correlation ID Mapping
* **Current Code:** `core.TraceID` and `core.IntentID` exist in `core/types.go`.
* **Feedback:** The proposal's identity mapping is clean and prevents layer conflation:
  * `IntentCore.EnvelopeID` $\rightarrow$ `Aetherium.external_message_id`
  * `IntentCore.TraceID` $\rightarrow$ `Aetherium.trace_id`
  * `Aetherium.signal_id` and `Aetherium.episode_id` stay strictly within Aetherium boundaries.

#### IC-006: Explicit Rejection Feedback
* **Feedback:** Ensure all rejection reasons from `admission/policy.go` map directly to standardized machine-readable responses:
  * `CodeInvalidVersion`
  * `CodeMissingIntentID`
  * `CodeMissingIdentity`
  * `CodeMissingSignature`
  * `CodeInvalidSignature`
  * `CodeMissingPayload`
  * `CodeInvalidTimestamp`

#### IC-007: Lifecycle Semantics Mapping
* **Current Code:** `lifecycle/machine.go` handles transitions (`StateUnspecified`, `StateSubmitted`, `StateAdmitted`, `StateExecuting`, `StateCompleted`, etc.).
* **Feedback:** Aetherium manifestation states (e.g., `RECEIVING`, `INTERPRETING`, `MANIFESTING`, `FADING`) must never overwrite or replace `IntentCore` authority states. An explicit mapping matrix should be documented.

---

## 4. Integration Contract & Security Boundary Recommendations

### 4.1 Security & Zero-Trust Boundary
1. **Untrusted Ingress:** AI model outputs must always be treated as untrusted data inputs entering through the IntentCore Admission Boundary.
2. **Resource Isolation:** External AI models must never have direct GPU execution, WebGPU shader injection, or DOM access. All rendering instructions must be synthesized downstream by Aetherium based on accepted IntentCore intents.

### 4.2 Cross-Repository Test Matrix Additions
To ensure robust integration, the test matrix in Section 12 of the proposal should be extended with:
* **Malformed UTF-8 Payload Test:** Ensures parser fails safely without panicking.
* **Clock Skew Test:** Ensures timestamps far in the future or past are rejected cleanly according to admission policy.
* **Replay Protection Test:** Ensures re-submitting an identical `EnvelopeID` is handled deterministically via CAS / Idempotency policies in `state/cas.go`.

---

## 5. Summary of Recommended Next Actions

1. **Adopt Proposal v0.1 with Architectural Refinements** outlined in this document.
2. **Align Development Timeline:** Align IC-002 (Cryptographic Verification) with Roadmap Phase 3.3 while providing interface stubs for Phase 3.1 stability testing.
3. **Formalize Envelope Canonical Serialization Spec:** Document how `SemanticEnvelope` is serialized before hashing/signing.
4. **Maintain Strict Boundary:** Keep IntentCore devoid of graphics/presentation concepts.
