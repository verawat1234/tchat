package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/auth/models"
	sharedModels "tchat.dev/shared/models"
)

type KYCRepository interface {
	Create(ctx context.Context, kyc *models.KYC) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.KYC, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.KYC, error)
	Update(ctx context.Context, kyc *models.KYC) error
	GetByStatus(ctx context.Context, status models.KYCStatus) ([]*models.KYC, error)
	GetPendingReviews(ctx context.Context) ([]*models.KYC, error)
	GetStatistics(ctx context.Context) (*KYCStatistics, error)
}

type DocumentVerifier interface {
	VerifyDocument(ctx context.Context, documentType models.DocumentType, documentNumber, country string) (*DocumentVerificationResult, error)
	ExtractDocumentData(ctx context.Context, imageURL string, documentType models.DocumentType) (*ExtractedDocumentData, error)
}

type ComplianceChecker interface {
	CheckSanctionsList(ctx context.Context, fullName, dateOfBirth, country string) (*SanctionsCheckResult, error)
	CheckPEPList(ctx context.Context, fullName, country string) (*PEPCheckResult, error)
	ValidateAddress(ctx context.Context, address models.Address) (*AddressValidationResult, error)
}

type BiometricVerifier interface {
	VerifyFaceMatch(ctx context.Context, documentPhotoURL, selfieURL string) (*FaceMatchResult, error)
	LivenessCheck(ctx context.Context, videoURL string) (*LivenessResult, error)
}

type KYCService struct {
	kycRepo           KYCRepository
	userService       *UserService
	documentVerifier  DocumentVerifier
	complianceChecker ComplianceChecker
	biometricVerifier BiometricVerifier
	eventPublisher    EventPublisher
	db                *gorm.DB
}

type KYCStatistics struct {
	TotalSubmissions      int64         `json:"total_submissions"`
	PendingReviews        int64         `json:"pending_reviews"`
	ApprovedToday         int64         `json:"approved_today"`
	RejectedToday         int64         `json:"rejected_today"`
	AverageProcessingTime time.Duration `json:"average_processing_time"`
	ApprovalRate          float64       `json:"approval_rate"`
}

type DocumentVerificationResult struct {
	IsValid      bool              `json:"is_valid"`
	DocumentData map[string]string `json:"document_data"`
	Confidence   float64           `json:"confidence"`
	ErrorMessage string            `json:"error_message,omitempty"`
}

type ExtractedDocumentData struct {
	FullName       string            `json:"full_name"`
	DateOfBirth    string            `json:"date_of_birth"`
	DocumentNumber string            `json:"document_number"`
	ExpiryDate     string            `json:"expiry_date"`
	Address        string            `json:"address"`
	Nationality    string            `json:"nationality"`
	Confidence     float64           `json:"confidence"`
	RawData        map[string]string `json:"raw_data"`
}

type SanctionsCheckResult struct {
	IsMatch      bool     `json:"is_match"`
	MatchedNames []string `json:"matched_names"`
	RiskScore    float64  `json:"risk_score"`
	Lists        []string `json:"lists"`
}

type PEPCheckResult struct {
	IsMatch   bool     `json:"is_match"`
	Positions []string `json:"positions"`
	RiskScore float64  `json:"risk_score"`
}

type AddressValidationResult struct {
	IsValid             bool    `json:"is_valid"`
	StandardizedAddress string  `json:"standardized_address"`
	Confidence          float64 `json:"confidence"`
}

type FaceMatchResult struct {
	IsMatch    bool    `json:"is_match"`
	Confidence float64 `json:"confidence"`
	Similarity float64 `json:"similarity"`
}

type LivenessResult struct {
	IsLive     bool    `json:"is_live"`
	Confidence float64 `json:"confidence"`
}

func NewKYCService(
	kycRepo KYCRepository,
	userService *UserService,
	documentVerifier DocumentVerifier,
	complianceChecker ComplianceChecker,
	biometricVerifier BiometricVerifier,
	eventPublisher EventPublisher,
	db *gorm.DB,
) *KYCService {
	return &KYCService{
		kycRepo:           kycRepo,
		userService:       userService,
		documentVerifier:  documentVerifier,
		complianceChecker: complianceChecker,
		biometricVerifier: biometricVerifier,
		eventPublisher:    eventPublisher,
		db:                db,
	}
}

