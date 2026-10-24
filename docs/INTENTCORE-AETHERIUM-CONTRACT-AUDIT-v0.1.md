# IntentCore ↔ Aetherium Integration Contract Audit v0.1

**Audit Date:** October 2026
**Repository:** `Orolar-CNR/IntentCore`
**Branch:** `jules-698959085642890864-3cd68eda`
**HEAD Commit:** `986dca4df4bb8c1df4a5fe5d21daf0e7cbefb2eb`
**Audit Scope:** Read-Only Architectural & Implementation Readiness Audit

---

## 1. Executive Summary

This report presents a read-only architectural and implementation audit evaluating **IntentCore**'s readiness for integration with **Aetherium-Manifest**. The primary objective is to verify whether IntentCore's existing contracts, admission boundaries, lifecycle state machines, and transport wrappers support integration while maintaining strict architectural boundaries:

1. **IntentCore** strictly retains exclusive authority over Admission, Identity, Cryptographic Verification, Lifecycle State Machine semantics, and State Persistence.
2. **Aetherium** retains exclusive ownership over Multimodal Interpretation, Modality-Specific Processing, and Manifestation/Rendering (Canvas2D, WebGPU, SDF, Particles).
3. Multimodal data crosses the boundary strictly as an **Opaque Payload** (`OpaquePayload []byte`).

### Key Audit Highlights
- **Architecture Integrity:** The core design and separation of responsibilities between IntentCore and Aetherium are sound and fully aligned with RFC-0000 through RFC-0004. IntentCore contains zero dependencies on graphics, WebGPU, SDF, or rendering logic.
- **Critical Security Finding (P0):** Cryptographic signature verification in `admission/policy.go` is currently a presence check (`len(env.Signatures) == 0`). It does not perform cryptographic signature verification, public-key resolution, or canonical byte hashing.
- **Contract & Traceability Gaps (P1):** `SemanticEnvelope` lacks the normative `expected_version` field mandated by RFC-0001, and `runtime/dispatcher.go` fails to propagate `TraceID` into transition requests, causing trace correlation breakdown in history recording.
- **Test Baseline:** All 18 pure Go unit tests across `proof`, `runtime`, `state`, `telemetry`, and `tests` pass 100%. Build failures in `transport/abtp` and `cmd/intentcored` stem strictly from missing toolchain eBPF dependencies (`llvm-strip` for `bpf2go`).

---

## 2. Repository and Test Baseline

### 2.1 Toolchain & Environment Baseline
- **Go Version:** `1.25.10` (`linux/amd64`)
- **Go Module:** `github.com/Orolar-CNR/IntentCore`
- **Direct Dependencies (`go.mod`):**
  - `github.com/google/uuid v1.6.0`
  - `github.com/cilium/ebpf v0.22.0`
  - `golang.org/x/sys v0.43.0` (indirect)

### 2.2 Execution Results (`go test ./...`)

#### Pure Go Kernel Test Suites (100% PASS)
| Package | Status | Duration | Tests Executed |
| :--- | :--- | :--- | :--- |
| `github.com/Orolar-CNR/IntentCore/proof` | **PASS** | 0.014s | `TestInMemoryProofRecorder` |
| `github.com/Orolar-CNR/IntentCore/runtime` | **PASS** | 0.310s | `TestDispatcher`, `TestPipelineProof`, `TestPipelineTelemetry`, `TestPipelineTelemetryReject`, `TestVerticalSlice`, `TestPipelineTimeoutAndCancellation` |
| `github.com/Orolar-CNR/IntentCore/state` | **PASS** | 0.007s | `TestInMemoryArchiveStore`, `TestDefaultSnapshotScheduler`, `TestDefaultRetentionPolicy`, `TestInMemorySnapshotStore_SaveLoad`, `TestRepository_Recover_ByID_NotLatest`, `TestRepository_Recover_ChecksumMismatch`, `TestRepository_Recover_NotFound`, `TestRepository_Recover_EmptyChecksum_SkipsVerification`, `TestRepository_SnapshotAndRecover` |
| `github.com/Orolar-CNR/IntentCore/telemetry` | **PASS** | 0.006s | `TestInMemoryRecorder` |
| `github.com/Orolar-CNR/IntentCore/tests` | **PASS** | 0.007s | `TestRFC0002DeterministicPolicies` (7 subtests), `TestRFC0003_CAS`, `TestRFC0004_TransitionMatrix` |

