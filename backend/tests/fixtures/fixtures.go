package fixtures

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// MasterFixtures provides access to all fixture types
type MasterFixtures struct {
	*BaseFixture
	Users     *UserFixtures
	Sessions  *SessionFixtures
	KYC       *KYCFixtures
	Content   *ContentFixtures
	Categories *ContentCategoryFixtures
	Versions  *ContentVersionFixtures
	Payments  *PaymentFixtures
	Wallets   *WalletFixtures
	Transactions *TransactionFixtures
	PaymentMethods *PaymentMethodFixtures
	Messaging *MessagingFixtures
}

// NewMasterFixtures creates a new master fixtures instance with all sub-fixtures
func NewMasterFixtures(seed ...int64) *MasterFixtures {
	base := NewBaseFixture(seed...)

	return &MasterFixtures{
		BaseFixture: base,
		Users:       NewUserFixtures(seed...),
		Sessions:    NewSessionFixtures(seed...),
		KYC:         NewKYCFixtures(seed...),
		Content:     NewContentFixtures(seed...),
		Categories:  NewContentCategoryFixtures(seed...),
		Versions:    NewContentVersionFixtures(seed...),
		Payments:    NewPaymentFixtures(seed...),
		Wallets:     NewWalletFixtures(seed...),
		Transactions: NewTransactionFixtures(seed...),
		PaymentMethods: NewPaymentMethodFixtures(seed...),
		Messaging:   NewMessagingFixtures(seed...),
	}
}

// Southeast Asian Countries for testing
var SEACountries = []string{"TH", "SG", "ID", "MY", "VN", "PH"}

// ComprehensiveTestData creates a complete set of test data for all services
func (m *MasterFixtures) ComprehensiveTestData() map[string]interface{} {
	// Create users for each SEA country
	users := make([]interface{}, 0)
	userIDs := make([]uuid.UUID, 0)

	for _, country := range SEACountries {
		// Create different user types per country
		basicUser := m.Users.BasicUser(country)
		verifiedUser := m.Users.VerifiedUser(country)
		premiumUser := m.Users.PremiumUser(country)

		users = append(users, basicUser, verifiedUser, premiumUser)
		userIDs = append(userIDs, basicUser.ID, verifiedUser.ID, premiumUser.ID)
	}

	// Create sessions for users
	sessions := make([]interface{}, 0)
	for i, userID := range userIDs[:6] { // Sessions for first 6 users
		platforms := []string{"ios", "android", "web"}
		platform := platforms[i%len(platforms)]

		activeSession := m.Sessions.ActiveSession(userID, platform)
		sessions = append(sessions, activeSession)

		// Add some expired/revoked sessions
		if i%3 == 0 {
			expiredSession := m.Sessions.ExpiredSession(userID, platform)
			sessions = append(sessions, expiredSession)
		}
	}

	// Create KYC data for verified users
	kycData := make([]interface{}, 0)
	for i, userID := range userIDs {
		country := SEACountries[i%len(SEACountries)]

		if i%3 == 1 { // Verified users
			kyc := m.KYC.VerifiedKYC(userID, country)
			kycData = append(kycData, kyc)
		} else if i%3 == 2 { // Premium users
			kyc := m.KYC.PremiumKYC(userID, country)
			kycData = append(kycData, kyc)
		} else { // Basic users
			kyc := m.KYC.BasicKYC(userID, country)
			kycData = append(kycData, kyc)
		}
	}

	// Create content data
	contentData := m.Content.TestContentData()

	// Create payment data for each user
	paymentData := make(map[string]interface{})
	wallets := make([]interface{}, 0)
	transactions := make([]interface{}, 0)
	paymentMethods := make([]interface{}, 0)

	for i, userID := range userIDs[:9] { // Payment data for first 9 users
		country := SEACountries[i%len(SEACountries)]
		userData := m.Payments.TestPaymentData(userID, country)

		if userWallets, ok := userData["wallets"].([]interface{}); ok {
			wallets = append(wallets, userWallets...)
		}
		if userTransactions, ok := userData["transactions"].([]interface{}); ok {
			transactions = append(transactions, userTransactions...)
		}
		if userPaymentMethods, ok := userData["payment_methods"].([]interface{}); ok {
			paymentMethods = append(paymentMethods, userPaymentMethods...)
		}
	}

	paymentData["wallets"] = wallets
	paymentData["transactions"] = transactions
	paymentData["payment_methods"] = paymentMethods

	// Create messaging data
	messagingData := make(map[string]interface{})
	for _, country := range SEACountries[:3] { // Messaging for first 3 countries
		countryUserIDs := make([]uuid.UUID, 0)
		for i, userID := range userIDs {
			if i%len(SEACountries) < 3 {
				countryUserIDs = append(countryUserIDs, userID)
			}
		}

		if len(countryUserIDs) >= 3 {
			data := m.Messaging.TestMessagingData(countryUserIDs[:3], country)
			messagingData[country] = data
		}
	}

	return map[string]interface{}{
		"users":           users,
		"sessions":        sessions,
		"kyc":             kycData,
		"content":         contentData,
		"payments":        paymentData,
		"messaging":       messagingData,
		"metadata": map[string]interface{}{
			"generated_at":    time.Now().UTC(),
			"seed":            m.seed,
			"countries":       SEACountries,
			"total_users":     len(users),
			"total_sessions":  len(sessions),
			"total_kyc":       len(kycData),
			"fixture_version": "1.0.0",
		},
	}
}

