package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// T036: Integration test external service failures
// Tests system resilience and fallback mechanisms when external services fail
type ServiceFailuresTestSuite struct {
	suite.Suite
	router         *gin.Engine
	serviceStatus  map[string]bool // service -> healthy
	fallbackData   map[string]interface{}
	circuitBreaker map[string]int // service -> failure count
}

func TestServiceFailuresSuite(t *testing.T) {
	suite.Run(t, new(ServiceFailuresTestSuite))
}

func (suite *ServiceFailuresTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.fallbackData = make(map[string]interface{})
	suite.circuitBreaker = make(map[string]int)

	suite.setupFailureTestEndpoints()
}

func (suite *ServiceFailuresTestSuite) SetupTest() {
	// Reset service status for each test to ensure isolation
	suite.serviceStatus = map[string]bool{
		"payment_gateway": true,
		"sms_service":     true,
		"email_service":   true,
		"kyc_provider":    true,
	}
	// Reset circuit breaker counters to ensure test isolation
	suite.circuitBreaker = make(map[string]int)
}

func (suite *ServiceFailuresTestSuite) setupFailureTestEndpoints() {
	// Payment with fallback
	suite.router.POST("/payments/process", func(c *gin.Context) {
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		// Check if payment gateway is healthy
		if !suite.serviceStatus["payment_gateway"] {
			// Increment failure count
			suite.circuitBreaker["payment_gateway"]++

			// Circuit breaker: after 3 failures, use fallback
			if suite.circuitBreaker["payment_gateway"] >= 3 {
				c.JSON(http.StatusOK, gin.H{
					"status":       "queued",
					"message":      "Payment queued for processing",
					"fallback":     true,
					"retry_after":  300, // 5 minutes
					"reference_id": "FALLBACK_" + time.Now().Format("20060102150405"),
				})
				return
			}

			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "payment_gateway_unavailable",
				"message": "Payment gateway is temporarily unavailable",
			})
			return
		}

		// Reset circuit breaker on success
		suite.circuitBreaker["payment_gateway"] = 0

		c.JSON(http.StatusOK, gin.H{
			"status":         "completed",
			"transaction_id": "TXN_" + time.Now().Format("20060102150405"),
		})
	})

	// SMS with email fallback
	suite.router.POST("/notifications/send", func(c *gin.Context) {
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		notificationType := req["type"].(string)

		// Try primary SMS service
		if notificationType == "sms" && !suite.serviceStatus["sms_service"] {
			// Fallback to email if SMS fails
			if suite.serviceStatus["email_service"] {
				c.JSON(http.StatusOK, gin.H{
					"status":          "sent",
					"fallback_method": "email",
					"message":         "SMS unavailable, sent via email instead",
				})
				return
			}

			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "all_notification_services_unavailable",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "sent",
			"method": notificationType,
		})
	})

	// KYC with manual review fallback
	suite.router.POST("/kyc/verify", func(c *gin.Context) {
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		if !suite.serviceStatus["kyc_provider"] {
			// Queue for manual review
			c.JSON(http.StatusAccepted, gin.H{
				"status":       "manual_review",
				"message":      "KYC provider unavailable, queued for manual review",
				"estimated_completion": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      "verified",
			"verified_at": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Service health control (for testing)
	suite.router.POST("/admin/service/:service/health", func(c *gin.Context) {
		service := c.Param("service")
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		healthy := req["healthy"].(bool)
		suite.serviceStatus[service] = healthy

		c.JSON(http.StatusOK, gin.H{
			"service": service,
			"healthy": healthy,
		})
	})
}

func (suite *ServiceFailuresTestSuite) TestPaymentGatewayFailure() {
	// Simulate payment gateway failure
	healthData := map[string]interface{}{"healthy": false}
	jsonData, _ := json.Marshal(healthData)
	req := httptest.NewRequest("POST", "/admin/service/payment_gateway/health", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Attempt payment - should fail initially
	paymentData := map[string]interface{}{
		"amount":   100.0,
		"currency": "THB",
	}

	for i := 0; i < 3; i++ {
		jsonData, _ = json.Marshal(paymentData)
		req = httptest.NewRequest("POST", "/payments/process", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		if i < 2 {
			assert.Equal(suite.T(), http.StatusServiceUnavailable, w.Code)
		} else {
			// Circuit breaker should activate on 3rd failure
			assert.Equal(suite.T(), http.StatusOK, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(suite.T(), "queued", response["status"])
			assert.Equal(suite.T(), true, response["fallback"])
		}
	}
}

func (suite *ServiceFailuresTestSuite) TestSMSFallbackToEmail() {
	// Disable SMS service
	healthData := map[string]interface{}{"healthy": false}
	jsonData, _ := json.Marshal(healthData)
	req := httptest.NewRequest("POST", "/admin/service/sms_service/health", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Try to send SMS notification
	notificationData := map[string]interface{}{
		"type":    "sms",
		"to":      "+66812345678",
		"message": "Your verification code is 123456",
	}

	jsonData, _ = json.Marshal(notificationData)
	req = httptest.NewRequest("POST", "/notifications/send", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	start := time.Now()
	suite.router.ServeHTTP(w, req)
	duration := time.Since(start)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.True(suite.T(), duration < 100*time.Millisecond, "Fallback should be fast")

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "sent", response["status"])
	assert.Equal(suite.T(), "email", response["fallback_method"])
}

func (suite *ServiceFailuresTestSuite) TestKYCProviderFailure() {
	// Disable KYC provider
	healthData := map[string]interface{}{"healthy": false}
	jsonData, _ := json.Marshal(healthData)
	req := httptest.NewRequest("POST", "/admin/service/kyc_provider/health", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Attempt KYC verification
	kycData := map[string]interface{}{
		"document_type":   "national_id",
		"document_number": "123456789",
		"country":         "TH",
	}

	jsonData, _ = json.Marshal(kycData)
	req = httptest.NewRequest("POST", "/kyc/verify", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusAccepted, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "manual_review", response["status"])
	assert.NotEmpty(suite.T(), response["estimated_completion"])
}

func (suite *ServiceFailuresTestSuite) TestServiceRecovery() {
	// Disable payment gateway
	healthData := map[string]interface{}{"healthy": false}
	jsonData, _ := json.Marshal(healthData)
	req := httptest.NewRequest("POST", "/admin/service/payment_gateway/health", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Verify service is down
	paymentData := map[string]interface{}{"amount": 100.0}
	jsonData, _ = json.Marshal(paymentData)
	req = httptest.NewRequest("POST", "/payments/process", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusServiceUnavailable, w.Code)

	// Restore service
	healthData["healthy"] = true
	jsonData, _ = json.Marshal(healthData)
	req = httptest.NewRequest("POST", "/admin/service/payment_gateway/health", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Verify service works again
	jsonData, _ = json.Marshal(paymentData)
	req = httptest.NewRequest("POST", "/payments/process", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "completed", response["status"])
}

func (suite *ServiceFailuresTestSuite) TestCascadingFailures() {
	// Disable multiple services
	services := []string{"sms_service", "email_service"}

	for _, service := range services {
		healthData := map[string]interface{}{"healthy": false}
		jsonData, _ := json.Marshal(healthData)
		req := httptest.NewRequest("POST", "/admin/service/"+service+"/health", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
	}

	// Try notification - should fail completely
	notificationData := map[string]interface{}{
		"type":    "sms",
		"to":      "+66812345678",
		"message": "Test message",
	}

	jsonData, _ := json.Marshal(notificationData)
	req := httptest.NewRequest("POST", "/notifications/send", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusServiceUnavailable, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "all_notification_services_unavailable", response["error"])
}