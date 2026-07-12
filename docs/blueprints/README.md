**Related RFCs:** RFC-0000 (Architectural Principles)

# Blueprints

Blueprints serve as the **Implementation-independent reference architecture** for IntentCore.

While the RFCs define the strict, mandatory, and normative contracts (using language like `MUST` and `SHALL`), the Blueprints in this directory are informative and descriptive. They illustrate *how* the system is expected to be built, the strategic technical decisions, and the long-term vision without violating the core architectural contracts.

## Purpose

*   **Reference Architecture:** Provides architectural diagrams, flows, and structural strategies (e.g., how to use XDP/eBPF for transport).
*   **Informative Guidance:** Uses descriptive language (e.g., "is expected to", "should consider") rather than strict normative rules.
*   **Bridge to Implementation:** Helps engineers understand the "why" and "how" behind the specifications before writing actual code.

## Relationship to RFCs

```text
RFCs
  │ (define strict rules)
  ▼
Architectural Contracts
  │ (guide)
  ▼
Blueprints
  │ (support)
  ▼
Implementation
```

Blueprints can evolve rapidly as new technologies emerge, provided they continue to respect the locked RFCs.
