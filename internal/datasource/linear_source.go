package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/manutej/maat-terminal/internal/graph"
)

// LinearSource fetches issues and projects from Linear API.
// Following Commandment #7 (Composition): Thin API client only.
type LinearSource struct {
	apiKey string
	teamID string
	client *http.Client
}

// NewLinearSource creates a Linear data source
// API key is read from LINEAR_API_KEY environment variable
func NewLinearSource(teamID string) *LinearSource {
	return &LinearSource{
		apiKey: os.Getenv("LINEAR_API_KEY"),
		teamID: teamID,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Name returns the data source identifier
func (l *LinearSource) Name() string {
	return "linear"
}

// SupportsRefresh returns true - Linear can be refreshed
func (l *LinearSource) SupportsRefresh() bool {
	return true
}

// Load fetches issues and projects from Linear
func (l *LinearSource) Load(ctx context.Context) ([]graph.Node, []graph.Edge, error) {
	if l.apiKey == "" {
		return nil, nil, fmt.Errorf("LINEAR_API_KEY environment variable not set")
	}

	var nodes []graph.Node
	var edges []graph.Edge

	// Fetch issues
	issues, err := l.fetchIssues(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("fetching issues: %w", err)
	}

	// Convert issues to nodes and collect edges
	for _, issue := range issues {
		node, issueEdges := l.issueToNode(issue)
		nodes = append(nodes, node)
		edges = append(edges, issueEdges...)
	}

	// Fetch projects
	projects, err := l.fetchProjects(ctx)
	if err != nil {
		// Log but continue - issues are more important
		fmt.Fprintf(os.Stderr, "Warning: failed to fetch projects: %v\n", err)
	} else {
		for _, project := range projects {
			node := l.projectToNode(project)
			nodes = append(nodes, node)
		}
	}

	return nodes, edges, nil
}

// LinearIssue represents the issue data from Linear API
type LinearIssue struct {
	ID          string   `json:"id"`
	Identifier  string   `json:"identifier"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    int      `json:"priority"`
	Status      string   `json:"status"`
	Labels      []string `json:"labels"`
	ProjectID   string   `json:"projectId"`
	ProjectName string   `json:"project"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
	URL         string   `json:"url"`
	// Relations
	BlockedBy []string `json:"blockedBy,omitempty"`
	Blocks    []string `json:"blocks,omitempty"`
	Related   []string `json:"relatedTo,omitempty"`
}

// LinearProject represents the project data from Linear API
type LinearProject struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	URL         string `json:"url"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// fetchIssues fetches issues from Linear GraphQL API
func (l *LinearSource) fetchIssues(ctx context.Context) ([]LinearIssue, error) {
	// Simplified query to stay under Linear's 10000 complexity limit
	// Removed: relations (high complexity), reduced first to 50
	query := `
	query IssuesByTeam($teamId: String!) {
		team(id: $teamId) {
			issues(first: 50) {
				nodes {
					id
					identifier
					title
					priority
					state { name }
					labels { nodes { name } }
					project { id name }
					createdAt
					updatedAt
					url
				}
			}
		}
	}`

	variables := map[string]interface{}{
		"teamId": l.teamID,
	}

	resp, err := l.graphqlRequest(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	// Parse response (simplified - no description or relations to stay under complexity limit)
	var result struct {
		Data struct {
			Team struct {
				Issues struct {
					Nodes []struct {
						ID         string `json:"id"`
						Identifier string `json:"identifier"`
						Title      string `json:"title"`
						Priority   int    `json:"priority"`
						State      struct {
							Name string `json:"name"`
						} `json:"state"`
						Labels struct {
							Nodes []struct {
								Name string `json:"name"`
							} `json:"nodes"`
						} `json:"labels"`
						Project *struct {
							ID   string `json:"id"`
							Name string `json:"name"`
						} `json:"project"`
						CreatedAt string `json:"createdAt"`
						UpdatedAt string `json:"updatedAt"`
						URL       string `json:"url"`
					} `json:"nodes"`
				} `json:"issues"`
			} `json:"team"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("Linear API error: %s", result.Errors[0].Message)
	}

	// Convert to LinearIssue slice
	var issues []LinearIssue
	for _, n := range result.Data.Team.Issues.Nodes {
		issue := LinearIssue{
			ID:         n.ID,
			Identifier: n.Identifier,
			Title:      n.Title,
			Priority:   n.Priority,
			Status:     n.State.Name,
			CreatedAt:  n.CreatedAt,
			UpdatedAt:  n.UpdatedAt,
			URL:        n.URL,
		}

		// Extract labels
		for _, label := range n.Labels.Nodes {
			issue.Labels = append(issue.Labels, label.Name)
		}

		// Extract project
		if n.Project != nil {
			issue.ProjectID = n.Project.ID
			issue.ProjectName = n.Project.Name
		}

		// Note: Relations fetched separately if needed to avoid query complexity limits

		issues = append(issues, issue)
	}

	return issues, nil
}

