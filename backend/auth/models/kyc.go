package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// KYCStatus represents the status of KYC verification
type KYCStatus string

const (
	KYCStatusPending          KYCStatus = "pending"
	KYCStatusSubmitted        KYCStatus = "submitted"
	KYCStatusUnderReview      KYCStatus = "under_review"
	KYCStatusMoreInfoRequired KYCStatus = "more_info_required"
	KYCStatusApproved         KYCStatus = "approved"
	KYCStatusRejected         KYCStatus = "rejected"
	KYCStatusExpired          KYCStatus = "expired"
)

// DocumentType represents the type of identity document
type DocumentType string

const (
	DocumentTypeNationalID    DocumentType = "national_id"
	DocumentTypePassport      DocumentType = "passport"
	DocumentTypeDriversLicense DocumentType = "drivers_license"
	DocumentTypeNRIC          DocumentType = "nric"           // Singapore
	DocumentTypeKTP           DocumentType = "ktp"            // Indonesia
	DocumentTypeMyKad         DocumentType = "mykad"          // Malaysia
	DocumentTypeUMID          DocumentType = "umid"           // Philippines
	DocumentTypeCCCD          DocumentType = "cccd"           // Vietnam
)

// VerificationTier represents the level of KYC verification
type VerificationTier int

const (
	TierNone     VerificationTier = 0 // No verification
	TierBasic    VerificationTier = 1 // Phone/Email verification
	TierIdentity VerificationTier = 2 // Identity document verification
	TierEnhanced VerificationTier = 3 // Enhanced due diligence
)

