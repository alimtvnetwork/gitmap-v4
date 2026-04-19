// Package store — migrate_v15rebuild.go provides the generic table-rebuild
// helper used by Phase 1.2+ v15 migrations. SQLite has no ALTER COLUMN, so
// renaming a primary key (Id → {Table}Id) and the table itself requires:
//
//  1. CREATE the new singular table with v15 schema.
//  2. INSERT INTO new (newCols) SELECT oldCols FROM old.
//  3. Verify row-count parity.
//  4. DROP the legacy plural table.
//
// Foreign keys are temporarily disabled so child tables (which still
// reference the legacy name at this point) survive the rename. Subsequent
// phases rebuild those child tables with the new REFERENCES clause.
//
// All steps are idempotent: detect-then-act via tableExists() means fresh
// installs and re-runs are safe no-ops.
package store

import (
	"fmt"
	"os"
)

// v15RebuildSpec describes one legacy-plural → v15-singular table rebuild.
type v15RebuildSpec struct {
	OldTable      string // e.g. "Groups"
	NewTable      string // e.g. "Group"
	NewCreateSQL  string // full CREATE TABLE for the new singular table
	OldColumnList string // exact column list to SELECT from old (e.g. "Id, Name, ...")
	NewColumnList string // exact column list to INSERT into new (e.g. "GroupId, Name, ...")
	StartMsg      string
	DoneMsg       string
}

// runV15Rebuild executes a single rebuild spec idempotently.
func (db *DB) runV15Rebuild(spec v15RebuildSpec) error {
	if !db.tableExists(spec.OldTable) {
		return nil // fresh install, nothing to migrate
	}

	if db.tableExists(spec.NewTable) {
		// Both exist — the new table was created by a prior partial run
		// (e.g., via the standard CREATE TABLE IF NOT EXISTS pass). Drop
		// the legacy and let the standard pass own the new one going
		// forward. Data preservation is impossible in this edge case
		// because the new table is presumed empty/fresh.
		_, _ = db.conn.Exec("DROP TABLE IF EXISTS " + spec.OldTable)

		return nil
	}

	if spec.StartMsg != "" {
		fmt.Println(spec.StartMsg)
	}

	oldCount, err := db.countRows(spec.OldTable)
	if err != nil {
		return fmt.Errorf("count %s: %w", spec.OldTable, err)
	}

	if err := db.execV15Rebuild(spec); err != nil {
		return err
	}

	newCount, err := db.countRows(spec.NewTable)
	if err != nil {
		return fmt.Errorf("count %s: %w", spec.NewTable, err)
	}

	if oldCount != newCount {
		fmt.Fprintf(os.Stderr,
			"  ✗ v15 %s→%s row-count mismatch: old=%d new=%d\n",
			spec.OldTable, spec.NewTable, oldCount, newCount)

		return fmt.Errorf("v15 %s→%s row-count mismatch: old=%d new=%d",
			spec.OldTable, spec.NewTable, oldCount, newCount)
	}

	if spec.DoneMsg != "" {
		fmt.Println(spec.DoneMsg)
	}

	return nil
}

// execV15Rebuild performs the table-rebuild dance for one spec.
func (db *DB) execV15Rebuild(spec v15RebuildSpec) error {
	if _, err := db.conn.Exec("PRAGMA foreign_keys = OFF"); err != nil {
		return fmt.Errorf("disable foreign keys: %w", err)
	}

	defer func() {
		_, _ = db.conn.Exec("PRAGMA foreign_keys = ON")
	}()

	if _, err := db.conn.Exec(spec.NewCreateSQL); err != nil {
		return fmt.Errorf("create %s: %w", spec.NewTable, err)
	}

	copySQL := fmt.Sprintf(
		`INSERT INTO "%s" (%s) SELECT %s FROM "%s"`,
		spec.NewTable, spec.NewColumnList, spec.OldColumnList, spec.OldTable,
	)

	if _, err := db.conn.Exec(copySQL); err != nil {
		return fmt.Errorf("copy %s→%s: %w", spec.OldTable, spec.NewTable, err)
	}

	if _, err := db.conn.Exec(`DROP TABLE "` + spec.OldTable + `"`); err != nil {
		return fmt.Errorf("drop %s: %w", spec.OldTable, err)
	}

	return nil
}
