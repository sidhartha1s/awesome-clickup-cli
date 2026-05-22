package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Server implements an MCP server for ClickUp integration
type Server struct {
	binaryPath string
}

// NewServer creates a new MCP server instance
func NewServer(binaryPath string) *Server {
	if binaryPath == "" {
		binaryPath = os.Args[0]
	}
	return &Server{binaryPath: binaryPath}
}

// Tool represents an MCP tool definition
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// Request represents an MCP JSON-RPC request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response represents an MCP JSON-RPC response
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error represents an MCP error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetTools returns the list of available MCP tools
func (s *Server) GetTools() []Tool {
	return []Tool{
		{
			Name:        "clickup_task_get",
			Description: "Get details for a ClickUp task including name, status, assignees, and due date",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"task_id": map[string]interface{}{
						"type":        "string",
						"description": "ClickUp task ID (e.g., 'abc123', 'CU-abc123')",
					},
				},
				"required": []string{"task_id"},
			},
		},
		{
			Name:        "clickup_task_update",
			Description: "Update a ClickUp task's status, priority, or assignee",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"task_id": map[string]interface{}{
						"type":        "string",
						"description": "ClickUp task ID",
					},
					"status": map[string]interface{}{
						"type":        "string",
						"description": "New status (e.g., 'in progress', 'complete')",
					},
					"priority": map[string]interface{}{
						"type":        "string",
						"description": "Priority (1=urgent, 2=high, 3=normal, 4=low)",
					},
				},
				"required": []string{"task_id"},
			},
		},
		{
			Name:        "clickup_search",
			Description: "Search ClickUp tasks by keyword",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum results (default 10)",
						"default":     10,
					},
				},
				"required": []string{"query"},
			},
		},
		{
			Name:        "clickup_comment_add",
			Description: "Add a comment to a ClickUp task",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"task_id": map[string]interface{}{
						"type":        "string",
						"description": "ClickUp task ID",
					},
					"text": map[string]interface{}{
						"type":        "string",
						"description": "Comment text (supports markdown)",
					},
				},
				"required": []string{"task_id", "text"},
			},
		},
		{
			Name:        "clickup_git_status",
			Description: "Detect ClickUp task from current git branch name",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "clickup_stale_tasks",
			Description: "Find tasks not updated within N days",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"days": map[string]interface{}{
						"type":        "integer",
						"description": "Days since last update (default 7)",
						"default":     7,
					},
				},
			},
		},
		{
			Name:        "clickup_workload",
			Description: "Get task distribution across team members",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "clickup_spaces_list",
			Description: "List all ClickUp spaces",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "clickup_lists_get",
			Description: "Get lists in a folder or space",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"folder_id": map[string]interface{}{
						"type":        "string",
						"description": "Folder ID to get lists from",
					},
				},
			},
		},
	}
}

// CallTool executes a tool and returns the result
func (s *Server) CallTool(name string, args map[string]interface{}) (interface{}, error) {
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
		if priority, ok := args["priority"].(string); ok && priority != "" {
			cmdArgs = append(cmdArgs, "--priority", priority)
		}
		cmdArgs = append(cmdArgs, "--agent")

	case "clickup_search":
		query, _ := args["query"].(string)
		cmdArgs = []string{"search", query, "--agent"}
		if limit, ok := args["limit"].(float64); ok {
			cmdArgs = append(cmdArgs, "--limit", fmt.Sprintf("%d", int(limit)))
		}

	case "clickup_comment_add":
		taskID, _ := args["task_id"].(string)
		text, _ := args["text"].(string)
		cmdArgs = []string{"comment", "create", "--task-id", taskID, "--text", text, "--agent"}

	case "clickup_git_status":
		cmdArgs = []string{"git", "status", "--json"}

	case "clickup_stale_tasks":
		cmdArgs = []string{"stale", "--agent"}
		if days, ok := args["days"].(float64); ok {
			cmdArgs = append(cmdArgs, "--days", fmt.Sprintf("%d", int(days)))
		}

	case "clickup_workload":
		cmdArgs = []string{"load", "--agent"}

	case "clickup_spaces_list":
		cmdArgs = []string{"space", "list", "--agent"}

	case "clickup_lists_get":
		cmdArgs = []string{"list", "list", "--agent"}
		if folderID, ok := args["folder_id"].(string); ok && folderID != "" {
			cmdArgs = append(cmdArgs, "--folder-id", folderID)
		}

	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}

	cmd := exec.Command(s.binaryPath, cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("command failed: %s", string(exitErr.Stderr))
		}
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return map[string]interface{}{"output": strings.TrimSpace(string(output))}, nil
	}
	return result, nil
}

// Run starts the MCP server (stdio transport)
func (s *Server) Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.sendError(encoder, nil, -32700, "Parse error")
			continue
		}

		s.handleRequest(encoder, &req)
	}

	return scanner.Err()
}

func (s *Server) handleRequest(encoder *json.Encoder, req *Request) {
	switch req.Method {
	case "initialize":
		encoder.Encode(Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]interface{}{
					"name":    "awesome-clickup-cli",
					"version": "1.0.0",
				},
			},
		})

	case "tools/list":
		encoder.Encode(Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"tools": s.GetTools(),
			},
		})

	case "tools/call":
		var params struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			s.sendError(encoder, req.ID, -32602, "Invalid params")
			return
		}

		result, err := s.CallTool(params.Name, params.Arguments)
		if err != nil {
			s.sendError(encoder, req.ID, -32000, err.Error())
			return
		}

		encoder.Encode(Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": mustMarshal(result),
					},
				},
			},
		})

	default:
		s.sendError(encoder, req.ID, -32601, "Method not found")
	}
}

func (s *Server) sendError(encoder *json.Encoder, id interface{}, code int, message string) {
	encoder.Encode(Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	})
}

func mustMarshal(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
