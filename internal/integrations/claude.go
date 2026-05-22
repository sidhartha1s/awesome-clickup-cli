package integrations

import (
	"fmt"
	"os"
	"path/filepath"
)

// ClaudeConfig holds configuration for Claude Code integration
type ClaudeConfig struct {
	BinaryPath  string
	ProjectRoot string
}

// GenerateClaudeMD creates a CLAUDE.md file for Claude Code integration
func GenerateClaudeMD(cfg ClaudeConfig) string {
	binaryPath := cfg.BinaryPath
	if binaryPath == "" {
		binaryPath = "awesome-clickup-cli"
	}

	return fmt.Sprintf(`# CLAUDE.md - ClickUp Integration

This project uses ClickUp for task management. The awesome-clickup-cli tool is available
for viewing and managing tasks directly from the command line.

## Quick Reference

`+"```bash"+`
# View task linked to current branch
%s git status

# Get task details
%s task get <task-id> --agent

# Update task status
%s task update <task-id> --status "in progress"

# Search tasks
%s search "keyword" --agent

# Find stale tasks
%s stale --days 7 --agent
`+"```"+`

## Branch Naming

Name branches with task IDs for automatic detection:
- `+"`feature/CU-abc123-description`"+`
- `+"`bugfix/CLICKUP-xyz789-fix`"+`
- `+"`#abc123-quick-fix`"+`

## Workflow

1. `+"`%s git status`"+` - Confirm which task you're on
2. Work on the code
3. `+"`%s task update <id> --status 'in progress'`"+` - Update status
4. `+"`%s git link-pr`"+` - Link PR to task after push

## MCP Server

This CLI can run as an MCP server for richer integration:

`+"```bash"+`
%s mcp-server
`+"```"+`

Add to your Claude MCP config (~/.claude.json):

`+"```json"+`
{
  "mcpServers": {
    "clickup": {
      "command": "%s",
      "args": ["mcp-server"]
    }
  }
}
`+"```"+`
`, binaryPath, binaryPath, binaryPath, binaryPath, binaryPath,
		binaryPath, binaryPath, binaryPath, binaryPath, binaryPath)
}

// GenerateMCPConfig creates MCP server configuration for Claude
func GenerateMCPConfig(binaryPath string) string {
	if binaryPath == "" {
		binaryPath = "awesome-clickup-cli"
	}

	return fmt.Sprintf(`{
  "mcpServers": {
    "clickup": {
      "command": "%s",
      "args": ["mcp-server"],
      "description": "ClickUp task management - view tasks, update status, search, git integration"
    }
  }
}`, binaryPath)
}

// WriteClaudeMD writes CLAUDE.md to the specified directory
func WriteClaudeMD(dir string, cfg ClaudeConfig) error {
	content := GenerateClaudeMD(cfg)
	path := filepath.Join(dir, "CLAUDE.md")
	return os.WriteFile(path, []byte(content), 0644)
}

// DetectClaudeCodeInstall checks if Claude Code is available
func DetectClaudeCodeInstall() (string, bool) {
	paths := []string{
		"claude",
		filepath.Join(os.Getenv("HOME"), ".claude", "local", "claude"),
		"/usr/local/bin/claude",
	}

	for _, p := range paths {
		if p == "claude" {
			if path, err := LookPath("claude"); err == nil {
				return path, true
			}
		} else if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	return "", false
}

// GetClaudeMCPConfigPath returns the path to Claude's MCP config
func GetClaudeMCPConfigPath() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".claude.json")
}
