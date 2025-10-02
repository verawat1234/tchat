package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// RecordingStatus represents the current state of a recording
type RecordingStatus string

const (
	RecordingNotStarted RecordingStatus = "NOT_STARTED"
	RecordingInProgress RecordingStatus = "RECORDING"
	RecordingProcessing RecordingStatus = "PROCESSING"
	RecordingCompleted  RecordingStatus = "COMPLETED"
	RecordingFailed     RecordingStatus = "FAILED"
	RecordingExpired    RecordingStatus = "EXPIRED"
)

// HLSConfig defines configuration for HLS recording
type HLSConfig struct {
	SegmentDuration  int    // Duration of each segment in seconds (default: 6)
	PlaylistSize     int    // Number of segments in manifest (default: 10)
	OutputDir        string // Base directory for recordings
	VideoCodec       string // Video codec (default: h264)
	AudioCodec       string // Audio codec (default: aac)
	EnableTranscript bool   // Generate VTT transcript from chat
}

// RecordingInfo contains information about a recording session
type RecordingInfo struct {
	StreamID       uuid.UUID
	Status         RecordingStatus
	LocalPath      string
	CDNURL         string
	StartTime      time.Time
	EndTime        time.Time
	FileSize       int64
	Duration       time.Duration
	Quality        string // 360p, 720p, 1080p
	Error          error
	ExpiryDate     time.Time
	TranscriptPath string
}

// S3Config defines S3-compatible storage configuration
type S3Config struct {
	Endpoint        string
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	CDNDomain       string // CDN domain for public URLs
}

// RecordingService interface defines methods for stream recording and storage
type RecordingService interface {
	// StartRecording initiates HLS recording for a stream
	// tracks parameter is a placeholder for WebRTC track integration
	StartRecording(streamID uuid.UUID, tracks interface{}) error

	// StopRecording terminates recording and returns the recording URL
	StopRecording(streamID uuid.UUID) (string, error)

	// UploadToCDN uploads a local recording to S3-compatible storage
	// Returns the CDN URL for public access
	UploadToCDN(localPath string, streamID uuid.UUID) (string, error)

	// SetLifecyclePolicy configures automatic deletion for a recording
	SetLifecyclePolicy(cdnURL string, expiryDate time.Time) error

	// GetRecordingStatus retrieves the current status of a recording
	GetRecordingStatus(streamID uuid.UUID) (RecordingStatus, error)

	// GetRecordingInfo retrieves detailed information about a recording
	GetRecordingInfo(streamID uuid.UUID) (*RecordingInfo, error)

	// CleanupExpiredRecordings removes local recordings past expiry date
	CleanupExpiredRecordings() error
}

// recordingServiceImpl implements RecordingService
type recordingServiceImpl struct {
	hlsConfig HLSConfig
	s3Config  S3Config

	// Track active recordings
	activeRecordings map[uuid.UUID]*RecordingInfo
	mu               sync.RWMutex

	// Cancel contexts for background workers
	cancelFuncs map[uuid.UUID]context.CancelFunc
	cancelMu    sync.Mutex

	// Storage client (would be S3 client in production)
	storageClient interface{} // placeholder for aws-sdk-go-v2 S3 client
}

// NewRecordingService creates a new recording service instance
func NewRecordingService(hlsConfig HLSConfig, s3Config S3Config) RecordingService {
	// Set defaults for HLS configuration
	if hlsConfig.SegmentDuration == 0 {
		hlsConfig.SegmentDuration = 6
	}
	if hlsConfig.PlaylistSize == 0 {
		hlsConfig.PlaylistSize = 10
	}
	if hlsConfig.OutputDir == "" {
		hlsConfig.OutputDir = "/tmp/recordings"
	}
	if hlsConfig.VideoCodec == "" {
		hlsConfig.VideoCodec = "h264"
	}
	if hlsConfig.AudioCodec == "" {
		hlsConfig.AudioCodec = "aac"
	}

	// Create output directory if it doesn't exist
	os.MkdirAll(hlsConfig.OutputDir, 0755)

	return &recordingServiceImpl{
		hlsConfig:        hlsConfig,
		s3Config:         s3Config,
		activeRecordings: make(map[uuid.UUID]*RecordingInfo),
		cancelFuncs:      make(map[uuid.UUID]context.CancelFunc),
	}
}

