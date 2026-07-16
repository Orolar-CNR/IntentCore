package proof

import (
	"context"
	"sync"

	"github.com/Orolar-CNR/IntentCore/contracts"
)

// InMemoryProofRecorder implements contracts.ProofRecorder for local testing.
type InMemoryProofRecorder struct {
	mu     sync.RWMutex
	proofs []contracts.Proof
}

// NewInMemoryProofRecorder creates a new proof recorder.
func NewInMemoryProofRecorder() *InMemoryProofRecorder {
	return &InMemoryProofRecorder{
		proofs: make([]contracts.Proof, 0),
	}
}

// RecordProof appends a proof to the in-memory log.
func (r *InMemoryProofRecorder) RecordProof(ctx context.Context, proof contracts.Proof) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.proofs = append(r.proofs, proof)
	return nil
}

// GetProofs returns all recorded proofs.
func (r *InMemoryProofRecorder) GetProofs() []contracts.Proof {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cpy := make([]contracts.Proof, len(r.proofs))
	copy(cpy, r.proofs)
	return cpy
}
