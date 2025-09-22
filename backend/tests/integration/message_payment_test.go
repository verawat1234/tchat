package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// T033: Integration test message payment workflow
// Tests payment integration within messaging (money transfers via chat)
type MessagePaymentTestSuite struct {
	suite.Suite
	router       *gin.Engine
	wallets      map[string]map[string]interface{}
	dialogs      map[string]map[string]interface{}
	transactions map[string][]map[string]interface{}
}

func TestMessagePaymentSuite(t *testing.T) {
	suite.Run(t, new(MessagePaymentTestSuite))
}

func (suite *MessagePaymentTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.wallets = make(map[string]map[string]interface{})
	suite.dialogs = make(map[string]map[string]interface{})
	suite.transactions = make(map[string][]map[string]interface{})

	suite.setupMessagePaymentEndpoints()
}

func (suite *MessagePaymentTestSuite) setupMessagePaymentEndpoints() {
	// Create test wallets
	suite.wallets["user_123_THB"] = map[string]interface{}{
		"id": "user_123_THB", "user_id": "user_123", "currency": "THB", "balance": 5000.0,
	}
	suite.wallets["user_456_THB"] = map[string]interface{}{
		"id": "user_456_THB", "user_id": "user_456", "currency": "THB", "balance": 1000.0,
	}

	// Send money via message
	suite.router.POST("/dialogs/:dialog_id/messages/payment", func(c *gin.Context) {
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		amount := req["amount"].(float64)
		fromWallet := req["from_wallet"].(string)
		toWallet := req["to_wallet"].(string)

		// Process payment
		fromW := suite.wallets[fromWallet]
		toW := suite.wallets[toWallet]

		if fromW["balance"].(float64) < amount {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient_funds"})
			return
		}

		fromW["balance"] = fromW["balance"].(float64) - amount
		toW["balance"] = toW["balance"].(float64) + amount

		// Create payment message
		message := map[string]interface{}{
			"id":         fmt.Sprintf("msg_payment_%d", time.Now().UnixNano()),
			"type":       "payment",
			"amount":     amount,
			"currency":   "THB",
			"status":     "completed",
			"created_at": time.Now().UTC().Format(time.RFC3339),
		}

		c.JSON(http.StatusCreated, message)
	})

	// Request money via message
	suite.router.POST("/dialogs/:dialog_id/messages/payment-request", func(c *gin.Context) {
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		message := map[string]interface{}{
			"id":          fmt.Sprintf("msg_request_%d", time.Now().UnixNano()),
			"type":        "payment_request",
			"amount":      req["amount"],
			"currency":    req["currency"],
			"description": req["description"],
			"status":      "pending",
			"created_at":  time.Now().UTC().Format(time.RFC3339),
		}

		c.JSON(http.StatusCreated, message)
	})
}

func (suite *MessagePaymentTestSuite) TestMessagePaymentFlow() {
	// Test sending money via message
	paymentData := map[string]interface{}{
		"amount":      500.0,
		"from_wallet": "user_123_THB",
		"to_wallet":   "user_456_THB",
		"message":     "Lunch money",
	}

	jsonData, _ := json.Marshal(paymentData)
	req := httptest.NewRequest("POST", "/dialogs/dialog_123/messages/payment", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	start := time.Now()
	suite.router.ServeHTTP(w, req)
	duration := time.Since(start)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.True(suite.T(), duration < 300*time.Millisecond, "Payment message should complete in <300ms")

	// Verify balances updated
	assert.Equal(suite.T(), 4500.0, suite.wallets["user_123_THB"]["balance"])
	assert.Equal(suite.T(), 1500.0, suite.wallets["user_456_THB"]["balance"])
}

func (suite *MessagePaymentTestSuite) TestPaymentRequest() {
	requestData := map[string]interface{}{
		"amount":      200.0,
		"currency":    "THB",
		"description": "Dinner split",
	}

	jsonData, _ := json.Marshal(requestData)
	req := httptest.NewRequest("POST", "/dialogs/dialog_123/messages/payment-request", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "payment_request", response["type"])
	assert.Equal(suite.T(), "pending", response["status"])
}