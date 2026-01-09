# MAAT Workflows

This directory contains executable workflow specifications for MAAT implementation using the OIS (Operator-based Integration System) agent orchestration framework.

## Phase 1: Foundation (Current)

### Workflows

| Workflow | File | Status | Duration |
|----------|------|--------|----------|
| **1A: SQLite Store** | `phase1-sqlite-store.yaml` | ✅ Ready | 60-90 min |
| **1B: Mock Data** | `phase1-mock-data.yaml` | ✅ Ready | 45-60 min |

### Quick Start

```bash
# Execute both Phase 1 workflows
./workflows/execute-phase1.sh

# Or with verbose output
./workflows/execute-phase1.sh --verbose

# Or dry run (preview without execution)
./workflows/execute-phase1.sh --dry-run
```

### Manual Execution

```bash
# Execute Phase 1A only
/ois-compose --workflow-plan workflows/phase1-sqlite-store.yaml \
  --input-file internal/graph/schema.go \
  --output internal/graph/store_generated.go \
  --verbose

# Execute Phase 1B only
/ois-compose --workflow-plan workflows/phase1-mock-data.yaml \
  --input-file "docs/CONSTITUTION.md specs/FUNCTIONAL-REQUIREMENTS.md" \
  --output internal/tui/mock_data.go
```

### Expected Outputs

After successful execution:

```
internal/graph/
├── schema.go              (existing - 170 lines)
├── store_generated.go     (generated - ~500 lines)
└── store_test.go          (generated - ~300 lines)

internal/tui/
├── model.go               (existing - 258 lines)
├── mock_data.go           (generated - ~300 lines)
└── ...
```

### Validation

```bash
# Run tests
make test

# Build binary
make build

# Test TUI with mock data
./maat
```

**Expected**: TUI displays mock graph with ~95 nodes and ~98 edges.

---

## Workflow Structure

Each YAML workflow contains:

### 1. Metadata
- Task description
- Complexity level (L1-L7)
- Token budget and duration
- Quality targets

### 2. Composition
- Block definitions (structure-observer, type-generator, etc.)
- Composition operators (→, ||, ×, ⊗, IF, UNTIL)
- Agent configurations and prompts

### 3. Type Flow
- Input → Intermediate → Output types
- Type safety validation

### 4. Success Criteria
- Functional requirements
- Non-functional requirements
- Constitutional alignment
- Testing requirements

### 5. Execution Instructions
- Prerequisites
- Commands
- Failure recovery

---

## Complexity Levels

| Level | Function | Phase 1 Usage |
|-------|----------|---------------|
| **L1** | OBSERVE → GENERATE | ✅ Both workflows |
| **L2** | + REASON | Phase 2 |
| **L3** | + Complex generation | Phase 3 |
| **L4** | + REFINE (human-in-loop) | Phase 4 |
| **L5** | + OPTIMIZE | Phase 5 |
| **L6** | + INTEGRATE | Future |
| **L7** | + REFLECT | Future |

**Phase 1 Philosophy**: Start at L1 (simplest), increase only when proven necessary.

---

## Atomic Blocks Used

### Phase 1A: SQLite Store
- `structure-observer` (systems-thinking × OBSERVE)
- `type-generator` (category-theory × GENERATE)

### Phase 1B: Mock Data
- `knowledge-gatherer` (knowledge-synthesis × OBSERVE)
- `mock-generator` (specification-driven × GENERATE)

---

## Composition Operators

### Phase 1: Sequential Only

```yaml
# Both workflows use simple sequential composition
composition:
  pattern: sequential
  operators: ["→"]

  step_1: BlockA
  step_2: BlockB
```

**Why sequential**:
- BlockB depends on BlockA output
- No independent operations to parallelize
- Token efficient (no synthesis overhead)
- Debugging simple (linear failure path)

**Parallel operators** (`||`) introduced in Phase 4 (Claude MCP).

---

## Token Budget

| Workflow | Observe | Generate | Total | % of Phase 1 |
|----------|---------|----------|-------|--------------|
| Phase 1A | 800 | 1,200 | 2,000 | 55.6% |
| Phase 1B | 600 | 1,000 | 1,600 | 44.4% |
| **Total** | **1,400** | **2,200** | **3,600** | **100%** |

**Phase 1 Budget**: 3,600 tokens / 35,400 total (10.2% of project)

---

## Constitutional Alignment

### Commandments Referenced

| Commandment | Phase 1A | Phase 1B |
|-------------|----------|----------|
| #1 Immutable Truth | ✅ Pure functions | ✅ Deterministic data |
| #2 Single Responsibility | ✅ Store interface | ✅ Entity focus |
| #3 Text Interface | ✅ JSON data | ✅ JSON data |
| #5 Controlled Effects | ✅ DB via methods | N/A |
| #9 Specification Constitution | ✅ This YAML | ✅ This YAML |

---

## Troubleshooting

### Workflow Fails to Execute

```bash
# Check prerequisites
go version            # Go 1.21+?
ls internal/graph/schema.go  # File exists?

# Check /ois-compose command available
which /ois-compose    # Command found?

# Validate YAML syntax
yamllint workflows/phase1-sqlite-store.yaml
```

### Generated Code Doesn't Compile

```bash
# Check imports
go mod tidy

# Check for missing dependencies
go get github.com/mattn/go-sqlite3

# Review compilation errors
go build ./internal/graph
```

### Tests Fail

```bash
# Run with verbose output
go test -v ./internal/graph

# Check test coverage
go test -cover ./internal/graph

# Review specific test
go test -run TestAddNode ./internal/graph
```

### Mock Data Unrealistic

```bash
# Re-run Phase 1B with refined prompts
/ois-compose --workflow-plan workflows/phase1-mock-data.yaml \
  --input-file "docs/CONSTITUTION.md specs/FUNCTIONAL-REQUIREMENTS.md" \
  --output internal/tui/mock_data.go

# Or manually edit internal/tui/mock_data.go
```

---

## Next Steps

After Phase 1 completion:

1. **Validate**: Run `./workflows/execute-phase1.sh` and verify success
2. **Test**: Run `./maat` and confirm TUI displays mock graph
3. **Proceed**: Move to Phase 2 workflows (graph rendering + keyboard navigation)

---

## References

- **Orchestration Plan**: `docs/AGENT-COMPOSITION-PLAN.md`
- **Atomic Blocks**: `docs/ATOMIC-BLOCKS.md`
- **Execution Plan**: `docs/EXECUTION-PLAN.md`
- **MAAT Spec**: `MAAT-SPEC.md`
- **Constitution**: `docs/CONSTITUTION.md`

---

**Status**: Phase 1 workflows ready to execute ✓
**Next Action**: Run `./workflows/execute-phase1.sh`
