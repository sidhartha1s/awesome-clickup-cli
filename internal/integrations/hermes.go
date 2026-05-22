package integrations

import (
	"fmt"
	"os"
	"path/filepath"
)

// HermesPluginConfig holds configuration for Hermes plugin generation
type HermesPluginConfig struct {
	PluginName  string
	Description string
	BinaryPath  string
}

// GenerateHermesPlugin creates a Hermes Agent plugin for ClickUp integration
func GenerateHermesPlugin(cfg HermesPluginConfig) map[string]string {
	if cfg.PluginName == "" {
		cfg.PluginName = "clickup"
	}
	if cfg.Description == "" {
		cfg.Description = "ClickUp task management integration"
	}

	files := make(map[string]string)

	// manifest.yaml
	files["manifest.yaml"] = fmt.Sprintf(`name: %s
version: 1.0.0
description: %s
author: awesome-clickup-cli
homepage: https://github.com/sidhartha1s/awesome-clickup-cli

capabilities:
  - tools
  - cli

dependencies: []

config:
  api_token:
    type: secret
    required: true
    description: ClickUp API token
`, cfg.PluginName, cfg.Description)

	// __init__.py
	files["__init__.py"] = `"""ClickUp integration plugin for Hermes Agent."""

from .tools import register
from .cli import register_cli

__all__ = ["register", "register_cli"]
`

	// tools.py
	binaryPath := cfg.BinaryPath
	if binaryPath == "" {
		binaryPath = "awesome-clickup-cli"
	}

	files["tools.py"] = fmt.Sprintf(`"""ClickUp tools for Hermes Agent."""

import subprocess
import json
from typing import Any

BINARY = "%s"


def _run_cli(*args: str) -> dict[str, Any]:
    """Run awesome-clickup-cli and return JSON output."""
    cmd = [BINARY] + list(args) + ["--agent"]
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=30)
        if result.returncode == 0:
            return json.loads(result.stdout)
        return {"error": result.stderr or "Command failed"}
    except subprocess.TimeoutExpired:
        return {"error": "Command timed out"}
    except json.JSONDecodeError:
        return {"output": result.stdout}


def task_get(task_id: str) -> dict[str, Any]:
    """Get details for a ClickUp task.

    Args:
        task_id: The ClickUp task ID (e.g., abc123)

    Returns:
        Task details including name, status, assignees, due date
    """
    return _run_cli("task", "get", task_id)


def task_update(task_id: str, status: str = None, priority: str = None) -> dict[str, Any]:
    """Update a ClickUp task's status or priority.

    Args:
        task_id: The ClickUp task ID
        status: New status (e.g., "in progress", "complete")
        priority: New priority (1=urgent, 2=high, 3=normal, 4=low)

    Returns:
        Updated task details
    """
    args = ["task", "update", task_id]
    if status:
        args.extend(["--status", status])
    if priority:
        args.extend(["--priority", priority])
    return _run_cli(*args)


def task_search(query: str, limit: int = 10) -> dict[str, Any]:
    """Search ClickUp tasks by keyword.

    Args:
        query: Search query
        limit: Maximum results to return

    Returns:
        List of matching tasks
    """
    return _run_cli("search", query, "--limit", str(limit))


def comment_create(task_id: str, text: str) -> dict[str, Any]:
    """Add a comment to a ClickUp task.

    Args:
        task_id: The ClickUp task ID
        text: Comment text

    Returns:
        Created comment details
    """
    return _run_cli("comment", "create", "--task-id", task_id, "--text", text)


def git_status() -> dict[str, Any]:
    """Get ClickUp task linked to current git branch.

    Returns:
        Task ID and URL if detected from branch name
    """
    return _run_cli("git", "status")


def stale_tasks(days: int = 7) -> dict[str, Any]:
    """Find tasks not updated in N days.

    Args:
        days: Number of days since last update

    Returns:
        List of stale tasks
    """
    return _run_cli("stale", "--days", str(days))


def workload() -> dict[str, Any]:
    """Get team workload distribution.

    Returns:
        Task counts per assignee
    """
    return _run_cli("load")


def register(ctx):
    """Register ClickUp tools with Hermes."""
    ctx.register_tool(
        name="clickup_task_get",
        description="Get details for a ClickUp task",
        handler=task_get,
        parameters={
            "task_id": {"type": "string", "required": True, "description": "ClickUp task ID"}
        }
    )

    ctx.register_tool(
        name="clickup_task_update",
        description="Update a ClickUp task's status or priority",
        handler=task_update,
        parameters={
            "task_id": {"type": "string", "required": True},
            "status": {"type": "string", "required": False},
            "priority": {"type": "string", "required": False}
        }
    )

    ctx.register_tool(
        name="clickup_search",
        description="Search ClickUp tasks by keyword",
        handler=task_search,
        parameters={
            "query": {"type": "string", "required": True},
            "limit": {"type": "integer", "required": False, "default": 10}
        }
    )

    ctx.register_tool(
        name="clickup_comment",
        description="Add a comment to a ClickUp task",
        handler=comment_create,
        parameters={
            "task_id": {"type": "string", "required": True},
            "text": {"type": "string", "required": True}
        }
    )

    ctx.register_tool(
        name="clickup_git_status",
        description="Get ClickUp task linked to current git branch",
        handler=git_status,
        parameters={}
    )

    ctx.register_tool(
        name="clickup_stale_tasks",
        description="Find tasks not updated recently",
        handler=stale_tasks,
        parameters={
            "days": {"type": "integer", "required": False, "default": 7}
        }
    )

    ctx.register_tool(
        name="clickup_workload",
        description="Get team workload distribution",
        handler=workload,
        parameters={}
    )
`, binaryPath)

	// cli.py
	files["cli.py"] = fmt.Sprintf(`"""CLI commands for ClickUp Hermes plugin."""

import subprocess
import sys

BINARY = "%s"


def _clickup_command(args):
    """Handler for hermes clickup <subcommand>."""
    sub = getattr(args, "clickup_cmd", None)
    extra = getattr(args, "extra_args", [])

    if sub == "status":
        _run(["git", "status"])
    elif sub == "task":
        if extra:
            _run(["task", "get"] + extra)
        else:
            print("Usage: hermes clickup task <task-id>")
    elif sub == "search":
        if extra:
            _run(["search"] + extra)
        else:
            print("Usage: hermes clickup search <query>")
    elif sub == "stale":
        _run(["stale", "--days", "7"])
    elif sub == "workload":
        _run(["load"])
    elif sub == "doctor":
        _run(["doctor"])
    else:
        print("Available commands: status, task, search, stale, workload, doctor")


def _run(args: list[str]):
    """Run awesome-clickup-cli with given arguments."""
    cmd = [BINARY] + args
    try:
        subprocess.run(cmd, check=False)
    except FileNotFoundError:
        print(f"Error: {BINARY} not found. Install it first.")
        sys.exit(1)


def register_cli(subparser):
    """Register hermes clickup CLI commands."""
    subparser.add_argument("clickup_cmd", nargs="?", default=None,
                          choices=["status", "task", "search", "stale", "workload", "doctor"],
                          help="ClickUp command to run")
    subparser.add_argument("extra_args", nargs="*", help="Additional arguments")
    subparser.set_defaults(func=_clickup_command)
`, binaryPath)

	return files
}

// WriteHermesPlugin writes the Hermes plugin files to a directory
func WriteHermesPlugin(baseDir string, cfg HermesPluginConfig) error {
	files := GenerateHermesPlugin(cfg)

	pluginDir := filepath.Join(baseDir, cfg.PluginName)
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

// DetectHermesInstall checks if Hermes Agent is installed
func DetectHermesInstall() (string, bool) {
	paths := []string{
		"hermes",
		filepath.Join(os.Getenv("HOME"), ".local", "bin", "hermes"),
		"/usr/local/bin/hermes",
	}

	for _, p := range paths {
		if p == "hermes" {
			if path, err := LookPath("hermes"); err == nil {
				return path, true
			}
		} else if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	return "", false
}
