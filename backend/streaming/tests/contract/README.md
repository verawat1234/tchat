# Pact Consumer Contract Tests for Streaming Service

## Overview

This directory contains Pact consumer contract tests for the POST /api/v1/streams endpoint. These tests define the contract between streaming clients (web/mobile) and the streaming service backend.

## Test Coverage

### `create_stream_test.go`

Comprehensive Pact consumer tests covering:

1. **CreateStoreStream_Success**: Successful store stream creation with KYC Tier 1
   - Validates that users with Standard KYC (Tier 1) can create store streams
   - Verifies response contains stream_key (RTMP URL) and webrtc_session_id
   - Validates max_capacity defaults to 50,000 viewers

2. **CreateStoreStream_InsufficientKYC**: KYC verification rejection
   - Validates that users without Standard KYC (Tier < 1) cannot create store streams
   - Verifies proper 403 Forbidden response
   - Checks error response includes required_kyc_tier, current_kyc_tier, and verification_url

3. **CreateVideoStream_NoKYCRequired**: Video stream creation without KYC
   - Validates that video streams can be created without KYC restrictions
   - Demonstrates business rule difference between store and video stream types
   - Verifies KYC Tier 0 users can create video streams

4. **CreateStream_MissingRequiredFields**: Request validation
   - Validates proper validation error handling
   - Verifies 400 Bad Request response for missing required fields
   - Checks detailed validation_errors object in response

## Prerequisites

### Install Pact FFI Library

The Pact Go v2 library requires the Pact FFI (Foreign Function Interface) library to be installed.

#### macOS (Homebrew)
```bash
brew tap pact-foundation/pact-ruby-standalone
brew install pact-ruby-standalone
```

Or download manually:
```bash
# Download the Pact FFI library for your platform
# https://github.com/pact-foundation/pact-reference/releases

# For macOS ARM64:
curl -LO https://github.com/pact-foundation/pact-reference/releases/download/libpact_ffi-v0.4.9/libpact_ffi-macos-aarch64-apple-darwin.tar.gz
tar -xzf libpact_ffi-macos-aarch64-apple-darwin.tar.gz
sudo cp lib/libpact_ffi.dylib /usr/local/lib/
sudo cp include/pact.h /usr/local/include/
```

#### Linux
```bash
# Download for Linux x86_64
curl -LO https://github.com/pact-foundation/pact-reference/releases/download/libpact_ffi-v0.4.9/libpact_ffi-linux-x86_64.tar.gz
tar -xzf libpact_ffi-linux-x86_64.tar.gz
sudo cp lib/libpact_ffi.so /usr/local/lib/
sudo ldconfig
```

## Running Tests

### Run Contract Tests
```bash
# From the streaming service root directory
go test -v ./tests/contract/create_stream_test.go -timeout 30s

# Run all contract tests in the directory
go test -v ./tests/contract/ -timeout 30s
```

### Expected Initial Behavior (TDD)

These tests are written following Test-Driven Development (TDD) principles and **MUST FAIL** initially because:

1. ❌ No handler implementation exists for POST /streams endpoint
2. ❌ No KYC validation logic is implemented
3. ❌ No stream_key generation service exists
4. ❌ No WebRTC session management is implemented
5. ❌ No database repository layer for live_streams table

This is the expected and correct behavior for TDD!

### Pact Files Output

Successful test execution generates Pact contract files in:
```
./pacts/streaming-web-client-streaming-service.json
```

These files can be:
- Shared with the provider (streaming service implementation team)
- Used for provider verification tests
- Published to a Pact Broker for contract management

## Contract Specifications

### Request Format

```json
{
  "stream_type": "store" | "video",
  "title": "Product Launch Event",
  "description": "Join us for an exclusive look at our new collection",
  "privacy_setting": "public" | "followers_only" | "private",
  "scheduled_start_time": "2025-09-30T15:00:00Z",
  "language": "en",
  "tags": ["commerce", "product-launch"]
}
```

### Success Response (201 Created)

```json
{
  "success": true,
  "message": "Stream created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "broadcaster_id": "550e8400-e29b-41d4-a716-446655440001",
    "broadcaster_kyc_tier": 1,
    "stream_type": "store",
    "title": "Product Launch Event",
    "description": "Join us for an exclusive look at our new collection",
    "privacy_setting": "public",
    "status": "scheduled",
    "stream_key": "rtmp://stream.tchat.dev/live/sk_abc123xyz789",
    "webrtc_session_id": "wrtc_session_abc123",
    "scheduled_start_time": "2025-09-30T15:00:00Z",
    "viewer_count": 0,
    "peak_viewer_count": 0,
    "max_capacity": 50000,
    "language": "en",
    "tags": ["commerce"],
    "created_at": "2025-09-30T14:00:00Z",
    "updated_at": "2025-09-30T14:00:00Z"
  }
}
```

### Error Response - Insufficient KYC (403 Forbidden)

```json
{
  "success": false,
  "error": "Store sellers require Standard KYC (Tier 1) verification",
  "required_kyc_tier": 1,
  "current_kyc_tier": 0,
  "verification_url": "https://tchat.dev/verify/kyc"
}
```

### Error Response - Validation Error (400 Bad Request)

```json
{
  "success": false,
  "error": "Validation error: missing required fields",
  "validation_errors": {
    "stream_type": "stream_type is required",
    "title": "title is required"
  }
}
```

## Business Rules Validated

1. **Store Stream KYC Requirement**: Users must have Standard KYC (Tier 1) verification to create store streams
2. **Video Stream Accessibility**: Video streams can be created by users with any KYC tier (including Tier 0)
3. **Stream Key Generation**: Each stream receives a unique RTMP stream key for broadcasting
4. **WebRTC Session**: Each stream receives a WebRTC session ID for browser-based streaming
5. **Default Capacity**: New streams default to 50,000 concurrent viewers
6. **Request Validation**: Required fields (stream_type, title) must be provided

## Next Steps (Implementation)

1. Implement POST /streams handler in `handlers/streaming_handler.go`
2. Add KYC validation middleware for store stream creation
3. Implement stream_key generation service (unique RTMP URLs)
4. Add WebRTC session management service
5. Create repository layer for live_streams database operations
6. Run provider verification tests against the implemented handler

## Related Files

- Contract specification: `/specs/029-implement-live-on/contracts/streaming-api.yaml`
- Database model: `/backend/streaming/models/live_stream.go`
- Handler (to be implemented): `/backend/streaming/handlers/streaming_handler.go`
- Repository (to be implemented): `/backend/streaming/repositories/stream_repository.go`

## Pact Documentation

- [Pact Go Documentation](https://github.com/pact-foundation/pact-go)
- [Pact Consumer Test Guide](https://docs.pact.io/implementation_guides/go/readme#consumer-testing)
- [Pact Matchers Reference](https://docs.pact.io/implementation_guides/go/readme#matchers)