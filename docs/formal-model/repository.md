**Related RFCs:** RFC-0000 (Architectural Principles)

# Formal Model: Repository

This document outlines the formal model for the IntentCore State Repository, ensuring single-source-of-truth and CAS (Compare-And-Swap) consistency.

## State Space ($R$)
The repository state $R$ at time $t$ is defined as a mapping from Intent IDs to a tuple of (State, Version).
$R_t : ID \to (State, \mathbb{N})$

## Operations

### 1. Compare-And-Swap (CAS)
```
CAS(id, v_exp, s_new):
    if R_t(id).Version == v_exp:
        R_{t+1}(id) = (s_new, v_exp + 1)
        History.Append(id, s_new, v_exp + 1)
        return SUCCESS
    else:
        return CONFLICT
```

## Consistency Invariants
1.  **Monotonic Versioning:** $\forall id, \text{Version}(id)$ is strictly monotonically increasing.
2.  **Immutable History:** The history ledger $L$ is an ordered sequence. $\forall i < |L|, L[i]$ is constant and immutable over time.
