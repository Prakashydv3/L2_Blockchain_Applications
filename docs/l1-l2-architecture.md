# L1/L2 Architecture — BHIV Blockchain

## Overview

BHIV operates as a two-layer blockchain system. L1 is the immutable truth spine. L2 chains are
application execution environments that periodically commit their state into L1.

---

## BHIV L1 Chain

Role: Sovereign Truth Spine

- Stores anchor commitments produced by L2 chains
- Does not execute application logic
- Does not validate L2 transactions
- Provides a permanent, tamper-proof record of L2 state roots
- Exposes an anchor contract that L2 clients call to register commitments

The L1 chain is not modified by this integration. Only the anchor contract deployed on L1 is
interacted with.

---

## L2 Application Chains

Role: Execution Layer

- Execute BHIV application logic (registries, projections, state transitions)
- Produce a state root at defined intervals (per block, per epoch, or per application event)
- State root summarises the full application state at a given point in time
- L2 chains do not write directly to L1 — they produce artifacts that the anchor client submits

---

## Anchor Contract

Location: Deployed on BHIV L1

Responsibilities:
- Accept anchor submissions from the off-chain anchor client
- Store: artifactHash, parentHash, submitter, timestamp, anchorId
- Emit an AnchorCommitted event on each successful submission
- Allow read queries to retrieve any anchor by anchorId

The anchor contract is the only L1 component touched by this integration.

---

## Off-Chain Anchor Client

Role: Bridge Layer

- Reads L2 state artifacts
- Hashes them deterministically using SHA-256
- Submits the hash to the L1 anchor contract
- Returns the anchorId for downstream verification

---

## Bucket Layer

Role: Artifact Storage

- Stores the full artifact JSON referenced by each anchor hash
- Allows any verifier to retrieve the original artifact and recompute the hash
- Not part of the L1 chain — operates as off-chain storage (S3, IPFS, or local store)

---

## State Commitment Flow

```
L2 Application State
        |
        v
  State Snapshot (artifact JSON)
        |
        v
  SHA-256 Hash  <-- deterministic, canonical field ordering
        |
        v
  Anchor Client
        |
        v
  L1 Anchor Contract  -->  AnchorCommitted Event
        |
        v
  anchorId stored on L1
        |
        v
  Off-Chain Verifier
    - retrieves artifact from Bucket Layer
    - recomputes hash
    - compares against on-chain record
    - verifies parentHash linkage
```

---

## Key Constraints

| Constraint | Detail |
|---|---|
| L1 consensus | Not modified |
| L1 block validation | Not modified |
| L1 transaction execution | Not modified |
| Governance logic | Not introduced |
| Anchor contract | Read/write via anchor client only |
| Hash algorithm | SHA-256, deterministic field ordering |

---

## Component Responsibilities Summary

| Component | Owned By | Writes To |
|---|---|---|
| L1 Chain | BHIV Core | Itself (consensus) |
| Anchor Contract | This integration | L1 state (via tx) |
| L2 Chain | BHIV App Layer | Its own state |
| Anchor Client | This integration | L1 via contract call |
| Bucket Layer | Infrastructure | Off-chain storage |
| Verifier | This integration | Read-only |
