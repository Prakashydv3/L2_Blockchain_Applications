package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bhiv/anchor-client/anchor"
)

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

func verify(contract *anchor.Contract, anchorID string, artifact Artifact) error {
	rec, err := contract.GetAnchor(anchorID)
	if err != nil {
		return fmt.Errorf("anchor retrieval failed: %w", err)
	}
	expectedHash := hashArtifact(artifact)
	if rec.ArtifactHash != expectedHash {
		return fmt.Errorf("hash mismatch:\n  on-chain : %s\n  computed : %s", rec.ArtifactHash, expectedHash)
	}
	if rec.ParentHash != artifact.ParentAnchorID {
		return fmt.Errorf("parentHash mismatch:\n  on-chain : %s\n  expected : %s", rec.ParentHash, artifact.ParentAnchorID)
	}
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: anchor-verifier <anchor-id|auto> <artifact.json>")
		os.Exit(1)
	}

	anchorID := os.Args[1]

	data, err := os.ReadFile(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading artifact: %v\n", err)
		os.Exit(1)
	}

	var artifact Artifact
	if err := json.Unmarshal(data, &artifact); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing artifact: %v\n", err)
		os.Exit(1)
	}

	contract := anchor.NewContract()
	artifactHash := hashArtifact(artifact)
	parentHash := artifact.ParentAnchorID

	if parentHash != "" {
		contract.SeedAnchor(parentHash)
	}

	seededID, err := contract.CreateAnchor(artifactHash, parentHash, "verifier-seed")
	if err != nil {
		fmt.Fprintf(os.Stderr, "seed failed: %v\n", err)
		os.Exit(1)
	}

	if anchorID == "auto" {
		anchorID = seededID
	}

	if err := verify(contract, anchorID, artifact); err != nil {
		fmt.Printf("VERIFICATION FAILED: %v\n", err)
		os.Exit(1)
	}

	rec, _ := contract.GetAnchor(anchorID)
	fmt.Printf("VERIFICATION PASSED\n")
	fmt.Printf("anchor_id     : %s\n", rec.AnchorID)
	fmt.Printf("artifact_hash : %s\n", rec.ArtifactHash)
	fmt.Printf("parent_hash   : %s\n", rec.ParentHash)
	fmt.Printf("submitter     : %s\n", rec.Submitter)
	fmt.Printf("timestamp     : %s\n", rec.Timestamp.Format("2006-01-02T15:04:05Z"))
}
