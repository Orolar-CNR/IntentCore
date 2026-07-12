# Formal Model: Repository

This document outlines the formal model for the IntentCore State Repository, ensuring single-source-of-truth and CAS (Compare-And-Swap) consistency.

## State Space ($R$)
The repository state $R$ at time $t$ is defined as a mapping from Intent IDs to a tuple of (State, Version).
$R_t : ID \to (State, \mathbb{N})$

## Operations

### Compare And Swap (CAS)
Given an Intent $id$, an expected version $v_{exp}$, and a new state $s_{new}$:

```text
CAS(id, v_exp, s_new):
    if R[id].version == v_exp:
        R[id].state = s_new
        R[id].version = v_exp + 1
        History.Append(id, v_exp, s_new)
        return SUCCESS
    else:
        return FAILURE(VERSION_CONFLICT)
```

## Consistency Invariants
1.  **Monotonic Versioning:** $\forall id, \text{Version}(id)$ is strictly monotonically increasing.
2.  **Immutable History:** The history ledger $L$ is an ordered sequence. $\forall i < |L|, L[i]$ is constant and immutable over time.
