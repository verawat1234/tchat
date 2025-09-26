package fixtures_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tchat-backend/tests/fixtures"
)

// TestFixturesBasicUsage demonstrates basic fixture usage
func TestFixturesBasicUsage(t *testing.T) {
	// Create fixtures with deterministic seed for reproducible tests
	f := fixtures.NewMasterFixtures(12345)

	// Test user fixtures
	user := f.Users.BasicUser("TH")
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "TH", string(user.Country))
	assert.NotNil(t, user.Phone)
	assert.NotNil(t, user.Email)

	// Test session fixtures
	session := f.Sessions.ActiveSession(user.ID, "mobile")
	assert.Equal(t, user.ID, session.UserID)
	assert.True(t, session.IsActive)
	assert.NotEmpty(t, session.AccessToken)

	// Test content fixtures
	content := f.Content.BasicContent("announcements", "TH")
	assert.Equal(t, "announcements", content.Category)
	assert.Contains(t, content.Tags, "TH")

	// Test payment fixtures
	wallet := f.Wallets.BasicWallet(user.ID, "TH")
	assert.Equal(t, user.ID, wallet.UserID)
	assert.Equal(t, "THB", string(wallet.Currency))
	assert.True(t, wallet.Balance > 0)
}

// TestFixturesMultiCountry demonstrates multi-country support
func TestFixturesMultiCountry(t *testing.T) {
	f := fixtures.NewMasterFixtures()

	countries := []string{"TH", "SG", "ID", "MY", "VN", "PH"}

	for _, country := range countries {
		t.Run(fmt.Sprintf("Country_%s", country), func(t *testing.T) {
			// Create user for each country
			user := f.Users.BasicUser(country)
			assert.Equal(t, country, string(user.Country))

			// Check phone number format
			assert.Contains(t, *user.Phone, "+")

			// Check locale
			expectedLocales := map[string]string{
				"TH": "th-TH",
				"SG": "en-SG",
				"ID": "id-ID",
				"MY": "ms-MY",
				"VN": "vi-VN",
				"PH": "en-PH",
			}
			assert.Equal(t, expectedLocales[country], user.Locale)

			// Test currency and amounts
			currency := f.Currency(country)
			amount := f.Amount(currency)
			assert.True(t, amount > 0)

			// Test localized content
			content := f.SEAContent(country, "greeting")
			assert.NotEmpty(t, content)
		})
	}
}

// TestFixturesComprehensiveData demonstrates comprehensive test data generation
func TestFixturesComprehensiveData(t *testing.T) {
	f := fixtures.NewMasterFixtures(54321)

	// Generate comprehensive test data
	data := f.ComprehensiveTestData()

	// Verify data structure
	require.Contains(t, data, "users")
	require.Contains(t, data, "sessions")
	require.Contains(t, data, "kyc")
	require.Contains(t, data, "content")
	require.Contains(t, data, "payments")
	require.Contains(t, data, "messaging")
	require.Contains(t, data, "metadata")

	// Check users
	users, ok := data["users"].([]interface{})
	require.True(t, ok)
	assert.True(t, len(users) > 0)

	// Check metadata
	metadata, ok := data["metadata"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, metadata, "generated_at")
	assert.Contains(t, metadata, "fixture_version")
	assert.Equal(t, int64(54321), metadata["seed"])

	// Print summary for visual verification
	summary := f.TestDataSummary(data)
	summaryJSON, _ := json.MarshalIndent(summary, "", "  ")
	t.Logf("Test Data Summary:\n%s", summaryJSON)
}

