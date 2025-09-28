package contract

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// TestConfig holds configuration for contract tests
type TestConfig struct {
	// Mock Server Configuration
	MockServerHost         string
	ContentServiceMockPort int
	AuthServiceMockPort    int
	CommerceServiceMockPort int
	MessagingServiceMockPort int
	PaymentServiceMockPort int
	NotificationServiceMockPort int

	// Provider Service Ports
	ContentServiceProviderPort int
	AuthServiceProviderPort    int
	CommerceServiceProviderPort int
	MessagingServiceProviderPort int
	PaymentServiceProviderPort int
	NotificationServiceProviderPort int

	// Pact Configuration
	PactLogLevel        string
	PactDir            string
	PactConsumerVersion string
	PactProviderVersion string

	// Consumer Names
	WebClientName    string
	MobileClientName string

	// Provider Names
	ContentServiceName     string
	AuthServiceName       string
	CommerceServiceName   string
	MessagingServiceName  string
	PaymentServiceName    string
	NotificationServiceName string

	// Test Configuration
	GinMode              string
	TestTimeout          string
	EnableVerboseOutput  bool

	// Performance Configuration
	MaxResponseTimeMs int
	TestConcurrency   int
	RequestTimeout    string
}

// LoadTestConfig loads configuration from .env.test file
func LoadTestConfig() (*TestConfig, error) {
	config := &TestConfig{
		// Default values
		MockServerHost:           "127.0.0.1",
		ContentServiceMockPort:   8090,
		AuthServiceMockPort:     8091,
		CommerceServiceMockPort: 8092,
		MessagingServiceMockPort: 8093,
		PaymentServiceMockPort:  8094,
		NotificationServiceMockPort: 8095,

		ContentServiceProviderPort: 8080,
		AuthServiceProviderPort:    8081,
		CommerceServiceProviderPort: 8082,
		MessagingServiceProviderPort: 8083,
		PaymentServiceProviderPort: 8084,
		NotificationServiceProviderPort: 8085,

		PactLogLevel:        "INFO",
		PactDir:            "./pacts",
		PactConsumerVersion: "1.0.0",
		PactProviderVersion: "1.0.0",

		WebClientName:    "content-web-client",
		MobileClientName: "content-mobile-client",

		ContentServiceName:     "content-service",
		AuthServiceName:       "auth-service",
		CommerceServiceName:   "commerce-service",
		MessagingServiceName:  "messaging-service",
		PaymentServiceName:    "payment-service",
		NotificationServiceName: "notification-service",

		GinMode:              "test",
		TestTimeout:          "30s",
		EnableVerboseOutput:  true,

		MaxResponseTimeMs: 200,
		TestConcurrency:   10,
		RequestTimeout:    "5s",
	}

	// Try to load from .env.test file
	if err := loadEnvFile(".env.test", config); err != nil {
		// If .env.test doesn't exist, try to load from environment variables
		loadEnvVars(config)
	}

	return config, nil
}

// loadEnvFile loads environment variables from the specified file
func loadEnvFile(filename string, config *TestConfig) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set configuration values
		setConfigValue(config, key, value)
	}

	return scanner.Err()
}

// loadEnvVars loads configuration from environment variables
func loadEnvVars(config *TestConfig) {
	if val := os.Getenv("PACT_MOCK_SERVER_HOST"); val != "" {
		config.MockServerHost = val
	}
	if val := os.Getenv("CONTENT_SERVICE_MOCK_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			config.ContentServiceMockPort = port
		}
	}
	if val := os.Getenv("AUTH_SERVICE_MOCK_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			config.AuthServiceMockPort = port
		}
	}
	if val := os.Getenv("COMMERCE_SERVICE_MOCK_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			config.CommerceServiceMockPort = port
		}
	}
	if val := os.Getenv("MESSAGING_SERVICE_MOCK_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			config.MessagingServiceMockPort = port
		}
	}
	if val := os.Getenv("PAYMENT_SERVICE_MOCK_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			config.PaymentServiceMockPort = port
		}
	}
	if val := os.Getenv("NOTIFICATION_SERVICE_MOCK_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			config.NotificationServiceMockPort = port
		}
	}

	// Load other configuration values
	if val := os.Getenv("PACT_LOG_LEVEL"); val != "" {
		config.PactLogLevel = val
	}
	if val := os.Getenv("PACT_DIR"); val != "" {
		config.PactDir = val
	}
	if val := os.Getenv("WEB_CLIENT_NAME"); val != "" {
		config.WebClientName = val
	}
	if val := os.Getenv("CONTENT_SERVICE_NAME"); val != "" {
		config.ContentServiceName = val
	}
}

