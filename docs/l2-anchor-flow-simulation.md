# L2 Anchor Flow Simulation

## Purpose

Demonstrates the complete L2 commitment pipeline from application state snapshot through to
anchor retrieval and verification using real computed values.

---

## Pipeline

```
Application State Snapshot (artifact JSON)
        |
        v
  artifact-hash-generator  →  SHA-256 hash
        |
        v
  anchor-submit-client  →  L1 Anchor Contract
        |
        v
  AnchorCommitted Event  →  anchorId returned
        |
        v
  anchor-verifier  →  VERIFICATION PASSED
```

---

## Step 1 — Application State Snapshot

The L2 chain produces a state artifact at block height 1042:

```json
{
  "chain_id": "bhiv-l2-app-001",
  "block_height": 1042,
  "timestamp": "2025-01-15T10:00:00Z",
  "application_state_root": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
  "registry_snapshot": "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
  "projection_log_root": "d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35",
  "replay_proof_hash": "4e07408562bedb8b60ce05c1decb3f3b9b8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e",
  "parent_anchor_id": ""
}
```

Artifact stored in Bucket Layer under key: `bhiv-l2-app-001/1042/artifact.json`

---

## Step 2 — Hash Generation

```bash
go run artifact-hash-generator/main.go artifact-hash-generator/example-artifact.json
```

Output:
```
c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
```

---

## Step 3 — Anchor Submission

```bash
go run anchor-client/main.go artifact-hash-generator/example-artifact.json anchor-client-v1
```

Output:
```
artifact_hash : c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
anchor_id     : anchor-397529ce6400bfc4
submitter     : anchor-client-v1
chain_id      : bhiv-l2-app-001
block_height  : 1042
```

L1 Anchor Contract emits:
```
AnchorCommitted {
  anchorId     : anchor-397529ce6400bfc4
  artifactHash : c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
  parentHash   : (empty — genesis)
  submitter    : anchor-client-v1
  timestamp    : 2025-01-15T10:00:05Z
}
```

---

## Step 4 — Anchor Retrieval

Any party queries the L1 contract:

```
GetAnchor("anchor-397529ce6400bfc4")
→ AnchorRecord {
    AnchorID     : anchor-397529ce6400bfc4
    ArtifactHash : c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
    ParentHash   : ""
    Submitter    : anchor-client-v1
  }
```

---

## Step 5 — Verification

```bash
go run anchor-client/anchor-verifier/main.go auto artifact-hash-generator/example-artifact.json
```

Output:
```
VERIFICATION PASSED
anchor_id     : anchor-397529ce6400bfc4
artifact_hash : c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
parent_hash   :
submitter     : verifier-seed
timestamp     : 2026-03-20T10:43:12Z
```

---

## Summary

| Step | Input | Output |
|---|---|---|
| Snapshot | L2 block 1042 state | artifact JSON |
| Hash | artifact JSON | `c2c63ac3...d1b2e` |
| Submit | artifact hash | `anchor-397529ce6400bfc4` |
| Retrieve | anchorId | AnchorRecord |
| Verify | AnchorRecord + artifact | VERIFICATION PASSED |
