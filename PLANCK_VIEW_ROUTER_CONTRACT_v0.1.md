# Planck View Router Contract v0.1 (Text-Only, Non-Authoritative)

## Purpose
Define the minimal mechanical interface for Planck.x1 to route a request for a *view* of a Directory asset.
Planck does not define meaning, rank results, or attest truth. It routes requests to surfaces that already exist.

## Entities (strict boundaries)
- Directory.x1: registrar of existence + pointers (no authority)
- Academic Surface: signed artifact producer (IRP bundles)
- Planck.x1: view router (no authority)
- Index: untouched reference layer (out of scope)

## Inputs
Planck receives a request to render a view of an existing asset:

### Request envelope (conceptual)
- scope: string              # e.g., "directory.x1"
- surface: string            # e.g., "directory"
- view: string               # e.g., "asset_card" | "integrity" | "snapshot" (names are non-semantic)
- asset_id: string           # e.g., "asset:<sha256>"
- params: object (optional)  # view-only parameters; must not change identity

## Routing rule (mechanical)
Given (scope, surface, view, asset_id), Planck chooses a route target:
- target_surface = scope
- target_endpoint = "/view/<view_name>"
Planck may apply access gating in the future, but v0.1 defines routing only.

## Required invariants
1) No authority
   - Planck must never claim verification of truth.
   - Planck may display integrity signals (hash/signature validity) only if verified by an invoked verifier.

2) Identity stability
   - Planck must never mint or modify asset_id.
   - View selection must not affect asset identity.

3) Surface integrity
   - Planck does not rewrite payloads or manifests.
   - Planck passes through pointers as provided by Directory/Surfaces.

4) No Search/Find
   - v0.1 contract does not include query, ranking, indexing, or discovery.

## Minimal view set (names only; semantics prohibited)
- asset_card: render Directory labels + pointers (no claims)
- integrity: render hash/signature fields + verification result if run
- listing: render listing status (private/public/payment_required) if provided by Directory
- snapshot: placeholder name only (export format is out of scope here)

## Non-goals (explicitly out of scope)
- paywalls, billing, or payment processing
- ranking, recommendations, or search
- identity attestation or credentialing
- governance, moderation, disputes, arbitration

