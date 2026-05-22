package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sidhartha1s/awesome-clickup-cli/internal/integrations"

	"github.com/spf13/cobra"
)

func newIntegrationsCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "integrations",
		Short: "Generate integrations for AI coding assistants",
		Long: `Generate configuration files and plugins for AI coding assistants.

Supported integrations:
  • codex      - OpenAI Codex CLI (AGENTS.md + MCP config)
  • hermes     - Hermes Agent (Python plugin)
  • openclaw   - OpenClaw Gateway (Python plugin)
  • claude     - Claude Code (CLAUDE.md + MCP config)
  • aider      - Aider (.aider.conf.yml)
  • mcp        - Generic MCP server setup
  • all        - Generate all integrations

Each integration enables the AI assistant to interact with ClickUp tasks,
update status, search, and use the unique git integration features.`,
	}

	cmd.AddCommand(newIntegrationsCodexCmd(flags))
	cmd.AddCommand(newIntegrationsHermesCmd(flags))
	cmd.AddCommand(newIntegrationsOpenClawCmd(flags))
	cmd.AddCommand(newIntegrationsClaudeCmd(flags))
	cmd.AddCommand(newIntegrationsAiderCmd(flags))
	cmd.AddCommand(newIntegrationsMCPCmd(flags))
	cmd.AddCommand(newIntegrationsAllCmd(flags))
	cmd.AddCommand(newIntegrationsDetectCmd(flags))

	return cmd
}

