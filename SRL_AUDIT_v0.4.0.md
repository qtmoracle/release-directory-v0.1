# SRL Compliance Audit — Directory MVP v0.4.0
Status: AUDIT ONLY (NO CHANGES)
Scope: Apply SRL v0.1 invariants to the existing Directory MVP v0.4.0 repository.
Repo: release_directory_v0.4
Branch: main
Head (post-REPRODUCE): ff63fe1
Tag (frozen milestone): v0.4.0-directory-mvp
IPFS root CID (repo snapshot): QmcAornreC2GsAirX2PATBiBoSaK8AetQpNtDuLL2cRNyg

---

## Audit Legend
- **Satisfied** — evidence exists in repo and/or demonstrated outputs
- **Ambiguous** — partially supported, but not explicitly proven
- **Missing** — invariant not met or no evidence exists

Evidence format: `path` + `command/output reference`

---

## SRL v0.1 — Core Invariants

### 1) Deterministic
**Result:** Satisfied  
**Evidence:**
- Deterministic identity + canonicalization described in brick docs:
  - `README_BRICK_A.md` (canonical subset selection, canonical JSON, sha256 fingerprint, asset_id derivation)
- Deterministic retrieval behavior described:
  - `README_BRICK_F.md` (read-only retrieval endpoint)
- Deterministic repo snapshot produced:
  - IPFS root CID exists and lists repo contents:
    - `ipfs ls QmcAornreC2GsAirX2PATBiBoSaK8AetQpNtDuLL2cRNyg`

**Notes:** Determinism at the algorithm/contract level is documented. Runtime determinism depends on reproduction verification.

---

### 2) Neutral
**Result:** Satisfied  
**Evidence:**
- Explicit non-goals and neutrality posture in repo docs:
  - `README.md` / `README_FREEZE.md` (no authority, no ranking, no semantics)
- Bricks C/D/E defined as mechanical/spec-only and not authoritative:
  - `README_BRICK_C.md`, `README_BRICK_D.md`, `README_BRICK_E.md`

---

### 3) Non-authoritative
**Result:** Satisfied  
**Evidence:**
- No ranking/scoring/governance claims in repo docs:
  - `README.md` / brick docs
- Directory acts as integrity producer + read-only retrieval:
  - `README_BRICK_F.md`
- Persistence model is append-only logs with restart-safe rebuild:
  - directory implementation + docs (see `directory_stub/` + README)

---

### 4) Reproducible by third parties
**Result:** Ambiguous  
**Evidence:**
- `REPRODUCE.md` now exists (added post-freeze milestone) and describes how to reproduce.
- IPFS CID exists for repo snapshot.

**Ambiguity drivers:**
- Third-party reproduction has not been demonstrated *in this repo* as an executed proof artifact (e.g., recorded outputs, hashes, or CI verification).
- Local environment assumptions may be sufficient, but audit does not yet include an independent reproduction transcript.

---

### 5) Frozen core with explicit versioning
**Result:** Satisfied  
**Evidence:**
- Frozen milestone tag exists:
  - `v0.4.0-directory-mvp`
- Freeze intent + boundaries in:
  - `README_FREEZE.md`
- IPFS freeze recorded (canonical repo snapshot excluding runtime data):
  - `IPFS_FREEZE.txt` (contains git commit, tag, CID, add flags, timestamp)
- Runtime data excluded (by policy):
  - `.gitignore` and freeze notes

**Note:** Tag is a milestone. Branch may advance with audit/support docs (e.g., REPRODUCE.md) without mutating frozen semantics.

---

## Normative Steps (SRL v0.1)

### Step 1) Surface Intent Declaration
**Result:** Satisfied  
**Evidence:**
- `README.md` and `README_FREEZE.md` describe intent: deterministic MVP directory, proof UI only, non-platform.

---

### Step 2) Input Contract Definition
**Result:** Ambiguous  
**Evidence:**
- `directory_mvd.json` exists (MVD schema / input data model)
- Brick A describes canonical subset selection and hashing

**Ambiguity drivers:**
- Input contract is present as files + brick docs, but a single explicit `INPUT_CONTRACT.md` artifact (SRL template style) does not exist in this repo.

---

### Step 3) Canonicalization
**Result:** Satisfied  
**Evidence:**
- Brick A defines canonicalization and hash derivations:
  - `README_BRICK_A.md`
- CLI-style determinism implied through canonical JSON handling

---

### Step 4) Freeze & Commit
**Result:** Satisfied  
**Evidence:**
- Tag exists: `v0.4.0-directory-mvp`
- IPFS CID exists + freeze receipt (repo snapshot excluding runtime data):
  - `IPFS_FREEZE.txt`
  - `ipfs ls <CID>` shows expected files

---

### Step 5) Minimal Read-Only UI
**Result:** Satisfied  
**Evidence:**
- UI exists in `ui/` and is static HTML/JS (no build step).
- Served at `/ui/` and confirmed by curl:
  - `curl -sS http://127.0.0.1:8080/ui/ | head` shows UI HTML.

---

### Step 6) Reproduction Proof (`REPRODUCE.md`)
**Result:** Satisfied (artifact present), Reproduction success = Ambiguous  
**Evidence:**
- `REPRODUCE.md` exists on main (post-push).

**Ambiguity drivers:**
- Audit does not include a recorded third-party run proving reproduction in an independent environment.

---

## Summary Matrix

| SRL Item | Result |
|---|---|
| Deterministic | Satisfied |
| Neutral | Satisfied |
| Non-authoritative | Satisfied |
| Reproducible by third parties | Ambiguous |
| Frozen core w/ explicit versioning | Satisfied |
| Surface Intent Declaration | Satisfied |
| Input Contract Definition | Ambiguous |
| Canonicalization | Satisfied |
| Freeze & Commit | Satisfied |
| Minimal Read-Only UI | Satisfied |
| REPRODUCE.md present | Satisfied |
| Independent reproduction proof | Ambiguous |

---

## Audit Stop Condition
This audit makes no changes and proposes no fixes.
Next actions (if authorized) are limited to SRL Step 2+ sequencing:
- strengthen reproduction proof evidence
- freeze SRL v0.1 only after successful validation