// TestFixturesQuickData demonstrates quick test data for rapid testing
func TestFixturesQuickData(t *testing.T) {
	f := fixtures.NewMasterFixtures()

	// Generate quick test data for Thailand
	data := f.QuickTestData("TH")

	// Verify all essential components are present
	assert.Contains(t, data, "user")
	assert.Contains(t, data, "session")
	assert.Contains(t, data, "kyc")
	assert.Contains(t, data, "content")
	assert.Contains(t, data, "wallet")
	assert.Contains(t, data, "transaction")
	assert.Contains(t, data, "payment_method")
	assert.Contains(t, data, "chat")
	assert.Contains(t, data, "message")

	// Verify relationships
	user := data["user"].(*fixtures.UserFixtures) // This would need proper type assertion in real usage
	session := data["session"]
	wallet := data["wallet"]

	// In real usage, you'd properly type assert and verify relationships
	assert.NotNil(t, user)
	assert.NotNil(t, session)
	assert.NotNil(t, wallet)
}

// TestFixturesValidationData demonstrates validation test data
func TestFixturesValidationData(t *testing.T) {
	f := fixtures.NewMasterFixtures()

	// Generate validation test data
	data := f.ValidationTestData()

	require.Contains(t, data, "validation_cases")
	cases, ok := data["validation_cases"].(map[string]interface{})
	require.True(t, ok)

	// Verify validation cases
	assert.Contains(t, cases, "valid_user")
	assert.Contains(t, cases, "invalid_user_no_contact")
	assert.Contains(t, cases, "edge_user_long_name")
	assert.Contains(t, cases, "invalid_content_empty_category")
	assert.Contains(t, cases, "invalid_self_transaction")

	t.Logf("Generated %d validation test cases", len(cases))
}

// TestFixturesPerformanceData demonstrates performance test data generation
func TestFixturesPerformanceData(t *testing.T) {
	f := fixtures.NewMasterFixtures()

	scales := []string{"small", "medium"}

	for _, scale := range scales {
		t.Run(fmt.Sprintf("Scale_%s", scale), func(t *testing.T) {
			data := f.PerformanceTestData(scale)

			require.Contains(t, data, "users")
			require.Contains(t, data, "messages")
			require.Contains(t, data, "metadata")

			users := data["users"].([]interface{})
			messages := data["messages"].([]interface{})
			metadata := data["metadata"].(map[string]interface{})

			expectedCounts := map[string]map[string]int{
				"small":  {"users": 100, "messages": 1000},
				"medium": {"users": 1000, "messages": 10000},
			}

			assert.Equal(t, expectedCounts[scale]["users"], len(users))
			assert.Equal(t, expectedCounts[scale]["messages"], len(messages))
			assert.Equal(t, scale, metadata["scale"])

			t.Logf("Generated %d users and %d messages for %s scale",
				len(users), len(messages), scale)
		})
	}
}

// TestFixturesDevData demonstrates development-friendly data
func TestFixturesDevData(t *testing.T) {
	f := fixtures.NewMasterFixtures()

	// Generate development data with predictable IDs
	data := f.DevUserData()

	require.Contains(t, data, "users")
	require.Contains(t, data, "metadata")

	metadata := data["metadata"].(map[string]interface{})
	assert.Equal(t, "development", metadata["type"])
	assert.Contains(t, metadata, "user_id")
	assert.Contains(t, metadata, "recipient_id")

	// Verify predictable UUIDs for development
	userID := metadata["user_id"].(string)
	recipientID := metadata["recipient_id"].(string)

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", userID)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440001", recipientID)

	t.Logf("Dev data generated with user ID: %s", userID)
}