func newIntegrationsCodexCmd(flags *rootFlags) *cobra.Command {
	var outputDir string
	var taskID string

	cmd := &cobra.Command{
		Use:   "codex",
		Short: "Generate Codex CLI integration (AGENTS.md)",
		Long: `Generate AGENTS.md for OpenAI Codex CLI integration.

Creates an AGENTS.md file that instructs Codex how to use this CLI
for ClickUp task management during coding sessions.`,
		Example: `  # Generate in current directory
  awesome-clickup-cli integrations codex

  # Generate in specific directory
  awesome-clickup-cli integrations codex --output ./my-project

  # Include current task context
  awesome-clickup-cli integrations codex --task abc123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputDir == "" {
				outputDir = "."
			}

			cfg := integrations.CodexConfig{
				ProjectRoot: outputDir,
				TaskID:      taskID,
			}

			if err := integrations.WriteAgentsMD(outputDir, cfg); err != nil {
				return fmt.Errorf("writing AGENTS.md: %w", err)
			}

			if flags.asJSON {
				return printJSONFiltered(cmd.OutOrStdout(), map[string]interface{}{
					"integration": "codex",
					"file":        filepath.Join(outputDir, "AGENTS.md"),
					"success":     true,
				}, flags)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Created %s\n", filepath.Join(outputDir, "AGENTS.md"))
			fmt.Fprintln(cmd.OutOrStdout(), "")
			fmt.Fprintln(cmd.OutOrStdout(), "Codex will now see ClickUp commands when working in this directory.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory (default: current)")
	cmd.Flags().StringVar(&taskID, "task", "", "Include specific task context")

	return cmd
}

func newIntegrationsHermesCmd(flags *rootFlags) *cobra.Command {
	var outputDir string
	var pluginName string

	cmd := &cobra.Command{
		Use:   "hermes",
		Short: "Generate Hermes Agent plugin",
		Long: `Generate a Python plugin for Hermes Agent integration.

Creates a plugin that registers ClickUp tools with Hermes,
allowing the agent to manage tasks, search, and use git integration.`,
		Example: `  # Generate plugin in ~/.hermes/plugins/
  awesome-clickup-cli integrations hermes

  # Generate in custom directory
  awesome-clickup-cli integrations hermes --output ./plugins`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputDir == "" {
				outputDir = filepath.Join(os.Getenv("HOME"), ".hermes", "plugins")
			}
			if pluginName == "" {
				pluginName = "clickup"
			}

			cfg := integrations.HermesPluginConfig{
				PluginName:  pluginName,
				Description: "ClickUp task management via awesome-clickup-cli",
				BinaryPath:  "awesome-clickup-cli",
			}

			if err := integrations.WriteHermesPlugin(outputDir, cfg); err != nil {
				return fmt.Errorf("writing Hermes plugin: %w", err)
			}

			pluginDir := filepath.Join(outputDir, pluginName)

			if flags.asJSON {
				return printJSONFiltered(cmd.OutOrStdout(), map[string]interface{}{
					"integration": "hermes",
					"plugin_dir":  pluginDir,
					"files":       []string{"manifest.yaml", "__init__.py", "tools.py", "cli.py"},
					"success":     true,
				}, flags)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Created Hermes plugin at %s\n", pluginDir)
			fmt.Fprintln(cmd.OutOrStdout(), "")
			fmt.Fprintln(cmd.OutOrStdout(), "To enable:")
			fmt.Fprintln(cmd.OutOrStdout(), "  hermes plugins enable clickup")
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory (default: ~/.hermes/plugins)")
	cmd.Flags().StringVar(&pluginName, "name", "clickup", "Plugin name")

	return cmd
}

func newIntegrationsOpenClawCmd(flags *rootFlags) *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "openclaw",
		Short: "Generate OpenClaw plugin",
		Long: `Generate a plugin for OpenClaw Gateway integration.

Creates a plugin that provides ClickUp tools to OpenClaw agents
across all connected messaging platforms.`,
		Example: `  # Generate plugin in ~/.openclaw/plugins/
  awesome-clickup-cli integrations openclaw

  # Generate in custom directory
  awesome-clickup-cli integrations openclaw --output ./plugins`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputDir == "" {
				outputDir = filepath.Join(os.Getenv("HOME"), ".openclaw", "plugins")
			}

			cfg := integrations.OpenClawConfig{
				AgentName:   "clickup",
				Description: "ClickUp task management via awesome-clickup-cli",
				BinaryPath:  "awesome-clickup-cli",
			}

			if err := integrations.WriteOpenClawPlugin(outputDir, cfg); err != nil {
				return fmt.Errorf("writing OpenClaw plugin: %w", err)
			}

			pluginDir := filepath.Join(outputDir, "clickup")

			if flags.asJSON {
				return printJSONFiltered(cmd.OutOrStdout(), map[string]interface{}{
					"integration": "openclaw",
					"plugin_dir":  pluginDir,
					"files":       []string{"manifest.yaml", "tools.py"},
					"success":     true,
				}, flags)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Created OpenClaw plugin at %s\n", pluginDir)
			fmt.Fprintln(cmd.OutOrStdout(), "")
			fmt.Fprintln(cmd.OutOrStdout(), "To enable:")
			fmt.Fprintln(cmd.OutOrStdout(), "  openclaw plugins enable clickup")
			fmt.Fprintln(cmd.OutOrStdout(), "  openclaw gateway restart")
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory (default: ~/.openclaw/plugins)")

	return cmd
}

func newIntegrationsClaudeCmd(flags *rootFlags) *cobra.Command {
	var outputDir string
	var mcpOnly bool

	cmd := &cobra.Command{
		Use:   "claude",
		Short: "Generate Claude Code integration (CLAUDE.md + MCP)",
		Long: `Generate CLAUDE.md and MCP configuration for Claude Code integration.

Creates files that enable Claude Code to use this CLI for ClickUp
task management, including the MCP server configuration.`,
		Example: `  # Generate in current directory
  awesome-clickup-cli integrations claude

  # Only show MCP config (don't write files)
  awesome-clickup-cli integrations claude --mcp-only`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputDir == "" {
				outputDir = "."
			}

			cfg := integrations.ClaudeConfig{
				BinaryPath:  "awesome-clickup-cli",
				ProjectRoot: outputDir,
			}

			if mcpOnly {
				mcpConfig := integrations.GenerateMCPConfig("awesome-clickup-cli")
				if flags.asJSON {
					var parsed interface{}
					json.Unmarshal([]byte(mcpConfig), &parsed)
					return printJSONFiltered(cmd.OutOrStdout(), parsed, flags)
				}
				fmt.Fprintln(cmd.OutOrStdout(), "Add to ~/.claude.json:")
				fmt.Fprintln(cmd.OutOrStdout(), mcpConfig)
				return nil
			}

			if err := integrations.WriteClaudeMD(outputDir, cfg); err != nil {
				return fmt.Errorf("writing CLAUDE.md: %w", err)
			}

			if flags.asJSON {
				return printJSONFiltered(cmd.OutOrStdout(), map[string]interface{}{
					"integration": "claude",
					"file":        filepath.Join(outputDir, "CLAUDE.md"),
					"mcp_config":  integrations.GenerateMCPConfig("awesome-clickup-cli"),
					"success":     true,
				}, flags)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Created %s\n", filepath.Join(outputDir, "CLAUDE.md"))
			fmt.Fprintln(cmd.OutOrStdout(), "")
			fmt.Fprintln(cmd.OutOrStdout(), "To enable MCP server, add to ~/.claude.json:")
			fmt.Fprintln(cmd.OutOrStdout(), integrations.GenerateMCPConfig("awesome-clickup-cli"))
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory (default: current)")
	cmd.Flags().BoolVar(&mcpOnly, "mcp-only", false, "Only print MCP config (don't write files)")

	return cmd
}

func newIntegrationsAiderCmd(flags *rootFlags) *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "aider",
		Short: "Generate Aider configuration",
		Long: `Generate .aider.conf.yml for Aider integration.

Creates a configuration file that includes ClickUp context
in Aider coding sessions.`,
		Example: `  # Generate in current directory
  awesome-clickup-cli integrations aider`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputDir == "" {
				outputDir = "."
			}

			cfg := integrations.AiderConfig{
				BinaryPath:  "awesome-clickup-cli",
				ProjectRoot: outputDir,
			}

			if err := integrations.WriteAiderConfig(outputDir, cfg); err != nil {
				return fmt.Errorf("writing .aider.conf.yml: %w", err)
			}

			if flags.asJSON {
				return printJSONFiltered(cmd.OutOrStdout(), map[string]interface{}{
					"integration": "aider",
					"file":        filepath.Join(outputDir, ".aider.conf.yml"),
					"success":     true,
				}, flags)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ Created %s\n", filepath.Join(outputDir, ".aider.conf.yml"))
			fmt.Fprintln(cmd.OutOrStdout(), "")
			fmt.Fprintln(cmd.OutOrStdout(), "Aider will now include ClickUp context. Use /run to execute CLI commands.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory (default: current)")

	return cmd
}

func newIntegrationsMCPCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp-server",
		Short: "Start as MCP server",
		Long: `Start the CLI as an MCP server for AI assistant integration.

This allows any MCP-compatible client (Claude, Codex, etc.) to use
ClickUp tools through the Model Context Protocol.

The server runs on stdio and provides these tools:
  • clickup_task_get      - Get task details
  • clickup_task_update   - Update task status/priority
  • clickup_search        - Search tasks
  • clickup_comment_add   - Add comments
  • clickup_git_status    - Detect task from branch
  • clickup_stale_tasks   - Find stale tasks
  • clickup_workload      - Team workload analysis`,
		Example: `  # Start MCP server
  awesome-clickup-cli integrations mcp-server

  # Add to Claude's ~/.claude.json:
  # {
  #   "mcpServers": {
  #     "clickup": {
  #       "command": "awesome-clickup-cli",
  #       "args": ["integrations", "mcp-server"]
  #     }
  #   }
  # }`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSimpleMCPServer()
		},
	}

	return cmd
}

