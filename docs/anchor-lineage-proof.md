# Anchor Lineage Proof

## Purpose

Demonstrates that sequential L2 state snapshots form a verifiable chain on L1, where each
anchor references the anchorId of the previous anchor via `parentHash`.

---

## Simulation Tool

```
anchor-client/lineage-sim/main.go
```

Run:
```bash
cd anchor-client
go run lineage-sim/main.go
```

---

## Lineage Chain — Verified Live Output

```
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

---

## Lineage Table

| Snapshot | Block | AnchorId | ParentAnchorId |
|---|---|---|---|
| snapshot1 (genesis) | 1042 | `anchor-397529ce6400bfc4` | _(none)_ |
| snapshot2 | 1087 | `anchor-229a802f0c560b5f` | `anchor-397529ce6400bfc4` |
| snapshot3 | 1134 | `anchor-6262a56e88e53761` | `anchor-229a802f0c560b5f` |

---

## Lineage Diagram

```
[genesis]
anchor-397529ce6400bfc4  (block 1042)
        |
        v
anchor-229a802f0c560b5f  (block 1087)
        |
        v
anchor-6262a56e88e53761  (block 1134)
```

---

## Continuity Guarantee

Each anchor's `parentHash` field on L1 points to the previous anchorId. This creates a
tamper-evident chain:

- To forge snapshot2, an attacker must also forge snapshot1 (already on L1 — immutable)
- To forge snapshot3, an attacker must forge snapshot2 and snapshot1
- The chain is only as weak as the L1 itself — which is the sovereign truth spine

---

## Verification Logic

The lineage verifier walks the chain backwards from the latest anchor:

```
anchor-6262a56e88e53761
  → parentHash: anchor-229a802f0c560b5f  ✓ exists on L1
      → parentHash: anchor-397529ce6400bfc4  ✓ exists on L1
          → parentHash: ""  ✓ genesis reached
```

All links resolve. Lineage is intact.
