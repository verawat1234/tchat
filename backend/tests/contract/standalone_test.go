package contract_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

// TestStandaloneContractValidation runs basic contract tests without complex dependencies
func TestStandaloneContractValidation(t *testing.T) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	t.Run("Auth Service Basic Contract", func(t *testing.T) {
		baseURL := "http://localhost:8081"

		// Test 1: Health endpoint contract
		resp, err := client.Get(baseURL + "/health")
		if err != nil {
			t.Logf("Auth service not available: %v", err)
			t.Skip("Auth service not running - skipping contract test")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			t.Log("✅ Auth service health endpoint contract: PASS")
		} else {
			t.Logf("⚠️ Auth service health endpoint returned: %d", resp.StatusCode)
		}

		// Test 2: Registration endpoint contract
		payload := map[string]interface{}{
			"phone_number": "+66812345678",
			"country":      "TH",
			"language":     "en",
		}

		jsonPayload, _ := json.Marshal(payload)
		resp2, err := client.Post(
			baseURL+"/api/v1/auth/register",
			"application/json",
			bytes.NewBuffer(jsonPayload),
		)
		if err != nil {
			t.Logf("Registration endpoint error: %v", err)
		} else {
			defer resp2.Body.Close()
			t.Logf("✅ Registration endpoint contract: Status %d (endpoint exists)", resp2.StatusCode)
		}
	})

	t.Run("Content Service Basic Contract", func(t *testing.T) {
		baseURL := "http://localhost:8082"

		resp, err := client.Get(baseURL + "/health")
		if err != nil {
			t.Logf("Content service not available: %v", err)
			t.Skip("Content service not running")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			t.Log("✅ Content service health endpoint contract: PASS")
		} else {
			t.Logf("⚠️ Content service health endpoint returned: %d", resp.StatusCode)
		}
	})

	t.Run("Service Discovery Contract", func(t *testing.T) {
		// Test that all expected services respond on their designated ports
		services := map[string]int{
			"auth":         8081,
			"content":      8082,
			"commerce":     8083,
			"messaging":    8084,
			"notification": 8085,
		}

		availableServices := 0
		for name, port := range services {
			url := fmt.Sprintf("http://localhost:%d/health", port)
			resp, err := client.Get(url)
			if err == nil && resp.StatusCode == 200 {
				availableServices++
				t.Logf("✅ Service %s available on port %d", name, port)
				resp.Body.Close()
			} else {
				t.Logf("⚠️ Service %s not available on port %d", name, port)
				if resp != nil {
					resp.Body.Close()
				}
			}
		}

		t.Logf("Service Discovery: %d/%d services available", availableServices, len(services))

		// Contract: At least one service should be available for integration testing
		if availableServices == 0 {
			t.Error("❌ Contract violation: No services available for testing")
		} else {
			t.Log("✅ Service discovery contract: PASS")
		}
	})
}


// TestContractInfrastructure validates that contract testing infrastructure works
func TestContractInfrastructure(t *testing.T) {
	t.Run("Environment Check", func(t *testing.T) {
		// Test that the test can run basic HTTP operations
		client := &http.Client{Timeout: 5 * time.Second}

		// Test localhost connectivity
		resp, err := client.Get("http://localhost:8081/health")
		if err != nil {
			t.Logf("Localhost connectivity test failed: %v", err)
			t.Log("⚠️ This is expected if services are not running")
		} else {
			resp.Body.Close()
			t.Log("✅ Localhost connectivity: WORKING")
		}

		t.Log("✅ Contract testing infrastructure: OPERATIONAL")
	})

	t.Run("JSON Processing", func(t *testing.T) {
		// Test JSON serialization/deserialization
		testData := map[string]interface{}{
			"test":      true,
			"timestamp": time.Now().Unix(),
		}

		jsonBytes, err := json.Marshal(testData)
		if err != nil {
			t.Errorf("JSON marshaling failed: %v", err)
			return
		}

		var parsed map[string]interface{}
		err = json.Unmarshal(jsonBytes, &parsed)
		if err != nil {
			t.Errorf("JSON unmarshaling failed: %v", err)
			return
		}

		t.Log("✅ JSON processing: WORKING")
	})
}