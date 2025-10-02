package agents

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// OrchestratorAgent coordinates multiple specialized subagents following Claude SDK patterns
// Implements the agent loop: gather context -> take action -> verify work -> repeat
type OrchestratorAgent struct {
	config          *AgentConfig
	subagents       map[string]Agent
	contextManager  *ContextManager
	toolRegistry    *ToolRegistry
	executionLogger *ExecutionLogger
	mu              sync.RWMutex
}

// AgentConfig defines configuration for agents
type AgentConfig struct {
	MaxIterations     int
	ContextLimit      int // Context window size limit
	ParallelSubagents int
	TimeoutSeconds    int
	EnableCompaction  bool
}

// Agent interface defines the contract for all agents
type Agent interface {
	// Execute runs the agent loop: gather -> act -> verify
	Execute(ctx context.Context, task *Task) (*AgentResult, error)

	// GetCapabilities returns what this agent can do
	GetCapabilities() []string

	// GetName returns agent identifier
	GetName() string
}

// Task represents work to be done by an agent
type Task struct {
	ID          string
	Type        string
	Description string
	Context     map[string]interface{}
	Priority    int
	CreatedAt   time.Time
}

// AgentResult contains the outcome of agent execution
type AgentResult struct {
	TaskID      string
	Success     bool
	Output      interface{}
	Artifacts   []string // Files created/modified
	Iterations  int
	Duration    time.Duration
	Errors      []error
	SubResults  []*AgentResult // Results from subagents
	ContextUsed int
}

// AgentLoop represents one iteration of the agent cycle
type AgentLoop struct {
	Iteration int
	Phase     string // "gather", "act", "verify"
	Input     interface{}
	Output    interface{}
	Timestamp time.Time
}

// NewOrchestratorAgent creates a new orchestrator following Claude SDK patterns
func NewOrchestratorAgent(config *AgentConfig) *OrchestratorAgent {
	return &OrchestratorAgent{
		config:          config,
		subagents:       make(map[string]Agent),
		contextManager:  NewContextManager(config.ContextLimit),
		toolRegistry:    NewToolRegistry(),
		executionLogger: NewExecutionLogger(),
	}
}

// RegisterSubagent adds a specialized subagent to the team
func (o *OrchestratorAgent) RegisterSubagent(name string, agent Agent) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.subagents[name] = agent
}

// Execute implements the main agent loop following Claude SDK patterns
func (o *OrchestratorAgent) Execute(ctx context.Context, task *Task) (*AgentResult, error) {
	startTime := time.Now()
	result := &AgentResult{
		TaskID:     task.ID,
		Artifacts:  make([]string, 0),
		SubResults: make([]*AgentResult, 0),
	}

	// Main agent loop: gather context -> take action -> verify work -> repeat
	for iteration := 0; iteration < o.config.MaxIterations; iteration++ {
		loopStart := time.Now()

		// Phase 1: Gather Context
		contextData, err := o.gatherContext(ctx, task, iteration)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("iteration %d gather failed: %w", iteration, err))
			continue
		}

		// Check if context needs compaction
		if o.config.EnableCompaction && o.contextManager.ShouldCompact() {
			if err := o.contextManager.Compact(ctx); err != nil {
				o.executionLogger.LogWarning("Context compaction failed", err)
			}
		}

		// Phase 2: Take Action
		actionResult, err := o.takeAction(ctx, task, contextData)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("iteration %d action failed: %w", iteration, err))
			continue
		}

		// Phase 3: Verify Work
		verified, err := o.verifyWork(ctx, task, actionResult)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("iteration %d verify failed: %w", iteration, err))
			continue
		}

		// Log iteration
		o.executionLogger.LogIteration(&AgentLoop{
			Iteration: iteration,
			Phase:     "complete",
			Input:     contextData,
			Output:    actionResult,
			Timestamp: loopStart,
		})

		// Check if work is complete
		if verified.Complete {
			result.Success = true
			result.Output = verified.Output
			result.Artifacts = append(result.Artifacts, verified.Artifacts...)
			break
		}

		// Continue to next iteration with feedback
		task.Context["previous_iteration"] = verified.Feedback
		result.Iterations = iteration + 1
	}

	result.Duration = time.Since(startTime)
	result.ContextUsed = o.contextManager.GetCurrentSize()

	return result, nil
}

