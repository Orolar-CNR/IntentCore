# RFC-0007 — Proof

Status: Draft
Category: Normative Specification
Scope: Verifiable evidence generation, derivation source, consistency guarantees with History/Ledger

## 1. Abstract

This document defines the normative contract for **Proof** within IntentCore. Proof is the verifiable evidence layer that attests to the correctness and authenticity of a committed state transition. Proof MUST be derived exclusively from the same authoritative record chain as History (RFC-0006) and Ledger (RFC-0003). Proof MUST NOT constitute an independent path of state truth.

This RFC depends on RFC-0006 (History) and, transitively, RFC-0003-state-repository (Repository/Ledger). Because RFC-0006 mandates synchronous History materialization (RFC-0006 Section 4.1), Proof MAY read from History directly without staleness concerns, per Section 3 of this document.

Filename note: this document is `docs/rfcs/RFC-0007-proof.md`, companion to `docs/rfcs/RFC-0006-history.md` and `docs/rfcs/RFC-0003-state-repository.md`. The straight-line reading order (0001→0002→0004→0006→0007) shown elsewhere in this RFC set is informative only; the binding dependency graph (see RFC-0006 §3.0) shows RFC-0006 depending directly on RFC-0003, not solely through RFC-0004.

## 2. Normative Language

MUST, MUST NOT, SHOULD, SHOULD NOT, and MAY are interpreted per RFC 2119.

## 3. Derivation Source

### 3.1 History as the Proof Input (Locked Decision)

MUST

- Proof generation MUST derive its evidence from History (RFC-0006) entries, which are themselves synchronous projections of committed Ledger entries.
- Because History materialization is synchronous with Ledger commit (RFC-0006 §4.1), Proof reading from History is guaranteed to reflect the current committed state with no lag window. Proof implementations MUST rely on this guarantee rather than re-deriving it independently.
- If a future revision of RFC-0006 permits asynchronous History materialization, this RFC MUST be revised in tandem to specify either (a) a bounded staleness SLA that Proof generation must tolerate or announce, or (b) a fallback path for Proof to read the Ledger directly during the staleness window. Until such a revision exists, Proof MUST assume synchronous availability.

MUST NOT

- Proof MUST NOT maintain or consult any transition record store independent of the Ledger/History chain.
- Proof MUST NOT be generated from Event Bus Contract (RFC-0005) artifacts (`trace_id`, `correlation_id`, delivery ACK/NACK state, or any other bus-level metadata). These are transport/delivery-layer observability constructs and carry no authority over state transitions.
- Proof MUST NOT be generated speculatively before the corresponding History entry (and therefore the underlying Ledger commit) exists.

### 3.2 Traceability

MUST

- Every Proof record MUST reference the specific History entry (and, transitively, the specific Ledger sequence number) from which it was derived.
- A Proof record MUST be reproducible: given the same referenced Ledger/History entry, re-deriving Proof MUST yield an equivalent, verifiable result.

## 4. Consistency Guarantee

MUST

- Proof MUST remain consistent with its source History/Ledger entry for the lifetime of the system. If the underlying Ledger is append-only and immutable (RFC-0003, RFC-0006 §5), Proof correctness is guaranteed to be stable once generated.
- Proof generation MUST NOT be retried against a different or superseding transition record once the original committed entry is referenced. A new transition (e.g., a compensating action) requires a new, separate Proof record tied to its own new History/Ledger entry — never a mutation of a prior Proof record.

MUST NOT

- Proof records, once generated, MUST NOT be edited or deleted, consistent with the immutability principle governing History and Ledger.

## 5. Generation Timing

MUST

- Proof generation MUST occur only after the corresponding History entry has been successfully materialized (RFC-0006 §3.2, step 4).
- Proof generation MUST NOT block or gate the completion signal of the underlying transition beyond what is already required by History materialization (RFC-0006 §4.1). In other words: transition completion depends on Ledger commit + History materialization; Proof generation is a subsequent, independently-retryable step and MUST NOT introduce additional blocking on the critical path of transition completion.

SHOULD

- Proof generation SHOULD be attempted promptly after History materialization to minimize the window in which a committed transition lacks verifiable evidence.
- If Proof generation fails transiently, it SHOULD be retried against the same referenced History/Ledger entry until it succeeds, without altering or bypassing the source record.

## 6. Verification Interface

SHOULD

- Proof SHOULD expose a verification interface that allows an external party to confirm that a given Proof record is consistent with the referenced History/Ledger entry, without requiring write access to Repository, Lifecycle, or History.

MAY

- Proof MAY use cryptographic signing or hashing over the referenced Ledger/History entry content to support non-repudiation and tamper-evidence, consistent with Architecture Landscape Section 22 (Security and Integrity) as applied to the Event Bus, extended here to the kernel's evidence layer.

## 7. Relationship Summary

```
Repository (Ledger)  ──canonical append-only commit──▶
History (RFC-0006)   ──synchronous projection──────────▶
Proof (this RFC)      ──derived, verifiable evidence────▶ external consumers / auditors
```

MUST

- No stage in this chain MUST be bypassed by a later stage. Proof MUST NOT read Ledger directly while History exists and is synchronously available (per the locked decision in Section 3.1); this keeps a single, well-defined derivation path rather than two parallel valid paths.

## 8. Summary of Mandatory Rules

The system SHALL enforce the following invariants:

1. Proof is derived exclusively from History (RFC-0006), which is itself a synchronous projection of the canonical Ledger (RFC-0003).
2. Proof MUST NOT be derived from Event Bus (RFC-0005) artifacts under any circumstance.
3. Every Proof record MUST be traceable to a specific, immutable History/Ledger entry.
4. Proof records are immutable once generated; corrections require new transitions with their own new Proof records.
5. Proof generation MUST NOT introduce additional blocking on transition completion beyond what History materialization already requires.
6. There is exactly one derivation path (Ledger → History → Proof); it MUST NOT be bypassed or duplicated.
