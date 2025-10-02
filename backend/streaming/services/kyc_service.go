// backend/streaming/services/kyc_service.go
// KYC Verification Service - Streaming authorization and background monitoring
// Implements T037: KYC verification integration for streaming authorization

package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/shared/models"
	streamModels "tchat.dev/streaming/models"
)

// KYCTier constants matching shared/models/user.go
const (
	KYCTierUnverified = 0 // No verification - video creators allowed
	KYCTierBasic      = 1 // Email + phone verified - store sellers minimum
	KYCTierStandard   = 2 // National ID verified (optional enhancement)
	KYCTierPremium    = 3 // Business registration (optional enhancement)
)

// KYCServiceInterface defines the contract for KYC verification operations
type KYCServiceInterface interface {
	// Core validation methods
	ValidateStoreSellerKYC(userID uuid.UUID) (bool, error)
	ValidateVideoCreatorAuth(userID uuid.UUID) (bool, error)
	GetKYCTier(userID uuid.UUID) (int, error)

	// Background monitoring
	MonitorKYCStatus(streamID, broadcasterID uuid.UUID) error
	TerminateStreamOnKYCFailure(streamID uuid.UUID) error

	// Stream lifecycle integration
	ValidateStreamCreation(broadcasterID uuid.UUID, streamType string) error
	StartKYCMonitoring(ctx context.Context, streamID, broadcasterID uuid.UUID, streamType string) error
	StopKYCMonitoring(streamID uuid.UUID)
}

// KYCService implements KYCServiceInterface
type KYCService struct {
	db                *gorm.DB
	monitoringTicker  *time.Ticker
	activeMonitoring  map[uuid.UUID]context.CancelFunc // streamID -> cancel function
	monitoringEnabled bool
}

// NewKYCService creates a new KYC service instance
func NewKYCService(db *gorm.DB) KYCServiceInterface {
	service := &KYCService{
		db:                db,
		activeMonitoring:  make(map[uuid.UUID]context.CancelFunc),
		monitoringEnabled: true,
	}
	return service
}

// ValidateStoreSellerKYC validates that a user meets store seller KYC requirements
// Store sellers require Tier 1+ (email + phone verified) as per FR-021
func (s *KYCService) ValidateStoreSellerKYC(userID uuid.UUID) (bool, error) {
	tier, err := s.GetKYCTier(userID)
	if err != nil {
		return false, fmt.Errorf("failed to get KYC tier: %w", err)
	}

	// Store sellers require minimum Tier 1 (Basic KYC)
	if tier < KYCTierBasic {
		return false, nil
	}

	return true, nil
}

// ValidateVideoCreatorAuth validates that a user meets video creator requirements
// Video creators require email OR phone verification as per FR-022
func (s *KYCService) ValidateVideoCreatorAuth(userID uuid.UUID) (bool, error) {
	var user models.User
	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, fmt.Errorf("user not found: %s", userID)
		}
		return false, fmt.Errorf("database error: %w", err)
	}

	// Video creators need either email OR phone verified (flexible requirement)
	if !user.EmailVerified && !user.PhoneVerified {
		return false, nil
	}

	return true, nil
}

// GetKYCTier returns the KYC verification tier for a user
// Returns 0-3 corresponding to tier levels in shared/models/user.go
func (s *KYCService) GetKYCTier(userID uuid.UUID) (int, error) {
	var user models.User
	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("user not found: %s", userID)
		}
		return 0, fmt.Errorf("database error: %w", err)
	}

	return user.KYCTier, nil
}

// ValidateStreamCreation validates KYC requirements before stream creation
// Called during POST /api/v1/streams validation as per FR-023
func (s *KYCService) ValidateStreamCreation(broadcasterID uuid.UUID, streamType string) error {
	if streamType == "store" {
		// Store sellers require Tier 1+ KYC
		valid, err := s.ValidateStoreSellerKYC(broadcasterID)
		if err != nil {
			return fmt.Errorf("KYC validation error: %w", err)
		}
		if !valid {
			return fmt.Errorf("store sellers require minimum KYC Tier 1 (email + phone verification)")
		}
	} else if streamType == "video" {
		// Video creators require email OR phone verification
		valid, err := s.ValidateVideoCreatorAuth(broadcasterID)
		if err != nil {
			return fmt.Errorf("authentication validation error: %w", err)
		}
		if !valid {
			return fmt.Errorf("video creators require email or phone verification")
		}
	} else {
		return fmt.Errorf("invalid stream type: %s", streamType)
	}

	return nil
}