// fetchProjects fetches projects from Linear GraphQL API
func (l *LinearSource) fetchProjects(ctx context.Context) ([]LinearProject, error) {
	query := `
	query ProjectsByTeam($teamId: String!) {
		team(id: $teamId) {
			projects(first: 50) {
				nodes {
					id
					name
					description
					state
					url
					createdAt
					updatedAt
				}
			}
		}
	}`

	variables := map[string]interface{}{
		"teamId": l.teamID,
	}

	resp, err := l.graphqlRequest(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Team struct {
				Projects struct {
					Nodes []LinearProject `json:"nodes"`
				} `json:"projects"`
			} `json:"team"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return result.Data.Team.Projects.Nodes, nil
}

// graphqlRequest makes a GraphQL request to Linear API
func (l *LinearSource) graphqlRequest(ctx context.Context, query string, variables map[string]interface{}) ([]byte, error) {
	body := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.linear.app/graphql", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", l.apiKey)

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Linear API returned %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// issueToNode converts a Linear issue to a graph node and edges
func (l *LinearSource) issueToNode(issue LinearIssue) (graph.Node, []graph.Edge) {
	// Build node data
	data := map[string]interface{}{
		"identifier":  issue.Identifier,
		"title":       issue.Title,
		"description": issue.Description,
		"priority":    issue.Priority,
		"status":      issue.Status,
		"labels":      issue.Labels,
		"project":     issue.ProjectName,
		"url":         issue.URL,
	}
	dataJSON, _ := json.Marshal(data)

	// Parse timestamps
	createdAt, _ := time.Parse(time.RFC3339, issue.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, issue.UpdatedAt)

	node := graph.Node{
		ID:     fmt.Sprintf("linear:%s", issue.Identifier),
		Type:   graph.NodeTypeIssue,
		Source: "linear",
		Data:   dataJSON,
		Metadata: graph.NodeMetadata{
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			AccessLevel: graph.RoleIC, // All issues visible to ICs
			SyncedAt:    time.Now(),
		},
	}

	// Build edges from relations
	var edges []graph.Edge

	// Blocks edges
	for _, blockedID := range issue.Blocks {
		edges = append(edges, graph.Edge{
			ID:       fmt.Sprintf("edge:%s-blocks-%s", issue.Identifier, blockedID),
			FromID:   node.ID,
			ToID:     fmt.Sprintf("linear:%s", blockedID),
			Relation: graph.EdgeBlocks,
		})
	}

	// Related edges
	for _, relatedID := range issue.Related {
		edges = append(edges, graph.Edge{
			ID:       fmt.Sprintf("edge:%s-related-%s", issue.Identifier, relatedID),
			FromID:   node.ID,
			ToID:     fmt.Sprintf("linear:%s", relatedID),
			Relation: graph.EdgeRelated,
		})
	}

	// Parent (project) edge
	if issue.ProjectID != "" {
		edges = append(edges, graph.Edge{
			ID:       fmt.Sprintf("edge:%s-in-project-%s", issue.Identifier, issue.ProjectID),
			FromID:   fmt.Sprintf("linear:project:%s", issue.ProjectID),
			ToID:     node.ID,
			Relation: graph.EdgeOwns,
		})
	}

	return node, edges
}

// projectToNode converts a Linear project to a graph node
func (l *LinearSource) projectToNode(project LinearProject) graph.Node {
	data := map[string]interface{}{
		"name":        project.Name,
		"description": project.Description,
		"status":      project.Status,
		"url":         project.URL,
	}
	dataJSON, _ := json.Marshal(data)

	createdAt, _ := time.Parse(time.RFC3339, project.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, project.UpdatedAt)

	return graph.Node{
		ID:     fmt.Sprintf("linear:project:%s", project.ID),
		Type:   graph.NodeTypeProject,
		Source: "linear",
		Data:   dataJSON,
		Metadata: graph.NodeMetadata{
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			AccessLevel: graph.RoleLead, // Projects visible to leads+
			SyncedAt:    time.Now(),
		},
	}
}