#### Environmental Toolchain Failures (Recorded Separately)
| Package | Result | Root Cause Analysis |
| :--- | :--- | :--- |
| `github.com/Orolar-CNR/IntentCore/transport/abtp` | **FAIL (Build)** | `loader.go:3` contains `//go:generate go run github.com/cilium/ebpf/cmd/bpf2go abtp xdp/abtp.c`. `bpf2go` requires `clang` and `llvm-strip`. Missing `llvm-strip` in test environment prevents eBPF bytecode compilation and auto-generation of `abtpObjects`. |
| `github.com/Orolar-CNR/IntentCore/cmd/intentcored` | **FAIL (Build)** | Depends on `transport/abtp`, failing compilation due to missing `abtpObjects` from eBPF generator. |
| `github.com/Orolar-CNR/IntentCore/internal/bootstrap` | **FAIL (Build)** | Depends on `transport/abtp`, failing compilation due to missing `abtpObjects` from eBPF generator. |

---

## 3. Relevant Files and Actual Code Entry Points

The actual code entry points and contract definitions in the repository are mapped below:

```
IntentCore Kernel Structure:
├── contracts/               <- Canonical Interface Contracts
│   ├── envelope.go          <- SemanticEnvelope & Signature definitions
│   ├── admission.go         <- AdmissionPolicy, AdmissionResult, RejectionCode
│   ├── lifecycle.go         <- Lifecycle interface & TransitionRequest
│   ├── intent.go            <- IntentRecord & IntentState
│   ├── repository.go        <- StateRepository interface & Snapshot
│   ├── runtime.go           <- Pipeline interface
│   ├── transport.go         <- Transport interface & EnvelopeHandler
│   ├── history.go           <- HistoryRecorder interface
│   ├── proof.go             <- ProofRecorder interface
│   └── telemetry.go         <- TelemetryRecorder interface
├── core/                    <- Primitive Types & Domain Errors
│   ├── types.go             <- IntentID, StateVersion, TraceID, Authority
│   └── errors.go            <- System errors (ErrVersionConflict, ErrAdmissionRejected, etc.)
├── admission/               <- Governance Gate
│   └── policy.go            <- DeterministicPolicy & PolicyEvaluator implementations
├── runtime/                 <- Pipeline Orchestration
│   ├── validation.go        <- ValidateEnvelope (JSON parser)
│   ├── normalization.go     <- NormalizeEnvelope (Default values)
│   ├── pipeline.go          <- DefaultPipeline (Validate -> Normalize -> Admit -> Dispatch)
│   └── dispatcher.go        <- Dispatcher & Worker goroutines
├── lifecycle/               <- Authoritative State Machine
│   ├── machine.go           <- StateMachine struct & Transition implementation
│   ├── transition.go        <- IsAllowed (RFC-0004 Transition Matrix)
│   └── authority.go         <- CheckAuthority placeholder
├── state/                   <- State Persistence & Checkpoint Engine
│   ├── repository.go        <- MemoryRepository (CAS & Snapshot execution)
│   ├── cas.go               <- Compare-And-Swap core logic
│   ├── snapshot.go          <- Snapshot creation & checksum verification
│   └── archive_store.go     <- Long-term cold storage
├── history/                 <- Audit Ledger
│   ├── ledger.go            <- Append-only evidence ledger
│   └── recorder.go          <- Recorder implementation
└── transport/               <- Boundary Adapters
    └── abtp/                <- ABTP XDP eBPF Transport Loader & Adapter
```

---

## 4. Contract-by-Contract Audit Matrix

