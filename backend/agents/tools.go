package agents

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// Tool interface defines the contract for all agent tools
// Following Claude SDK pattern: "tools are the primary building blocks of execution"
type Tool interface {
	// Execute runs the tool with given input
	Execute(ctx context.Context, input *ToolInput) (*ToolOutput, error)

	// GetName returns tool identifier
	GetName() string

	// GetDescription returns what this tool does
	GetDescription() string
}

// ToolRegistry manages available tools for agents
type ToolRegistry struct {
	tools map[string]Tool
	mu    sync.RWMutex
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	registry := &ToolRegistry{
		tools: make(map[string]Tool),
	}

	// Register default tools
	registry.RegisterTool(NewBashTool())
	registry.RegisterTool(NewReadFileTool())
	registry.RegisterTool(NewWriteFileTool())
	registry.RegisterTool(NewCodeGeneratorTool())
	registry.RegisterTool(NewTestGeneratorTool())

	return registry
}

// RegisterTool adds a tool to the registry
func (tr *ToolRegistry) RegisterTool(tool Tool) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.tools[tool.GetName()] = tool
}

// GetTool retrieves a tool by name
func (tr *ToolRegistry) GetTool(name string) (Tool, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	tool, exists := tr.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	return tool, nil
}

// Execute runs a tool by name
func (tr *ToolRegistry) Execute(ctx context.Context, toolName string, task *Task) (*ToolResult, error) {
	tool, err := tr.GetTool(toolName)
	if err != nil {
		return nil, err
	}

	input := &ToolInput{
		Task: task,
	}

	output, err := tool.Execute(ctx, input)
	if err != nil {
		return &ToolResult{
			ToolName: toolName,
			Success:  false,
			Error:    err,
		}, err
	}

	return &ToolResult{
		ToolName: toolName,
		Success:  true,
		Output:   output,
	}, nil
}

// ============================================================================
// BashTool - Execute shell commands
// ============================================================================

type BashTool struct {
	name        string
	description string
}

func NewBashTool() *BashTool {
	return &BashTool{
		name:        "bash",
		description: "Execute shell commands",
	}
}

func (b *BashTool) GetName() string {
	return b.name
}

func (b *BashTool) GetDescription() string {
	return b.description
}

func (b *BashTool) Execute(ctx context.Context, input *ToolInput) (*ToolOutput, error) {
	if input.Command == "" {
		return nil, fmt.Errorf("bash command required")
	}

	cmd := exec.CommandContext(ctx, "bash", "-c", input.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("bash execution failed: %w, output: %s", err, string(output))
	}

	return &ToolOutput{
		Content: map[string]string{
			"stdout": string(output),
		},
	}, nil
}

// ============================================================================
// ReadFileTool - Read file contents
// ============================================================================

type ReadFileTool struct {
	name        string
	description string
}

func NewReadFileTool() *ReadFileTool {
	return &ReadFileTool{
		name:        "read-file",
		description: "Read contents of a file",
	}
}

func (r *ReadFileTool) GetName() string {
	return r.name
}

func (r *ReadFileTool) GetDescription() string {
	return r.description
}

