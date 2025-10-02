# Claude Agent SDK Team Implementation

Complete implementation of agent team architecture following patterns from [Anthropic's Claude Agent SDK](https://www.anthropic.com/engineering/building-agents-with-the-claude-agent-sdk).

## Architecture Overview

This implementation follows the core Claude SDK patterns:

1. **Agent Loop**: gather context → take action → verify work → repeat
2. **Subagent Coordination**: Orchestrator delegates to specialized agents
3. **Agentic Search**: Primary search using file system structure
4. **Context Management**: Automatic compaction with LRU scoring
5. **Tool System**: Primary building blocks for agent execution

```
┌─────────────────────────────────────────────────────────────┐
│                    OrchestratorAgent                        │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │         Main Agent Loop (Max Iterations)           │   │
│  │                                                     │   │
│  │  ┌──────────────────────────────────────────────┐ │   │
│  │  │  Phase 1: Gather Context                     │ │   │
│  │  │  - SearchRelevantFiles (agentic search)      │ │   │
│  │  │  - Get previous iteration feedback           │ │   │
│  │  │  - Identify required subagents               │ │   │
│  │  └──────────────────────────────────────────────┘ │   │
│  │              ↓                                      │   │
│  │  ┌──────────────────────────────────────────────┐ │   │
│  │  │  Context Compaction Check                    │ │   │
│  │  │  - Trigger at 80% capacity                   │ │   │
│  │  │  - Remove bottom 30% by LRU score            │ │   │
│  │  └──────────────────────────────────────────────┘ │   │
│  │              ↓                                      │   │
│  │  ┌──────────────────────────────────────────────┐ │   │
│  │  │  Phase 2: Take Action                        │ │   │
│  │  │  - Execute subagents (parallel/sequential)   │ │   │
│  │  │  - Execute required tools                    │ │   │
│  │  └──────────────────────────────────────────────┘ │   │
│  │              ↓                                      │   │
│  │  ┌──────────────────────────────────────────────┐ │   │
│  │  │  Phase 3: Verify Work                        │ │   │
│  │  │  - Rule-based validation                     │ │   │
│  │  │  - Collect artifacts (files created)         │ │   │
│  │  │  - Check subagent success                    │ │   │
│  │  └──────────────────────────────────────────────┘ │   │
│  │              ↓                                      │   │
│  │  ┌──────────────────────────────────────────────┐ │   │
│  │  │  Complete? → Exit : Continue with Feedback   │ │   │
│  │  └──────────────────────────────────────────────┘ │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │         Subagent Coordination                      │   │
│  │                                                     │   │
│  │  Parallel Execution (when >1 subagent needed):     │   │
│  │  ┌──────────────┐ ┌──────────────┐ ┌────────────┐│   │
│  │  │ SearchAgent  │ │  CodeAgent   │ │ TestAgent  ││   │
│  │  │              │ │              │ │            ││   │
│  │  │ Isolated     │ │ Isolated     │ │ Isolated   ││   │
│  │  │ Context      │ │ Context      │ │ Context    ││   │
│  │  └──────────────┘ └──────────────┘ └────────────┘│   │
│  │         ↓                ↓                ↓         │   │
│  │  ┌──────────────────────────────────────────────┐ │   │
│  │  │    Semaphore-based Concurrency Control       │ │   │
│  │  │    (max concurrent = config.ParallelSubagents)│ │   │
│  │  └──────────────────────────────────────────────┘ │   │
│  └────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                    ContextManager                           │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │         Agentic Search (Primary Method)            │   │
│  │                                                     │   │
│  │  Strategy 1: File System Structure Search          │   │
│  │  ┌──────────────────────────────────────────────┐ │   │
│  │  │ Keyword → Directory Mapping:                 │ │   │
│  │  │  "test"      → tests/, __tests__/, spec/     │ │   │
│  │  │  "component" → components/, src/components/  │ │   │
│  │  │  "service"   → services/, backend/services/  │ │   │
│  │  │  "model"     → models/, src/models/          │ │   │
│  │  └──────────────────────────────────────────────┘ │   │
│  │              ↓                                      │   │
│  │  Strategy 2: Search Index (Supplementary)          │   │
│  │  ┌──────────────────────────────────────────────┐ │   │
│  │  │ Keyword-based content indexing               │ │   │
│  │  └──────────────────────────────────────────────┘ │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │         Context Compaction System                  │   │
│  │                                                     │   │
│  │  Trigger: Usage >= 80% of contextLimit             │   │
│  │  ┌──────────────────────────────────────────────┐ │   │
│  │  │ Retention Score Calculation:                 │ │   │
│  │  │  score = (AccessCount × 10.0)                │ │   │
│  │  │        + Relevance                           │ │   │
│  │  │        - (FileSize / 1000.0)                 │ │   │
│  │  └──────────────────────────────────────────────┘ │   │
│  │              ↓                                      │   │
│  │  Remove bottom 30% of files by score               │   │
│  │  Target: Reduce to 70% of contextLimit             │   │
│  └────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                    ToolRegistry                             │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │         Registered Tools                           │   │
│  │                                                     │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────┐│   │
│  │  │  BashTool    │  │ ReadFileTool │  │WriteTool ││   │
│  │  │              │  │              │  │          ││   │
│  │  │ Execute shell│  │ Read file    │  │Write file││   │
│  │  │ commands     │  │ contents     │  │contents  ││   │
│  │  └──────────────┘  └──────────────┘  └──────────┘│   │
│  │                                                     │   │
│  │  ┌──────────────┐  ┌──────────────┐               │   │
│  │  │CodeGenerator │  │TestGenerator │               │   │
│  │  │              │  │              │               │   │
│  │  │Generate code │  │Generate tests│               │   │
│  │  │from context  │  │for code      │               │   │
│  │  └──────────────┘  └──────────────┘               │   │
│  └────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Core Components

### 1. OrchestratorAgent (`orchestrator.go`)

Main coordinator implementing the agent loop pattern.

**Key Features**:
- Three-phase agent loop (gather → act → verify)
- Parallel and sequential subagent execution
- Automatic context compaction
- Iteration with feedback loops
- Validation rules for task types

**Configuration**:
```go
config := &AgentConfig{
    MaxIterations:     5,      // Maximum agent loop iterations
    ContextLimit:      100000, // Context window size limit
    ParallelSubagents: 3,      // Max concurrent subagents
    TimeoutSeconds:    300,    // Operation timeout
    EnableCompaction:  true,   // Auto context compaction
}
```

### 2. Specialized Subagents (`subagents.go`)

Three specialized agents that can operate independently or be orchestrated:

#### SearchAgent
- **Purpose**: Find relevant files using agentic search
- **Strategy**: File system structure-based (primary), content search (supplementary)
- **Capabilities**: file-search, agentic-search, content-search, directory-traversal

#### CodeAgent
- **Purpose**: Generate and modify code
- **Strategy**: Context-aware code generation with iterative refinement
- **Capabilities**: code-generation, code-modification, refactoring, linting

#### TestAgent
- **Purpose**: Generate and execute tests
- **Strategy**: Automated test generation with coverage analysis
- **Capabilities**: test-generation, test-execution, coverage-analysis

### 3. ContextManager (`context_manager.go`)

Manages agent context following Claude SDK patterns.

**Key Features**:
- Agentic search using file system structure (primary method)
- Automatic context compaction at 80% capacity
- LRU-based file removal scoring
- Search index for keyword-based lookup
- Relevance scoring for search results

**Compaction Strategy**:
```
Trigger: contextSize >= 80% of contextLimit
Score = (AccessCount × 10.0) + Relevance - (FileSize / 1000.0)
Action: Remove bottom 30% of cached files by score
Target: Reduce to 70% of contextLimit
```

### 4. Tool System (`tools.go`)

Five core tools for agent execution:

1. **BashTool**: Execute shell commands
2. **ReadFileTool**: Read file contents
3. **WriteFileTool**: Write files with directory creation
4. **CodeGeneratorTool**: Generate code from task descriptions
5. **TestGeneratorTool**: Generate tests for code files

## Usage Examples

### Basic Agent Execution

```go
package main

import (
    "context"
    "fmt"
    "github.com/tchat/backend/agents"
)

func main() {
    // Create orchestrator with configuration
    config := &agents.AgentConfig{
        MaxIterations:     5,
        ContextLimit:      100000,
        ParallelSubagents: 3,
        EnableCompaction:  true,
    }
    orchestrator := agents.NewOrchestratorAgent(config)

    // Register specialized subagents
    contextManager := agents.NewContextManager(config.ContextLimit)
    toolRegistry := agents.NewToolRegistry()

    orchestrator.RegisterSubagent("search-agent", agents.NewSearchAgent(contextManager))
    orchestrator.RegisterSubagent("code-agent", agents.NewCodeAgent(contextManager, toolRegistry))
    orchestrator.RegisterSubagent("test-agent", agents.NewTestAgent(contextManager, toolRegistry))

    // Create task
    task := &agents.Task{
        ID:          "task-001",
        Type:        "code",
        Description: "Create a user authentication function",
        Context:     make(map[string]interface{}),
        Priority:    1,
    }

    // Execute task
    result, err := orchestrator.Execute(context.Background(), task)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Check results
    if result.Success {
        fmt.Printf("Task completed successfully in %d iterations\n", result.Iterations)
        fmt.Printf("Artifacts created: %v\n", result.Artifacts)
    } else {
        fmt.Printf("Task failed after %d iterations\n", result.Iterations)
        for _, err := range result.Errors {
            fmt.Printf("Error: %v\n", err)
        }
    }
}
```

### Using Individual Subagents

```go
// Use SearchAgent directly
searchAgent := agents.NewSearchAgent(contextManager)
searchTask := &agents.Task{
    ID:          "search-001",
    Description: "Find all test files",
}

searchResult, err := searchAgent.Execute(context.Background(), searchTask)
if err == nil {
    fmt.Printf("Found files: %v\n", searchResult.Artifacts)
}

// Use CodeAgent directly
codeAgent := agents.NewCodeAgent(contextManager, toolRegistry)
codeTask := &agents.Task{
    ID:          "code-001",
    Description: "Generate user service in Go",
}

codeResult, err := codeAgent.Execute(context.Background(), codeTask)
if err == nil {
    fmt.Printf("Generated code files: %v\n", codeResult.Artifacts)
}
```

### Using Tools Directly

```go
// Use BashTool
toolRegistry := agents.NewToolRegistry()
bashTool, _ := toolRegistry.GetTool("bash")

input := &agents.ToolInput{
    Command: "go test ./...",
}

output, err := bashTool.Execute(context.Background(), input)
if err == nil {
    fmt.Printf("Test output: %s\n", output.Content["stdout"])
}

// Use CodeGeneratorTool
codeGenTool, _ := toolRegistry.GetTool("code-generator")

input := &agents.ToolInput{
    Task: &agents.Task{
        Description: "Create authentication middleware",
    },
    RelevantFiles: []string{"backend/auth/user.go"},
}

output, err := codeGenTool.Execute(context.Background(), input)
if err == nil {
    fmt.Printf("Generated files: %v\n", output.Files)
}
```

### Context Management

```go
// Create context manager
contextManager := agents.NewContextManager(100000) // 100KB limit

// Build search index
if err := contextManager.BuildIndex("."); err != nil {
    log.Fatal(err)
}

// Search for relevant files
files, err := contextManager.SearchRelevantFiles("test authentication")
if err == nil {
    fmt.Printf("Found %d relevant files\n", len(files))
}

// Cache a file
content, _ := os.ReadFile("backend/auth/service.go")
contextManager.CacheFile("backend/auth/service.go", string(content))

// Check if compaction is needed
if contextManager.ShouldCompact() {
    fmt.Println("Context compaction triggered")
    contextManager.Compact(context.Background())
}

// Get utilization
fmt.Printf("Context utilization: %.2f%%\n", contextManager.GetUtilization())
```

## Claude SDK Pattern Implementation

### 1. Agent Loop ✅

```go
for iteration := 0; iteration < maxIterations; iteration++ {
    // Phase 1: Gather Context
    contextData, err := o.gatherContext(ctx, task, iteration)

    // Phase 2: Take Action
    actionResult, err := o.takeAction(ctx, task, contextData)

    // Phase 3: Verify Work
    verified, err := o.verifyWork(ctx, task, actionResult)

    if verified.Complete {
        break
    }

    // Continue with feedback
    task.Context["previous_iteration"] = verified.Feedback
}
```

### 2. Agentic Search ✅

Primary method using file system structure:

```go
// Map keywords to directory structures
switch keyword {
case "test":
    searchPaths = []string{"tests", "__tests__", "test", "spec"}
case "component":
    searchPaths = []string{"components", "src/components", "ui"}
case "service":
    searchPaths = []string{"services", "src/services", "backend/services"}
}
```

### 3. Context Compaction ✅

Automatic compaction when approaching context limits:

```go
// Trigger at 80% capacity
if utilization >= 0.8 {
    // Calculate retention score
    score := (AccessCount × 10.0) + Relevance - (Size / 1000.0)

    // Remove bottom 30% by score
    // Target: 70% of contextLimit
}
```

### 4. Parallel Execution ✅

Goroutine-based parallel subagent execution:

```go
semaphore := make(chan struct{}, maxConcurrent)

for _, agentName := range agentNames {
    go func(name string) {
        semaphore <- struct{}{}
        defer func() { <-semaphore }()

        result, err := agent.Execute(ctx, subTask)
    }(agentName)
}
```

### 5. Tools as Building Blocks ✅

Tools are the primary execution mechanism:

```go
type Tool interface {
    Execute(ctx context.Context, input *ToolInput) (*ToolOutput, error)
    GetName() string
    GetDescription() string
}

// Tools: bash, read-file, write-file, code-generator, test-generator
```

## Configuration Options

```go
type AgentConfig struct {
    MaxIterations     int     // Max agent loop iterations (default: 5)
    ContextLimit      int     // Context size limit in bytes (default: 100000)
    ParallelSubagents int     // Max concurrent subagents (default: 3)
    TimeoutSeconds    int     // Operation timeout (default: 300)
    EnableCompaction  bool    // Auto context compaction (default: true)
}
```

## Testing

```bash
# Run all agent tests
cd backend/agents
go test ./... -v

# Run with coverage
go test ./... -cover

# Run specific test
go test -v -run TestOrchestratorExecute
```

## Integration with Tchat

The agent system can be integrated into the Tchat backend for:

1. **Automated Code Analysis**: Analyze codebase and suggest improvements
2. **Test Generation**: Automatically generate tests for new code
3. **Code Review**: Review pull requests and identify issues
4. **Documentation**: Generate documentation from code
5. **Refactoring**: Suggest and implement refactorings

## Performance Considerations

- **Context Size**: Keep context limit appropriate for your use case
- **Parallel Execution**: Adjust `ParallelSubagents` based on system resources
- **Compaction Frequency**: Triggered at 80% by default, adjustable via `compactionRatio`
- **Iteration Limit**: Set `MaxIterations` to prevent infinite loops
- **Tool Timeout**: Set appropriate `TimeoutSeconds` for long-running operations

## Future Enhancements

- [ ] Add more specialized agents (RefactorAgent, SecurityAgent, PerformanceAgent)
- [ ] Implement agent learning from previous executions
- [ ] Add more sophisticated validation rules
- [ ] Implement agent collaboration protocols
- [ ] Add metrics and monitoring
- [ ] Implement agent persistence across sessions
- [ ] Add support for custom tools
- [ ] Implement cost tracking for agent operations

## References

- [Anthropic Claude Agent SDK](https://www.anthropic.com/engineering/building-agents-with-the-claude-agent-sdk)
- [Agentic Systems Architecture](https://www.anthropic.com/research/agentic-systems)
- [Claude SDK Patterns](https://docs.anthropic.com/claude/docs/agents-best-practices)