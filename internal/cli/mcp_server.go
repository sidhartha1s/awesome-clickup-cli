package cli

import (
	"github.com/spf13/cobra"
)

func newMCPServerCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp-server",
		Short: "Start as MCP server for AI assistant integration",
		Long: `Start the CLI as an MCP (Model Context Protocol) server.

This allows any MCP-compatible AI assistant to use ClickUp tools:
  • Claude Code / Claude Desktop
  • OpenAI Codex CLI
  • Cursor
  • Continue
  • Any MCP client

The server communicates over stdio using JSON-RPC 2.0.

Available tools:
  • clickup_task_get      - Get task details
  • clickup_task_update   - Update task status/priority
  • clickup_search        - Search tasks by keyword
  • clickup_comment_add   - Add comments to tasks
  • clickup_git_status    - Detect task from git branch
  • clickup_stale_tasks   - Find stale tasks
  • clickup_workload      - Team workload analysis`,
		Example: `  # Start MCP server (for debugging)
  awesome-clickup-cli mcp-server

  # Add to Claude Code ~/.claude.json:
  {
    "mcpServers": {
      "clickup": {
        "command": "awesome-clickup-cli",
        "args": ["mcp-server"]
      }
    }
  }

  # Add to Codex:
  codex mcp add clickup --command "awesome-clickup-cli mcp-server"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSimpleMCPServer()
		},
	}

	return cmd
}