| Contract Area | RFC Specification | Actual Implementation Status | Compliance Assessment |
| :--- | :--- | :--- | :--- |
| **Semantic Envelope Wire Schema** | RFC-0001: SchemaVersion, EnvelopeID, AgentIdentity, EventTimestamp, Signatures, OpaquePayload, `intent_id`, `expected_version` | `contracts/envelope.go:85-93`: Struct `SemanticEnvelope` contains SchemaVersion, EnvelopeID, AgentIdentity, EventTimestamp, Signatures, TelemetryClass, OpaquePayload. | **PARTIAL**: `expected_version` field is missing from struct. |
| **Opaque Payload Handling** | RFC-0001: JSON-encoded opaque payload preserved without alteration | `contracts/envelope.go:92` (`OpaquePayload []byte`) & `runtime/normalization.go:142-147` | **COMPLIANT**: Payload remains raw byte slice; kernel performs no decoding of multimodal data. |
| **Admission Boundary Governance** | RFC-0002: Deterministic evaluation pipeline (`Authenticate` -> `Authorize` -> `Policy Evaluation` -> `Accept/Reject`) | `admission/policy.go:16-60` (`DeterministicPolicy`) & `admission/policy.go:78-115` (`PolicyEvaluator`) | **PARTIAL**: Evaluates version, missing IDs, identity, timestamp, and signature presence, but lacks actual signature verification. |
| **Cryptographic Signature Verification** | RFC-0015 & RFC-0002: Cryptographic authenticity, Ed25519/ECDSA signature verification | `admission/policy.go:37-39`: `if len(env.Signatures) == 0 { return dp.reject(...) }` | **SIMULATION / GAP**: Only verifies slice length $>0$. No public key resolution or byte signature checking. |
| **Structural Validation** | RFC-0001: Strict structural validation prior to admission | `runtime/validation.go:126-136`: Unmarshals JSON into `SemanticEnvelope`. Lines 132-133 have commented-out checks. | **PARTIAL**: Structural validation relies entirely on JSON unmarshaling error. |
| **Traceability & Correlation** | RFC-0000 & RFC-0001: TraceID propagation across transport, pipeline, and history ledger | `core/types.go:13` (`TraceID string`) & `history/ledger.go:18` (`TraceID core.TraceID`) | **GAP**: `runtime/dispatcher.go:401-411` leaves `req.TraceID` unpopulated when invoking `TransitionRequest`. |
| **Lifecycle State Machine** | RFC-0004: State Machine matrix (`Pending` -> `Validated` -> `Admitted` -> `Scheduled` -> `Executing` -> `Completed`/`Failed` -> `RolledBack`) | `lifecycle/transition.go:6-26` (`IsAllowed`) & `lifecycle/machine.go:28-94` (`StateMachine.Transition`) | **COMPLIANT**: Matrix strictly enforced; state mutations atomically committed via Repository CAS. |
| **Authority Evaluation** | RFC-0004: Check authority for transition authorization | `lifecycle/authority.go:10-13`: `func CheckAuthority(req TransitionRequest) bool { return true }` | **STUB**: Unconditional `true` placeholder. |
| **State Mutation & CAS** | RFC-0003: Atomic CAS versioning and immutable ledger | `state/cas.go` & `state/repository.go` | **COMPLIANT**: Thread-safe atomic Compare-And-Swap with monotonic version increments. |

---

## 5. Detailed Findings & Priority Classification

Each finding is classified into one of four severity levels based on evidence:
- **P0:** Critical security or correctness blocker (must be fixed before untrusted integration).
- **P1:** High priority integration blocker (required before first production integration milestone).
- **P2:** Medium priority roadmap item (required for later roadmap phases).
- **P3:** Minor enhancement or code hygiene.

---

### Finding F-01: Cryptographic Signature Verification is Simulated (P0)

- **Severity:** **P0 (Critical Security Blocker)**
- **Category:** Security / Admission Boundary
- **Location:** `admission/policy.go`, Lines 37–39
```go
// 4. Missing Signature (Simulating empty slice check for now)
if len(env.Signatures) == 0 {
    return dp.reject(contracts.CodeMissingSignature, "Signature is required"), nil
}
```
- **Verified Fact:** `DeterministicPolicy` only verifies that `len(env.Signatures) > 0`. It does not perform cryptographic signature verification against the payload or header bytes, nor does it resolve public keys or verify key-signature pairs.
- **Inferred Risk:** Any unauthenticated caller can pass admission by supplying arbitrary non-empty byte arrays in `Signatures` (e.g., `[]Signature{{Signature: []byte("fake"), PublicKey: []byte("fake")}}`), bypassing the governance gate.
- **Recommendation:** Implement a pluggable `SignatureVerifier` interface in the `admission` package supporting Ed25519 cryptographic verification over a canonical byte serialization of `SemanticEnvelope`.

---

### Finding F-02: Missing `expected_version` in `SemanticEnvelope` Wire Contract (P1)

