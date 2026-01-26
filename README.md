QTM Directory — Minimal Viable Demonstration (v0.3)

This repository contains a minimal, non-authoritative demonstration of the QTM Directory:
a registrar that records the existence of assets and points to integrity signals, without asserting truth, governance, ranking, or meaning.

This work is intended for academic review, grant evaluation, and technical inspection.

What This Demonstrates

This MVP demonstrates that it is possible to coordinate across systems while cleanly separating:

Identity (what exists)

Visibility (who can see it)

Interfaces (how it is viewed)

Exports (what is captured)

Events (what changes)

Specifically, it proves:

Deterministic Identity

Asset identifiers (asset_id) are content-addressed

Same input → same identifier

Changed input → new identifier

Identity is reproducible and idempotent

Visibility Is Orthogonal

Assets may be private, public, or payment_required

Visibility does not alter identity

Commercial surfaces do not corrupt registration

Interfaces Are Not Authority

Planck acts only as a view router

No ranking, scoring, or interpretation occurs

Snapshots Are Immutable Views

Exports are signed, view-specific captures

Snapshots are not truth claims or governance records

Events Are Mechanical Signals

Webhooks indicate state change only

No semantics, no policy, no inference


What This Explicitly Does NOT Do

This system does not:

Verify correctness, truth, or validity of assets

Rank, score, or recommend assets

Perform governance or policy enforcement

Act as an identity provider

Monetize data or users

Assert authority over academic or commercial claims

All interpretation is intentionally left outside the system.

System Boundaries
Component	Role
Academic Surface	Produces signed artifacts (integrity only)
Directory	Registers existence + pointers
Planck	Routes views (no meaning)
Index	Untouched reference layer
Each component is isolated by design to prevent capture or semantic drift.


release_directory_v0.3/
├── academic/                      # Signed academic artifacts (examples)
├── directory_stub/                # Minimal registrar stub
├── fixtures/                      # Determinism test cases
├── README_BRICK_*.md              # Brick-level specifications
├── SNAPSHOT_EXPORT_SPEC_v0.1.md
├── WEBHOOK_EVENT_SCHEMA_v0.1.md
├── PLANCK_VIEW_ROUTER_CONTRACT_v0.1.md
└── README_FREEZE.md                # Local freeze declaration


How to Reproduce (Minimal)

Inspect directory_mvd.json

Modify a field → observe a new asset_id

Revert → observe the original asset_id

Review fixtures to confirm idempotency

Inspect brick READMEs for scope boundaries

No trust assumptions are required.

License & Notice

This repository is provided as a reference implementation.
It makes no claims of fitness for production, governance, or authority.

See LICENSE for terms.

Status

Directory MVP v0.3 — Local Freeze

Publication to GitHub and IPFS pending

Contact / Attribution

Maintained by QTM Benefits Corporation
For research, grants, or inspection only.
