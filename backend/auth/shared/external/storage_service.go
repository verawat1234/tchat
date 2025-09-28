package external

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"path/filepath"
	"strings"
	"time"
)

// StorageProvider represents different storage service providers
type StorageProvider string

const (
	// Cloud storage providers
	AWSS3Provider      StorageProvider = "aws_s3"
	GoogleCloudProvider StorageProvider = "google_cloud"
	AzureBlobProvider  StorageProvider = "azure_blob"

	// CDN providers
	CloudflareProvider StorageProvider = "cloudflare"
	CloudinaryProvider StorageProvider = "cloudinary"

	// Local/Regional providers
	AliyunOSSProvider  StorageProvider = "aliyun_oss"  // Popular in Southeast Asia
	QCloudProvider     StorageProvider = "qcloud"      // Tencent Cloud
	LocalProvider      StorageProvider = "local"
)

// StorageConfig holds storage service configuration
type StorageConfig struct {
	Provider        StorageProvider `mapstructure:"provider" validate:"required"`
	AccessKey       string          `mapstructure:"access_key" validate:"required"`
	SecretKey       string          `mapstructure:"secret_key" validate:"required"`
	Region          string          `mapstructure:"region" validate:"required"`
	Bucket          string          `mapstructure:"bucket" validate:"required"`
	Endpoint        string          `mapstructure:"endpoint"`
	BaseURL         string          `mapstructure:"base_url"`
	CDNDomain       string          `mapstructure:"cdn_domain"`
	UseSSL          bool            `mapstructure:"use_ssl"`
	PublicRead      bool            `mapstructure:"public_read"`
	Timeout         time.Duration   `mapstructure:"timeout"`
	MaxFileSize     int64           `mapstructure:"max_file_size"`     // bytes
	AllowedTypes    []string        `mapstructure:"allowed_types"`    // MIME types
	CompressImages  bool            `mapstructure:"compress_images"`
	GenerateThumbs  bool            `mapstructure:"generate_thumbs"`
	ThumbnailSizes  []ImageSize     `mapstructure:"thumbnail_sizes"`
}

// ImageSize represents thumbnail dimensions
type ImageSize struct {
	Name   string `mapstructure:"name"`
	Width  int    `mapstructure:"width"`
	Height int    `mapstructure:"height"`
	Quality int   `mapstructure:"quality"`
}

// DefaultStorageConfig returns default storage configuration
func DefaultStorageConfig() *StorageConfig {
	return &StorageConfig{
		Provider:       AWSS3Provider,
		Region:         "ap-southeast-1", // Singapore
		UseSSL:         true,
		PublicRead:     true,
		Timeout:        60 * time.Second,
		MaxFileSize:    50 * 1024 * 1024, // 50MB
		AllowedTypes:   []string{"image/jpeg", "image/png", "image/gif", "image/webp", "video/mp4", "audio/mpeg", "application/pdf"},
		CompressImages: true,
		GenerateThumbs: true,
		ThumbnailSizes: []ImageSize{
			{Name: "small", Width: 150, Height: 150, Quality: 80},
			{Name: "medium", Width: 300, Height: 300, Quality: 85},
			{Name: "large", Width: 800, Height: 600, Quality: 90},
		},
	}
}

// FileUpload represents a file upload request
type FileUpload struct {
	FileName    string            `json:"file_name" validate:"required"`
	ContentType string            `json:"content_type" validate:"required"`
	Size        int64             `json:"size" validate:"required,gt=0"`
	Data        io.Reader         `json:"-"`
	Folder      string            `json:"folder,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Public      bool              `json:"public,omitempty"`
	ExpiresAt   *time.Time        `json:"expires_at,omitempty"`
	UserID      string            `json:"user_id,omitempty"`
}

// FileInfo represents stored file information
type FileInfo struct {
	ID          string                 `json:"id"`
	FileName    string                 `json:"file_name"`
	OriginalName string                `json:"original_name"`
	ContentType string                 `json:"content_type"`
	Size        int64                  `json:"size"`
	URL         string                 `json:"url"`
	CDNUrl      string                 `json:"cdn_url,omitempty"`
	Folder      string                 `json:"folder,omitempty"`
	Provider    StorageProvider        `json:"provider"`
	Bucket      string                 `json:"bucket"`
	Key         string                 `json:"key"`
	ETag        string                 `json:"etag,omitempty"`
	Public      bool                   `json:"public"`
	Compressed  bool                   `json:"compressed,omitempty"`
	Thumbnails  map[string]string      `json:"thumbnails,omitempty"`
	Tags        map[string]string      `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
}

