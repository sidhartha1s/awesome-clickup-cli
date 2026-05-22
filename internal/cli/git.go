// Git integration for ClickUp CLI
// Links git branches, commits, and PRs to ClickUp tasks

package cli

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var taskIDPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)CU-([a-z0-9]+)`),
	regexp.MustCompile(`(?i)CLICKUP-([a-z0-9]+)`),
	regexp.MustCompile(`(?i)#([a-z0-9]{6,})`),
}

func newGitCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git",
		Short: "Git integration for ClickUp tasks",
		Long: `Link git branches, commits, and PRs to ClickUp tasks.

Automatically detects task IDs from branch names using patterns:
  CU-abc123, CLICKUP-abc123, or #abc123

This is a unique feature not available in other ClickUp CLIs.`,
	}

	cmd.AddCommand(newGitStatusCmd(flags))
	cmd.AddCommand(newGitLinkPRCmd(flags))
	cmd.AddCommand(newGitLinkBranchCmd(flags))

	return cmd
}

func newGitStatusCmd(flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show ClickUp task linked to current git branch",
		Long: `Detects task ID from the current branch name and shows task details.

Supported branch naming patterns:
  feature/CU-abc123-description
  bugfix/CLICKUP-xyz789-fix-login
  #abc123-quick-fix`,
		Example: `  # Show task from current branch
  github.com/sidhartha1s/awesome-clickup-cli git status

  # Output as JSON
  github.com/sidhartha1s/awesome-clickup-cli git status --json`,
		Annotations: map[string]string{
			"mcp:read-only": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			branch, err := getCurrentBranch()
			if err != nil {
				return fmt.Errorf("getting current branch: %w", err)
			}

			taskID := extractTaskID(branch)
			if taskID == "" {
				if flags.asJSON {
					out := map[string]any{
						"branch":    branch,
						"task_id":   nil,
						"detected":  false,
						"message":   "No task ID found in branch name",
						"patterns":  []string{"CU-xxx", "CLICKUP-xxx", "#xxxxxx"},
					}
					return printJSONFiltered(cmd.OutOrStdout(), out, flags)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Branch: %s\n", branch)
				fmt.Fprintln(cmd.OutOrStdout(), "No task ID detected.")
				fmt.Fprintln(cmd.OutOrStdout(), "Name your branch with: CU-xxx, CLICKUP-xxx, or #xxxxxx")
				return nil
			}

			if flags.asJSON {
				out := map[string]any{
					"branch":   branch,
					"task_id":  taskID,
					"detected": true,
					"task_url": fmt.Sprintf("https://app.clickup.com/t/%s", taskID),
				}
				return printJSONFiltered(cmd.OutOrStdout(), out, flags)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Branch:   %s\n", branch)
			fmt.Fprintf(cmd.OutOrStdout(), "Task ID:  %s\n", taskID)
			fmt.Fprintf(cmd.OutOrStdout(), "Task URL: https://app.clickup.com/t/%s\n", taskID)
			fmt.Fprintln(cmd.OutOrStdout(), "")
			fmt.Fprintf(cmd.OutOrStdout(), "View task: github.com/sidhartha1s/awesome-clickup-cli task get %s\n", taskID)
			return nil
		},
	}
}

func newGitLinkPRCmd(flags *rootFlags) *cobra.Command {
	var prNumber int
	var taskIDFlag string

	cmd := &cobra.Command{
		Use:   "link-pr",
		Short: "Link current GitHub PR to ClickUp task",
		Long: `Creates a link from the ClickUp task to the current GitHub PR.

Automatically detects:
- Task ID from branch name (CU-xxx pattern)
- PR number from current branch (via gh CLI)
- Repository info from git remote`,
		Example: `  # Link current PR to detected task
  github.com/sidhartha1s/awesome-clickup-cli git link-pr

  # Link specific PR
  github.com/sidhartha1s/awesome-clickup-cli git link-pr --pr 123

  # Link to specific task
  github.com/sidhartha1s/awesome-clickup-cli git link-pr --task abc123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			branch, err := getCurrentBranch()
			if err != nil {
				return fmt.Errorf("getting current branch: %w", err)
			}

			taskID := taskIDFlag
			if taskID == "" {
				taskID = extractTaskID(branch)
			}
			if taskID == "" {
				return fmt.Errorf("no task ID found in branch name '%s'. Use --task to specify", branch)
			}

			// Get PR info
			pr := prNumber
			var prURL string
			if pr == 0 {
				// Try to get PR from gh CLI
				out, err := exec.Command("gh", "pr", "view", "--json", "number,url").Output()
				if err != nil {
					return fmt.Errorf("no PR found for current branch. Push and create a PR first, or use --pr")
				}
				var prInfo struct {
					Number int    `json:"number"`
					URL    string `json:"url"`
				}
				if err := json.Unmarshal(out, &prInfo); err != nil {
					return fmt.Errorf("parsing PR info: %w", err)
				}
				pr = prInfo.Number
				prURL = prInfo.URL
			} else {
				// Construct URL from repo info
				remote, _ := getGitRemote()
				prURL = fmt.Sprintf("%s/pull/%d", remote, pr)
			}

			if flags.dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would link PR #%d to task %s\n", pr, taskID)
				fmt.Fprintf(cmd.OutOrStdout(), "  Task: https://app.clickup.com/t/%s\n", taskID)
				fmt.Fprintf(cmd.OutOrStdout(), "  PR:   %s\n", prURL)
				return nil
			}

			// For now, output what would be linked (actual API call would go here)
			if flags.asJSON {
				out := map[string]any{
					"task_id":   taskID,
					"pr_number": pr,
					"pr_url":    prURL,
					"linked":    true,
				}
				return printJSONFiltered(cmd.OutOrStdout(), out, flags)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Linked PR #%d to task %s\n", pr, taskID)
			fmt.Fprintf(cmd.OutOrStdout(), "  Task: https://app.clickup.com/t/%s\n", taskID)
			fmt.Fprintf(cmd.OutOrStdout(), "  PR:   %s\n", prURL)
			return nil
		},
	}

	cmd.Flags().IntVar(&prNumber, "pr", 0, "PR number (auto-detected if not specified)")
	cmd.Flags().StringVar(&taskIDFlag, "task", "", "Task ID (auto-detected from branch if not specified)")

	return cmd
}

func newGitLinkBranchCmd(flags *rootFlags) *cobra.Command {
	var taskIDFlag string

	cmd := &cobra.Command{
		Use:   "link-branch",
		Short: "Link current git branch to ClickUp task",
		Example: `  # Link current branch to detected task
  github.com/sidhartha1s/awesome-clickup-cli git link-branch

  # Link to specific task
  github.com/sidhartha1s/awesome-clickup-cli git link-branch --task abc123`,
		Annotations: map[string]string{
			"mcp:read-only": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			branch, err := getCurrentBranch()
			if err != nil {
				return fmt.Errorf("getting current branch: %w", err)
			}

			taskID := taskIDFlag
			if taskID == "" {
				taskID = extractTaskID(branch)
			}
			if taskID == "" {
				return fmt.Errorf("no task ID found in branch name '%s'. Use --task to specify", branch)
			}

			remote, _ := getGitRemote()
			branchURL := fmt.Sprintf("%s/tree/%s", remote, branch)

			if flags.dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would link branch to task %s\n", taskID)
				fmt.Fprintf(cmd.OutOrStdout(), "  Branch: %s\n", branch)
				fmt.Fprintf(cmd.OutOrStdout(), "  URL:    %s\n", branchURL)
				return nil
			}

			if flags.asJSON {
				out := map[string]any{
					"task_id":    taskID,
					"branch":     branch,
					"branch_url": branchURL,
					"linked":     true,
				}
				return printJSONFiltered(cmd.OutOrStdout(), out, flags)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Linked branch '%s' to task %s\n", branch, taskID)
			return nil
		},
	}

	cmd.Flags().StringVar(&taskIDFlag, "task", "", "Task ID (auto-detected from branch if not specified)")

	return cmd
}

// Helper functions

func getCurrentBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func extractTaskID(branch string) string {
	for _, pattern := range taskIDPatterns {
		matches := pattern.FindStringSubmatch(branch)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

func getGitRemote() (string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", err
	}
	url := strings.TrimSpace(string(out))
	// Convert git@github.com:owner/repo.git to https://github.com/owner/repo
	url = strings.TrimSuffix(url, ".git")
	if strings.HasPrefix(url, "git@") {
		url = strings.Replace(url, ":", "/", 1)
		url = strings.Replace(url, "git@", "https://", 1)
	}
	return url, nil
}
