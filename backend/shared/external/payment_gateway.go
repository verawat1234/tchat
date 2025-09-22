package external

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// PaymentGateway represents different payment gateway providers
type PaymentGateway string

const (
	// Southeast Asian payment gateways
	OmiseGateway      PaymentGateway = "omise"       // Thailand, Japan
	TwoC2PGateway     PaymentGateway = "2c2p"        // Southeast Asia
	MidtransGateway   PaymentGateway = "midtrans"    // Indonesia
	XenditGateway     PaymentGateway = "xendit"      // Indonesia, Philippines
	PayNowGateway     PaymentGateway = "paynow"      // Singapore
	GrabPayGateway    PaymentGateway = "grabpay"     // Southeast Asia
	ShopeepayGateway  PaymentGateway = "shopeepay"   // Southeast Asia
	TrueMoneyGateway  PaymentGateway = "truemoney"   // Thailand
	PromptpayGateway  PaymentGateway = "promptpay"   // Thailand
	DanaGateway       PaymentGateway = "dana"        // Indonesia
	OvoGateway        PaymentGateway = "ovo"         // Indonesia
	GCashGateway      PaymentGateway = "gcash"       // Philippines
	PayMayaGateway    PaymentGateway = "paymaya"     // Philippines

	// International gateways
	StripeGateway     PaymentGateway = "stripe"
	PayPalGateway     PaymentGateway = "paypal"
	AdyenGateway      PaymentGateway = "adyen"
)

// PaymentConfig holds payment gateway configuration
type PaymentConfig struct {
	Gateway           PaymentGateway `mapstructure:"gateway" validate:"required"`
	APIKey            string         `mapstructure:"api_key" validate:"required"`
	SecretKey         string         `mapstructure:"secret_key" validate:"required"`
	MerchantID        string         `mapstructure:"merchant_id"`
	Environment       string         `mapstructure:"environment"` // sandbox, production
	WebhookSecret     string         `mapstructure:"webhook_secret"`
	WebhookURL        string         `mapstructure:"webhook_url"`
	Currency          string         `mapstructure:"currency"`
	CountryCode       string         `mapstructure:"country_code"`
	ReturnURL         string         `mapstructure:"return_url"`
	CancelURL         string         `mapstructure:"cancel_url"`
	Timeout           time.Duration  `mapstructure:"timeout"`
	MaxRetries        int            `mapstructure:"max_retries"`
	RetryDelay        time.Duration  `mapstructure:"retry_delay"`
	EnableLogging     bool           `mapstructure:"enable_logging"`
	SupportedMethods  []string       `mapstructure:"supported_methods"`
}

// DefaultPaymentConfig returns default payment configuration
func DefaultPaymentConfig() *PaymentConfig {
	return &PaymentConfig{
		Gateway:          OmiseGateway,
		Environment:      "sandbox",
		Currency:         "THB",
		CountryCode:      "TH",
		Timeout:          30 * time.Second,
		MaxRetries:       3,
		RetryDelay:       2 * time.Second,
		EnableLogging:    true,
		SupportedMethods: []string{"card", "internet_banking", "mobile_banking", "e_wallet"},
	}
}

// PaymentMethod represents different payment methods
type PaymentMethod string

const (
	CreditCardMethod     PaymentMethod = "credit_card"
	DebitCardMethod      PaymentMethod = "debit_card"
	BankTransferMethod   PaymentMethod = "bank_transfer"
	InternetBankingMethod PaymentMethod = "internet_banking"
	MobileBankingMethod  PaymentMethod = "mobile_banking"
	EWalletMethod        PaymentMethod = "e_wallet"
	QRCodeMethod         PaymentMethod = "qr_code"
	CryptoMethod         PaymentMethod = "crypto"
	InstallmentMethod    PaymentMethod = "installment"
)