// gatherContext implements Phase 1: Gather Context
func (o *OrchestratorAgent) gatherContext(ctx context.Context, task *Task, iteration int) (*ContextData, error) {
	contextData := &ContextData{
		Task:      task,
		Iteration: iteration,
		Timestamp: time.Now(),
	}

	// Agentic search: use file system structure for context
	files, err := o.contextManager.SearchRelevantFiles(task.Description)
	if err != nil {
		return nil, fmt.Errorf("file search failed: %w", err)
	}
	contextData.RelevantFiles = files

	// Gather from previous iterations
	if prevContext, ok := task.Context["previous_iteration"]; ok {
		contextData.PreviousFeedback = prevContext.(string)
	}

	// Identify which subagents are needed
	requiredAgents := o.identifyRequiredSubagents(task)
	contextData.RequiredAgents = requiredAgents

	return contextData, nil
}

// takeAction implements Phase 2: Take Action
func (o *OrchestratorAgent) takeAction(ctx context.Context, task *Task, contextData *ContextData) (*ActionResult, error) {
	actionResult := &ActionResult{
		Timestamp: time.Now(),
		Actions:   make([]string, 0),
	}

	// Determine if we need parallel subagent execution
	if len(contextData.RequiredAgents) > 1 && o.config.ParallelSubagents > 1 {
		// Parallel subagent execution
		subResults, err := o.executeSubagentsParallel(ctx, task, contextData.RequiredAgents)
		if err != nil {
			return nil, fmt.Errorf("parallel execution failed: %w", err)
		}
		actionResult.SubResults = subResults
	} else {
		// Sequential execution
		for _, agentName := range contextData.RequiredAgents {
			subResult, err := o.executeSubagent(ctx, task, agentName)
			if err != nil {
				o.executionLogger.LogWarning(fmt.Sprintf("Subagent %s failed", agentName), err)
				continue
			}
			actionResult.SubResults = append(actionResult.SubResults, subResult)
		}
	}

	// Execute tools as needed
	for _, tool := range contextData.RequiredTools {
		toolResult, err := o.toolRegistry.Execute(ctx, tool, task)
		if err != nil {
			return nil, fmt.Errorf("tool execution failed: %w", err)
		}
		actionResult.ToolResults = append(actionResult.ToolResults, toolResult)
		actionResult.Actions = append(actionResult.Actions, fmt.Sprintf("tool:%s", tool))
	}

	return actionResult, nil
}

// verifyWork implements Phase 3: Verify Work
func (o *OrchestratorAgent) verifyWork(ctx context.Context, task *Task, actionResult *ActionResult) (*VerificationResult, error) {
	verification := &VerificationResult{
		Timestamp: time.Now(),
		Artifacts: make([]string, 0),
	}

	// Rule-based validation
	rules := o.getValidationRules(task.Type)
	for _, rule := range rules {
		passed, err := rule.Validate(actionResult)
		if err != nil {
			return nil, fmt.Errorf("rule validation error: %w", err)
		}
		if !passed {
			verification.Complete = false
			verification.Feedback = fmt.Sprintf("Validation rule failed: %s", rule.Name)
			return verification, nil
		}
	}

	// Collect artifacts (files created/modified)
	artifacts, err := o.collectArtifacts(actionResult)
	if err != nil {
		return nil, fmt.Errorf("artifact collection failed: %w", err)
	}
	verification.Artifacts = artifacts

	// Check if all subagents succeeded
	allSuccess := true
	for _, subResult := range actionResult.SubResults {
		if !subResult.Success {
			allSuccess = false
			verification.Feedback = fmt.Sprintf("Subagent failed: %s", subResult.TaskID)
			break
		}
	}

	if allSuccess && len(verification.Artifacts) > 0 {
		verification.Complete = true
		verification.Output = actionResult
	}

	return verification, nil
}

// executeSubagentsParallel executes multiple subagents in parallel following Claude SDK patterns
func (o *OrchestratorAgent) executeSubagentsParallel(ctx context.Context, task *Task, agentNames []string) ([]*AgentResult, error) {
	results := make([]*AgentResult, len(agentNames))
	errors := make([]error, len(agentNames))
	var wg sync.WaitGroup

	// Limit concurrency
	semaphore := make(chan struct{}, o.config.ParallelSubagents)

	for i, agentName := range agentNames {
		wg.Add(1)
		go func(idx int, name string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			agent, exists := o.subagents[name]
			if !exists {
				errors[idx] = fmt.Errorf("subagent not found: %s", name)
				return
			}

			// Create isolated context for subagent
			subTask := o.createSubTask(task, name)
			result, err := agent.Execute(ctx, subTask)
			if err != nil {
				errors[idx] = err
				return
			}

			results[idx] = result
		}(i, agentName)
	}

	wg.Wait()

	// Filter out failed subagents
	successfulResults := make([]*AgentResult, 0)
	for i, result := range results {
		if errors[i] == nil && result != nil {
			successfulResults = append(successfulResults, result)
		}
	}

	return successfulResults, nil
}

