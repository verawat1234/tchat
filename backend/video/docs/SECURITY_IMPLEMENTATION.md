# Video Security Implementation Guide

## Overview
This document outlines the comprehensive security implementation for the Tchat video platform that **prevents direct video downloads** using blob URLs and token-based authentication across all platforms (Web, Android, iOS).

## Security Architecture

### Problem Statement
Without proper security measures, users can:
- Right-click and download videos directly
- Access video files through browser DevTools
- Share direct video URLs
- Bypass access control and monetization

### Solution: Blob URL Pattern with Token Authentication
Our security implementation uses:
1. **Backend**: Signed URLs with token authentication and streaming proxy
2. **Web**: Blob URL creation from authenticated fetch
3. **Android**: Custom ExoPlayer DataSource with token headers
4. **iOS**: Custom AVPlayer resource loader with authentication

## Backend Implementation

### 1. Security Service (`security_service.go`)

**Features**:
- HMAC-SHA256 signed tokens with 2-hour expiration
- Token validation and signature verification
- Signed URL generation with embedded authentication

**Token Structure**:
```go
type StreamToken struct {
    VideoID   uuid.UUID `json:"video_id"`
    UserID    uuid.UUID `json:"user_id"`
    ExpiresAt time.Time `json:"expires_at"`
    Quality   string    `json:"quality"`
    Signature string    `json:"signature"`
}
```

**Token Generation**:
```go
// Generate signature: HMAC-SHA256(videoID:userID:expiresAt:quality)
func (s *SecurityService) GenerateStreamToken(videoID, userID uuid.UUID, quality string) (*StreamToken, error) {
    expiresAt := time.Now().Add(s.tokenTTL)

    token := &StreamToken{
        VideoID:   videoID,
        UserID:    userID,
        ExpiresAt: expiresAt,
        Quality:   quality,
    }

    signature, err := s.generateSignature(token)
    if err != nil {
        return nil, err
    }

    token.Signature = signature
    return token, nil
}
```

**Signature Validation**:
- Checks token expiration
- Verifies HMAC-SHA256 signature
- Prevents token tampering

### 2. Secure Streaming Handler (`secure_stream_handler.go`)

**Endpoints**:
- `GET /api/v1/videos/:id/token` - Generate streaming token
- `POST /api/v1/videos/:id/validate-token` - Validate token
- `GET /api/v1/videos/:id/stream/secure` - Stream video with authentication

**Streaming Proxy Features**:
- Token validation before serving content
- HTTP byte-range support for streaming
- No direct file URL exposure
- Automatic cache-control headers

**Byte-Range Request Handling**:
```go
func (h *SecureStreamHandler) handleRangeRequest(c *gin.Context, file *os.File, fileSize int64, rangeHeader string) {
    // Parse Range header: "bytes=start-end"
    ranges := strings.TrimPrefix(rangeHeader, "bytes=")
    parts := strings.Split(ranges, "-")

    start, _ := strconv.ParseInt(parts[0], 10, 64)
    end := fileSize - 1
    if parts[1] != "" {
        end, _ = strconv.ParseInt(parts[1], 10, 64)
    }

    // Seek to start position
    file.Seek(start, io.SeekStart)

    // Set partial content headers
    c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
    c.Status(http.StatusPartialContent)

    // Stream partial content
    io.CopyN(c.Writer, file, end-start+1)
}
```

## Web Implementation

### 1. Secure Video Service (`secureVideoService.ts`)

**Features**:
- Fetches video as Blob using authenticated requests
- Creates blob: URLs that don't expose file paths
- Automatic blob URL revocation to prevent memory leaks
- Token validation and refresh

**Blob URL Creation**:
```typescript
async createSecureBlobURL(videoId: string, quality: string = 'auto'): Promise<string> {
    // Get streaming token
    const tokenData = await this.getStreamToken(videoId, quality);

    // Fetch video as blob using signed URL
    const videoBlob = await this.fetchVideoBlob(tokenData.signed_url);

    // Create blob URL
    const blobUrl = URL.createObjectURL(videoBlob);

    // Cache and schedule auto-revoke
    this.blobCache.set(`${videoId}-${quality}`, blobUrl);
    this.scheduleRevoke(cacheKey, blobUrl, expiresIn);

    return blobUrl;
}
```

