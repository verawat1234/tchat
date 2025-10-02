# Using Agent Team via Claude Code

Complete guide for using the agent team system through Claude Code interface.

## Quick Start

### 1. Basic Task Execution via Claude Code

You can ask Claude Code to use the agent system for various tasks:

```
"Use the agent system to find all authentication-related files"
"Use the agent system to generate tests for the social repository"
"Use the agent system to refactor the video service"
```

Claude Code will internally:
1. Create the appropriate task
2. Execute the orchestrator
3. Coordinate subagents
4. Return results to you

### 2. Direct Command Examples

#### Example 1: Search for Files
```bash
# Ask Claude Code:
"Use SearchAgent to find all test files in the project"

# Claude Code will execute:
# - Create SearchAgent
# - Build search index
# - Use agentic search to find test files
# - Return list of matching files
```

#### Example 2: Generate Code
```bash
# Ask Claude Code:
"Use CodeAgent to create a new authentication middleware in Go"

# Claude Code will:
# - Search for related authentication files
# - Analyze existing patterns
# - Generate new middleware code
# - Save to appropriate location
```

#### Example 3: Generate Tests
```bash
# Ask Claude Code:
"Use TestAgent to generate tests for SocialRepository.kt"

# Claude Code will:
# - Read SocialRepository.kt
# - Analyze the code structure
# - Generate comprehensive test file
# - Save as SocialRepositoryTest.kt
```

## Integration with Claude Code Workflows

### Using with Claude Code Commands

#### `/analyze` Command with Agent System
```
/analyze --use-agents backend/social

# Claude Code will:
# 1. Use SearchAgent to find all social service files
# 2. Use CodeAgent to analyze code quality
# 3. Use TestAgent to check test coverage
# 4. Provide comprehensive analysis report
```

#### `/improve` Command with Agent System
```
/improve --use-agents apps/kmp/composeApp/src/commonMain/kotlin/com/tchat/mobile/repositories/

# Claude Code will:
# 1. SearchAgent finds all repository files
# 2. CodeAgent analyzes and suggests improvements
# 3. Orchestrator coordinates refactoring
# 4. TestAgent generates/updates tests
```

#### `/test` Command with Agent System
```
/test --use-agents --generate backend/video

# Claude Code will:
# 1. SearchAgent finds video service files
# 2. TestAgent generates missing tests
# 3. Orchestrator runs tests
# 4. Reports coverage and results
```

## Practical Use Cases

### Use Case 1: Automated Code Review

**Request to Claude Code:**
```
"Use the agent system to review the social repository implementation at
apps/kmp/composeApp/src/commonMain/kotlin/com/tchat/mobile/repositories/SocialRepository.kt"
```

**What happens:**
1. SearchAgent finds related files (models, services, handlers)
2. CodeAgent analyzes code quality, patterns, best practices
3. TestAgent checks test coverage
4. Orchestrator compiles comprehensive review report

**Output:**
```
Agent Review Report:
- Code Quality: 8/10
- Test Coverage: 75%
- Issues Found: 3
  1. Missing error handling in followUser (line 45)
  2. No unit tests for unfollowUser
  3. Consider adding retry logic for network errors
- Suggestions:
  1. Add @throws documentation
  2. Extract error handling to utility
  3. Implement exponential backoff
```

### Use Case 2: Test Generation

**Request to Claude Code:**
```
"Generate comprehensive tests for SocialRepository using the agent system"
```

**Agent Execution Flow:**
```
Orchestrator.Execute(task)
├── Phase 1: Gather Context
│   ├── SearchAgent finds SocialRepository.kt
│   └── SearchAgent finds related models and dependencies
├── Phase 2: Take Action
│   ├── CodeAgent analyzes SocialRepository methods
│   └── TestAgent generates test cases
└── Phase 3: Verify Work
    ├── Validates test syntax
    └── Checks test coverage
```

**Generated Output:**
- `SocialRepositoryTest.kt` with unit tests
- Mock implementations for dependencies
- Test coverage report

### Use Case 3: Feature Implementation

**Request to Claude Code:**
```
"Use the agent system to implement a new 'blockUser' feature in SocialRepository"
```

**Agent Workflow:**
1. **SearchAgent**: Find existing follow/unfollow implementations
2. **CodeAgent**:
   - Analyze patterns
   - Generate blockUser method
   - Generate unblockUser method
   - Update repository interface
3. **TestAgent**: Generate tests for block functionality
4. **Orchestrator**: Verify all components work together

