# Memory: features/v15-rename-progress
Updated: now

Phase 1 of the v15 database naming alignment is in flight (DB schema migration only, NO new commands or features yet).

## Spec
PascalCase + **singular** table names + `{TableName}Id` primary keys + FKs match referenced PK name. Per https://github.com/alimtvnetwork/coding-guidelines-v15/blob/main/spec/04-database-conventions/01-naming-conventions.md.

## Phase 1 sub-phase status

| Phase | Scope | Status |
|---|---|---|
| 1.1 | `Repos` → `Repo`, `Id` → `RepoId`, `GroupRepos` → `GroupRepo`, index → `IdxRepo_AbsolutePath`, child-FK clauses updated in 4 sibling constants files (Aliases, DetectedProjects, RepoVersionHistory, GroupRepo), `migrateV15Repo()` added, version bump to v3.1.0 | DONE (sandbox, awaits user `go build` + `go test` verification) |
| 1.2 | `Groups` → `"Group"` (Id → GroupId, double-quoted because reserved word), `Releases` → `Release` (ReleaseId), `Aliases` → `Alias` (AliasId), `Bookmarks` → `Bookmark` (BookmarkId), `migrateV15Phase2()` orchestrator + 4 idempotent migration funcs sharing `execV15Rebuild()` helper, `SQLImportInsertGroup` + `SQLImportInsertBookmark` constants for import.go, GroupRepo FK now references `"Group"(GroupId)`, version bump to v3.2.0 | DONE (sandbox, awaits user `go build` + `go test` verification) |
| 1.3 | `Amendments` → `Amendment`, `CommitTemplates` → `CommitTemplate`, `Settings` → `Setting`, `SSHKeys` → `SshKey`, `InstalledTools` → `InstalledTool`, `TempReleases` → `TempRelease` | TODO |
| 1.4 | ZipGroup family + Project family (incl. `CSharp` → `Csharp` strict v15) + Task family + History tables | TODO |
| 1.5 | Boolean prefix fixes (`Release.Draft` → `IsDraft`, `Release.PreRelease` → `IsPreRelease`) | TODO |
| 1.6 | Update spec/12, regenerate both ERDs, bump CHANGELOG entry, update mem://index core, **explicit GroupRepo + RepoVersionHistory FK rebuild** to clean stale REFERENCES text on existing-user DBs | TODO |

## Phase 1.1 — what changed (recap)
See the previous version of this file. Net result: Repos/Id → Repo/RepoId with atomic table-rebuild migration; child FK clauses updated to `REFERENCES Repo(RepoId)` in Aliases, DetectedProjects, RepoVersionHistory, GroupRepo constants files.

## Phase 1.2 — what changed

### Files edited
- `gitmap/constants/constants_store.go` — added `TableGroup`/`TableRelease`, `SQLCreateGroup` (quoted `"Group"` because SQL reserved word), `SQLCreateRelease`, updated `SQLCreateGroupRepo` FK to `"Group"(GroupId)`, rewrote all Group/Release SQL to v15 names + PKs, added `SQLImportInsertGroup`, added `SQLDropGroup`/`SQLDropRelease`, added 4×4 migration message constants (Group/Release/Alias/Bookmark Start/Done/Migration/CountMismatch). Removed `SQLCreateGroups`, `SQLCreateReleases`.
- `gitmap/constants/constants_alias.go` — full rewrite to v15: `TableAlias`, `LegacyTableAliases`, `SQLCreateAlias` with `AliasId` PK, every SQL renamed `Aliases` → `Alias` and `a.Id` → `a.AliasId`, `SQLDropAlias` added.
- `gitmap/constants/constants_bookmark.go` — full rewrite to v15: `TableBookmark`, `LegacyTableBookmarks`, `SQLCreateBookmark` with `BookmarkId` PK, every SQL renamed `Bookmarks` → `Bookmark`, `SQLDropBookmark` added, `SQLImportInsertBookmark` added.
- `gitmap/store/store.go` — Migrate() chain now calls `migrateV15Phase2()` after `migrateV15Repo()`. CREATE list uses `SQLCreateGroup`/`SQLCreateRelease`/`SQLCreateAlias`/`SQLCreateBookmark` (and removed duplicate). Reset() drops add new singular drops alongside legacy plural drops.
- `gitmap/store/import.go` — 2 raw SQL strings replaced with `constants.SQLImportInsertGroup` + `constants.SQLImportInsertBookmark`.
- `gitmap/store/migrate_v15phase2.go` — NEW. Orchestrator `migrateV15Phase2()` + 4 migration funcs (`migrateV15Group/Release/Alias/Bookmark`) + shared `execV15Rebuild(createSQL, copySQL, legacyTable)` helper. All idempotent.
- `gitmap/constants/constants.go` — Version bumped to `3.2.0`.
- `spec/01-app/gitmap-core-schema-simplified.mmd` — ERD updated with all 4 v15 tables.

### Migration flow on first launch after Phase 1.2 upgrade
1. `migrateLegacyIDs()` (no-op unless very old UUID install).
2. `migrateV15Repo()` (no-op on fresh installs and on already-Phase-1.1-migrated DBs).
3. `migrateV15Phase2()` — runs Group, Release, Alias, Bookmark migrations sequentially. Each is detect-then-act idempotent: skips if legacy table absent, skips if v15 table already exists. Uses `execV15Rebuild()` shared helper which handles PRAGMA fk=OFF → CREATE → INSERT...SELECT (preserving original Id as new {Table}Id) → DROP legacy → fk=ON. Each migration also runs an explicit row-count parity check.
4. Standard CREATE TABLE pass with `IF NOT EXISTS` — safe.

### Known caveats inherited / created
- **Aliases-table FK refresh**: existing user DBs that were on Phase 1.1 had an `Aliases` table with `REFERENCES Repo(RepoId)` but PK column `Id`. Phase 1.2's `migrateV15Alias` rebuilds this to `Alias(AliasId)` with the same FK clause, fully fixing the stale text.
- **GroupRepo FK still references `"Group"(GroupId)` correctly on fresh installs.** For existing-Phase-1.1 users, GroupRepo still has stale text `REFERENCES Groups(Id)`. SQLite checks FKs against the stored schema text, so GroupRepo inserts may fail with `foreign_keys=ON` until GroupRepo is itself rebuilt. **Mitigation scheduled for Phase 1.6**: explicit GroupRepo FK rebuild as part of cleanup.
- **RepoVersionHistory FK** still points at `Repo(RepoId)` correctly (refreshed in Phase 1.1). Will be cleanly rebuilt if/when its own table name changes (no plan to rename — already singular).

## What's still NOT done in Phase 1
- Phases 1.3 / 1.4 / 1.5 / 1.6.
- ScanFolder table (Phase 2).
- VersionProbe table (Phase 2).
- `gitmap find-next` command (Phase 2).
- `gitmap pull` parallel runner (Phase 3).
- `gitmap cn next all` bulk update (Phase 3).
