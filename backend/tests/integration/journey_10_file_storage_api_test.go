// Journey 10: File Management & Storage API Integration Tests
// Comprehensive testing of file upload, storage, CDN distribution, media processing,
// backup systems, and file lifecycle management across Southeast Asian regions

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// AuthenticatedUser represents an authenticated user session
type AuthenticatedUser struct {
	UserID       string `json:"userId"`
	Email        string `json:"email"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Country      string `json:"country"`
	Language     string `json:"language"`
}

// FileMetadata represents file information and metadata
type FileMetadata struct {
	ID               string                 `json:"id,omitempty"`
	OriginalName     string                 `json:"original_name"`
	GeneratedName    string                 `json:"generated_name"`
	MimeType         string                 `json:"mime_type"`
	FileSize         int64                  `json:"file_size"`
	Hash             string                 `json:"hash"`
	StorageProvider  string                 `json:"storage_provider"`
	StoragePath      string                 `json:"storage_path"`
	CDNUrl           string                 `json:"cdn_url"`
	ThumbnailUrl     string                 `json:"thumbnail_url,omitempty"`
	ProcessingStatus string                 `json:"processing_status"`
	Metadata         map[string]interface{} `json:"metadata"`
	Tags             []string               `json:"tags,omitempty"`
	OwnerID          string                 `json:"owner_id"`
	Privacy          string                 `json:"privacy"`
	ExpiresAt        time.Time             `json:"expires_at,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	Region           string                 `json:"region"`
	BackupStatus     string                 `json:"backup_status"`
}

// StorageProvider represents a storage backend configuration
type StorageProvider struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Type            string                 `json:"type"`
	Region          string                 `json:"region"`
	Config          map[string]interface{} `json:"config"`
	Status          string                 `json:"status"`
	Capacity        int64                  `json:"capacity_bytes"`
	UsedSpace       int64                  `json:"used_bytes"`
	Priority        int                    `json:"priority"`
	CostPerGB       float64               `json:"cost_per_gb"`
	Performance     map[string]float64     `json:"performance_metrics"`
	LastHealthCheck time.Time             `json:"last_health_check"`
}

// MediaProcessingJob represents a media processing task
type MediaProcessingJob struct {
	ID          string                 `json:"id,omitempty"`
	FileID      string                 `json:"file_id"`
	JobType     string                 `json:"job_type"`
	Status      string                 `json:"status"`
	Progress    float64               `json:"progress"`
	Parameters  map[string]interface{} `json:"parameters"`
	OutputFiles []string               `json:"output_files,omitempty"`
	StartedAt   time.Time             `json:"started_at,omitempty"`
	CompletedAt time.Time             `json:"completed_at,omitempty"`
	ErrorMsg    string                 `json:"error_message,omitempty"`
	Priority    int                    `json:"priority"`
	EstimatedDuration time.Duration   `json:"estimated_duration"`
	ActualDuration    time.Duration   `json:"actual_duration,omitempty"`
}

// BackupJob represents a file backup operation
type BackupJob struct {
	ID              string                 `json:"id,omitempty"`
	FileID          string                 `json:"file_id"`
	BackupType      string                 `json:"backup_type"`
	SourceRegion    string                 `json:"source_region"`
	DestinationRegion string               `json:"destination_region"`
	Status          string                 `json:"status"`
	Progress        float64               `json:"progress"`
	StartedAt       time.Time             `json:"started_at,omitempty"`
	CompletedAt     time.Time             `json:"completed_at,omitempty"`
	VerificationHash string               `json:"verification_hash,omitempty"`
	RetentionPolicy string                `json:"retention_policy"`
	Cost            float64               `json:"cost"`
	ErrorMsg        string                 `json:"error_message,omitempty"`
}