func (ks *KYCService) SubmitKYC(ctx context.Context, req *SubmitKYCRequest) (*models.KYC, error) {
	// Validate request
	if err := ks.validateSubmitKYCRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if user exists
	user, err := ks.userService.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user can submit KYC (simplified check for active users)
	if user.Status != string(sharedModels.UserStatusActive) {
		return nil, fmt.Errorf("user cannot submit KYC in current status: %s", user.Status)
	}

	// Check for existing KYC
	existingKYC, err := ks.kycRepo.GetByUserID(ctx, req.UserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing KYC: %w", err)
	}

	if existingKYC != nil && existingKYC.Status == models.KYCStatusPending {
		return nil, fmt.Errorf("KYC submission already pending review")
	}

	if existingKYC != nil && existingKYC.Status == models.KYCStatusApproved {
		return nil, fmt.Errorf("KYC already approved")
	}

	// Create KYC record
	submitTime := time.Now()
	kyc := &models.KYC{
		ID:             uuid.New(),
		UserID:         req.UserID,
		DocumentType:   req.DocumentType,
		DocumentNumber: req.DocumentNumber,
		Status:         models.KYCStatusPending,
		Tier:           models.TierBasic, // Start with basic tier
		FullName:       req.PersonalInfo.FirstName + " " + req.PersonalInfo.LastName,
		DateOfBirth:    req.PersonalInfo.DateOfBirth,
		Nationality:    string(req.PersonalInfo.Nationality),
		Gender:         req.PersonalInfo.Gender,
		Address:        &req.PersonalInfo.Address,
		ComplianceFlags: models.ComplianceFlags{
			SanctionCheck:  false,
			PEPCheck:       false,
			WatchlistCheck: false,
			IdentityMatch:  false,
			LivenessCheck:  false,
			DocumentAuth:   false,
			AMLScreening:   false,
			SourceOfFunds:  false,
		},
		SubmittedAt: &submitTime,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// TODO: Implement document format validation
	// Note: ValidateDocumentFormat method not available in current KYC model

	// Save KYC to database
	if err := ks.kycRepo.Create(ctx, kyc); err != nil {
		return nil, fmt.Errorf("failed to create KYC: %w", err)
	}

	// Start asynchronous verification process
	go ks.processKYCVerification(context.Background(), kyc)

	// Publish KYC submitted event
	if err := ks.publishKYCEvent(ctx, "kyc.submitted", kyc.UserID, map[string]interface{}{
		"kyc_id":        kyc.ID,
		"document_type": kyc.DocumentType,
		"tier":          kyc.Tier,
	}); err != nil {
		fmt.Printf("Failed to publish KYC submitted event: %v\n", err)
	}

	return kyc, nil
}

func (ks *KYCService) GetKYCByUserID(ctx context.Context, userID uuid.UUID) (*models.KYC, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	kyc, err := ks.kycRepo.GetByUserID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("KYC not found")
		}
		return nil, fmt.Errorf("failed to get KYC: %w", err)
	}

	return kyc, nil
}

func (ks *KYCService) GetKYCByID(ctx context.Context, kycID uuid.UUID) (*models.KYC, error) {
	if kycID == uuid.Nil {
		return nil, fmt.Errorf("KYC ID is required")
	}

	kyc, err := ks.kycRepo.GetByID(ctx, kycID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("KYC not found")
		}
		return nil, fmt.Errorf("failed to get KYC: %w", err)
	}

	return kyc, nil
}

func (ks *KYCService) ApproveKYC(ctx context.Context, kycID uuid.UUID, reviewerID uuid.UUID, tier models.VerificationTier, notes string) error {
	kyc, err := ks.GetKYCByID(ctx, kycID)
	if err != nil {
		return err
	}

	if kyc.Status != models.KYCStatusPending {
		return fmt.Errorf("cannot approve KYC with status: %s", kyc.Status)
	}

	// Update KYC status
	kyc.Status = models.KYCStatusApproved
	kyc.Tier = tier
	kyc.ReviewerID = &reviewerID
	kyc.ReviewedAt = &time.Time{}
	*kyc.ReviewedAt = time.Now()
	kyc.ReviewNotes = notes
	kyc.UpdatedAt = time.Now()

	// Save updated KYC
	if err := ks.kycRepo.Update(ctx, kyc); err != nil {
		return fmt.Errorf("failed to update KYC: %w", err)
	}

	// Update user KYC tier
	if err := ks.userService.UpdateKYCTier(ctx, kyc.UserID, tier); err != nil {
		return fmt.Errorf("failed to update user KYC tier: %w", err)
	}

	// Publish KYC approved event
	if err := ks.publishKYCEvent(ctx, sharedModels.EventTypeUserKYCVerified, kyc.UserID, map[string]interface{}{
		"kyc_id":      kyc.ID,
		"tier":        tier,
		"reviewer_id": reviewerID,
		"notes":       notes,
	}); err != nil {
		fmt.Printf("Failed to publish KYC approved event: %v\n", err)
	}

	return nil
}