// setConfigValue sets a configuration value based on the key
func setConfigValue(config *TestConfig, key, value string) {
	switch key {
	case "PACT_MOCK_SERVER_HOST":
		config.MockServerHost = value
	case "CONTENT_SERVICE_MOCK_PORT":
		if port, err := strconv.Atoi(value); err == nil {
			config.ContentServiceMockPort = port
		}
	case "AUTH_SERVICE_MOCK_PORT":
		if port, err := strconv.Atoi(value); err == nil {
			config.AuthServiceMockPort = port
		}
	case "COMMERCE_SERVICE_MOCK_PORT":
		if port, err := strconv.Atoi(value); err == nil {
			config.CommerceServiceMockPort = port
		}
	case "MESSAGING_SERVICE_MOCK_PORT":
		if port, err := strconv.Atoi(value); err == nil {
			config.MessagingServiceMockPort = port
		}
	case "PAYMENT_SERVICE_MOCK_PORT":
		if port, err := strconv.Atoi(value); err == nil {
			config.PaymentServiceMockPort = port
		}
	case "NOTIFICATION_SERVICE_MOCK_PORT":
		if port, err := strconv.Atoi(value); err == nil {
			config.NotificationServiceMockPort = port
		}
	case "PACT_LOG_LEVEL":
		config.PactLogLevel = value
	case "PACT_DIR":
		config.PactDir = value
	case "PACT_CONSUMER_VERSION":
		config.PactConsumerVersion = value
	case "PACT_PROVIDER_VERSION":
		config.PactProviderVersion = value
	case "WEB_CLIENT_NAME":
		config.WebClientName = value
	case "MOBILE_CLIENT_NAME":
		config.MobileClientName = value
	case "CONTENT_SERVICE_NAME":
		config.ContentServiceName = value
	case "AUTH_SERVICE_NAME":
		config.AuthServiceName = value
	case "GIN_MODE":
		config.GinMode = value
	case "TEST_TIMEOUT":
		config.TestTimeout = value
	case "ENABLE_VERBOSE_OUTPUT":
		config.EnableVerboseOutput = value == "true"
	case "MAX_RESPONSE_TIME_MS":
		if val, err := strconv.Atoi(value); err == nil {
			config.MaxResponseTimeMs = val
		}
	case "TEST_CONCURRENCY":
		if val, err := strconv.Atoi(value); err == nil {
			config.TestConcurrency = val
		}
	case "REQUEST_TIMEOUT":
		config.RequestTimeout = value
	}
}

// GetMockServerURL returns the full URL for a service's mock server
func (c *TestConfig) GetMockServerURL(service string) string {
	var port int
	switch service {
	case "content":
		port = c.ContentServiceMockPort
	case "auth":
		port = c.AuthServiceMockPort
	case "commerce":
		port = c.CommerceServiceMockPort
	case "messaging":
		port = c.MessagingServiceMockPort
	case "payment":
		port = c.PaymentServiceMockPort
	case "notification":
		port = c.NotificationServiceMockPort
	default:
		port = c.ContentServiceMockPort
	}

	return fmt.Sprintf("http://%s:%d", c.MockServerHost, port)
}

// GetProviderURL returns the full URL for a service's provider
func (c *TestConfig) GetProviderURL(service string) string {
	var port int
	switch service {
	case "content":
		port = c.ContentServiceProviderPort
	case "auth":
		port = c.AuthServiceProviderPort
	case "commerce":
		port = c.CommerceServiceProviderPort
	case "messaging":
		port = c.MessagingServiceProviderPort
	case "payment":
		port = c.PaymentServiceProviderPort
	case "notification":
		port = c.NotificationServiceProviderPort
	default:
		port = c.ContentServiceProviderPort
	}

	return fmt.Sprintf("http://%s:%d", c.MockServerHost, port)
}