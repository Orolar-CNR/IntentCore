# IntentCore → Aetherium-Manifest Integration Plan
## Coordination Adapter Proposal v0.1

## 1. Objective

IntentCore should act as the **coordination and admission boundary** for Aetherium-Manifest, not as the visual renderer and not as the semantic manifestation engine.

The desired separation is:

```text
External AI / Human / Client
          ↓
     IntentCore
  validation/admission
  identity/policy/lifecycle
          ↓
 Aetherium Adapter
          ↓
 Aetherium Manifest
  semantic interpretation /
  governed manifestation
          ↓
 Presence / Visual Contract
          ↓
 Renderer
```

For the current repository state, the practical Phase-0-compatible path is:

```text
External Message
      ↓
IntentCore SemanticEnvelope
      ↓
Aetherium IntentCore Adapter
      ↓
Phase-0 Local Interpreter
      ↓
Visual State Contract
      ↓
Canvas2D / WebGPU
```

The Phase-1 target can later become:

```text
External Intent
      ↓
IntentCore
      ↓
Aetherium AETH boundary
      ↓
Presence IR
      ↓
Governor
      ↓
Manifestation Runtime
      ↓
Renderer
```

---

# 2. What IntentCore Already Provides

The supplied IntentCore repository already contains executable foundations for:

### Transport

```text
transport/abtp
```

- UDP transport adapter
- `contracts.Transport`
- envelope delivery into runtime

### Validation

```text
runtime/validation.go
```

- JSON decoding into `SemanticEnvelope`

### Normalization

```text
runtime/normalization.go
```

- default normalization

### Admission

```text
admission/policy.go
```

- deterministic admission policy
- schema version check
- identity requirement
- signature presence requirement
- payload requirement
- timestamp requirement

### Lifecycle

```text
lifecycle/
```

- transition authority
- transition matrix
- authority checking
- CAS-backed persistence

### State

```text
state/
```

- repository
- CAS
- snapshots
- recovery mechanisms

### History / Proof / Telemetry

```text
history/
proof/
telemetry/
```

These are useful coordination primitives for Aetherium.

---

# 3. What Is Missing for Aetherium

The two repositories currently have **no executable integration boundary**.

A repository-wide search of the supplied Aetherium snapshot found no actual IntentCore client/adapter, SemanticEnvelope producer, ABTP client, or Presence-IR exchange.

Therefore the current integration status is:

```text
Aetherium ↔ IntentCore = NOT INTEGRATED
```

This is the first fact that should be recorded in the status tracker.

---

# 4. The Main Architectural Gap

IntentCore's current payload model is:

```go
SemanticEnvelope {
    SchemaVersion
    EnvelopeID
    AgentIdentity
    EventTimestamp
    Signatures
    TelemetryClass
    OpaquePayload
}
```

Aetherium Phase 0 currently expects:

```text
message text / signal
        ↓
interpretIntent()
        ↓
Visual State Candidate
```

There is no canonical adapter translating:

```text
SemanticEnvelope.OpaquePayload
```

into:

```text
Aetherium input message / signal
```

That adapter should be explicit.

---

# 5. Recommended New Boundary

Do NOT make IntentCore know about:

- Canvas2D
- WebGPU
- WGSL
- particle counts
- SDF
- visual colors
- renderer internals

Instead introduce an Aetherium-facing payload contract.

Example conceptual payload:

```json
{
  "message_type": "aetherium.input",
  "schema_version": "0.1",
  "message_id": "…",
  "source": "external_ai",
  "content_type": "text",
  "content": {
    "text": "สร้างสามเหลี่ยม"
  },
  "context": {
    "trace_id": "…",
    "parent_trace_id": null
  }
}
```

IntentCore should carry and govern this payload.

Aetherium should interpret it.

---

# 6. Do Not Put Visual State Inside IntentCore

Avoid this:

```text
IntentCore
  → hue
  → energy
  → density
  → turbulence
  → particle_count
  → WebGPU
```

That would collapse the architectural boundary.

Preferred:

