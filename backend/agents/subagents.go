package agents

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// SearchAgent implements agentic search using file system structure
// Following Claude SDK pattern: "agentic search" through directory traversal
type SearchAgent struct {
	name           string
	contextManager *ContextManager
	searchDepth    int
}

// NewSearchAgent creates a new search agent
func NewSearchAgent(contextManager *ContextManager) *SearchAgent {
	return &SearchAgent{
		name:           "search-agent",
		contextManager: contextManager,
		searchDepth:    5,
	}
}

// Execute implements the agent loop for search operations
func (s *SearchAgent) Execute(ctx context.Context, task *Task) (*AgentResult, error) {
	startTime := time.Now()
	result := &AgentResult{
		TaskID:    task.ID,
		Artifacts: make([]string, 0),
	}

	// Phase 1: Gather Context - Analyze search query
	searchQuery := task.Description
	searchScope := s.determineSearchScope(searchQuery)

	// Phase 2: Take Action - Execute agentic search
	files, err := s.agenticSearch(ctx, searchQuery, searchScope)
	if err != nil {
		result.Errors = append(result.Errors, err)
		return result, err
	}

	// Phase 3: Verify Work - Validate search results
	if len(files) == 0 {
		result.Success = false
		result.Output = "No files found matching search criteria"
		return result, nil
	}

	result.Success = true
	result.Output = files
	result.Artifacts = files
	result.Duration = time.Since(startTime)
	result.Iterations = 1

	return result, nil
}

// agenticSearch implements file system-based search following Claude SDK patterns
func (s *SearchAgent) agenticSearch(ctx context.Context, query string, scope string) ([]string, error) {
	matchedFiles := make([]string, 0)

	// Strategy 1: Directory structure search (primary method per Claude SDK)
	if strings.Contains(query, "test") {
		// Search in test directories
		testDirs := []string{"tests", "__tests__", "test", "spec"}
		for _, dir := range testDirs {
			files, _ := s.searchDirectory(scope, dir)
			matchedFiles = append(matchedFiles, files...)
		}
	}

	if strings.Contains(query, "component") || strings.Contains(query, "ui") {
		// Search in component directories
		componentDirs := []string{"components", "src/components", "ui"}
		for _, dir := range componentDirs {
			files, _ := s.searchDirectory(scope, dir)
			matchedFiles = append(matchedFiles, files...)
		}
	}

	// Strategy 2: Content search using grep (supplementary)
	if len(matchedFiles) == 0 {
		grepResults, _ := s.contentSearch(scope, query)
		matchedFiles = append(matchedFiles, grepResults...)
	}

	return matchedFiles, nil
}

// searchDirectory searches for files in a specific directory
func (s *SearchAgent) searchDirectory(basePath, dir string) ([]string, error) {
	searchPath := filepath.Join(basePath, dir)
	cmd := exec.Command("find", searchPath, "-type", "f", "-maxdepth", fmt.Sprintf("%d", s.searchDepth))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	return files, nil
}

