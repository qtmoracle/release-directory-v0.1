# Brick C — Planck View Router (Routing Only, No Authority)

## Scope
Planck.x1 is a **view router** only.
It does not mint assets, verify truth, rank results, or define meaning.

## What is locked in v0.1
- Planck routes requests based on:
  (scope, surface, view, asset_id)
- Planck never modifies or recomputes asset_id
- Planck passes through pointers exactly as provided by Directory/Surfaces

## Explicit boundaries
- No Search/Find
- No paywalls or billing
- No identity attestation
- No governance or moderation
- No payload rewriting

## Source of truth
The authoritative contract for this brick is:
PLANCK_VIEW_ROUTER_CONTRACT_v0.1.md

This README exists to lock intent and prevent scope creep.