```text
IntentCore
  → accepted semantic message
  → lifecycle / identity / policy evidence
          ↓
Aetherium
  → interpretation
  → semantic-to-visual mapping
  → Visual State / Presence IR
          ↓
Renderer
```

IntentCore coordinates the message's authority and lifecycle.

Aetherium owns manifestation semantics.

---

# 7. Proposed Adapter Package

Create a dedicated integration layer in Aetherium:

```text
integrations/
└── intentcore/
    ├── README.md
    ├── envelope-adapter.js
    ├── payload-schema.json
    ├── transport-client.js
    └── tests/
        ├── envelope-adapter.test.js
        └── integration-contract.test.js
```

The adapter's responsibility should be narrowly defined:

```text
IntentCore SemanticEnvelope
        ↓
verify expected Aetherium payload type
        ↓
extract accepted message
        ↓
attach trace / provenance
        ↓
submit to Aetherium input boundary
```

It must not directly invoke a renderer.

---

# 8. Recommended IntentCore Changes

## IC-001 — Add explicit application payload typing

Current:

```go
OpaquePayload []byte
```

Keep `OpaquePayload` at the kernel contract level, but define an application-level typed payload convention.

Recommended:

```text
content_type
payload_schema
payload_version
```

Example:

```json
{
  "content_type": "aetherium.input",
  "payload_version": "0.1"
}
```

Do not make the kernel depend directly on Aetherium-specific visual structures.

---

## IC-002 — Strengthen signature verification

Current admission checks:

```go
len(env.Signatures) == 0
```

This proves signature presence, not cryptographic verification.

Therefore the current implementation should NOT be described as complete cryptographic authentication.

Before production Aetherium integration, add:

```text
signature verification
      ↓
canonical envelope representation
      ↓
public-key resolution
      ↓
cryptographic verification
      ↓
admission
```

The exact algorithm and key-management model should be specified before implementation.

---

## IC-003 — Add explicit payload schema validation

`ValidateEnvelope()` currently JSON-unmarshals the envelope but does not fully enforce the documented envelope contract.

It should validate:

- UUID format/version
- required fields
- timestamp validity
- signatures structure
- payload encoding
- schema version
- application payload schema

This should remain deterministic and bounded.

---

## IC-004 — Add trace propagation

IntentCore already has:

```go
TraceID
```

Aetherium's Presence IR specifies:

```text
trace_id
parent_trace_id
trace_seed
```

Define an explicit mapping:

```text
IntentCore.TraceID
        ↓
Aetherium.trace_id
```

Do not invent `parent_trace_id` if no parent exists.

---

## IC-005 — Define a correlation ID / message ID mapping

IntentCore already has:

```text
EnvelopeID / IntentID
```

Aetherium has:

```text
signal_id
episode_id
trace_id
```

These should not be conflated.

Recommended:

```text
IntentCore EnvelopeID → Aetherium external_message_id
IntentCore TraceID    → Aetherium trace_id
Aetherium signal_id   → local input event identity
Aetherium episode_id  → local interaction grouping identity
```

This preserves layer-specific identity.

---

## IC-006 — Add explicit rejection feedback

The Aetherium adapter needs a machine-readable distinction between:

```text
ACCEPTED
REJECTED_POLICY
REJECTED_SCHEMA
REJECTED_SIGNATURE
REJECTED_VERSION
REJECTED_PAYLOAD
TIMEOUT
UNAVAILABLE
```

Do not convert all failures into a generic error.

---

## IC-007 — Make lifecycle semantics explicit

Current IntentCore dispatch creates a `Pending` transition.

Aetherium should not interpret:

```text
Pending
```

as:

```text
PROCESSING
```

unless an explicit mapping is ratified.

These are different state domains:

```text
IntentCore lifecycle state
        ≠
Aetherium manifestation state
```

A mapping contract is required.

---

# 9. Recommended Cross-Repository Contract

Create a small independent contract:

```text
AetheriumIntegrationEnvelope v0.1
```

Conceptual schema:

```json
{
  "contract": "aetherium.integration",
  "version": "0.1",
  "message_id": "uuid",
  "trace_id": "string",
  "source": {
    "type": "external_ai",
    "identity": "string"
  },
  "payload": {
    "content_type": "text",
    "data": "..."
  },
  "temporal": {
    "event_time": "number"
  }
}
```