// contentSearch performs grep-based content search
func (s *SearchAgent) contentSearch(basePath, query string) ([]string, error) {
	cmd := exec.Command("grep", "-r", "-l", query, basePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	return files, nil
}

// determineSearchScope analyzes query to determine search scope
func (s *SearchAgent) determineSearchScope(query string) string {
	// Default to current working directory
	scope := "."

	// Intelligent scope determination based on query
	if strings.Contains(query, "backend") || strings.Contains(query, "server") {
		scope = "./backend"
	} else if strings.Contains(query, "frontend") || strings.Contains(query, "web") {
		scope = "./apps/web"
	} else if strings.Contains(query, "mobile") || strings.Contains(query, "kmp") {
		scope = "./apps/kmp"
	}

	return scope
}

func (s *SearchAgent) GetCapabilities() []string {
	return []string{"file-search", "agentic-search", "content-search", "directory-traversal"}
}

func (s *SearchAgent) GetName() string {
	return s.name
}

// ============================================================================
// CodeAgent - Handles code generation and modification
// ============================================================================

type CodeAgent struct {
	name           string
	contextManager *ContextManager
	toolRegistry   *ToolRegistry
}

func NewCodeAgent(contextManager *ContextManager, toolRegistry *ToolRegistry) *CodeAgent {
	return &CodeAgent{
		name:           "code-agent",
		contextManager: contextManager,
		toolRegistry:   toolRegistry,
	}
}

func (c *CodeAgent) Execute(ctx context.Context, task *Task) (*AgentResult, error) {
	startTime := time.Now()
	result := &AgentResult{
		TaskID:    task.ID,
		Artifacts: make([]string, 0),
	}

	maxIterations := 3
	for iteration := 0; iteration < maxIterations; iteration++ {
		// Phase 1: Gather Context - Read relevant files
		relevantFiles, err := c.gatherCodeContext(ctx, task)
		if err != nil {
			result.Errors = append(result.Errors, err)
			continue
		}

		// Phase 2: Take Action - Generate/modify code
		codeOutput, err := c.generateCode(ctx, task, relevantFiles)
		if err != nil {
			result.Errors = append(result.Errors, err)
			continue
		}

		// Phase 3: Verify Work - Validate generated code
		verified, err := c.verifyCode(ctx, codeOutput)
		if err != nil {
			result.Errors = append(result.Errors, err)
			continue
		}

		if verified {
			result.Success = true
			result.Output = codeOutput
			result.Artifacts = codeOutput.Files
			result.Iterations = iteration + 1
			break
		}

		// Continue to next iteration with feedback
		task.Context["previous_code"] = codeOutput
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

func (c *CodeAgent) gatherCodeContext(ctx context.Context, task *Task) ([]string, error) {
	// Use search to find related code files
	searchQuery := task.Description
	files, err := c.contextManager.SearchRelevantFiles(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("context gathering failed: %w", err)
	}

	return files, nil
}

func (c *CodeAgent) generateCode(ctx context.Context, task *Task, relevantFiles []string) (*CodeOutput, error) {
	// Code generation logic following Claude SDK patterns
	// Use bash tools, code generation templates, MCP integration

	output := &CodeOutput{
		Files:   make([]string, 0),
		Content: make(map[string]string),
	}

	// Execute code generation tool
	tool, err := c.toolRegistry.GetTool("code-generator")
	if err != nil {
		return nil, fmt.Errorf("tool not found: %w", err)
	}

	toolResult, err := tool.Execute(ctx, &ToolInput{
		Task:          task,
		RelevantFiles: relevantFiles,
	})
	if err != nil {
		return nil, fmt.Errorf("code generation failed: %w", err)
	}

	output.Files = toolResult.Files
	output.Content = toolResult.Content

	return output, nil
}

func (c *CodeAgent) verifyCode(ctx context.Context, output *CodeOutput) (bool, error) {
	// Rule-based validation
	if len(output.Files) == 0 {
		return false, nil
	}

	// Check for common issues
	for _, content := range output.Content {
		if strings.Contains(content, "TODO") || strings.Contains(content, "FIXME") {
			return false, fmt.Errorf("generated code contains TODO/FIXME markers")
		}
	}

	// Optional: Run linter/formatter
	// This would use bash tools to run linters

	return true, nil
}

func (c *CodeAgent) GetCapabilities() []string {
	return []string{"code-generation", "code-modification", "refactoring", "linting"}
}

func (c *CodeAgent) GetName() string {
	return c.name
}

// ============================================================================
// TestAgent - Handles test generation and execution
// ============================================================================

type TestAgent struct {
	name           string
	contextManager *ContextManager
	toolRegistry   *ToolRegistry
}

func NewTestAgent(contextManager *ContextManager, toolRegistry *ToolRegistry) *TestAgent {
	return &TestAgent{
		name:           "test-agent",
		contextManager: contextManager,
		toolRegistry:   toolRegistry,
	}
}

func (t *TestAgent) Execute(ctx context.Context, task *Task) (*AgentResult, error) {
	startTime := time.Now()
	result := &AgentResult{
		TaskID:    task.ID,
		Artifacts: make([]string, 0),
	}

	// Phase 1: Gather Context - Find code to test
	filesToTest, err := t.identifyFilesToTest(ctx, task)
	if err != nil {
		result.Errors = append(result.Errors, err)
		return result, err
	}

	// Phase 2: Take Action - Generate and run tests
	testResults, err := t.generateAndRunTests(ctx, filesToTest)
	if err != nil {
		result.Errors = append(result.Errors, err)
		return result, err
	}

	// Phase 3: Verify Work - Check test coverage and results
	verified, err := t.verifyTests(ctx, testResults)
	if err != nil {
		result.Errors = append(result.Errors, err)
		return result, err
	}

	result.Success = verified
	result.Output = testResults
	result.Artifacts = testResults.TestFiles
	result.Duration = time.Since(startTime)
	result.Iterations = 1

	return result, nil
}

func (t *TestAgent) identifyFilesToTest(ctx context.Context, task *Task) ([]string, error) {
	// Use search to find source files
	files, err := t.contextManager.SearchRelevantFiles(task.Description)
	if err != nil {
		return nil, err
	}

	// Filter out test files
	sourceFiles := make([]string, 0)
	for _, file := range files {
		if !strings.Contains(file, "test") && !strings.Contains(file, "spec") {
			sourceFiles = append(sourceFiles, file)
		}
	}

	return sourceFiles, nil
}

func (t *TestAgent) generateAndRunTests(ctx context.Context, files []string) (*TestResults, error) {
	results := &TestResults{
		TestFiles: make([]string, 0),
		Passed:    0,
		Failed:    0,
		Coverage:  0.0,
	}

	// Use test generation tool
	tool, err := t.toolRegistry.GetTool("test-generator")
	if err != nil {
		return nil, fmt.Errorf("test tool not found: %w", err)
	}

	for _, file := range files {
		toolResult, err := tool.Execute(ctx, &ToolInput{
			TargetFile: file,
		})
		if err != nil {
			results.Failed++
			continue
		}

		results.TestFiles = append(results.TestFiles, toolResult.TestFile)
		results.Passed++
	}

	// Run tests using bash
	cmd := exec.Command("go", "test", "./...", "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("test execution failed: %w", err)
	}

	results.Output = string(output)

	return results, nil
}

func (t *TestAgent) verifyTests(ctx context.Context, results *TestResults) (bool, error) {
	// Validation rules
	if results.Failed > 0 {
		return false, fmt.Errorf("tests failed: %d", results.Failed)
	}

	if len(results.TestFiles) == 0 {
		return false, fmt.Errorf("no tests generated")
	}

	// Check coverage threshold (optional)
	coverageThreshold := 80.0
	if results.Coverage < coverageThreshold {
		return false, fmt.Errorf("coverage too low: %.2f%% (minimum: %.2f%%)", results.Coverage, coverageThreshold)
	}

	return true, nil
}

func (t *TestAgent) GetCapabilities() []string {
	return []string{"test-generation", "test-execution", "coverage-analysis"}
}

func (t *TestAgent) GetName() string {
	return t.name
}

// Supporting types

type CodeOutput struct {
	Files   []string
	Content map[string]string
}

type TestResults struct {
	TestFiles []string
	Passed    int
	Failed    int
	Coverage  float64
	Output    string
}

type ToolInput struct {
	Task          *Task
	RelevantFiles []string
	TargetFile    string
}

type ToolOutput struct {
	Files    []string
	Content  map[string]string
	TestFile string
}