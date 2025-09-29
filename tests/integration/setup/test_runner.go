package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// TestRunnerConfig holds configuration for the test runner
type TestRunnerConfig struct {
	TestPattern     string
	TestTimeout     time.Duration
	Parallel        bool
	MaxParallel     int
	Verbose         bool
	SetupDocker     bool
	CleanupAfter    bool
	FixturesDir     string
	MigrationsDir   string
	DatabaseURL     string
	CoverageOutput  string
	TestSuites      []string
}

// TestSuite represents a test suite configuration
type TestSuite struct {
	Name        string
	Path        string
	Description string
	Dependencies []string
	Timeout     time.Duration
}

// Available test suites
var testSuites = []TestSuite{
	{
		Name:        "backend-integration",
		Path:        "../backend",
		Description: "Backend integration tests for all commerce endpoints",
		Dependencies: []string{"postgres", "redis"},
		Timeout:     10 * time.Minute,
	},
	{
		Name:        "frontend-integration",
		Path:        "../frontend",
		Description: "Frontend integration tests for RTK Query and KMP",
		Dependencies: []string{"postgres"},
		Timeout:     5 * time.Minute,
	},
	{
		Name:        "cross-platform",
		Path:        "../cross-platform",
		Description: "Cross-platform synchronization and consistency tests",
		Dependencies: []string{"postgres", "redis", "kafka"},
		Timeout:     15 * time.Minute,
	},
	{
		Name:        "performance",
		Path:        "../performance",
		Description: "Performance and load testing suite",
		Dependencies: []string{"postgres", "redis", "kafka"},
		Timeout:     30 * time.Minute,
	},
	{
		Name:        "all",
		Path:        "..",
		Description: "Run all test suites",
		Dependencies: []string{"postgres", "redis", "kafka", "scylla", "minio"},
		Timeout:     60 * time.Minute,
	},
}

func main() {
	config := parseFlags()

	log.Printf("🚀 Starting Tchat Integration Test Runner")
	log.Printf("Configuration: %+v", config)

	// Setup test environment if requested
	if config.SetupDocker {
		if err := setupDockerEnvironment(); err != nil {
			log.Fatalf("Failed to setup Docker environment: %v", err)
		}
	}

	// Run tests
	if err := runTests(config); err != nil {
		log.Fatalf("Tests failed: %v", err)
	}

	// Cleanup if requested
	if config.CleanupAfter {
		if err := cleanupEnvironment(); err != nil {
			log.Printf("Warning: Failed to cleanup environment: %v", err)
		}
	}

	log.Printf("✅ All tests completed successfully!")
}

func parseFlags() TestRunnerConfig {
	var config TestRunnerConfig

	flag.StringVar(&config.TestPattern, "pattern", "", "Test pattern to run (e.g., TestCart*)")
	flag.DurationVar(&config.TestTimeout, "timeout", 30*time.Minute, "Test timeout duration")
	flag.BoolVar(&config.Parallel, "parallel", true, "Run tests in parallel")
	flag.IntVar(&config.MaxParallel, "max-parallel", 4, "Maximum parallel test processes")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose output")
	flag.BoolVar(&config.SetupDocker, "setup", false, "Setup Docker test environment")
	flag.BoolVar(&config.CleanupAfter, "cleanup", false, "Cleanup environment after tests")
	flag.StringVar(&config.FixturesDir, "fixtures", "../fixtures", "Path to test fixtures")
	flag.StringVar(&config.MigrationsDir, "migrations", "../../../backend/migrations", "Path to database migrations")
	flag.StringVar(&config.DatabaseURL, "db", "", "Database URL (default: from environment)")
	flag.StringVar(&config.CoverageOutput, "coverage", "", "Coverage output file")

	var suites string
	flag.StringVar(&suites, "suites", "all", "Test suites to run (comma-separated)")

	flag.Parse()

	// Parse test suites
	if suites != "" {
		config.TestSuites = strings.Split(suites, ",")
	} else {
		config.TestSuites = []string{"all"}
	}

	// Set default database URL if not provided
	if config.DatabaseURL == "" {
		config.DatabaseURL = os.Getenv("TEST_DATABASE_URL")
		if config.DatabaseURL == "" {
			config.DatabaseURL = "postgres://tchat_test:tchat_test_password@localhost:5433/tchat_test?sslmode=disable"
		}
	}

	return config
}

func setupDockerEnvironment() error {
	log.Printf("🐳 Setting up Docker test environment...")

	// Check if Docker is available
	if err := exec.Command("docker", "--version").Run(); err != nil {
		return fmt.Errorf("Docker is not available: %w", err)
	}

	// Check if docker-compose is available
	if err := exec.Command("docker-compose", "--version").Run(); err != nil {
		return fmt.Errorf("docker-compose is not available: %w", err)
	}

	// Start Docker services
	cmd := exec.Command("docker-compose", "-f", "docker-compose.test.yml", "up", "-d", "--build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Docker services: %w", err)
	}

	// Wait for services to be ready
	log.Printf("⏳ Waiting for services to be ready...")
	time.Sleep(30 * time.Second)

	// Validate services
	if err := validateServices(); err != nil {
		return fmt.Errorf("service validation failed: %w", err)
	}

	log.Printf("✅ Docker environment setup completed!")
	return nil
}

