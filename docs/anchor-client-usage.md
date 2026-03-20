# Anchor Client Usage

## Purpose

The anchor submit client is the off-chain bridge that takes an L2 artifact hash and submits it
to the L1 anchor contract, returning an anchorId.

---

## Location

```
anchor-client/main.go
```

---

## Build

```bash
cd anchor-client
go build -o anchor-submit-client .
```

---

## Usage

```bash
./anchor-submit-client <artifact.json> [submitter]
```

| Argument | Required | Description |
|---|---|---|
| artifact.json | yes | Path to the canonical L2 state artifact JSON file |
| submitter | no | Identifier of the submitting node (default: anchor-client-v1) |

---

## What It Does

1. Reads and parses the artifact JSON into the canonical Artifact struct
2. Serialises the struct to compact JSON (deterministic field order)
3. Computes SHA-256 hash of the canonical bytes
4. Reads `parent_anchor_id` from the artifact for chain linkage
5. Calls `contract.CreateAnchor(artifactHash, parentHash, submitter)`
6. Prints the artifact hash and returned anchorId

---

## Example Run

Input artifact (`example-artifact.json`):
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

Command:
```bash
./anchor-submit-client example-artifact.json anchor-client-v1
```

Output (verified live):
```
artifact_hash : c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
anchor_id     : anchor-397529ce6400bfc4
submitter     : anchor-client-v1
chain_id      : bhiv-l2-app-001
block_height  : 1042
```

---

## anchor_contract.createAnchor() Mapping

| Client Parameter | Contract Parameter | Value |
|---|---|---|
| artifactHash | artifactHash | SHA-256 hex of canonical artifact |
| parentHash | parentHash | parent_anchor_id from artifact (empty for genesis) |
| submitter | submitter | caller identifier string |

Return value: `anchorId` — deterministic identifier derived from artifactHash + sequence number.

---

## Production Integration Note

In production the `anchor.Contract` in-process simulation is replaced by an RPC or HTTP call
to the deployed L1 anchor contract. The interface remains identical:

```
CreateAnchor(artifactHash, parentHash, submitter) → anchorId
GetAnchor(anchorId) → AnchorRecord
```
