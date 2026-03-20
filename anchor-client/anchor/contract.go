package anchor

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
	"time"
)

// AnchorRecord mirrors what the L1 anchor contract stores per commitment.
type AnchorRecord struct {
	AnchorID     string
	ArtifactHash string
	ParentHash   string
	Submitter    string
	Timestamp    time.Time
}

// Contract is an in-process simulation of the L1 anchor contract.
// In production this is replaced by an actual on-chain contract call.
type Contract struct {
	mu      sync.RWMutex
	store   map[string]AnchorRecord
	counter int
}

func NewContract() *Contract {
	return &Contract{store: make(map[string]AnchorRecord)}
}

// CreateAnchor stores a new anchor and returns its anchorId.
// Mirrors anchor_contract.createAnchor(artifactHash, parentHash, submitter).
func (c *Contract) CreateAnchor(artifactHash, parentHash, submitter string) (string, error) {
	if artifactHash == "" {
		return "", errors.New("artifactHash must not be empty")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if parentHash != "" {
		if _, ok := c.store[parentHash]; !ok {
			return "", fmt.Errorf("parentHash %q not found in contract", parentHash)
		}
	}

	c.counter++
	anchorID := deriveAnchorID(artifactHash, c.counter)

	c.store[anchorID] = AnchorRecord{
		AnchorID:     anchorID,
		ArtifactHash: artifactHash,
		ParentHash:   parentHash,
		Submitter:    submitter,
		Timestamp:    time.Now().UTC(),
	}

	return anchorID, nil
}

// GetAnchor retrieves an anchor record by anchorId.
func (c *Contract) GetAnchor(anchorID string) (AnchorRecord, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	rec, ok := c.store[anchorID]
	if !ok {
		return AnchorRecord{}, fmt.Errorf("anchor %q not found", anchorID)
	}
	return rec, nil
}

// SeedAnchor inserts a pre-existing anchor record used in simulation to represent
// an anchor that was committed in a prior session.
func (c *Contract) SeedAnchor(anchorID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.store[anchorID]; !ok {
		c.store[anchorID] = AnchorRecord{AnchorID: anchorID}
	}
}

func deriveAnchorID(artifactHash string, seq int) string {
	raw := fmt.Sprintf("%s:%d", artifactHash, seq)
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("anchor-%x", sum[:8])
}
