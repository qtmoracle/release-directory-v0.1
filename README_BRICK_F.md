# Brick F — Directory Retrieval Endpoint (Fixture-Backed, Read-Only)

## Scope
Adds a minimal read-only retrieval endpoint for demo/view rendering without persistence.

## Endpoint
GET /asset?asset_id=...

## Behavior (v0.1)
- Looks up asset_id in a hard-coded fixture map (no DB)
- Loads the corresponding Directory MVD JSON from disk
- Returns view payload:
  - view: "asset_card"
  - asset_id
  - directory_mvd (as loaded)
  - listing (deterministic default: private, unpaid)
  - rendered_at_utc

## Explicit boundaries
- No persistence (no DB)
- No search / ranking
- No authority / truth verification
- No listing state storage (listing is default-only in this endpoint)

## Purpose
Enables minimal Directory UI + Planck UI routing in the next phase without contaminating identity.
