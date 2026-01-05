# ADR-007: Plugin Architecture

**Status**: Accepted
**Date**: 2026-01-05
**Deciders**: MAAT Architecture Team
**Commandments**: #2 Single Responsibility, #5 Controlled Effects

---

## Context

Gap analysis identified the need for **extensibility**:
- Organizations use different tools (Jira, Azure DevOps, GitLab)
- Custom data sources (internal wikis, databases)
- Future integrations we can't predict
- Community contributions

MAAT must define a clear **platform boundary**:
- What is core (non-negotiable)?
- What is extensible (plugins)?
- How do plugins integrate safely?

## Decision

Implement **plugin architecture** with clear interfaces and lifecycle:

### Platform Boundary

```
┌──────────────────────────────────────────────────────┐
│                      MAAT CORE                        │
├──────────────────────────────────────────────────────┤
│  Knowledge Graph  │  TUI Shell  │  Command Palette   │
│  (SQLite + Views) │  (Bubble Tea)│  (Actions)        │
├──────────────────────────────────────────────────────┤
│                   PLUGIN INTERFACE                    │
├──────────────────────────────────────────────────────┤
│  DataSourcePlugin │  ViewPlugin  │  ActionPlugin     │
│  (Jira, GitLab)   │  (Custom UI) │  (Custom Commands)│
└──────────────────────────────────────────────────────┘
```

### Core Components (Not Pluggable)

| Component | Rationale |
|-----------|-----------|
| Knowledge Graph | Single source of truth; consistency required |
| TUI Shell | UX consistency; Bubble Tea contract |
| State Management | Elm architecture; purity required |
| Confirmation Gates | Security; human-in-loop non-negotiable |
| Audit Trail | Compliance; all actions logged |

### Plugin Types

#### 1. DataSourcePlugin (Primary)

```go
// Provides nodes and edges from external systems
type DataSourcePlugin interface {
    // Identity
    Name() string
    Version() string

    // Capabilities
    NodeTypes() []NodeType
    EdgeTypes() []EdgeType
    SupportsWrite() bool

    // Lifecycle
    Init(ctx context.Context, config PluginConfig) error
    Shutdown(ctx context.Context) error
    HealthCheck(ctx context.Context) error

    // Read operations (always allowed)
    FetchNodes(ctx context.Context, filter NodeFilter) ([]Node, error)
    FetchEdges(ctx context.Context, nodeIDs []string) ([]Edge, error)

    // Write operations (optional, require confirmation)
    WriteNode(ctx context.Context, node Node) (ConfirmRequest, error)
    DeleteNode(ctx context.Context, nodeID string) (ConfirmRequest, error)
}

type PluginConfig struct {
    APIKey    string            `json:"api_key"`
    BaseURL   string            `json:"base_url"`
    TeamID    string            `json:"team_id"`
    Extra     map[string]string `json:"extra"`
}

type NodeFilter struct {
    Types        []NodeType
    UpdatedAfter time.Time
    Limit        int
    Cursor       string
}
```

#### 2. ViewPlugin (Advanced)

```go
// Custom visualization for specific node types
type ViewPlugin interface {
    Name() string

    // Which nodes this view handles
    HandlesNodeType(NodeType) bool

    // Render custom view for a node
    RenderNode(node Node, width, height int) string

    // Handle input for custom view
    HandleInput(node Node, key tea.KeyMsg) (tea.Cmd, bool)
}
```

#### 3. ActionPlugin (Advanced)

```go
// Custom commands for the command palette
type ActionPlugin interface {
    Name() string

    // Commands this plugin provides
    Commands() []Command

    // Execute a command
    Execute(ctx context.Context, cmd string, args []string, model Model) (tea.Cmd, error)
}
```

### Plugin Discovery & Loading

```go
type PluginManager struct {
    dataSources []DataSourcePlugin
    views       []ViewPlugin
    actions     []ActionPlugin
    configPath  string
}

func NewPluginManager(configPath string) (*PluginManager, error) {
    pm := &PluginManager{configPath: configPath}

    // Load built-in plugins
    pm.RegisterDataSource(&LinearPlugin{})
    pm.RegisterDataSource(&GitHubPlugin{})
    pm.RegisterDataSource(&FilesystemPlugin{})

    // Load external plugins from config directory
    // ~/.config/maat/plugins/
    pluginDir := filepath.Join(configPath, "plugins")

    // Go plugins (.so files)
    goPlugins, _ := filepath.Glob(filepath.Join(pluginDir, "*.so"))
    for _, path := range goPlugins {
        if err := pm.loadGoPlugin(path); err != nil {
            log.Warn("Failed to load plugin", "path", path, "error", err)
        }
    }

    // WASM plugins (.wasm files) - future
    // wasmPlugins, _ := filepath.Glob(filepath.Join(pluginDir, "*.wasm"))

    return pm, nil
}

func (pm *PluginManager) loadGoPlugin(path string) error {
    p, err := plugin.Open(path)
    if err != nil {
        return err
    }

    // Look for NewDataSourcePlugin symbol
    sym, err := p.Lookup("NewDataSourcePlugin")
    if err == nil {
        factory := sym.(func() DataSourcePlugin)
        pm.RegisterDataSource(factory())
    }

    // Look for NewViewPlugin symbol
    sym, err = p.Lookup("NewViewPlugin")
    if err == nil {
        factory := sym.(func() ViewPlugin)
        pm.RegisterView(factory())
    }

    return nil
}
```