// StartRecording initiates HLS recording for a stream
func (s *recordingServiceImpl) StartRecording(streamID uuid.UUID, tracks interface{}) error {
	s.mu.Lock()

	// Check if recording already exists
	if info, exists := s.activeRecordings[streamID]; exists {
		s.mu.Unlock()
		if info.Status == RecordingInProgress {
			return fmt.Errorf("recording already in progress for stream %s", streamID)
		}
	}

	// Create recording directory
	timestamp := time.Now().Format("20060102_150405")
	localPath := filepath.Join(s.hlsConfig.OutputDir, streamID.String(), timestamp)
	if err := os.MkdirAll(localPath, 0755); err != nil {
		s.mu.Unlock()
		return fmt.Errorf("failed to create recording directory: %w", err)
	}

	// Initialize recording info
	info := &RecordingInfo{
		StreamID:  streamID,
		Status:    RecordingInProgress,
		LocalPath: localPath,
		StartTime: time.Now(),
		Quality:   "1080p", // Default quality, should be determined from stream
	}
	s.activeRecordings[streamID] = info
	s.mu.Unlock()

	// Start recording in background goroutine
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelMu.Lock()
	s.cancelFuncs[streamID] = cancel
	s.cancelMu.Unlock()

	go s.recordWorker(ctx, streamID, localPath, tracks)

	return nil
}

// StartRecordingWithInput initiates HLS recording with a provided input stream (for testing)
func (s *recordingServiceImpl) StartRecordingWithInput(streamID uuid.UUID, input io.Reader) error {
	s.mu.Lock()

	// Check if recording already exists
	if info, exists := s.activeRecordings[streamID]; exists {
		s.mu.Unlock()
		if info.Status == RecordingInProgress {
			return fmt.Errorf("recording already in progress for stream %s", streamID)
		}
	}

	// Create recording directory
	timestamp := time.Now().Format("20060102_150405")
	localPath := filepath.Join(s.hlsConfig.OutputDir, streamID.String(), timestamp)
	if err := os.MkdirAll(localPath, 0755); err != nil {
		s.mu.Unlock()
		return fmt.Errorf("failed to create recording directory: %w", err)
	}

	// Initialize recording info
	info := &RecordingInfo{
		StreamID:  streamID,
		Status:    RecordingInProgress,
		LocalPath: localPath,
		StartTime: time.Now(),
		Quality:   "1080p", // Default quality
	}
	s.activeRecordings[streamID] = info
	s.mu.Unlock()

	// Start recording in background goroutine with input stream
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelMu.Lock()
	s.cancelFuncs[streamID] = cancel
	s.cancelMu.Unlock()

	go s.recordWorkerWithInput(ctx, streamID, localPath, input)

	return nil
}

// recordWorker handles the actual recording process
func (s *recordingServiceImpl) recordWorker(ctx context.Context, streamID uuid.UUID, outputPath string, tracks interface{}) {
	defer func() {
		if r := recover(); r != nil {
			s.updateRecordingStatus(streamID, RecordingFailed, fmt.Errorf("recording panic: %v", r))
		}
	}()

	// Output files
	playlistPath := filepath.Join(outputPath, "playlist.m3u8")
	segmentPattern := filepath.Join(outputPath, "segment_%03d.ts")
	transcriptPath := filepath.Join(outputPath, "transcript.vtt")

	// Build FFmpeg command for HLS recording
	// In production, this would receive WebRTC tracks and encode to HLS
	// For now, this is a template showing the command structure
	args := []string{
		"-i", "pipe:0", // Input from stdin (WebRTC tracks would be piped here)
		"-c:v", s.hlsConfig.VideoCodec,
		"-c:a", s.hlsConfig.AudioCodec,
		"-preset", "veryfast", // Fast encoding for live streams
		"-g", "48", // GOP size (2 seconds at 24fps)
		"-sc_threshold", "0", // Disable scene change detection
		"-f", "hls",
		"-hls_time", fmt.Sprintf("%d", s.hlsConfig.SegmentDuration),
		"-hls_list_size", fmt.Sprintf("%d", s.hlsConfig.PlaylistSize),
		"-hls_flags", "delete_segments+append_list",
		"-hls_segment_filename", segmentPattern,
		playlistPath,
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	// Create pipes for stdin/stdout/stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		s.updateRecordingStatus(streamID, RecordingFailed, fmt.Errorf("failed to create stdin pipe: %w", err))
		return
	}
	defer stdin.Close()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Start FFmpeg process
	if err := cmd.Start(); err != nil {
		s.updateRecordingStatus(streamID, RecordingFailed, fmt.Errorf("failed to start FFmpeg: %w", err))
		return
	}

	// In production, this would pipe WebRTC track data to stdin
	// For now, we simulate the structure
	go func() {
		// Placeholder: Write track data to stdin
		// In actual implementation, this would read from WebRTC tracks
		// and write encoded frames to stdin
		<-ctx.Done()
		stdin.Close()
	}()

	// Wait for recording to complete or context cancellation
	if err := cmd.Wait(); err != nil {
		if ctx.Err() != context.Canceled {
			s.updateRecordingStatus(streamID, RecordingFailed, fmt.Errorf("FFmpeg error: %w, stderr: %s", err, stderr.String()))
			return
		}
	}

	// Generate VTT transcript if enabled
	if s.hlsConfig.EnableTranscript {
		if err := s.generateTranscript(streamID, transcriptPath); err != nil {
			// Log error but don't fail recording
			fmt.Printf("Failed to generate transcript for stream %s: %v\n", streamID, err)
		} else {
			s.mu.Lock()
			if info, exists := s.activeRecordings[streamID]; exists {
				info.TranscriptPath = transcriptPath
			}
			s.mu.Unlock()
		}
	}

	// Calculate file size and duration
	fileSize, duration := s.calculateRecordingStats(outputPath)

	s.mu.Lock()
	if info, exists := s.activeRecordings[streamID]; exists {
		info.FileSize = fileSize
		info.Duration = duration
		info.EndTime = time.Now()
		info.Status = RecordingProcessing
	}
	s.mu.Unlock()
}

