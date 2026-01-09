package tui

import (
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
)

// Commands describe effects, runtime executes (Commandment #8: Async Purity)
// No goroutines - only tea.Cmd (Commandment #5: Controlled Effects)

// doNothing is a no-op command
func doNothing() tea.Msg {
	return nil
}

// fetchData loads mock graph data
// In Phase 2+, this will call Linear/GitHub APIs
func fetchData() tea.Cmd {
	return func() tea.Msg {
		// Load mock graph for testing
		nodes, edges := GetMockGraph()

		// Convert to display format
		displayNodes := make([]DisplayNode, len(nodes))
		for i, node := range nodes {
			displayNodes[i] = DisplayNode{
				ID:     node.ID,
				Type:   node.Type,
				Title:  node.Title(),
				Status: node.Status(),
			}
		}

		displayEdges := make([]DisplayEdge, len(edges))
		for i, edge := range edges {
			displayEdges[i] = DisplayEdge{
				FromID:   edge.FromID,
				ToID:     edge.ToID,
				Relation: edge.Relation,
			}
		}

		return GraphDataLoadedMsg{
			Nodes: displayNodes,
			Edges: displayEdges,
		}
	}
}

// executeConfirmedAction runs a user-confirmed external write
func executeConfirmedAction(action func() error) tea.Cmd {
	return func() tea.Msg {
		if err := action(); err != nil {
			return ErrorOccurred{Err: err}
		}
		return DataLoadedMsg{Data: "Action completed successfully"}
	}
}

// refreshData re-fetches current view's data
func refreshData() tea.Cmd {
	return func() tea.Msg {
		// Placeholder: Will re-query based on current view state
		return DataLoadedMsg{Data: "Data refreshed"}
	}
}

// openInBrowser opens a URL in the default browser (read-only action)
func openInBrowser(url string) tea.Cmd {
	return func() tea.Msg {
		if url == "" {
			return StatusMsg{Message: "No URL available for this node", IsError: true}
		}

		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", url)
		case "linux":
			cmd = exec.Command("xdg-open", url)
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", url)
		default:
			return StatusMsg{Message: "Unsupported platform for opening browser", IsError: true}
		}

		if err := cmd.Start(); err != nil {
			return StatusMsg{Message: "Failed to open browser: " + err.Error(), IsError: true}
		}

		return StatusMsg{Message: "Opened in browser", IsError: false}
	}
}

// copyToClipboard copies text to the system clipboard (read-only action)
func copyToClipboard(text string) tea.Cmd {
	return func() tea.Msg {
		if text == "" {
			return StatusMsg{Message: "No URL to copy", IsError: true}
		}

		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("pbcopy")
		case "linux":
			cmd = exec.Command("xclip", "-selection", "clipboard")
		case "windows":
			cmd = exec.Command("clip")
		default:
			return StatusMsg{Message: "Unsupported platform for clipboard", IsError: true}
		}

		stdin, err := cmd.StdinPipe()
		if err != nil {
			return StatusMsg{Message: "Clipboard error: " + err.Error(), IsError: true}
		}

		if err := cmd.Start(); err != nil {
			return StatusMsg{Message: "Clipboard error: " + err.Error(), IsError: true}
		}

		_, _ = stdin.Write([]byte(text))
		_ = stdin.Close()

		if err := cmd.Wait(); err != nil {
			return StatusMsg{Message: "Clipboard error: " + err.Error(), IsError: true}
		}

		return StatusMsg{Message: "URL copied to clipboard", IsError: false}
	}
}
