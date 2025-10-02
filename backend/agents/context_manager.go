package agents

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// ContextManager manages agent context following Claude SDK patterns
// Implements context engineering through file system organization and automatic compaction
type ContextManager struct {
	contextLimit    int
	currentSize     int
	fileCache       map[string]*FileContext
	searchIndex     map[string][]string // keyword -> file paths
	mu              sync.RWMutex
	compactionRatio float64 // When to trigger compaction (0.8 = 80% full)
}

// FileContext represents cached file information
type FileContext struct {
	Path         string
	Content      string
	Size         int
	LastAccessed int64
	AccessCount  int
	Relevance    float64
}

// NewContextManager creates a new context manager
func NewContextManager(contextLimit int) *ContextManager {
	return &ContextManager{
		contextLimit:    contextLimit,
		fileCache:       make(map[string]*FileContext),
		searchIndex:     make(map[string][]string),
		compactionRatio: 0.8,
	}
}

// SearchRelevantFiles implements agentic search using file system structure
// This is the primary search method per Claude SDK
func (cm *ContextManager) SearchRelevantFiles(query string) ([]string, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	matchedFiles := make([]string, 0)
	keywords := cm.extractKeywords(query)

	// Strategy 1: Use file system structure (primary method)
	structuralMatches := cm.searchByStructure(keywords)
	matchedFiles = append(matchedFiles, structuralMatches...)

	// Strategy 2: Use search index (supplementary)
	if len(matchedFiles) < 5 {
		indexedMatches := cm.searchByIndex(keywords)
		matchedFiles = append(matchedFiles, indexedMatches...)
	}

	// Deduplicate and rank by relevance
	uniqueFiles := cm.deduplicateAndRank(matchedFiles, keywords)

	return uniqueFiles, nil
}

// searchByStructure implements agentic search through directory structure
func (cm *ContextManager) searchByStructure(keywords []string) []string {
	matchedFiles := make([]string, 0)

	// Map keywords to likely directory structures
	for _, keyword := range keywords {
		var searchPaths []string

		switch {
		case strings.Contains(keyword, "test"):
			searchPaths = []string{"tests", "__tests__", "test", "spec"}
		case strings.Contains(keyword, "component"):
			searchPaths = []string{"components", "src/components", "ui", "views"}
		case strings.Contains(keyword, "service"):
			searchPaths = []string{"services", "src/services", "backend/services"}
		case strings.Contains(keyword, "model"):
			searchPaths = []string{"models", "src/models", "backend/models"}
		case strings.Contains(keyword, "util"):
			searchPaths = []string{"utils", "src/utils", "lib"}
		case strings.Contains(keyword, "api"):
			searchPaths = []string{"api", "src/api", "backend/api", "handlers"}
		case strings.Contains(keyword, "config"):
			searchPaths = []string{"config", "src/config"}
		default:
			// Try to find files with keyword in path
			searchPaths = []string{"."}
		}

		// Search in identified paths
		for _, path := range searchPaths {
			files, _ := cm.findFilesInPath(path, keyword)
			matchedFiles = append(matchedFiles, files...)
		}
	}

	return matchedFiles
}

// searchByIndex uses pre-built search index
func (cm *ContextManager) searchByIndex(keywords []string) []string {
	matchedFiles := make([]string, 0)

	for _, keyword := range keywords {
		if files, exists := cm.searchIndex[keyword]; exists {
			matchedFiles = append(matchedFiles, files...)
		}
	}

	return matchedFiles
}

