package constants

// Database location.
const (
	DBDir  = "data"
	DBFile = "gitmap.db"
)

// Lock file.
const (
	LockFileName       = "gitmap.lock"
	LockFilePermission = 0o644
	ErrLockHeld        = "another gitmap process is running (PID %d).\n  If incorrect, delete: %s"
)

// Table names (v15: PascalCase + singular).
const (
	TableRepo      = "Repo"
	TableGroup     = "Group"
	TableGroupRepo = "GroupRepo"
	TableRelease   = "Release"
)

// Legacy table names retained only for migration detection (do not use in new SQL).
const (
	LegacyTableRepos      = "Repos"
	LegacyTableGroups     = "Groups"
	LegacyTableGroupRepos = "GroupRepos"
	LegacyTableReleases   = "Releases"
)

// SQL: create Repo table (v15: singular + RepoId PK).
const SQLCreateRepo = `CREATE TABLE IF NOT EXISTS Repo (
	RepoId           INTEGER PRIMARY KEY AUTOINCREMENT,
	Slug             TEXT NOT NULL,
	RepoName         TEXT NOT NULL,
	HttpsUrl         TEXT NOT NULL,
	SshUrl           TEXT NOT NULL,
	Branch           TEXT NOT NULL,
	RelativePath     TEXT NOT NULL,
	AbsolutePath     TEXT NOT NULL,
	CloneInstruction TEXT NOT NULL,
	Notes            TEXT DEFAULT '',
	CreatedAt        TEXT DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt        TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: create Groups table (legacy plural — renamed in Phase 1.2).
const SQLCreateGroups = `CREATE TABLE IF NOT EXISTS Groups (
	Id          INTEGER PRIMARY KEY AUTOINCREMENT,
	Name        TEXT NOT NULL UNIQUE,
	Description TEXT DEFAULT '',
	Color       TEXT DEFAULT '',
	CreatedAt   TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: create GroupRepo join table (v15: singular). RepoId FK now points to Repo(RepoId).
const SQLCreateGroupRepo = `CREATE TABLE IF NOT EXISTS GroupRepo (
	GroupId INTEGER NOT NULL REFERENCES Groups(Id) ON DELETE CASCADE,
	RepoId  INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE,
	PRIMARY KEY (GroupId, RepoId)
)`

// SQL: create Releases table (legacy plural — renamed in Phase 1.2).
const SQLCreateReleases = `CREATE TABLE IF NOT EXISTS Releases (
	Id           INTEGER PRIMARY KEY AUTOINCREMENT,
	Version      TEXT NOT NULL,
	Tag          TEXT NOT NULL UNIQUE,
	Branch       TEXT NOT NULL,
	SourceBranch TEXT NOT NULL,
	CommitSha    TEXT NOT NULL,
	Changelog    TEXT DEFAULT '',
	Notes        TEXT DEFAULT '',
	Draft        INTEGER DEFAULT 0,
	PreRelease   INTEGER DEFAULT 0,
	IsLatest     INTEGER DEFAULT 0,
	Source       TEXT DEFAULT 'release',
	CreatedAt    TEXT DEFAULT CURRENT_TIMESTAMP
)`

// SQL: add Source column to existing Releases table.
const SQLAddSourceColumn = "ALTER TABLE Releases ADD COLUMN Source TEXT DEFAULT 'release'"

// SQL: enable foreign keys.
const SQLEnableFK = "PRAGMA foreign_keys = ON"

// SQL: repo operations (v15: Repo table, RepoId PK).
const (
	SQLUpsertRepo = `INSERT INTO Repo (Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(AbsolutePath) DO UPDATE SET
			Slug=excluded.Slug, RepoName=excluded.RepoName, HttpsUrl=excluded.HttpsUrl,
			SshUrl=excluded.SshUrl, Branch=excluded.Branch, RelativePath=excluded.RelativePath,
			CloneInstruction=excluded.CloneInstruction, Notes=excluded.Notes, UpdatedAt=CURRENT_TIMESTAMP`

	SQLSelectAllRepos = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes FROM Repo ORDER BY Slug"

	SQLSelectRepoBySlug = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes FROM Repo WHERE Slug = ?"

	SQLSelectRepoByPath = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes FROM Repo WHERE AbsolutePath = ?"
)

// SQL: upsert by AbsolutePath (spec requirement).
const SQLUpsertRepoByPath = `INSERT INTO Repo (Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(AbsolutePath) DO UPDATE SET
		Slug=excluded.Slug, RepoName=excluded.RepoName, HttpsUrl=excluded.HttpsUrl,
		SshUrl=excluded.SshUrl, Branch=excluded.Branch, RelativePath=excluded.RelativePath,
		CloneInstruction=excluded.CloneInstruction, Notes=excluded.Notes, UpdatedAt=CURRENT_TIMESTAMP`

// SQL: create unique index on AbsolutePath for upsert-by-path (v15: IdxRepo_AbsolutePath).
const SQLCreateAbsPathIndex = "CREATE UNIQUE INDEX IF NOT EXISTS IdxRepo_AbsolutePath ON Repo(AbsolutePath)"

// SQL: drop the legacy index name from pre-v15 installs.
const SQLDropLegacyAbsPathIndex = "DROP INDEX IF EXISTS idx_Repos_AbsolutePath"

// SQL: group operations (still on legacy Groups table — renamed in Phase 1.2).
const (
	SQLInsertGroup = "INSERT INTO Groups (Name, Description, Color) VALUES (?, ?, ?)"

	SQLSelectAllGroups = "SELECT Id, Name, Description, Color, CreatedAt FROM Groups ORDER BY Name"

	SQLSelectGroupByName = "SELECT Id, Name, Description, Color, CreatedAt FROM Groups WHERE Name = ?"

	SQLDeleteGroup = "DELETE FROM Groups WHERE Name = ?"

	SQLInsertGroupRepo = "INSERT OR IGNORE INTO GroupRepo (GroupId, RepoId) VALUES (?, ?)"

	SQLDeleteGroupRepo = "DELETE FROM GroupRepo WHERE GroupId = ? AND RepoId = ?"

	SQLSelectGroupRepos = `SELECT r.RepoId, r.Slug, r.RepoName, r.HttpsUrl, r.SshUrl, r.Branch,
		r.RelativePath, r.AbsolutePath, r.CloneInstruction, r.Notes
		FROM Repo r JOIN GroupRepo gr ON r.RepoId = gr.RepoId WHERE gr.GroupId = ? ORDER BY r.Slug`

	SQLCountGroupRepos = "SELECT COUNT(*) FROM GroupRepo WHERE GroupId = ?"
)

// SQL: release operations (still on legacy Releases table — renamed in Phase 1.2).
const (
	SQLUpsertRelease = `INSERT INTO Releases (Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, Draft, PreRelease, IsLatest, Source, CreatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(Tag) DO UPDATE SET
			Version=excluded.Version, Branch=excluded.Branch, SourceBranch=excluded.SourceBranch,
			CommitSha=excluded.CommitSha, Changelog=excluded.Changelog, Notes=excluded.Notes, Draft=excluded.Draft,
			PreRelease=excluded.PreRelease, IsLatest=excluded.IsLatest, Source=excluded.Source`

	SQLSelectAllReleases = `SELECT Id, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, Draft, PreRelease, IsLatest, Source, CreatedAt
		FROM Releases ORDER BY CreatedAt DESC`

	SQLSelectReleaseByTag = `SELECT Id, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, Draft, PreRelease, IsLatest, Source, CreatedAt
		FROM Releases WHERE Tag = ?`

	SQLClearLatestRelease = "UPDATE Releases SET IsLatest = 0 WHERE IsLatest = 1"

	SQLAddNotesColumn = "ALTER TABLE Releases ADD COLUMN Notes TEXT DEFAULT ''"
)

// SQL: reset operations.
const (
	SQLDropGroupRepo  = "DROP TABLE IF EXISTS GroupRepo"
	SQLDropGroupRepos = "DROP TABLE IF EXISTS GroupRepos" // legacy
	SQLDropGroups     = "DROP TABLE IF EXISTS Groups"
	SQLDropRepo       = "DROP TABLE IF EXISTS Repo"
	SQLDropRepos      = "DROP TABLE IF EXISTS Repos" // legacy, kept for migrateLegacyIDs
	SQLDropReleases   = "DROP TABLE IF EXISTS Releases"
)

// Store error messages.
const (
	ErrDBOpen          = "failed to open database at %s: %v (operation: open)"
	ErrDBMigrate       = "failed to initialize tables: %v"
	ErrDBUpsert        = "failed to upsert repo: %v"
	ErrDBQuery         = "failed to query repos: %v"
	ErrDBNoMatch       = "no repo matches slug: %s\n"
	ErrDBCreateDir     = "failed to create database directory at %s: %v (operation: mkdir)"
	ErrDBGroupCreate   = "failed to create group: %v"
	ErrDBGroupQuery    = "failed to query groups: %v"
	ErrDBGroupAdd      = "failed to add repo to group: %v"
	ErrDBGroupRemove   = "failed to remove repo from group: %v"
	ErrDBGroupDelete   = "failed to delete group: %v"
	ErrDBGroupNone     = "no group found: %s"
	ErrDBGroupExists   = "group already exists: %s"
	ErrDBReleaseUpsert = "failed to upsert release: %v"
	ErrDBReleaseQuery  = "failed to query releases: %v"
)

// Phase 1 v15 migration messages.
const (
	MsgV15RepoMigrationStart = "→ Migrating database to v15 schema (Repos → Repo)..."
	MsgV15RepoMigrationDone  = "✓ Migrated Repos → Repo (RepoId PK). Existing data preserved."
	ErrV15RepoMigration      = "v15 Repo migration failed: %v"
	ErrV15RepoCountMismatch  = "v15 Repo migration count mismatch: old=%d new=%d"
)
