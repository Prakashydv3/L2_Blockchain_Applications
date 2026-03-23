# REVIEW_PACKET.md
## BHIV L1/L2 Anchoring Bridge

---

## Phase 1 — Entry Points

**Backend Entry Point 1 — Hash Generator**
Path: `artifact-hash-generator/main.go`
Purpose: Accepts artifact JSON, produces deterministic SHA-256 state root hash

**Backend Entry Point 2 — Anchor Submit Client**
Path: `anchor-client/main.go`
Purpose: Accepts artifact JSON, submits hash to L1 anchor contract, returns anchorId

---

## Phase 2 — Core Execution Files (3 Files)

**File 1 — artifact-hash-generator/main.go**
Parses artifact JSON into a typed struct, serialises to canonical compact JSON, applies SHA-256.
Output: 64-char lowercase hex hash.

**File 2 — anchor-client/main.go**
Reads artifact, computes hash, calls `contract.CreateAnchor(artifactHash, parentHash, submitter)`.
Output: anchorId returned from L1 contract.

**File 3 — anchor-client/anchor/contract.go**
In-process simulation of the L1 anchor contract. Stores AnchorRecords, validates parentHash linkage, derives deterministic anchorIds.
Output: AnchorRecord with anchorId, artifactHash, parentHash, submitter, timestamp.

---

## Phase 3 — Live Execution Flow

User action: L2 chain produces state snapshot at block 1042

```
artifact JSON (example-artifact.json)
  → artifact-hash-generator/main.go       [SHA-256 canonical hash]
  → anchor-client/main.go                 [CreateAnchor call]
  → anchor/contract.go                    [stores record, returns anchorId]
  → anchor-verifier/main.go               [retrieves record, recomputes hash, compares]
  → VERIFICATION PASSED
```

---

## Phase 4 — Real Output Proof

**Hash Generator:**
```
c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
```

**Anchor Submit Client:**
```
artifact_hash : c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
anchor_id     : anchor-397529ce6400bfc4
submitter     : anchor-client-v1
chain_id      : bhiv-l2-app-001
block_height  : 1042
```

**Anchor Verifier:**
```
VERIFICATION PASSED
anchor_id     : anchor-397529ce6400bfc4
artifact_hash : c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
parent_hash   :
submitter     : verifier-seed
timestamp     : 2026-03-23T04:28:06Z
```

**Lineage Simulation:**
```
snapshot1  block 1042  artifact_hash: c2c63ac3...  anchor_id: anchor-397529ce6400bfc4
snapshot2  block 1087  artifact_hash: 9a655a60...  anchor_id: anchor-229a802f0c560b5f  parent: anchor-397529ce6400bfc4
snapshot3  block 1134  artifact_hash: 13de4972...  anchor_id: anchor-6262a56e88e53761  parent: anchor-229a802f0c560b5f

anchor-6262a56e88e53761 <- anchor-229a802f0c560b5f <- anchor-397529ce6400bfc4 <- (genesis)

Lineage intact: all parent references verified.
```

---

## Phase 5 — Task Contribution Summary

**Built:**
- `artifact-hash-generator/main.go` — deterministic SHA-256 artifact hashing utility
- `anchor-client/main.go` — off-chain anchor submit client
- `anchor-client/anchor/contract.go` — L1 anchor contract simulation
- `anchor-client/anchor-verifier/main.go` — anchor integrity verifier
- `anchor-client/lineage-sim/main.go` — 3-snapshot parent chain continuity simulation
- All 10 documentation files across Phases 1–3

**Modified:**
- Nothing modified — greenfield build

**Did NOT touch:**
- L1 consensus logic
- L1 block validation logic
- L1 transaction execution logic
- Any governance logic

---

## Phase 6 — Failure Cases

| Scenario | System Behavior |
|---|---|
| Empty `artifactHash` submitted | `contract.CreateAnchor` returns error: `artifactHash must not be empty` — submission rejected |
| `parentHash` references non-existent anchorId | `contract.CreateAnchor` returns error: `parentHash not found in contract` — submission rejected |
| Artifact JSON missing required field | Go struct zero-value fills field — hash is deterministic but semantically invalid — verifier will catch mismatch |
| Artifact tampered after anchoring | Verifier recomputes hash, detects mismatch — outputs `VERIFICATION FAILED: hash mismatch` |
| Wrong `parent_anchor_id` in artifact | Verifier detects parentHash mismatch — outputs `VERIFICATION FAILED: parentHash mismatch` |
| anchorId not found on chain | Verifier outputs `VERIFICATION FAILED: anchor retrieval failed` |

---

## Phase 7 — Proof of Execution

All four tools executed live. Console output:

```
$ go run artifact-hash-generator/main.go example-artifact.json
c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e

$ go run anchor-client/main.go example-artifact.json anchor-client-v1
artifact_hash : c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
anchor_id     : anchor-397529ce6400bfc4
submitter     : anchor-client-v1
chain_id      : bhiv-l2-app-001
block_height  : 1042

$ go run anchor-client/anchor-verifier/main.go auto example-artifact.json
VERIFICATION PASSED
anchor_id     : anchor-397529ce6400bfc4
artifact_hash : c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
parent_hash   :
submitter     : verifier-seed
timestamp     : 2026-03-23T04:28:06Z

$ go run anchor-client/lineage-sim/main.go
=== BHIV Anchor Lineage Simulation ===

snapshot1
  block_height  : 1042
  artifact_hash : c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
  parent_anchor :
  anchor_id     : anchor-397529ce6400bfc4

snapshot2
  block_height  : 1087
  artifact_hash : 9a655a600a6b51bd2b208578a257f9f5676f8bfe04038e7844e9a70673bafa96
  parent_anchor : anchor-397529ce6400bfc4
  anchor_id     : anchor-229a802f0c560b5f

snapshot3
  block_height  : 1134
  artifact_hash : 13de4972e2b0b633f945ca943607432c8b848a4c7ccbf641d8213074e23de928
  parent_anchor : anchor-229a802f0c560b5f
  anchor_id     : anchor-6262a56e88e53761

=== Lineage Verification ===

  anchor-6262a56e88e53761 <- anchor-229a802f0c560b5f
  anchor-229a802f0c560b5f <- anchor-397529ce6400bfc4
  anchor-397529ce6400bfc4 <-

Lineage intact: all parent references verified.
```

Both modules build with zero errors:
```
$ go build ./...   [artifact-hash-generator]  →  BUILD_OK
$ go build ./...   [anchor-client]             →  BUILD_OK
```