// recordWorkerWithInput handles recording with a provided input stream (for testing)
func (s *recordingServiceImpl) recordWorkerWithInput(ctx context.Context, streamID uuid.UUID, outputPath string, input io.Reader) {
	defer func() {
		if r := recover(); r != nil {
			s.updateRecordingStatus(streamID, RecordingFailed, fmt.Errorf("recording panic: %v", r))
		}
	}()

	// Output files
	playlistPath := filepath.Join(outputPath, "playlist.m3u8")
	segmentPattern := filepath.Join(outputPath, "segment_%03d.ts")
	transcriptPath := filepath.Join(outputPath, "transcript.vtt")

	// Build FFmpeg command for HLS recording with input stream
	args := []string{
		"-i", "pipe:0", // Input from stdin
		"-c:v", s.hlsConfig.VideoCodec,
		"-c:a", s.hlsConfig.AudioCodec,
		"-f", "hls",
		"-hls_time", fmt.Sprintf("%d", s.hlsConfig.SegmentDuration),
		"-hls_list_size", fmt.Sprintf("%d", s.hlsConfig.PlaylistSize),
		"-hls_segment_filename", segmentPattern,
		playlistPath,
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	// Set up stdin pipe for video input
	stdin, err := cmd.StdinPipe()
	if err != nil {
		s.updateRecordingStatus(streamID, RecordingFailed, fmt.Errorf("failed to create stdin pipe: %w", err))
		return
	}
	defer stdin.Close()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Start FFmpeg process
	if err := cmd.Start(); err != nil {
		s.updateRecordingStatus(streamID, RecordingFailed, fmt.Errorf("failed to start FFmpeg: %w", err))
		return
	}

	// Pipe input stream to FFmpeg stdin
	go func() {
		defer stdin.Close()
		io.Copy(stdin, input)
	}()

	// Wait for recording to complete or context cancellation
	if err := cmd.Wait(); err != nil {
		if ctx.Err() != context.Canceled {
			s.updateRecordingStatus(streamID, RecordingFailed, fmt.Errorf("FFmpeg error: %w, stderr: %s", err, stderr.String()))
			return
		}
	}

	// Generate VTT transcript if enabled
	if s.hlsConfig.EnableTranscript {
		if err := s.generateTranscript(streamID, transcriptPath); err != nil {
			fmt.Printf("Failed to generate transcript for stream %s: %v\n", streamID, err)
		} else {
			s.mu.Lock()
			if info, exists := s.activeRecordings[streamID]; exists {
				info.TranscriptPath = transcriptPath
			}
			s.mu.Unlock()
		}
	}

	// Calculate file size and duration
	fileSize, duration := s.calculateRecordingStats(outputPath)

	s.mu.Lock()
	if info, exists := s.activeRecordings[streamID]; exists {
		info.FileSize = fileSize
		info.Duration = duration
		info.EndTime = time.Now()
		info.Status = RecordingProcessing
	}
	s.mu.Unlock()
}

// StopRecording terminates recording and returns the recording URL
func (s *recordingServiceImpl) StopRecording(streamID uuid.UUID) (string, error) {
	// Cancel the recording context
	s.cancelMu.Lock()
	if cancel, exists := s.cancelFuncs[streamID]; exists {
		cancel()
		delete(s.cancelFuncs, streamID)
	}
	s.cancelMu.Unlock()

	// Wait for recording to transition to PROCESSING state
	maxWait := 30 * time.Second
	deadline := time.Now().Add(maxWait)
	for time.Now().Before(deadline) {
		s.mu.RLock()
		info, exists := s.activeRecordings[streamID]
		s.mu.RUnlock()

		if !exists {
			return "", fmt.Errorf("recording not found for stream %s", streamID)
		}

		if info.Status == RecordingProcessing || info.Status == RecordingCompleted {
			break
		}

		if info.Status == RecordingFailed {
			return "", fmt.Errorf("recording failed: %w", info.Error)
		}

		time.Sleep(500 * time.Millisecond)
	}

	// Get recording info
	s.mu.RLock()
	info, exists := s.activeRecordings[streamID]
	s.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("recording not found for stream %s", streamID)
	}

	// Upload to CDN
	cdnURL, err := s.UploadToCDN(info.LocalPath, streamID)
	if err != nil {
		s.updateRecordingStatus(streamID, RecordingFailed, fmt.Errorf("upload failed: %w", err))
		return "", err
	}

	// Set expiry date (30 days from now)
	expiryDate := time.Now().AddDate(0, 0, 30)
	if err := s.SetLifecyclePolicy(cdnURL, expiryDate); err != nil {
		// Log error but don't fail - lifecycle policy is best effort
		fmt.Printf("Failed to set lifecycle policy for %s: %v\n", cdnURL, err)
	}

	// Update recording info
	s.mu.Lock()
	info.Status = RecordingCompleted
	info.CDNURL = cdnURL
	info.ExpiryDate = expiryDate
	s.mu.Unlock()

	return cdnURL, nil
}