// StorageService interface for storage operations
type StorageService interface {
	Upload(ctx context.Context, upload *FileUpload) (*FileInfo, error)
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string, expires time.Duration) (string, error)
	GetInfo(ctx context.Context, key string) (*FileInfo, error)
	List(ctx context.Context, folder string, limit int, marker string) ([]*FileInfo, string, error)
	Copy(ctx context.Context, sourceKey, destKey string) error
	Move(ctx context.Context, sourceKey, destKey string) error
	HealthCheck(ctx context.Context) error
}

// StorageManager manages multiple storage providers with routing
type StorageManager struct {
	providers      map[StorageProvider]StorageService
	primaryConfig  *StorageConfig
	routingRules   map[string]StorageProvider // file type -> provider mapping
	fallbackOrder  []StorageProvider
	stats          *StorageStats
}

// StorageStats tracks storage statistics
type StorageStats struct {
	TotalFiles      int64                          `json:"total_files"`
	TotalSize       int64                          `json:"total_size"`
	TotalBandwidth  int64                          `json:"total_bandwidth"`
	ProviderStats   map[StorageProvider]*ProviderStorageStats `json:"provider_stats"`
	TypeStats       map[string]*TypeStats          `json:"type_stats"`
}

// ProviderStorageStats tracks per-provider statistics
type ProviderStorageStats struct {
	Files      int64   `json:"files"`
	Size       int64   `json:"size"`
	Bandwidth  int64   `json:"bandwidth"`
	Requests   int64   `json:"requests"`
	AvgLatency time.Duration `json:"avg_latency"`
	LastUsed   time.Time     `json:"last_used"`
}

// TypeStats tracks per-file-type statistics
type TypeStats struct {
	Files    int64     `json:"files"`
	Size     int64     `json:"size"`
	LastUsed time.Time `json:"last_used"`
}

// NewStorageManager creates a new storage manager
func NewStorageManager(config *StorageConfig) *StorageManager {
	manager := &StorageManager{
		providers:     make(map[StorageProvider]StorageService),
		primaryConfig: config,
		routingRules:  make(map[string]StorageProvider),
		fallbackOrder: []StorageProvider{AWSS3Provider, GoogleCloudProvider},
		stats: &StorageStats{
			ProviderStats: make(map[StorageProvider]*ProviderStorageStats),
			TypeStats:     make(map[string]*TypeStats),
		},
	}

	// Setup file type routing
	manager.setupTypeRouting()

	return manager
}

// setupTypeRouting configures file type specific routing
func (sm *StorageManager) setupTypeRouting() {
	// Images to CDN providers for better performance
	sm.routingRules["image/jpeg"] = CloudinaryProvider
	sm.routingRules["image/png"] = CloudinaryProvider
	sm.routingRules["image/gif"] = CloudinaryProvider
	sm.routingRules["image/webp"] = CloudinaryProvider

	// Videos to cloud storage with CDN
	sm.routingRules["video/mp4"] = AWSS3Provider
	sm.routingRules["video/webm"] = AWSS3Provider

	// Documents to standard cloud storage
	sm.routingRules["application/pdf"] = AWSS3Provider
	sm.routingRules["text/plain"] = AWSS3Provider

	// Audio files
	sm.routingRules["audio/mpeg"] = AWSS3Provider
	sm.routingRules["audio/wav"] = AWSS3Provider
}

// RegisterProvider registers a storage provider
func (sm *StorageManager) RegisterProvider(provider StorageProvider, service StorageService) {
	sm.providers[provider] = service
	sm.stats.ProviderStats[provider] = &ProviderStorageStats{}
	log.Printf("Registered storage provider: %s", provider)
}

