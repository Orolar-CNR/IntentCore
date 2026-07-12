**Related RFCs:** RFC-0000 (Architectural Principles)

# Formal Model: Lifecycle

This document defines the formal mathematical model for the IntentCore Lifecycle State Machine, suitable for translation into formal verification languages like TLA+ or Alloy.

## States ($S$)
$S = \{ Pending, Validated, Admitted, Scheduled, Executing, Completed, Failed, RolledBack \}$

## Terminal States ($T$)
$T \subset S$
$T = \{ Completed, RolledBack \}$

## Allowed Transitions ($\delta$)
$\delta \subseteq S \times S$
$\delta = \{
    (Pending, Validated),
    (Validated, Admitted),
    (Admitted, Scheduled),
    (Scheduled, Executing),
    (Executing, Completed),
    (Executing, Failed),
    (Failed, RolledBack)
\}$

## Invariants ($I$)
1.  **Terminality:** $\forall t \in T, \nexists s \in S : (t, s) \in \delta$
    *(No transitions are allowed out of a terminal state).*
2.  **Determinism:** Given a state $s \in S$, an event $e$, and transition function $F(s, e) \to s'$, $F$ must be a pure function.

## Preconditions and Postconditions
For any state transition $(s_{current}, s_{next}) \in \delta$:
*   **Precondition:** `Repository.Version == ExpectedVersion`
*   **Postcondition:** `Repository.Version == ExpectedVersion + 1` AND `History.Append(TransitionRecord)`
