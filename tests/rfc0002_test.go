package tests

import (
	"context"
	"testing"
	"time"

	"github.com/Orolar-CNR/IntentCore/admission"
	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/core"
	"github.com/google/uuid"
)

// TestRFC0002DeterministicPolicies verifies that envelopes violating deterministic
// constraints are rejected with the proper code and reason, testing compliance with RFC-0002.
func TestRFC0002DeterministicPolicies(t *testing.T) {
	evaluator := admission.NewPolicyEvaluator()

	validEnvelope := func() contracts.SemanticEnvelope {
		return contracts.SemanticEnvelope{
			EnvelopeID:     core.IntentID(uuid.New()),
			SchemaVersion:  admission.DefaultSchemaVersion,
			AgentIdentity:  "agent-007",
			Signatures:     []contracts.Signature{{Signature: []byte("sig")}},
			OpaquePayload:  []byte(`{"action":"test"}`),
			EventTimestamp: time.Now(),
		}
	}

	tests := []struct {
		name         string
		mutateEnv    func(env *contracts.SemanticEnvelope)
		expectReject bool
		expectedCode contracts.RejectionCode
	}{
		{
			name:         "Valid Envelope",
			mutateEnv:    func(env *contracts.SemanticEnvelope) {},
			expectReject: false,
		},
		{
			name: "Invalid Schema Version",
			mutateEnv: func(env *contracts.SemanticEnvelope) {
				env.SchemaVersion = "99.9.9"
			},
			expectReject: true,
			expectedCode: contracts.CodeInvalidVersion,
		},
		{
			name: "Missing IntentID",
			mutateEnv: func(env *contracts.SemanticEnvelope) {
				env.EnvelopeID = core.IntentID(uuid.Nil)
			},
			expectReject: true,
			expectedCode: contracts.CodeMissingIntentID,
		},
		{
			name: "Missing AgentIdentity",
			mutateEnv: func(env *contracts.SemanticEnvelope) {
				env.AgentIdentity = ""
			},
			expectReject: true,
			expectedCode: contracts.CodeMissingIdentity,
		},
		{
			name: "Missing Signature",
			mutateEnv: func(env *contracts.SemanticEnvelope) {
				env.Signatures = nil
			},
			expectReject: true,
			expectedCode: contracts.CodeMissingSignature,
		},
		{
			name: "Missing Payload",
			mutateEnv: func(env *contracts.SemanticEnvelope) {
				env.OpaquePayload = nil
			},
			expectReject: true,
			expectedCode: contracts.CodeMissingPayload,
		},
		{
			name: "Missing Timestamp",
			mutateEnv: func(env *contracts.SemanticEnvelope) {
				env.EventTimestamp = time.Time{}
			},
			expectReject: true,
			expectedCode: contracts.CodeInvalidTimestamp,
		},
	}

	ctx := context.Background()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			env := validEnvelope()
			tc.mutateEnv(&env)

			res, err := evaluator.Evaluate(ctx, env)
			if err != nil {
				t.Fatalf("Unexpected error from evaluation: %v", err)
			}

			if tc.expectReject {
				if res.Decision != contracts.DecisionReject {
					t.Errorf("Expected rejection but was accepted")
				}
				if res.Evidence.Code != tc.expectedCode {
					t.Errorf("Expected code %s, got %s", tc.expectedCode, res.Evidence.Code)
				}
				if res.Error == nil {
					t.Errorf("Expected Error to be populated on rejection")
				}
			} else {
				if res.Decision != contracts.DecisionAccept {
					t.Errorf("Expected acceptance but was rejected: %v", res.Evidence.Reason)
				}
			}
		})
	}
}
