// Shared types for integration tests
// This file contains all common struct definitions to avoid redeclaration conflicts

package integration

import "time"

// ===== User Types =====

// TestUser represents a test user for integration tests
type TestUser struct {
	ID           string `json:"id"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Country      string `json:"country"`
	Language     string `json:"language"`
	Timezone     string `json:"timezone"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AuthenticatedUser represents an authenticated user across all journeys
type AuthenticatedUser struct {
	ID          string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Country     string `json:"country"`
	Language    string `json:"language"`
	Token       string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// ===== Request Types =====

// RegistrationRequest represents user registration data
type RegistrationRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Country   string `json:"country"`
	Language  string `json:"language"`
}

// CreateProductRequest represents product creation data
type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Category    string  `json:"category"`
	Images      []string `json:"images"`
	Stock       int     `json:"stock"`
}

// CreateContentRequest represents content creation data
type CreateContentRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Category    string `json:"category"`
	Tags        []string `json:"tags"`
	Status      string `json:"status"`
}

// ===== Response Types =====

// RegistrationResponse represents registration API response
type RegistrationResponse struct {
	UserID     string `json:"userId"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Token      string `json:"token"`
	Expires    string `json:"expires"`
}

// ContentResponse represents content API response
type ContentResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	ContentType string    `json:"content_type"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Author      string    `json:"author"`
}

// ===== Commerce Types =====

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
}

// ===== Location & Geo Types =====

// LocationData represents geographical location information
type LocationData struct {
	Address     string      `json:"address"`
	City        string      `json:"city"`
	Country     string      `json:"country"`
	Coordinates Coordinates `json:"coordinates"`
	Timezone    string      `json:"timezone"`
}

// Coordinates represents latitude and longitude
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ===== Service Configuration =====

// ServicePort represents the port configuration for each microservice
type ServicePort struct {
	Auth         int
	Content      int
	Commerce     int
	Messaging    int
	Notification int
	Payment      int
}

// DefaultServicePorts returns the default port configuration
func DefaultServicePorts() ServicePort {
	return ServicePort{
		Auth:         8081,
		Content:      8082,
		Commerce:     8083,
		Messaging:    8084,
		Notification: 8085,
		Payment:      8086,
	}
}

// ===== Common API Types =====

// APIResponse represents a generic API response structure
type APIResponse struct {
	Success   bool                   `json:"success"`
	Status    string                 `json:"status"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
	Timestamp string                 `json:"timestamp"`
	Error     *APIError              `json:"error,omitempty"`
}

// APIError represents API error details
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ServiceHealthCheck represents health check response
type ServiceHealthCheck struct {
	Service   string `json:"service"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// ===== Test Utility Types =====

// TestContext represents shared test context
type TestContext struct {
	BaseURL     string
	TestUser    *TestUser
	AuthToken   string
	TestID      string
	StartTime   time.Time
}

// PerformanceMetrics represents performance measurement data
type PerformanceMetrics struct {
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Duration     int64     `json:"duration_ms"`
	RequestCount int       `json:"request_count"`
	ErrorCount   int       `json:"error_count"`
	SuccessRate  float64   `json:"success_rate"`
}