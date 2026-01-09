package graph

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Store provides persistent storage for the knowledge graph using SQLite
type Store struct {
	db *sql.DB
}

// NewStore creates a new graph store at the specified database path
// If dbPath is ":memory:", an in-memory database is used
func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	store := &Store{db: db}

	if err := store.CreateTables(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return store, nil
}

// CreateTables initializes the database schema per ADR-003
func (s *Store) CreateTables() error {
	schema := `
	-- Core nodes table
	CREATE TABLE IF NOT EXISTS nodes (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		source TEXT NOT NULL,
		data JSON NOT NULL,
		metadata JSON NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Core edges table
	CREATE TABLE IF NOT EXISTS edges (
		id TEXT PRIMARY KEY,
		from_id TEXT NOT NULL,
		to_id TEXT NOT NULL,
		relation TEXT NOT NULL,
		metadata JSON,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (from_id) REFERENCES nodes(id) ON DELETE CASCADE,
		FOREIGN KEY (to_id) REFERENCES nodes(id) ON DELETE CASCADE,
		UNIQUE(from_id, to_id, relation)
	);

	-- Indexes for graph traversal performance
	CREATE INDEX IF NOT EXISTS idx_nodes_type ON nodes(type);
	CREATE INDEX IF NOT EXISTS idx_nodes_source ON nodes(source);
	CREATE INDEX IF NOT EXISTS idx_edges_from ON edges(from_id);
	CREATE INDEX IF NOT EXISTS idx_edges_to ON edges(to_id);
	CREATE INDEX IF NOT EXISTS idx_edges_relation ON edges(relation);

	-- Graph views for common queries
	CREATE VIEW IF NOT EXISTS issue_dependencies AS
	SELECT
		n1.id as issue_id,
		json_extract(n1.data, '$.title') as issue_title,
		n2.id as blocks_id,
		json_extract(n2.data, '$.title') as blocks_title
	FROM nodes n1
	JOIN edges e ON n1.id = e.from_id AND e.relation = 'blocks'
	JOIN nodes n2 ON e.to_id = n2.id
	WHERE n1.type = 'Issue';

	CREATE VIEW IF NOT EXISTS pr_file_map AS
	SELECT
		n1.id as pr_id,
		json_extract(n1.data, '$.number') as pr_number,
		n2.id as file_id,
		json_extract(n2.data, '$.path') as file_path
	FROM nodes n1
	JOIN edges e ON n1.id = e.from_id AND e.relation = 'modifies'
	JOIN nodes n2 ON e.to_id = n2.id
	WHERE n1.type = 'PR' AND n2.type = 'File';
	`

	_, err := s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

// AddNode inserts a new node into the graph
// Returns error if node with same ID already exists
func (s *Store) AddNode(node Node) error {
	// Set default metadata if not provided
	if node.Metadata.CreatedAt.IsZero() {
		node.Metadata.CreatedAt = time.Now()
	}
	if node.Metadata.UpdatedAt.IsZero() {
		node.Metadata.UpdatedAt = time.Now()
	}

	// Validate node type
	if !ValidateNodeType(string(node.Type)) {
		return fmt.Errorf("invalid node type: %s", node.Type)
	}

	// Marshal metadata to JSON
	metadataJSON, err := json.Marshal(node.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Insert node
	_, err = s.db.Exec(`
		INSERT INTO nodes (id, type, source, data, metadata)
		VALUES (?, ?, ?, ?, ?)
	`, node.ID, node.Type, node.Source, node.Data, metadataJSON)

	if err != nil {
		return fmt.Errorf("failed to insert node: %w", err)
	}

	return nil
}

// UpsertNode inserts or updates a node (idempotent operation)
func (s *Store) UpsertNode(node Node) error {
	// Update timestamp
	node.Metadata.UpdatedAt = time.Now()
	if node.Metadata.CreatedAt.IsZero() {
		node.Metadata.CreatedAt = time.Now()
	}

	// Validate node type
	if !ValidateNodeType(string(node.Type)) {
		return fmt.Errorf("invalid node type: %s", node.Type)
	}

	// Marshal metadata to JSON
	metadataJSON, err := json.Marshal(node.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Upsert node (SQLite 3.24.0+)
	_, err = s.db.Exec(`
		INSERT INTO nodes (id, type, source, data, metadata)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			type = excluded.type,
			source = excluded.source,
			data = excluded.data,
			metadata = excluded.metadata
	`, node.ID, node.Type, node.Source, node.Data, metadataJSON)

	if err != nil {
		return fmt.Errorf("failed to upsert node: %w", err)
	}

	return nil
}

// AddEdge inserts a new edge into the graph
// Returns error if edge with same (from_id, to_id, relation) already exists
func (s *Store) AddEdge(edge Edge) error {
	// Validate edge type
	if !ValidateEdgeType(string(edge.Relation)) {
		return fmt.Errorf("invalid edge relation: %s", edge.Relation)
	}

	// Generate ID if not provided
	if edge.ID == "" {
		edge.ID = fmt.Sprintf("%s-%s-%s", edge.FromID, edge.Relation, edge.ToID)
	}

	// Set default metadata
	if edge.Metadata.CreatedAt.IsZero() {
		edge.Metadata.CreatedAt = time.Now()
	}

	// Marshal metadata to JSON
	var metadataJSON []byte
	var err error
	if edge.Metadata.Data != nil || !edge.Metadata.CreatedAt.IsZero() {
		metadataJSON, err = json.Marshal(edge.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal edge metadata: %w", err)
		}
	}

	// Insert edge
	_, err = s.db.Exec(`
		INSERT INTO edges (id, from_id, to_id, relation, metadata)
		VALUES (?, ?, ?, ?, ?)
	`, edge.ID, edge.FromID, edge.ToID, edge.Relation, metadataJSON)

	if err != nil {
		return fmt.Errorf("failed to insert edge: %w", err)
	}

	return nil
}

// UpsertEdge inserts or updates an edge (idempotent operation)
func (s *Store) UpsertEdge(edge Edge) error {
	// Validate edge type
	if !ValidateEdgeType(string(edge.Relation)) {
		return fmt.Errorf("invalid edge relation: %s", edge.Relation)
	}

	// Generate ID if not provided
	if edge.ID == "" {
		edge.ID = fmt.Sprintf("%s-%s-%s", edge.FromID, edge.Relation, edge.ToID)
	}

	// Set default metadata
	if edge.Metadata.CreatedAt.IsZero() {
		edge.Metadata.CreatedAt = time.Now()
	}

	// Marshal metadata to JSON
	var metadataJSON []byte
	var err error
	if edge.Metadata.Data != nil || !edge.Metadata.CreatedAt.IsZero() {
		metadataJSON, err = json.Marshal(edge.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal edge metadata: %w", err)
		}
	}

	// Upsert edge
	_, err = s.db.Exec(`
		INSERT INTO edges (id, from_id, to_id, relation, metadata)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(from_id, to_id, relation) DO UPDATE SET
			metadata = excluded.metadata
	`, edge.ID, edge.FromID, edge.ToID, edge.Relation, metadataJSON)

	if err != nil {
		return fmt.Errorf("failed to upsert edge: %w", err)
	}

	return nil
}

// GetNode retrieves a node by ID
func (s *Store) GetNode(id string) (*Node, error) {
	var node Node
	var metadataJSON []byte

	err := s.db.QueryRow(`
		SELECT id, type, source, data, metadata
		FROM nodes
		WHERE id = ?
	`, id).Scan(&node.ID, &node.Type, &node.Source, &node.Data, &metadataJSON)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("node not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query node: %w", err)
	}

	// Unmarshal metadata
	if err := json.Unmarshal(metadataJSON, &node.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &node, nil
}

// GetNeighbors returns all nodes connected to the given node
// regardless of edge direction or relation type
func (s *Store) GetNeighbors(nodeID string) ([]Node, error) {
	rows, err := s.db.Query(`
		SELECT DISTINCT n.id, n.type, n.source, n.data, n.metadata
		FROM nodes n
		JOIN edges e ON (e.to_id = n.id OR e.from_id = n.id)
		WHERE (e.from_id = ? OR e.to_id = ?)
		AND n.id != ?
	`, nodeID, nodeID, nodeID)

	if err != nil {
		return nil, fmt.Errorf("failed to query neighbors: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var neighbors []Node
	for rows.Next() {
		var node Node
		var metadataJSON []byte

		err := rows.Scan(&node.ID, &node.Type, &node.Source, &node.Data, &metadataJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}

		// Unmarshal metadata
		if err := json.Unmarshal(metadataJSON, &node.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		neighbors = append(neighbors, node)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return neighbors, nil
}

// GetEdges returns all edges connected to a node (both incoming and outgoing)
func (s *Store) GetEdges(nodeID string) ([]Edge, error) {
	rows, err := s.db.Query(`
		SELECT id, from_id, to_id, relation, metadata
		FROM edges
		WHERE from_id = ? OR to_id = ?
	`, nodeID, nodeID)

	if err != nil {
		return nil, fmt.Errorf("failed to query edges: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var edges []Edge
	for rows.Next() {
		var edge Edge
		var metadataJSON sql.NullString

		err := rows.Scan(&edge.ID, &edge.FromID, &edge.ToID, &edge.Relation, &metadataJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan edge: %w", err)
		}

		// Unmarshal metadata if present
		if metadataJSON.Valid {
			if err := json.Unmarshal([]byte(metadataJSON.String), &edge.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal edge metadata: %w", err)
			}
		}

		edges = append(edges, edge)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating edge rows: %w", err)
	}

	return edges, nil
}

// DeleteNode removes a node and all connected edges (cascade delete)
func (s *Store) DeleteNode(id string) error {
	result, err := s.db.Exec("DELETE FROM nodes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("node not found: %s", id)
	}

	return nil
}

// DeleteEdge removes a specific edge by ID
func (s *Store) DeleteEdge(id string) error {
	result, err := s.db.Exec("DELETE FROM edges WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete edge: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("edge not found: %s", id)
	}

	return nil
}

// ListNodes returns all nodes, optionally filtered
func (s *Store) ListNodes(filter *NodeFilter) ([]Node, error) {
	query := "SELECT id, type, source, data, metadata FROM nodes WHERE 1=1"
	args := []interface{}{}

	if filter != nil {
		if len(filter.Types) > 0 {
			placeholders := ""
			for i, t := range filter.Types {
				if i > 0 {
					placeholders += ","
				}
				placeholders += "?"
				args = append(args, t)
			}
			query += " AND type IN (" + placeholders + ")"
		}

		if len(filter.Sources) > 0 {
			placeholders := ""
			for i, s := range filter.Sources {
				if i > 0 {
					placeholders += ","
				}
				placeholders += "?"
				args = append(args, s)
			}
			query += " AND source IN (" + placeholders + ")"
		}

		if !filter.UpdatedAfter.IsZero() {
			query += " AND json_extract(metadata, '$.updated_at') > ?"
			args = append(args, filter.UpdatedAfter.Format(time.RFC3339))
		}
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query nodes: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var nodes []Node
	for rows.Next() {
		var node Node
		var metadataJSON []byte

		err := rows.Scan(&node.ID, &node.Type, &node.Source, &node.Data, &metadataJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}

		if err := json.Unmarshal(metadataJSON, &node.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		nodes = append(nodes, node)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return nodes, nil
}

// Close closes the database connection
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
