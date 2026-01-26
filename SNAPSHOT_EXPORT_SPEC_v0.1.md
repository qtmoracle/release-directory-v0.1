# Snapshot / Export Specification v0.1 (Text-Only, Signed Snapshot)

## Purpose
Define a minimal, non-authoritative snapshot format for exporting the *state of a view* at a specific time.
A snapshot is an immutable, signed record of what was rendered — not a claim of truth.

## What a snapshot is
- A snapshot is a **read-only export**
- It captures:
  - asset_id
  - view name
  - rendered fields (as provided by Directory/Surfaces)
  - integrity pointers (hashes, signatures, URLs)
  - timestamp
- A snapshot does NOT:
  - change asset identity
  - verify truth
  - rank or interpret data

## Snapshot envelope (conceptual)
snapshot_version: "0.1"
snapshot_id: "snapshot:<sha256>"
created_at_utc: "<RFC3339>"
asset_id: "asset:<sha256>"
scope: "<surface scope>"
view: "<view name>"
payload:
  rendered_fields: { ... }       # view output only
  integrity_pointers: { ... }    # hashes, signatures, manifests
signing:
  algorithm: "ed25519"
  public_key_base64: "<base64>"
  signature_base64: "<base64>"
determinism:
  canonicalization: "encoding/json canonical"
  hash_algorithm: "sha256"

## Determinism rules
- Snapshot payload must be canonicalized before hashing/signing
- snapshot_id = sha256(canonical_snapshot_bytes)
- snapshot_id MUST NOT be reused for different content

## Retrieval (v0.1 constraints)
- Snapshots may be:
  - returned inline
  - written to file
  - pinned to IPFS
- Retrieval mechanism is out of scope for v0.1
- v0.1 defines the **format only**, not transport

## Required invariants
1) Identity isolation
   - snapshot_id ≠ asset_id
   - snapshots never mint or modify assets

2) Non-authority
   - snapshots record what was rendered, not what is true

3) Reproducibility
   - given identical inputs + view, the same snapshot bytes MUST be reproducible

## Non-goals
- delta snapshots
- streaming updates
- search, indexing, ranking
- legal or governance claims

