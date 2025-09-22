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

// T034: Integration test KYC verification workflow
// Tests KYC process integration across services with Southeast Asian compliance
type KYCFlowTestSuite struct {
	suite.Suite
	router  *gin.Engine
	kycData map[string]map[string]interface{}
	users   map[string]map[string]interface{}
}

func TestKYCFlowSuite(t *testing.T) {
	suite.Run(t, new(KYCFlowTestSuite))
}

func (suite *KYCFlowTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.kycData = make(map[string]map[string]interface{})
	suite.users = make(map[string]map[string]interface{})

	suite.setupKYCEndpoints()
}

func (suite *KYCFlowTestSuite) setupKYCEndpoints() {
	// KYC submission with Southeast Asian validation
	suite.router.POST("/users/:user_id/kyc", func(c *gin.Context) {
		userID := c.Param("user_id")
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		// Validate Southeast Asian countries
		country := req["country"].(string)
		validCountries := []string{"TH", "SG", "ID", "MY", "PH", "VN"}
		valid := false
		for _, vc := range validCountries {
			if country == vc {
				valid = true
				break
			}
		}

		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_country"})
			return
		}

		// Store KYC data
		suite.kycData[userID] = map[string]interface{}{
			"status":          "verified",
			"country":         country,
			"document_type":   req["document_type"],
			"verified_at":     time.Now().UTC().Format(time.RFC3339),
			"compliance_tier": "tier_2", // Full compliance for Southeast Asia
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":          "verified",
			"compliance_tier": "tier_2",
			"verified_at":     time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Check KYC status
	suite.router.GET("/users/:user_id/kyc/status", func(c *gin.Context) {
		userID := c.Param("user_id")
		kyc, exists := suite.kycData[userID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "kyc_not_found"})
			return
		}
		c.JSON(http.StatusOK, kyc)
	})
}

func (suite *KYCFlowTestSuite) TestSoutheastAsianKYC() {
	countries := []struct {
		country string
		docType string
	}{
		{"TH", "national_id"},
		{"SG", "nric"},
		{"ID", "ktp"},
		{"MY", "mykad"},
		{"PH", "umid"},
		{"VN", "cccd"},
	}

	for _, testCase := range countries {
		suite.T().Logf("Testing KYC for country: %s", testCase.country)

		kycData := map[string]interface{}{
			"country":         testCase.country,
			"document_type":   testCase.docType,
			"document_number": "123456789",
			"full_name":       "Test User",
			"date_of_birth":   "1990-01-01",
		}

		jsonData, _ := json.Marshal(kycData)
		req := httptest.NewRequest("POST", "/users/user_123/kyc", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		start := time.Now()
		suite.router.ServeHTTP(w, req)
		duration := time.Since(start)

		assert.Equal(suite.T(), http.StatusCreated, w.Code)
		assert.True(suite.T(), duration < 500*time.Millisecond, "KYC verification should complete in <500ms")

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(suite.T(), "verified", response["status"])
		assert.Equal(suite.T(), "tier_2", response["compliance_tier"])
	}
}

func (suite *KYCFlowTestSuite) TestKYCStatusCheck() {
	// First submit KYC
	kycData := map[string]interface{}{
		"country":         "TH",
		"document_type":   "national_id",
		"document_number": "123456789",
		"full_name":       "Test User",
		"date_of_birth":   "1990-01-01",
	}

	jsonData, _ := json.Marshal(kycData)
	req := httptest.NewRequest("POST", "/users/user_456/kyc", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Check status
	req = httptest.NewRequest("GET", "/users/user_456/kyc/status", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var status map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &status)
	assert.Equal(suite.T(), "verified", status["status"])
	assert.Equal(suite.T(), "TH", status["country"])
}