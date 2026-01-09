package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/manutej/maat-terminal/internal/graph"
)

// GitScanner scans a local git repository for commits and branches.
// Uses git CLI for simplicity and broad compatibility.
type GitScanner struct {
	repoPath string
	maxCommits int
}

// NewGitScanner creates a new git repository scanner
func NewGitScanner(repoPath string) *GitScanner {
	return &GitScanner{
		repoPath:   repoPath,
		maxCommits: 50, // Limit to recent commits for performance
	}
}

// SetMaxCommits sets the maximum number of commits to load
func (g *GitScanner) SetMaxCommits(n int) {
	g.maxCommits = n
}

// Name returns the data source identifier
func (g *GitScanner) Name() string {
	return "git:" + filepath.Base(g.repoPath)
}

// SupportsRefresh returns true - git repos can always be refreshed
func (g *GitScanner) SupportsRefresh() bool {
	return true
}

// Load scans the git repository and returns nodes and edges
func (g *GitScanner) Load(ctx context.Context) ([]graph.Node, []graph.Edge, error) {
	var nodes []graph.Node
	var edges []graph.Edge

	// Check if directory is a git repo
	if !g.isGitRepo() {
		return nil, nil, fmt.Errorf("not a git repository: %s", g.repoPath)
	}

	// Create project node
	projectNode := g.createProjectNode()
	nodes = append(nodes, projectNode)

	// Load commits
	commits, commitEdges, err := g.loadCommits(projectNode.ID)
	if err == nil {
		nodes = append(nodes, commits...)
		edges = append(edges, commitEdges...)
	}

	// Load branches as service nodes
	branches, branchEdges, err := g.loadBranches(projectNode.ID)
	if err == nil {
		nodes = append(nodes, branches...)
		edges = append(edges, branchEdges...)
	}

	return nodes, edges, nil
}