### Plugin Configuration

```yaml
# ~/.config/maat/plugins.yaml
plugins:
  # Built-in plugins
  linear:
    enabled: true
    config:
      api_key: ${LINEAR_API_KEY}
      team_id: "TEAM123"

  github:
    enabled: true
    config:
      token: ${GITHUB_TOKEN}
      org: "myorg"

  # External plugins
  jira:
    enabled: true
    path: "~/.config/maat/plugins/jira.so"
    config:
      base_url: "https://mycompany.atlassian.net"
      api_token: ${JIRA_TOKEN}
      project: "PROJ"

  custom_wiki:
    enabled: false
    path: "~/.config/maat/plugins/wiki.so"
    config:
      url: "https://wiki.internal"
```

### Example: Jira Plugin

```go
// plugins/jira/plugin.go
package main

import (
    "context"
    "maat/plugin"
)

type JiraPlugin struct {
    client *jira.Client
    config plugin.PluginConfig
}

func NewDataSourcePlugin() plugin.DataSourcePlugin {
    return &JiraPlugin{}
}

func (p *JiraPlugin) Name() string { return "jira" }
func (p *JiraPlugin) Version() string { return "1.0.0" }

func (p *JiraPlugin) NodeTypes() []plugin.NodeType {
    return []plugin.NodeType{
        plugin.NodeTypeIssue,
        plugin.NodeTypeEpic,
        plugin.NodeTypeSprint,
    }
}

func (p *JiraPlugin) Init(ctx context.Context, config plugin.PluginConfig) error {
    p.config = config
    p.client = jira.NewClient(config.BaseURL, config.APIKey)
    return nil
}

func (p *JiraPlugin) FetchNodes(ctx context.Context, filter plugin.NodeFilter) ([]plugin.Node, error) {
    jql := fmt.Sprintf("project = %s AND updated > %s",
        p.config.Extra["project"],
        filter.UpdatedAfter.Format("2006-01-02"),
    )

    issues, err := p.client.Search(ctx, jql, filter.Limit)
    if err != nil {
        return nil, err
    }

    nodes := make([]plugin.Node, len(issues))
    for i, issue := range issues {
        nodes[i] = plugin.Node{
            ID:     "jira:" + issue.Key,
            Type:   plugin.NodeTypeIssue,
            Source: "jira",
            Data:   mustMarshal(issue),
            Metadata: plugin.NodeMetadata{
                CreatedAt: issue.Fields.Created,
                UpdatedAt: issue.Fields.Updated,
            },
        }
    }
    return nodes, nil
}

// Build with: go build -buildmode=plugin -o jira.so
```

### Security Considerations

```go
// Plugins are sandboxed via interface contract
// - Cannot access core state directly
// - Must return ConfirmRequest for writes
// - All actions logged via core audit trail

type PluginSandbox struct {
    plugin      DataSourcePlugin
    rateLimiter *rate.Limiter
    timeout     time.Duration
}

func (s *PluginSandbox) FetchNodes(ctx context.Context, filter NodeFilter) ([]Node, error) {
    // Rate limiting
    if err := s.rateLimiter.Wait(ctx); err != nil {
        return nil, err
    }

    // Timeout enforcement
    ctx, cancel := context.WithTimeout(ctx, s.timeout)
    defer cancel()

    // Delegate to plugin
    nodes, err := s.plugin.FetchNodes(ctx, filter)

    // Validate returned nodes
    for _, node := range nodes {
        if err := validateNode(node); err != nil {
            return nil, fmt.Errorf("plugin %s returned invalid node: %w",
                s.plugin.Name(), err)
        }
    }

    return nodes, err
}
```

## Consequences

### Positive
- **Extensibility**: Support any data source
- **Community**: Enable third-party contributions
- **Enterprise**: Custom internal tool integration
- **Maintenance**: Core team focuses on core

### Negative
- **Complexity**: Plugin API is additional surface
- **Security**: Plugins run with app permissions
- **Compatibility**: API versioning needed
- **Testing**: Plugin matrix grows

### Mitigations
- Start with DataSourcePlugin only (V1)
- Plugin review process for marketplace
- Semantic versioning for plugin API
- Plugin test harness provided

## Compliance

This ADR enforces:
- **Commandment #2**: Plugins have single responsibility
- **Commandment #5**: Plugin effects go through core

## References

- Go Plugin System: https://pkg.go.dev/plugin
- WASM Plugins: https://github.com/nicholasjackson/wasm-plugins
- MAAT Constitution: `/Users/manu/Documents/LUXOR/MAAT/docs/CONSTITUTION.md`