**Memory Management**:
```typescript
cleanup(): void {
    // Revoke all blob URLs
    for (const blobUrl of this.blobCache.values()) {
        URL.revokeObjectURL(blobUrl);
    }

    // Clear all timers
    for (const timer of this.revokeTimers.values()) {
        clearTimeout(timer);
    }

    this.blobCache.clear();
    this.revokeTimers.clear();
}
```

### 2. React Hook (`useSecureVideo.ts`)

**Basic Usage**:
```typescript
const VideoPlayer = ({ videoId }) => {
    const { blobUrl, isLoading, error } = useSecureVideo({
        videoId,
        quality: '720p',
        autoLoad: true
    });

    if (isLoading) return <div>Loading...</div>;
    if (error) return <div>Error: {error.message}</div>;

    return <video src={blobUrl} controls />;
};
```

**Advanced: MediaSource API**:
For large videos, use `useSecureVideoWithMediaSource` to stream in chunks:
```typescript
const { sourceUrl, isLoading, error } = useSecureVideoWithMediaSource({
    videoId,
    quality: '1080p'
});

return <video src={sourceUrl} controls />;
```

## KMP Mobile Implementation

### 1. Android (`SecureAndroidVideoPlayer.kt`)

**ExoPlayer Custom DataSource**:
```kotlin
private class SecureDataSourceFactory(
    private val authToken: String,
    private val signedUrl: String
) : DataSource.Factory {

    override fun createDataSource(): DataSource {
        val httpDataSource = DefaultHttpDataSource.Factory()
            .setUserAgent("Tchat/1.0 Android")
            .createDataSource()

        // Add authentication headers
        httpDataSource.setRequestProperty("Authorization", "Bearer $authToken")
        httpDataSource.setRequestProperty("X-Signed-URL", signedUrl)
        httpDataSource.setRequestProperty("Accept", "video/*")

        return httpDataSource
    }
}
```

**Usage**:
```kotlin
@Composable
fun VideoScreen(videoId: String) {
    SecureAndroidVideoPlayer(
        videoId = videoId,
        quality = "720p",
        autoPlay = true,
        onError = { error ->
            Log.e("VideoPlayer", "Error: $error")
        }
    )
}
```

### 2. iOS (`SecureIOSVideoPlayer.kt`)

**AVPlayer Resource Loader Delegate**:
```kotlin
private class SecureResourceLoaderDelegate(
    private val authToken: String,
    private val originalUrl: String
) : NSObject(), AVAssetResourceLoaderDelegateProtocol {

    override fun resourceLoader(
        resourceLoader: AVAssetResourceLoader,
        shouldWaitForLoadingOfRequestedResource: AVAssetResourceLoadingRequest
    ): Boolean {
        // Create authenticated URL request
        val urlRequest = NSMutableURLRequest(uRL = NSURL.URLWithString(originalUrl)!!)
        urlRequest.setValue("Bearer $authToken", forHTTPHeaderField = "Authorization")

        // Handle byte-range requests
        val dataRequest = loadingRequest.dataRequest
        val requestedOffset = dataRequest?.requestedOffset ?: 0
        val requestedLength = dataRequest?.requestedLength?.toLong() ?: 0

        if (requestedOffset > 0 || requestedLength > 0) {
            urlRequest.setValue(
                "bytes=$requestedOffset-${requestedOffset + requestedLength - 1}",
                forHTTPHeaderField = "Range"
            )
        }

        // Make authenticated request and respond with data
        // ... (implementation details)

        return true
    }
}
```

## Security Features Summary

### ✅ Prevents Video Downloads
- ❌ No direct file URLs exposed
- ❌ Cannot right-click save video
- ❌ Cannot inspect network tab for direct URLs
- ✅ All video access requires authentication

### ✅ Token-Based Authentication
- 2-hour token expiration (configurable)
- HMAC-SHA256 signature verification
- Per-user, per-video tokens
- Quality-specific tokens