// isGitRepo checks if the path is a git repository
func (g *GitScanner) isGitRepo() bool {
	cmd := exec.Command("git", "-C", g.repoPath, "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// createProjectNode creates a project node from the repo
func (g *GitScanner) createProjectNode() graph.Node {
	repoName := filepath.Base(g.repoPath)

	// Get remote URL if available
	remoteURL := ""
	cmd := exec.Command("git", "-C", g.repoPath, "remote", "get-url", "origin")
	if output, err := cmd.Output(); err == nil {
		remoteURL = strings.TrimSpace(string(output))
	}

	data := map[string]interface{}{
		"name":        repoName,
		"description": fmt.Sprintf("Git repository at %s", g.repoPath),
		"status":      "active",
		"remote":      remoteURL,
		"path":        g.repoPath,
	}
	dataJSON, _ := json.Marshal(data)

	return graph.Node{
		ID:     fmt.Sprintf("project:%s", repoName),
		Type:   graph.NodeTypeProject,
		Source: "git",
		Data:   dataJSON,
		Metadata: graph.NodeMetadata{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			CreatedBy:   "git-scanner",
			AccessLevel: graph.RoleExec,
			SyncedAt:    time.Now(),
		},
	}
}

// loadCommits loads recent commits from the repository
func (g *GitScanner) loadCommits(projectID string) ([]graph.Node, []graph.Edge, error) {
	var nodes []graph.Node
	var edges []graph.Edge

	// Get commit log in a parseable format
	// Format: hash|author|date|subject
	cmd := exec.Command("git", "-C", g.repoPath, "log",
		fmt.Sprintf("--max-count=%d", g.maxCommits),
		"--format=%H|%an|%aI|%s",
	)
	output, err := cmd.Output()
	if err != nil {
		return nil, nil, fmt.Errorf("git log failed: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var prevCommitID string

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 4)
		if len(parts) < 4 {
			continue
		}

		hash := parts[0]
		author := parts[1]
		dateStr := parts[2]
		message := parts[3]

		commitID := fmt.Sprintf("commit:%s", hash[:8])

		commitDate, _ := time.Parse(time.RFC3339, dateStr)

		data := map[string]interface{}{
			"message": message,
			"author":  author,
			"hash":    hash,
			"date":    dateStr,
		}
		dataJSON, _ := json.Marshal(data)

		node := graph.Node{
			ID:     commitID,
			Type:   graph.NodeTypeCommit,
			Source: "git",
			Data:   dataJSON,
			Metadata: graph.NodeMetadata{
				CreatedAt:   commitDate,
				UpdatedAt:   commitDate,
				CreatedBy:   author,
				AccessLevel: graph.RoleIC,
				SyncedAt:    time.Now(),
			},
		}
		nodes = append(nodes, node)

		// Edge: project owns commit
		edges = append(edges, graph.Edge{
			ID:       fmt.Sprintf("edge:project-commit:%s", hash[:8]),
			FromID:   projectID,
			ToID:     commitID,
			Relation: graph.EdgeOwns,
			Metadata: graph.EdgeMetadata{CreatedAt: commitDate},
		})

		// Edge: commit parent relationship (sequential)
		if prevCommitID != "" {
			edges = append(edges, graph.Edge{
				ID:       fmt.Sprintf("edge:commit-parent:%s-%s", hash[:8], prevCommitID[7:]),
				FromID:   prevCommitID,
				ToID:     commitID,
				Relation: graph.EdgeParentOf,
				Metadata: graph.EdgeMetadata{CreatedAt: commitDate},
			})
		}
		prevCommitID = commitID

		// Check for issue references in commit message (e.g., #123, fixes #456)
		issueRefs := extractIssueReferences(message)
		for _, issueNum := range issueRefs {
			edges = append(edges, graph.Edge{
				ID:       fmt.Sprintf("edge:commit-mentions:%s-%d", hash[:8], issueNum),
				FromID:   commitID,
				ToID:     fmt.Sprintf("issue:%d", issueNum),
				Relation: graph.EdgeMentions,
				Metadata: graph.EdgeMetadata{CreatedAt: commitDate},
			})
		}
	}

	return nodes, edges, nil
}

// loadBranches loads git branches as service nodes
func (g *GitScanner) loadBranches(projectID string) ([]graph.Node, []graph.Edge, error) {
	var nodes []graph.Node
	var edges []graph.Edge

	// Get all branches
	cmd := exec.Command("git", "-C", g.repoPath, "branch", "-a", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, nil, fmt.Errorf("git branch failed: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, branch := range lines {
		branch = strings.TrimSpace(branch)
		if branch == "" || strings.Contains(branch, "HEAD") {
			continue
		}

		branchID := fmt.Sprintf("service:branch:%s", sanitizeID(branch))

		data := map[string]interface{}{
			"name": branch,
			"type": "branch",
		}
		dataJSON, _ := json.Marshal(data)

		node := graph.Node{
			ID:     branchID,
			Type:   graph.NodeTypeService,
			Source: "git",
			Data:   dataJSON,
			Metadata: graph.NodeMetadata{
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				CreatedBy:   "git-scanner",
				AccessLevel: graph.RoleIC,
				SyncedAt:    time.Now(),
			},
		}
		nodes = append(nodes, node)

		// Edge: project owns branch
		edges = append(edges, graph.Edge{
			ID:       fmt.Sprintf("edge:project-branch:%s", sanitizeID(branch)),
			FromID:   projectID,
			ToID:     branchID,
			Relation: graph.EdgeOwns,
			Metadata: graph.EdgeMetadata{CreatedAt: time.Now()},
		})
	}

	return nodes, edges, nil
}

// extractIssueReferences finds issue numbers in commit messages
func extractIssueReferences(message string) []int {
	var refs []int
	// Simple regex-free approach: look for #N patterns
	parts := strings.Fields(message)
	for _, part := range parts {
		// Remove punctuation
		part = strings.Trim(part, ".,;:!?()[]")
		if strings.HasPrefix(part, "#") {
			numStr := strings.TrimPrefix(part, "#")
			var num int
			if _, err := fmt.Sscanf(numStr, "%d", &num); err == nil && num > 0 {
				refs = append(refs, num)
			}
		}
	}
	return refs
}

// sanitizeID makes a string safe for use as an ID
func sanitizeID(s string) string {
	// Replace problematic characters
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
