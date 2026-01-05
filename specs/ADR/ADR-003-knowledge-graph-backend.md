# ADR-003: Knowledge Graph Backend

**Status**: Accepted
**Date**: 2026-01-05
**Deciders**: MAAT Architecture Team
**Commandments**: #1 Immutable Truth, #3 Text Interface

---

## Context

The TUI graph visualization (ADR-002) needs a **persistent data layer**:
- Cache fetched data for offline access
- Store relationships across systems (Linear ↔ GitHub ↔ Files)
- Enable semantic queries ("show all issues blocking this PR")
- Support AI context assembly

Options considered:
1. **Neo4j**: Full graph database (heavy, external dependency)
2. **SQLite with JSON**: Simple but loses graph semantics
3. **SQLite with graph views**: Local-first, graph-capable, single binary

## Decision

Adopt **SQLite with graph views** for the knowledge graph backend:

```sql
-- Core schema
CREATE TABLE nodes (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,           -- Issue | PR | File | Commit | Service
    source TEXT NOT NULL,         -- linear | github | local | plugin:<name>
    data JSON NOT NULL,           -- Source-specific payload
    metadata JSON NOT NULL,       -- created_at, updated_at, access_level
    embedding BLOB                -- Optional: vector for semantic search
);

CREATE TABLE edges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_id TEXT NOT NULL REFERENCES nodes(id),
    to_id TEXT NOT NULL REFERENCES nodes(id),
    relation TEXT NOT NULL,       -- blocks | related | implements | calls | owns
    metadata JSON,
    UNIQUE(from_id, to_id, relation)
);

-- Indexes for graph traversal
CREATE INDEX idx_nodes_type ON nodes(type);
CREATE INDEX idx_nodes_source ON nodes(source);
CREATE INDEX idx_edges_from ON edges(from_id);
CREATE INDEX idx_edges_to ON edges(to_id);
CREATE INDEX idx_edges_relation ON edges(relation);

-- Graph views for common queries
CREATE VIEW issue_dependencies AS
SELECT
    n1.id as issue_id,
    n1.data->>'title' as issue_title,
    n2.id as blocks_id,
    n2.data->>'title' as blocks_title
FROM nodes n1
JOIN edges e ON n1.id = e.from_id AND e.relation = 'blocks'
JOIN nodes n2 ON e.to_id = n2.id
WHERE n1.type = 'Issue';

CREATE VIEW pr_file_map AS
SELECT
    n1.id as pr_id,
    n1.data->>'number' as pr_number,
    n2.id as file_id,
    n2.data->>'path' as file_path
FROM nodes n1
JOIN edges e ON n1.id = e.from_id AND e.relation = 'modifies'
JOIN nodes n2 ON e.to_id = n2.id
WHERE n1.type = 'PR' AND n2.type = 'File';
```

### Go Data Structures

```go
type KnowledgeGraph struct {
    db      *sql.DB
    cache   *NodeCache  // In-memory LRU for hot nodes
}

type Node struct {
    ID       string          `json:"id"`
    Type     NodeType        `json:"type"`
    Source   string          `json:"source"`
    Data     json.RawMessage `json:"data"`
    Metadata NodeMetadata    `json:"metadata"`
}

type NodeMetadata struct {
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    CreatedBy   string    `json:"created_by"`   // user | ai:<session_id>
    AccessLevel Role      `json:"access_level"` // exec | lead | ic
    SyncedAt    time.Time `json:"synced_at"`    // Last API sync
}

type Edge struct {
    FromID   string       `json:"from_id"`
    ToID     string       `json:"to_id"`
    Relation EdgeType     `json:"relation"`
    Metadata EdgeMetadata `json:"metadata,omitempty"`
}

type EdgeType string
const (
    EdgeBlocks     EdgeType = "blocks"
    EdgeRelated    EdgeType = "related"
    EdgeImplements EdgeType = "implements"
    EdgeCalls      EdgeType = "calls"
    EdgeOwns       EdgeType = "owns"
    EdgeModifies   EdgeType = "modifies"
    EdgeMentions   EdgeType = "mentions"
)
```

### Query Patterns

```go
// Traverse graph from a node
func (kg *KnowledgeGraph) Neighbors(nodeID string, relations []EdgeType) ([]Node, error) {
    query := `
        SELECT n.* FROM nodes n
        JOIN edges e ON (e.to_id = n.id OR e.from_id = n.id)
        WHERE (e.from_id = ? OR e.to_id = ?)
        AND e.relation IN (?)
        AND n.id != ?
    `
    // Execute and return nodes
}

// Find path between nodes (BFS)
func (kg *KnowledgeGraph) FindPath(fromID, toID string) ([]Edge, error) {
    // BFS implementation using edges table
}

// Semantic query (if embeddings enabled)
func (kg *KnowledgeGraph) SemanticSearch(query string, limit int) ([]Node, error) {
    // Vector similarity search using embedding column
}
```

### Sync Strategy

```go
type SyncManager struct {
    kg       *KnowledgeGraph
    plugins  []DataSourcePlugin
    interval time.Duration  // Default: 5 minutes
}

func (sm *SyncManager) Sync(ctx context.Context) error {
    for _, plugin := range sm.plugins {
        nodes, err := plugin.FetchNodes(ctx, NodeFilter{
            UpdatedAfter: sm.lastSync,
        })
        if err != nil {
            continue  // Log and continue with other plugins
        }

        for _, node := range nodes {
            sm.kg.UpsertNode(node)  // Insert or update
        }

        edges, err := plugin.FetchEdges(ctx, nodeIDs(nodes))
        for _, edge := range edges {
            sm.kg.UpsertEdge(edge)
        }
    }
    sm.lastSync = time.Now()
    return nil
}
```

## Consequences

### Positive
- **Local-First**: Works offline; fast queries
- **Single Binary**: SQLite embedded; no external database
- **Graph Capable**: Views and recursive CTEs enable graph queries
- **Extensible**: JSON columns allow schema evolution

### Negative
- **Not True Graph DB**: Complex graph algorithms (PageRank) harder
- **Embedding Size**: Vector search adds storage overhead
- **Sync Complexity**: Must handle conflicts and staleness

### Mitigations
- Use recursive CTEs for multi-hop traversals
- Embeddings optional (enable via config)
- Sync manager with TTL and conflict resolution

## Compliance

This ADR enforces:
- **Commandment #1**: Immutable data storage (append-mostly pattern)
- **Commandment #3**: JSON as text interface for node data

## References

- [SQLite as Application File Format](https://www.sqlite.org/appfileformat.html)
- [Graph Queries in SQLite](https://www.sqlite.org/lang_with.html)
- MAAT Constitution: `/Users/manu/Documents/LUXOR/MAAT/docs/CONSTITUTION.md`