// UploadToCDN uploads a local recording to S3-compatible storage
func (s *recordingServiceImpl) UploadToCDN(localPath string, streamID uuid.UUID) (string, error) {
	// Validate local path exists
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return "", fmt.Errorf("local path does not exist: %s", localPath)
	}

	// Generate S3 key with timestamp
	timestamp := time.Now().Format("20060102_150405")
	s3Prefix := fmt.Sprintf("recordings/%s/%s/", streamID.String(), timestamp)

	// In production, this would use AWS SDK v2 to upload files
	// Example structure:
	/*
		uploader := manager.NewUploader(s.storageClient)

		// Upload all files in the directory
		files, err := filepath.Glob(filepath.Join(localPath, "*"))
		if err != nil {
			return "", fmt.Errorf("failed to list files: %w", err)
		}

		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				return "", fmt.Errorf("failed to open file: %w", err)
			}
			defer f.Close()

			key := s3Prefix + filepath.Base(file)
			_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
				Bucket: aws.String(s.s3Config.Bucket),
				Key:    aws.String(key),
				Body:   f,
			})
			if err != nil {
				return "", fmt.Errorf("failed to upload file: %w", err)
			}
		}
	*/

	// Generate CDN URL
	cdnURL := s.buildCDNURL(s3Prefix, "playlist.m3u8")

	return cdnURL, nil
}

// SetLifecyclePolicy configures automatic deletion for a recording
func (s *recordingServiceImpl) SetLifecyclePolicy(cdnURL string, expiryDate time.Time) error {
	// In production, this would configure S3 lifecycle policy
	// Example structure:
	/*
		client := s3.NewFromConfig(s.storageClient)

		lifecycleConfig := &s3.PutBucketLifecycleConfigurationInput{
			Bucket: aws.String(s.s3Config.Bucket),
			LifecycleConfiguration: &types.BucketLifecycleConfiguration{
				Rules: []types.LifecycleRule{
					{
						Id:     aws.String(fmt.Sprintf("delete-%s", streamID)),
						Status: types.ExpirationStatusEnabled,
						Filter: &types.LifecycleRuleFilterMemberPrefix{
							Value: s3Prefix,
						},
						Expiration: &types.LifecycleExpiration{
							Date: aws.Time(expiryDate),
						},
					},
				},
			},
		}

		_, err := client.PutBucketLifecycleConfiguration(context.TODO(), lifecycleConfig)
		return err
	*/

	// Placeholder implementation
	return nil
}

