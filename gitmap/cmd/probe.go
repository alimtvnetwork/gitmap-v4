package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/gitmap/constants"
	"github.com/user/gitmap/model"
	"github.com/user/gitmap/probe"
	"github.com/user/gitmap/store"
)

// runProbe dispatches `gitmap probe [<repo-path>|--all]`. With no args or
// `--all`, every repo in the database is probed sequentially. Phase 2.5
// will replace this with a parallel worker pool.
func runProbe(args []string) {
	checkHelp("probe", args)

	db := openSfDB()
	defer db.Close()

	targets, err := resolveProbeTargets(db, args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if len(targets) == 0 {
		fmt.Print(constants.MsgProbeNoTargets)
		return
	}

	probeAndReport(db, targets)
}

// resolveProbeTargets converts CLI args into a list of repos to probe.
func resolveProbeTargets(db *store.DB, args []string) ([]model.ScanRecord, error) {
	if len(args) == 0 || args[0] == constants.ProbeFlagAll {
		return db.ListRepos()
	}

	absPath, err := filepath.Abs(args[0])
	if err != nil {
		return nil, fmt.Errorf(constants.ErrSFAbsResolve, args[0], err)
	}

	matches, err := db.FindByPath(absPath)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf(constants.ErrProbeNoRepo, absPath)
	}

	return matches, nil
}

// probeAndReport executes RunOne for every target, persists results, and
// prints a per-repo line plus a final summary.
func probeAndReport(db *store.DB, targets []model.ScanRecord) {
	fmt.Printf(constants.MsgProbeStartFmt, len(targets))

	available, unchanged, failed := 0, 0, 0
	for _, repo := range targets {
		url := pickProbeURL(repo)
		if url == "" {
			fmt.Fprintf(os.Stderr, constants.ErrProbeMissingURL+"\n", repo.Slug)
			failed++
			continue
		}

		result := probe.RunOne(url)
		recordProbeResult(db, repo, result)
		available, unchanged, failed = tallyProbe(repo, result, available, unchanged, failed)
	}

	fmt.Printf(constants.MsgProbeDoneFmt, available, unchanged, failed)
}

// pickProbeURL prefers HTTPS (less auth friction in CI), falls back to SSH.
func pickProbeURL(r model.ScanRecord) string {
	if r.HTTPSUrl != "" {
		return r.HTTPSUrl
	}

	return r.SSHUrl
}

// recordProbeResult persists the probe row, logging-but-not-exiting on error.
func recordProbeResult(db *store.DB, repo model.ScanRecord, result probe.Result) {
	if err := db.RecordVersionProbe(result.AsModel(repo.ID)); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

// tallyProbe prints the per-repo line and updates the running counters.
func tallyProbe(repo model.ScanRecord, r probe.Result, ok, none, fail int) (int, int, int) {
	if r.Error != "" {
		fmt.Printf(constants.MsgProbeFailFmt, repo.Slug, r.Error)
		return ok, none, fail + 1
	}
	if r.IsAvailable {
		fmt.Printf(constants.MsgProbeOkFmt, repo.Slug, r.NextVersionTag, r.Method)
		return ok + 1, none, fail
	}
	fmt.Printf(constants.MsgProbeNoneFmt, repo.Slug, r.Method)

	return ok, none + 1, fail
}
