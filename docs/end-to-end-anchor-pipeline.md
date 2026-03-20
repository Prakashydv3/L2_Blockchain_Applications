# End-to-End Anchor Pipeline

## Purpose

Demonstrates the complete BHIV anchoring system from raw artifact JSON through to verified
on-chain commitment, using real computed values from all pipeline components.

---

## Full Pipeline

```
Artifact JSON
    │
    ▼
artifact-hash-generator
    │  SHA-256(canonical JSON)
    ▼
artifact_hash: c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
    │
    ▼
anchor-submit-client
    │  contract.CreateAnchor(artifactHash, parentHash, submitter)
    ▼
L1 Anchor Contract
    │  stores AnchorRecord, emits AnchorCommitted event
    ▼
anchor_id: anchor-397529ce6400bfc4
    │
    ▼
anchor-verifier
    │  GetAnchor(anchorId) → recompute hash → compare
    ▼
VERIFICATION PASSED
```

---

## Stage 1 — Artifact JSON

File: `artifact-hash-generator/example-artifact.json`

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

---

## Stage 2 — Hash Generator

```bash
go run artifact-hash-generator/main.go artifact-hash-generator/example-artifact.json
```

```
c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
```

Canonical JSON serialised by generator:
```
{"chain_id":"bhiv-l2-app-001","block_height":1042,"timestamp":"2025-01-15T10:00:00Z","application_state_root":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","registry_snapshot":"6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b","projection_log_root":"d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35","replay_proof_hash":"4e07408562bedb8b60ce05c1decb3f3b9b8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e","parent_anchor_id":""}
```

---

## Stage 3 — Anchor Client

```bash
go run anchor-client/main.go artifact-hash-generator/example-artifact.json anchor-client-v1
```

```
artifact_hash : c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
anchor_id     : anchor-397529ce6400bfc4
submitter     : anchor-client-v1
chain_id      : bhiv-l2-app-001
block_height  : 1042
```

---

## Stage 4 — Anchor Contract (L1 Event)

The L1 anchor contract stores the record and emits:

```
AnchorCommitted {
  anchorId     : anchor-397529ce6400bfc4
  artifactHash : c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
  parentHash   : ""
  submitter    : anchor-client-v1
}
```

---

## Stage 5 — Anchor Verification

```bash
go run anchor-client/anchor-verifier/main.go auto artifact-hash-generator/example-artifact.json
```

```
VERIFICATION PASSED
anchor_id     : anchor-397529ce6400bfc4
artifact_hash : c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
parent_hash   :
submitter     : verifier-seed
timestamp     : 2026-03-20T10:43:12Z
```

---

## Lineage Extension (3-Snapshot Chain)

```bash
go run anchor-client/lineage-sim/main.go
```

```
snapshot1  block 1042  →  anchor-397529ce6400bfc4  (genesis)
snapshot2  block 1087  →  anchor-229a802f0c560b5f  (parent: anchor-397529ce6400bfc4)
snapshot3  block 1134  →  anchor-6262a56e88e53761  (parent: anchor-229a802f0c560b5f)

Lineage intact: all parent references verified.
```

---

## Component Map

| Component | File | Role |
|---|---|---|
| Hash generator | `artifact-hash-generator/main.go` | Deterministic SHA-256 of artifact |
| Anchor contract | `anchor-client/anchor/contract.go` | L1 anchor store simulation |
| Submit client | `anchor-client/main.go` | Submits hash to contract |
| Verifier | `anchor-client/anchor-verifier/main.go` | Verifies on-chain record |
| Lineage sim | `anchor-client/lineage-sim/main.go` | Proves parent chain continuity |

---

## Proof Summary

| Artifact | Hash | AnchorId |
|---|---|---|
| example-artifact.json (block 1042) | `c2c63ac3...d1b2e` | `anchor-397529ce6400bfc4` |
| snapshot2 (block 1087) | `9a655a60...fa96` | `anchor-229a802f0c560b5f` |
| snapshot3 (block 1134) | `13de4972...e928` | `anchor-6262a56e88e53761` |
