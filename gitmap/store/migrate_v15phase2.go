// Package store — migrate_v15phase2.go performs Phase 1.2 of the v15 rename:
// Groups -> "Group", Releases -> Release, Aliases -> Alias, Bookmarks -> Bookmark.
//
// All four migrations follow the same atomic table-rebuild pattern:
//
//  1. Detect-then-act: skipped on fresh installs and on already-migrated DBs.
//  2. PRAGMA foreign_keys=OFF for the duration of each rebuild.
//  3. CREATE new singular table -> INSERT...SELECT preserving Id values as
//     {Table}Id (column rename via column list) -> row-count parity check ->
//     DROP legacy plural table -> PRAGMA foreign_keys=ON.
//
// Note: "Group" is a SQL reserved word — quoted with double quotes everywhere.
package store

import (
	"fmt"
	"os"

	"github.com/user/gitmap/constants"
)

// migrateV15Phase2 runs all four Phase 1.2 table renames in sequence.
// Each is independently idempotent.
func (db *DB) migrateV15Phase2() error {
	if err := db.migrateV15Group(); err != nil {
		return fmt.Errorf(constants.ErrV15GroupMigration, err)
	}

	if err := db.migrateV15Release(); err != nil {
		return fmt.Errorf(constants.ErrV15ReleaseMigration, err)
	}

	if err := db.migrateV15Alias(); err != nil {
		return fmt.Errorf(constants.ErrV15AliasMigration, err)
	}

	if err := db.migrateV15Bookmark(); err != nil {
		return fmt.Errorf(constants.ErrV15BookmarkMigration, err)
	}

	return nil
}

// migrateV15Group rebuilds Groups -> "Group" with GroupId PK.
func (db *DB) migrateV15Group() error {
	if !db.tableExists(constants.LegacyTableGroups) {
		return nil
	}

	if db.tableExists(constants.TableGroup) {
		_, _ = db.conn.Exec("DROP TABLE IF EXISTS Groups")

		return nil
	}

	fmt.Println(constants.MsgV15GroupMigrationStart)

	oldCount, err := db.countRows(constants.LegacyTableGroups)
	if err != nil {
		return fmt.Errorf("count Groups: %w", err)
	}

	const copySQL = `INSERT INTO "Group"
		(GroupId, Name, Description, Color, CreatedAt)
		SELECT Id, Name, Description, Color, CreatedAt
		FROM Groups`

	if err := db.execV15Rebuild(constants.SQLCreateGroup, copySQL, "Groups"); err != nil {
		return err
	}

	newCount, err := db.countRows(constants.TableGroup)
	if err != nil {
		return fmt.Errorf(`count "Group": %w`, err)
	}

	if oldCount != newCount {
		fmt.Fprintf(os.Stderr, "  ✗ v15 Group row-count mismatch: old=%d new=%d\n", oldCount, newCount)

		return fmt.Errorf(constants.ErrV15GroupCountMismatch, oldCount, newCount)
	}

	fmt.Println(constants.MsgV15GroupMigrationDone)

	return nil
}

// migrateV15Release rebuilds Releases -> Release with ReleaseId PK.
func (db *DB) migrateV15Release() error {
	if !db.tableExists(constants.LegacyTableReleases) {
		return nil
	}

	if db.tableExists(constants.TableRelease) {
		_, _ = db.conn.Exec("DROP TABLE IF EXISTS Releases")

		return nil
	}

	fmt.Println(constants.MsgV15ReleaseMigrationStart)

	oldCount, err := db.countRows(constants.LegacyTableReleases)
	if err != nil {
		return fmt.Errorf("count Releases: %w", err)
	}

	const copySQL = `INSERT INTO Release
		(ReleaseId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog,
		 Notes, Draft, PreRelease, IsLatest, Source, CreatedAt)
		SELECT Id, Version, Tag, Branch, SourceBranch, CommitSha, Changelog,
		 Notes, Draft, PreRelease, IsLatest, Source, CreatedAt
		FROM Releases`

	if err := db.execV15Rebuild(constants.SQLCreateRelease, copySQL, "Releases"); err != nil {
		return err
	}

	newCount, err := db.countRows(constants.TableRelease)
	if err != nil {
		return fmt.Errorf("count Release: %w", err)
	}

	if oldCount != newCount {
		fmt.Fprintf(os.Stderr, "  ✗ v15 Release row-count mismatch: old=%d new=%d\n", oldCount, newCount)

		return fmt.Errorf(constants.ErrV15ReleaseCountMismatch, oldCount, newCount)
	}

	fmt.Println(constants.MsgV15ReleaseMigrationDone)

	return nil
}

