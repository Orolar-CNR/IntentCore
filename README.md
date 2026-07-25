# IntentCore

A language-agnostic intent coordination kernel that provides deterministic lifecycle management, state consistency, policy enforcement, and proof-oriented coordination for distributed autonomous systems.

IntentCore uses AetherBus as its transport layer while keeping validation, normalization, admission, lifecycle control, state mutation, history, proof, and telemetry inside the kernel.

## Architectural Identity

The project is intentionally split into four stable names:

| Name | Architectural role |
| --- | --- |
| IntentCore | Kernel and center of system coordination |
| AetherBus | Transport layer and protocol boundary |
| SemanticEnvelope | Wire format carried by AetherBus |
| RFC | Locked contract for stable implementation behavior |

This naming model prevents the early-project ambiguity where AetherBus could mean the project, the broker, the transport, or the protocol. IntentCore is now the coordination kernel; AetherBus is the transport protocol that carries semantic messages into it.

## Architecture Blueprint

The repository blueprint is maintained in [IntentCore Architecture Landscape](docs/architecture-landscape.md). It documents the architecture boundaries, frozen core contracts, implementation phases, internal data flow, architectural principles, historical naming context, and next-phase development direction.