// QuickTestData creates a minimal set of test data for rapid testing
func (m *MasterFixtures) QuickTestData(country string) map[string]interface{} {
	// Single user with complete data
	user := m.Users.BasicUser(country)
	session := m.Sessions.ActiveSession(user.ID, "web")
	kyc := m.KYC.BasicKYC(user.ID, country)

	// Basic content
	content := m.Content.BasicContent("test", country)
	category := m.Categories.RootCategory("test")

	// Basic payment data
	wallet := m.Wallets.BasicWallet(user.ID, country)
	transaction := m.Transactions.TopUpTransaction(user.ID, m.Currency(country))
	paymentMethod := m.PaymentMethods.BankAccountPaymentMethod(user.ID, country)

	// Basic messaging data
	chat := m.Messaging.DirectChat(user.ID, m.UUID("test-recipient"))
	message := m.Messaging.BasicTextMessage(chat.ID, user.ID, country)

	return map[string]interface{}{
		"user":           user,
		"session":        session,
		"kyc":            kyc,
		"content":        content,
		"category":       category,
		"wallet":         wallet,
		"transaction":    transaction,
		"payment_method": paymentMethod,
		"chat":           chat,
		"message":        message,
		"metadata": map[string]interface{}{
			"generated_at": time.Now().UTC(),
			"country":      country,
			"type":         "quick_test",
		},
	}
}

// DevUserData creates test data for development with known IDs
func (m *MasterFixtures) DevUserData() map[string]interface{} {
	// Predictable UUIDs for development
	devUserID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	devRecipientID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

	// Override base fixture to use predictable UUIDs
	devFixtures := &MasterFixtures{
		BaseFixture: &BaseFixture{seed: 12345}, // Fixed seed
		Users:       &UserFixtures{BaseFixture: &BaseFixture{seed: 12345}},
		Sessions:    &SessionFixtures{BaseFixture: &BaseFixture{seed: 12345}},
		KYC:         &KYCFixtures{BaseFixture: &BaseFixture{seed: 12345}},
		Content:     &ContentFixtures{BaseFixture: &BaseFixture{seed: 12345}},
		Categories:  &ContentCategoryFixtures{BaseFixture: &BaseFixture{seed: 12345}},
		Versions:    &ContentVersionFixtures{BaseFixture: &BaseFixture{seed: 12345}},
		Payments:    &PaymentFixtures{BaseFixture: &BaseFixture{seed: 12345}},
		Wallets:     &WalletFixtures{BaseFixture: &BaseFixture{seed: 12345}},
		Transactions: &TransactionFixtures{BaseFixture: &BaseFixture{seed: 12345}},
		PaymentMethods: &PaymentMethodFixtures{BaseFixture: &BaseFixture{seed: 12345}},
		Messaging:   &MessagingFixtures{BaseFixture: &BaseFixture{seed: 12345}},
	}

	// Create dev user data
	user := devFixtures.Users.BasicUser("TH")
	user.ID = devUserID
	user.Name = "Dev Test User"
	email := "dev@tchat-test.com"
	user.Email = &email

	recipient := devFixtures.Users.BasicUser("TH")
	recipient.ID = devRecipientID
	recipient.Name = "Dev Recipient"
	recipientEmail := "recipient@tchat-test.com"
	recipient.Email = &recipientEmail

	// Create associated data
	session := devFixtures.Sessions.ActiveSession(devUserID, "web")
	kyc := devFixtures.KYC.VerifiedKYC(devUserID, "TH")
	wallet := devFixtures.Wallets.BasicWallet(devUserID, "TH")
	paymentMethod := devFixtures.PaymentMethods.BankAccountPaymentMethod(devUserID, "TH")
	chat := devFixtures.Messaging.DirectChat(devUserID, devRecipientID)
	messages := devFixtures.Messaging.MessageThread(chat.ID, []uuid.UUID{devUserID, devRecipientID}, 5, "TH")

	return map[string]interface{}{
		"users":          []interface{}{user, recipient},
		"session":        session,
		"kyc":            kyc,
		"wallet":         wallet,
		"payment_method": paymentMethod,
		"chat":           chat,
		"messages":       messages,
		"metadata": map[string]interface{}{
			"type":           "development",
			"user_id":        devUserID.String(),
			"recipient_id":   devRecipientID.String(),
			"generated_at":   time.Now().UTC(),
		},
	}
}