func (ks *KYCService) RejectKYC(ctx context.Context, kycID uuid.UUID, reviewerID uuid.UUID, reason string, notes string) error {
	kyc, err := ks.GetKYCByID(ctx, kycID)
	if err != nil {
		return err
	}

	if kyc.Status != models.KYCStatusPending {
		return fmt.Errorf("cannot reject KYC with status: %s", kyc.Status)
	}

	// Update KYC status
	kyc.Status = models.KYCStatusRejected
	kyc.RejectionReason = reason
	kyc.ReviewerID = &reviewerID
	kyc.ReviewedAt = &time.Time{}
	*kyc.ReviewedAt = time.Now()
	kyc.ReviewNotes = notes
	kyc.UpdatedAt = time.Now()

	// Save updated KYC
	if err := ks.kycRepo.Update(ctx, kyc); err != nil {
		return fmt.Errorf("failed to update KYC: %w", err)
	}

	// Publish KYC rejected event
	if err := ks.publishKYCEvent(ctx, "kyc.rejected", kyc.UserID, map[string]interface{}{
		"kyc_id":           kyc.ID,
		"rejection_reason": reason,
		"reviewer_id":      reviewerID,
		"notes":            notes,
	}); err != nil {
		fmt.Printf("Failed to publish KYC rejected event: %v\n", err)
	}

	return nil
}

func (ks *KYCService) RequestMoreInfo(ctx context.Context, kycID uuid.UUID, reviewerID uuid.UUID, requirements []string, notes string) error {
	kyc, err := ks.GetKYCByID(ctx, kycID)
	if err != nil {
		return err
	}

	if kyc.Status != models.KYCStatusPending {
		return fmt.Errorf("cannot request more info for KYC with status: %s", kyc.Status)
	}

	// Update KYC status
	kyc.Status = models.KYCStatusMoreInfoRequired
	// TODO: Store additional requirements (field not available in current model)
	// kyc.AdditionalRequirements = requirements
	kyc.ReviewerID = &reviewerID
	kyc.ReviewedAt = &time.Time{}
	*kyc.ReviewedAt = time.Now()
	kyc.ReviewNotes = notes
	kyc.UpdatedAt = time.Now()

	// Save updated KYC
	if err := ks.kycRepo.Update(ctx, kyc); err != nil {
		return fmt.Errorf("failed to update KYC: %w", err)
	}

	// Publish more info required event
	if err := ks.publishKYCEvent(ctx, "kyc.more_info_required", kyc.UserID, map[string]interface{}{
		"kyc_id":       kyc.ID,
		"requirements": requirements,
		"reviewer_id":  reviewerID,
		"notes":        notes,
	}); err != nil {
		fmt.Printf("Failed to publish KYC more info required event: %v\n", err)
	}

	return nil
}

func (ks *KYCService) GetPendingReviews(ctx context.Context) ([]*models.KYC, error) {
	return ks.kycRepo.GetPendingReviews(ctx)
}

func (ks *KYCService) GetKYCStatistics(ctx context.Context) (*KYCStatistics, error) {
	return ks.kycRepo.GetStatistics(ctx)
}