// Upload uploads a file with automatic provider selection
func (sm *StorageManager) Upload(ctx context.Context, upload *FileUpload) (*FileInfo, error) {
	// Validate file
	if err := sm.validateFile(upload); err != nil {
		return nil, fmt.Errorf("file validation failed: %w", err)
	}

	// Generate unique file name
	upload.FileName = sm.generateFileName(upload.FileName, upload.UserID)

	// Determine best provider for this file
	provider := sm.selectProvider(upload)

	// Try primary provider
	if service, exists := sm.providers[provider]; exists {
		fileInfo, err := sm.uploadWithProvider(ctx, service, provider, upload)
		if err == nil {
			sm.updateStats(provider, upload.ContentType, upload.Size, true)
			return fileInfo, nil
		}
		log.Printf("Primary storage provider %s failed: %v", provider, err)
		sm.updateStats(provider, upload.ContentType, upload.Size, false)
	}

	// Try fallback providers
	for _, fallbackProvider := range sm.fallbackOrder {
		if fallbackProvider == provider {
			continue // Skip primary provider
		}

		if service, exists := sm.providers[fallbackProvider]; exists {
			fileInfo, err := sm.uploadWithProvider(ctx, service, fallbackProvider, upload)
			if err == nil {
				sm.updateStats(fallbackProvider, upload.ContentType, upload.Size, true)
				log.Printf("File uploaded successfully using fallback provider: %s", fallbackProvider)
				return fileInfo, nil
			}
			log.Printf("Fallback storage provider %s failed: %v", fallbackProvider, err)
			sm.updateStats(fallbackProvider, upload.ContentType, upload.Size, false)
		}
	}

	return nil, fmt.Errorf("all storage providers failed to upload file")
}

// uploadWithProvider uploads file using a specific provider
func (sm *StorageManager) uploadWithProvider(ctx context.Context, service StorageService, provider StorageProvider, upload *FileUpload) (*FileInfo, error) {
	startTime := time.Now()

	fileInfo, err := service.Upload(ctx, upload)
	if err != nil {
		return nil, err
	}

	fileInfo.Provider = provider
	fileInfo.CreatedAt = startTime

	// Generate CDN URLs if available
	if sm.primaryConfig.CDNDomain != "" {
		fileInfo.CDNUrl = sm.generateCDNURL(fileInfo.Key)
	}

	return fileInfo, nil
}

// selectProvider selects the best storage provider based on file type and other factors
func (sm *StorageManager) selectProvider(upload *FileUpload) StorageProvider {
	// Check file type specific routing
	if provider, exists := sm.routingRules[upload.ContentType]; exists {
		if _, providerExists := sm.providers[provider]; providerExists {
			return provider
		}
	}

	// Check file size routing (large files to specialized providers)
	if upload.Size > 100*1024*1024 { // > 100MB
		if _, exists := sm.providers[AWSS3Provider]; exists {
			return AWSS3Provider
		}
	}

	// Fall back to primary provider
	return sm.primaryConfig.Provider
}

// validateFile validates file upload requirements
func (sm *StorageManager) validateFile(upload *FileUpload) error {
	// Check file size
	if upload.Size > sm.primaryConfig.MaxFileSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", upload.Size, sm.primaryConfig.MaxFileSize)
	}

	// Check file type
	if len(sm.primaryConfig.AllowedTypes) > 0 {
		allowed := false
		for _, allowedType := range sm.primaryConfig.AllowedTypes {
			if upload.ContentType == allowedType {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("file type %s is not allowed", upload.ContentType)
		}
	}

	// Validate file name
	if upload.FileName == "" {
		return fmt.Errorf("file name is required")
	}

	return nil
}

// generateFileName generates a unique file name
func (sm *StorageManager) generateFileName(originalName, userID string) string {
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)

	// Clean the name
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ToLower(name)

	// Generate unique suffix
	timestamp := time.Now().UnixNano()

	if userID != "" {
		return fmt.Sprintf("%s/%s_%d%s", userID, name, timestamp, ext)
	}

	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}

// generateCDNURL generates CDN URL for a file
func (sm *StorageManager) generateCDNURL(key string) string {
	if sm.primaryConfig.CDNDomain == "" {
		return ""
	}

	return fmt.Sprintf("https://%s/%s", sm.primaryConfig.CDNDomain, key)
}