// PaymentRequest represents a payment request
type PaymentRequest struct {
	Amount          float64                `json:"amount" validate:"required,gt=0"`
	Currency        string                 `json:"currency" validate:"required"`
	Description     string                 `json:"description" validate:"required"`
	OrderID         string                 `json:"order_id" validate:"required"`
	CustomerID      string                 `json:"customer_id,omitempty"`
	CustomerEmail   string                 `json:"customer_email,omitempty"`
	CustomerPhone   string                 `json:"customer_phone,omitempty"`
	PaymentMethod   PaymentMethod          `json:"payment_method" validate:"required"`
	ReturnURL       string                 `json:"return_url,omitempty"`
	CancelURL       string                 `json:"cancel_url,omitempty"`
	WebhookURL      string                 `json:"webhook_url,omitempty"`
	ExpiresAt       *time.Time             `json:"expires_at,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	BillingAddress  *Address               `json:"billing_address,omitempty"`
	ShippingAddress *Address               `json:"shipping_address,omitempty"`
	Items           []PaymentItem          `json:"items,omitempty"`
	SaveCard        bool                   `json:"save_card,omitempty"`
	CardToken       string                 `json:"card_token,omitempty"`
}

// Address represents billing/shipping address
type Address struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Company     string `json:"company,omitempty"`
	Address1    string `json:"address1,omitempty"`
	Address2    string `json:"address2,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	Country     string `json:"country,omitempty"`
	Phone       string `json:"phone,omitempty"`
}

