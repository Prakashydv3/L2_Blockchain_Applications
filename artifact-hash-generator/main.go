package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
)

// Artifact defines the canonical L2 state root artifact.
// Field order here matches the canonical serialisation order.
type Artifact struct {
	ChainID              string `json:"chain_id"`
	BlockHeight          uint64 `json:"block_height"`
	Timestamp            string `json:"timestamp"`
	ApplicationStateRoot string `json:"application_state_root"`
	RegistrySnapshot     string `json:"registry_snapshot"`
	ProjectionLogRoot    string `json:"projection_log_root"`
	ReplayProofHash      string `json:"replay_proof_hash"`
	ParentAnchorID       string `json:"parent_anchor_id"`
}

func hashArtifact(a Artifact) (string, error) {
	canonical, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canonical)
	return fmt.Sprintf("%x", sum), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: artifact-hash-generator <artifact.json>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}

	var artifact Artifact
	if err := json.Unmarshal(data, &artifact); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing artifact JSON: %v\n", err)
		os.Exit(1)
	}

	hash, err := hashArtifact(artifact)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error hashing artifact: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(hash)
}