// Download downloads a file from storage
func (sm *StorageManager) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	// Try to download from all providers (since we don't know which one has it)
	for provider, service := range sm.providers {
		reader, err := service.Download(ctx, key)
		if err == nil {
			log.Printf("File downloaded from provider: %s", provider)
			return reader, nil
		}
	}

	return nil, fmt.Errorf("file not found in any storage provider: %s", key)
}

// Delete deletes a file from storage
func (sm *StorageManager) Delete(ctx context.Context, key string) error {
	var errors []string

	// Try to delete from all providers
	for provider, service := range sm.providers {
		if err := service.Delete(ctx, key); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", provider, err))
		} else {
			log.Printf("File deleted from provider: %s", provider)
		}
	}

	if len(errors) > 0 {
		log.Printf("Some delete operations failed: %s", strings.Join(errors, "; "))
	}

	return nil // Don't fail if file is deleted from at least one provider
}

// GetURL gets a signed URL for file access
func (sm *StorageManager) GetURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	// Try to get URL from all providers
	for provider, service := range sm.providers {
		url, err := service.GetURL(ctx, key, expires)
		if err == nil {
			log.Printf("URL generated from provider: %s", provider)
			return url, nil
		}
	}

	return "", fmt.Errorf("failed to generate URL for file: %s", key)
}

// GetInfo gets file information
func (sm *StorageManager) GetInfo(ctx context.Context, key string) (*FileInfo, error) {
	// Try to get info from all providers
	for provider, service := range sm.providers {
		info, err := service.GetInfo(ctx, key)
		if err == nil && info != nil {
			info.Provider = provider
			return info, nil
		}
	}

	return nil, fmt.Errorf("file info not found: %s", key)
}

// updateStats updates storage statistics
func (sm *StorageManager) updateStats(provider StorageProvider, contentType string, size int64, success bool) {
	// Update provider stats
	providerStats := sm.stats.ProviderStats[provider]
	if providerStats == nil {
		providerStats = &ProviderStorageStats{}
		sm.stats.ProviderStats[provider] = providerStats
	}

	providerStats.Requests++
	if success {
		providerStats.Files++
		providerStats.Size += size
	}
	providerStats.LastUsed = time.Now()

	// Update type stats
	typeStats := sm.stats.TypeStats[contentType]
	if typeStats == nil {
		typeStats = &TypeStats{}
		sm.stats.TypeStats[contentType] = typeStats
	}

	if success {
		typeStats.Files++
		typeStats.Size += size
	}
	typeStats.LastUsed = time.Now()

	// Update totals
	if success {
		sm.stats.TotalFiles++
		sm.stats.TotalSize += size
	}
}

// GetStats returns storage statistics
func (sm *StorageManager) GetStats() *StorageStats {
	return sm.stats
}

// HealthCheck checks the health of all storage providers
func (sm *StorageManager) HealthCheck(ctx context.Context) error {
	var errors []string

	for provider, service := range sm.providers {
		if err := service.HealthCheck(ctx); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", provider, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("storage provider health check failures: %s", strings.Join(errors, "; "))
	}

	return nil
}

// Helper functions

// GetMimeType gets MIME type from file extension
func GetMimeType(filename string) string {
	ext := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream"
	}
	return mimeType
}

// IsImageFile checks if file is an image
func IsImageFile(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}

// IsVideoFile checks if file is a video
func IsVideoFile(contentType string) bool {
	return strings.HasPrefix(contentType, "video/")
}

// IsAudioFile checks if file is audio
func IsAudioFile(contentType string) bool {
	return strings.HasPrefix(contentType, "audio/")
}

// GetFileCategory categorizes file by content type
func GetFileCategory(contentType string) string {
	switch {
	case IsImageFile(contentType):
		return "image"
	case IsVideoFile(contentType):
		return "video"
	case IsAudioFile(contentType):
		return "audio"
	case strings.HasPrefix(contentType, "text/"):
		return "document"
	case contentType == "application/pdf":
		return "document"
	default:
		return "other"
	}
}

// ValidateFileName validates file name format
func ValidateFileName(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	if len(filename) > 255 {
		return fmt.Errorf("filename too long")
	}

	// Check for invalid characters
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(filename, char) {
			return fmt.Errorf("filename contains invalid character: %s", char)
		}
	}

	return nil
}

// FormatFileSize formats file size in human readable format
func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}