- **Severity:** **P1 (Required for First Integration Milestone)**
- **Category:** Contract Inconsistency
- **Location:** `contracts/envelope.go`, Lines 85–93 vs. `docs/rfcs/RFC-0001-semantic-envelope.md`, Section 4
- **Verified Fact:** RFC-0001 Section 4 explicitly mandates:
  > *"The SemanticEnvelope MUST include the following two kernel-owned fields, both present simultaneously... `intent_id`... `expected_version` — the version of the target state the Intent was authored against."*
  However, `contracts.SemanticEnvelope` in `contracts/envelope.go` omits `ExpectedVersion core.StateVersion`.
- **Inferred Risk:** External clients (such as Aetherium integration adapters) cannot specify optimistic concurrency expectations at envelope creation time. The dispatcher currently hardcodes `ExpectedVersion: 0` when queueing new transitions in `runtime/dispatcher.go:405`.
- **Recommendation:** Add `ExpectedVersion core.StateVersion` to `contracts.SemanticEnvelope` in `contracts/envelope.go` and update `runtime/dispatcher.go` to forward `env.ExpectedVersion`.

---

### Finding F-03: `TraceID` Not Propagated in Runtime Dispatcher (P1)

- **Severity:** **P1 (Required for First Integration Milestone)**
- **Category:** Traceability & Correlation
- **Location:** `runtime/dispatcher.go`, Lines 401–412
```go
tReq := contracts.TransitionRequest{
    IntentID:        req.Envelope.EnvelopeID,
    FromState:       "", // indicates new
    ToState:         contracts.StatePending,
    ExpectedVersion: 0,
    ActorID:         req.Envelope.AgentIdentity,
    Metadata: map[string]any{
        "AgentIdentity": req.Envelope.AgentIdentity,
        "Payload":       req.Envelope.OpaquePayload,
    },
}
```
- **Verified Fact:** `core.TraceID` is defined in `core/types.go` and supported in `TransitionRequest` (`contracts/lifecycle.go`) and `EvidenceRecord` (`history/ledger.go`). However, `dispatcher.go` does not populate `req.TraceID` when constructing `TransitionRequest`.
- **Inferred Risk:** When transitions are recorded in the history ledger, `TraceID` is empty. End-to-end distributed tracing between IntentCore and Aetherium (`IntentCore.TraceID -> Aetherium.trace_id`) is broken at the runtime dispatch boundary.
- **Recommendation:** Add `TraceID core.TraceID` to `contracts.SemanticEnvelope` (or extract from context) and pass it to `TransitionRequest.TraceID` in `runtime/dispatcher.go`.

---

### Finding F-04: Structural Envelope Validation Relies Exclusively on JSON Unmarshaling (P2)

- **Severity:** **P2 (Required for Roadmap Phase 3)**
- **Category:** Validation & Input Hygiene
- **Location:** `runtime/validation.go`, Lines 126–136
```go
func ValidateEnvelope(payload []byte) (contracts.SemanticEnvelope, error) {
    var env contracts.SemanticEnvelope
    if err := json.Unmarshal(payload, &env); err != nil {
        return env, core.ErrValidationFailed
    }

    // Basic structural checks could go here (e.g. UUID empty check)
    // if env.EnvelopeID == uuid.Nil { return env, core.ErrValidationFailed }

    return env, nil
}
```
- **Verified Fact:** `ValidateEnvelope` only checks JSON syntactic validity. Structural invariants (e.g., non-nil UUIDv4, non-zero timestamp, valid schema version) are commented out in `validation.go` and deferred to `admission/policy.go`.
- **Inferred Risk:** Malformed or incomplete envelopes pass through the Validation and Normalization stages unhindered, relying entirely on Admission for structural rejection.
- **Recommendation:** Re-enable basic structural checks in `ValidateEnvelope` (e.g., verifying `EnvelopeID != uuid.Nil`) to fail fast at the Validation stage per RFC-0001 Section 2.

---

### Finding F-05: Missing Rejection Code `ERR_INVALID_SIGNATURE` (P2)

- **Severity:** **P2 (Required for Roadmap Phase 3)**
- **Category:** Contract Completeness
- **Location:** `contracts/admission.go`, Lines 18–26
- **Verified Fact:** `contracts.RejectionCode` defines `CodeMissingSignature` (`ERR_MISSING_SIGNATURE`), but lacks a specific code for invalid/forged signatures (e.g., `CodeInvalidSignature = "ERR_INVALID_SIGNATURE"`).
- **Inferred Risk:** When cryptographic signature evaluation is introduced, invalid signatures would have to fall back to `CodeGeneralPolicyViolation`, obscuring audit metrics.
- **Recommendation:** Define `CodeInvalidSignature RejectionCode = "ERR_INVALID_SIGNATURE"` in `contracts/admission.go`.

