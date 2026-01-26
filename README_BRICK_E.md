# Brick E — Webhook Events (Mechanical Signals Only)

## Scope
Webhook events are **mechanical signals** indicating that a state change occurred.
They do not imply correctness, trust, ranking, or interpretation.

## What is locked in v0.1
- Events reference existing identifiers only (e.g., asset_id)
- event_id is content-addressed and unique per occurrence
- Event payloads are minimal and redacted

## Explicit boundaries
- Events never mint or modify asset_id
- No delivery guarantees, retries, or auth defined here
- No consumer semantics or policy enforcement
- No governance, moderation, or adjudication

## Determinism guarantees
- Canonical JSON required before hashing
- event_id = sha256(canonical_event_bytes)
- Same action at different times produces different event_id

## Source of truth
The authoritative specification for this brick is:
WEBHOOK_EVENT_SCHEMA_v0.1.md

This README exists to lock semantics and prevent scope creep.