**Files Created/Modified:**
- `SocialRepository.kt` - added blockUser/unblockUser
- `SocialRepositoryTest.kt` - added block tests
- `backend/social/handlers/user_handler.go` - added block endpoints

### Use Case 4: Refactoring

**Request to Claude Code:**
```
"Use the agent system to refactor error handling in the video service"
```

**Agent Process:**
```
1. SearchAgent:
   - Finds all video service files
   - Identifies error handling patterns

2. CodeAgent:
   - Analyzes current error handling
   - Designs unified error handling strategy
   - Refactors code to use new pattern

3. TestAgent:
   - Updates existing tests
   - Adds new error case tests

4. Orchestrator:
   - Verifies: max 3 iterations with feedback
   - Ensures no breaking changes
```

## Advanced Usage Patterns

### Pattern 1: Iterative Refinement

**Request:**
```
"Use the agent system to implement authentication middleware,
iterate until all tests pass"
```

**Agent Loop (Max 5 iterations):**
```
Iteration 1:
├── Gather: Find existing auth patterns
├── Action: Generate initial middleware
└── Verify: Run tests → 2 failures

Iteration 2:
├── Gather: Analyze test failures (feedback from iteration 1)
├── Action: Fix token validation
└── Verify: Run tests → 1 failure

Iteration 3:
├── Gather: Analyze remaining failure
├── Action: Fix error response format
└── Verify: Run tests → All pass ✅
```

### Pattern 2: Parallel Agent Execution

**Request:**
```
"Use the agent system to analyze the entire social service in parallel"
```

**Parallel Execution:**
```
Orchestrator spawns 3 agents concurrently:

┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│  SearchAgent    │  │   CodeAgent     │  │   TestAgent     │
│                 │  │                 │  │                 │
│ Find all files  │  │ Analyze quality │  │ Check coverage  │
│ Build file map  │  │ Find issues     │  │ Find gaps       │
└─────────────────┘  └─────────────────┘  └─────────────────┘
         ↓                    ↓                    ↓
    ┌────────────────────────────────────────────────┐
    │         Orchestrator combines results          │
    └────────────────────────────────────────────────┘
```

### Pattern 3: Context-Aware Operations

**Request:**
```
"Use the agent system to add social features to the mobile app,
using existing backend patterns"
```

**Context Management:**
```
1. SearchAgent builds context:
   - backend/social/handlers/*.go (backend patterns)
   - apps/web/src/components/social/*.tsx (web patterns)
   - apps/kmp/.../repositories/SocialRepository.kt (existing mobile)

2. Context Manager:
   - Caches relevant files
   - Scores by relevance
   - Triggers compaction if > 80%

3. CodeAgent generates:
   - Mobile UI components following web patterns
   - Repository methods following existing patterns
   - API calls matching backend contracts
```

## Claude Code Integration Examples

### Example 1: Full Feature Implementation

**Claude Code Command:**
```
You: "Implement a 'share post' feature across all platforms using the agent system"

Claude Code will:
1. Create orchestrator task
2. Execute agents in sequence:
   - SearchAgent: Find post-related code
   - CodeAgent: Generate share functionality
   - TestAgent: Generate tests
3. Report results with file locations
```

**Terminal Output:**
```bash
🤖 Agent System Execution Started

Phase 1: Gathering Context
  ✅ SearchAgent found 12 relevant files
  ✅ Context: 45,892 bytes (45% capacity)

Phase 2: Taking Action
  ✅ CodeAgent generated 3 files:
     - backend/social/handlers/post_handler.go (added sharePost)
     - apps/web/src/components/social/ShareButton.tsx
     - apps/kmp/.../repositories/SocialRepository.kt (added sharePost)

Phase 3: Verifying Work
  ✅ TestAgent generated 3 test files
  ✅ All tests passed

✅ Task completed in 2 iterations
📦 Artifacts: 6 files created/modified
```

### Example 2: Code Analysis

**Claude Code Command:**
```
You: "Analyze SocialRepository.kt using the agent system and suggest improvements"

Claude Code executes:
- SearchAgent: Find related files
- CodeAgent: Analyze code quality
- Orchestrator: Compile report
```

