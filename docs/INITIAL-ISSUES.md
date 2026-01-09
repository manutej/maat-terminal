# MAAT Initial Issues Template

These issues will be created in Linear once MCP integration is active.

## Issue 1: Add GitHub API Integration

**Type**: Feature
**Priority**: High
**Labels**: enhancement, datasource

**Title**: Add GitHub API Integration

**Description**:
Implement GitHubClient data source to fetch real issues and PRs from GitHub repositories for display in the knowledge graph.

**Requirements**:
- Implement `internal/github/client.go` following thin API client pattern (Commandment #7)
- Support GITHUB_TOKEN environment variable for authentication
- Fetch issues, PRs, and commits for a given repository
- Map GitHub entities to graph nodes/edges (Commandment #2)
- Return data in read-only format (no write operations per Commandment #10)
- Use tea.Cmd for async fetching (Commandment #5)

**Acceptance Criteria**:
- [ ] Can authenticate with GitHub using GITHUB_TOKEN
- [ ] Can fetch issues from a repository
- [ ] Can fetch PRs from a repository
- [ ] Can fetch recent commits
- [ ] All entities mapped to graph nodes
- [ ] No mutations to GitHub state
- [ ] All async operations use tea.Cmd

**Reference**: `specs/FUNCTIONAL-REQUIREMENTS.md` FR-003

---

## Issue 2: Add Linear API Integration

**Type**: Feature
**Priority**: High
**Labels**: enhancement, datasource

**Title**: Add Linear API Integration

**Description**:
Implement LinearSource data source to fetch issues, projects, and cycles from Linear for display in the knowledge graph.

**Requirements**:
- Implement `internal/linear/client.go` following thin API client pattern
- Use Linear MCP or GraphQL API directly
- Support LINEAR_API_KEY environment variable
- Fetch issues, projects, cycles, and relationships
- Map Linear entities to graph nodes/edges
- Return data in read-only format initially
- Use tea.Cmd for async fetching

**Acceptance Criteria**:
- [ ] Can authenticate with Linear using LINEAR_API_KEY or MCP
- [ ] Can fetch issues from workspace
- [ ] Can fetch projects and cycles
- [ ] Can fetch issue relationships (blocks, related)
- [ ] All entities mapped to graph nodes
- [ ] No mutations to Linear state (read-only initially)
- [ ] All async operations use tea.Cmd

**Reference**: `specs/FUNCTIONAL-REQUIREMENTS.md` FR-003

---

## Issue 3: Improve Search/Filter Capabilities

**Type**: Enhancement
**Priority**: Medium
**Labels**: enhancement, ux

**Title**: Improve Search/Filter Capabilities

**Description**:
Add comprehensive search and filter functionality for navigating large knowledge graphs efficiently.

**Requirements**:
- Text search across node titles and descriptions
- Filter by entity type (issue, PR, commit, file)
- Filter by status (open, closed, merged, etc.)
- Filter by assignee/author
- Filter by date range
- Maintain keyboard-driven interface (Commandment #4)
- Persist filter state in Model immutably (Commandment #1)

**Acceptance Criteria**:
- [ ] Text search with incremental filtering
- [ ] Type filter (issue/PR/commit/file)
- [ ] Status filter with multi-select
- [ ] Assignee/author filter
- [ ] Date range filter
- [ ] All filters composable (can combine multiple)
- [ ] Clear visual indication of active filters
- [ ] Keyboard shortcuts for common filters

**Reference**: `specs/FUNCTIONAL-REQUIREMENTS.md` FR-001

---

## Issue 4: Add Node Actions

**Type**: Enhancement
**Priority**: Medium
**Labels**: enhancement, ux

**Title**: Add Node Actions Menu

**Description**:
Implement action menu for selected nodes, allowing users to interact with entities while respecting safety constraints.

**Requirements**:
- Action menu accessible via keyboard shortcut
- Available actions based on node type
- All write operations require confirmation (Commandment #10)
- Actions implemented as tea.Cmd (Commandment #5)
- Maintain immutable Model state (Commandment #1)

**Actions to Support**:
1. Open in browser (read-only)
2. Copy URL to clipboard (read-only)
3. View full details (read-only)
4. Change status (requires confirmation)
5. Add comment (requires confirmation)
6. Link to another node (requires confirmation)

**Acceptance Criteria**:
- [ ] Action menu UI component
- [ ] Keyboard shortcut to open actions (e.g., 'a')
- [ ] Context-aware actions based on node type
- [ ] All write actions show ConfirmRequest dialog
- [ ] Read-only actions execute immediately
- [ ] Success/error feedback messages
- [ ] All actions as pure tea.Cmd functions

**Reference**: `specs/FUNCTIONAL-REQUIREMENTS.md` FR-005, Commandment #10

---

## Issue 5: Fix Graph View Scrolling

**Type**: Bug
**Priority**: Low
**Labels**: bug, ux

**Title**: Fix Graph View Scrolling for Large Graphs

**Description**:
Large knowledge graphs overflow the terminal viewport, making navigation difficult. Implement scrolling or pagination to handle large datasets.

**Current Behavior**:
- Graphs with 50+ nodes overflow the viewport
- No way to scroll to see all nodes
- Bottom nodes are cut off

**Requirements**:
- Implement viewport-based scrolling
- Or implement pagination with keyboard navigation
- Maintain current focus when scrolling
- Visual indicators for scroll position
- Smooth keyboard-driven scrolling (Commandment #4)

**Acceptance Criteria**:
- [ ] Can view all nodes in graphs with 100+ nodes
- [ ] Scroll indicators show position in graph
- [ ] Focused node stays visible when scrolling
- [ ] Keyboard shortcuts for page up/down
- [ ] Performance remains smooth with large graphs
- [ ] Current view state immutable in Model

**Reference**: `specs/FUNCTIONAL-REQUIREMENTS.md` FR-001

---

## Issue 6: Add README and Usage Guide

**Type**: Documentation
**Priority**: Low
**Labels**: documentation

**Title**: Add README and Usage Guide

**Description**:
Create comprehensive documentation for users to understand MAAT's purpose, setup, and usage.

**Requirements**:
- README.md with project overview
- Quick start guide
- Integration setup (Linear, GitHub)
- Keybindings reference
- Configuration options
- Architecture overview
- Contributing guide

**Sections Needed**:
1. Project overview and philosophy
2. Installation instructions
3. Quick start (first run)
4. Keybindings reference table
5. Integration setup (Linear, GitHub, Claude)
6. Configuration file format
7. Architecture overview (reference CONSTITUTION.md)
8. Troubleshooting
9. Contributing guidelines
10. License information

**Acceptance Criteria**:
- [ ] README.md covers all sections
- [ ] Keybindings documented in table format
- [ ] Integration setup includes LINEAR_API_KEY and GITHUB_TOKEN
- [ ] Architecture references 10 Commandments
- [ ] Quick start gets user to working TUI in < 5 minutes
- [ ] Screenshots or ASCII demo of TUI
- [ ] Links to relevant spec documents

**Reference**: `docs/CONSTITUTION.md`, `specs/FUNCTIONAL-REQUIREMENTS.md`

---

## Team Configuration

**Team Name**: MAAT
**Project**: MAAT Terminal (if creating new project)
**Workflow States**:
- Todo
- In Progress
- In Review
- Done

**Labels to Create**:
- enhancement
- bug
- documentation
- datasource
- ux
- architecture
- testing
