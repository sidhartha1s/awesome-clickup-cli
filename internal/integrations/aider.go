package integrations

import (
	"fmt"
	"os"
	"path/filepath"
)

// AiderConfig holds configuration for Aider integration
type AiderConfig struct {
	BinaryPath  string
	ProjectRoot string
}

// GenerateAiderConfig creates .aider.conf.yml content for ClickUp integration
func GenerateAiderConfig(cfg AiderConfig) string {
	binaryPath := cfg.BinaryPath
	if binaryPath == "" {
		binaryPath = "awesome-clickup-cli"
	}

	return fmt.Sprintf(`# Aider configuration with ClickUp integration
# https://aider.chat/docs/config.html

# Auto-commit with task context
auto-commits: true

# Include ClickUp context in prompts
read:
  - CLAUDE.md
  - AGENTS.md

# Custom commands for ClickUp integration
# Use /run to execute these

# Example workflow:
# /run %s git status
# /run %s task get <task-id> --agent
# /run %s task update <task-id> --status "in progress"
`, binaryPath, binaryPath, binaryPath)
}

// GenerateAiderInstructions creates a README for Aider + ClickUp workflow
func GenerateAiderInstructions(cfg AiderConfig) string {
	binaryPath := cfg.BinaryPath
	if binaryPath == "" {
		binaryPath = "awesome-clickup-cli"
	}

	return fmt.Sprintf(`# Aider + ClickUp Workflow

## Setup

1. Install awesome-clickup-cli and authenticate:
   `+"```bash"+`
   %s auth set-token YOUR_API_TOKEN
   %s doctor
   `+"```"+`

2. Add the .aider.conf.yml to your project (already done if you ran the integration command)

## Workflow

### Start a session with task context

`+"```bash"+`
# Get current task from branch
%s git status

# Start aider with task context
aider --message "Working on ClickUp task $(` + "`%s git status --json | jq -r .task_id`" + `)"
`+"```"+`

### During development

Use Aider's /run command to interact with ClickUp:

`+"```"+`
/run %s task get abc123 --agent
/run %s task update abc123 --status "in progress"
/run %s search "related feature" --agent
`+"```"+`

### After completing work

`+"```bash"+`
# Update task status
%s task update <task-id> --status "complete"

# Link PR if created
%s git link-pr
`+"```"+`

## Tips

- Name branches with task IDs: `+"`feature/CU-abc123-description`"+`
- The CLI auto-detects task IDs from branch names
- Use `+"`--agent`"+` flag for JSON output that Aider can parse
`, binaryPath, binaryPath, binaryPath, binaryPath, binaryPath, binaryPath, binaryPath, binaryPath, binaryPath)
}

// WriteAiderConfig writes .aider.conf.yml to the specified directory
func WriteAiderConfig(dir string, cfg AiderConfig) error {
	content := GenerateAiderConfig(cfg)
	path := filepath.Join(dir, ".aider.conf.yml")
	return os.WriteFile(path, []byte(content), 0644)
}

// DetectAiderInstall checks if Aider is installed
func DetectAiderInstall() (string, bool) {
	paths := []string{
		"aider",
		filepath.Join(os.Getenv("HOME"), ".local", "bin", "aider"),
		"/usr/local/bin/aider",
	}

	for _, p := range paths {
		if p == "aider" {
			if path, err := LookPath("aider"); err == nil {
				return path, true
			}
		} else if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	return "", false
}
