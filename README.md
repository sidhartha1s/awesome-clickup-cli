# Awesome ClickUp CLI

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

**The most feature-complete ClickUp CLI with git integration and AI assistant support.**

Built with security-first design. Integrates seamlessly with Codex, Claude, Hermes, OpenClaw, and Aider.

## Why This CLI?

| Feature | This CLI | Others |
|---------|----------|--------|
| **Git integration** | ✅ Auto-detect task from branch, link PRs | ❌ |
| **AI integrations** | ✅ Codex, Claude, Hermes, OpenClaw, Aider | ❌ |
| **MCP server** | ✅ Model Context Protocol support | ❌ |
| **Secure credentials** | ✅ 0o600 permissions, keyring | ⚠️ Plaintext |
| **Offline search** | ✅ FTS5 SQLite | Partial |
| **API coverage** | 82+ endpoints | ~45 |

## Quick Start

```bash
# Install
go install github.com/sidhartha1s/awesome-clickup-cli@latest

# Authenticate
awesome-clickup-cli auth set-token YOUR_API_TOKEN

# Verify
awesome-clickup-cli doctor
```

## Git Integration

Auto-detect ClickUp tasks from your branch name:

```bash
# On branch: feature/CU-abc123-add-login
awesome-clickup-cli git status
# → Task ID: abc123
# → URL: https://app.clickup.com/t/abc123

# Link your PR to the task
awesome-clickup-cli git link-pr

# Link branch to task
awesome-clickup-cli git link-branch
```

Supported patterns:
- `feature/CU-abc123-description`
- `bugfix/CLICKUP-xyz789-fix`
- `#abc123-quick-fix`

## AI Assistant Integrations

Generate integration configs for your favorite AI coding assistants:

```bash
# Detect installed AI tools
awesome-clickup-cli integrations detect

# Generate all integrations
awesome-clickup-cli integrations all

# Or generate specific ones:
awesome-clickup-cli integrations codex     # AGENTS.md for Codex CLI
awesome-clickup-cli integrations claude    # CLAUDE.md + MCP config
awesome-clickup-cli integrations hermes    # Python plugin for Hermes Agent
awesome-clickup-cli integrations openclaw  # Python plugin for OpenClaw Gateway
awesome-clickup-cli integrations aider     # .aider.conf.yml
```

### MCP Server Mode

Run as an MCP server for any Model Context Protocol client:

```bash
awesome-clickup-cli mcp-server
```

Add to Claude's `~/.claude.json`:
```json
{
  "mcpServers": {
    "clickup": {
      "command": "awesome-clickup-cli",
      "args": ["mcp-server"]
    }
  }
}
```

Add to Codex:
```bash
codex mcp add clickup --command "awesome-clickup-cli mcp-server"
```

## Core Features

### Task Management

```bash
# Get task details
awesome-clickup-cli task get TASK_ID

# List tasks in a list
awesome-clickup-cli list task get LIST_ID

# Create a task in a list
awesome-clickup-cli list task create LIST_ID --name "New task"

# Update task status
awesome-clickup-cli task update TASK_ID --status "in progress"

# Add a comment
awesome-clickup-cli task comment create TASK_ID --comment-text "Your comment"

# Search tasks
awesome-clickup-cli search "keyword" --agent
```

### Analytics

```bash
# Find stale tasks
awesome-clickup-cli stale --days 7

# Team workload analysis
awesome-clickup-cli load

# Find orphaned tasks
awesome-clickup-cli orphans
```

### Offline Capabilities

```bash
# Sync data locally
awesome-clickup-cli sync

# Search offline
awesome-clickup-cli search "query" --data-source local
```

## Agent Mode

All commands support `--agent` flag for AI assistant integration:

```bash
awesome-clickup-cli task get TASK_ID --agent
```

This enables:
- JSON output
- Non-interactive mode
- Compact response
- No color codes
- Auto-confirm prompts

## Security

- Credentials stored with **0o600 permissions** (owner-only)
- No plaintext fallback
- OS keyring support via [go-keyring](https://github.com/zalando/go-keyring)
- Token never exposed in command output

## All Commands

### Task Management
- `task get <id>` - Get task details
- `task update <id>` - Update task
- `task delete <id>` - Delete task
- `list task get <list_id>` - List tasks in a list
- `list task create <list_id>` - Create task in list
- `task comment create <task_id>` - Add comment
- `task comment get <task_id>` - List comments

### Organization
- `team list` - List workspaces
- `team space get <team_id>` - List spaces in workspace
- `space get <space_id>` - Get space details
- `space folder get <space_id>` - List folders in space
- `folder list get <folder_id>` - List lists in folder
- `list get <list_id>` - Get list details

### Views & Goals
- `view list`, `view get`, `view create`
- `goal list`, `goal get`, `goal create`

### Git Integration
- `git status` - Detect task from branch
- `git link-pr` - Link PR to task
- `git link-branch` - Link branch to task

### Analytics
- `stale` - Find stale tasks
- `load` - Workload analysis
- `orphans` - Find orphaned tasks
- `search` - Full-text search

### AI Integrations
- `integrations detect` - Find installed AI tools
- `integrations all` - Generate all configs
- `integrations codex` - Codex AGENTS.md
- `integrations claude` - Claude CLAUDE.md + MCP
- `integrations hermes` - Hermes plugin
- `integrations openclaw` - OpenClaw plugin
- `integrations aider` - Aider config
- `mcp-server` - Run as MCP server

### Utility
- `doctor` - Health check
- `sync` - Sync data locally
- `auth set-token` - Set API token
- `auth status` - Check auth status
- `profile save/list` - Save command profiles

## Configuration

```bash
# Set default profile
awesome-clickup-cli profile save default --compact --json

# Use profile
awesome-clickup-cli task list --profile default
```

## Contributing

Contributions welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) first.

## License

Apache-2.0. See [LICENSE](LICENSE).