// KYC represents a Know Your Customer verification record
type KYC struct {
	ID               uuid.UUID        `json:"id" gorm:"column:id;primaryKey;type:varchar(36)"`
	UserID           uuid.UUID        `json:"user_id" gorm:"column:user_id;type:varchar(36);not null;index"`
	DocumentType     DocumentType     `json:"document_type" gorm:"column:document_type;type:varchar(20);not null"`
	DocumentNumber   string           `json:"document_number" gorm:"column:document_number;type:varchar(50);not null"`
	FullName         string           `json:"full_name" gorm:"column:full_name;type:varchar(200);not null"`
	DateOfBirth      time.Time        `json:"date_of_birth" gorm:"column:date_of_birth;not null"`
	Nationality      string           `json:"nationality" gorm:"column:nationality;type:varchar(2);not null"` // ISO 3166-1 alpha-2
	PlaceOfBirth     string           `json:"place_of_birth" gorm:"column:place_of_birth;type:varchar(100)"`
	Gender           string           `json:"gender" gorm:"column:gender;type:varchar(10)"`
	Address          *Address         `json:"address,omitempty" gorm:"embedded;embeddedPrefix:address_"`
	DocumentImages   DocumentImages   `json:"document_images" gorm:"column:document_images;type:json"`
	SelfieImage      string           `json:"selfie_image" gorm:"column:selfie_image;type:varchar(500)"`
	Status           KYCStatus        `json:"status" gorm:"column:status;type:varchar(20);default:'pending'"`
	Tier             VerificationTier `json:"tier" gorm:"column:tier;default:0"`
	SubmittedAt      *time.Time       `json:"submitted_at" gorm:"column:submitted_at"`
	ReviewedAt       *time.Time       `json:"reviewed_at" gorm:"column:reviewed_at"`
	ApprovedAt       *time.Time       `json:"approved_at" gorm:"column:approved_at"`
	RejectedAt       *time.Time       `json:"rejected_at" gorm:"column:rejected_at"`
	ExpiresAt        *time.Time       `json:"expires_at" gorm:"column:expires_at"`
	ReviewerID       *uuid.UUID       `json:"reviewer_id,omitempty" gorm:"column:reviewer_id;type:varchar(36)"`
	RejectionReason  string           `json:"rejection_reason,omitempty" gorm:"column:rejection_reason;type:text"`
	ReviewNotes      string           `json:"review_notes,omitempty" gorm:"column:review_notes;type:text"`
	ComplianceFlags  ComplianceFlags  `json:"compliance_flags" gorm:"column:compliance_flags;type:json"`
	CreatedAt        time.Time        `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time        `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// Address represents a physical address for KYC verification
type Address struct {
	Street     string `json:"street" gorm:"column:street;type:varchar(200)"`
	City       string `json:"city" gorm:"column:city;type:varchar(100)"`
	State      string `json:"state" gorm:"column:state;type:varchar(100)"`
	PostalCode string `json:"postal_code" gorm:"column:postal_code;type:varchar(20)"`
	Country    string `json:"country" gorm:"column:country;type:varchar(2)"` // ISO 3166-1 alpha-2
}

// DocumentImages contains URLs to document images
type DocumentImages struct {
	Front string `json:"front"` // Front side of document
	Back  string `json:"back"`  // Back side of document (if applicable)
}

// ComplianceFlags represents various compliance checks
type ComplianceFlags struct {
	PEPCheck        bool   `json:"pep_check"`        // Politically Exposed Person
	SanctionCheck   bool   `json:"sanction_check"`   // Sanctions list check
	WatchlistCheck  bool   `json:"watchlist_check"`  // Watchlist verification
	IdentityMatch   bool   `json:"identity_match"`   // Document-selfie match
	LivenessCheck   bool   `json:"liveness_check"`   // Selfie liveness detection
	DocumentAuth    bool   `json:"document_auth"`    // Document authenticity
	AMLScreening    bool   `json:"aml_screening"`    // Anti-Money Laundering
	SourceOfFunds   bool   `json:"source_of_funds"`  // Source of funds verification
}

// KYCRequest represents a KYC submission request
type KYCRequest struct {
	DocumentType     DocumentType   `json:"document_type" validate:"required"`
	DocumentNumber   string         `json:"document_number" validate:"required,min=5,max=50"`
	FullName         string         `json:"full_name" validate:"required,min=2,max=200"`
	DateOfBirth      string         `json:"date_of_birth" validate:"required"` // YYYY-MM-DD format
	Nationality      string         `json:"nationality" validate:"required,len=2"`
	PlaceOfBirth     string         `json:"place_of_birth" validate:"omitempty,max=100"`
	Gender           string         `json:"gender" validate:"omitempty,oneof=male female other"`
	Address          *Address       `json:"address,omitempty"`
	DocumentImages   DocumentImages `json:"document_images" validate:"required"`
	SelfieImage      string         `json:"selfie_image" validate:"required,url"`
}

// KYCResponse represents the response after KYC submission
type KYCResponse struct {
	ID              uuid.UUID        `json:"id"`
	Status          KYCStatus        `json:"status"`
	Tier            VerificationTier `json:"tier"`
	SubmittedAt     *time.Time       `json:"submitted_at"`
	EstimatedTime   string           `json:"estimated_time,omitempty"`
	RequiredActions []string         `json:"required_actions,omitempty"`
	ExpiresAt       *time.Time       `json:"expires_at,omitempty"`
}

// Country-specific document validation patterns
var documentPatterns = map[string]map[DocumentType]*regexp.Regexp{
	"TH": {
		DocumentTypeNationalID: regexp.MustCompile(`^[0-9]{13}$`),                        // 13 digits
		DocumentTypePassport:   regexp.MustCompile(`^[A-Z][0-9]{7}$`),                   // A1234567
	},
	"SG": {
		DocumentTypeNRIC:     regexp.MustCompile(`^[STFG][0-9]{7}[A-Z]$`),              // S1234567A
		DocumentTypePassport: regexp.MustCompile(`^[A-Z][0-9]{7}$`),                    // A1234567
	},
	"ID": {
		DocumentTypeKTP:      regexp.MustCompile(`^[0-9]{16}$`),                        // 16 digits
		DocumentTypePassport: regexp.MustCompile(`^[A-Z][0-9]{7}$`),                   // A1234567
	},
	"MY": {
		DocumentTypeMyKad:    regexp.MustCompile(`^[0-9]{6}-[0-9]{2}-[0-9]{4}$`),       // 123456-12-3456
		DocumentTypePassport: regexp.MustCompile(`^[A-Z][0-9]{8}$`),                   // A12345678
	},
	"PH": {
		DocumentTypeUMID:     regexp.MustCompile(`^[0-9]{4}-[0-9]{7}-[0-9]$`),          // 1234-1234567-1
		DocumentTypePassport: regexp.MustCompile(`^[A-Z]{2}[0-9]{7}$`),                // AB1234567
	},
	"VN": {
		DocumentTypeCCCD:     regexp.MustCompile(`^[0-9]{12}$`),                        // 12 digits
		DocumentTypePassport: regexp.MustCompile(`^[A-Z][0-9]{7}$`),                   // A1234567
	},
}

// Valid document types for each country
var countryDocuments = map[string][]DocumentType{
	"TH": {DocumentTypeNationalID, DocumentTypePassport, DocumentTypeDriversLicense},
	"SG": {DocumentTypeNRIC, DocumentTypePassport},
	"ID": {DocumentTypeKTP, DocumentTypePassport},
	"MY": {DocumentTypeMyKad, DocumentTypePassport},
	"PH": {DocumentTypeUMID, DocumentTypePassport},
	"VN": {DocumentTypeCCCD, DocumentTypePassport},
}

// NewKYC creates a new KYC record from a request
func NewKYC(userID uuid.UUID, req KYCRequest) (*KYC, error) {
	// Validate the request
	if err := ValidateKYCRequest(req); err != nil {
		return nil, fmt.Errorf("invalid KYC request: %w", err)
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, fmt.Errorf("invalid date of birth format: %w", err)
	}

	kyc := &KYC{
		ID:             uuid.New(),
		UserID:         userID,
		DocumentType:   req.DocumentType,
		DocumentNumber: strings.TrimSpace(req.DocumentNumber),
		FullName:       strings.TrimSpace(req.FullName),
		DateOfBirth:    dob,
		Nationality:    strings.ToUpper(req.Nationality),
		PlaceOfBirth:   strings.TrimSpace(req.PlaceOfBirth),
		Gender:         req.Gender,
		Address:        req.Address,
		DocumentImages: req.DocumentImages,
		SelfieImage:    req.SelfieImage,
		Status:         KYCStatusPending,
		Tier:           TierNone,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	return kyc, nil
}

// ValidateKYCRequest validates a KYC submission request
func ValidateKYCRequest(req KYCRequest) error {
	var errs []string

	// Validate required fields
	if req.DocumentType == "" {
		errs = append(errs, "document_type is required")
	}

	if strings.TrimSpace(req.DocumentNumber) == "" {
		errs = append(errs, "document_number is required")
	}

	if strings.TrimSpace(req.FullName) == "" {
		errs = append(errs, "full_name is required")
	}

	if req.DateOfBirth == "" {
		errs = append(errs, "date_of_birth is required")
	}

	if req.Nationality == "" {
		errs = append(errs, "nationality is required")
	}

	if req.DocumentImages.Front == "" {
		errs = append(errs, "document_front_image is required")
	}

	if req.SelfieImage == "" {
		errs = append(errs, "selfie_image is required")
	}

	// Validate nationality (must be Southeast Asian)
	nationality := strings.ToUpper(req.Nationality)
	validNationality := false
	validCountries := []string{"TH", "SG", "ID", "MY", "PH", "VN"}
	for _, country := range validCountries {
		if country == nationality {
			validNationality = true
			break
		}
	}
	if !validNationality {
		errs = append(errs, "nationality must be a Southeast Asian country")
	}

	// Validate document type for nationality
	if validDocs, exists := countryDocuments[nationality]; exists {
		validDocType := false
		for _, docType := range validDocs {
			if docType == req.DocumentType {
				validDocType = true
				break
			}
		}
		if !validDocType {
			errs = append(errs, fmt.Sprintf("document_type %s not valid for country %s", req.DocumentType, nationality))
		}
	}

	// Validate document number format
	if patterns, exists := documentPatterns[nationality]; exists {
		if pattern, exists := patterns[req.DocumentType]; exists {
			if !pattern.MatchString(req.DocumentNumber) {
				errs = append(errs, fmt.Sprintf("invalid document number format for %s %s", nationality, req.DocumentType))
			}
		}
	}

	// Validate date of birth
	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			errs = append(errs, "invalid date_of_birth format (use YYYY-MM-DD)")
		} else {
			// Must be at least 18 years old
			if time.Since(dob).Hours()/24/365 < 18 {
				errs = append(errs, "user must be at least 18 years old")
			}
			// Cannot be born in the future
			if dob.After(time.Now()) {
				errs = append(errs, "date_of_birth cannot be in the future")
			}
			// Cannot be older than 120 years
			if time.Since(dob).Hours()/24/365 > 120 {
				errs = append(errs, "invalid date_of_birth (too old)")
			}
		}
	}

	// Validate gender if provided
	if req.Gender != "" {
		validGenders := []string{"male", "female", "other"}
		validGender := false
		for _, gender := range validGenders {
			if req.Gender == gender {
				validGender = true
				break
			}
		}
		if !validGender {
			errs = append(errs, "gender must be one of: male, female, other")
		}
	}

	// Validate URLs
	urlPattern := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	if req.DocumentImages.Front != "" && !urlPattern.MatchString(req.DocumentImages.Front) {
		errs = append(errs, "invalid document_front_image URL")
	}
	if req.DocumentImages.Back != "" && !urlPattern.MatchString(req.DocumentImages.Back) {
		errs = append(errs, "invalid document_back_image URL")
	}
	if req.SelfieImage != "" && !urlPattern.MatchString(req.SelfieImage) {
		errs = append(errs, "invalid selfie_image URL")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

// Submit marks the KYC as submitted for review
func (k *KYC) Submit() error {
	if k.Status != KYCStatusPending {
		return errors.New("can only submit pending KYC records")
	}

	now := time.Now().UTC()
	k.Status = KYCStatusSubmitted
	k.SubmittedAt = &now
	k.UpdatedAt = now

	return nil
}

// StartReview moves KYC to under review status
func (k *KYC) StartReview(reviewerID uuid.UUID) error {
	if k.Status != KYCStatusSubmitted {
		return errors.New("can only start review for submitted KYC records")
	}

	now := time.Now().UTC()
	k.Status = KYCStatusUnderReview
	k.ReviewerID = &reviewerID
	k.ReviewedAt = &now
	k.UpdatedAt = now

	return nil
}

// Approve approves the KYC and sets the verification tier
func (k *KYC) Approve(tier VerificationTier, notes string) error {
	if k.Status != KYCStatusUnderReview {
		return errors.New("can only approve KYC records under review")
	}

	now := time.Now().UTC()
	k.Status = KYCStatusApproved
	k.Tier = tier
	k.ApprovedAt = &now
	k.ReviewNotes = notes
	k.UpdatedAt = now

	// Set expiry based on tier (higher tiers last longer)
	switch tier {
	case TierBasic:
		expiresAt := now.Add(365 * 24 * time.Hour) // 1 year
		k.ExpiresAt = &expiresAt
	case TierIdentity:
		expiresAt := now.Add(3 * 365 * 24 * time.Hour) // 3 years
		k.ExpiresAt = &expiresAt
	case TierEnhanced:
		expiresAt := now.Add(5 * 365 * 24 * time.Hour) // 5 years
		k.ExpiresAt = &expiresAt
	}

	return nil
}

// Reject rejects the KYC with a reason
func (k *KYC) Reject(reason string, notes string) error {
	if k.Status != KYCStatusUnderReview {
		return errors.New("can only reject KYC records under review")
	}

	now := time.Now().UTC()
	k.Status = KYCStatusRejected
	k.RejectedAt = &now
	k.RejectionReason = reason
	k.ReviewNotes = notes
	k.UpdatedAt = now

	return nil
}

// IsExpired checks if the KYC verification has expired
func (k *KYC) IsExpired() bool {
	if k.ExpiresAt == nil {
		return false
	}
	return time.Now().UTC().After(*k.ExpiresAt)
}

// IsApproved checks if the KYC is approved and not expired
func (k *KYC) IsApproved() bool {
	return k.Status == KYCStatusApproved && !k.IsExpired()
}

// CanTransact checks if the user can perform transactions based on KYC tier
func (k *KYC) CanTransact() bool {
	return k.IsApproved() && k.Tier >= TierIdentity
}

// GetTransactionLimits returns transaction limits based on KYC tier
func (k *KYC) GetTransactionLimits() (daily, monthly, annual int64) {
	if !k.IsApproved() {
		return 0, 0, 0
	}

	switch k.Tier {
	case TierBasic:
		return 100000, 500000, 2000000 // THB equivalent in cents
	case TierIdentity:
		return 1000000, 5000000, 20000000 // THB equivalent in cents
	case TierEnhanced:
		return 5000000, 25000000, 100000000 // THB equivalent in cents
	default:
		return 0, 0, 0
	}
}

// UpdateComplianceFlags updates the compliance check results
func (k *KYC) UpdateComplianceFlags(flags ComplianceFlags) {
	k.ComplianceFlags = flags
	k.UpdatedAt = time.Now().UTC()
}

// ToResponse converts KYC to response format
func (k *KYC) ToResponse() KYCResponse {
	response := KYCResponse{
		ID:          k.ID,
		Status:      k.Status,
		Tier:        k.Tier,
		SubmittedAt: k.SubmittedAt,
		ExpiresAt:   k.ExpiresAt,
	}

	// Add estimated processing time based on status
	switch k.Status {
	case KYCStatusPending:
		response.EstimatedTime = "Please submit your documents"
		response.RequiredActions = []string{"Upload document images", "Take selfie photo", "Submit for review"}
	case KYCStatusSubmitted:
		response.EstimatedTime = "1-3 business days"
	case KYCStatusUnderReview:
		response.EstimatedTime = "Under review"
	case KYCStatusRejected:
		response.RequiredActions = []string{"Fix issues and resubmit"}
	}

	return response
}

// ToPublicProfile returns sanitized KYC info for public APIs
func (k *KYC) ToPublicProfile() map[string]interface{} {
	return map[string]interface{}{
		"id":           k.ID,
		"status":       k.Status,
		"tier":         k.Tier,
		"submitted_at": k.SubmittedAt,
		"approved_at":  k.ApprovedAt,
		"expires_at":   k.ExpiresAt,
		"can_transact": k.CanTransact(),
	}
}

// KYCFilter represents filters for KYC queries
type KYCFilter struct {
	Status      []KYCStatus        `json:"status"`
	Tier        []VerificationTier `json:"tier"`
	Nationality []string           `json:"nationality"`
	ReviewerID  *uuid.UUID         `json:"reviewer_id"`
	CreatedAfter  *time.Time       `json:"created_after"`
	CreatedBefore *time.Time       `json:"created_before"`
}

// KYCSearchQuery represents a KYC search query
type KYCSearchQuery struct {
	Query  string    `json:"query"`  // Search in name, document number
	Filter KYCFilter `json:"filter"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}

// KYCStats represents KYC processing statistics
type KYCStats struct {
	TotalSubmissions int64            `json:"total_submissions"`
	PendingCount     int64            `json:"pending_count"`
	UnderReviewCount int64            `json:"under_review_count"`
	ApprovedCount    int64            `json:"approved_count"`
	RejectedCount    int64            `json:"rejected_count"`
	ByTier           map[string]int64 `json:"by_tier"`
	ByCountry        map[string]int64 `json:"by_country"`
	AvgProcessingTime string          `json:"avg_processing_time"`
}

// PersonalInfo represents personal information for KYC verification
type PersonalInfo struct {
	FirstName     string    `json:"first_name" binding:"required"`
	LastName      string    `json:"last_name" binding:"required"`
	DateOfBirth   time.Time `json:"date_of_birth" binding:"required"`
	Gender        string    `json:"gender,omitempty"`
	Nationality   string   `json:"nationality" binding:"required"`
	Address       Address   `json:"address" binding:"required"`
	PlaceOfBirth  string    `json:"place_of_birth,omitempty"`
	Occupation    string    `json:"occupation,omitempty"`
	MotherName    string    `json:"mother_name,omitempty"`
	FatherName    string    `json:"father_name,omitempty"`
}

// KYCDocuments represents submitted KYC documents
type KYCDocuments struct {
	PrimaryID   DocumentSubmission   `json:"primary_id" binding:"required"`
	ProofOfAddress DocumentSubmission `json:"proof_of_address,omitempty"`
	Selfie      DocumentSubmission   `json:"selfie" binding:"required"`
	Additional  []DocumentSubmission `json:"additional,omitempty"`
}

// DocumentSubmission represents a submitted document
type DocumentSubmission struct {
	Type         DocumentType `json:"type" binding:"required"`
	Number       string       `json:"number" binding:"required"`
	IssueDate    time.Time    `json:"issue_date"`
	ExpiryDate   time.Time    `json:"expiry_date"`
	IssuingAuth  string       `json:"issuing_authority"`
	FileURL      string       `json:"file_url" binding:"required"`
	Status       string       `json:"status"`
	ReviewNotes  string       `json:"review_notes,omitempty"`
}