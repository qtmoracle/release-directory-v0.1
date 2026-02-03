# REPRODUCE.md — QTM Directory MVP v0.4.0 (SRL v0.1)

This document describes how to reproduce and verify the **Directory MVP v0.4.0** surface from source, **without authority** and **without modifying semantics**.

## Scope

- Verifies the **frozen repo snapshot** for `v0.4.0-directory-mvp`
- Verifies the **IPFS root CID** corresponds to that snapshot
- Provides a minimal runtime smoke check for:
  - health
  - mint
  - retrieve
  - set listing
  - restart persistence behavior (shape-level)

## Non-goals

- No enhancements, refactors, or feature additions
- No ranking, inference, governance, or monetization logic
- No mutation beyond the surface's intended append-only runtime logs

---

## 0) Assumptions (Clean VM)

Tested assumptions for a clean Linux VM:

- git installed
- go installed (sufficient to run the stub)
- IPFS (kubo) installed and initialized (`ipfs init`)
- Network access for cloning the repo (or an equivalent transfer method)

> If you cannot install IPFS, you can still verify tag/commit and run the surface,
> but you will not be able to verify the IPFS CID.

---

## 1) Obtain source and checkout the frozen release

Clone the repository (replace the URL with your actual remote if needed):

```bash
git clone <REPO_URL> release_directory_v0.4
cd release_directory_v0.4