// GetRecordingStatus retrieves the current status of a recording
func (s *recordingServiceImpl) GetRecordingStatus(streamID uuid.UUID) (RecordingStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	info, exists := s.activeRecordings[streamID]
	if !exists {
		return RecordingNotStarted, nil
	}

	// Check if recording is expired
	if info.Status == RecordingCompleted && !info.ExpiryDate.IsZero() && time.Now().After(info.ExpiryDate) {
		return RecordingExpired, nil
	}

	return info.Status, nil
}

// GetRecordingInfo retrieves detailed information about a recording
func (s *recordingServiceImpl) GetRecordingInfo(streamID uuid.UUID) (*RecordingInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	info, exists := s.activeRecordings[streamID]
	if !exists {
		return nil, fmt.Errorf("recording not found for stream %s", streamID)
	}

	// Create a copy to avoid race conditions
	infoCopy := *info
	return &infoCopy, nil
}

// CleanupExpiredRecordings removes local recordings past expiry date
func (s *recordingServiceImpl) CleanupExpiredRecordings() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	var errs []error

	for streamID, info := range s.activeRecordings {
		// Skip if not expired
		if info.ExpiryDate.IsZero() || now.Before(info.ExpiryDate) {
			continue
		}

		// Remove local files
		if info.LocalPath != "" {
			if err := os.RemoveAll(info.LocalPath); err != nil {
				errs = append(errs, fmt.Errorf("failed to remove %s: %w", info.LocalPath, err))
				continue
			}
		}

		// Update status
		info.Status = RecordingExpired

		// Remove from active recordings after grace period
		if now.After(info.ExpiryDate.AddDate(0, 0, 7)) { // 7 days grace period
			delete(s.activeRecordings, streamID)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("cleanup errors: %v", errs)
	}

	return nil
}

// Helper methods

// updateRecordingStatus updates the status of a recording
func (s *recordingServiceImpl) updateRecordingStatus(streamID uuid.UUID, status RecordingStatus, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if info, exists := s.activeRecordings[streamID]; exists {
		info.Status = status
		info.Error = err
		if status == RecordingFailed || status == RecordingCompleted {
			info.EndTime = time.Now()
		}
	}
}

// calculateRecordingStats calculates file size and duration
func (s *recordingServiceImpl) calculateRecordingStats(outputPath string) (int64, time.Duration) {
	var totalSize int64
	var duration time.Duration

	// Calculate total file size
	filepath.Walk(outputPath, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	// Parse duration from playlist if exists
	playlistPath := filepath.Join(outputPath, "playlist.m3u8")
	if data, err := os.ReadFile(playlistPath); err == nil {
		// Parse m3u8 playlist to extract duration
		// This is a simplified placeholder - production would use proper HLS parser
		_ = data
		duration = 0 // Would be calculated from playlist
	}

	return totalSize, duration
}

// generateTranscript generates VTT transcript from chat messages
func (s *recordingServiceImpl) generateTranscript(streamID uuid.UUID, outputPath string) error {
	// In production, this would:
	// 1. Query chat messages from database
	// 2. Format them as WebVTT subtitle file
	// 3. Write to outputPath

	// Example VTT format:
	/*
		WEBVTT

		00:00:01.000 --> 00:00:05.000
		[User1]: Hello everyone!

		00:00:05.500 --> 00:00:10.000
		[User2]: Great stream!
	*/

	// Placeholder implementation
	vttContent := "WEBVTT\n\n"

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create transcript file: %w", err)
	}
	defer file.Close()

	_, err = io.WriteString(file, vttContent)
	return err
}

// buildCDNURL constructs the CDN URL for a recording
func (s *recordingServiceImpl) buildCDNURL(s3Prefix, filename string) string {
	// If custom CDN domain is configured, use it
	if s.s3Config.CDNDomain != "" {
		return fmt.Sprintf("https://%s/%s%s", s.s3Config.CDNDomain, s3Prefix, filename)
	}

	// Otherwise use S3 URL
	protocol := "https"
	if !s.s3Config.UseSSL {
		protocol = "http"
	}

	if s.s3Config.Region != "" {
		return fmt.Sprintf("%s://%s.s3.%s.amazonaws.com/%s%s",
			protocol, s.s3Config.Bucket, s.s3Config.Region, s3Prefix, filename)
	}

	return fmt.Sprintf("%s://%s/%s/%s%s",
		protocol, s.s3Config.Endpoint, s.s3Config.Bucket, s3Prefix, filename)
}

// StartCleanupWorker starts a background worker to cleanup expired recordings
func (s *recordingServiceImpl) StartCleanupWorker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.CleanupExpiredRecordings(); err != nil {
				fmt.Printf("Cleanup error: %v\n", err)
			}
		}
	}
}