---

### Finding F-06: Authority Check is an Unconditional Placeholder (P2)

- **Severity:** **P2 (Required for Roadmap Phase 3)**
- **Category:** Governance
- **Location:** `lifecycle/authority.go`, Lines 9–14
```go
func CheckAuthority(req contracts.TransitionRequest) bool {
    // For Phase 1, we simply return true. Real implementations would verify
    // the requested authority against the actor ID and target state.
    return true
}
```
- **Verified Fact:** `CheckAuthority` returns `true` unconditionally for all transition requests.
- **Inferred Risk:** Role-based or capability-based authorization is not enforced at the lifecycle transition boundary.
- **Recommendation:** Maintain as documented stub until Phase 3.3 Zero-Trust rollout, but ensure callers pass valid `Authority` identifiers.

---

### Finding F-07: Missing Toolchain Prerequisites for eBPF Transport Build (P3)

- **Severity:** **P3 (Environment / Build Documentation)**
- **Category:** Toolchain Prerequisites
- **Location:** `transport/abtp/loader.go`, Line 3
- **Verified Fact:** Running `go generate ./...` or `go test ./transport/abtp` fails in standard Go test environments unless `clang` and `llvm-strip` are installed.
- **Inferred Risk:** Developers or CI environments without LLVM toolchains experience package build failures in `transport/abtp`, `cmd/intentcored`, and `internal/bootstrap`.
- **Recommendation:** Add a build tag (e.g., `//go:build ebpf`) or document LLVM toolchain setup prerequisites in `README.md`.

---

## 6. Boundary Analysis: IntentCore Kernel vs. Aetherium

To preserve system stability, the responsibility boundaries between IntentCore and Aetherium are categorized below into Verified Contracts, Proposed Enhancements, and Violations to Avoid.

```
+-------------------------------------------------------------------+
|                        INTENTCORE KERNEL                          |
|  - Owns Admission, Identity, Authority, Coordination, Lifecycle   |
|  - Owns CAS State Repository, Snapshots, History Ledger           |
|  - Stores Payload strictly as Opaque Bytes (OpaquePayload []byte) |
+-------------------------------------------------------------------+
                                  |
                   SemanticEnvelope (ABTP Transport)
                                  |
                                  v
+-------------------------------------------------------------------+
|                        AETHERIUM-MANIFEST                         |
|  - Owns Multimodal Interpretation (Text, Image, Video, Trajectory)|
|  - Owns Modality-Specific Parsing & Manifestation State Machine   |
|  - Owns Downstream Rendering (Canvas2D, WebGPU, SDF, Particles)  |
+-------------------------------------------------------------------+
```

### Boundary Verification Rules

1. **Kernel Opaque Payload Invariant:**
   - **Fact:** `contracts.SemanticEnvelope` carries `OpaquePayload []byte`.
   - **Rule:** IntentCore kernel code MUST NOT unmarshal or inspect the interior JSON schema of `OpaquePayload` for Aetherium rendering fields (such as shader parameters, SDF distance fields, canvas dimensions, or particle counts).

2. **Lifecycle State Authority Invariant:**
   - **Fact:** IntentCore defines authoritative states in `contracts/intent.go`: `Pending`, `Validated`, `Admitted`, `Scheduled`, `Executing`, `Completed`, `Failed`, `RolledBack`.
   - **Rule:** Aetherium's internal processing states (e.g., `INTERPRETING`, `MANIFESTING`, `FADING`) MUST NOT overwrite or redefine IntentCore lifecycle states. Aetherium processing states map strictly as sub-states or metadata attached to IntentCore's `Executing` or `Completed` states.

3. **Identity & Traceability Mapping:**
   - `IntentCore.EnvelopeID` $\rightarrow$ `Aetherium.external_message_id`
   - `IntentCore.TraceID` $\rightarrow$ `Aetherium.trace_id`
   - Aetherium-internal identifiers (`signal_id`, `episode_id`) remain strictly within Aetherium boundaries.

---

## 7. Multimodal Ingress Assessment