// PaymentItem represents an item in the payment
type PaymentItem struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description,omitempty"`
	Quantity    int     `json:"quantity" validate:"required,gt=0"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Currency    string  `json:"currency" validate:"required"`
	Category    string  `json:"category,omitempty"`
	SKU         string  `json:"sku,omitempty"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	PaymentID       string                 `json:"payment_id"`
	Status          string                 `json:"status"`
	Gateway         PaymentGateway         `json:"gateway"`
	Amount          float64                `json:"amount"`
	Currency        string                 `json:"currency"`
	OrderID         string                 `json:"order_id"`
	PaymentURL      string                 `json:"payment_url,omitempty"`
	QRCode          string                 `json:"qr_code,omitempty"`
	TransactionID   string                 `json:"transaction_id,omitempty"`
	AuthorizationID string                 `json:"authorization_id,omitempty"`
	Fee             float64                `json:"fee,omitempty"`
	NetAmount       float64                `json:"net_amount,omitempty"`
	ExpiresAt       *time.Time             `json:"expires_at,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	Error           string                 `json:"error,omitempty"`
	ErrorCode       string                 `json:"error_code,omitempty"`
}

// PaymentStatus represents payment status
type PaymentStatus struct {
	PaymentID       string                 `json:"payment_id"`
	Status          string                 `json:"status"` // pending, processing, completed, failed, cancelled, expired, refunded
	Gateway         PaymentGateway         `json:"gateway"`
	Amount          float64                `json:"amount"`
	Currency        string                 `json:"currency"`
	TransactionID   string                 `json:"transaction_id,omitempty"`
	FailureReason   string                 `json:"failure_reason,omitempty"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty"`
	FailedAt        *time.Time             `json:"failed_at,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// RefundRequest represents a refund request
type RefundRequest struct {
	PaymentID   string  `json:"payment_id" validate:"required"`
	Amount      float64 `json:"amount,omitempty"` // Partial refund amount, leave empty for full refund
	Reason      string  `json:"reason,omitempty"`
	Description string  `json:"description,omitempty"`
}

// RefundResponse represents a refund response
type RefundResponse struct {
	RefundID    string    `json:"refund_id"`
	PaymentID   string    `json:"payment_id"`
	Status      string    `json:"status"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Reason      string    `json:"reason,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
	Error       string    `json:"error,omitempty"`
}

// PaymentGatewayService interface for payment operations
type PaymentGatewayService interface {
	CreatePayment(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error)
	GetPaymentStatus(ctx context.Context, paymentID string) (*PaymentStatus, error)
	CapturePayment(ctx context.Context, paymentID string, amount float64) (*PaymentResponse, error)
	CancelPayment(ctx context.Context, paymentID string) error
	RefundPayment(ctx context.Context, request *RefundRequest) (*RefundResponse, error)
	ValidateWebhook(ctx context.Context, payload []byte, signature string) (bool, error)
	GetSupportedMethods(ctx context.Context, countryCode, currency string) ([]PaymentMethod, error)
	HealthCheck(ctx context.Context) error
}

// PaymentGatewayManager manages multiple payment gateways with routing
type PaymentGatewayManager struct {
	gateways       map[PaymentGateway]PaymentGatewayService
	primaryConfig  *PaymentConfig
	routingRules   map[string]PaymentGateway // country -> gateway mapping
	fallbackOrder  []PaymentGateway
	stats          *PaymentStats
}

// PaymentStats tracks payment statistics
type PaymentStats struct {
	TotalPayments   int64                          `json:"total_payments"`
	TotalAmount     float64                        `json:"total_amount"`
	TotalFees       float64                        `json:"total_fees"`
	SuccessRate     float64                        `json:"success_rate"`
	GatewayStats    map[PaymentGateway]*GatewayStats `json:"gateway_stats"`
	MethodStats     map[PaymentMethod]*MethodStats   `json:"method_stats"`
}

// GatewayStats tracks per-gateway statistics
type GatewayStats struct {
	Payments    int64   `json:"payments"`
	Amount      float64 `json:"amount"`
	Fees        float64 `json:"fees"`
	SuccessRate float64 `json:"success_rate"`
	AvgLatency  time.Duration `json:"avg_latency"`
	LastUsed    time.Time     `json:"last_used"`
}

// MethodStats tracks per-method statistics
type MethodStats struct {
	Payments    int64   `json:"payments"`
	Amount      float64 `json:"amount"`
	SuccessRate float64 `json:"success_rate"`
	LastUsed    time.Time `json:"last_used"`
}

// NewPaymentGatewayManager creates a new payment gateway manager
func NewPaymentGatewayManager(config *PaymentConfig) *PaymentGatewayManager {
	manager := &PaymentGatewayManager{
		gateways:      make(map[PaymentGateway]PaymentGatewayService),
		primaryConfig: config,
		routingRules:  make(map[string]PaymentGateway),
		fallbackOrder: []PaymentGateway{StripeGateway, AdyenGateway},
		stats: &PaymentStats{
			GatewayStats: make(map[PaymentGateway]*GatewayStats),
			MethodStats:  make(map[PaymentMethod]*MethodStats),
		},
	}

	// Setup country-specific routing for Southeast Asia
	manager.setupCountryRouting()

	return manager
}

// setupCountryRouting configures country-specific payment gateway routing
func (pm *PaymentGatewayManager) setupCountryRouting() {
	pm.routingRules["TH"] = OmiseGateway      // Thailand - Omise, 2C2P
	pm.routingRules["SG"] = TwoC2PGateway     // Singapore - 2C2P, Stripe
	pm.routingRules["ID"] = MidtransGateway   // Indonesia - Midtrans, Xendit
	pm.routingRules["MY"] = TwoC2PGateway     // Malaysia - 2C2P, Stripe
	pm.routingRules["PH"] = XenditGateway     // Philippines - Xendit, PayMaya
	pm.routingRules["VN"] = TwoC2PGateway     // Vietnam - 2C2P, Adyen
}

// RegisterGateway registers a payment gateway
func (pm *PaymentGatewayManager) RegisterGateway(gateway PaymentGateway, service PaymentGatewayService) {
	pm.gateways[gateway] = service
	pm.stats.GatewayStats[gateway] = &GatewayStats{}
	log.Printf("Registered payment gateway: %s", gateway)
}

// CreatePayment creates a payment with automatic gateway selection
func (pm *PaymentGatewayManager) CreatePayment(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error) {
	// Determine best gateway for this payment
	gateway := pm.selectGateway(request)

	// Try primary gateway
	if service, exists := pm.gateways[gateway]; exists {
		response, err := pm.createWithGateway(ctx, service, gateway, request)
		if err == nil {
			pm.updateStats(gateway, request.PaymentMethod, true, request.Amount, response.Fee)
			return response, nil
		}
		log.Printf("Primary payment gateway %s failed: %v", gateway, err)
		pm.updateStats(gateway, request.PaymentMethod, false, request.Amount, 0)
	}

	// Try fallback gateways
	for _, fallbackGateway := range pm.fallbackOrder {
		if fallbackGateway == gateway {
			continue // Skip primary gateway
		}

		if service, exists := pm.gateways[fallbackGateway]; exists {
			response, err := pm.createWithGateway(ctx, service, fallbackGateway, request)
			if err == nil {
				pm.updateStats(fallbackGateway, request.PaymentMethod, true, request.Amount, response.Fee)
				log.Printf("Payment created successfully using fallback gateway: %s", fallbackGateway)
				return response, nil
			}
			log.Printf("Fallback payment gateway %s failed: %v", fallbackGateway, err)
			pm.updateStats(fallbackGateway, request.PaymentMethod, false, request.Amount, 0)
		}
	}

	return nil, fmt.Errorf("all payment gateways failed to create payment")
}

// createWithGateway creates payment using a specific gateway
func (pm *PaymentGatewayManager) createWithGateway(ctx context.Context, service PaymentGatewayService, gateway PaymentGateway, request *PaymentRequest) (*PaymentResponse, error) {
	startTime := time.Now()

	response, err := service.CreatePayment(ctx, request)
	if err != nil {
		return nil, err
	}

	response.Gateway = gateway
	response.CreatedAt = startTime
	return response, nil
}

// selectGateway selects the best payment gateway based on various factors
func (pm *PaymentGatewayManager) selectGateway(request *PaymentRequest) PaymentGateway {
	// Extract country from customer data or use default routing
	countryCode := pm.extractCountryCode(request)

	// Check country-specific routing
	if countryCode != "" {
		if gateway, exists := pm.routingRules[strings.ToUpper(countryCode)]; exists {
			if _, gatewayExists := pm.gateways[gateway]; gatewayExists {
				return gateway
			}
		}
	}

	// Check payment method specific routing
	gateway := pm.selectByPaymentMethod(request.PaymentMethod, countryCode)
	if gateway != "" {
		if _, gatewayExists := pm.gateways[gateway]; gatewayExists {
			return gateway
		}
	}

	// Fall back to primary gateway
	return pm.primaryConfig.Gateway
}

// selectByPaymentMethod selects gateway based on payment method and country
func (pm *PaymentGatewayManager) selectByPaymentMethod(method PaymentMethod, countryCode string) PaymentGateway {
	// E-wallet routing by country
	if method == EWalletMethod {
		switch strings.ToUpper(countryCode) {
		case "TH":
			return TrueMoneyGateway
		case "SG":
			return GrabPayGateway
		case "ID":
			return OvoGateway
		case "PH":
			return GCashGateway
		}
	}

	// QR code routing by country
	if method == QRCodeMethod {
		switch strings.ToUpper(countryCode) {
		case "TH":
			return PromptpayGateway
		case "SG":
			return PayNowGateway
		case "ID":
			return DanaGateway
		}
	}

	return ""
}

// extractCountryCode extracts country code from payment request
func (pm *PaymentGatewayManager) extractCountryCode(request *PaymentRequest) string {
	// Check billing address first
	if request.BillingAddress != nil && request.BillingAddress.Country != "" {
		return request.BillingAddress.Country
	}

	// Extract from phone number
	if request.CustomerPhone != "" {
		return pm.extractCountryFromPhone(request.CustomerPhone)
	}

	// Use gateway's default country
	return pm.primaryConfig.CountryCode
}

// extractCountryFromPhone extracts country code from phone number
func (pm *PaymentGatewayManager) extractCountryFromPhone(phoneNumber string) string {
	number := strings.TrimPrefix(phoneNumber, "+")

	countryPrefixes := map[string]string{
		"66": "TH", // Thailand
		"65": "SG", // Singapore
		"62": "ID", // Indonesia
		"60": "MY", // Malaysia
		"63": "PH", // Philippines
		"84": "VN", // Vietnam
	}

	for prefix, country := range countryPrefixes {
		if strings.HasPrefix(number, prefix) {
			return country
		}
	}

	return ""
}

// GetPaymentStatus gets payment status from the appropriate gateway
func (pm *PaymentGatewayManager) GetPaymentStatus(ctx context.Context, paymentID string) (*PaymentStatus, error) {
	// Try to get status from all gateways (since we don't know which one processed it)
	for gateway, service := range pm.gateways {
		status, err := service.GetPaymentStatus(ctx, paymentID)
		if err == nil && status != nil {
			status.Gateway = gateway
			return status, nil
		}
	}

	return nil, fmt.Errorf("payment status not found for payment ID: %s", paymentID)
}

// updateStats updates gateway and method statistics
func (pm *PaymentGatewayManager) updateStats(gateway PaymentGateway, method PaymentMethod, success bool, amount, fee float64) {
	// Update gateway stats
	gatewayStats := pm.stats.GatewayStats[gateway]
	if gatewayStats == nil {
		gatewayStats = &GatewayStats{}
		pm.stats.GatewayStats[gateway] = gatewayStats
	}

	gatewayStats.Payments++
	if success {
		gatewayStats.Amount += amount
		gatewayStats.Fees += fee
	}
	gatewayStats.LastUsed = time.Now()

	// Update method stats
	methodStats := pm.stats.MethodStats[method]
	if methodStats == nil {
		methodStats = &MethodStats{}
		pm.stats.MethodStats[method] = methodStats
	}

	methodStats.Payments++
	if success {
		methodStats.Amount += amount
	}
	methodStats.LastUsed = time.Now()

	// Update totals
	pm.stats.TotalPayments++
	if success {
		pm.stats.TotalAmount += amount
		pm.stats.TotalFees += fee
	}
}

// GetStats returns payment statistics
func (pm *PaymentGatewayManager) GetStats() *PaymentStats {
	return pm.stats
}

// HealthCheck checks the health of all payment gateways
func (pm *PaymentGatewayManager) HealthCheck(ctx context.Context) error {
	var errors []string

	for gateway, service := range pm.gateways {
		if err := service.HealthCheck(ctx); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", gateway, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("payment gateway health check failures: %s", strings.Join(errors, "; "))
	}

	return nil
}

// GetSupportedMethods returns supported payment methods for a country and currency
func (pm *PaymentGatewayManager) GetSupportedMethods(ctx context.Context, countryCode, currency string) ([]PaymentMethod, error) {
	gateway := pm.routingRules[strings.ToUpper(countryCode)]
	if gateway == "" {
		gateway = pm.primaryConfig.Gateway
	}

	if service, exists := pm.gateways[gateway]; exists {
		return service.GetSupportedMethods(ctx, countryCode, currency)
	}

	return nil, fmt.Errorf("no payment gateway available for country: %s", countryCode)
}

// Payment method mappings for different regions
var PaymentMethodsByCountry = map[string][]PaymentMethod{
	"TH": {CreditCardMethod, DebitCardMethod, InternetBankingMethod, MobileBankingMethod, EWalletMethod, QRCodeMethod},
	"SG": {CreditCardMethod, DebitCardMethod, BankTransferMethod, EWalletMethod, QRCodeMethod},
	"ID": {CreditCardMethod, DebitCardMethod, BankTransferMethod, EWalletMethod, QRCodeMethod, InstallmentMethod},
	"MY": {CreditCardMethod, DebitCardMethod, InternetBankingMethod, EWalletMethod},
	"PH": {CreditCardMethod, DebitCardMethod, BankTransferMethod, EWalletMethod, InstallmentMethod},
	"VN": {CreditCardMethod, DebitCardMethod, BankTransferMethod, EWalletMethod},
}

// ValidatePaymentRequest validates a payment request
func ValidatePaymentRequest(request *PaymentRequest) error {
	if request.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if request.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	if request.OrderID == "" {
		return fmt.Errorf("order ID is required")
	}

	if request.Description == "" {
		return fmt.Errorf("description is required")
	}

	return nil
}