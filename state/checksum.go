package state

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// computeChecksum returns the SHA-256 digest of the JSON encoding
// produced by encoding/json for the snapshot state.
// The checksum is intended to detect corruption or modification of
// snapshot data under the repository's defined serialization format.
func computeChecksum(data map[string]*stateEntry) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