// migrateV15Alias rebuilds Aliases -> Alias with AliasId PK.
func (db *DB) migrateV15Alias() error {
	if !db.tableExists(constants.LegacyTableAliases) {
		return nil
	}

	if db.tableExists(constants.TableAlias) {
		_, _ = db.conn.Exec("DROP TABLE IF EXISTS Aliases")

		return nil
	}

	fmt.Println(constants.MsgV15AliasMigrationStart)

	oldCount, err := db.countRows(constants.LegacyTableAliases)
	if err != nil {
		return fmt.Errorf("count Aliases: %w", err)
	}

	const copySQL = `INSERT INTO Alias
		(AliasId, Alias, RepoId, CreatedAt)
		SELECT Id, Alias, RepoId, CreatedAt
		FROM Aliases`

	if err := db.execV15Rebuild(constants.SQLCreateAlias, copySQL, "Aliases"); err != nil {
		return err
	}

	newCount, err := db.countRows(constants.TableAlias)
	if err != nil {
		return fmt.Errorf("count Alias: %w", err)
	}

	if oldCount != newCount {
		fmt.Fprintf(os.Stderr, "  ✗ v15 Alias row-count mismatch: old=%d new=%d\n", oldCount, newCount)

		return fmt.Errorf(constants.ErrV15AliasCountMismatch, oldCount, newCount)
	}

	fmt.Println(constants.MsgV15AliasMigrationDone)

	return nil
}

// migrateV15Bookmark rebuilds Bookmarks -> Bookmark with BookmarkId PK.
func (db *DB) migrateV15Bookmark() error {
	if !db.tableExists(constants.LegacyTableBookmarks) {
		return nil
	}

	if db.tableExists(constants.TableBookmark) {
		_, _ = db.conn.Exec("DROP TABLE IF EXISTS Bookmarks")

		return nil
	}

	fmt.Println(constants.MsgV15BookmarkMigrationStart)

	oldCount, err := db.countRows(constants.LegacyTableBookmarks)
	if err != nil {
		return fmt.Errorf("count Bookmarks: %w", err)
	}

	const copySQL = `INSERT INTO Bookmark
		(BookmarkId, Name, Command, Args, Flags, CreatedAt)
		SELECT Id, Name, Command, Args, Flags, CreatedAt
		FROM Bookmarks`

	if err := db.execV15Rebuild(constants.SQLCreateBookmark, copySQL, "Bookmarks"); err != nil {
		return err
	}

	newCount, err := db.countRows(constants.TableBookmark)
	if err != nil {
		return fmt.Errorf("count Bookmark: %w", err)
	}

	if oldCount != newCount {
		fmt.Fprintf(os.Stderr, "  ✗ v15 Bookmark row-count mismatch: old=%d new=%d\n", oldCount, newCount)

		return fmt.Errorf(constants.ErrV15BookmarkCountMismatch, oldCount, newCount)
	}

	fmt.Println(constants.MsgV15BookmarkMigrationDone)

	return nil
}

// execV15Rebuild is the shared table-rebuild dance: PRAGMA fk=OFF -> CREATE new
// -> INSERT...SELECT -> DROP legacy -> PRAGMA fk=ON.
func (db *DB) execV15Rebuild(createSQL, copySQL, legacyTable string) error {
	if _, err := db.conn.Exec("PRAGMA foreign_keys = OFF"); err != nil {
		return fmt.Errorf("disable foreign keys: %w", err)
	}

	defer func() {
		_, _ = db.conn.Exec("PRAGMA foreign_keys = ON")
	}()

	if _, err := db.conn.Exec(createSQL); err != nil {
		return fmt.Errorf("create v15 table: %w", err)
	}

	if _, err := db.conn.Exec(copySQL); err != nil {
		return fmt.Errorf("copy %s -> v15: %w", legacyTable, err)
	}

	dropSQL := fmt.Sprintf("DROP TABLE %q", legacyTable)
	if _, err := db.conn.Exec(dropSQL); err != nil {
		return fmt.Errorf("drop legacy %s: %w", legacyTable, err)
	}

	return nil
}
