package proof_test

import (
	"context"
	"testing"

	"github.com/Orolar-CNR/IntentCore/contracts"
	"github.com/Orolar-CNR/IntentCore/proof"
)

func TestInMemoryProofRecorder(t *testing.T) {
	recorder := proof.NewInMemoryProofRecorder()

	p := contracts.Proof{
		IntentID: "test-intent",
		ProofID:  "proof-123",
	}

	err := recorder.RecordProof(context.Background(), p)
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	proofs := recorder.GetProofs()
	if len(proofs) != 1 {
		t.Fatalf("Expected 1 proof, got %d", len(proofs))
	}

	if proofs[0].IntentID != "test-intent" {
		t.Errorf("Expected IntentID test-intent, got %s", proofs[0].IntentID)
	}
}
