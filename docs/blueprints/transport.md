# Transport Blueprint

**Phase:** Phase 0 — Specification & Blueprint
**Related RFCs:** RFC-0001 (Semantic Envelope)

## Purpose
This blueprint describes the implementation strategy for the transport boundary, specifically focusing on how high-performance technologies like eBPF and XDP are expected to be utilized without leaking into the core IntentCore logic.

## Transport Pipeline
The transport layer's sole responsibility is delivering the `SemanticEnvelope` into the kernel space.

```text
Network Interface (NIC)
        │
        ▼
XDP / eBPF (Hardware/Kernel space fast-path)
        │ (Drops invalid packets instantly)
        ▼
AF_XDP Socket / Ring Buffer
        │
        ▼
Go Userspace (Decoder / Deserializer)
        │
        ▼
SemanticEnvelope
        │
        ▼
IntentCore Kernel (Validation Layer)
```

## Hardware Acceleration (eBPF / XDP)
IntentCore is designed for extreme throughput. The reference implementation uses eBPF (Extended Berkeley Packet Filter) and XDP (eXpress Data Path).
*   **Early Drop:** Malformed packets, invalid magic bytes, or unsupported protocol versions are expected to be dropped directly in the Linux kernel network stack (or even on the NIC) before they consume userspace CPU cycles.
*   **Zero-copy:** By utilizing AF_XDP, packet data can be shared directly with the Go userspace application, minimizing memory copies.

## Protocol Independence
While ABTP (AetherBus Transport Protocol) over UDP/XDP is the reference, the architecture strictly decouples transport from intent processing. The transport adapter in Go merely wraps the network I/O, ensuring that if a future implementation uses DPDK, RDMA, or standard TCP sockets, the core kernel remains completely untouched.
