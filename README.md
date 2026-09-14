# awesome-clickup-cli

A Go CLI for ClickUp with git branch task detection, AI assistant integrations, an MCP server, offline search, and workload analytics.

## What it does

- Wraps 82+ ClickUp API endpoints: tasks, lists, folders, spaces, teams, views, goals, comments.
- Detects the ClickUp task ID from the current git branch name and links PRs/branches to it.
- Generates integration configs for Codex, Claude Code, Hermes Agent, OpenClaw Gateway, and Aider, and can run itself as an MCP server.
- Syncs data to a local SQLite store (FTS5) for offline full-text search.
- Runs workload analytics: stale-task detection, team load distribution, orphaned-task detection.
- Every command accepts `--agent` for JSON output, no color, non-interactive mode, and auto-confirm.

## Install

```bash
go install github.com/sidhartha1s/awesome-clickup-cli@latest
awesome-clickup-cli auth set-token YOUR_API_TOKEN
awesome-clickup-cli doctor
```

## Usage

```bash
# Git integration
awesome-clickup-cli git status          # detect task from branch
awesome-clickup-cli git link-pr         # link current PR to the detected task
awesome-clickup-cli git link-branch     # link branch to task

# Task management
awesome-clickup-cli task get TASK_ID
awesome-clickup-cli task update TASK_ID --status "in progress"
awesome-clickup-cli list task create LIST_ID --name "New task"
awesome-clickup-cli task comment create TASK_ID --comment-text "Your comment"

# Search and analytics
awesome-clickup-cli search "keyword" --agent
awesome-clickup-cli stale --days 7
awesome-clickup-cli load
awesome-clickup-cli orphans

# Offline
awesome-clickup-cli sync
awesome-clickup-cli search "query" --data-source local

# AI assistant integrations
awesome-clickup-cli integrations detect
awesome-clickup-cli integrations all
awesome-clickup-cli mcp-server
```

Branch naming patterns that auto-detect a task: `feature/CU-abc123-description`, `bugfix/CLICKUP-xyz789-fix`, `#abc123-quick-fix`.

## Layout

| Path | Role |
|------|------|
| `cmd/` | Binary entry points |
| `internal/` | Core CLI logic (256 files: commands, API client, cache, MCP tool defs) |
| `CLAUDE.md`, `AGENTS.md`, `SKILL.md` | Generated integration docs for Claude Code, Codex, and Claude Code skills |
| `CONTRIBUTING.md`, `LICENSE`, `NOTICE` | Apache-2.0 project files |
| `Makefile` | Build, test, lint, install targets |

## Notes / Gotchas

- Credentials are stored with `0o600` permissions via OS keyring (`go-keyring`); no plaintext fallback, token never printed.
- Config profiles: `awesome-clickup-cli profile save default --compact --json`, then `--profile default` on later commands.
- MCP server mode: `awesome-clickup-cli mcp-server`, registered in `~/.claude.json` under `mcpServers.clickup`.
- The Makefile's `build`/`install` targets point at `./cmd/clickup-reference-pp-cli` and `./cmd/clickup-reference-pp-mcp`, not `./cmd/awesome-clickup-cli`. A commit renaming these directories ("Fix: rename cmd directories for proper go install") was later reverted, so confirm the actual `cmd/` subdirectory names before trusting `make build` or `make install`.
- License: Apache-2.0 (`LICENSE`, `NOTICE`).
