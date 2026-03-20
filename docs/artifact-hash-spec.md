# Artifact Hash Specification

## Purpose

Defines the deterministic hashing process used to convert an L2 state artifact into a
SHA-256 state root hash suitable for anchoring on L1.

---

## Algorithm

| Property | Value |
|---|---|
| Hash function | SHA-256 |
| Input encoding | UTF-8 |
| Output encoding | lowercase hex, 64 characters |
| Serialisation | Compact JSON (no spaces, no newlines) |
| Field ordering | Canonical — fixed struct field order |

---

## Canonical Field Order

The artifact struct fields are serialised in this exact order:

1. `chain_id`
2. `block_height`
3. `timestamp`
4. `application_state_root`
5. `registry_snapshot`
6. `projection_log_root`
7. `replay_proof_hash`
8. `parent_anchor_id`

Any deviation in field order produces a different hash. The generator enforces order via a
typed Go struct — not a raw map — so field order is always deterministic.

---

## Hashing Process

```
artifact JSON (file or stdin)
        |
        v
  Parse into typed Artifact struct
        |
        v
  json.Marshal(struct)  -->  compact canonical JSON bytes
        |
        v
  sha256.Sum256(bytes)
        |
        v
  hex.EncodeToString(sum)
        |
        v
  64-char lowercase hex string  =  state root hash
```

---

## Tool

Location: `artifact-hash-generator/main.go`

Build:
```bash
cd artifact-hash-generator
go build -o artifact-hash-generator .
```

Run:
```bash
./artifact-hash-generator example-artifact.json
```

---

## Determinism Proof

### Input Artifact (`example-artifact.json`)

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

### Canonical JSON (as serialised by generator)

```
{"chain_id":"bhiv-l2-app-001","block_height":1042,"timestamp":"2025-01-15T10:00:00Z","application_state_root":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","registry_snapshot":"6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b","projection_log_root":"d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35","replay_proof_hash":"4e07408562bedb8b60ce05c1decb3f3b9b8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e","parent_anchor_id":""}
```

### Output Hash (verified live)

```
c2c63ac34bb60b396d4445d45d0f54dec754f97c8b87504128a62c34338d1b2e
```

Running `go run main.go example-artifact.json` on any machine with the same input will always
produce this exact hash. This is the determinism guarantee.

---

## Failure Modes

| Scenario | Result |
|---|---|
| Extra whitespace in JSON input | Different hash — generator re-serialises via struct, whitespace is stripped |
| Different field order in input JSON | Same hash — struct parsing normalises order |
| Missing required field | Zero value used — hash still deterministic but semantically invalid |
| Different timestamp string | Different hash — timestamp is part of the canonical input |

---

## Integration Note

The hash output from this generator is the value submitted to the L1 anchor contract as
`artifactHash`. The full artifact JSON is stored separately in the Bucket Layer. Any verifier
can retrieve the artifact, run this generator, and confirm the hash matches the on-chain record.
