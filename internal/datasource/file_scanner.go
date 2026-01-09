package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/manutej/maat-terminal/internal/graph"
)

// FileScanner scans a directory for source code files.
type FileScanner struct {
	rootPath   string
	projectID  string
	maxFiles   int
	extensions []string
}

// NewFileScanner creates a new file system scanner
func NewFileScanner(rootPath, projectID string) *FileScanner {
	return &FileScanner{
		rootPath:  rootPath,
		projectID: projectID,
		maxFiles:  200, // Limit for performance
		extensions: []string{
			".go", ".js", ".ts", ".tsx", ".jsx",
			".py", ".rb", ".rs", ".java", ".kt",
			".c", ".cpp", ".h", ".hpp",
			".md", ".yaml", ".yml", ".json", ".toml",
			".html", ".css", ".scss",
		},
	}
}

// SetMaxFiles sets the maximum number of files to scan
func (f *FileScanner) SetMaxFiles(n int) {
	f.maxFiles = n
}

// Name returns the data source identifier
func (f *FileScanner) Name() string {
	return "files:" + filepath.Base(f.rootPath)
}

// SupportsRefresh returns true
func (f *FileScanner) SupportsRefresh() bool {
	return true
}

// Load scans the directory and returns file nodes
func (f *FileScanner) Load(ctx context.Context) ([]graph.Node, []graph.Edge, error) {
	var nodes []graph.Node
	var edges []graph.Edge

	// Track directories for parent_of relationships
	dirs := make(map[string]string) // dir path -> node ID
	fileCount := 0

	err := filepath.Walk(f.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors, continue walking
		}

		// Skip hidden directories and common ignore patterns
		if info.IsDir() {
			base := filepath.Base(path)
			if strings.HasPrefix(base, ".") || f.shouldSkipDir(base) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check file limit
		if fileCount >= f.maxFiles {
			return filepath.SkipAll
		}

		// Check extension
		ext := strings.ToLower(filepath.Ext(path))
		if !f.isValidExtension(ext) {
			return nil
		}

		fileCount++

		// Create file node
		relPath, _ := filepath.Rel(f.rootPath, path)
		node, edge := f.createFileNode(relPath, path, info)
		nodes = append(nodes, node)
		edges = append(edges, edge)

		// Track parent directory
		dir := filepath.Dir(relPath)
		if dir != "." && dir != "" {
			if _, exists := dirs[dir]; !exists {
				dirNode, dirEdge := f.createDirNode(dir)
				nodes = append(nodes, dirNode)
				edges = append(edges, dirEdge)
				dirs[dir] = dirNode.ID
			}
			// File belongs to directory
			edges = append(edges, graph.Edge{
				ID:       fmt.Sprintf("edge:dir-file:%s", sanitizeID(relPath)),
				FromID:   dirs[dir],
				ToID:     node.ID,
				Relation: graph.EdgeOwns,
				Metadata: graph.EdgeMetadata{CreatedAt: time.Now()},
			})
		}

		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("walk failed: %w", err)
	}

	return nodes, edges, nil
}

// shouldSkipDir returns true for directories that should be ignored
func (f *FileScanner) shouldSkipDir(name string) bool {
	skipDirs := []string{
		"node_modules", "vendor", "dist", "build", "target",
		"__pycache__", ".git", ".svn", ".hg",
		"coverage", ".next", ".nuxt", ".cache",
	}
	for _, skip := range skipDirs {
		if name == skip {
			return true
		}
	}
	return false
}

// isValidExtension checks if file extension should be included
func (f *FileScanner) isValidExtension(ext string) bool {
	for _, valid := range f.extensions {
		if ext == valid {
			return true
		}
	}
	return false
}

// createFileNode creates a graph node for a file
func (f *FileScanner) createFileNode(relPath, fullPath string, info os.FileInfo) (graph.Node, graph.Edge) {
	// Detect language from extension
	lang := detectLanguage(filepath.Ext(relPath))

	// Count lines (simple approach - count newlines)
	lines := 0
	if content, err := os.ReadFile(fullPath); err == nil {
		lines = strings.Count(string(content), "\n") + 1
	}

	data := map[string]interface{}{
		"path":     relPath,
		"language": lang,
		"lines":    lines,
		"size":     info.Size(),
	}
	dataJSON, _ := json.Marshal(data)

	nodeID := fmt.Sprintf("file:%s", sanitizeID(relPath))

	node := graph.Node{
		ID:     nodeID,
		Type:   graph.NodeTypeFile,
		Source: "filesystem",
		Data:   dataJSON,
		Metadata: graph.NodeMetadata{
			CreatedAt:   info.ModTime(),
			UpdatedAt:   info.ModTime(),
			CreatedBy:   "file-scanner",
			AccessLevel: graph.RoleIC,
			SyncedAt:    time.Now(),
		},
	}

	// Edge: project owns file
	edge := graph.Edge{
		ID:       fmt.Sprintf("edge:project-file:%s", sanitizeID(relPath)),
		FromID:   f.projectID,
		ToID:     nodeID,
		Relation: graph.EdgeOwns,
		Metadata: graph.EdgeMetadata{CreatedAt: info.ModTime()},
	}

	return node, edge
}

// createDirNode creates a service node for a directory
func (f *FileScanner) createDirNode(dir string) (graph.Node, graph.Edge) {
	data := map[string]interface{}{
		"name": filepath.Base(dir),
		"path": dir,
		"type": "directory",
	}
	dataJSON, _ := json.Marshal(data)

	nodeID := fmt.Sprintf("service:dir:%s", sanitizeID(dir))

	node := graph.Node{
		ID:     nodeID,
		Type:   graph.NodeTypeService,
		Source: "filesystem",
		Data:   dataJSON,
		Metadata: graph.NodeMetadata{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			CreatedBy:   "file-scanner",
			AccessLevel: graph.RoleIC,
			SyncedAt:    time.Now(),
		},
	}

	// Edge: project owns directory
	edge := graph.Edge{
		ID:       fmt.Sprintf("edge:project-dir:%s", sanitizeID(dir)),
		FromID:   f.projectID,
		ToID:     nodeID,
		Relation: graph.EdgeOwns,
		Metadata: graph.EdgeMetadata{CreatedAt: time.Now()},
	}

	return node, edge
}

// detectLanguage returns the programming language for a file extension
func detectLanguage(ext string) string {
	languages := map[string]string{
		".go":   "Go",
		".js":   "JavaScript",
		".ts":   "TypeScript",
		".tsx":  "TypeScript",
		".jsx":  "JavaScript",
		".py":   "Python",
		".rb":   "Ruby",
		".rs":   "Rust",
		".java": "Java",
		".kt":   "Kotlin",
		".c":    "C",
		".cpp":  "C++",
		".h":    "C",
		".hpp":  "C++",
		".md":   "Markdown",
		".yaml": "YAML",
		".yml":  "YAML",
		".json": "JSON",
		".toml": "TOML",
		".html": "HTML",
		".css":  "CSS",
		".scss": "SCSS",
	}
	if lang, ok := languages[strings.ToLower(ext)]; ok {
		return lang
	}
	return "Unknown"
}
