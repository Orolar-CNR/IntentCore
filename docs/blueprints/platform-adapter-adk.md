# Blueprint: ADK Platform Adapter

**Status:** Draft
**Category:** Blueprint
**Purpose:** Informative architectural reference
**Scope:** Integration of Google's Agent Development Kit (ADK) as a Platform Adapter for IntentCore

---

## 1. Executive Summary
This document specifies the integration architecture for using **Google Agent Development Kit (ADK) for Go** as a Platform Adapter for **IntentCore**.

To participate in hyper-scale, cross-network agent ecosystems, IntentCore relies on Platform Adapters to handle external transport abstraction. By implementing an ADK Adapter, IntentCore can receive standardized open protocols such as **Agent-to-Agent (A2A)** and **Model Context Protocol (MCP)**.

The ADK Adapter strictly adheres to the Platform Adapter Architecture. It acts exclusively as an outer integration shell, translating network signals into canonical `SemanticEnvelope`s and pushing them through the standard IntentCore runtime pipeline. It does not bypass any kernel stages or mutate state directly.

```text
External Agent Network
        │
        ▼
   ADK (A2A / MCP)
        │
        ▼
 Platform Adapter (ADK Implementation)
        │
        ▼
  SemanticEnvelope
        │
        ▼
    IntentCore
 (Pipeline -> Validation -> Admission -> Dispatcher -> Lifecycle -> Ledger)
```

## 2. Architectural Protocol Mapping
To establish alignment, inbound protocol signals received via ADK are converted directly into strongly-typed `SemanticEnvelope` structures before entering the pipeline.

| External Protocol Event | ADK Go Abstraction | IntentCore Core Mapping |
| :--- | :--- | :--- |
| **A2A Task Delegation** | `adk.AgentMessage` | Translated to `SemanticEnvelope` -> Injected into Pipeline |
| **MCP Tool Call** | `adk.ToolRequest` | Translated to `SemanticEnvelope` -> Injected into Pipeline |
| **State Sync Signal** | `adk.SessionState` | Translated to `SemanticEnvelope` -> Injected into Pipeline |

*Note: The ADK Adapter MUST NOT interact directly with the `CommitmentLedger` or the `AdmissionPolicy`. All external state syncs or tool requests must be packaged as envelopes and subjected to standard pipeline processing.*

## 3. Go Adapter Implementation Specification

The integration architecture utilizes the Adapter Structural Pattern. The adapter encapsulates the IntentCore Runtime Pipeline and exposes it as a node within an ADK Go distributed topology.

```go
package adk_adapter

import (
	"context"
	"fmt"
	"time"

	// Mocking references to the internal system specification for clear compilation context
	// "github.com/Orolar-CNR/IntentCore/contracts" // Replace with actual module path if different
	// "github.com/Orolar-CNR/IntentCore/core"      // Assuming a core package holds these definitions
)

// ADKIntentAdapter bridges Google's Go ADK ecosystem and the IntentCore Core Runtime.
type ADKIntentAdapter struct {
	pipeline contracts.Pipeline // The entry point into the IntentCore kernel
	publisher contracts.EventPublisher // Interface for broadcasting outbound events
}

// NewADKIntentAdapter initializes the cross-network adapter.
func NewADKIntentAdapter(pipeline contracts.Pipeline, publisher contracts.EventPublisher) *ADKIntentAdapter {
	return &ADKIntentAdapter{
		pipeline:  pipeline,
		publisher: publisher,
	}
}

// HandleInboundA2AMessage processes external Agent-to-Agent (A2A) protocol envelopes
// and translates them into language-agnostic IntentCore SemanticEnvelopes.
func (a *ADKIntentAdapter) HandleInboundA2AMessage(ctx context.Context, senderID string, intentType string, route string, payload []byte) error {
	fmt.Printf("[ADK Adapter]: Intercepted network A2A message from Agent: %s\n", senderID)

	// Translate network wire variables into strongly-typed canonical core definitions
	envelope := core.SemanticEnvelope{
		TraceID:      fmt.Sprintf("TR-ADK-%d", time.Now().UnixNano()),
		IntentID:     fmt.Sprintf("INTENT-%s-%d", intentType, time.Now().Unix()),
		PolicyDomain: core.NamespaceAgent, // Automatically bound via Namespace Enum
		RoutePath:    route,
		Priority:     128, // Default balanced priority for network intents
		DeadlineNano: uint64(time.Now().Add(15 * time.Second).UnixNano()),
		QoSLevel:     core.QoSAtLeastOnce, // Enforce network acknowledgment matrix
		PayloadSize:  len(payload),
		Payload:      payload,
	}

	// Dispatch message into the canonical Pipeline.
	// The Pipeline handles Validation -> Normalization -> Admission -> Dispatcher -> Lifecycle.
	err := a.pipeline.Execute(ctx, envelope)
	if err != nil {
		return fmt.Errorf("pipeline execution failed for intent %s: %w", envelope.IntentID, err)
	}

	fmt.Printf("[ADK Adapter]: Successfully injected external intent %s into Pipeline\n", envelope.IntentID)
	return nil
}
```

## 4. Operational Pipeline Phase Alignment

### Phase 1: Platform Adapter
* Bind network transports (e.g., **ZMQ ROUTER/DEALER**) to ADK's network transport abstractions.
* Ensure serialization blocks cleanly pass raw `[]byte` arrays without modifying the payload context.

### Phase 2: Semantic Translation
* Convert external protocol signals (`adk.AgentMessage`, `adk.ToolRequest`) into IntentCore `SemanticEnvelope` instances.
* Map external routing IDs to IntentCore namespaces.

### Phase 3: Pipeline Integration
* Inject the constructed `SemanticEnvelope` into `Pipeline.Execute()`.
* The IntentCore kernel natively handles the Validation, Normalization, Admission, and Lifecycle phases. The adapter waits for pipeline execution results.

### Phase 4: Outbound Publisher & Federation
* The `EventPublisher` interface is utilized by the adapter to listen to lifecycle events emitted by the kernel.
* The adapter translates canonical outbound events back into ADK protocol formats for Cross-Network Routing and Federation.

## 5. Verification Checklist for Developers
- [ ] Verify that `adk-go` module dependencies are pinned in `go.mod`.
- [ ] Ensure the adapter strictly uses `Pipeline.Execute()` and does NOT import or call `Validation`, `AdmissionPolicy`, or `Ledger` components directly.
- [ ] Confirm all `DeadlineNano` conversions accurately handle Unix Epoch time formats across differing languages.
- [ ] Validate that the outbound Event Publisher accurately translates core kernel state changes into network-safe broadcast messages.
