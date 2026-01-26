# Webhook Event Schema v0.1 (Mechanical Signals Only)

## Purpose
Define a minimal, non-authoritative event schema for notifying external systems that
a **mechanical state change occurred** in the Directory or related surfaces.

Events convey *that something happened*, not *what it means*.

## Event principles
- Events are signals, not assertions
- Events do not verify truth
- Events do not carry interpretation, ranking, or policy
- Events reference existing identifiers only

## Event envelope (conceptual)
event_version: "0.1"
event_id: "event:<sha256>"
event_type: "<string>"
occurred_at_utc: "<RFC3339>"
source:
  surface: "<directory | planck | surface-name>"
  scope: "<scope identifier>"
subject:
  asset_id: "asset:<sha256>"
data:
  before: { ... }   # optional, redacted/minimal
  after: { ... }    # optional, redacted/minimal
integrity:
  hash_algorithm: "sha256"
  event_hash: "<hex>"
determinism:
  canonicalization: "encoding/json canonical"

## Event ID rules
- event_id = sha256(canonical_event_bytes)
- event_id MUST be unique per event occurrence
- Same semantic action at different times => different event_id

## Allowed event types (v0.1)
- asset_minted
- listing_changed
- snapshot_created
- integrity_verified

## Required invariants
1) Identity isolation
   - Events never mint or modify asset_id
   - Events reference existing identifiers only

2) Minimal disclosure
   - data.before / data.after MUST be minimal
   - No payload duplication unless required for integrity

3) Non-authority
   - Events do not imply correctness, trust, or endorsement

## Delivery (out of scope)
- HTTP, queues, retries, signing, auth are NOT defined here
- v0.1 defines schema only

## Non-goals
- subscriptions, filters, or routing logic
- guarantees of delivery
- consumer semantics
- governance or policy enforcement

