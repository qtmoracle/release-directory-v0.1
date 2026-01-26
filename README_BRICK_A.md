# Brick A — Directory Idempotency (MVD fingerprinting + deterministic asset_id)

## Scope
Directory does NOT verify truth. It produces/returns deterministic identifiers based on integrity pointer inputs.

## Algorithm (MVD v0.1)
1) Build fingerprint input as a subset of the mint request:
   - mvd_version (required)
   - provenance (required)
   - labels (if present)
   - p5 (if present)
   - declared (if present)
   - validity (if present)
   - external_ids (if present)
   - actors (if present)

2) Canonicalize JSON:
   - sort map keys recursively
   - preserve array order

3) fingerprint_sha256:
   - sha256(canonical_json_bytes)

4) asset_id:
   - asset_id = "asset:" + sha256("qtm.directory.mvd.0.1|" + fingerprint_sha256_hex)

## Invariants proven
- Same mint request bytes (semantically identical) => same fingerprint_sha256 => same asset_id
- Change one included field => fingerprint_sha256 changes => asset_id changes
- created_at_utc is server-generated and NOT included in fingerprint