// TestFixturesSpecificScenarios demonstrates specific testing scenarios
func TestFixturesSpecificScenarios(t *testing.T) {
	f := fixtures.NewMasterFixtures()

	t.Run("KYC_Verification_Flow", func(t *testing.T) {
		userID := f.UUID("kyc-test-user")

		// Test different KYC stages
		basicKYC := f.KYC.BasicKYC(userID, "TH")
		verifiedKYC := f.KYC.VerifiedKYC(userID, "TH")
		premiumKYC := f.KYC.PremiumKYC(userID, "TH")
		rejectedKYC := f.KYC.RejectedKYC(userID, "TH")

		assert.Equal(t, "pending", string(basicKYC.Status))
		assert.Equal(t, "approved", string(verifiedKYC.Status))
		assert.Equal(t, 3, int(premiumKYC.Tier))
		assert.Equal(t, "rejected", string(rejectedKYC.Status))
	})

	t.Run("Payment_Transaction_Flow", func(t *testing.T) {
		fromUserID := f.UUID("sender")
		toUserID := f.UUID("recipient")

		// Test different transaction types
		transfer := f.Transactions.BasicTransaction(fromUserID, toUserID, "THB")
		pending := f.Transactions.PendingTransaction(fromUserID, toUserID, "THB")
		failed := f.Transactions.FailedTransaction(fromUserID, toUserID, "THB")
		topup := f.Transactions.TopUpTransaction(fromUserID, "THB")

		assert.Equal(t, "completed", string(transfer.Status))
		assert.Equal(t, "pending", string(pending.Status))
		assert.Equal(t, "failed", string(failed.Status))
		assert.Equal(t, "topup", string(topup.Type))
	})

	t.Run("Messaging_Thread", func(t *testing.T) {
		chatID := f.UUID("test-chat")
		userIDs := []string{
			f.UUID("user1").String(),
			f.UUID("user2").String(),
		}

		// Test message thread generation
		messages := f.Messaging.MessageThread(
			f.UUID(chatID.String()),
			[]string{userIDs[0], userIDs[1]},
			10,
			"TH",
		)

		assert.Len(t, messages, 10)

		// Verify chronological order
		for i := 1; i < len(messages); i++ {
			assert.True(t, messages[i].CreatedAt.After(messages[i-1].CreatedAt),
				"Messages should be in chronological order")
		}
	})
}

// TestFixturesCleanup demonstrates cleanup utilities
func TestFixturesCleanup(t *testing.T) {
	cleanup := fixtures.NewCleanupFixtures()

	// Generate cleanup SQL
	sqlStatements := cleanup.GenerateCleanupSQL()

	assert.True(t, len(sqlStatements) > 0)

	for _, sql := range sqlStatements {
		assert.Contains(t, sql, "DELETE FROM")
		t.Logf("Cleanup SQL: %s", sql)
	}
}

// TestFixturesReproducibility demonstrates reproducible test data
func TestFixturesReproducibility(t *testing.T) {
	seed := int64(99999)

	// Create two fixtures with same seed
	f1 := fixtures.NewMasterFixtures(seed)
	f2 := fixtures.NewMasterFixtures(seed)

	// Generate same data
	user1 := f1.Users.BasicUser("TH")
	user2 := f2.Users.BasicUser("TH")

	// Should be identical
	assert.Equal(t, user1.ID, user2.ID)
	assert.Equal(t, user1.Name, user2.Name)
	assert.Equal(t, user1.Phone, user2.Phone)
	assert.Equal(t, user1.Email, user2.Email)

	t.Logf("Reproducible user ID: %s", user1.ID)
}

// BenchmarkFixturesGeneration benchmarks fixture generation performance
func BenchmarkFixturesGeneration(b *testing.B) {
	f := fixtures.NewMasterFixtures()

	b.Run("Single_User", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = f.Users.BasicUser("TH")
		}
	})

	b.Run("User_With_Session", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			user := f.Users.BasicUser("TH")
			_ = f.Sessions.ActiveSession(user.ID, "web")
		}
	})

	b.Run("Complete_User_Data", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			user := f.Users.BasicUser("TH")
			_ = f.Sessions.ActiveSession(user.ID, "web")
			_ = f.KYC.BasicKYC(user.ID, "TH")
			_ = f.Wallets.BasicWallet(user.ID, "TH")
		}
	})

	b.Run("Quick_Test_Data", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = f.QuickTestData("TH")
		}
	})
}