func runSimpleMCPServer() error {
	scanner := bufio.NewScanner(os.Stdin)
	binaryPath := os.Args[0]

	tools := []map[string]interface{}{
		{"name": "clickup_task_get", "description": "Get task details", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"task_id": map[string]interface{}{"type": "string"}}, "required": []string{"task_id"}}},
		{"name": "clickup_task_update", "description": "Update task status", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"task_id": map[string]interface{}{"type": "string"}, "status": map[string]interface{}{"type": "string"}}, "required": []string{"task_id"}}},
		{"name": "clickup_search", "description": "Search tasks", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"query": map[string]interface{}{"type": "string"}}, "required": []string{"query"}}},
		{"name": "clickup_git_status", "description": "Detect task from branch", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{}}},
		{"name": "clickup_stale_tasks", "description": "Find stale tasks", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"days": map[string]interface{}{"type": "integer", "default": 7}}}},
		{"name": "clickup_workload", "description": "Team workload", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{}}},
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req map[string]interface{}
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			sendMCPError(nil, -32700, "Parse error")
			continue
		}

		method, _ := req["method"].(string)
		id := req["id"]

		switch method {
		case "initialize":
			sendMCPResponse(id, map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]interface{}{"tools": map[string]interface{}{}},
				"serverInfo":      map[string]interface{}{"name": "awesome-clickup-cli", "version": "1.0.0"},
			})
		case "tools/list":
			sendMCPResponse(id, map[string]interface{}{"tools": tools})
		case "tools/call":
			params, _ := req["params"].(map[string]interface{})
			toolName, _ := params["name"].(string)
			toolArgs, _ := params["arguments"].(map[string]interface{})
			result := callMCPTool(binaryPath, toolName, toolArgs)
			sendMCPResponse(id, map[string]interface{}{
				"content": []map[string]interface{}{{"type": "text", "text": result}},
			})
		default:
			sendMCPError(id, -32601, "Method not found")
		}
	}
	return scanner.Err()
}

