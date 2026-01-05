# Go-Git Integration

**Library**: `github.com/go-git/go-git/v5`
**Version**: v5.11+
**Purpose**: Pure Go git implementation for FR-010

---

## DELTA FORCE: Exact Code for MAAT

### 1. Repository Access

**MAAT Requirement**: FR-010 Git commit history
**Pattern**: Open repo, traverse commits, extract metadata

```go
// maat/internal/git/repository.go
package git

import (
    "context"
    "time"

    "github.com/go-git/go-git/v5"
    "github.com/go-git/go-git/v5/plumbing"
    "github.com/go-git/go-git/v5/plumbing/object"
)

type Repository struct {
    repo *git.Repository
    path string
}

// Open or discover repository
func OpenRepository(path string) (*Repository, error) {
    // Try to open at path
    repo, err := git.PlainOpen(path)
    if err != nil {
        // Try to discover (walk up to find .git)
        repo, err = git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
            DetectDotGit: true,
        })
        if err != nil {
            return nil, fmt.Errorf("no git repository found: %w", err)
        }
    }

    return &Repository{repo: repo, path: path}, nil
}

// Get repository root path
func (r *Repository) RootPath() (string, error) {
    wt, err := r.repo.Worktree()
    if err != nil {
        return "", err
    }
    return wt.Filesystem.Root(), nil
}
```

### 2. Commit History Fetching

**MAAT Requirement**: FR-010 Commit graph with snippets
**Pattern**: Iterate commits, extract data, limit results

```go
// maat/internal/git/commits.go
package git

import (
    "github.com/go-git/go-git/v5"
    "github.com/go-git/go-git/v5/plumbing/object"
    "github.com/go-git/go-git/v5/plumbing/storer"
)

type Commit struct {
    Hash       string
    ShortHash  string
    Message    string
    FirstLine  string
    Author     string
    AuthorEmail string
    Date       time.Time
    Parents    []string
    Stats      DiffStats
    IsHead     bool
    Branches   []string
    Tags       []string
}

type DiffStats struct {
    FilesChanged int
    Additions    int
    Deletions    int
    Files        []FileChange
}

type FileChange struct {
    Path       string
    Status     ChangeStatus  // Added, Modified, Deleted
    Additions  int
    Deletions  int
    Snippet    string  // First few lines of diff
}

type ChangeStatus string
const (
    StatusAdded    ChangeStatus = "A"
    StatusModified ChangeStatus = "M"
    StatusDeleted  ChangeStatus = "D"
    StatusRenamed  ChangeStatus = "R"
)

// Fetch commit history with limit
func (r *Repository) FetchCommits(ctx context.Context, limit int) ([]Commit, error) {
    // Get HEAD reference
    head, err := r.repo.Head()
    if err != nil {
        return nil, fmt.Errorf("failed to get HEAD: %w", err)
    }

    // Get commit iterator
    iter, err := r.repo.Log(&git.LogOptions{
        From:  head.Hash(),
        Order: git.LogOrderCommitterTime,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get log: %w", err)
    }
    defer iter.Close()

    // Collect commits
    var commits []Commit
    isFirst := true

    err = iter.ForEach(func(c *object.Commit) error {
        if len(commits) >= limit {
            return storer.ErrStop
        }

        commit := Commit{
            Hash:        c.Hash.String(),
            ShortHash:   c.Hash.String()[:7],
            Message:     c.Message,
            FirstLine:   firstLine(c.Message),
            Author:      c.Author.Name,
            AuthorEmail: c.Author.Email,
            Date:        c.Author.When,
            IsHead:      isFirst,
        }

        // Get parent hashes
        for _, parent := range c.ParentHashes {
            commit.Parents = append(commit.Parents, parent.String()[:7])
        }

        commits = append(commits, commit)
        isFirst = false
        return nil
    })

    if err != nil && err != storer.ErrStop {
        return nil, err
    }

    // Fetch branch/tag info in parallel
    r.enrichWithRefs(commits)

    return commits, nil
}

// Fetch commits for a specific file
func (r *Repository) FetchFileHistory(ctx context.Context, path string, limit int) ([]Commit, error) {
    head, err := r.repo.Head()
    if err != nil {
        return nil, err
    }

    iter, err := r.repo.Log(&git.LogOptions{
        From:     head.Hash(),
        FileName: &path,
    })
    if err != nil {
        return nil, err
    }
    defer iter.Close()

    var commits []Commit
    err = iter.ForEach(func(c *object.Commit) error {
        if len(commits) >= limit {
            return storer.ErrStop
        }
        commits = append(commits, commitToStruct(c, false))
        return nil
    })

    return commits, err
}
```

