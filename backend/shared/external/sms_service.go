package external

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// SMSProvider represents different SMS service providers
type SMSProvider string

const (
	TwilioProvider   SMSProvider = "twilio"
	AWSProvider      SMSProvider = "aws"
	VonageProvider   SMSProvider = "vonage"
	ThaiSMSProvider  SMSProvider = "thai_sms"  // For Thailand
	SingtelProvider  SMSProvider = "singtel"   // For Singapore
	TelkomselProvider SMSProvider = "telkomsel" // For Indonesia
)

// SMSConfig holds SMS service configuration
type SMSConfig struct {
	Provider        SMSProvider `mapstructure:"provider" validate:"required"`
	APIKey          string      `mapstructure:"api_key" validate:"required"`
	APISecret       string      `mapstructure:"api_secret"`
	FromNumber      string      `mapstructure:"from_number"`
	FromName        string      `mapstructure:"from_name"`
	Timeout         time.Duration `mapstructure:"timeout"`
	MaxRetries      int         `mapstructure:"max_retries"`
	RetryDelay      time.Duration `mapstructure:"retry_delay"`
	RateLimitRPS    int         `mapstructure:"rate_limit_rps"`
	EnableWebhook   bool        `mapstructure:"enable_webhook"`
	WebhookURL      string      `mapstructure:"webhook_url"`
	CountryPriority []string    `mapstructure:"country_priority"` // Country codes for provider priority
}

// DefaultSMSConfig returns default SMS configuration
func DefaultSMSConfig() *SMSConfig {
	return &SMSConfig{
		Provider:        TwilioProvider,
		FromName:        "Tchat",
		Timeout:         30 * time.Second,
		MaxRetries:      3,
		RetryDelay:      2 * time.Second,
		RateLimitRPS:    10,
		EnableWebhook:   true,
		CountryPriority: []string{"TH", "SG", "ID", "MY", "PH", "VN"},
	}
}

// SMSMessage represents an SMS message
type SMSMessage struct {
	To          string            `json:"to" validate:"required"`
	Message     string            `json:"message" validate:"required,max=160"`
	From        string            `json:"from,omitempty"`
	CountryCode string            `json:"country_code,omitempty"`
	Template    string            `json:"template,omitempty"`
	Variables   map[string]string `json:"variables,omitempty"`
	Priority    string            `json:"priority,omitempty"` // high, normal, low
	ScheduledAt *time.Time        `json:"scheduled_at,omitempty"`
	ValidUntil  *time.Time        `json:"valid_until,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SMSResponse represents SMS sending response
type SMSResponse struct {
	MessageID   string            `json:"message_id"`
	Status      string            `json:"status"`
	Provider    SMSProvider       `json:"provider"`
	Cost        float64           `json:"cost,omitempty"`
	Parts       int               `json:"parts"`
	Encoding    string            `json:"encoding"`
	CountryCode string            `json:"country_code"`
	Error       string            `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	SentAt      time.Time         `json:"sent_at"`
}

// SMSStatus represents message delivery status
type SMSStatus struct {
	MessageID     string    `json:"message_id"`
	Status        string    `json:"status"` // sent, delivered, failed, undelivered
	StatusDetails string    `json:"status_details"`
	DeliveredAt   *time.Time `json:"delivered_at,omitempty"`
	FailedAt      *time.Time `json:"failed_at,omitempty"`
	ErrorCode     string    `json:"error_code,omitempty"`
	ErrorMessage  string    `json:"error_message,omitempty"`
	Cost          float64   `json:"cost,omitempty"`
}

// SMSService interface for SMS operations
type SMSService interface {
	SendSMS(ctx context.Context, message *SMSMessage) (*SMSResponse, error)
	SendBulkSMS(ctx context.Context, messages []*SMSMessage) ([]*SMSResponse, error)
	GetStatus(ctx context.Context, messageID string) (*SMSStatus, error)
	GetBalance(ctx context.Context) (float64, error)
	ValidateNumber(ctx context.Context, phoneNumber, countryCode string) (bool, error)
	HealthCheck(ctx context.Context) error
}