**Structured Output:**
```markdown
# Agent Analysis Report: SocialRepository.kt

## Overview
- File: apps/kmp/.../repositories/SocialRepository.kt
- Lines: 156
- Complexity: Medium

## Code Quality Analysis
- ✅ Follows Kotlin conventions
- ✅ Good error handling
- ⚠️  Missing documentation on 3 methods
- ⚠️  Could benefit from retry logic

## Test Coverage
- Current: 75%
- Missing tests:
  - unfollowUser error cases
  - getFollowersList pagination

## Suggestions
1. Add KDoc comments to public methods
2. Implement exponential backoff for API calls
3. Add integration tests for pagination
4. Consider using Result type instead of exceptions

## Automated Fixes Available
Agent can automatically:
- Generate missing tests (/generate-tests)
- Add KDoc comments (/add-docs)
- Implement retry logic (/add-retry)
```

### Example 3: Debugging with Agents

**Claude Code Command:**
```
You: "The social repository is throwing errors. Use the agent system to debug."

Claude Code:
1. SearchAgent finds error logs and repository code
2. CodeAgent analyzes error patterns
3. Orchestrator suggests fixes
```

**Debug Report:**
```
🔍 Agent Debug Session

Error Pattern Detected:
- NullPointerException in followUser() at line 45
- Occurs when response.body() is null

Root Cause Analysis:
- Missing null check on API response
- No fallback for network failures

Suggested Fix:
response.body()?.let { body ->
    // process body
} ?: throw NetworkException("Empty response body")

Would you like me to apply this fix? (y/n)
```

## Best Practices

### 1. Clear Task Descriptions
✅ **Good:** "Use the agent system to generate tests for SocialRepository.kt focusing on error cases"
❌ **Bad:** "Make tests"

### 2. Specify Context When Needed
✅ **Good:** "Refactor video service using patterns from auth service"
❌ **Bad:** "Refactor video"

### 3. Leverage Iterative Refinement
✅ **Good:** "Implement feature X, iterate until tests pass"
❌ **Bad:** "Implement feature X" (single attempt)

### 4. Use Parallel Execution for Analysis
✅ **Good:** "Analyze entire social service in parallel"
❌ **Bad:** Sequential analysis of each file

### 5. Trust Context Compaction
- Let the system manage context automatically
- Compaction triggers at 80% capacity
- Removes least-accessed files automatically

## Monitoring Agent Execution

### Check Agent Progress
```
You: "Show me what the agent system is doing"

Claude Code shows:
- Current phase (Gather/Action/Verify)
- Active subagents
- Context utilization
- Iteration count
- Estimated time remaining
```

### View Agent Logs
```
You: "Show agent execution logs"

Output:
[09:15:32] Orchestrator started task-001
[09:15:33] Phase 1: Gather Context
[09:15:35] SearchAgent found 8 files
[09:15:36] Context: 62% utilized
[09:15:37] Phase 2: Take Action
[09:15:38] CodeAgent generating...
[09:15:42] TestAgent generating tests...
[09:15:45] Phase 3: Verify Work
[09:15:46] Validation passed ✅
[09:15:46] Task completed in 1 iteration
```

## Troubleshooting

### Issue: Agent Takes Too Long
**Solution:**
```
You: "Use faster agent execution mode"
Claude Code: Reduces MaxIterations to 3, increases ParallelSubagents
```

### Issue: Context Capacity Exceeded
**Solution:**
```
Agent automatically triggers compaction at 80%
Removes 30% least-accessed files
Continues execution
```

### Issue: Generated Code Doesn't Compile
**Solution:**
```
Agent loop detects compilation failure
Iteration 2 attempts fix with feedback
Max 5 iterations until success or report failure
```

## Integration with Existing Claude Code Features

### Works With All Commands
- `/analyze` - Enhanced with agent-based analysis
- `/improve` - Uses agents for refactoring
- `/test` - Leverages TestAgent
- `/document` - Generates docs using CodeAgent
- `/git` - Can use agents to analyze changes

### Persona System Integration
Claude Code automatically activates appropriate agents based on persona:
- **Architect Persona** → SearchAgent + CodeAgent (analysis focus)
- **QA Persona** → TestAgent (testing focus)
- **Refactorer Persona** → CodeAgent (refactoring focus)

## Summary

The agent system is fully integrated with Claude Code and can be used naturally through conversation:

**Simply ask Claude Code to:**
- "Use agents to find X"
- "Generate Y using the agent system"
- "Analyze Z with agents"
- "Refactor A using agents"

Claude Code will automatically:
1. Create appropriate tasks
2. Orchestrate subagents
3. Manage context
4. Iterate until complete
5. Return results

No manual setup required - just describe what you want, and the agent system handles the complexity!