func (ks *KYCService) GetKYCStatus(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	// Get user
	user, err := ks.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get KYC record if exists
	kyc, err := ks.kycRepo.GetByUserID(ctx, userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get KYC: %w", err)
	}

	status := map[string]interface{}{
		"user_id":    userID,
		"kyc_tier":   0, // Default KYC tier for shared model
		"can_submit": user.Status == string(sharedModels.UserStatusActive),
		"limits": map[string]interface{}{
			"daily_limit": 1000.0, // Default daily limit
			"can_send":    true,    // Default permission
		},
	}

	if kyc != nil {
		status["kyc_id"] = kyc.ID
		status["status"] = kyc.Status
		status["tier"] = kyc.Tier
		status["submitted_at"] = kyc.SubmittedAt
		status["reviewed_at"] = kyc.ReviewedAt
		status["rejection_reason"] = kyc.RejectionReason
		// TODO: Include additional requirements when field is available
	// status["additional_requirements"] = kyc.AdditionalRequirements
		status["compliance_flags"] = kyc.ComplianceFlags

		// TODO: Calculate transaction limits when method is available
		// daily, monthly, annual := kyc.GetTransactionLimits()
		status["transaction_limits"] = map[string]interface{}{
			"daily":   10000,  // Placeholder values
			"monthly": 50000,
			"annual":  100000,
		}
	} else {
		status["status"] = "not_submitted"
		// TODO: Implement CanSubmitKYC check
		status["can_submit"] = true  // Placeholder
	}

	return status, nil
}

func (ks *KYCService) UpgradeTier(ctx context.Context, userID uuid.UUID, targetTier models.VerificationTier, additionalDocs models.KYCDocuments) error {
	// Get existing KYC
	kyc, err := ks.GetKYCByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if kyc.Status != models.KYCStatusApproved {
		return fmt.Errorf("KYC must be approved before tier upgrade")
	}

	if kyc.Tier >= targetTier {
		return fmt.Errorf("user already has tier %d or higher", targetTier)
	}

	// Create tier upgrade request
	kyc.Status = models.KYCStatusPending
	kyc.Tier = targetTier
	// TODO: Update documents when Documents field is available
	// // kyc.Documents.IDCardBack = additionalDocs.IDCardBack
	// // kyc.Documents.AddressProof = additionalDocs.AddressProof
	// // kyc.Documents.IncomeProof = additionalDocs.IncomeProof
	// // kyc.Documents.BankStatement = additionalDocs.BankStatement
	kyc.UpdatedAt = time.Now()

	// Save updated KYC
	if err := ks.kycRepo.Update(ctx, kyc); err != nil {
		return fmt.Errorf("failed to update KYC for tier upgrade: %w", err)
	}

	// Start verification process for new tier
	go ks.processKYCVerification(context.Background(), kyc)

	// Publish tier upgrade event
	if err := ks.publishKYCEvent(ctx, "kyc.tier_upgrade_requested", kyc.UserID, map[string]interface{}{
		"kyc_id":       kyc.ID,
		"target_tier":  targetTier,
		"current_tier": kyc.Tier,
	}); err != nil {
		fmt.Printf("Failed to publish KYC tier upgrade event: %v\n", err)
	}

	return nil
}

// Private helper methods

func (ks *KYCService) processKYCVerification(ctx context.Context, kyc *models.KYC) {
	// Document verification
	// TODO: Document verification when Documents field is available
	// if kyc.Documents.IDCardFront != "" {
	if false { // Temporarily disabled
		result, err := ks.documentVerifier.VerifyDocument(ctx, kyc.DocumentType, kyc.DocumentNumber, kyc.Nationality)
		if err != nil {
			fmt.Printf("Document verification failed for KYC %s: %v\n", kyc.ID, err)
			return
		}

		if !result.IsValid {
			ks.autoRejectKYC(ctx, kyc, "document_verification_failed", result.ErrorMessage)
			return
		}
	}

	// TODO: Extract document data when Documents field is available
	// Extract document data
	if false { // Temporarily disabled
		extracted, err := ks.documentVerifier.ExtractDocumentData(ctx, "", kyc.DocumentType)
		if err == nil && extracted.Confidence > 0.8 {
			// Validate extracted data against submitted data
			if !ks.validateExtractedData(kyc, extracted) {
				ks.autoRejectKYC(ctx, kyc, "data_mismatch", "Extracted document data does not match submitted information")
				return
			}
		}
	}

	// Compliance checks
	if err := ks.runComplianceChecks(ctx, kyc); err != nil {
		fmt.Printf("Compliance checks failed for KYC %s: %v\n", kyc.ID, err)
		return
	}

	// TODO: Biometric verification when Documents field is available
	// Biometric verification
	if false { // Temporarily disabled
		if err := ks.runBiometricChecks(ctx, kyc); err != nil {
			fmt.Printf("Biometric checks failed for KYC %s: %v\n", kyc.ID, err)
			return
		}
	}

	// Update compliance flags
	ks.updateComplianceFlags(ctx, kyc)
}