// SMSManager manages multiple SMS providers with failover
type SMSManager struct {
	providers      map[SMSProvider]SMSService
	primaryConfig  *SMSConfig
	fallbackOrder  []SMSProvider
	countryRouting map[string]SMSProvider
	stats          *SMSStats
}

// SMSStats tracks SMS statistics
type SMSStats struct {
	TotalSent     int64             `json:"total_sent"`
	TotalFailed   int64             `json:"total_failed"`
	TotalCost     float64           `json:"total_cost"`
	ProviderStats map[SMSProvider]*ProviderStats `json:"provider_stats"`
}

// ProviderStats tracks per-provider statistics
type ProviderStats struct {
	Sent        int64   `json:"sent"`
	Failed      int64   `json:"failed"`
	Cost        float64 `json:"cost"`
	AvgLatency  time.Duration `json:"avg_latency"`
	SuccessRate float64 `json:"success_rate"`
	LastUsed    time.Time `json:"last_used"`
}

// NewSMSManager creates a new SMS manager
func NewSMSManager(config *SMSConfig) *SMSManager {
	manager := &SMSManager{
		providers:      make(map[SMSProvider]SMSService),
		primaryConfig:  config,
		fallbackOrder:  []SMSProvider{TwilioProvider, AWSProvider, VonageProvider},
		countryRouting: make(map[string]SMSProvider),
		stats: &SMSStats{
			ProviderStats: make(map[SMSProvider]*ProviderStats),
		},
	}

	// Setup country-specific routing for Southeast Asia
	manager.setupCountryRouting()

	return manager
}

// setupCountryRouting configures country-specific SMS provider routing
func (sm *SMSManager) setupCountryRouting() {
	sm.countryRouting["TH"] = ThaiSMSProvider  // Thailand
	sm.countryRouting["SG"] = SingtelProvider  // Singapore
	sm.countryRouting["ID"] = TelkomselProvider // Indonesia
	sm.countryRouting["MY"] = TwilioProvider   // Malaysia - use Twilio
	sm.countryRouting["PH"] = TwilioProvider   // Philippines - use Twilio
	sm.countryRouting["VN"] = TwilioProvider   // Vietnam - use Twilio
}

// RegisterProvider registers an SMS provider
func (sm *SMSManager) RegisterProvider(provider SMSProvider, service SMSService) {
	sm.providers[provider] = service
	sm.stats.ProviderStats[provider] = &ProviderStats{}
	log.Printf("Registered SMS provider: %s", provider)
}

// SendSMS sends an SMS with automatic provider selection and failover
func (sm *SMSManager) SendSMS(ctx context.Context, message *SMSMessage) (*SMSResponse, error) {
	// Determine best provider for this message
	provider := sm.selectProvider(message)

	// Try primary provider
	if service, exists := sm.providers[provider]; exists {
		response, err := sm.sendWithProvider(ctx, service, provider, message)
		if err == nil {
			sm.updateStats(provider, true, response.Cost, time.Since(response.SentAt))
			return response, nil
		}
		log.Printf("Primary SMS provider %s failed: %v", provider, err)
		sm.updateStats(provider, false, 0, 0)
	}

	// Try fallback providers
	for _, fallbackProvider := range sm.fallbackOrder {
		if fallbackProvider == provider {
			continue // Skip primary provider
		}

		if service, exists := sm.providers[fallbackProvider]; exists {
			response, err := sm.sendWithProvider(ctx, service, fallbackProvider, message)
			if err == nil {
				sm.updateStats(fallbackProvider, true, response.Cost, time.Since(response.SentAt))
				log.Printf("SMS sent successfully using fallback provider: %s", fallbackProvider)
				return response, nil
			}
			log.Printf("Fallback SMS provider %s failed: %v", fallbackProvider, err)
			sm.updateStats(fallbackProvider, false, 0, 0)
		}
	}

	return nil, fmt.Errorf("all SMS providers failed to send message")
}

