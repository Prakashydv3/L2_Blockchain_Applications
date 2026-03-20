package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

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

func main() {
	contract := anchor.NewContract()

	snapshots := []Artifact{
		{
			ChainID:              "bhiv-l2-app-001",
			BlockHeight:          1042,
			Timestamp:            "2025-01-15T10:00:00Z",
			ApplicationStateRoot: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			RegistrySnapshot:     "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
			ProjectionLogRoot:    "d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35",
			ReplayProofHash:      "4e07408562bedb8b60ce05c1decb3f3b9b8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e",
			ParentAnchorID:       "",
		},
		{
			ChainID:              "bhiv-l2-app-001",
			BlockHeight:          1087,
			Timestamp:            "2025-01-15T10:05:00Z",
			ApplicationStateRoot: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			RegistrySnapshot:     "1b4f0e9851971998e732078544c96b36c3d01cedf7caa332359d6f1d83567014",
			ProjectionLogRoot:    "60303ae22b998861bce3b28f33eec1be758a213c86c93c076dbe9f558c11c752",
			ReplayProofHash:      "fd61a03af4f77d870fc21e05e7e80678095c92d808cfb3b5c279ee04c74aca13",
			ParentAnchorID:       "", // filled after snapshot1 anchored
		},
		{
			ChainID:              "bhiv-l2-app-001",
			BlockHeight:          1134,
			Timestamp:            "2025-01-15T10:10:00Z",
			ApplicationStateRoot: "3f79bb7b435b05321651daefd374cdc681dc06faa65e374e38337b88ca046dea",
			RegistrySnapshot:     "a87ff679a2f3e71d9181a67b7542122c04521bad634b0a3e5c9b8e8e8e8e8e8e",
			ProjectionLogRoot:    "e4da3b7fbbce2345d7772b0674a318d5e0f2df5f5e8e8e8e8e8e8e8e8e8e8e8e",
			ReplayProofHash:      "1679091c5a880faf6fb5e6087eb1b2dc6b8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e",
			ParentAnchorID:       "", // filled after snapshot2 anchored
		},
	}

	fmt.Println("=== BHIV Anchor Lineage Simulation ===")
	fmt.Println()

	prevAnchorID := ""
	for i, snap := range snapshots {
		snap.ParentAnchorID = prevAnchorID
		h := hashArtifact(snap)

		anchorID, err := contract.CreateAnchor(h, prevAnchorID, "lineage-sim")
		if err != nil {
			fmt.Printf("ERROR snapshot%d: %v\n", i+1, err)
			return
		}

		fmt.Printf("snapshot%d\n", i+1)
		fmt.Printf("  block_height  : %d\n", snap.BlockHeight)
		fmt.Printf("  artifact_hash : %s\n", h)
		fmt.Printf("  parent_anchor : %s\n", prevAnchorID)
		fmt.Printf("  anchor_id     : %s\n", anchorID)
		fmt.Println()

		prevAnchorID = anchorID
	}

	fmt.Println("=== Lineage Verification ===")
	fmt.Println()

	// Walk the chain backwards to verify linkage
	current := prevAnchorID
	chain := []string{}
	for current != "" {
		rec, err := contract.GetAnchor(current)
		if err != nil {
			fmt.Printf("ERROR retrieving %s: %v\n", current, err)
			return
		}
		chain = append(chain, fmt.Sprintf("%s <- %s", rec.AnchorID, rec.ParentHash))
		current = rec.ParentHash
	}

	for _, link := range chain {
		fmt.Println(" ", link)
	}
	fmt.Println()
	fmt.Println("Lineage intact: all parent references verified.")
}
