# probe

Run the hybrid HEAD-then-clone version probe against one or every repo.

## Synopsis

```
gitmap probe                   # probe every repo in the database
gitmap probe --all             # explicit form of the above
gitmap probe <repo-path>       # probe a single repo by absolute path
gitmap probe --all --json      # emit a JSON array (CI-friendly)
```

## What it does

For each target repo, the probe inspects the remote and looks for the
highest semver-style tag (`vN.N.N` or `N.N.N`). Results are persisted
into the `VersionProbe` table — one row per probe, never overwriting
history. `gitmap find-next` then surfaces every repo whose **latest**
probe row reports an available update.

## Strategy

Two strategies, tried in order. The probe falls through to the next
only when the previous one fails or returns zero tags:

| # | Method | Command | Why it can fail |
|---|---|---|---|
| 1 | `ls-remote` | `git ls-remote --tags --sort=-v:refname <url>` | Server rejects unauthenticated probes, returns zero tags, or git exits non-zero |
| 2 | `shallow-clone` | `git clone --depth 1 --filter=blob:none --no-checkout <url>` then `git tag --sort=-v:refname` | Network/auth failure |

The shallow-clone fallback is **treeless** (`--filter=blob:none`) and
**checkout-less** (`--no-checkout`) so we only pay for the refs
database — no working tree, no blobs. The temp directory is removed
before return.

## URL preference

`HTTPSUrl` is preferred over `SSHUrl` — HTTPS has less auth friction
in CI and on first-time-ever clones. SSH only kicks in when HTTPS is
empty.

## Output

Per-repo line format:

| Symbol | Meaning |
|---|---|
| `✓ <slug> → v1.2.3 (method=ls-remote)` | New tag found |
| `· <slug> → no new version (method=ls-remote)` | Probe ran, no tag |
| `✗ <slug> → <error>` | Probe failed (error captured in `VersionProbe.Error`) |

Final summary:

```
✓ Probe complete: <available> available, <unchanged> unchanged, <failed> failed.
```

## Examples

```
$ gitmap probe --all
→ Probing 12 repo(s)...
  ✓ awesome-cli → v2.4.0 (method=ls-remote)
  · helper-lib → no new version (method=ls-remote)
  ✗ private-repo → ls-remote failed: exit status 128
  ✓ infra-tools → v0.9.1 (method=shallow-clone)
✓ Probe complete: 2 available, 1 unchanged, 1 failed.

$ gitmap probe E:\src\awesome-cli
→ Probing 1 repo(s)...
  ✓ awesome-cli → v2.4.0 (method=ls-remote)
✓ Probe complete: 1 available, 0 unchanged, 0 failed.
```

## Database

Each invocation appends one row per probed repo to `VersionProbe` —
never overwrites. Inspect history with the SQLite CLI:

```
sqlite> SELECT RepoId, NextVersionTag, Method, ProbedAt
        FROM VersionProbe ORDER BY ProbedAt DESC LIMIT 10;
```

## See also

- `gitmap find-next` (alias `fn`) — read the latest probe results
- `gitmap sf list` — show scan folders / repo membership
- `gitmap pull` / `gitmap cn next all` — apply the upgrades
