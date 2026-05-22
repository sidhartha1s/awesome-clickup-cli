package integrations

import (
	"fmt"
	"os"
	"path/filepath"
)

// OpenClawConfig holds configuration for OpenClaw integration
type OpenClawConfig struct {
	AgentName   string
	BinaryPath  string
	Description string
}

// GenerateOpenClawPlugin creates an OpenClaw plugin for ClickUp integration
func GenerateOpenClawPlugin(cfg OpenClawConfig) map[string]string {
	if cfg.AgentName == "" {
		cfg.AgentName = "clickup"
	}
	if cfg.Description == "" {
		cfg.Description = "ClickUp task management via awesome-clickup-cli"
	}

	files := make(map[string]string)

	binaryPath := cfg.BinaryPath
	if binaryPath == "" {
		binaryPath = "awesome-clickup-cli"
	}

	// manifest.yaml
	files["manifest.yaml"] = fmt.Sprintf(`name: clickup
version: 1.0.0
description: %s
author: sidhartha1s
homepage: https://github.com/sidhartha1s/awesome-clickup-cli

type: tool-provider

capabilities:
  tools:
    - clickup_task_get
    - clickup_task_update
    - clickup_task_search
    - clickup_comment_add
    - clickup_git_status
    - clickup_stale_tasks
    - clickup_workload

config:
  binary_path:
    type: string
    default: "%s"
    description: Path to awesome-clickup-cli binary
`, cfg.Description, binaryPath)

	// tools.py
	files["tools.py"] = fmt.Sprintf(`"""ClickUp tools for OpenClaw Gateway."""

import subprocess
import json
from typing import Any, Optional

BINARY = "%s"


def _run_cli(*args: str, timeout: int = 30) -> dict[str, Any]:
    """Execute awesome-clickup-cli and parse JSON output."""
    cmd = [BINARY] + list(args) + ["--agent"]
    try:
        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            timeout=timeout
        )
        if result.returncode == 0:
            try:
                return json.loads(result.stdout)
            except json.JSONDecodeError:
                return {"output": result.stdout.strip()}
        return {"error": result.stderr.strip() or f"Exit code {result.returncode}"}
    except subprocess.TimeoutExpired:
        return {"error": f"Command timed out after {timeout}s"}
    except FileNotFoundError:
        return {"error": f"Binary not found: {BINARY}"}


class ClickUpTools:
    """ClickUp tool implementations for OpenClaw."""

    @staticmethod
    def task_get(task_id: str) -> dict[str, Any]:
        """Retrieve full details for a ClickUp task.

        Args:
            task_id: ClickUp task identifier (e.g., 'abc123', 'CU-abc123')

        Returns:
            Task object with name, description, status, assignees, dates, etc.
        """
        # Strip common prefixes
        task_id = task_id.upper().replace("CU-", "").replace("CLICKUP-", "").lstrip("#")
        return _run_cli("task", "get", task_id)

    @staticmethod
    def task_update(
        task_id: str,
        status: Optional[str] = None,
        priority: Optional[str] = None,
        assignee: Optional[str] = None
    ) -> dict[str, Any]:
        """Update a ClickUp task.

        Args:
            task_id: Task identifier
            status: New status name (e.g., 'in progress', 'complete', 'blocked')
            priority: Priority level (1=urgent, 2=high, 3=normal, 4=low)
            assignee: User ID or email to assign

        Returns:
            Updated task object
        """
        args = ["task", "update", task_id]
        if status:
            args.extend(["--status", status])
        if priority:
            args.extend(["--priority", str(priority)])
        if assignee:
            args.extend(["--assignee", assignee])
        return _run_cli(*args)

    @staticmethod
    def task_search(query: str, limit: int = 10, status: Optional[str] = None) -> dict[str, Any]:
        """Search ClickUp tasks by keyword.

        Args:
            query: Search terms
            limit: Maximum results (default 10)
            status: Filter by status

        Returns:
            List of matching tasks with id, name, status, url
        """
        args = ["search", query, "--limit", str(limit)]
        if status:
            args.extend(["--status", status])
        return _run_cli(*args)

    @staticmethod
    def comment_add(task_id: str, text: str, notify: bool = False) -> dict[str, Any]:
        """Add a comment to a ClickUp task.

        Args:
            task_id: Task identifier
            text: Comment body (supports markdown)
            notify: Whether to notify assignees

        Returns:
            Created comment object
        """
        args = ["comment", "create", "--task-id", task_id, "--text", text]
        if notify:
            args.append("--notify")
        return _run_cli(*args)

    @staticmethod
    def git_status() -> dict[str, Any]:
        """Detect ClickUp task from current git branch.

        Parses branch names for patterns like:
        - feature/CU-abc123-description
        - bugfix/CLICKUP-xyz789-fix
        - #abc123-quick-fix

        Returns:
            Object with branch, task_id, task_url, detected (bool)
        """
        return _run_cli("git", "status")

    @staticmethod
    def stale_tasks(days: int = 7, assignee: Optional[str] = None) -> dict[str, Any]:
        """Find tasks not updated within N days.

        Args:
            days: Staleness threshold (default 7)
            assignee: Filter by assignee

        Returns:
            List of stale tasks with last_updated timestamp
        """
        args = ["stale", "--days", str(days)]
        if assignee:
            args.extend(["--assignee", assignee])
        return _run_cli(*args)

    @staticmethod
    def workload(space_id: Optional[str] = None) -> dict[str, Any]:
        """Get task distribution across team members.

        Args:
            space_id: Limit to specific space

        Returns:
            Workload per assignee with task counts by status
        """
        args = ["load"]
        if space_id:
            args.extend(["--space", space_id])
        return _run_cli(*args)


def register(ctx):
    """Register ClickUp tools with OpenClaw Gateway."""
    tools = ClickUpTools()

    ctx.register_tool(
        "clickup_task_get",
        tools.task_get,
        description="Get full details for a ClickUp task",
        parameters={
            "task_id": {"type": "string", "required": True}
        }
    )

    ctx.register_tool(
        "clickup_task_update",
        tools.task_update,
        description="Update a ClickUp task's status, priority, or assignee",
        parameters={
            "task_id": {"type": "string", "required": True},
            "status": {"type": "string"},
            "priority": {"type": "string"},
            "assignee": {"type": "string"}
        }
    )

    ctx.register_tool(
        "clickup_task_search",
        tools.task_search,
        description="Search ClickUp tasks by keyword",
        parameters={
            "query": {"type": "string", "required": True},
            "limit": {"type": "integer", "default": 10},
            "status": {"type": "string"}
        }
    )

    ctx.register_tool(
        "clickup_comment_add",
        tools.comment_add,
        description="Add a comment to a ClickUp task",
        parameters={
            "task_id": {"type": "string", "required": True},
            "text": {"type": "string", "required": True},
            "notify": {"type": "boolean", "default": False}
        }
    )

    ctx.register_tool(
        "clickup_git_status",
        tools.git_status,
        description="Detect ClickUp task from current git branch name"
    )

    ctx.register_tool(
        "clickup_stale_tasks",
        tools.stale_tasks,
        description="Find tasks not updated within N days",
        parameters={
            "days": {"type": "integer", "default": 7},
            "assignee": {"type": "string"}
        }
    )

    ctx.register_tool(
        "clickup_workload",
        tools.workload,
        description="Get task distribution across team members",
        parameters={
            "space_id": {"type": "string"}
        }
    )
`, binaryPath)

	return files
}