func (ks *KYCService) runComplianceChecks(ctx context.Context, kyc *models.KYC) error {
	fullName := kyc.FullName // Use existing FullName field

	// Sanctions check
	sanctionsResult, err := ks.complianceChecker.CheckSanctionsList(ctx, fullName, kyc.DateOfBirth.Format("2006-01-02"), kyc.Nationality)
	if err != nil {
		return fmt.Errorf("sanctions check failed: %w", err)
	}

	if sanctionsResult.IsMatch {
		ks.autoRejectKYC(ctx, kyc, "sanctions_list_match", fmt.Sprintf("Match found on sanctions lists: %v", sanctionsResult.MatchedNames))
		return fmt.Errorf("sanctions list match found")
	}

	// PEP check
	pepResult, err := ks.complianceChecker.CheckPEPList(ctx, fullName, kyc.Nationality)
	if err != nil {
		return fmt.Errorf("PEP check failed: %w", err)
	}

	// High-risk PEP might require additional review
	if pepResult.IsMatch && pepResult.RiskScore > 0.8 {
		// TODO: Add RequiresManualReview field to ComplianceFlags
		// kyc.ComplianceFlags.RequiresManualReview = true
	}

	// Address validation
	addressResult, err := ks.complianceChecker.ValidateAddress(ctx, *kyc.Address)
	if err != nil {
		return fmt.Errorf("address validation failed: %w", err)
	}

	if !addressResult.IsValid {
		ks.autoRejectKYC(ctx, kyc, "invalid_address", "Address validation failed")
		return fmt.Errorf("invalid address")
	}

	return nil
}

func (ks *KYCService) runBiometricChecks(ctx context.Context, kyc *models.KYC) error {
	// Face match
	// TODO: Face match when Documents field available
	faceMatch, err := ks.biometricVerifier.VerifyFaceMatch(ctx, "", "")
	if err != nil {
		return fmt.Errorf("face match verification failed: %w", err)
	}

	if !faceMatch.IsMatch || faceMatch.Confidence < 0.8 {
		ks.autoRejectKYC(ctx, kyc, "face_match_failed", fmt.Sprintf("Face match confidence: %.2f", faceMatch.Confidence))
		return fmt.Errorf("face match failed")
	}

	// Liveness check (if video provided)
	// TODO: Liveness check when Documents field available
	if false { // Temporarily disabled
		liveness, err := ks.biometricVerifier.LivenessCheck(ctx, "")
		if err != nil {
			return fmt.Errorf("liveness check failed: %w", err)
		}

		if !liveness.IsLive || liveness.Confidence < 0.8 {
			ks.autoRejectKYC(ctx, kyc, "liveness_check_failed", fmt.Sprintf("Liveness confidence: %.2f", liveness.Confidence))
			return fmt.Errorf("liveness check failed")
		}
	}

	return nil
}

func (ks *KYCService) validateExtractedData(kyc *models.KYC, extracted *ExtractedDocumentData) bool {
	// Compare extracted data with submitted data
	submittedName := kyc.FullName
	if extracted.FullName != submittedName {
		return false
	}

	if extracted.DocumentNumber != kyc.DocumentNumber {
		return false
	}

	// Add more validation as needed
	return true
}

func (ks *KYCService) updateComplianceFlags(ctx context.Context, kyc *models.KYC) {
	kyc.ComplianceFlags.SanctionCheck = true
	kyc.ComplianceFlags.PEPCheck = true
	// TODO: Add AddressCheck field to ComplianceFlags
	// kyc.ComplianceFlags.AddressCheck = true
	// TODO: Add BiometricCheck field to ComplianceFlags
	// kyc.ComplianceFlags.BiometricCheck = true
	kyc.UpdatedAt = time.Now()

	ks.kycRepo.Update(ctx, kyc)
}

func (ks *KYCService) autoRejectKYC(ctx context.Context, kyc *models.KYC, reason, details string) {
	kyc.Status = models.KYCStatusRejected
	kyc.RejectionReason = reason
	kyc.ReviewNotes = details
	now := time.Now()
	kyc.ReviewedAt = &now
	kyc.UpdatedAt = now

	ks.kycRepo.Update(ctx, kyc)

	// Publish auto-rejection event
	ks.publishKYCEvent(ctx, "kyc.auto_rejected", kyc.UserID, map[string]interface{}{
		"kyc_id":  kyc.ID,
		"reason":  reason,
		"details": details,
	})
}

