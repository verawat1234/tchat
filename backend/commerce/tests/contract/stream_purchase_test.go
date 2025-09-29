package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// StreamPurchaseTestSuite tests the Stream content purchase API
type StreamPurchaseTestSuite struct {
	suite.Suite
	client  *http.Client
	baseURL string
}

// StreamPurchaseRequest represents the purchase request payload
type StreamPurchaseRequest struct {
	MediaContentID string `json:"mediaContentId"`
	Quantity       int    `json:"quantity"`
	MediaLicense   string `json:"mediaLicense"`
	DownloadFormat string `json:"downloadFormat"`
	CartID         string `json:"cartId,omitempty"`
}

// StreamPurchaseResponse represents the purchase response
type StreamPurchaseResponse struct {
	OrderID     string  `json:"orderId"`
	TotalAmount float64 `json:"totalAmount"`
	Currency    string  `json:"currency"`
	Success     bool    `json:"success"`
	Message     string  `json:"message"`
}

func (suite *StreamPurchaseTestSuite) SetupSuite() {
	suite.client = &http.Client{}
	suite.baseURL = "http://localhost:8083" // Commerce service port
}

// TestStreamContentPurchase tests POST /api/v1/stream/content/purchase
func (suite *StreamPurchaseTestSuite) TestStreamContentPurchase() {
	// This test MUST FAIL until the endpoint is implemented
	url := suite.baseURL + "/api/v1/stream/content/purchase"

	// Create valid purchase request
	purchaseReq := StreamPurchaseRequest{
		MediaContentID: "test-book-001",
		Quantity:       1,
		MediaLicense:   "personal",
		DownloadFormat: "PDF",
	}

	jsonData, err := json.Marshal(purchaseReq)
	suite.Require().NoError(err)

	resp, err := suite.client.Post(url, "application/json", bytes.NewBuffer(jsonData))

	// Expected behavior: endpoint should exist and return proper response
	suite.Require().NoError(err, "Request should not fail")
	defer resp.Body.Close()

	// Expected: 200 OK or 201 Created status
	suite.True(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated,
		"Should return 200 OK or 201 Created")

	// Expected: proper Content-Type
	suite.Equal("application/json", resp.Header.Get("Content-Type"), "Should return JSON")

	// Expected: valid response structure
	var response StreamPurchaseResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err, "Response should be valid JSON")

	// Expected: successful purchase response
	suite.True(response.Success, "Purchase should be successful")
	suite.NotEmpty(response.OrderID, "Order ID should not be empty")
	suite.GreaterOrEqual(response.TotalAmount, float64(0), "Total amount should be non-negative")
	suite.NotEmpty(response.Currency, "Currency should not be empty")
}

// TestStreamContentPurchaseValidation tests input validation
func (suite *StreamPurchaseTestSuite) TestStreamContentPurchaseValidation() {
	// This test MUST FAIL until validation is implemented
	url := suite.baseURL + "/api/v1/stream/content/purchase"

	// Test with invalid request (missing required fields)
	invalidReq := StreamPurchaseRequest{
		Quantity: 1,
		// Missing MediaContentID
	}

	jsonData, err := json.Marshal(invalidReq)
	suite.Require().NoError(err)

	resp, err := suite.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Expected: 400 Bad Request for invalid input
	suite.Equal(http.StatusBadRequest, resp.StatusCode, "Should return 400 for invalid input")
}

// TestStreamContentPurchaseWithCartIntegration tests cart integration
func (suite *StreamPurchaseTestSuite) TestStreamContentPurchaseWithCartIntegration() {
	// This test MUST FAIL until cart integration is implemented
	url := suite.baseURL + "/api/v1/stream/content/purchase"

	purchaseReq := StreamPurchaseRequest{
		MediaContentID: "test-podcast-001",
		Quantity:       1,
		MediaLicense:   "family",
		DownloadFormat: "MP3",
		CartID:         "test-cart-123",
	}

	jsonData, err := json.Marshal(purchaseReq)
	suite.Require().NoError(err)

	resp, err := suite.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Expected: should handle cart integration
	suite.True(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated,
		"Should handle cart integration successfully")

	var response StreamPurchaseResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err)

	// Expected: successful integration with cart
	suite.True(response.Success, "Cart integration should be successful")
}

// TestStreamContentPurchaseAuth tests authentication requirements
func (suite *StreamPurchaseTestSuite) TestStreamContentPurchaseAuth() {
	// This test MUST FAIL until authentication is implemented
	url := suite.baseURL + "/api/v1/stream/content/purchase"

	purchaseReq := StreamPurchaseRequest{
		MediaContentID: "test-book-001",
		Quantity:       1,
		MediaLicense:   "personal",
		DownloadFormat: "PDF",
	}

	jsonData, err := json.Marshal(purchaseReq)
	suite.Require().NoError(err)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	// Request without authentication header
	resp, err := suite.client.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Expected: should require authentication
	suite.Equal(http.StatusUnauthorized, resp.StatusCode, "Should require authentication")
}

// TestStreamContentPurchasePerformance tests response time requirements
func (suite *StreamPurchaseTestSuite) TestStreamContentPurchasePerformance() {
	// This test MUST FAIL until performance requirements are met
	url := suite.baseURL + "/api/v1/stream/content/purchase"

	purchaseReq := StreamPurchaseRequest{
		MediaContentID: "test-book-001",
		Quantity:       1,
		MediaLicense:   "personal",
		DownloadFormat: "PDF",
	}

	jsonData, err := json.Marshal(purchaseReq)
	suite.Require().NoError(err)

	start := time.Now()
	resp, err := suite.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	elapsed := time.Since(start)

	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Expected: <200ms response time
	suite.Less(elapsed.Milliseconds(), int64(200), "Response time should be less than 200ms")
}

func TestStreamPurchaseTestSuite(t *testing.T) {
	suite.Run(t, new(StreamPurchaseTestSuite))
}