// findFilesInPath searches for files in a specific path
func (cm *ContextManager) findFilesInPath(basePath, keyword string) ([]string, error) {
	files := make([]string, 0)

	// Walk directory
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on error
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Match by filename or path
		if strings.Contains(strings.ToLower(path), strings.ToLower(keyword)) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// deduplicateAndRank removes duplicates and ranks by relevance
func (cm *ContextManager) deduplicateAndRank(files []string, keywords []string) []string {
	// Deduplicate
	seen := make(map[string]bool)
	unique := make([]string, 0)

	for _, file := range files {
		if !seen[file] {
			seen[file] = true
			unique = append(unique, file)
		}
	}

	// Calculate relevance scores
	type scoredFile struct {
		path  string
		score float64
	}

	scored := make([]scoredFile, 0, len(unique))
	for _, file := range unique {
		score := cm.calculateRelevance(file, keywords)
		scored = append(scored, scoredFile{path: file, score: score})
	}

	// Sort by score (descending)
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Return top results
	maxResults := 20
	result := make([]string, 0, maxResults)
	for i, sf := range scored {
		if i >= maxResults {
			break
		}
		result = append(result, sf.path)
	}

	return result
}

// calculateRelevance scores file relevance based on keywords
func (cm *ContextManager) calculateRelevance(filePath string, keywords []string) float64 {
	score := 0.0

	// Check cached relevance
	if fc, exists := cm.fileCache[filePath]; exists {
		score += fc.Relevance
	}

	// Score based on keyword matches in path
	lowerPath := strings.ToLower(filePath)
	for _, keyword := range keywords {
		lowerKeyword := strings.ToLower(keyword)
		if strings.Contains(lowerPath, lowerKeyword) {
			score += 10.0

			// Bonus for exact filename match
			filename := filepath.Base(lowerPath)
			if strings.Contains(filename, lowerKeyword) {
				score += 20.0
			}
		}
	}

	// Bonus for recently accessed files
	if fc, exists := cm.fileCache[filePath]; exists {
		score += float64(fc.AccessCount) * 2.0
	}

	// Penalty for deeply nested files
	depth := strings.Count(filePath, string(filepath.Separator))
	score -= float64(depth) * 0.5

	return score
}

// extractKeywords extracts searchable keywords from query
func (cm *ContextManager) extractKeywords(query string) []string {
	// Simple tokenization
	words := strings.Fields(strings.ToLower(query))

	// Filter stop words
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true,
	}

	keywords := make([]string, 0)
	for _, word := range words {
		if !stopWords[word] && len(word) > 2 {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

// CacheFile adds a file to the context cache
func (cm *ContextManager) CacheFile(path string, content string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	size := len(content)

	// Check if adding this file would exceed limit
	if cm.currentSize+size > cm.contextLimit {
		// Attempt compaction first
		if err := cm.compactInternal(); err != nil {
			return fmt.Errorf("context full and compaction failed: %w", err)
		}

		// Check again after compaction
		if cm.currentSize+size > cm.contextLimit {
			return fmt.Errorf("context limit exceeded: %d + %d > %d", cm.currentSize, size, cm.contextLimit)
		}
	}

	// Cache the file
	fc := &FileContext{
		Path:         path,
		Content:      content,
		Size:         size,
		LastAccessed: 0, // Will be set on access
		AccessCount:  0,
		Relevance:    0.0,
	}

	cm.fileCache[path] = fc
	cm.currentSize += size

	// Update search index
	cm.updateSearchIndex(path, content)

	return nil
}

// updateSearchIndex indexes file for keyword search
func (cm *ContextManager) updateSearchIndex(path string, content string) {
	// Extract keywords from content
	words := strings.Fields(strings.ToLower(content))

	// Index by keywords
	for _, word := range words {
		if len(word) > 3 { // Only index meaningful words
			if _, exists := cm.searchIndex[word]; !exists {
				cm.searchIndex[word] = make([]string, 0)
			}
			cm.searchIndex[word] = append(cm.searchIndex[word], path)
		}
	}
}

// ShouldCompact checks if context should be compacted
func (cm *ContextManager) ShouldCompact() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	utilization := float64(cm.currentSize) / float64(cm.contextLimit)
	return utilization >= cm.compactionRatio
}

// Compact implements automatic context compaction following Claude SDK patterns
func (cm *ContextManager) Compact(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	return cm.compactInternal()
}

// compactInternal performs the actual compaction (must be called with lock held)
func (cm *ContextManager) compactInternal() error {
	// Strategy: Remove least recently used and least relevant files

	// Convert cache to slice for sorting
	type cacheEntry struct {
		path    string
		context *FileContext
		score   float64
	}

	entries := make([]cacheEntry, 0, len(cm.fileCache))
	for path, fc := range cm.fileCache {
		// Calculate retention score (higher = keep)
		score := float64(fc.AccessCount)*10.0 + fc.Relevance - float64(fc.Size)/1000.0
		entries = append(entries, cacheEntry{
			path:    path,
			context: fc,
			score:   score,
		})
	}

	// Sort by score (ascending - lowest scores removed first)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].score < entries[j].score
	})

	// Remove bottom 30% of files
	targetSize := int(float64(cm.contextLimit) * 0.7)
	removedSize := 0

	for _, entry := range entries {
		if cm.currentSize-removedSize <= targetSize {
			break
		}

		// Remove from cache
		delete(cm.fileCache, entry.path)
		removedSize += entry.context.Size

		// Remove from search index
		for keyword := range cm.searchIndex {
			files := cm.searchIndex[keyword]
			filtered := make([]string, 0)
			for _, f := range files {
				if f != entry.path {
					filtered = append(filtered, f)
				}
			}
			cm.searchIndex[keyword] = filtered
		}
	}

	cm.currentSize -= removedSize

	return nil
}

// GetCurrentSize returns current context size
func (cm *ContextManager) GetCurrentSize() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.currentSize
}

// GetUtilization returns context utilization percentage
func (cm *ContextManager) GetUtilization() float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return float64(cm.currentSize) / float64(cm.contextLimit) * 100.0
}

// GetCachedFile retrieves a cached file and updates access stats
func (cm *ContextManager) GetCachedFile(path string) (*FileContext, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	fc, exists := cm.fileCache[path]
	if exists {
		fc.AccessCount++
		fc.LastAccessed = nowMillis()
	}

	return fc, exists
}

// BuildIndex builds the search index from file system
func (cm *ContextManager) BuildIndex(rootPath string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Walk file system and build index
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on error
		}

		// Skip directories and binary files
		if info.IsDir() {
			return nil
		}

		// Read file content (for text files only)
		if cm.isTextFile(path) {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil // Continue on error
			}

			cm.updateSearchIndex(path, string(content))
		}

		return nil
	})

	return err
}

// isTextFile checks if file is text (simple heuristic)
func (cm *ContextManager) isTextFile(path string) bool {
	ext := filepath.Ext(path)
	textExts := map[string]bool{
		".go":   true,
		".ts":   true,
		".tsx":  true,
		".js":   true,
		".jsx":  true,
		".py":   true,
		".java": true,
		".kt":   true,
		".swift": true,
		".md":   true,
		".txt":  true,
		".yaml": true,
		".yml":  true,
		".json": true,
	}

	return textExts[ext]
}

// Helper function to get current time in milliseconds
func nowMillis() int64 {
	return 0 // Placeholder - would use time.Now().UnixMilli() in real implementation
}