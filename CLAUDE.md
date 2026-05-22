# CLAUDE.md - ClickUp Integration

This project uses ClickUp for task management. The awesome-clickup-cli tool is available
for viewing and managing tasks directly from the command line.

## Quick Reference

```bash
# View task linked to current branch
awesome-clickup-cli git status

# Get task details
awesome-clickup-cli task get <task-id> --agent

# Update task status
awesome-clickup-cli task update <task-id> --status "in progress"

# Search tasks
awesome-clickup-cli search "keyword" --agent

# Find stale tasks
awesome-clickup-cli stale --days 7 --agent
```

## Branch Naming

Name branches with task IDs for automatic detection:
- `feature/CU-abc123-description`
- `bugfix/CLICKUP-xyz789-fix`
- `#abc123-quick-fix`

## Workflow

1. `awesome-clickup-cli git status` - Confirm which task you're on
2. Work on the code
3. `awesome-clickup-cli task update <id> --status 'in progress'` - Update status
4. `awesome-clickup-cli git link-pr` - Link PR to task after push

## MCP Server

This CLI can run as an MCP server for richer integration:

```bash
awesome-clickup-cli mcp-server
```

Add to your Claude MCP config (~/.claude.json):

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
