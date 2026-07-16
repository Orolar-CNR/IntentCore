# IntentCore Baseline State Report (Milestone 2.8)

## 1. Executive Summary

This document captures the frozen baseline state of the IntentCore repository following the completion of Milestone 2.8. It serves as a verifiable snapshot of the core architecture implementing RFC-0001 through RFC-0004.

The system enforces a strict one-way data dependency flow, successfully receiving raw payloads via the Transport Boundary (ABTP), validating them into the canonical `SemanticEnvelope`, evaluating them via Deterministic Admission Policies, queueing them in the Runtime Dispatcher, and executing state transitions via the Lifecycle engine backed by a CAS-enabled State Repository.

## 2. Test Execution Report

A full test suite was executed across the codebase (`go test -v -race ./...`).

**Summary of Results:**
- All Go modules compiled successfully (after generating C datapath binaries via `go generate`).
- No race conditions were detected.

| Package | Status | Notable Details |
| :--- | :--- | :--- |
| `proof` | PASS | Verified `TestInMemoryProofRecorder`. |
| `runtime` | PASS | Verified `TestDispatcher`, `TestVerticalSlice`, and pipeline behaviors including timeouts and cancellations. |
| `state` | PASS | Verified `TestRepository_SnapshotAndRecover`. |
| `telemetry` | PASS | Verified `TestInMemoryRecorder`. |
| `tests` | PASS | Verified end-to-end deterministic policies (`TestRFC0002DeterministicPolicies`), CAS constraints (`TestRFC0003_CAS`), and state transitions (`TestRFC0004_TransitionMatrix`). |
| `transport/abtp` | PASS | Verified `TestAdapter` and eBPF Loader attach/detach logic. |
| `admission` | Skipped | No dedicated test files found. |
| `lifecycle` | Skipped | No dedicated test files found. |
| `history` | Skipped | No dedicated test files found. |
| `core` | Skipped | No test files. |
| `contracts` | Skipped | No test files. |

*(Note: The full raw output of the test execution is preserved in `docs/raw_test_report.log` for deep traceability.)*

## 3. Subsystem Breakdown

### 3.1 Transport Boundary (ABTP)
The `transport/abtp` package implements the `contracts.Transport` interface. It provides an `Adapter` that listens for UDP packets on port 10000. It frames the incoming bytes and asynchronously delegates them to the provided `contracts.EnvelopeHandler` (the Runtime Pipeline) without evaluating semantic intent, honoring the RFC-0000 separation of concerns.

### 3.2 Admission Policy (Deterministic Rejection)
The `admission` package evaluates the payload using deterministic rules. The default `DeterministicPolicy` evaluates `contracts.SemanticEnvelope` to ensure:
- The schema version is `1.0.0`
- Essential fields (`EnvelopeID`, `AgentIdentity`, `Signatures`, `OpaquePayload`, and `EventTimestamp`) are present.
Failure on any constraint immediately short-circuits the pipeline with a specific `RejectionCode`, rejecting the payload before execution.

### 3.3 Runtime Dispatch Core
The `runtime` package manages execution logic through a `Pipeline` orchestrator and a `Dispatcher`.
- **Pipeline:** Implements timeouts, context propagation, and retries. It coordinates the execution flow: Validation -> Normalization -> Admission -> Dispatch.
- **Dispatcher:** Provides isolated, non-blocking asynchronous execution via a configurable buffered channel (`chan DispatchRequest`). A background worker reads queued, admitted envelopes and requests state transitions from the Lifecycle engine.

### 3.4 State & Recovery
The core of state mutation occurs in the `lifecycle` and `state` packages.
- **Lifecycle Engine:** `StateMachine` serves as the sole authority for state mutations. It checks `IsAllowed` against the target state (e.g., `Pending` -> `Validated`) preventing invalid skips or updates from terminal states, checking authority, and coordinating persistence.
- **State Repository:** Implements strict Compare-and-Swap (CAS) in memory with thread safety (`sync.RWMutex`). It creates Snapshots containing offsets and intent counts, simulating a ledger recovery model where recovery is driven by an internal checkpoint payload.

## 4. Intent Flow

The trace of a typical intent as it flows through the current IntentCore implementation is as follows:

1. **Transport Entry:** A raw byte payload arrives via the ABTP UDP socket (`transport/abtp.Adapter.listen()`).
2. **Execution Triggered:** The adapter pushes the payload bytes into `DefaultPipeline.Execute()`.
3. **Validation & Normalization:** The raw payload is unmarshaled into the canonical `contracts.SemanticEnvelope`.
4. **Admission Evaluation:** The `PolicyEvaluator` loops over policies (primarily `DeterministicPolicy`), verifying the `SemanticEnvelope` schema and data integrity.
5. **Telemetry / Proof (Optional):** If configured, telemetry records acceptance and the proof recorder registers a baseline trace.
6. **Dispatch Queuing:** The envelope enters the `Dispatcher`'s channel queue.
7. **Worker Processing:** A background dispatcher goroutine pulls the envelope and converts it into a `TransitionRequest` targeting the `Pending` state.
8. **Lifecycle Transition:** `lifecycle.StateMachine.Transition()` checks the transition matrix and builds the next state version.
9. **State Commit:** The `state.Repository.CompareAndSwap()` ensures no race conditions occur and atomicly updates the internal memory store.
10. **History Audit:** Finally, a history record is emitted for full traceability.