// MonitorKYCStatus performs a single KYC status check for an active stream
// Called every 5 minutes by background monitoring as per FR-026
func (s *KYCService) MonitorKYCStatus(streamID, broadcasterID uuid.UUID) error {
	// Get stream details
	var stream streamModels.LiveStream
	err := s.db.Where("id = ?", streamID).First(&stream).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("stream not found: %s", streamID)
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Only monitor live streams
	if stream.Status != "live" {
		return nil // Not an error, stream just ended naturally
	}

	// Validate KYC based on stream type
	if stream.StreamType == "store" {
		valid, err := s.ValidateStoreSellerKYC(broadcasterID)
		if err != nil {
			log.Printf("KYC monitoring error for stream %s: %v", streamID, err)
			return err
		}

		// Terminate stream if KYC drops below Tier 1
		if !valid {
			log.Printf("KYC failure detected for stream %s (broadcaster %s) - terminating", streamID, broadcasterID)
			if err := s.TerminateStreamOnKYCFailure(streamID); err != nil {
				return fmt.Errorf("failed to terminate stream: %w", err)
			}
		}
	} else if stream.StreamType == "video" {
		// For video creators, we could monitor email/phone verification status
		// but typically this is less critical than store seller KYC
		valid, err := s.ValidateVideoCreatorAuth(broadcasterID)
		if err != nil {
			log.Printf("Auth monitoring error for stream %s: %v", streamID, err)
			return err
		}

		if !valid {
			log.Printf("Auth failure detected for video stream %s (broadcaster %s) - terminating", streamID, broadcasterID)
			if err := s.TerminateStreamOnKYCFailure(streamID); err != nil {
				return fmt.Errorf("failed to terminate stream: %w", err)
			}
		}
	}

	return nil
}

// TerminateStreamOnKYCFailure terminates a stream due to KYC/auth failure
// Implements graceful stream termination with notification as per FR-026
func (s *KYCService) TerminateStreamOnKYCFailure(streamID uuid.UUID) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":     "terminated",
		"end_time":   now,
		"updated_at": now,
	}

	result := s.db.Model(&streamModels.LiveStream{}).
		Where("id = ?", streamID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("database error during termination: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("stream not found: %s", streamID)
	}

	log.Printf("Stream %s terminated due to KYC/auth failure", streamID)

	// TODO: Send notification to viewers about stream termination
	// This would integrate with the notification service
	// notification.SendStreamTerminationNotice(streamID, "Stream terminated due to verification issues")

	return nil
}

// StartKYCMonitoring starts background monitoring for an active stream
// Checks KYC status every 5 minutes during stream lifetime
func (s *KYCService) StartKYCMonitoring(ctx context.Context, streamID, broadcasterID uuid.UUID, streamType string) error {
	if !s.monitoringEnabled {
		return nil // Monitoring disabled
	}

	// Create cancellable context for this stream
	monitorCtx, cancel := context.WithCancel(ctx)

	// Store cancel function for cleanup
	s.activeMonitoring[streamID] = cancel

	// Start monitoring goroutine
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		log.Printf("Starting KYC monitoring for stream %s (broadcaster %s, type %s)", streamID, broadcasterID, streamType)

		for {
			select {
			case <-ticker.C:
				// Perform KYC check
				if err := s.MonitorKYCStatus(streamID, broadcasterID); err != nil {
					log.Printf("KYC monitoring error for stream %s: %v", streamID, err)
				}

			case <-monitorCtx.Done():
				log.Printf("Stopping KYC monitoring for stream %s", streamID)
				return
			}
		}
	}()

	return nil
}

// StopKYCMonitoring stops background monitoring for a stream
// Called when stream ends or is terminated
func (s *KYCService) StopKYCMonitoring(streamID uuid.UUID) {
	if cancel, exists := s.activeMonitoring[streamID]; exists {
		cancel()
		delete(s.activeMonitoring, streamID)
		log.Printf("KYC monitoring stopped for stream %s", streamID)
	}
}

// KYCValidationResult contains detailed validation information
type KYCValidationResult struct {
	UserID        uuid.UUID `json:"user_id"`
	IsValid       bool      `json:"is_valid"`
	CurrentTier   int       `json:"current_tier"`
	RequiredTier  int       `json:"required_tier"`
	EmailVerified bool      `json:"email_verified"`
	PhoneVerified bool      `json:"phone_verified"`
	Message       string    `json:"message"`
}

// GetDetailedKYCValidation returns detailed KYC validation information
// Useful for API responses and debugging
func (s *KYCService) GetDetailedKYCValidation(userID uuid.UUID, streamType string) (*KYCValidationResult, error) {
	var user models.User
	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %s", userID)
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	result := &KYCValidationResult{
		UserID:        userID,
		CurrentTier:   user.KYCTier,
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
	}

	if streamType == "store" {
		result.RequiredTier = KYCTierBasic
		result.IsValid = user.KYCTier >= KYCTierBasic

		if !result.IsValid {
			result.Message = fmt.Sprintf("Store sellers require KYC Tier %d (current: %d). Please complete email and phone verification.",
				KYCTierBasic, user.KYCTier)
		} else {
			result.Message = "KYC requirements met for store streaming"
		}
	} else if streamType == "video" {
		result.RequiredTier = 0 // No specific tier required
		result.IsValid = user.EmailVerified || user.PhoneVerified

		if !result.IsValid {
			result.Message = "Video creators require email or phone verification"
		} else {
			result.Message = "Authentication requirements met for video streaming"
		}
	} else {
		return nil, fmt.Errorf("invalid stream type: %s", streamType)
	}

	return result, nil
}

// Cleanup performs cleanup of monitoring resources
// Should be called on service shutdown
func (s *KYCService) Cleanup() {
	for streamID, cancel := range s.activeMonitoring {
		cancel()
		log.Printf("Cleaned up KYC monitoring for stream %s", streamID)
	}
	s.activeMonitoring = make(map[uuid.UUID]context.CancelFunc)
}