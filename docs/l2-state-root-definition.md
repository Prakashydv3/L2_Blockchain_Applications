# L2 State Root Definition

## Purpose

Defines exactly what an L2 chain commits to L1. Every anchor on L1 references one state root.
A state root is a deterministic hash of a structured artifact that captures the full application
state at a specific point in time.

---

## State Root Format

A state root is a hex-encoded SHA-256 hash:

```
"state_root": "a3f1c2d4e5b6789012345678abcdef01234567890abcdef1234567890abcdef12"
```

Length: 64 hex characters (32 bytes)
Algorithm: SHA-256
Encoding: lowercase hex, no 0x prefix

---

## Artifact Structure

The artifact is the canonical JSON document that is hashed to produce the state root.
Fields must be serialised in the exact order defined below. No extra whitespace. No trailing commas.

### Required Fields

| Field | Type | Description |
|---|---|---|
| chain_id | string | Identifier of the L2 chain producing this artifact |
| block_height | uint64 | L2 block height at time of snapshot |
| timestamp | string | ISO 8601 UTC timestamp of snapshot |
| application_state_root | string | Root hash of the application state trie |
| registry_snapshot | string | Hash of the current registry state |
| projection_log_root | string | Root hash of the projection event log |
| replay_proof_hash | string | Hash of the replay proof for this state transition |
| parent_anchor_id | string | anchorId of the previous anchor (empty string for genesis) |

---

## Canonical JSON Example

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

Genesis anchor has `parent_anchor_id` set to empty string `""`.

---

## Subsequent Snapshot Example (with parent linkage)

```json
{
  "chain_id": "bhiv-l2-app-001",
  "block_height": 1087,
  "timestamp": "2025-01-15T10:05:00Z",
  "application_state_root": "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
  "registry_snapshot": "1b4f0e9851971998e732078544c96b36c3d01cedf7caa332359d6f1d83567014",
  "projection_log_root": "60303ae22b998861bce3b28f33eec1be758a213c86c93c076dbe9f558c11c752",
  "replay_proof_hash": "fd61a03af4f77d870fc21e05e7e80678095c92d808cfb3b5c279ee04c74aca13",
  "parent_anchor_id": "anchor-0001"
}
```

---

## Hashing Rules

1. Serialise the artifact as compact JSON — no extra spaces, no newlines
2. Field order must follow the canonical order defined in the table above exactly
3. Apply SHA-256 to the UTF-8 encoded bytes of the compact JSON string
4. Encode the result as lowercase hex

Any deviation in field order, whitespace, or encoding produces a different hash.
The artifact-hash-generator enforces these rules.

---

## Field Definitions

### application_state_root
The Merkle root of the L2 application state trie at the snapshot block height.
Produced by the L2 chain's state commitment module.

### registry_snapshot
SHA-256 hash of the serialised registry state (all registered entities and their current values).

### projection_log_root
Merkle root of the ordered projection event log up to and including the snapshot block.

### replay_proof_hash
SHA-256 hash of the replay proof bundle. The replay proof allows any verifier to reconstruct
the state transition from the previous snapshot to this one.

### parent_anchor_id
The anchorId returned by the L1 anchor contract for the immediately preceding anchor.
Creates a verifiable chain of custody from genesis to current state.
Set to empty string `""` for the first (genesis) anchor.