A future **Aetherium Multimodal Ingress Contract** will bring rich multimodal inputs (text, image, video, action_trajectory) into the system. This audit assesses IntentCore's readiness for such ingress:

### Readiness Matrix

| Modality / Aspect | Ingress Path | IntentCore Kernel Handling | Missing Boundary Protection |
| :--- | :--- | :--- | :--- |
| **Text Ingress** | Prompt / Command payload | Carried via `OpaquePayload []byte` | Maximum payload byte size enforcement in `admission/policy.go` |
| **Image Ingress** | Compressed image buffer / URI | Carried via `OpaquePayload []byte` | Content-Type header metadata (`ContentType string`) on `SemanticEnvelope` |
| **Video Ingress** | Frame buffer / Streaming URI | Carried via `OpaquePayload []byte` | Payload size guardrails to prevent memory exhaustion DoS attacks |
| **Action Trajectory** | Vector / Spatial trajectory JSON | Carried via `OpaquePayload []byte` | Schema validation at external Adapter boundary before reaching kernel |

### Architectural Recommendations for Multimodal Ingress
1. **Envelope Content-Type Metadata:** Add optional `ContentType string` and `PayloadVersion string` fields to `SemanticEnvelope` (e.g. `application/vnd.aetherium.multimodal+json`) so ingress adapters can route payloads without opening or parsing the payload byte slice inside the kernel.
2. **Payload Size Guardrails:** Add a configurable `MaxPayloadBytes` policy check in `admission/policy.go` (e.g., 10 MB default) to prevent oversized video/image buffers from causing kernel Out-Of-Memory (OOM) conditions.

---

## 8. Smallest Recommended Next Implementation Task

To advance IntentCore's integration readiness with maximum safety and smallest scope, we recommend executing **exactly one bounded task**:

### Task Title: Pluggable Signature Verification & TraceID Propagation (Task IC-INT-01)

#### Objective
Replace the simulated signature check in `admission/policy.go` with a pluggable `SignatureVerifier` interface (supporting Ed25519 verification) and ensure `TraceID` and `ExpectedVersion` are correctly propagated from `SemanticEnvelope` through the dispatcher to the lifecycle and history ledger.

#### Scope of Changes
1. **`contracts/envelope.go`**: Add `ExpectedVersion core.StateVersion` and optional `ContentType string` to `SemanticEnvelope`.
2. **`contracts/admission.go`**: Add `CodeInvalidSignature RejectionCode = "ERR_INVALID_SIGNATURE"`.
3. **`admission/policy.go`**:
   - Define interface `SignatureVerifier interface { Verify(env contracts.SemanticEnvelope) bool }`.
   - Update `DeterministicPolicy` to evaluate signatures via `SignatureVerifier` (defaulting to a mock/noop verifier when unconfigured, and rejecting invalid signatures with `CodeInvalidSignature`).
4. **`runtime/dispatcher.go`**: Update worker dispatch logic to forward `TraceID` and `ExpectedVersion` from `SemanticEnvelope` into `TransitionRequest`.

#### Acceptance Criteria
1. `DeterministicPolicy` rejects envelopes with malformed/invalid cryptographic signatures when a real `SignatureVerifier` is configured, returning `CodeInvalidSignature`.
2. `TransitionRequest` produced by `runtime/dispatcher.go` contains the non-empty `TraceID` and `ExpectedVersion` supplied in `SemanticEnvelope`.
3. `history/recorder.go` successfully records the non-empty `TraceID` in the append-only `EvidenceRecord`.
4. All existing 18 pure Go unit tests continue to pass with zero regressions.
5. New unit tests verifying signature rejection, `TraceID` propagation, and `ExpectedVersion` forwarding pass cleanly.

#### Required Test Plan
- **Unit Test 1 (`admission/policy_test.go`):** Test `DeterministicPolicy` with a mock `SignatureVerifier` returning `false`, verifying `CodeInvalidSignature` rejection.
- **Unit Test 2 (`runtime/dispatcher_test.go`):** Verify `dispatcher.Dispatch` forwards `TraceID` and `ExpectedVersion` to `contracts.TransitionRequest`.
- **Unit Test 3 (`tests/integration_trace_test.go`):** End-to-end test verifying `SemanticEnvelope` -> `Pipeline.Execute` -> `Dispatcher` -> `Lifecycle.Transition` -> `History.RecordTransition` preserves `TraceID` in `EvidenceRecord`.

---