// CDNDistribution represents CDN configuration and stats
type CDNDistribution struct {
	ID               string            `json:"id"`
	Domain           string            `json:"domain"`
	OriginDomain     string            `json:"origin_domain"`
	Status           string            `json:"status"`
	Regions          []string          `json:"regions"`
	CacheSettings    map[string]interface{} `json:"cache_settings"`
	SecuritySettings map[string]interface{} `json:"security_settings"`
	Analytics        map[string]interface{} `json:"analytics"`
	BandwidthUsage   map[string]float64 `json:"bandwidth_usage"`
	HitRate          float64           `json:"hit_rate"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

// FileAccessLog represents file access tracking
type FileAccessLog struct {
	ID         string                 `json:"id,omitempty"`
	FileID     string                 `json:"file_id"`
	UserID     string                 `json:"user_id,omitempty"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	Region     string                 `json:"region"`
	AccessType string                 `json:"access_type"`
	Timestamp  time.Time             `json:"timestamp"`
	Duration   time.Duration          `json:"duration,omitempty"`
	BytesServed int64                 `json:"bytes_served"`
	StatusCode int                    `json:"status_code"`
	Referer    string                 `json:"referer,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// StorageQuota represents user storage quota and usage
type StorageQuota struct {
	UserID       string            `json:"user_id"`
	TotalQuota   int64            `json:"total_quota_bytes"`
	UsedSpace    int64            `json:"used_bytes"`
	FileCount    int              `json:"file_count"`
	QuotaType    string           `json:"quota_type"`
	ExpiresAt    time.Time        `json:"expires_at,omitempty"`
	Breakdown    map[string]int64 `json:"usage_breakdown"`
	Overage      int64            `json:"overage_bytes"`
	OverageCost  float64          `json:"overage_cost"`
	LastUpdated  time.Time        `json:"last_updated"`
}

// Journey10FileStorageAPISuite tests comprehensive file management and storage systems
type Journey10FileStorageAPISuite struct {
	suite.Suite
	baseURL         string
	httpClient      *http.Client
	user1           *AuthenticatedUser
	user2           *AuthenticatedUser
	admin           *AuthenticatedUser
	testFiles       []FileMetadata
	storageProviders []StorageProvider
	processingJobs   []MediaProcessingJob
	backupJobs       []BackupJob
	tempFiles       []string
}

func (suite *Journey10FileStorageAPISuite) SetupSuite() {
	suite.baseURL = "http://localhost:8081"
	suite.httpClient = &http.Client{Timeout: 60 * time.Second} // Longer timeout for file operations

	// Create test users
	suite.user1 = suite.createTestUser("storage_user1@tchat.com", "password123")
	suite.user2 = suite.createTestUser("storage_user2@tchat.com", "password456")
	suite.admin = suite.createTestUser("storage_admin@tchat.com", "admin789")

	// Set up storage providers
	suite.setupStorageProviders()

	// Create temporary test files
	suite.createTestFiles()

	// Set up storage quotas
	suite.setupStorageQuotas()
}

func (suite *Journey10FileStorageAPISuite) TearDownSuite() {
	suite.cleanupTestFiles()
	suite.cleanupTestData()
}

func (suite *Journey10FileStorageAPISuite) createTestUser(email, password string) *AuthenticatedUser {
	registerData := map[string]interface{}{
		"email":         email,
		"password":      password,
		"firstName":     "Test",
		"lastName":      "User",
		"country":       "TH",
		"language":      "en",
		"storage_quota": 1073741824, // 1GB in bytes
	}

	jsonData, _ := json.Marshal(registerData)
	resp, err := suite.httpClient.Post(
		fmt.Sprintf("%s/api/v1/auth/register", suite.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	suite.NoError(err)
	defer resp.Body.Close()

	var registerResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&registerResponse)
	suite.NoError(err)

	// Check if registration was successful
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		suite.FailNow("Registration failed", "Status: %d, Response: %+v", resp.StatusCode, registerResponse)
	}

	// Safely extract user_id and token with validation
	userID, userIDOk := registerResponse["user_id"].(string)
	token, tokenOk := registerResponse["token"].(string)

	if !userIDOk || !tokenOk {
		suite.FailNow("Invalid registration response", "Expected user_id and token, got: %+v", registerResponse)
	}

	return &AuthenticatedUser{
		UserID:      userID,
		AccessToken: token,
		Email:       email,
	}
}

func (suite *Journey10FileStorageAPISuite) setupStorageProviders() {
	providers := []StorageProvider{
		{
			Name:     "Singapore S3",
			Type:     "s3",
			Region:   "ap-southeast-1",
			Status:   "active",
			Capacity: 10737418240, // 10GB
			Priority: 1,
			CostPerGB: 0.023,
			Config: map[string]interface{}{
				"bucket":     "tchat-files-sg",
				"access_key": "test_access_key",
				"endpoint":   "s3.ap-southeast-1.amazonaws.com",
			},
		},
		{
			Name:     "Thailand CDN Storage",
			Type:     "cdn",
			Region:   "ap-southeast-2",
			Status:   "active",
			Capacity: 5368709120, // 5GB
			Priority: 2,
			CostPerGB: 0.019,
			Config: map[string]interface{}{
				"distribution_id": "test_distribution",
				"origin_domain":   "files-th.tchat.com",
			},
		},
	}

	for _, provider := range providers {
		provider.LastHealthCheck = time.Now()
		suite.storageProviders = append(suite.storageProviders, provider)
	}
}

func (suite *Journey10FileStorageAPISuite) createTestFiles() {
	// Create temporary test files
	testFileContents := map[string]string{
		"test_image.jpg": "fake_jpg_content_for_testing",
		"test_document.pdf": "fake_pdf_content_for_testing",
		"test_video.mp4": "fake_mp4_content_for_testing",
		"test_audio.mp3": "fake_mp3_content_for_testing",
	}

	for filename, content := range testFileContents {
		tempFile, err := os.CreateTemp("", filename)
		if err == nil {
			tempFile.WriteString(content)
			tempFile.Close()
			suite.tempFiles = append(suite.tempFiles, tempFile.Name())
		}
	}
}

func (suite *Journey10FileStorageAPISuite) setupStorageQuotas() {
	quotas := []StorageQuota{
		{
			UserID:     suite.user1.UserID,
			TotalQuota: 1073741824, // 1GB
			QuotaType:  "premium",
			ExpiresAt:  time.Now().AddDate(1, 0, 0),
		},
		{
			UserID:     suite.user2.UserID,
			TotalQuota: 536870912, // 512MB
			QuotaType:  "free",
		},
	}

	for _, quota := range quotas {
		quota.LastUpdated = time.Now()
		jsonData, _ := json.Marshal(quota)
		req, _ := http.NewRequest("PUT",
			fmt.Sprintf("%s/api/v1/storage/quota", suite.baseURL),
			bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		_, _ = suite.httpClient.Do(req)
	}
}

// Test file upload functionality
func (suite *Journey10FileStorageAPISuite) TestFileUploadFunctionality() {
	if len(suite.tempFiles) == 0 {
		suite.T().Skip("No temporary files available for testing")
	}

	// Test single file upload
	testFile := suite.tempFiles[0]
	file, err := os.Open(testFile)
	suite.NoError(err)
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add file field
	part, err := writer.CreateFormFile("file", "test_image.jpg")
	suite.NoError(err)
	io.Copy(part, file)

	// Add metadata fields
	writer.WriteField("privacy", "private")
	writer.WriteField("tags", "test,image,upload")
	writer.WriteField("description", "Test image upload")

	writer.Close()

	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/files/upload", suite.baseURL),
		&body)
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var uploadResponse FileMetadata
	json.NewDecoder(resp.Body).Decode(&uploadResponse)

	suite.NotEmpty(uploadResponse.ID)
	suite.Equal("test_image.jpg", uploadResponse.OriginalName)
	suite.Equal(suite.user1.UserID, uploadResponse.OwnerID)
	suite.Equal("private", uploadResponse.Privacy)
	suite.NotEmpty(uploadResponse.Hash)
	suite.testFiles = append(suite.testFiles, uploadResponse)

	// Test chunked upload for large files
	chunkSize := 1024 * 1024 // 1MB chunks
	uploadSession := map[string]interface{}{
		"filename":    "large_video.mp4",
		"file_size":   5242880, // 5MB
		"chunk_size":  chunkSize,
		"mime_type":   "video/mp4",
		"privacy":     "public",
	}

	jsonData, _ := json.Marshal(uploadSession)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/files/upload/chunked/init", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var sessionResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&sessionResponse)

	sessionID, ok := sessionResponse["session_id"].(string)
	suite.True(ok, "Expected session_id in session response: %+v", sessionResponse)
	suite.NotEmpty(sessionID)

	// Upload first chunk
	chunkData := strings.Repeat("test_chunk_data", 1000) // ~13KB chunk
	req, _ = http.NewRequest("PUT",
		fmt.Sprintf("%s/api/v1/files/upload/chunked/%s/chunk/0", suite.baseURL, sessionID),
		strings.NewReader(chunkData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	// Complete chunked upload
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/files/upload/chunked/%s/complete", suite.baseURL, sessionID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var completeResponse FileMetadata
	json.NewDecoder(resp.Body).Decode(&completeResponse)
	suite.NotEmpty(completeResponse.ID)
}

// Test file download and access control
func (suite *Journey10FileStorageAPISuite) TestFileDownloadAndAccessControl() {
	if len(suite.testFiles) == 0 {
		suite.T().Skip("No test files available")
	}

	fileID := suite.testFiles[0].ID

	// Test authorized download
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/files/%s/download", suite.baseURL, fileID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Contains(resp.Header.Get("Content-Type"), "image")

	// Test unauthorized download
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/files/%s/download", suite.baseURL, fileID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.user2.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusForbidden, resp.StatusCode)

	// Test public file access (no auth required)
	// Update file to public
	updateData := map[string]interface{}{
		"privacy": "public",
	}

	jsonData, _ := json.Marshal(updateData)
	req, _ = http.NewRequest("PATCH",
		fmt.Sprintf("%s/api/v1/files/%s", suite.baseURL, fileID),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	// Test public download without authentication
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/files/%s/download", suite.baseURL, fileID), nil)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	// Test signed URL generation
	urlRequest := map[string]interface{}{
		"expiry_hours": 24,
		"permissions":  []string{"read"},
	}

	jsonData, _ = json.Marshal(urlRequest)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/files/%s/signed-url", suite.baseURL, fileID),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var urlResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&urlResponse)

	suite.Contains(urlResponse, "signed_url")
	suite.Contains(urlResponse, "expires_at")
}

// Test media processing and thumbnails
func (suite *Journey10FileStorageAPISuite) TestMediaProcessingAndThumbnails() {
	if len(suite.testFiles) == 0 {
		suite.T().Skip("No test files available")
	}

	fileID := suite.testFiles[0].ID

	// Request thumbnail generation
	thumbnailJob := MediaProcessingJob{
		FileID:  fileID,
		JobType: "thumbnail",
		Parameters: map[string]interface{}{
			"width":   200,
			"height":  200,
			"quality": 85,
		},
		Priority: 1,
	}

	jsonData, _ := json.Marshal(thumbnailJob)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/files/%s/process", suite.baseURL, fileID),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusAccepted, resp.StatusCode)

	var jobResponse MediaProcessingJob
	json.NewDecoder(resp.Body).Decode(&jobResponse)
	suite.NotEmpty(jobResponse.ID)
	suite.Equal("queued", jobResponse.Status)

	// Check job status
	time.Sleep(2 * time.Second)

	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/files/processing/%s", suite.baseURL, jobResponse.ID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var statusResponse MediaProcessingJob
	json.NewDecoder(resp.Body).Decode(&statusResponse)

	suite.Contains([]string{"queued", "processing", "completed", "failed"}, statusResponse.Status)

	// Test video transcoding
	videoTranscode := MediaProcessingJob{
		FileID:  fileID,
		JobType: "video_transcode",
		Parameters: map[string]interface{}{
			"format":     "mp4",
			"resolution": "720p",
			"bitrate":    "2000k",
			"codec":      "h264",
		},
		Priority: 2,
	}

	jsonData, _ = json.Marshal(videoTranscode)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/files/%s/process", suite.baseURL, fileID),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusAccepted, resp.StatusCode)
}

// Test CDN distribution and caching
func (suite *Journey10FileStorageAPISuite) TestCDNDistributionAndCaching() {
	// Get CDN distribution info
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/storage/cdn/distributions", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var distributionsResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&distributionsResponse)

	suite.Contains(distributionsResponse, "distributions")

	// Test cache invalidation
	if len(suite.testFiles) > 0 {
		invalidateRequest := map[string]interface{}{
			"paths":        []string{"/files/" + suite.testFiles[0].ID},
			"invalidation_type": "immediate",
		}

		jsonData, _ := json.Marshal(invalidateRequest)
		req, _ = http.NewRequest("POST",
			fmt.Sprintf("%s/api/v1/storage/cdn/invalidate", suite.baseURL),
			bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err = suite.httpClient.Do(req)
		suite.NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusAccepted, resp.StatusCode)

		var invalidateResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&invalidateResponse)
		suite.Contains(invalidateResponse, "invalidation_id")
	}

	// Test CDN analytics
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/storage/cdn/analytics?period=24h", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var analyticsResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&analyticsResponse)

	suite.Contains(analyticsResponse, "bandwidth_usage")
	suite.Contains(analyticsResponse, "requests_count")
	suite.Contains(analyticsResponse, "hit_rate")
}

// Test backup and disaster recovery
func (suite *Journey10FileStorageAPISuite) TestBackupAndDisasterRecovery() {
	if len(suite.testFiles) == 0 {
		suite.T().Skip("No test files available")
	}

	fileID := suite.testFiles[0].ID

	// Create backup job
	backupJob := BackupJob{
		FileID:            fileID,
		BackupType:        "cross_region",
		SourceRegion:      "ap-southeast-1",
		DestinationRegion: "ap-southeast-2",
		RetentionPolicy:   "30d",
	}

	jsonData, _ := json.Marshal(backupJob)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/storage/backup", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var jobResponse BackupJob
	json.NewDecoder(resp.Body).Decode(&jobResponse)
	suite.NotEmpty(jobResponse.ID)

	// Check backup status
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/storage/backup/%s", suite.baseURL, jobResponse.ID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var backupStatus BackupJob
	json.NewDecoder(resp.Body).Decode(&backupStatus)

	suite.Equal(jobResponse.ID, backupStatus.ID)
	suite.Contains([]string{"queued", "running", "completed", "failed"}, backupStatus.Status)

	// Test backup verification
	verifyRequest := map[string]interface{}{
		"verify_integrity": true,
		"verify_accessibility": true,
	}

	jsonData, _ = json.Marshal(verifyRequest)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/storage/backup/%s/verify", suite.baseURL, jobResponse.ID),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusAccepted, resp.StatusCode)

	// List all backups for file
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/storage/files/%s/backups", suite.baseURL, fileID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var backupsResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&backupsResponse)

	suite.Contains(backupsResponse, "backups")
}

// Test storage analytics and monitoring
func (suite *Journey10FileStorageAPISuite) TestStorageAnalyticsAndMonitoring() {
	// Get storage usage statistics
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/storage/analytics/usage", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var usageResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&usageResponse)

	suite.Contains(usageResponse, "total_storage_used")
	suite.Contains(usageResponse, "files_count")
	suite.Contains(usageResponse, "storage_by_type")

	// Get user storage quota
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/storage/quota", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var quotaResponse StorageQuota
	json.NewDecoder(resp.Body).Decode(&quotaResponse)

	suite.Equal(suite.user1.UserID, quotaResponse.UserID)
	suite.Greater(quotaResponse.TotalQuota, int64(0))

	// Get file access analytics
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/storage/analytics/access?period=7d", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var accessResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&accessResponse)

	suite.Contains(accessResponse, "total_requests")
	suite.Contains(accessResponse, "bandwidth_used")
	suite.Contains(accessResponse, "regional_breakdown")

	// Get storage provider health
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/storage/providers/health", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var healthResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&healthResponse)

	suite.Contains(healthResponse, "providers")
	providers := healthResponse["providers"].([]interface{})
	suite.Greater(len(providers), 0)

	for _, provider := range providers {
		p := provider.(map[string]interface{})
		suite.Contains(p, "status")
		suite.Contains(p, "response_time")
		suite.Contains(p, "capacity_used")
	}
}

// Test file lifecycle management
func (suite *Journey10FileStorageAPISuite) TestFileLifecycleManagement() {
	if len(suite.testFiles) == 0 {
		suite.T().Skip("No test files available")
	}

	fileID := suite.testFiles[0].ID

	// Set file expiration
	expirationData := map[string]interface{}{
		"expires_at": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
		"auto_delete": true,
	}

	jsonData, _ := json.Marshal(expirationData)
	req, _ := http.NewRequest("PUT",
		fmt.Sprintf("%s/api/v1/files/%s/expiration", suite.baseURL, fileID),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	// Test file archiving
	archiveData := map[string]interface{}{
		"archive_tier": "cold_storage",
		"reason":       "Infrequent access",
	}

	jsonData, _ = json.Marshal(archiveData)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/files/%s/archive", suite.baseURL, fileID),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusAccepted, resp.StatusCode)

	// Test file restoration from archive
	time.Sleep(2 * time.Second)

	restoreData := map[string]interface{}{
		"priority":     "standard",
		"restore_days": 7,
	}

	jsonData, _ = json.Marshal(restoreData)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/files/%s/restore", suite.baseURL, fileID),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusAccepted, resp.StatusCode)

	// Test bulk operations
	bulkData := map[string]interface{}{
		"operation": "delete",
		"file_ids":  []string{fileID},
		"reason":    "Cleanup test files",
	}

	jsonData, _ = json.Marshal(bulkData)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/files/bulk", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusAccepted, resp.StatusCode)

	var bulkResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&bulkResponse)
	suite.Contains(bulkResponse, "job_id")
}

// Test file search and metadata management
func (suite *Journey10FileStorageAPISuite) TestFileSearchAndMetadataManagement() {
	// Search files by various criteria
	searchQueries := []map[string]string{
		{"query": "test", "type": "filename"},
		{"mime_type": "image/*", "privacy": "public"},
		{"owner_id": suite.user1.UserID, "limit": "10"},
		{"tags": "test,image", "sort": "created_at", "order": "desc"},
	}

	for _, query := range searchQueries {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/files/search", suite.baseURL), nil)
		req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)

		q := req.URL.Query()
		for key, value := range query {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()

		resp, err := suite.httpClient.Do(req)
		suite.NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusOK, resp.StatusCode)

		var searchResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&searchResponse)

		suite.Contains(searchResponse, "files")
		suite.Contains(searchResponse, "total_count")
	}

	// Update file metadata
	if len(suite.testFiles) > 0 {
		fileID := suite.testFiles[0].ID

		metadataUpdate := map[string]interface{}{
			"tags":        []string{"updated", "test", "metadata"},
			"description": "Updated description for testing",
			"custom_metadata": map[string]interface{}{
				"category":    "test_images",
				"source":      "automated_test",
				"version":     "1.1",
			},
		}

		jsonData, _ := json.Marshal(metadataUpdate)
		req, _ := http.NewRequest("PATCH",
			fmt.Sprintf("%s/api/v1/files/%s/metadata", suite.baseURL, fileID),
			bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.httpClient.Do(req)
		suite.NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusOK, resp.StatusCode)

		// Verify metadata update
		req, _ = http.NewRequest("GET",
			fmt.Sprintf("%s/api/v1/files/%s", suite.baseURL, fileID), nil)
		req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)

		resp, err = suite.httpClient.Do(req)
		suite.NoError(err)
		defer resp.Body.Close()

		var fileResponse FileMetadata
		json.NewDecoder(resp.Body).Decode(&fileResponse)

		suite.Contains(fileResponse.Tags, "updated")
		suite.Equal("Updated description for testing", fileResponse.Metadata["description"])
	}
}

func (suite *Journey10FileStorageAPISuite) cleanupTestFiles() {
	// Remove temporary files
	for _, tempFile := range suite.tempFiles {
		os.Remove(tempFile)
	}
}

func (suite *Journey10FileStorageAPISuite) cleanupTestData() {
	// Clean up test files and storage data
	req, _ := http.NewRequest("DELETE",
		fmt.Sprintf("%s/api/v1/storage/cleanup/test-data", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	_, _ = suite.httpClient.Do(req)
}

func TestJourney10FileStorageAPISuite(t *testing.T) {
	suite.Run(t, new(Journey10FileStorageAPISuite))
}