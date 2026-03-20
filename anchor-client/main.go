package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bhiv/anchor-client/anchor"
)

// Artifact matches the canonical L2 state root artifact structure.
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

func hashArtifact(a Artifact) string {
	b, _ := json.Marshal(a)
	sum := sha256.Sum256(b)
	return fmt.Sprintf("%x", sum)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: anchor-submit-client <artifact.json> [submitter]")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}

	var artifact Artifact
	if err := json.Unmarshal(data, &artifact); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing artifact: %v\n", err)
		os.Exit(1)
	}

	submitter := "anchor-client-v1"
	if len(os.Args) >= 3 {
		submitter = os.Args[2]
	}

	artifactHash := hashArtifact(artifact)
	parentHash := artifact.ParentAnchorID

	contract := anchor.NewContract()

	// If parentHash is an anchorId reference, it must already exist in the contract.
	// For simulation, we seed the parent if provided so linkage can be verified.
	if parentHash != "" {
		contract.SeedAnchor(parentHash)
	}

	anchorID, err := contract.CreateAnchor(artifactHash, parentHash, submitter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "anchor submission failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("artifact_hash : %s\n", artifactHash)
	fmt.Printf("anchor_id     : %s\n", anchorID)
	fmt.Printf("submitter     : %s\n", submitter)
	fmt.Printf("chain_id      : %s\n", artifact.ChainID)
	fmt.Printf("block_height  : %d\n", artifact.BlockHeight)
}