func sendMCPResponse(id interface{}, result interface{}) {
	resp := map[string]interface{}{"jsonrpc": "2.0", "id": id, "result": result}
	b, _ := json.Marshal(resp)
	fmt.Println(string(b))
}

func sendMCPError(id interface{}, code int, message string) {
	resp := map[string]interface{}{"jsonrpc": "2.0", "id": id, "error": map[string]interface{}{"code": code, "message": message}}
	b, _ := json.Marshal(resp)
	fmt.Println(string(b))
}

func callMCPTool(binaryPath, name string, args map[string]interface{}) string {
	var cmdArgs []string
	switch name {
	case "clickup_task_get":
		taskID, _ := args["task_id"].(string)
		cmdArgs = []string{"task", "get", taskID, "--agent"}
	case "clickup_task_update":
		taskID, _ := args["task_id"].(string)
		cmdArgs = []string{"task", "update", taskID}
		if status, ok := args["status"].(string); ok && status != "" {
			cmdArgs = append(cmdArgs, "--status", status)
		}
		cmdArgs = append(cmdArgs, "--agent")
	case "clickup_search":
		query, _ := args["query"].(string)
		cmdArgs = []string{"search", query, "--agent"}
	case "clickup_git_status":
		cmdArgs = []string{"git", "status", "--json"}
	case "clickup_stale_tasks":
		cmdArgs = []string{"stale", "--agent"}
		if days, ok := args["days"].(float64); ok {
			cmdArgs = append(cmdArgs, "--days", fmt.Sprintf("%d", int(days)))
		}
	case "clickup_workload":
		cmdArgs = []string{"load", "--agent"}
	default:
		return fmt.Sprintf(`{"error": "unknown tool: %s"}`, name)
	}

	cmd := exec.Command(binaryPath, cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, strings.ReplaceAll(err.Error(), `"`, `\"`))
	}
	return strings.TrimSpace(string(output))
}