// ValidationTestData creates test data specifically for validation testing
func (m *MasterFixtures) ValidationTestData() map[string]interface{} {
	validationCases := make(map[string]interface{})

	// Valid data cases
	validUser := m.Users.VerifiedUser("TH")
	validationCases["valid_user"] = validUser

	// Invalid data cases for testing validation
	invalidUser := m.Users.BasicUser("TH")
	invalidUser.Phone = nil // Invalid: no phone
	invalidUser.Email = nil // Invalid: no email
	validationCases["invalid_user_no_contact"] = invalidUser

	// Edge cases
	edgeUser := m.Users.BasicUser("TH")
	longName := ""
	for i := 0; i < 100; i++ {
		longName += "Test "
	}
	edgeUser.Name = longName // Very long name
	validationCases["edge_user_long_name"] = edgeUser

	// Invalid content
	invalidContent := m.Content.BasicContent("test")
	invalidContent.Category = "" // Invalid: empty category
	validationCases["invalid_content_empty_category"] = invalidContent

	// Invalid transaction
	userID := m.UUID("validation-user")
	invalidTransaction := m.Transactions.BasicTransaction(userID, userID, "THB") // Self-transaction
	validationCases["invalid_self_transaction"] = invalidTransaction

	return map[string]interface{}{
		"validation_cases": validationCases,
		"metadata": map[string]interface{}{
			"type":        "validation",
			"generated_at": time.Now().UTC(),
			"purpose":     "Input validation and error handling testing",
		},
	}
}

// PerformanceTestData creates large datasets for performance testing
func (m *MasterFixtures) PerformanceTestData(scale string) map[string]interface{} {
	var userCount, messageCount int

	switch scale {
	case "small":
		userCount, messageCount = 100, 1000
	case "medium":
		userCount, messageCount = 1000, 10000
	case "large":
		userCount, messageCount = 10000, 100000
	default:
		userCount, messageCount = 100, 1000
	}

	users := make([]interface{}, 0, userCount)
	userIDs := make([]uuid.UUID, 0, userCount)

	// Generate users
	for i := 0; i < userCount; i++ {
		country := SEACountries[i%len(SEACountries)]
		user := m.Users.BasicUser(country)
		user.ID = m.UUID(fmt.Sprintf("perf-user-%d", i))
		users = append(users, user)
		userIDs = append(userIDs, user.ID)
	}

	// Generate messages for performance testing
	messages := make([]interface{}, 0, messageCount)
	chatID := m.UUID("performance-chat")

	for i := 0; i < messageCount; i++ {
		senderID := userIDs[i%len(userIDs)]
		country := SEACountries[i%len(SEACountries)]

		message := m.Messaging.BasicTextMessage(chatID, senderID, country)
		message.ID = m.UUID(fmt.Sprintf("perf-message-%d", i))
		message.CreatedAt = m.PastTime(messageCount - i) // Chronological order
		messages = append(messages, message)
	}

	return map[string]interface{}{
		"users":    users,
		"messages": messages,
		"metadata": map[string]interface{}{
			"type":          "performance",
			"scale":         scale,
			"user_count":    userCount,
			"message_count": messageCount,
			"generated_at":  time.Now().UTC(),
		},
	}
}

// CleanupFixtures provides utilities for cleaning up test data
type CleanupFixtures struct {
	*BaseFixture
}

// NewCleanupFixtures creates a new cleanup fixtures instance
func NewCleanupFixtures() *CleanupFixtures {
	return &CleanupFixtures{
		BaseFixture: NewBaseFixture(),
	}
}

// GenerateCleanupSQL generates SQL statements for cleaning up test data
func (c *CleanupFixtures) GenerateCleanupSQL() []string {
	return []string{
		"DELETE FROM sessions WHERE access_token LIKE 'test-token-%';",
		"DELETE FROM kyc WHERE verification_notes = 'Automatic verification';",
		"DELETE FROM content_items WHERE tags @> '[\"test\"]';",
		"DELETE FROM content_categories WHERE name LIKE 'Test %';",
		"DELETE FROM transactions WHERE description LIKE 'Test %';",
		"DELETE FROM wallets WHERE user_id IN (SELECT id FROM users WHERE email LIKE '%@tchat-test.com');",
		"DELETE FROM users WHERE email LIKE '%@tchat-test.com' OR name LIKE 'Test %';",
	}
}

// TestDataSummary provides a summary of generated test data
func (m *MasterFixtures) TestDataSummary(data map[string]interface{}) map[string]interface{} {
	summary := make(map[string]interface{})

	// Count items in each category
	for key, value := range data {
		switch v := value.(type) {
		case []interface{}:
			summary[key+"_count"] = len(v)
		case map[string]interface{}:
			if key == "metadata" {
				summary["metadata"] = v
			} else {
				summary[key+"_sections"] = len(v)
			}
		default:
			summary[key+"_type"] = fmt.Sprintf("%T", v)
		}
	}

	summary["total_sections"] = len(data)
	summary["summary_generated_at"] = time.Now().UTC()

	return summary
}