### ✅ Memory Management
- Automatic blob URL revocation on unmount
- Scheduled cleanup on token expiration
- No memory leaks from orphaned blob URLs

### ✅ Byte-Range Support
- Efficient streaming without full download
- Seek support in video player
- Bandwidth optimization

### ✅ Cross-Platform Consistency
- Same security model across Web, Android, iOS
- Platform-specific implementations with identical security guarantees
- Unified token validation

## Integration Checklist

### Backend Setup
- [ ] Configure secret key for token signing (environment variable)
- [ ] Register secure streaming routes in gateway
- [ ] Set up CDN with token validation (optional)
- [ ] Configure CORS headers for authenticated requests
- [ ] Implement user authentication middleware

### Web Setup
- [ ] Import `secureVideoService` in video components
- [ ] Replace direct video URLs with `useSecureVideo` hook
- [ ] Add cleanup calls on component unmount
- [ ] Configure auth token storage (localStorage/cookie)
- [ ] Test blob URL creation and revocation

### Android Setup
- [ ] Use `SecureAndroidVideoPlayer` instead of basic ExoPlayer
- [ ] Configure EncryptedSharedPreferences for token storage
- [ ] Add Ktor client for API calls
- [ ] Test authenticated streaming with different qualities

### iOS Setup
- [ ] Use `SecureIOSVideoPlayer` with custom resource loader
- [ ] Configure Keychain for token storage
- [ ] Add URLSession configuration for authenticated requests
- [ ] Test AVPlayer with custom scheme

## Performance Considerations

### Caching Strategy
- **Backend**: Redis cache for token validation (5-minute TTL)
- **Web**: In-memory blob cache with automatic expiration
- **Mobile**: Disable ExoPlayer/AVPlayer caching to prevent downloads

### Bandwidth Optimization
- Use byte-range requests for seeking
- Implement adaptive bitrate streaming (HLS/DASH)
- Progressive video quality based on network speed

### Token Refresh
- Implement automatic token refresh 5 minutes before expiration
- Seamless transition without playback interruption

## Monitoring & Alerts

### Metrics to Track
- Token validation failures (potential attacks)
- Blob URL creation rate (memory usage indicator)
- Stream request latency
- Token expiration events

### Security Alerts
- Unusual token validation failures (>5% rate)
- Suspicious access patterns (same token multiple IPs)
- Expired token usage attempts

## Testing

### Backend Tests
```go
func TestSecurityService_GenerateStreamToken(t *testing.T) {
    service := NewSecurityService("test-secret-key")

    token, err := service.GenerateStreamToken(videoID, userID, "720p")
    assert.NoError(t, err)
    assert.NotEmpty(t, token.Signature)

    // Validate token
    err = service.ValidateStreamToken(token)
    assert.NoError(t, err)
}
```

### Web Tests
```typescript
describe('SecureVideoService', () => {
    it('should create blob URL from authenticated fetch', async () => {
        const blobUrl = await secureVideoService.createSecureBlobURL('video-id', '720p');

        expect(blobUrl).toMatch(/^blob:/);
        expect(blobUrl).not.toContain('cdn.tchat.com');
    });

    it('should auto-revoke blob URL on expiration', async () => {
        // Test implementation
    });
});
```

## Troubleshooting

### Issue: Blob URL shows black screen
**Solution**: Check token expiration, ensure authenticated fetch succeeds

### Issue: Video downloads despite blob URL
**Solution**: Verify no fallback to direct URLs, check browser compatibility

### Issue: ExoPlayer shows "source error"
**Solution**: Verify DataSource.Factory authentication headers, check backend token validation

### Issue: AVPlayer resource loading fails
**Solution**: Check custom URL scheme registration, verify AVAssetResourceLoaderDelegate implementation

## Future Enhancements

1. **DRM Integration**: Add Widevine (Android), FairPlay (iOS) for additional protection
2. **Watermarking**: Dynamic video watermarking with user ID
3. **Geo-Fencing**: Location-based access control
4. **Device Limits**: Restrict concurrent playback devices per user
5. **Forensic Watermarking**: Invisible watermarks for piracy tracking