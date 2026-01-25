# Directory v0.1 + Academic MVP — Freeze

**Freeze date:** 2026-01-25

## What this bundle demonstrates

This release freezes two completed MVPs:

1. **Academic Surface MVP (v0.1)**
   - Deterministic IRP packaging (`qtm irp pack`)
   - Canonical JSON payload
   - SHA-256 hashing
   - Ed25519 signing
   - Strict verification (`qtm irp verify --strict`)
   - Explicit non-assertive claims boundary

2. **Directory MVP (v0.1)**
   - Minimum Viable Data (MVD) schema
   - Base Asset mint endpoint
   - External IDs referencing signed academic IRP artifacts
   - Provenance nonce for replay protection
   - No semantic interpretation or enforcement

## What is explicitly NOT claimed

- No claims of truth, accuracy, or quality of the academic paper
- No ranking, scoring, moderation, or endorsement
- No persistence, idempotency, or storage guarantees
- No search, discovery, or monetization logic
- No identity verification or attestation
- No authority delegated to the Directory

## Architectural intent

- Academic surfaces produce signed artifacts
- Directory registers existence and pointers only
- Index remains untouched and non-operational here
- Planck routing, search, and monetization are deferred

## Status

As of this freeze, **Academic MVP v0.1** and **Directory MVP v0.1** are complete.
All future work must build forward from these artifacts without revision.