func (r *ReadFileTool) Execute(ctx context.Context, input *ToolInput) (*ToolOutput, error) {
	if input.TargetFile == "" {
		return nil, fmt.Errorf("target file required")
	}

	content, err := os.ReadFile(input.TargetFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return &ToolOutput{
		Files: []string{input.TargetFile},
		Content: map[string]string{
			input.TargetFile: string(content),
		},
	}, nil
}

// ============================================================================
// WriteFileTool - Write file contents
// ============================================================================

type WriteFileTool struct {
	name        string
	description string
}

func NewWriteFileTool() *WriteFileTool {
	return &WriteFileTool{
		name:        "write-file",
		description: "Write contents to a file",
	}
}

func (w *WriteFileTool) GetName() string {
	return w.name
}

func (w *WriteFileTool) GetDescription() string {
	return w.description
}

func (w *WriteFileTool) Execute(ctx context.Context, input *ToolInput) (*ToolOutput, error) {
	if input.TargetFile == "" || input.FileContent == "" {
		return nil, fmt.Errorf("target file and content required")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(input.TargetFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(input.TargetFile, []byte(input.FileContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &ToolOutput{
		Files: []string{input.TargetFile},
		Content: map[string]string{
			input.TargetFile: input.FileContent,
		},
	}, nil
}

// ============================================================================
// CodeGeneratorTool - Generate code following Claude SDK patterns
// ============================================================================

type CodeGeneratorTool struct {
	name        string
	description string
}

func NewCodeGeneratorTool() *CodeGeneratorTool {
	return &CodeGeneratorTool{
		name:        "code-generator",
		description: "Generate code based on task description and context",
	}
}

func (c *CodeGeneratorTool) GetName() string {
	return c.name
}

func (c *CodeGeneratorTool) GetDescription() string {
	return c.description
}

func (c *CodeGeneratorTool) Execute(ctx context.Context, input *ToolInput) (*ToolOutput, error) {
	if input.Task == nil {
		return nil, fmt.Errorf("task required for code generation")
	}

	// Extract code generation requirements from task
	description := input.Task.Description
	relevantFiles := input.RelevantFiles

	// Analyze relevant files for context
	contextCode := make(map[string]string)
	for _, file := range relevantFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue // Skip files that can't be read
		}
		contextCode[file] = string(content)
	}

	// Generate code structure based on task description
	generatedFiles := make(map[string]string)

	// Determine target language and structure
	language := c.detectLanguage(description)
	targetFile := c.determineTargetFile(description, language)

	// Generate code content
	codeContent := c.generateCodeContent(description, language, contextCode)

	generatedFiles[targetFile] = codeContent

	return &ToolOutput{
		Files:   []string{targetFile},
		Content: generatedFiles,
	}, nil
}

func (c *CodeGeneratorTool) detectLanguage(description string) string {
	description = strings.ToLower(description)

	if strings.Contains(description, "go") || strings.Contains(description, "golang") {
		return "go"
	}
	if strings.Contains(description, "typescript") || strings.Contains(description, "tsx") {
		return "typescript"
	}
	if strings.Contains(description, "javascript") || strings.Contains(description, "jsx") {
		return "javascript"
	}
	if strings.Contains(description, "python") || strings.Contains(description, ".py") {
		return "python"
	}

	// Default to Go for backend tasks
	return "go"
}

func (c *CodeGeneratorTool) determineTargetFile(description string, language string) string {
	description = strings.ToLower(description)

	// Extract potential file name from description
	words := strings.Fields(description)
	var fileName string

	for i, word := range words {
		if word == "create" || word == "generate" || word == "implement" {
			if i+1 < len(words) {
				fileName = words[i+1]
				break
			}
		}
	}

	if fileName == "" {
		fileName = "generated"
	}

	// Add appropriate extension
	switch language {
	case "go":
		return fmt.Sprintf("backend/%s.go", fileName)
	case "typescript":
		return fmt.Sprintf("apps/web/src/%s.ts", fileName)
	case "javascript":
		return fmt.Sprintf("apps/web/src/%s.js", fileName)
	case "python":
		return fmt.Sprintf("scripts/%s.py", fileName)
	default:
		return fmt.Sprintf("generated/%s.txt", fileName)
	}
}

func (c *CodeGeneratorTool) generateCodeContent(description string, language string, contextCode map[string]string) string {
	// Basic code generation based on language
	switch language {
	case "go":
		return c.generateGoCode(description, contextCode)
	case "typescript", "javascript":
		return c.generateTypeScriptCode(description, contextCode)
	case "python":
		return c.generatePythonCode(description, contextCode)
	default:
		return "// Generated code placeholder\n"
	}
}

func (c *CodeGeneratorTool) generateGoCode(description string, contextCode map[string]string) string {
	// Basic Go code template
	return `package main

import (
	"context"
	"fmt"
)

// Generated function based on task description
func GeneratedFunction(ctx context.Context) error {
	// TODO: Implement based on: ` + description + `
	return nil
}

func main() {
	if err := GeneratedFunction(context.Background()); err != nil {
		fmt.Println("Error:", err)
	}
}
`
}

func (c *CodeGeneratorTool) generateTypeScriptCode(description string, contextCode map[string]string) string {
	return `// Generated TypeScript code
// Task: ` + description + `

export function generatedFunction(): void {
  // TODO: Implement based on task description
}
`
}

func (c *CodeGeneratorTool) generatePythonCode(description string, contextCode map[string]string) string {
	return `# Generated Python code
# Task: ` + description + `

def generated_function():
    # TODO: Implement based on task description
    pass

if __name__ == "__main__":
    generated_function()
`
}

// ============================================================================
// TestGeneratorTool - Generate tests for code
// ============================================================================

type TestGeneratorTool struct {
	name        string
	description string
}

func NewTestGeneratorTool() *TestGeneratorTool {
	return &TestGeneratorTool{
		name:        "test-generator",
		description: "Generate tests for code files",
	}
}

func (t *TestGeneratorTool) GetName() string {
	return t.name
}

func (t *TestGeneratorTool) GetDescription() string {
	return t.description
}

func (t *TestGeneratorTool) Execute(ctx context.Context, input *ToolInput) (*ToolOutput, error) {
	if input.TargetFile == "" {
		return nil, fmt.Errorf("target file required for test generation")
	}

	// Read target file
	content, err := os.ReadFile(input.TargetFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read target file: %w", err)
	}

	// Determine test file path
	testFile := t.generateTestFilePath(input.TargetFile)

	// Generate test content
	testContent := t.generateTestContent(input.TargetFile, string(content))

	return &ToolOutput{
		Files:    []string{testFile},
		TestFile: testFile,
		Content: map[string]string{
			testFile: testContent,
		},
	}, nil
}

func (t *TestGeneratorTool) generateTestFilePath(targetFile string) string {
	ext := filepath.Ext(targetFile)
	base := strings.TrimSuffix(targetFile, ext)

	switch ext {
	case ".go":
		return base + "_test.go"
	case ".ts", ".tsx":
		return strings.Replace(targetFile, "/src/", "/tests/", 1)
	case ".js", ".jsx":
		return base + ".test.js"
	case ".py":
		return base + "_test.py"
	default:
		return base + ".test" + ext
	}
}

func (t *TestGeneratorTool) generateTestContent(targetFile string, sourceCode string) string {
	ext := filepath.Ext(targetFile)

	switch ext {
	case ".go":
		return t.generateGoTest(targetFile, sourceCode)
	case ".ts", ".tsx", ".js", ".jsx":
		return t.generateJavaScriptTest(targetFile, sourceCode)
	case ".py":
		return t.generatePythonTest(targetFile, sourceCode)
	default:
		return "// Generated test placeholder\n"
	}
}

func (t *TestGeneratorTool) generateGoTest(targetFile string, sourceCode string) string {
	packageName := "main" // Default, should be extracted from source
	return `package ` + packageName + `

import (
	"context"
	"testing"
)

func TestGeneratedFunction(t *testing.T) {
	ctx := context.Background()

	// TODO: Add test cases based on source code
	err := GeneratedFunction(ctx)
	if err != nil {
		t.Errorf("GeneratedFunction() error = %v", err)
	}
}
`
}

func (t *TestGeneratorTool) generateJavaScriptTest(targetFile string, sourceCode string) string {
	return `import { describe, it, expect } from 'vitest';
import { generatedFunction } from './` + filepath.Base(targetFile) + `';

describe('generatedFunction', () => {
  it('should work correctly', () => {
    // TODO: Add test cases
    expect(generatedFunction).toBeDefined();
  });
});
`
}

func (t *TestGeneratorTool) generatePythonTest(targetFile string, sourceCode string) string {
	return `import unittest
from ` + strings.TrimSuffix(filepath.Base(targetFile), ".py") + ` import generated_function

class TestGeneratedFunction(unittest.TestCase):
    def test_generated_function(self):
        # TODO: Add test cases
        result = generated_function()
        self.assertIsNotNone(result)

if __name__ == '__main__':
    unittest.main()
`
}

// ExecutionLogger logs agent execution for debugging and monitoring
type ExecutionLogger struct {
	logs []string
	mu   sync.RWMutex
}

func NewExecutionLogger() *ExecutionLogger {
	return &ExecutionLogger{
		logs: make([]string, 0),
	}
}

func (e *ExecutionLogger) LogIteration(loop *AgentLoop) {
	e.mu.Lock()
	defer e.mu.Unlock()

	logEntry := fmt.Sprintf("[Iteration %d] Phase: %s, Time: %s",
		loop.Iteration, loop.Phase, loop.Timestamp.Format("15:04:05"))
	e.logs = append(e.logs, logEntry)
}

func (e *ExecutionLogger) LogWarning(message string, err error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	logEntry := fmt.Sprintf("[WARNING] %s: %v", message, err)
	e.logs = append(e.logs, logEntry)
}

func (e *ExecutionLogger) GetLogs() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return append([]string{}, e.logs...)
}