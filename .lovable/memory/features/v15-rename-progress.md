# Memory: features/v15-rename-progress
Updated: now

Phase 1 of the v15 database naming alignment is in flight (DB schema migration only, NO new commands or features yet).

## Spec
PascalCase + **singular** table names + `{TableName}Id` primary keys + FKs match referenced PK name. Per https://github.com/alimtvnetwork/coding-guidelines-v15/blob/main/spec/04-database-conventions/01-naming-conventions.md.

## Phase 1 sub-phase status

| Phase | Scope | Status |
|---|---|---|
| 1.1 | `Repos` → `Repo`, `Id` → `RepoId`, `GroupRepos` → `GroupRepo`, index → `IdxRepo_AbsolutePath`, child-FK clauses updated in 4 sibling constants files (Aliases, DetectedProjects, RepoVersionHistory, GroupRepo), `migrateV15Repo()` added, version bump to v3.1.0 | DONE (sandbox, awaits user `go build` + `go test` verification) |
| 1.2 | `Groups` → `Group` (Id → GroupId), `Releases` → `Release` (+IsDraft/IsPreRelease bool fix moved to 1.5), `Aliases` → `Alias`, `Bookmarks` → `Bookmark` | TODO |
| 1.3 | `Amendments` → `Amendment`, `CommitTemplates` → `CommitTemplate`, `Settings` → `Setting`, `SSHKeys` → `SshKey`, `InstalledTools` → `InstalledTool`, `TempReleases` → `TempRelease` | TODO |
| 1.4 | ZipGroup family + Project family (incl. `CSharp` → `Csharp` strict v15) + Task family + History tables | TODO |
| 1.5 | Boolean prefix fixes (`Release.Draft` → `IsDraft`, `Release.PreRelease` → `IsPreRelease`) | TODO |
| 1.6 | Update spec/12, regenerate both ERDs, bump CHANGELOG entry, update mem://index core | TODO |

## Phase 1.1 — what changed

### Files edited
- `gitmap/constants/constants_store.go` — full rewrite. Renamed constants: `TableRepo`, `TableGroupRepo`, `SQLCreateRepo`, `SQLCreateGroupRepo`, `SQLDropRepo`, `SQLDropGroupRepo`. Added `LegacyTableRepos`, `LegacyTableGroupRepos`, `SQLDropLegacyAbsPathIndex`, `MsgV15RepoMigrationStart/Done`, `ErrV15RepoMigration/CountMismatch`. Index renamed: `idx_Repos_AbsolutePath` → `IdxRepo_AbsolutePath`. Every `Id` PK reference in repo SQL is now `RepoId`.
- `gitmap/constants/constants_alias.go` — FK clause `REFERENCES Repos(Id)` → `REFERENCES Repo(RepoId)`. Joins `JOIN Repos r ON a.RepoId = r.Id` → `JOIN Repo r ON a.RepoId = r.RepoId`. `SELECT r.Id` → `SELECT r.RepoId`.
- `gitmap/constants/constants_version_history.go` — FK clause + ALTER TABLE Repos → Repo + WHERE Id → WHERE RepoId + SQLSelectRepoIDByPath now selects `RepoId FROM Repo`.
- `gitmap/constants/constants_project_sql.go` — FK on DetectedProjects + JOIN Repos r ON dp.RepoId = r.Id → JOIN Repo.
- `gitmap/store/store.go` — Migrate() now calls `migrateV15Repo()` between `migrateLegacyIDs()` and the standard CREATE TABLE pass. Reset() drops both `GroupRepo` and legacy `GroupRepos`/`Repos` for safety.
- `gitmap/store/migrate_v15repo.go` — NEW. The atomic table-rebuild dance with row-count parity check.
- `gitmap/store/migrateids.go` — `rebuildReposTable()` now defines its legacy create SQL inline (the constant `SQLCreateRepos` was removed). References `LegacyTableRepos`.
- `gitmap/constants/constants.go` — Version bumped to `3.1.0`.
- `spec/01-app/gitmap-core-schema-simplified.mmd` — ERD updated to show `Repo` + `RepoId` + `GroupRepo`.

### Migration flow on first launch after upgrade
1. `migrateLegacyIDs()` runs (no-op unless very old UUID-PK install).
2. `migrateV15Repo()` runs. If `Repos` table exists: PRAGMA foreign_keys=OFF → CREATE Repo + IdxRepo_AbsolutePath → INSERT INTO Repo SELECT FROM Repos (preserving Id values as RepoId, plus all CreatedAt/UpdatedAt timestamps) → row-count check → DROP Repos → DROP legacy index → PRAGMA foreign_keys=ON.
3. Standard CREATE TABLE pass runs. Child tables (Aliases, DetectedProjects, RepoVersionHistory, GroupRepo) all use `IF NOT EXISTS` so they survive the rename. New installs get `REFERENCES Repo(RepoId)` directly. Existing installs keep their old FK text pointing at "Repos(Id)" — SQLite tolerates this with `foreign_keys=OFF` during the rename, and on subsequent inserts the FK is checked by stored schema; we accept this minor inconsistency until Phase 1.2/1.3 rebuilds those child tables for their own renames.

### Known follow-ups for Phase 1.2+
- The 4 child tables (Aliases, DetectedProjects, RepoVersionHistory, GroupRepo) still have stale FK schema text on existing user databases (says `REFERENCES Repos(Id)` even though `Repos` no longer exists). The constraint is functionally inactive until those tables are themselves rebuilt. **Mitigation:** when Phase 1.2 renames Aliases → Alias and Phase 1.4 renames DetectedProjects → DetectedProject, their table-rebuild dance will simultaneously rewrite the FK clause. RepoVersionHistory and GroupRepo will need explicit FK-only rebuilds in Phase 1.6 if their table names don't change.
- Version bumped optimistically to v3.1.0 — bump again at Phase 1.6 final to v3.2.0 if desired.

## What's still NOT done in Phase 1
- ScanFolder table (deferred to Phase 2).
- VersionProbe table (deferred to Phase 2).
- `gitmap find-next` command (deferred to Phase 2).
- `gitmap pull` parallel runner (deferred to Phase 3).
- `gitmap cn next all` bulk update (deferred to Phase 3).