// executeSubagent executes a single subagent
func (o *OrchestratorAgent) executeSubagent(ctx context.Context, task *Task, agentName string) (*AgentResult, error) {
	agent, exists := o.subagents[agentName]
	if !exists {
		return nil, fmt.Errorf("subagent not found: %s", agentName)
	}

	subTask := o.createSubTask(task, agentName)
	return agent.Execute(ctx, subTask)
}

// identifyRequiredSubagents determines which subagents are needed for the task
func (o *OrchestratorAgent) identifyRequiredSubagents(task *Task) []string {
	requiredAgents := make([]string, 0)

	switch task.Type {
	case "search":
		requiredAgents = append(requiredAgents, "search-agent")
	case "code":
		requiredAgents = append(requiredAgents, "code-agent")
	case "test":
		requiredAgents = append(requiredAgents, "test-agent")
	case "analysis":
		requiredAgents = append(requiredAgents, "search-agent", "code-agent")
	case "full-stack":
		requiredAgents = append(requiredAgents, "search-agent", "code-agent", "test-agent")
	default:
		// Default to all available subagents
		o.mu.RLock()
		for name := range o.subagents {
			requiredAgents = append(requiredAgents, name)
		}
		o.mu.RUnlock()
	}

	return requiredAgents
}

// createSubTask creates an isolated task for a subagent
func (o *OrchestratorAgent) createSubTask(parentTask *Task, agentName string) *Task {
	return &Task{
		ID:          fmt.Sprintf("%s-%s", parentTask.ID, agentName),
		Type:        agentName,
		Description: parentTask.Description,
		Context: map[string]interface{}{
			"parent_task": parentTask.ID,
			"agent_name":  agentName,
		},
		Priority:  parentTask.Priority,
		CreatedAt: time.Now(),
	}
}

// getValidationRules returns validation rules for a task type
func (o *OrchestratorAgent) getValidationRules(taskType string) []*ValidationRule {
	// Return task-specific validation rules
	rules := make([]*ValidationRule, 0)

	// Common rules
	rules = append(rules, &ValidationRule{
		Name: "artifacts-exist",
		Validate: func(result *ActionResult) (bool, error) {
			return len(result.ToolResults) > 0 || len(result.SubResults) > 0, nil
		},
	})

	return rules
}

// collectArtifacts collects all files created or modified during action
func (o *OrchestratorAgent) collectArtifacts(actionResult *ActionResult) ([]string, error) {
	artifacts := make([]string, 0)

	// Collect from tool results
	for _, toolResult := range actionResult.ToolResults {
		if files, ok := toolResult.Output.([]string); ok {
			artifacts = append(artifacts, files...)
		}
	}

	// Collect from subagent results
	for _, subResult := range actionResult.SubResults {
		artifacts = append(artifacts, subResult.Artifacts...)
	}

	return artifacts, nil
}

// GetCapabilities returns orchestrator capabilities
func (o *OrchestratorAgent) GetCapabilities() []string {
	return []string{
		"subagent-coordination",
		"parallel-execution",
		"context-management",
		"iterative-refinement",
	}
}

// GetName returns agent identifier
func (o *OrchestratorAgent) GetName() string {
	return "orchestrator"
}

// Supporting types

type ContextData struct {
	Task             *Task
	Iteration        int
	Timestamp        time.Time
	RelevantFiles    []string
	PreviousFeedback string
	RequiredAgents   []string
	RequiredTools    []string
}

type ActionResult struct {
	Timestamp   time.Time
	Actions     []string
	SubResults  []*AgentResult
	ToolResults []*ToolResult
}

type VerificationResult struct {
	Complete  bool
	Output    interface{}
	Artifacts []string
	Feedback  string
	Timestamp time.Time
}

type ValidationRule struct {
	Name     string
	Validate func(*ActionResult) (bool, error)
}

type ToolResult struct {
	ToolName string
	Success  bool
	Output   interface{}
	Error    error
}