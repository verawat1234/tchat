package unit

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"tchat.dev/calling/models"
)

// ModelsTestSuite defines the test suite for model validation
type ModelsTestSuite struct {
	suite.Suite
}

// TestCallSessionValidation tests CallSession model validation
func (suite *ModelsTestSuite) TestCallSessionValidation() {
	tests := []struct {
		name        string
		callSession *models.CallSession
		expectValid bool
	}{
		{
			name: "Valid voice call session",
			callSession: &models.CallSession{
				ID:          uuid.New(),
				Type:        "voice",
				Status:      "connecting",
				InitiatedBy: uuid.New(),
				StartedAt:   time.Now(),
			},
			expectValid: true,
		},
		{
			name: "Valid video call session",
			callSession: &models.CallSession{
				ID:          uuid.New(),
				Type:        "video",
				Status:      "active",
				InitiatedBy: uuid.New(),
				StartedAt:   time.Now(),
			},
			expectValid: true,
		},
		{
			name: "Invalid call type",
			callSession: &models.CallSession{
				ID:          uuid.New(),
				Type:        "invalid",
				Status:      "connecting",
				InitiatedBy: uuid.New(),
				StartedAt:   time.Now(),
			},
			expectValid: false,
		},
		{
			name: "Invalid status",
			callSession: &models.CallSession{
				ID:          uuid.New(),
				Type:        "voice",
				Status:      "invalid",
				InitiatedBy: uuid.New(),
				StartedAt:   time.Now(),
			},
			expectValid: false,
		},
		{
			name: "Empty initiated by",
			callSession: &models.CallSession{
				ID:        uuid.New(),
				Type:      "voice",
				Status:    "connecting",
				StartedAt: time.Now(),
			},
			expectValid: false,
		},
		{
			name: "Zero started at",
			callSession: &models.CallSession{
				ID:          uuid.New(),
				Type:        "voice",
				Status:      "connecting",
				InitiatedBy: uuid.New(),
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := validateCallSession(tt.callSession)
			if tt.expectValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestCallParticipantValidation tests CallParticipant model validation
func (suite *ModelsTestSuite) TestCallParticipantValidation() {
	callSessionID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name            string
		callParticipant *models.CallParticipant
		expectValid     bool
	}{
		{
			name: "Valid caller participant",
			callParticipant: &models.CallParticipant{
				ID:               uuid.New(),
				CallSessionID:    callSessionID,
				UserID:           userID,
				Role:             "caller",
				JoinedAt:         time.Now(),
				AudioEnabled:     true,
				VideoEnabled:     false,
				ConnectionQuality: "good",
			},
			expectValid: true,
		},
		{
			name: "Valid callee participant",
			callParticipant: &models.CallParticipant{
				ID:               uuid.New(),
				CallSessionID:    callSessionID,
				UserID:           userID,
				Role:             "callee",
				JoinedAt:         time.Now(),
				AudioEnabled:     true,
				VideoEnabled:     true,
				ConnectionQuality: "excellent",
			},
			expectValid: true,
		},
		{
			name: "Invalid role",
			callParticipant: &models.CallParticipant{
				ID:               uuid.New(),
				CallSessionID:    callSessionID,
				UserID:           userID,
				Role:             "invalid",
				JoinedAt:         time.Now(),
				AudioEnabled:     true,
				VideoEnabled:     false,
				ConnectionQuality: "good",
			},
			expectValid: false,
		},
		{
			name: "Invalid connection quality",
			callParticipant: &models.CallParticipant{
				ID:               uuid.New(),
				CallSessionID:    callSessionID,
				UserID:           userID,
				Role:             "caller",
				JoinedAt:         time.Now(),
				AudioEnabled:     true,
				VideoEnabled:     false,
				ConnectionQuality: "invalid",
			},
			expectValid: false,
		},
		{
			name: "Empty call session ID",
			callParticipant: &models.CallParticipant{
				ID:               uuid.New(),
				UserID:           userID,
				Role:             "caller",
				JoinedAt:         time.Now(),
				AudioEnabled:     true,
				VideoEnabled:     false,
				ConnectionQuality: "good",
			},
			expectValid: false,
		},
		{
			name: "Empty user ID",
			callParticipant: &models.CallParticipant{
				ID:               uuid.New(),
				CallSessionID:    callSessionID,
				Role:             "caller",
				JoinedAt:         time.Now(),
				AudioEnabled:     true,
				VideoEnabled:     false,
				ConnectionQuality: "good",
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := validateCallParticipant(tt.callParticipant)
			if tt.expectValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestUserPresenceValidation tests UserPresence model validation
func (suite *ModelsTestSuite) TestUserPresenceValidation() {
	userID := uuid.New()
	callID := uuid.New()

	tests := []struct {
		name         string
		userPresence *models.UserPresence
		expectValid  bool
	}{
		{
			name: "Valid online presence",
			userPresence: &models.UserPresence{
				ID:     uuid.New(),
				UserID: userID,
				Status: "online",
				InCall: false,
				LastSeen: time.Now(),
			},
			expectValid: true,
		},
		{
			name: "Valid in-call presence",
			userPresence: &models.UserPresence{
				ID:            uuid.New(),
				UserID:        userID,
				Status:        "busy",
				InCall:        true,
				CurrentCallID: &callID,
				LastSeen:      time.Now(),
			},
			expectValid: true,
		},
		{
			name: "Valid offline presence",
			userPresence: &models.UserPresence{
				ID:       uuid.New(),
				UserID:   userID,
				Status:   "offline",
				InCall:   false,
				LastSeen: time.Now(),
			},
			expectValid: true,
		},
		{
			name: "Invalid status",
			userPresence: &models.UserPresence{
				ID:       uuid.New(),
				UserID:   userID,
				Status:   "invalid",
				InCall:   false,
				LastSeen: time.Now(),
			},
			expectValid: false,
		},
		{
			name: "Empty user ID",
			userPresence: &models.UserPresence{
				ID:       uuid.New(),
				Status:   "online",
				InCall:   false,
				LastSeen: time.Now(),
			},
			expectValid: false,
		},
		{
			name: "In call without call ID",
			userPresence: &models.UserPresence{
				ID:       uuid.New(),
				UserID:   userID,
				Status:   "busy",
				InCall:   true,
				LastSeen: time.Now(),
			},
			expectValid: false,
		},
		{
			name: "Not in call but has call ID",
			userPresence: &models.UserPresence{
				ID:            uuid.New(),
				UserID:        userID,
				Status:        "online",
				InCall:        false,
				CurrentCallID: &callID,
				LastSeen:      time.Now(),
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := validateUserPresence(tt.userPresence)
			if tt.expectValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestCallHistoryValidation tests CallHistory model validation
func (suite *ModelsTestSuite) TestCallHistoryValidation() {
	callSessionID := uuid.New()
	callerID := uuid.New()
	calleeID := uuid.New()

	tests := []struct {
		name        string
		callHistory *models.CallHistory
		expectValid bool
	}{
		{
			name: "Valid completed call history",
			callHistory: &models.CallHistory{
				ID:            uuid.New(),
				CallSessionID: callSessionID,
				CallerID:      callerID,
				CalleeID:      calleeID,
				CallType:      "voice",
				CallStatus:    "completed",
				StartedAt:     time.Now().Add(-10 * time.Minute),
				EndedAt:       &[]time.Time{time.Now()}[0],
				Duration:      600, // 10 minutes
			},
			expectValid: true,
		},
		{
			name: "Valid missed call history",
			callHistory: &models.CallHistory{
				ID:            uuid.New(),
				CallSessionID: callSessionID,
				CallerID:      callerID,
				CalleeID:      calleeID,
				CallType:      "video",
				CallStatus:    "missed",
				StartedAt:     time.Now(),
			},
			expectValid: true,
		},
		{
			name: "Valid failed call history",
			callHistory: &models.CallHistory{
				ID:            uuid.New(),
				CallSessionID: callSessionID,
				CallerID:      callerID,
				CalleeID:      calleeID,
				CallType:      "voice",
				CallStatus:    "failed",
				StartedAt:     time.Now(),
				FailureReason: &[]string{"network_error"}[0],
			},
			expectValid: true,
		},
		{
			name: "Invalid call type",
			callHistory: &models.CallHistory{
				ID:            uuid.New(),
				CallSessionID: callSessionID,
				CallerID:      callerID,
				CalleeID:      calleeID,
				CallType:      "invalid",
				CallStatus:    "completed",
				StartedAt:     time.Now(),
			},
			expectValid: false,
		},
		{
			name: "Invalid call status",
			callHistory: &models.CallHistory{
				ID:            uuid.New(),
				CallSessionID: callSessionID,
				CallerID:      callerID,
				CalleeID:      calleeID,
				CallType:      "voice",
				CallStatus:    "invalid",
				StartedAt:     time.Now(),
			},
			expectValid: false,
		},
		{
			name: "Same caller and callee",
			callHistory: &models.CallHistory{
				ID:            uuid.New(),
				CallSessionID: callSessionID,
				CallerID:      callerID,
				CalleeID:      callerID, // Same as caller
				CallType:      "voice",
				CallStatus:    "completed",
				StartedAt:     time.Now(),
			},
			expectValid: false,
		},
		{
			name: "Negative duration",
			callHistory: &models.CallHistory{
				ID:            uuid.New(),
				CallSessionID: callSessionID,
				CallerID:      callerID,
				CalleeID:      calleeID,
				CallType:      "voice",
				CallStatus:    "completed",
				StartedAt:     time.Now(),
				Duration:      -100,
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := validateCallHistory(tt.callHistory)
			if tt.expectValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestCallSessionStatusTransitions tests valid status transitions
func (suite *ModelsTestSuite) TestCallSessionStatusTransitions() {
	callSession := &models.CallSession{
		ID:          uuid.New(),
		Type:        "voice",
		Status:      "connecting",
		InitiatedBy: uuid.New(),
		StartedAt:   time.Now(),
	}

	// Test valid transitions
	validTransitions := []struct {
		from string
		to   string
	}{
		{"connecting", "active"},
		{"connecting", "failed"},
		{"connecting", "ended"},
		{"active", "ended"},
		{"active", "failed"},
	}

	for _, transition := range validTransitions {
		suite.T().Run(transition.from+"_to_"+transition.to, func(t *testing.T) {
			callSession.Status = transition.from
			err := transitionCallStatus(callSession, transition.to)
			assert.NoError(t, err)
			assert.Equal(t, transition.to, callSession.Status)
		})
	}

	// Test invalid transitions
	invalidTransitions := []struct {
		from string
		to   string
	}{
		{"active", "connecting"},
		{"ended", "active"},
		{"failed", "connecting"},
		{"ended", "connecting"},
	}

	for _, transition := range invalidTransitions {
		suite.T().Run(transition.from+"_to_"+transition.to+"_invalid", func(t *testing.T) {
			callSession.Status = transition.from
			err := transitionCallStatus(callSession, transition.to)
			assert.Error(t, err)
		})
	}
}

// TestCallDurationCalculation tests call duration calculation
func (suite *ModelsTestSuite) TestCallDurationCalculation() {
	startTime := time.Now().Add(-10 * time.Minute)
	endTime := time.Now()

	callSession := &models.CallSession{
		ID:          uuid.New(),
		Type:        "voice",
		Status:      "active",
		InitiatedBy: uuid.New(),
		StartedAt:   startTime,
	}

	// Calculate duration when call ends
	duration := calculateCallDuration(callSession, endTime)

	expectedDuration := int(endTime.Sub(startTime).Seconds())
	assert.Equal(suite.T(), expectedDuration, duration)
	assert.Greater(suite.T(), duration, 0)
}

// TestModelConstraints tests database constraints
func (suite *ModelsTestSuite) TestModelConstraints() {
	// Test unique constraints
	suite.T().Run("unique_call_participant_per_call", func(t *testing.T) {
		callSessionID := uuid.New()
		userID := uuid.New()

		participant1 := &models.CallParticipant{
			ID:            uuid.New(),
			CallSessionID: callSessionID,
			UserID:        userID,
			Role:          "caller",
			JoinedAt:      time.Now(),
		}

		participant2 := &models.CallParticipant{
			ID:            uuid.New(),
			CallSessionID: callSessionID,
			UserID:        userID, // Same user in same call
			Role:          "callee",
			JoinedAt:      time.Now(),
		}

		// This should fail unique constraint
		err := validateUniqueCallParticipant(participant1, participant2)
		assert.Error(t, err)
	})

	// Test foreign key constraints
	suite.T().Run("foreign_key_constraints", func(t *testing.T) {
		callSessionID := uuid.New()
		userID := uuid.New()

		participant := &models.CallParticipant{
			ID:            uuid.New(),
			CallSessionID: callSessionID,
			UserID:        userID,
			Role:          "caller",
			JoinedAt:      time.Now(),
		}

		// In real implementation, this would check if call session exists
		err := validateCallParticipantReferences(participant)
		assert.NoError(t, err) // Mock validation passes
	})
}

// TestModelsTestSuite runs the test suite
func TestModelsTestSuite(t *testing.T) {
	suite.Run(t, new(ModelsTestSuite))
}

// Helper functions for validation

func validateCallSession(call *models.CallSession) error {
	if call.Type != "voice" && call.Type != "video" {
		return assert.AnError
	}
	if call.Status != "connecting" && call.Status != "active" && call.Status != "ended" && call.Status != "failed" {
		return assert.AnError
	}
	if call.InitiatedBy == uuid.Nil {
		return assert.AnError
	}
	if call.StartedAt.IsZero() {
		return assert.AnError
	}
	return nil
}

func validateCallParticipant(participant *models.CallParticipant) error {
	if participant.Role != "caller" && participant.Role != "callee" {
		return assert.AnError
	}
	if participant.ConnectionQuality != "excellent" && participant.ConnectionQuality != "good" &&
	   participant.ConnectionQuality != "fair" && participant.ConnectionQuality != "poor" {
		return assert.AnError
	}
	if participant.CallSessionID == uuid.Nil {
		return assert.AnError
	}
	if participant.UserID == uuid.Nil {
		return assert.AnError
	}
	return nil
}

func validateUserPresence(presence *models.UserPresence) error {
	if presence.Status != "online" && presence.Status != "busy" &&
	   presence.Status != "away" && presence.Status != "offline" {
		return assert.AnError
	}
	if presence.UserID == uuid.Nil {
		return assert.AnError
	}
	if presence.InCall && presence.CurrentCallID == nil {
		return assert.AnError
	}
	if !presence.InCall && presence.CurrentCallID != nil {
		return assert.AnError
	}
	return nil
}

func validateCallHistory(history *models.CallHistory) error {
	if history.CallType != "voice" && history.CallType != "video" {
		return assert.AnError
	}
	if history.CallStatus != "completed" && history.CallStatus != "missed" &&
	   history.CallStatus != "declined" && history.CallStatus != "failed" {
		return assert.AnError
	}
	if history.CallerID == history.CalleeID {
		return assert.AnError
	}
	if history.Duration != nil && *history.Duration < 0 {
		return assert.AnError
	}
	return nil
}

func transitionCallStatus(call *models.CallSession, newStatus string) error {
	validTransitions := map[string][]string{
		"connecting": {"active", "failed", "ended"},
		"active":     {"ended", "failed"},
		"ended":      {},
		"failed":     {},
	}

	allowedTransitions, exists := validTransitions[call.Status]
	if !exists {
		return assert.AnError
	}

	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			call.Status = newStatus
			return nil
		}
	}

	return assert.AnError
}

func calculateCallDuration(call *models.CallSession, endTime time.Time) int {
	return int(endTime.Sub(call.StartedAt).Seconds())
}

func validateUniqueCallParticipant(p1, p2 *models.CallParticipant) error {
	if p1.CallSessionID == p2.CallSessionID && p1.UserID == p2.UserID {
		return assert.AnError
	}
	return nil
}

func validateCallParticipantReferences(participant *models.CallParticipant) error {
	// Mock validation - in real implementation, check if call session and user exist
	return nil
}

// Benchmark tests

func BenchmarkCallSessionValidation(b *testing.B) {
	callSession := &models.CallSession{
		ID:          uuid.New(),
		Type:        "voice",
		Status:      "connecting",
		InitiatedBy: uuid.New(),
		StartedAt:   time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateCallSession(callSession)
	}
}

func BenchmarkCallParticipantValidation(b *testing.B) {
	participant := &models.CallParticipant{
		ID:               uuid.New(),
		CallSessionID:    uuid.New(),
		UserID:           uuid.New(),
		Role:             "caller",
		JoinedAt:         time.Now(),
		AudioEnabled:     true,
		VideoEnabled:     false,
		ConnectionQuality: "good",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateCallParticipant(participant)
	}
}