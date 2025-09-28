package contract

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
)

func TestPactInstallation(t *testing.T) {
	// Test that pact-go v2 is properly installed and FFI library is available
	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "content-consumer",
		Provider: "content-service",
		Host:     "127.0.0.1",
		Port:     8000,
	})
	if err != nil {
		t.Fatalf("Error creating Pact: %v", err)
	}

	// Add a simple interaction to test FFI integration
	mockProvider.
		AddInteraction().
		Given("Pact FFI library is available").
		UponReceiving("a health check request").
		WithRequest("GET", "/health").
		WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
			b.Header("Content-Type", matchers.String("application/json"))
			b.JSONBody(map[string]interface{}{
				"status":  matchers.String("ok"),
				"service": matchers.String("content-service"),
			})
		})

	// Test the mock service
	err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/health", config.Host, config.Port), nil)
		if err != nil {
			return err
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		log.Printf("✅ Pact FFI library installation successful!")
		log.Printf("✅ pact-go v2 API working correctly!")
		log.Printf("✅ Mock server response: %d", resp.StatusCode)

		return nil
	})

	if err != nil {
		t.Fatalf("Test execution failed: %v", err)
	}

	t.Log("🎉 Pact installation test passed! FFI library is working correctly.")
}