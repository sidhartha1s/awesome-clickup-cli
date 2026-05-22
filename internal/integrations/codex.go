package integrations

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CodexConfig holds configuration for Codex CLI integration
type CodexConfig struct {
	ProjectRoot string
	TaskID      string
	ListID      string
	SpaceID     string
}

// GenerateAgentsMD creates an AGENTS.md file for Codex CLI integration
func GenerateAgentsMD(cfg CodexConfig) (string, error) {
	content := `# AGENTS.md - ClickUp Integration for Codex CLI

This file configures Codex CLI to work with ClickUp tasks in this repository.

## ClickUp Context

When working on code in this repository, you have access to ClickUp task management
through the awesome-clickup-cli tool. Use it to:

- View task details and requirements
- Update task status as you work
- Add comments with progress updates
- Link commits and PRs to tasks

## Available Commands

` + "```bash" + `
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
` + "```" + `

## Branch Naming Convention

Name branches with ClickUp task IDs for automatic linking:
- ` + "`feature/CU-abc123-description`" + `
- ` + "`bugfix/CLICKUP-xyz789-fix`" + `
- ` + "`#abc123-quick-fix`" + `

## Workflow

1. Before starting: ` + "`awesome-clickup-cli git status`" + ` to confirm task context
2. While working: Update status with ` + "`task update`" + `
3. On commit: Task ID in branch auto-links
4. On PR: ` + "`awesome-clickup-cli git link-pr`" + ` to connect PR to task
`

	if cfg.TaskID != "" {
		content += fmt.Sprintf(`
## Current Task

Task ID: %s
View: https://app.clickup.com/t/%s

To get task details:
`+"```bash"+`
awesome-clickup-cli task get %s --agent
`+"```"+`
`, cfg.TaskID, cfg.TaskID, cfg.TaskID)
	}

	return content, nil
}

// GenerateCodexMCPConfig creates MCP configuration for Codex CLI
func GenerateCodexMCPConfig(binaryPath string) string {
	absPath, _ := filepath.Abs(binaryPath)
	return fmt.Sprintf(`# Add to your codex MCP configuration
# Run: codex mcp add clickup --command "%s mcp-server"

mcp:
  servers:
    clickup:
      command: "%s"
      args: ["mcp-server"]
      description: "ClickUp task management - view tasks, update status, search"
`, absPath, absPath)
}

// WriteAgentsMD writes AGENTS.md to the specified directory
func WriteAgentsMD(dir string, cfg CodexConfig) error {
	content, err := GenerateAgentsMD(cfg)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, "AGENTS.md")
	return os.WriteFile(path, []byte(content), 0644)
}

// DetectCodexInstall checks if Codex CLI is installed
func DetectCodexInstall() (string, bool) {
	paths := []string{
		"codex",
		filepath.Join(os.Getenv("HOME"), ".cargo", "bin", "codex"),
		"/usr/local/bin/codex",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
		// Check in PATH
		if p == "codex" {
			if path, err := LookPath("codex"); err == nil {
				return path, true
			}
		}
	}
	return "", false
}

// LookPath searches for an executable in PATH
func LookPath(file string) (string, error) {
	pathEnv := os.Getenv("PATH")
	for _, dir := range strings.Split(pathEnv, string(os.PathListSeparator)) {
		path := filepath.Join(dir, file)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("executable not found: %s", file)
}