### 3. Diff and Stats Extraction

**MAAT Requirement**: FR-010 Code diff snippets
**Pattern**: Compare commits, extract changes, truncate for preview

```go
// maat/internal/git/diff.go
package git

import (
    "strings"

    "github.com/go-git/go-git/v5/plumbing/object"
)

// Get diff stats for a commit
func (r *Repository) GetCommitStats(ctx context.Context, hash string) (DiffStats, error) {
    commit, err := r.repo.CommitObject(plumbing.NewHash(hash))
    if err != nil {
        return DiffStats{}, err
    }

    stats, err := commit.Stats()
    if err != nil {
        return DiffStats{}, err
    }

    result := DiffStats{
        FilesChanged: len(stats),
    }

    for _, stat := range stats {
        result.Additions += stat.Addition
        result.Deletions += stat.Deletion

        fc := FileChange{
            Path:      stat.Name,
            Additions: stat.Addition,
            Deletions: stat.Deletion,
            Status:    inferStatus(stat),
        }
        result.Files = append(result.Files, fc)
    }

    return result, nil
}

// Get full diff for a commit
func (r *Repository) GetCommitDiff(ctx context.Context, hash string) (string, error) {
    commit, err := r.repo.CommitObject(plumbing.NewHash(hash))
    if err != nil {
        return "", err
    }

    // Get parent (for diff comparison)
    var parentTree *object.Tree
    if commit.NumParents() > 0 {
        parent, err := commit.Parent(0)
        if err == nil {
            parentTree, _ = parent.Tree()
        }
    }

    commitTree, err := commit.Tree()
    if err != nil {
        return "", err
    }

    // Generate patch
    changes, err := parentTree.Diff(commitTree)
    if err != nil {
        return "", err
    }

    patch, err := changes.Patch()
    if err != nil {
        return "", err
    }

    return patch.String(), nil
}

// Get diff snippet for preview (first N lines)
func (r *Repository) GetDiffSnippet(ctx context.Context, hash, filePath string, maxLines int) (string, error) {
    fullDiff, err := r.GetCommitDiff(ctx, hash)
    if err != nil {
        return "", err
    }

    // Find the section for this file
    lines := strings.Split(fullDiff, "\n")
    var snippet []string
    inFile := false
    lineCount := 0

    for _, line := range lines {
        if strings.HasPrefix(line, "diff --git") {
            if strings.Contains(line, filePath) {
                inFile = true
            } else {
                inFile = false
            }
            continue
        }

        if inFile && (strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-")) {
            if !strings.HasPrefix(line, "+++") && !strings.HasPrefix(line, "---") {
                snippet = append(snippet, line)
                lineCount++
                if lineCount >= maxLines {
                    break
                }
            }
        }
    }

    return strings.Join(snippet, "\n"), nil
}

func inferStatus(stat object.FileStat) ChangeStatus {
    if stat.Addition > 0 && stat.Deletion == 0 {
        return StatusAdded
    }
    if stat.Addition == 0 && stat.Deletion > 0 {
        return StatusDeleted
    }
    return StatusModified
}
```

### 4. Branch and Tag Resolution

**MAAT Requirement**: FR-010 Branch visualization
**Pattern**: Map refs to commits

```go
// maat/internal/git/refs.go
package git

import (
    "github.com/go-git/go-git/v5/plumbing"
)

// Enrich commits with branch/tag info
func (r *Repository) enrichWithRefs(commits []Commit) {
    // Build hash -> commit index map
    hashIndex := make(map[string]int)
    for i, c := range commits {
        hashIndex[c.Hash] = i
    }

    // Get branches
    branches, _ := r.repo.Branches()
    branches.ForEach(func(ref *plumbing.Reference) error {
        hash := ref.Hash().String()
        if idx, ok := hashIndex[hash]; ok {
            commits[idx].Branches = append(commits[idx].Branches, ref.Name().Short())
        }
        return nil
    })

    // Get tags
    tags, _ := r.repo.Tags()
    tags.ForEach(func(ref *plumbing.Reference) error {
        hash := ref.Hash().String()
        if idx, ok := hashIndex[hash]; ok {
            commits[idx].Tags = append(commits[idx].Tags, ref.Name().Short())
        }
        return nil
    })
}

// Get current branch name
func (r *Repository) CurrentBranch() (string, error) {
    head, err := r.repo.Head()
    if err != nil {
        return "", err
    }

    if head.Name().IsBranch() {
        return head.Name().Short(), nil
    }

    // Detached HEAD - return short hash
    return head.Hash().String()[:7], nil
}

// List all branches
func (r *Repository) ListBranches() ([]string, error) {
    var branches []string

    iter, err := r.repo.Branches()
    if err != nil {
        return nil, err
    }

    iter.ForEach(func(ref *plumbing.Reference) error {
        branches = append(branches, ref.Name().Short())
        return nil
    })

    return branches, nil
}
```

