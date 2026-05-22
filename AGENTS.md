# AGENTS.md - ClickUp Integration for Codex CLI

This file configures Codex CLI to work with ClickUp tasks in this repository.

## ClickUp Context

When working on code in this repository, you have access to ClickUp task management
through the awesome-clickup-cli tool. Use it to:

- View task details and requirements
- Update task status as you work
- Add comments with progress updates
- Link commits and PRs to tasks

## Available Commands

```bash
# View current task (auto-detected from branch name)
awesome-clickup-cli git status

# Get full task details
awesome-clickup-cli task get <task-id> --agent

# Update task status
awesome-clickup-cli task update <task-id> --status "in progress"

# Add a comment
awesome-clickup-cli comment create --task-id <task-id> --text "Started implementation"

# Search related tasks
awesome-clickup-cli search "related feature" --agent
```

## Branch Naming Convention

Name branches with ClickUp task IDs for automatic linking:
- `feature/CU-abc123-description`
- `bugfix/CLICKUP-xyz789-fix`
- `#abc123-quick-fix`

## Workflow

1. Before starting: `awesome-clickup-cli git status` to confirm task context
2. While working: Update status with `task update`
3. On commit: Task ID in branch auto-links
4. On PR: `awesome-clickup-cli git link-pr` to connect PR to task
