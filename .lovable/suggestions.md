# Suggestions Tracker

Pending improvement ideas — not yet approved for implementation.

---

## Active Suggestions

### Add `version-history` to docs site sidebar/commands navigation
- **Status:** Pending
- **Priority:** Low
- **Description:** Page exists at `/version-history` but is not linked from the sidebar or commands page; users won't discover it organically.
- **Added:** v2.76.0 session

### Add `clone` page to docs site
- **Status:** Pending
- **Priority:** Low
- **Description:** Cover both file-based and direct-URL clone documentation.
- **Added:** v2.76.0 session

### Add `--dry-run` flag to `clone-next`
- **Status:** Pending
- **Priority:** Medium
- **Description:** Spec 87-clone-next-flatten.md mentions previewing destructive folder removal; not yet implemented.
- **Added:** v2.75.0 session

### Expand `install` command with database tools
- **Status:** Pending
- **Priority:** Medium
- **Description:** MySQL, MariaDB, PostgreSQL, SQLite, MongoDB, CouchDB, Redis, Cassandra, Neo4j, Elasticsearch, DuckDB, Chocolatey, Winget. Full table in `.lovable/plan.md` Part B.
- **Added:** plan.md Part B

### Add `gitmap uninstall <tool>` command
- **Status:** Pending
- **Priority:** Medium
- **Description:** Per-package-manager uninstall builders + DB record cleanup. Spec in `.lovable/plan.md` Part E.
- **Added:** plan.md Part E

### Enhanced `install --list` grouped by category with installed status
- **Status:** Pending
- **Priority:** Low
- **Description:** Group by Core/Databases/Utilities; show ✓/✗ + version from new `InstalledTools` DB table.
- **Added:** plan.md Part F

### Unit tests for task, env, and install commands
- **Status:** Pending
- **Priority:** Low
- **Description:** Coverage gap open since v2.49.0.
- **Added:** v2.49.0

### Update `helptext/env.md` with `--shell` flag examples
- **Status:** Pending
- **Priority:** Low
- **Description:** `--shell` is wired but not demonstrated in `gitmap help env`.
- **Added:** v2.49.0

### Create `spec/01-spec-authoring-guide/` with spec writing conventions
- **Status:** Pending
- **Priority:** Low
- **Description:** Document the spec authoring conventions used across `spec/`.
- **Added:** v3.3.0 session

### Add the `version` badge to `Index.tsx` landing page hero
- **Status:** Pending
- **Priority:** Low
- **Description:** Visitors should see the current version on the homepage, not only inside the docs layout header.
- **Added:** v3.0.0 session (this session)

### Centralize `VERSION` constant in `src/constants/index.ts`
- **Status:** Pending
- **Priority:** Medium
- **Description:** Currently hardcoded in `DocsLayout.tsx`. Move to one place so the docs layout, landing page, and any future footer all import the same value.
- **Added:** v3.0.0 session (this session)

### Lint rule for missing `// gitmap:cmd top-level` markers
- **Status:** Pending
- **Priority:** Medium
- **Description:** Scan `constants/*.go` for files containing `Cmd[A-Z]` string constants without the marker; warn contributors at PR time.
- **Added:** v3.0.0 session (this session)

### Integration test for `release-alias` auto-stash round-trip
- **Status:** Pending
- **Priority:** Medium
- **Description:** Create a temp Git repo, register via `runAs`, dirty the tree, assert `autoStashIfDirty` + `popAutoStash` round-trip leaves the working tree byte-identical.
- **Added:** v3.0.0 session (this session)

---

## Implemented Suggestions

### `--flatten` for `clone-next` → default behavior
- **Implemented:** v2.75.0
- **Notes:** Clones into base-name folder by default; tracked in `RepoVersionHistory`.

### `gitmap clone <url>` auto-flatten versioned URLs
- **Implemented:** v2.75.0

### `RepoVersionHistory` table for tracking version transitions
- **Implemented:** v2.75.0

### `gitmap version-history` (`vh`) command
- **Implemented:** v2.76.0

### Database ERD covering all 22 tables
- **Implemented:** v2.76.0

### Tab completion for `version-history`/`vh`
- **Implemented:** v2.76.0

### Docs site page for version-history
- **Implemented:** v2.76.0

### `gitmap doctor setup` checks
- **Implemented:** v2.74.0

### Shell wrapper `GITMAP_WRAPPER=1` detection
- **Implemented:** v2.74.0

### VS Code admin-mode bypass with 3-tier launch strategy
- **Implemented:** v2.72.0

### `spec/12-consolidated-guidelines/` with 18 unified guideline documents
- **Implemented:** v3.3.0

### `gitmap as` / `release-alias` (`ra`) / `release-alias-pull` (`rap`)
- **Implemented:** v3.0.0
- **Notes:** Auto-stash labeled `gitmap-release-alias autostash <alias>-<version>-<unix-ts>`, popped via label-match against `git stash list` for concurrent safety. Files: `cmd/{as,asops,releasealias,releasealias_git}.go`. Spec: `spec/01-app/98-as-and-release-alias.md`.

### `gitmap db-migrate` (`dbm`)
- **Implemented:** v3.0.0
- **Notes:** Idempotent re-run of every `CREATE TABLE IF NOT EXISTS` and column migration; auto-invoked at end of `gitmap update`.

### Marker-comment opt-in for completion generator
- **Implemented:** v3.0.0
- **Notes:** Replaces `sourceFiles`/`skipNames` with `// gitmap:cmd top-level` and `// gitmap:cmd skip`. See `mem://features/marker-comments`.

### CI `generate-check` drift detection
- **Implemented:** v3.0.0
- **Notes:** `.github/workflows/ci.yml` runs `go generate ./...` + `git diff --exit-code`; wired into `test-summary` needs.

### `migrateTRCommitSha` switched to detect-then-act
- **Implemented:** v3.0.0
- **Notes:** Uses `PRAGMA table_info(TempReleases)` instead of brittle string-matching on `"no such column"`; kills the cosmetic warning on Unix builds.

### Migration guide section in CHANGELOG
- **Implemented:** v3.0.0 session (this session)
- **Notes:** Added `## Migration guide — v2.x → v3.0.0 (constants contributors)` block at the top of `CHANGELOG.md` with marker-comment example and verification steps.

### Spec doc `spec/01-app/98-as-and-release-alias.md`
- **Implemented:** v3.0.0 session (this session)
- **Notes:** Matches the 97-move-and-merge.md format; covers dispatcher wiring, auto-stash semantics, exit codes.

### v3.0.0 badge in docs site header
- **Implemented:** v3.0.0 session (this session)
- **Notes:** `src/components/docs/DocsLayout.tsx` — Tailwind classes `ml-2 px-2 py-0.5 text-xs font-mono bg-primary/10 text-primary rounded`.