### 5. Issue/PR Link Parsing

**MAAT Requirement**: FR-010 Link commits to issues
**Pattern**: Regex parse commit messages

```go
// maat/internal/git/links.go
package git

import (
    "regexp"
    "strings"
)

var (
    // Match patterns: LIN-123, #456, fixes #789, closes LINEAR-101
    linearPattern = regexp.MustCompile(`(?i)LIN-(\d+)`)
    githubPattern = regexp.MustCompile(`(?i)(?:fixes|closes|resolves)?\s*#(\d+)`)
    genericPattern = regexp.MustCompile(`(?i)([A-Z]+-\d+)`)
)

type LinkedRef struct {
    Type   string  // "linear", "github", "generic"
    ID     string  // "LIN-123", "#456"
    Action string  // "fixes", "closes", "mentions"
}

// Parse issue/PR references from commit message
func ParseLinkedRefs(message string) []LinkedRef {
    var refs []LinkedRef

    // Linear issues (LIN-123)
    for _, match := range linearPattern.FindAllStringSubmatch(message, -1) {
        refs = append(refs, LinkedRef{
            Type:   "linear",
            ID:     "LIN-" + match[1],
            Action: inferAction(message, match[0]),
        })
    }

    // GitHub issues (#123)
    for _, match := range githubPattern.FindAllStringSubmatch(message, -1) {
        action := "mentions"
        if strings.Contains(strings.ToLower(match[0]), "fixes") ||
           strings.Contains(strings.ToLower(match[0]), "closes") {
            action = "closes"
        }
        refs = append(refs, LinkedRef{
            Type:   "github",
            ID:     "#" + match[1],
            Action: action,
        })
    }

    return refs
}

func inferAction(message, ref string) string {
    lower := strings.ToLower(message)
    refLower := strings.ToLower(ref)

    // Find position of ref and look backwards for action keywords
    pos := strings.Index(lower, refLower)
    if pos > 0 {
        prefix := lower[:pos]
        if strings.HasSuffix(prefix, "fixes ") || strings.HasSuffix(prefix, "fix ") {
            return "fixes"
        }
        if strings.HasSuffix(prefix, "closes ") || strings.HasSuffix(prefix, "close ") {
            return "closes"
        }
        if strings.HasSuffix(prefix, "resolves ") || strings.HasSuffix(prefix, "resolve ") {
            return "resolves"
        }
    }
    return "mentions"
}
```

---

## Integration with Knowledge Graph

```go
// maat/internal/git/adapter.go
package git

import "maat/internal/graph"

// Convert commits to graph nodes
func (c Commit) ToNode() graph.Node {
    return graph.Node{
        ID:     "git:" + c.ShortHash,
        Type:   graph.NodeTypeCommit,
        Source: "git",
        Data:   mustMarshal(c),
        Metadata: graph.NodeMetadata{
            CreatedAt:   c.Date,
            AccessLevel: graph.RoleIC,
        },
    }
}

// Generate edges from commit relationships
func (c Commit) ToEdges() []graph.Edge {
    var edges []graph.Edge

    // Parent relationships
    for _, parent := range c.Parents {
        edges = append(edges, graph.Edge{
            FromID:   "git:" + c.ShortHash,
            ToID:     "git:" + parent,
            Relation: graph.EdgeParentOf,
        })
    }

    // Linked issue relationships
    for _, ref := range ParseLinkedRefs(c.Message) {
        var targetID string
        switch ref.Type {
        case "linear":
            targetID = "linear:" + ref.ID
        case "github":
            targetID = "github:" + ref.ID
        }

        if targetID != "" {
            edges = append(edges, graph.Edge{
                FromID:   "git:" + c.ShortHash,
                ToID:     targetID,
                Relation: graph.EdgeImplements,
            })
        }
    }

    return edges
}
```

---

## File Structure

```
maat/internal/git/
├── repository.go  # Repo access
├── commits.go     # Commit fetching
├── diff.go        # Diff extraction
├── refs.go        # Branch/tag resolution
├── links.go       # Issue link parsing
└── adapter.go     # Graph node conversion
```

---

## Performance Considerations

1. **Limit iteration**: Always use `storer.ErrStop` to break early
2. **Lazy diff loading**: Only fetch diffs when user expands commit
3. **Cache refs**: Branch/tag resolution is expensive, cache results
4. **Parallel enrichment**: Run stats fetching in goroutines via tea.Cmd
