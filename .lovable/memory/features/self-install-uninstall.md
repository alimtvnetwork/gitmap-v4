---
name: Self Install Uninstall
description: gitmap self-install / self-uninstall manage the gitmap binary itself (separate from third-party install/uninstall). Embedded scripts via go:embed, Windows handoff, marker-block PATH cleanup.
type: feature
---

`gitmap self-install` and `gitmap self-uninstall` were added because
`install`/`uninstall` are reserved for third-party tools (npp, vscode,
dev tools).

## self-uninstall removes
- Binary + deploy dir artefacts (`gitmap`, `gitmap.exe`, `gitmap-handoff-*`, `*.old`, `gitmap-completion.*`)
- `.gitmap/` data dir (skip `--keep-data`)
- PATH snippet marker block in shell profile (skip `--keep-snippet`)
- Completion files in deploy dir
- Requires `--confirm` or interactive `yes`

## Windows handoff
When the running .exe lives in the deploy dir, copies itself to
`%TEMP%\gitmap-handoff-<pid>.exe`, re-execs hidden
`self-uninstall-runner` verb, then schedules self-deletion via
`cmd.exe /C ping ... & del`.

## self-install
- Default dir: Windows `D:\gitmap`, Unix `~/.local/bin/gitmap`
- Always prompts unless `--dir` / `--yes`
- Loads installer from embedded `go:embed` first, falls back to
  `raw.githubusercontent.com/alimtvnetwork/gitmap-v3/main/gitmap/scripts/install.{ps1,sh}`
- PowerShell scripts written with UTF-8 BOM
- Forwards `--version <tag>` to the installer

## Files
- `gitmap/constants/constants_selfinstall.go` — IDs, messages, defaults
- `gitmap/scripts/embed.go` — `go:embed install.ps1 install.sh uninstall.ps1`
- `gitmap/cmd/selfinstall.go` — install command
- `gitmap/cmd/selfuninstall.go`, `selfuninstallparts.go`, `selfuninstallhandoff.go` — uninstall command (split for <200 line rule)
- `gitmap/helptext/self-install.md`, `self-uninstall.md`
- `spec/01-app/90-self-install-uninstall.md`

Spec: spec/01-app/90-self-install-uninstall.md
