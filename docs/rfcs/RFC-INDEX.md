# IntentCore RFC Index

Status: Informative (non-normative index; each linked RFC remains the normative source)

## 1. Canonical Filenames

| RFC | Filename | Responsibility |
|---|---|---|
| RFC-0001 | `docs/rfcs/RFC-0001-semantic-envelope.md` | Canonical wire contract into the kernel; kernel-owned identity fields (`intent_id`, `expected_version`) |
| RFC-0002 | `docs/rfcs/RFC-0002-admission-boundary.md` | What Admission may and may not know; duplication/structural validation only |
| RFC-0003 | `docs/rfcs/RFC-0003-state-repository.md` | Single authoritative store; CAS; Ledger; snapshot; recovery (renamed from `state-topology` to avoid clashing with Phase 3 federation/topology terminology) |
| RFC-0004 | `docs/rfcs/RFC-0004-lifecycle-state-machine.md` | Deterministic lifecycle FSM; transition authorization |
| RFC-0005 | `docs/rfcs/RFC-0005-event-bus-contract.md` | Runtime event plane; delivery, retry, DLQ, ordering — transport boundary only, never enters the kernel |
| RFC-0006 | `docs/rfcs/RFC-0006-history.md` | Immutable evidence; synchronous projection of the Ledger |
| RFC-0007 | `docs/rfcs/RFC-0007-proof.md` | Verifiable evidence derived exclusively from History/Ledger |

## 2. Reading Order (Informative Only)

The diagram below is an **informative reading order** for newcomers to the spec set. It is not the full dependency graph — see Section 3.

```
RFC-0001 SemanticEnvelope
        │
        ▼
RFC-0002 Admission Boundary
        │
        ▼
RFC-0003 State Repository
        │
        ▼
RFC-0004 Lifecycle State Machine
        │
        ▼
RFC-0006 History
        │
        ▼
RFC-0007 Proof

RFC-0005 Event Bus Contract
        │
        └────────────► (Transport Boundary only)
             wraps SemanticEnvelope
             never enters the kernel
```

## 3. True Dependency Graph (Binding)

RFC-0006 (History) depends **directly** on RFC-0003 (State Repository) — it projects from the Ledger, not from Lifecycle's transition logic. RFC-0004 (Lifecycle) independently depends on RFC-0003 to issue CAS mutations. The reading-order chain above should not be mistaken for implying History depends on Lifecycle for its data; it depends on the Ledger.

```
RFC-0001 ──▶ RFC-0002 ──▶ RFC-0004 ──▶ RFC-0006 ──▶ RFC-0007
                              │            ▲
                              ▼            │
                          RFC-0003 ────────┘
```

## 4. Boundary Notes (Non-Redundant Summary)

- RFC-0005 (Event Bus) never enters the kernel. Its `idempotency_key`, `trace_id`, `correlation_id` are transport/observability constructs only, and MUST be unwrapped/translated into kernel-owned SemanticEnvelope fields at the transport/adapter boundary before reaching Admission.
- RFC-0002 (Admission) checks duplication and structural validity of `intent_id`/`expected_version` only. It does not read Repository and does not resolve version conflicts — that happens exclusively at RFC-0004 → RFC-0003 CAS time.
- RFC-0006 (History) is a synchronous projection of RFC-0003's Ledger — never an independent write store.
- RFC-0007 (Proof) is derived exclusively from RFC-0006 — never from RFC-0005 artifacts, and never speculatively ahead of a materialized History entry.
