---
name: awesome-clickup-cli
description: "ClickUp CLI with git integration, AI assistant support (Codex, Claude, Hermes, OpenClaw, Aider), MCP server, offline search, and workload analytics"
author: "sidhartha1s"
license: "Apache-2.0"
argument-hint: "<command> [args]"
allowed-tools: "Read Bash"
trigger-phrases:
  - clickup task
  - link PR to clickup
  - check clickup branch
  - clickup workload
  - use clickup
  - clickup integration
metadata:
  openclaw:
    requires:
      bins:
        - awesome-clickup-cli
---

# Awesome ClickUp CLI

Secure ClickUp CLI with git integration, AI assistant integrations, offline search, and workload analytics.

## Prerequisites: Install the CLI

```bash
go install github.com/sidhartha1s/awesome-clickup-cli/cmd/awesome-clickup-cli@latest
```

Then authenticate:
```bash
awesome-clickup-cli auth set-token YOUR_API_TOKEN
awesome-clickup-cli doctor
```

## When to Use This CLI

Use when you need:
- **Git-integrated workflows** — Auto-detect task from branch, link PRs
- **AI assistant integration** — MCP server, Codex, Claude, Hermes, OpenClaw, Aider
- **Secure credential storage** — 0o600 permissions, no plaintext
- **Offline search** — FTS5 SQLite-backed search
- **Team analytics** — Workload distribution, stale task detection

## Unique Capabilities

### Git Integration

```bash
# Detect task from current branch (feature/CU-abc123-description)
awesome-clickup-cli git status

# Link current PR to detected task
awesome-clickup-cli git link-pr

# Link branch to task
awesome-clickup-cli git link-branch
```

### AI Assistant Integrations

```bash
# Detect installed AI tools
awesome-clickup-cli integrations detect

# Generate all integration configs
awesome-clickup-cli integrations all

# Individual integrations
awesome-clickup-cli integrations codex     # AGENTS.md
awesome-clickup-cli integrations claude    # CLAUDE.md + MCP config
awesome-clickup-cli integrations hermes    # Python plugin
awesome-clickup-cli integrations openclaw  # Python plugin
awesome-clickup-cli integrations aider     # .aider.conf.yml

# Run as MCP server
awesome-clickup-cli mcp-server
```

### Analytics

```bash
# Find stale tasks
awesome-clickup-cli stale --days 7 --agent

# Team workload distribution
awesome-clickup-cli load --agent

# Find orphaned tasks
awesome-clickup-cli orphans --agent
```

### Offline Search

```bash
# Sync data locally
awesome-clickup-cli sync

# Search offline
awesome-clickup-cli search "keyword" --data-source local --agent
```

## Command Reference

### Git Integration
- `awesome-clickup-cli git status` — Show ClickUp task from current branch
- `awesome-clickup-cli git link-pr` — Link GitHub PR to task
- `awesome-clickup-cli git link-branch` — Link branch to task

### Task Management
- `awesome-clickup-cli task get <id>` — Get task details
- `awesome-clickup-cli list task get <list_id>` — List tasks in a list
- `awesome-clickup-cli list task create <list_id> --name "Name"` — Create task
- `awesome-clickup-cli task update <id> --status "in progress"` — Update task
- `awesome-clickup-cli task comment create <id> --comment-text "Comment"` — Add comment

### Search & Analytics
- `awesome-clickup-cli search "query"` — Full-text search
- `awesome-clickup-cli stale --days 7` — Find stale tasks
- `awesome-clickup-cli load` — Workload analysis
- `awesome-clickup-cli orphans` — Find orphaned tasks

### Organization
- `awesome-clickup-cli team list` — List workspaces
- `awesome-clickup-cli team space get <team_id>` — List spaces
- `awesome-clickup-cli space folder get <space_id>` — List folders
- `awesome-clickup-cli folder list get <folder_id>` — List lists

### AI Integrations
- `awesome-clickup-cli integrations detect` — Detect installed AI tools
- `awesome-clickup-cli integrations all` — Generate all configs
- `awesome-clickup-cli integrations codex` — Generate AGENTS.md
- `awesome-clickup-cli integrations claude` — Generate CLAUDE.md + MCP
- `awesome-clickup-cli integrations hermes` — Generate Hermes plugin
- `awesome-clickup-cli integrations openclaw` — Generate OpenClaw plugin
- `awesome-clickup-cli integrations aider` — Generate .aider.conf.yml
- `awesome-clickup-cli mcp-server` — Run as MCP server

### Utility
- `awesome-clickup-cli doctor` — Health check
- `awesome-clickup-cli sync` — Sync data locally
- `awesome-clickup-cli auth set-token <token>` — Set API token
- `awesome-clickup-cli auth status` — Check authentication

## Agent Mode

All commands support `--agent` for JSON output and non-interactive mode:

```bash
awesome-clickup-cli task get abc123 --agent
awesome-clickup-cli search "bug" --agent --select id,name,status
```
