# L2 Integration Guidelines

## Purpose

Defines how BHIV L2 application chains integrate with the L1 anchoring pipeline. This document
is the reference for any L2 chain operator or application developer connecting to the anchor system.

---

## 1. How Applications Produce State Roots

An L2 application produces a state root by assembling a canonical artifact at a defined
commitment point (per block, per epoch, or per application event).

The artifact must contain exactly these fields in this order:

```json
{
  "chain_id": "<unique L2 chain identifier>",
  "block_height": <uint64 — L2 block height at snapshot>,
  "timestamp": "<ISO 8601 UTC — time of snapshot>",
  "application_state_root": "<hex — Merkle root of application state trie>",
  "registry_snapshot": "<hex — SHA-256 of serialised registry state>",
  "projection_log_root": "<hex — Merkle root of projection event log>",
  "replay_proof_hash": "<hex — SHA-256 of replay proof bundle>",
  "parent_anchor_id": "<anchorId of previous anchor, or empty string for genesis>"
}
```

The L2 chain is responsible for computing each sub-hash correctly before assembling the artifact.
The anchor system does not validate sub-hash contents — it only hashes and anchors the artifact.

---

## 2. How Roots Are Hashed

The artifact is hashed using the artifact-hash-generator:

```bash
go run artifact-hash-generator/main.go <artifact.json>
```

Rules:
- The artifact is parsed into a typed Go struct (canonical field order enforced)
- Serialised to compact JSON — no spaces, no newlines
- SHA-256 applied to UTF-8 bytes
- Output: 64-character lowercase hex string

The resulting hash is the state root submitted to L1.

Do not hash raw JSON strings directly — always parse through the typed struct to guarantee
field order and eliminate whitespace variation.

---

## 3. When Anchors Are Submitted

Anchors should be submitted at consistent, predictable intervals. Recommended triggers:

| Trigger | Description |
|---|---|
| Per epoch | Submit once per defined epoch (e.g. every 100 L2 blocks) |
| Per application event | Submit after significant state transitions (registry update, projection flush) |
| On demand | Submit when a verifier or consumer requests a fresh commitment |

Minimum requirement: at least one anchor per L2 chain per 24-hour period to maintain
continuous lineage on L1.

The anchor client must always set `parent_anchor_id` to the anchorId returned by the
previous successful anchor submission. Never skip a parent reference — gaps break lineage.

---

## 4. How Replay Proofs Reference Anchors

A replay proof allows any party to reconstruct the L2 state transition from snapshot N to
snapshot N+1. The proof must reference the anchor that committed snapshot N.

Replay proof structure:

```json
{
  "from_anchor_id": "anchor-397529ce6400bfc4",
  "to_anchor_id": "anchor-229a802f0c560b5f",
  "from_block_height": 1042,
  "to_block_height": 1087,
  "state_delta_root": "<hex — root of state delta between snapshots>",
  "event_log": ["<ordered list of state transition events>"],
  "proof_hash": "<SHA-256 of this replay proof document>"
}
```

The `proof_hash` of this document is what goes into the artifact's `replay_proof_hash` field
for the next snapshot. This creates a bidirectional link:

```
anchor N  →  replay_proof_hash in anchor N+1
anchor N+1  →  from_anchor_id in replay proof
```

Any verifier can:
1. Retrieve anchor N and anchor N+1 from L1
2. Retrieve the replay proof from the Bucket Layer
3. Verify `proof_hash` matches `replay_proof_hash` in anchor N+1
4. Replay the state delta to confirm `application_state_root` in anchor N+1

---

## 5. Chain ID Convention

Each L2 chain must use a unique, stable `chain_id`. Recommended format:

```
bhiv-l2-<application-name>-<instance>
```

Examples:
```
bhiv-l2-registry-001
bhiv-l2-projection-001
bhiv-l2-app-001
```

The chain_id is part of the hashed artifact — changing it produces a different hash and
breaks lineage continuity.

---

## 6. Bucket Layer Storage Convention

Full artifacts must be stored in the Bucket Layer at a deterministic path so verifiers can
retrieve them:

```
<chain_id>/<block_height>/artifact.json
```

Example:
```
bhiv-l2-app-001/1042/artifact.json
```

The anchor on L1 stores only the hash. The Bucket Layer stores the full artifact.
Verifiers use the hash to confirm the artifact has not been tampered with after storage.

---

## 7. Integration Checklist for L2 Chain Operators

- [ ] Assign a unique stable `chain_id`
- [ ] Implement state trie to produce `application_state_root` per snapshot
- [ ] Implement registry serialisation to produce `registry_snapshot`
- [ ] Implement projection log Merkle tree to produce `projection_log_root`
- [ ] Implement replay proof bundle to produce `replay_proof_hash`
- [ ] Store full artifact JSON in Bucket Layer at canonical path
- [ ] Run artifact-hash-generator to produce state root hash
- [ ] Run anchor-submit-client to commit hash to L1
- [ ] Record returned anchorId for use as `parent_anchor_id` in next snapshot
- [ ] Verify each anchor using anchor-verifier after submission
