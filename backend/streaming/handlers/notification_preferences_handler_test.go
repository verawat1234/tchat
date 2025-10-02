package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"tchat.dev/streaming/models"
)

// MockNotificationPreferenceRepository is a mock implementation
type MockNotificationPreferenceRepository struct {
	mock.Mock
}

func (m *MockNotificationPreferenceRepository) Create(ctx context.Context, pref *models.NotificationPreference) error {
	args := m.Called(ctx, pref)
	return args.Error(0)
}

func (m *MockNotificationPreferenceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.NotificationPreference, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NotificationPreference), args.Error(1)
}

func (m *MockNotificationPreferenceRepository) Update(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	args := m.Called(ctx, userID, updates)
	return args.Error(0)
}

func (m *MockNotificationPreferenceRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockNotificationPreferenceRepository) GetOrCreateDefault(ctx context.Context, userID uuid.UUID) (*models.NotificationPreference, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NotificationPreference), args.Error(1)
}

func (m *MockNotificationPreferenceRepository) Upsert(ctx context.Context, pref *models.NotificationPreference) error {
	args := m.Called(ctx, pref)
	return args.Error(0)
}

func TestHandleGet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockNotificationPreferenceRepository)
	handler := NewNotificationPreferencesHandler(mockRepo)

	userID := uuid.New()
	now := time.Now()
	quietStart := time.Date(2025, 1, 1, 22, 0, 0, 0, time.UTC)
	quietEnd := time.Date(2025, 1, 1, 8, 0, 0, 0, time.UTC)

	prefs := &models.NotificationPreference{
		UserID:              userID,
		PushEnabled:         true,
		EmailEnabled:        true,
		InAppEnabled:        true,
		StoreStreamsEnabled: false,
		VideoStreamsEnabled: false,
		QuietHoursStart:     &quietStart,
		QuietHoursEnd:       &quietEnd,
		Timezone:            sql.NullString{String: "Asia/Bangkok", Valid: true},
		UpdatedAt:           now,
	}

	mockRepo.On("GetOrCreateDefault", mock.Anything, userID).Return(prefs, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.GET("/api/v1/notification-preferences", handler.HandleGet)

	req, _ := http.NewRequest("GET", "/api/v1/notification-preferences", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response NotificationPreferencesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, userID, response.Data.UserID)
	assert.True(t, response.Data.PushEnabled)
	assert.True(t, response.Data.EmailEnabled)
	assert.False(t, response.Data.SMSEnabled)
	assert.True(t, response.Data.InAppEnabled)
	assert.NotNil(t, response.Data.QuietHoursStart)
	assert.Equal(t, "22:00:00", *response.Data.QuietHoursStart)
	assert.Equal(t, "Asia/Bangkok", response.Data.Timezone)

	mockRepo.AssertExpectations(t)
}

func TestHandlePut_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockNotificationPreferenceRepository)
	handler := NewNotificationPreferencesHandler(mockRepo)

	userID := uuid.New()
	now := time.Now()

	existingPrefs := &models.NotificationPreference{
		UserID:              userID,
		PushEnabled:         true,
		EmailEnabled:        true,
		InAppEnabled:        true,
		StoreStreamsEnabled: true,
		VideoStreamsEnabled: true,
		QuietHoursStart:     nil,
		QuietHoursEnd:       nil,
		Timezone:            sql.NullString{String: "UTC", Valid: true},
		UpdatedAt:           now,
	}

	mockRepo.On("GetOrCreateDefault", mock.Anything, userID).Return(existingPrefs, nil)
	mockRepo.On("Update", mock.Anything, userID, mock.Anything).Return(nil)
	mockRepo.On("GetByUserID", mock.Anything, userID).Return(existingPrefs, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.PUT("/api/v1/notification-preferences", handler.HandlePut)

	reqBody := map[string]interface{}{
		"push_enabled": false,
		"quiet_hours_start": "22:00:00",
		"quiet_hours_end": "08:00:00",
		"timezone": "Asia/Bangkok",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/notification-preferences", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response NotificationPreferencesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, "Notification preferences updated successfully", response.Message)

	mockRepo.AssertExpectations(t)
}

func TestHandleGet_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockNotificationPreferenceRepository)
	handler := NewNotificationPreferencesHandler(mockRepo)

	router := gin.New()
	// Don't set user_id to simulate unauthorized request
	router.GET("/api/v1/notification-preferences", handler.HandleGet)

	req, _ := http.NewRequest("GET", "/api/v1/notification-preferences", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["error"].(string), "JWT token")
}

func TestHandlePut_InvalidTimeFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockNotificationPreferenceRepository)
	handler := NewNotificationPreferencesHandler(mockRepo)

	userID := uuid.New()
	now := time.Now()

	existingPrefs := &models.NotificationPreference{
		UserID:              userID,
		PushEnabled:         true,
		EmailEnabled:        true,
		InAppEnabled:        true,
		StoreStreamsEnabled: true,
		VideoStreamsEnabled: true,
		UpdatedAt:           now,
	}

	mockRepo.On("GetOrCreateDefault", mock.Anything, userID).Return(existingPrefs, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	router.PUT("/api/v1/notification-preferences", handler.HandlePut)

	reqBody := map[string]interface{}{
		"quiet_hours_start": "25:00:00", // Invalid hour
		"quiet_hours_end": "08:00:00",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/notification-preferences", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["error"].(string), "Invalid")

	mockRepo.AssertExpectations(t)
}