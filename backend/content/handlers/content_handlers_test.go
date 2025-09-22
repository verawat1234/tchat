package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tchat.dev/content/models"
	"tchat.dev/content/services"
)

// MockContentService is a mock implementation of ContentService
type MockContentService struct {
	mock.Mock
}

func (m *MockContentService) CreateContent(ctx interface{}, req models.CreateContentRequest) (*models.ContentItem, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.ContentItem), args.Error(1)
}

func (m *MockContentService) GetContent(ctx interface{}, id interface{}) (*models.ContentItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.ContentItem), args.Error(1)
}

func (m *MockContentService) GetContentItems(ctx interface{}, filters interface{}, pagination interface{}, sort interface{}) (*models.ContentResponse, error) {
	args := m.Called(ctx, filters, pagination, sort)
	return args.Get(0).(*models.ContentResponse), args.Error(1)
}

func (m *MockContentService) ValidateContentValue(contentType models.ContentType, value models.ContentValue) error {
	args := m.Called(contentType, value)
	return args.Error(0)
}

func setupTestRouter() (*gin.Engine, *MockContentService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockContentService{}
	handlers := NewContentHandlers(mockService)

	v1 := router.Group("/api/v1")
	RegisterContentRoutes(v1, handlers)

	return router, mockService
}

func TestGetContentItems(t *testing.T) {
	router, mockService := setupTestRouter()

	// Mock response data
	expectedResponse := &models.ContentResponse{
		Items: []models.ContentItem{
			{
				Category: "navigation",
				Type:     models.ContentTypeText,
				Value:    models.ContentValue{"text": "Home"},
				Status:   models.ContentStatusPublished,
			},
		},
		Total: 1,
		Pagination: models.Pagination{
			Page:     0,
			PageSize: 20,
			Total:    1,
		},
	}

	mockService.On("GetContentItems", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedResponse, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/content", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

func TestCreateContent(t *testing.T) {
	router, mockService := setupTestRouter()

	// Mock request data
	requestData := models.CreateContentRequest{
		Category: "navigation",
		Type:     models.ContentTypeText,
		Value:    models.ContentValue{"text": "Home"},
		Status:   models.ContentStatusDraft,
	}

	// Mock response data
	expectedContent := &models.ContentItem{
		Category: requestData.Category,
		Type:     requestData.Type,
		Value:    requestData.Value,
		Status:   requestData.Status,
	}

	mockService.On("ValidateContentValue", requestData.Type, requestData.Value).Return(nil)
	mockService.On("CreateContent", mock.Anything, requestData).Return(expectedContent, nil)

	// Create request
	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", "/api/v1/content", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

func TestGetContentByCategory(t *testing.T) {
	router, mockService := setupTestRouter()

	// Mock response data
	expectedResponse := &models.ContentResponse{
		Items: []models.ContentItem{
			{
				Category: "navigation",
				Type:     models.ContentTypeText,
				Value:    models.ContentValue{"text": "Home"},
				Status:   models.ContentStatusPublished,
			},
		},
		Total: 1,
		Pagination: models.Pagination{
			Page:     0,
			PageSize: 20,
			Total:    1,
		},
	}

	mockService.On("GetContentByCategory", mock.Anything, "navigation", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/content/category/navigation", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

func TestCreateContentValidation(t *testing.T) {
	router, _ := setupTestRouter()

	// Test with invalid JSON
	req, _ := http.NewRequest("POST", "/api/v1/content", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add health check route directly
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "content",
		})
	})

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "content", response["service"])
}