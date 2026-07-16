package contracts

import (
	"context"
)

// Proof represents cryptographic or auditable evidence of a system event or state.
type Proof struct {
	ProofID   string
	IntentID  string
	StateHash string
	Signature []byte
}

// ProofRecorder defines an interface for persistently recording cryptographic proofs.
// This is strictly an OPTIONAL dependency for the Runtime Kernel.
type ProofRecorder interface {
	RecordProof(ctx context.Context, proof Proof) error
}
