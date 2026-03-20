# Anchor Verification Flow

## Purpose

Describes how any party can independently verify that an on-chain anchor record correctly
represents a given L2 state artifact.

---

## Location

```
anchor-client/anchor-verifier/main.go
```

---

## Build

```bash
cd anchor-client
go build -o anchor-verifier ./anchor-verifier/
```

---

## Usage

```bash
./anchor-verifier <anchor-id|auto> <artifact.json>
```

| Argument | Description |
|---|---|
| anchor-id | The anchorId to verify against. Pass `auto` in simulation mode. |
| artifact.json | The full artifact JSON retrieved from the Bucket Layer. |

---

## Verification Steps

```
1. Retrieve AnchorRecord from L1 contract using anchorId
        |
        v
2. Read artifact JSON from Bucket Layer
        |
        v
3. Parse artifact into canonical Artifact struct
        |
        v
4. Serialise struct → compact JSON → SHA-256 → hex
        |
        v
5. Compare computed hash against AnchorRecord.ArtifactHash
        |
        v
6. Compare artifact.parent_anchor_id against AnchorRecord.ParentHash
        |
        v
7. PASS if both match — FAIL if either differs
```

---

## What Is Verified

| Check | Pass Condition |
|---|---|
| artifactHash integrity | SHA-256(canonical artifact) == on-chain ArtifactHash |
| parentHash linkage | artifact.parent_anchor_id == on-chain ParentHash |

---

## Example Run (verified live)

```bash
./anchor-verifier auto example-artifact.json
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

## Failure Cases

| Scenario | Output |
|---|---|
| Artifact tampered after anchoring | `VERIFICATION FAILED: hash mismatch` |
| Wrong parent reference | `VERIFICATION FAILED: parentHash mismatch` |
| anchorId not found on chain | `VERIFICATION FAILED: anchor retrieval failed` |

---

## Verifier Guarantees

- Stateless: verifier holds no persistent state, reads only from contract and Bucket Layer
- Deterministic: same artifact always produces same hash, so verification is repeatable
- Non-repudiable: once anchored on L1, the record cannot be altered