// WriteOpenClawPlugin writes the OpenClaw plugin files to a directory
func WriteOpenClawPlugin(baseDir string, cfg OpenClawConfig) error {
	files := GenerateOpenClawPlugin(cfg)

	pluginDir := filepath.Join(baseDir, "clickup")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("creating plugin directory: %w", err)
	}

	for filename, content := range files {
		path := filepath.Join(pluginDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", filename, err)
		}
	}

	return nil
}

// GenerateOpenClawInstallInstructions returns setup instructions
func GenerateOpenClawInstallInstructions(pluginDir string) string {
	return fmt.Sprintf(`# OpenClaw ClickUp Integration

## Installation

1. Copy the plugin to your OpenClaw plugins directory:
   ` + "```bash" + `
   cp -r %s/clickup ~/.openclaw/plugins/
   ` + "```" + `

2. Enable the plugin:
   ` + "```bash" + `
   openclaw plugins enable clickup
   openclaw gateway restart
   ` + "```" + `

3. Configure your agent to use ClickUp tools:
   ` + "```bash" + `
   openclaw agents config --agent main --tools clickup_task_get,clickup_search
   ` + "```" + `

## Usage

The plugin provides these tools to your OpenClaw agents:

- ` + "`clickup_task_get`" + ` - Get task details
- ` + "`clickup_task_update`" + ` - Update task status/priority
- ` + "`clickup_task_search`" + ` - Search tasks
- ` + "`clickup_comment_add`" + ` - Add comments
- ` + "`clickup_git_status`" + ` - Detect task from branch
- ` + "`clickup_stale_tasks`" + ` - Find stale tasks
- ` + "`clickup_workload`" + ` - Team workload analysis
`, pluginDir)
}

// DetectOpenClawInstall checks if OpenClaw is installed
func DetectOpenClawInstall() (string, bool) {
	paths := []string{
		"openclaw",
		filepath.Join(os.Getenv("HOME"), ".local", "bin", "openclaw"),
		"/usr/local/bin/openclaw",
	}

	for _, p := range paths {
		if p == "openclaw" {
			if path, err := LookPath("openclaw"); err == nil {
				return path, true
			}
		} else if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	return "", false
}
