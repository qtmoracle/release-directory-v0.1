# Brick D — Snapshot / Export (Signed View Capture, No Authority)

## Scope
A snapshot is a **signed, immutable export of a rendered view** at a point in time.
It records what was shown — not what is true.

## What is locked in v0.1
- Snapshots are derived from a specific:
  (scope, view, asset_id)
- snapshot_id is deterministic and content-addressed
- snapshot signing uses ed25519 over canonical JSON bytes

## Explicit boundaries
- snapshot_id ≠ asset_id
- snapshots never mint, modify, or supersede assets
- no ranking, interpretation, or verification of truth
- no transport or storage mechanism defined here

## Determinism guarantees
- canonical JSON serialization required
- identical inputs MUST reproduce identical snapshot bytes
- different content MUST produce a different snapshot_id

## Source of truth
The authoritative specification for this brick is:
SNAPSHOT_EXPORT_SPEC_v0.1.md

This README exists to lock semantics and prevent scope creep.