// sendWithProvider sends SMS using a specific provider
func (sm *SMSManager) sendWithProvider(ctx context.Context, service SMSService, provider SMSProvider, message *SMSMessage) (*SMSResponse, error) {
	startTime := time.Now()

	response, err := service.SendSMS(ctx, message)
	if err != nil {
		return nil, err
	}

	response.Provider = provider
	response.SentAt = startTime
	return response, nil
}

// selectProvider selects the best SMS provider based on country and other factors
func (sm *SMSManager) selectProvider(message *SMSMessage) SMSProvider {
	// Check country-specific routing first
	if message.CountryCode != "" {
		if provider, exists := sm.countryRouting[strings.ToUpper(message.CountryCode)]; exists {
			if _, providerExists := sm.providers[provider]; providerExists {
				return provider
			}
		}
	}

	// Extract country code from phone number if not provided
	if message.CountryCode == "" && len(message.To) > 0 {
		countryCode := sm.extractCountryCode(message.To)
		if provider, exists := sm.countryRouting[countryCode]; exists {
			if _, providerExists := sm.providers[provider]; providerExists {
				return provider
			}
		}
	}

	// Fall back to primary provider
	return sm.primaryConfig.Provider
}

// extractCountryCode extracts country code from phone number
func (sm *SMSManager) extractCountryCode(phoneNumber string) string {
	// Remove leading + if present
	number := strings.TrimPrefix(phoneNumber, "+")

	// Southeast Asian country codes
	countryPrefixes := map[string]string{
		"66":  "TH", // Thailand
		"65":  "SG", // Singapore
		"62":  "ID", // Indonesia
		"60":  "MY", // Malaysia
		"63":  "PH", // Philippines
		"84":  "VN", // Vietnam
	}

	for prefix, country := range countryPrefixes {
		if strings.HasPrefix(number, prefix) {
			return country
		}
	}

	return ""
}

// SendBulkSMS sends multiple SMS messages
func (sm *SMSManager) SendBulkSMS(ctx context.Context, messages []*SMSMessage) ([]*SMSResponse, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages to send")
	}

	responses := make([]*SMSResponse, len(messages))
	errors := make([]error, len(messages))

	// Group messages by provider for batch sending
	providerGroups := make(map[SMSProvider][]*SMSMessage)
	providerIndexes := make(map[SMSProvider][]int)

	for i, message := range messages {
		provider := sm.selectProvider(message)
		providerGroups[provider] = append(providerGroups[provider], message)
		providerIndexes[provider] = append(providerIndexes[provider], i)
	}

	// Send messages in groups by provider
	for provider, groupMessages := range providerGroups {
		if service, exists := sm.providers[provider]; exists {
			groupResponses, err := service.SendBulkSMS(ctx, groupMessages)
			if err != nil {
				// Mark all messages in this group as failed
				for _, index := range providerIndexes[provider] {
					errors[index] = err
				}
				continue
			}

			// Map responses back to original indexes
			for i, response := range groupResponses {
				if i < len(providerIndexes[provider]) {
					originalIndex := providerIndexes[provider][i]
					responses[originalIndex] = response
					response.Provider = provider
				}
			}
		}
	}

	// Check if any messages succeeded
	successCount := 0
	for i, response := range responses {
		if response != nil {
			successCount++
		} else {
			// Create error response for failed messages
			responses[i] = &SMSResponse{
				Status:   "failed",
				Provider: sm.primaryConfig.Provider,
				Error:    "failed to send",
				SentAt:   time.Now(),
			}
		}
	}

	if successCount == 0 {
		return responses, fmt.Errorf("all bulk SMS messages failed")
	}

	log.Printf("Bulk SMS: %d/%d messages sent successfully", successCount, len(messages))
	return responses, nil
}

// GetStatus gets the delivery status of an SMS
func (sm *SMSManager) GetStatus(ctx context.Context, messageID string) (*SMSStatus, error) {
	// Try to get status from all providers (since we don't know which one sent it)
	for provider, service := range sm.providers {
		status, err := service.GetStatus(ctx, messageID)
		if err == nil && status != nil {
			log.Printf("Found SMS status in provider %s: %s", provider, status.Status)
			return status, nil
		}
	}

	return nil, fmt.Errorf("SMS status not found for message ID: %s", messageID)
}