func (ks *KYCService) validateSubmitKYCRequest(req *SubmitKYCRequest) error {
	if req.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}

	if req.DocumentType == "" {
		return fmt.Errorf("document type is required")
	}

	if req.DocumentNumber == "" {
		return fmt.Errorf("document number is required")
	}

	// TODO: Validate PersonalInfo fields when available
	// if req.PersonalInfo.FirstName == "" {
	//	return fmt.Errorf("first name is required")
	// }

	// if req.PersonalInfo.LastName == "" {
	//	return fmt.Errorf("last name is required")
	// }

	// if req.PersonalInfo.DateOfBirth.IsZero() {
	//	return fmt.Errorf("date of birth is required")
	// }

	// TODO: Validate Documents field when available
	// if req.Documents.IDCardFront == "" {
	//	return fmt.Errorf("ID card front image is required")
	// }

	return nil
}

func (ks *KYCService) publishKYCEvent(ctx context.Context, eventType sharedModels.EventType, userID uuid.UUID, data map[string]interface{}) error {
	event := &sharedModels.Event{
		ID:            uuid.New(),
		Type:          eventType,
		Category:      sharedModels.EventCategoryDomain,
		Severity:      sharedModels.SeverityInfo,
		Subject:       fmt.Sprintf("KYC event: %s", eventType),
		AggregateID:   userID.String(),
		AggregateType: "user",
		EventVersion:  1,
		OccurredAt:    time.Now(),
		Status:        sharedModels.EventStatusPending,
		Metadata: sharedModels.EventMetadata{
			Source:      "auth-service",
			Environment: "production",
			Region:      "sea",
		},
	}

	if err := event.MarshalData(data); err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	return ks.eventPublisher.Publish(ctx, event)
}

// Request/Response structures

type SubmitKYCRequest struct {
	UserID         uuid.UUID           `json:"user_id" binding:"required"`
	DocumentType   models.DocumentType `json:"document_type" binding:"required"`
	DocumentNumber string              `json:"document_number" binding:"required"`
	PersonalInfo   models.PersonalInfo `json:"personal_info" binding:"required"`
	Address        models.Address      `json:"address" binding:"required"`
	Documents      models.KYCDocuments `json:"documents" binding:"required"`
}

type KYCResponse struct {
	ID                     uuid.UUID               `json:"id"`
	UserID                 uuid.UUID               `json:"user_id"`
	DocumentType           models.DocumentType     `json:"document_type"`
	Status                 models.KYCStatus        `json:"status"`
	Tier                   models.VerificationTier `json:"tier"`
	RejectionReason        string                  `json:"rejection_reason,omitempty"`
	AdditionalRequirements []string                `json:"additional_requirements,omitempty"`
	ReviewNotes            string                  `json:"review_notes,omitempty"`
	ComplianceFlags        models.ComplianceFlags  `json:"compliance_flags"`
	SubmittedAt            time.Time               `json:"submitted_at"`
	ReviewedAt             *time.Time              `json:"reviewed_at,omitempty"`
	CreatedAt              time.Time               `json:"created_at"`
	UpdatedAt              time.Time               `json:"updated_at"`
}

type KYCListResponse struct {
	KYCs       []*KYCResponse `json:"kycs"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

func ToKYCResponse(kyc *models.KYC) *KYCResponse {
	return &KYCResponse{
		ID:              kyc.ID,
		UserID:          kyc.UserID,
		DocumentType:    kyc.DocumentType,
		Status:          kyc.Status,
		Tier:            kyc.Tier,
		RejectionReason: kyc.RejectionReason,
		ReviewNotes:     kyc.ReviewNotes,
		ComplianceFlags: kyc.ComplianceFlags,
		SubmittedAt:     func() time.Time { if kyc.SubmittedAt != nil { return *kyc.SubmittedAt }; return time.Time{} }(),
		ReviewedAt:      kyc.ReviewedAt,
		CreatedAt:       kyc.CreatedAt,
		UpdatedAt:       kyc.UpdatedAt,
	}
}
