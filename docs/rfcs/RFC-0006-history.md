# RFC-0006 — History

Status: Draft
Category: Normative Specification
Scope: Authoritative transition evidence, Ledger relationship, projection semantics, consistency guarantees

## 1. Abstract

This document defines the normative contract for **History** within IntentCore. History is the authoritative, immutable record of every successful Lifecycle transition. History MUST NOT function as an independent write store. History is a **synchronous derived projection** of the canonical Ledger maintained by the State Repository (RFC-0003).

This RFC exists to close a specific architectural risk: the possibility of two independently-written "authoritative" logs (Ledger and History) drifting out of sync, which would violate the Single Source of Truth invariant (Section 3.4 of the Architecture Landscape).

## 2. Normative Language

MUST, MUST NOT, SHOULD, SHOULD NOT, and MAY are interpreted per RFC 2119.

## 3. Relationship to Ledger (Repository, RFC-0003-state-repository)

### 3.0 Dependency Note

This RFC depends **directly** on RFC-0003 (State Repository), not merely by inheritance through RFC-0004 (Lifecycle State Machine). Lifecycle orders and authorizes transitions, but the Ledger itself — the record History projects from — is owned and persisted exclusively by Repository. Any reading-order diagram elsewhere in this document set that shows a single straight chain (0001→0002→0004→0006→0007) is informative only; the binding dependency graph is:

```
RFC-0001 ──▶ RFC-0002 ──▶ RFC-0004 ──▶ RFC-0006 ──▶ RFC-0007
                              │            ▲
                              ▼            │
                          RFC-0003 ────────┘
```


### 3.1 Single Write Path

MUST

- Lifecycle MUST request exactly one authoritative mutation per transition, executed by Repository via Compare-And-Swap (CAS).
- Repository MUST be the only component that performs the durable, authoritative append to the Ledger.
- History MUST NOT be written to independently of a committed Ledger entry.
- Lifecycle MUST NOT emit evidence to History in parallel with, or as an alternative to, the Repository's Ledger commit.

MUST NOT

- There MUST NOT be two independent write paths that both claim authority over the same transition record.
- History MUST NOT accept writes from any source other than the Ledger projection mechanism defined in this RFC.

### 3.2 Write Sequence (Normative Order)

The write sequence for every successful transition MUST occur in this exact order:

```
1. Lifecycle evaluates and authorizes a transition
2. Lifecycle issues a CAS mutation request to Repository
3. Repository commits the transition to the Ledger (append-only)
4. Only upon successful Ledger commit: History projection is materialized
5. Lifecycle emits Proof and Telemetry (RFC-0007, RFC-0008), consistent with the committed Ledger entry
```

MUST

- Step 4 (History materialization) MUST NOT occur before step 3 (Ledger commit) succeeds.
- If step 3 fails or is rejected (e.g., CAS conflict), no History entry MUST be created for that attempt.
- History entries MUST be traceable 1:1 to a specific committed Ledger entry (by ledger sequence number or equivalent monotonic identifier).

### 3.3 History as Projection, Not Store

MUST

- History MUST be implemented as a read-model / projection derived from the Ledger.
- History MUST NOT introduce fields, facts, or claims that are not derivable from the committed Ledger entry and its associated Intent context.
- Deleting or modifying a History entry MUST NOT be possible through any interface. Corrections, if ever required, MUST be represented as new compensating Ledger entries, never as mutation of existing History records.

SHOULD

- History SHOULD denormalize or enrich the raw Ledger entry (e.g., resolved human-readable transition names, correlated Intent metadata) for audit and query convenience, provided all enrichment is derived from the same committed record and does not introduce independent facts.

## 4. Consistency Model

### 4.1 Synchronous Materialization (Locked Decision)

MUST

- History projection MUST be materialized **synchronously**, within the same logical write sequence as the Ledger commit (see 3.2), before the mutation request is considered complete from Lifecycle's perspective.
- A transition MUST NOT be reported as "Completed" to any external observer until both the Ledger commit and the corresponding History materialization have succeeded.
- If History materialization fails after a successful Ledger commit, the overall transition MUST be treated as **not yet complete**, and Lifecycle MUST retry History materialization from the committed Ledger entry until it succeeds, rather than proceeding as if History were optional.

This RFC file is stored as `docs/rfcs/RFC-0006-history.md` and its canonical companion is `docs/rfcs/RFC-0003-state-repository.md` (renamed from an earlier `state-topology` naming to avoid clashing with Phase 3 federation/cluster-topology terminology).

Rationale

- Synchronous materialization is required so that Proof (RFC-0007) can read from History without staleness or lag ambiguity. This avoids introducing an eventual-consistency window between "transition committed" and "evidence available," which would otherwise force Proof to choose between reading a possibly-stale History or bypassing History to read the Ledger directly.

SHOULD NOT

- Implementations SHOULD NOT introduce asynchronous, batched, or lazy History materialization. If a future implementation requires this for scale reasons, it MUST be proposed as a revision to this RFC, and RFC-0007 MUST be revised in tandem to specify how Proof handles the resulting staleness window.

### 4.2 Failure Handling

MUST

- A failed CAS mutation attempt (rejected by Repository) MUST NOT produce a History entry.
- A failed History materialization after a successful Ledger commit MUST NOT be silently dropped. It MUST be retried, and MUST be surfaced as an operational fault (e.g., via Telemetry, RFC-0008) if retries are exhausted, until resolved.
- History MUST eventually be consistent with every committed Ledger entry with no permanent gaps. Gaps are an operational incident, not an accepted steady state.

## 5. Immutability

MUST (restates and specializes Architecture Landscape Section 3.5)

- Every History entry, once materialized, MUST be append-only.
- History entries MUST NOT be modified or deleted by any component, operator tooling, or administrative interface.
- Any correction MUST take the form of a new Ledger entry (e.g., a compensating transition) that produces a new, additional History entry — never an edit of a prior one.

## 6. Relationship to Event Bus Observability (RFC-0005)

MUST

- `trace_id` and `correlation_id`, as defined in the Event Bus Contract (RFC-0005), are observability identifiers only.
- `trace_id` and `correlation_id` MUST NOT be treated as substitutes for History records.
- History entries MAY reference a `trace_id`/`correlation_id` as enrichment metadata (see 3.3, SHOULD) to aid cross-referencing with bus-level observability, but the authoritative content of a History entry MUST be derived solely from the committed Ledger entry, never from Event Bus envelope fields.

## 7. Access and Query

SHOULD

- History SHOULD expose a query interface indexed by `intent_id`, ledger sequence number, and transition timestamp.
- History SHOULD support ordered replay of transitions for a given `intent_id` for audit purposes.

MAY

- History MAY be exposed to external consumers (e.g., audit systems, dashboards) as a read-only interface. Such exposure MUST NOT provide any mutation path back into Repository or Lifecycle.

## 8. Summary of Mandatory Rules

The system SHALL enforce the following invariants:

1. Repository's Ledger is the sole authoritative append-only write path for transition records.
2. History MUST NOT be written independently of a committed Ledger entry.
3. History materialization MUST be synchronous with Ledger commit, within the same write sequence.
4. A transition is not "Completed" until both Ledger commit and History materialization succeed.
5. History entries are immutable; corrections are new Ledger entries, never edits.

6. Event Bus observability identifiers (`trace_id`, `correlation_id`) are never authoritative evidence and never a substitute for History.
