package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// IntegrationTestSuite runs all microservice integration tests
type IntegrationTestSuite struct {
	suite.Suite
}

// TestAuthIntegration runs Auth service integration tests
func (suite *IntegrationTestSuite) TestAuthIntegration() {
	RunAuthIntegrationTests(suite.T())
}

// TestContentIntegration runs Content service integration tests
func (suite *IntegrationTestSuite) TestContentIntegration() {
	RunContentIntegrationTests(suite.T())
}

// TestCommerceIntegration runs Commerce service integration tests
func (suite *IntegrationTestSuite) TestCommerceIntegration() {
	RunCommerceIntegrationTests(suite.T())
}

// TestNotificationIntegration runs Notification service integration tests
func (suite *IntegrationTestSuite) TestNotificationIntegration() {
	RunNotificationIntegrationTests(suite.T())
}

// TestMessagingIntegration runs Messaging service integration tests
func (suite *IntegrationTestSuite) TestMessagingIntegration() {
	RunMessagingIntegrationTests(suite.T())
}

// TestPaymentIntegration runs Payment service integration tests
func (suite *IntegrationTestSuite) TestPaymentIntegration() {
	RunPaymentIntegrationTests(suite.T())
}

// TestRunAllIntegrationTests runs the complete integration test suite
func TestRunAllIntegrationTests(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}