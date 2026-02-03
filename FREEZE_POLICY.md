# Freeze Policy — release_directory_v0.*

## Rule
Frozen releases are identified by **git tags** (e.g., `v0.4.0-directory-mvp`) and the corresponding **IPFS root CID** recorded in `IPFS_FREEZE.txt`.

## Allowed changes after a freeze tag
The `main` branch may advance **only** to add:
- reproduction documentation (e.g., `REPRODUCE.md`)
- audits / checklists (e.g., `SRL_AUDIT_*`)
- navigation maps / meta docs that do not alter system behavior

## Prohibited changes after a freeze tag (unless a new version is declared)
- modifying Brick A/B/F semantics
- modifying runtime behavior in `directory_stub/`
- modifying UI behavior beyond documentation-only changes
- refactoring that changes deterministic outputs

## How to cut a new release
If any behavior changes are needed, a **new version** must be declared and frozen under a new tag and CID.