func validateServices() error {
	services := []struct {
		name    string
		command []string
	}{
		{
			name:    "PostgreSQL",
			command: []string{"docker", "exec", "tchat-postgres-test", "pg_isready", "-U", "tchat_test", "-d", "tchat_test"},
		},
		{
			name:    "Redis",
			command: []string{"docker", "exec", "tchat-redis-test", "redis-cli", "-a", "tchat_test_redis_password", "ping"},
		},
		{
			name:    "Kafka",
			command: []string{"docker", "exec", "tchat-kafka-test", "kafka-topics", "--bootstrap-server", "localhost:9092", "--list"},
		},
	}

	for _, service := range services {
		log.Printf("Validating %s...", service.name)

		for i := 0; i < 30; i++ {
			if err := exec.Command(service.command[0], service.command[1:]...).Run(); err == nil {
				log.Printf("✅ %s is ready", service.name)
				break
			}

			if i == 29 {
				return fmt.Errorf("%s failed to become ready", service.name)
			}

			time.Sleep(2 * time.Second)
		}
	}

	return nil
}

func runTests(config TestRunnerConfig) error {
	log.Printf("🧪 Running test suites: %v", config.TestSuites)

	for _, suiteName := range config.TestSuites {
		suite := findTestSuite(suiteName)
		if suite == nil {
			return fmt.Errorf("unknown test suite: %s", suiteName)
		}

		log.Printf("🏃 Running test suite: %s", suite.Name)
		log.Printf("📝 Description: %s", suite.Description)

		if err := runTestSuite(suite, config); err != nil {
			return fmt.Errorf("test suite '%s' failed: %w", suite.Name, err)
		}

		log.Printf("✅ Test suite '%s' completed successfully", suite.Name)
	}

	return nil
}

func findTestSuite(name string) *TestSuite {
	for _, suite := range testSuites {
		if suite.Name == name {
			return &suite
		}
	}
	return nil
}

func runTestSuite(suite *TestSuite, config TestRunnerConfig) error {
	// Build go test command
	args := []string{"test"}

	// Add test path
	testPath, err := filepath.Abs(suite.Path)
	if err != nil {
		return fmt.Errorf("failed to resolve test path: %w", err)
	}
	args = append(args, testPath+"/...")

	// Add test pattern if specified
	if config.TestPattern != "" {
		args = append(args, "-run", config.TestPattern)
	}

	// Add timeout
	timeout := config.TestTimeout
	if suite.Timeout > 0 && suite.Timeout < timeout {
		timeout = suite.Timeout
	}
	args = append(args, "-timeout", timeout.String())

	// Add verbose flag
	if config.Verbose {
		args = append(args, "-v")
	}

	// Add parallel execution
	if config.Parallel {
		args = append(args, "-parallel", fmt.Sprintf("%d", config.MaxParallel))
	}

	// Add coverage if specified
	if config.CoverageOutput != "" {
		coverageFile := fmt.Sprintf("%s_%s.out", config.CoverageOutput, suite.Name)
		args = append(args, "-coverprofile", coverageFile)
	}

	// Set environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("TEST_DATABASE_URL=%s", config.DatabaseURL))
	env = append(env, fmt.Sprintf("FIXTURES_DIR=%s", config.FixturesDir))
	env = append(env, fmt.Sprintf("MIGRATIONS_DIR=%s", config.MigrationsDir))

	// Create and run command
	cmd := exec.Command("go", args...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("🔧 Running command: go %s", strings.Join(args, " "))

	start := time.Now()
	err = cmd.Run()
	duration := time.Since(start)

	if err != nil {
		return fmt.Errorf("test execution failed after %v: %w", duration, err)
	}

	log.Printf("⏱️  Test suite completed in %v", duration)
	return nil
}

func cleanupEnvironment() error {
	log.Printf("🧹 Cleaning up test environment...")

	// Stop and remove Docker containers
	cmd := exec.Command("docker-compose", "-f", "docker-compose.test.yml", "down", "-v", "--remove-orphans")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to cleanup Docker environment: %w", err)
	}

	log.Printf("✅ Environment cleanup completed!")
	return nil
}

func printUsage() {
	fmt.Println("Tchat Integration Test Runner")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Printf("  %s [flags]\n", os.Args[0])
	fmt.Println("")
	fmt.Println("Available test suites:")
	for _, suite := range testSuites {
		fmt.Printf("  %-20s %s\n", suite.Name, suite.Description)
	}
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Printf("  %s -setup -suites=backend-integration -v\n", os.Args[0])
	fmt.Printf("  %s -pattern=TestCart* -suites=backend-integration\n", os.Args[0])
	fmt.Printf("  %s -setup -suites=all -coverage=coverage -cleanup\n", os.Args[0])
	fmt.Println("")
	flag.PrintDefaults()
}