// updateStats updates provider statistics
func (sm *SMSManager) updateStats(provider SMSProvider, success bool, cost float64, latency time.Duration) {
	stats := sm.stats.ProviderStats[provider]
	if stats == nil {
		stats = &ProviderStats{}
		sm.stats.ProviderStats[provider] = stats
	}

	if success {
		stats.Sent++
		sm.stats.TotalSent++
	} else {
		stats.Failed++
		sm.stats.TotalFailed++
	}

	stats.Cost += cost
	sm.stats.TotalCost += cost
	stats.LastUsed = time.Now()

	// Update average latency
	if success && latency > 0 {
		if stats.AvgLatency == 0 {
			stats.AvgLatency = latency
		} else {
			stats.AvgLatency = (stats.AvgLatency + latency) / 2
		}
	}

	// Update success rate
	total := stats.Sent + stats.Failed
	if total > 0 {
		stats.SuccessRate = float64(stats.Sent) / float64(total) * 100
	}
}

// GetStats returns SMS statistics
func (sm *SMSManager) GetStats() *SMSStats {
	return sm.stats
}

// HealthCheck checks the health of all SMS providers
func (sm *SMSManager) HealthCheck(ctx context.Context) error {
	var errors []string

	for provider, service := range sm.providers {
		if err := service.HealthCheck(ctx); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", provider, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("SMS provider health check failures: %s", strings.Join(errors, "; "))
	}

	return nil
}

// SMS message templates for different purposes
var SMSTemplates = map[string]string{
	"otp_verification": "Your Tchat verification code is {{code}}. This code will expire in {{expiry}} minutes. Do not share this code with anyone.",
	"login_alert":      "New login detected on your Tchat account from {{device}} at {{time}}. If this wasn't you, please secure your account immediately.",
	"payment_received": "You received {{amount}} {{currency}} from {{sender}}. Your new balance is {{balance}} {{currency}}.",
	"payment_sent":     "You sent {{amount}} {{currency}} to {{recipient}}. Your new balance is {{balance}} {{currency}}.",
	"welcome":          "Welcome to Tchat! Your account has been successfully created. Start messaging with friends across Southeast Asia.",
}

// FormatTemplate formats an SMS template with variables
func FormatTemplate(template string, variables map[string]string) string {
	message := template
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		message = strings.ReplaceAll(message, placeholder, value)
	}
	return message
}

// ValidatePhoneNumber validates a phone number format
func ValidatePhoneNumber(phoneNumber, countryCode string) bool {
	// Remove spaces and special characters
	number := strings.ReplaceAll(phoneNumber, " ", "")
	number = strings.ReplaceAll(number, "-", "")
	number = strings.ReplaceAll(number, "(", "")
	number = strings.ReplaceAll(number, ")", "")

	// Remove leading + if present
	number = strings.TrimPrefix(number, "+")

	// Check minimum and maximum length
	if len(number) < 8 || len(number) > 15 {
		return false
	}

	// Check if all characters are digits
	for _, char := range number {
		if char < '0' || char > '9' {
			return false
		}
	}

	// Country-specific validation for Southeast Asia
	switch strings.ToUpper(countryCode) {
	case "TH": // Thailand: +66
		return strings.HasPrefix(number, "66") && len(number) == 11
	case "SG": // Singapore: +65
		return strings.HasPrefix(number, "65") && len(number) == 10
	case "ID": // Indonesia: +62
		return strings.HasPrefix(number, "62") && (len(number) >= 11 && len(number) <= 13)
	case "MY": // Malaysia: +60
		return strings.HasPrefix(number, "60") && (len(number) >= 10 && len(number) <= 11)
	case "PH": // Philippines: +63
		return strings.HasPrefix(number, "63") && len(number) == 12
	case "VN": // Vietnam: +84
		return strings.HasPrefix(number, "84") && (len(number) >= 10 && len(number) <= 11)
	}

	return true // Default to valid for other countries
}