func newIntegrationsAllCmd(flags *rootFlags) *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "all",
		Short: "Generate all integrations",
		Long: `Generate integration files for all supported AI assistants.

Creates AGENTS.md, CLAUDE.md, .aider.conf.yml, and plugin directories
for Hermes and OpenClaw.`,
		Example: `  # Generate all in current directory
  awesome-clickup-cli integrations all

  # Generate in specific directory
  awesome-clickup-cli integrations all --output ./my-project`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputDir == "" {
				outputDir = "."
			}

			results := make(map[string]interface{})

			// Codex (AGENTS.md)
			if err := integrations.WriteAgentsMD(outputDir, integrations.CodexConfig{ProjectRoot: outputDir}); err != nil {
				results["codex"] = map[string]interface{}{"success": false, "error": err.Error()}
			} else {
				results["codex"] = map[string]interface{}{"success": true, "file": filepath.Join(outputDir, "AGENTS.md")}
				if !flags.asJSON {
					fmt.Fprintf(cmd.OutOrStdout(), "✓ Codex: %s\n", filepath.Join(outputDir, "AGENTS.md"))
				}
			}

			// Claude (CLAUDE.md)
			if err := integrations.WriteClaudeMD(outputDir, integrations.ClaudeConfig{BinaryPath: "awesome-clickup-cli"}); err != nil {
				results["claude"] = map[string]interface{}{"success": false, "error": err.Error()}
			} else {
				results["claude"] = map[string]interface{}{"success": true, "file": filepath.Join(outputDir, "CLAUDE.md")}
				if !flags.asJSON {
					fmt.Fprintf(cmd.OutOrStdout(), "✓ Claude: %s\n", filepath.Join(outputDir, "CLAUDE.md"))
				}
			}

			// Aider
			if err := integrations.WriteAiderConfig(outputDir, integrations.AiderConfig{BinaryPath: "awesome-clickup-cli"}); err != nil {
				results["aider"] = map[string]interface{}{"success": false, "error": err.Error()}
			} else {
				results["aider"] = map[string]interface{}{"success": true, "file": filepath.Join(outputDir, ".aider.conf.yml")}
				if !flags.asJSON {
					fmt.Fprintf(cmd.OutOrStdout(), "✓ Aider: %s\n", filepath.Join(outputDir, ".aider.conf.yml"))
				}
			}

			// Hermes
			hermesDir := filepath.Join(outputDir, ".hermes-plugin")
			if err := integrations.WriteHermesPlugin(hermesDir, integrations.HermesPluginConfig{
				PluginName: "clickup",
				BinaryPath: "awesome-clickup-cli",
			}); err != nil {
				results["hermes"] = map[string]interface{}{"success": false, "error": err.Error()}
			} else {
				results["hermes"] = map[string]interface{}{"success": true, "dir": filepath.Join(hermesDir, "clickup")}
				if !flags.asJSON {
					fmt.Fprintf(cmd.OutOrStdout(), "✓ Hermes: %s\n", filepath.Join(hermesDir, "clickup"))
				}
			}

			// OpenClaw
			openclawDir := filepath.Join(outputDir, ".openclaw-plugin")
			if err := integrations.WriteOpenClawPlugin(openclawDir, integrations.OpenClawConfig{
				BinaryPath: "awesome-clickup-cli",
			}); err != nil {
				results["openclaw"] = map[string]interface{}{"success": false, "error": err.Error()}
			} else {
				results["openclaw"] = map[string]interface{}{"success": true, "dir": filepath.Join(openclawDir, "clickup")}
				if !flags.asJSON {
					fmt.Fprintf(cmd.OutOrStdout(), "✓ OpenClaw: %s\n", filepath.Join(openclawDir, "clickup"))
				}
			}

			if flags.asJSON {
				return printJSONFiltered(cmd.OutOrStdout(), results, flags)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "")
			fmt.Fprintln(cmd.OutOrStdout(), "All integrations generated. See 'integrations <name> --help' for setup instructions.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory (default: current)")

	return cmd
}

func newIntegrationsDetectCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detect",
		Short: "Detect installed AI assistants",
		Long:  `Scan for installed AI coding assistants and show which integrations are available.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			detected := make(map[string]interface{})

			if path, ok := integrations.DetectCodexInstall(); ok {
				detected["codex"] = map[string]interface{}{"installed": true, "path": path}
			} else {
				detected["codex"] = map[string]interface{}{"installed": false}
			}

			if path, ok := integrations.DetectHermesInstall(); ok {
				detected["hermes"] = map[string]interface{}{"installed": true, "path": path}
			} else {
				detected["hermes"] = map[string]interface{}{"installed": false}
			}

			if path, ok := integrations.DetectOpenClawInstall(); ok {
				detected["openclaw"] = map[string]interface{}{"installed": true, "path": path}
			} else {
				detected["openclaw"] = map[string]interface{}{"installed": false}
			}

			if path, ok := integrations.DetectClaudeCodeInstall(); ok {
				detected["claude"] = map[string]interface{}{"installed": true, "path": path}
			} else {
				detected["claude"] = map[string]interface{}{"installed": false}
			}

			if path, ok := integrations.DetectAiderInstall(); ok {
				detected["aider"] = map[string]interface{}{"installed": true, "path": path}
			} else {
				detected["aider"] = map[string]interface{}{"installed": false}
			}

			if flags.asJSON {
				return printJSONFiltered(cmd.OutOrStdout(), detected, flags)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Detected AI Assistants:")
			fmt.Fprintln(cmd.OutOrStdout(), "")

			for name, info := range detected {
				infoMap := info.(map[string]interface{})
				if infoMap["installed"].(bool) {
					fmt.Fprintf(cmd.OutOrStdout(), "  ✓ %-10s %s\n", name, infoMap["path"])
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "  ✗ %-10s not found\n", name)
				}
			}

			fmt.Fprintln(cmd.OutOrStdout(), "")
			fmt.Fprintln(cmd.OutOrStdout(), "Generate integration: awesome-clickup-cli integrations <name>")
			return nil
		},
	}

	return cmd
}