The contract should deliberately NOT contain:

```text
hue
energy
density
shape
turbulence
particle_count
shader
WebGPU
```

Those belong downstream.

---

# 10. External AI Gateway Boundary

The intended security topology should be:

```text
External AI
    │
    │ untrusted semantic payload
    ▼
IntentCore Admission Boundary
    │
    ├── schema validation
    ├── identity verification
    ├── signature verification
    ├── policy admission
    ├── lifecycle
    └── trace/proof
    │
    ▼
Aetherium
    │
    ├── semantic interpretation
    ├── manifestation contract
    └── renderer
```

The external model must never obtain:

```text
GPU access
renderer access
shader execution authority
direct Visual State mutation authority
```

---

# 11. What Should NOT Be Added to IntentCore

Do not add these to the IntentCore kernel merely to integrate Aetherium:

- Canvas2D code
- WebGPU code
- WGSL
- SDF implementation
- particle simulation
- visual color mapping
- manifestation morphology
- Temporal Signal Fusion implementation
- browser DOM logic
- UI components

Those belong to Aetherium or its outer interaction layer.

---

# 12. Integration Test Matrix

| Test | Expected |
|---|---|
| Valid Aetherium message | Accepted and delivered |
| Missing signature | Rejected before Aetherium |
| Invalid signature | Rejected before Aetherium |
| Unknown payload schema | Rejected |
| Unsupported version | Rejected |
| Valid payload with renderer fields | Renderer fields ignored/rejected at ingress |
| Duplicate message ID | Explicit idempotency behavior |
| Trace propagation | Stable trace mapping |
| Policy rejection | No manifestation |
| Renderer failure | Does not mutate IntentCore authority state |
| Aetherium unavailable | Explicit integration failure |
| Replay | Same accepted payload yields same semantic input |
| Malformed payload | No renderer invocation |
| Timeout | Bounded failure |
| Unauthorized source | Rejected |

---

# 13. Recommended Development Sequence

### Step 1 — Contract

Define:

```text
AetheriumIntegrationEnvelope v0.1
```

### Step 2 — IntentCore

Implement:

```text
payload schema validation
signature verification
trace propagation
typed application payload metadata
```

### Step 3 — Aetherium

Implement:

```text
integrations/intentcore/
```

with no renderer dependency.

### Step 4 — Integration test

Prove:

```text
IntentCore
    ↓
accepted Aetherium message
    ↓
Aetherium interpreter
    ↓
Visual State Contract
```

### Step 5 — Only then

Consider:

```text
IntentCore
    ↓
AETH
    ↓
Presence IR
    ↓
Governor
```

That is a Phase-1 evolution, not a prerequisite for the current Phase-0 visual prototype.

---

# 14. Final Architectural Position

The cleanest division is:

```text
┌─────────────────────────────────────────────┐
│ External AI / Human / Client                │
└──────────────────────┬──────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────┐
│ IntentCore                                  │
│                                             │
│ Identity                                    │
│ Validation                                  │
│ Admission                                   │
│ Policy                                      │
│ Lifecycle                                   │
│ State / CAS                                 │
│ History / Proof / Telemetry                 │
└──────────────────────┬──────────────────────┘
                       │
                 Accepted Message
                       │
                       ▼
┌─────────────────────────────────────────────┐
│ Aetherium-Manifest                         │
│                                             │
│ Temporal Signal / Interaction Boundary      │
│ Semantic Interpretation                     │
│ AETH / Presence IR (future Phase 1)         │
│ Visual State Contract (current Phase 0)     │
│ Manifestation Runtime                       │
└──────────────────────┬──────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────┐
│ Renderer                                    │
│ Canvas2D / WebGPU                           │
└─────────────────────────────────────────────┘
```

### Core rule

> **IntentCore governs whether a message is allowed to enter and how it is coordinated. Aetherium governs how an accepted message becomes environmental manifestation. The renderer only executes downstream governed representation.**

This preserves the architectural separation already present in both repositories while creating a concrete integration boundary instead of merging their responsibilities.
