# Calling Service API Documentation

## Overview

The Calling Service provides voice and video calling functionality for the Tchat application. It includes WebRTC-based real-time communication, presence management, call session handling, and signaling coordination.

**Service Port**: 8093
**API Version**: v1
**Base URL**: `/api/v1`

## Authentication

All endpoints require JWT authentication via the `Authorization: Bearer <token>` header.

## Core Endpoints

### Call Management

#### POST /calls/initiate
Initiate a new voice or video call.

**Request Body:**
```json
{
  "caller_id": "string (UUID)",
  "callee_id": "string (UUID)",
  "call_type": "voice|video"
}
```

**Response:**
```json
{
  "call_id": "string (UUID)",
  "status": "connecting",
  "signaling_url": "ws://localhost:8093/ws/signaling/{call_id}",
  "created_at": "2024-01-01T12:00:00Z"
}
```

**Status Codes:**
- `201`: Call initiated successfully
- `400`: Invalid request parameters
- `401`: Unauthorized
- `409`: User already in an active call

#### PUT /calls/{call_id}/answer
Answer an incoming call.

**Request Body:**
```json
{
  "user_id": "string (UUID)"
}
```

**Response:**
```json
{
  "call_id": "string (UUID)",
  "status": "active",
  "participants": [
    {
      "user_id": "string (UUID)",
      "role": "caller|callee",
      "joined_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

#### DELETE /calls/{call_id}
End an active call.

**Request Body:**
```json
{
  "user_id": "string (UUID)",
  "reason": "user_hangup|connection_failed|timeout"
}
```

**Response:**
```json
{
  "call_id": "string (UUID)",
  "status": "ended",
  "duration": 180,
  "ended_at": "2024-01-01T12:03:00Z"
}
```

#### GET /calls/{call_id}/status
Get current call status and participants.

**Response:**
```json
{
  "call_id": "string (UUID)",
  "status": "connecting|active|ended|failed",
  "type": "voice|video",
  "initiated_by": "string (UUID)",
  "started_at": "2024-01-01T12:00:00Z",
  "duration": 180,
  "participants": [
    {
      "user_id": "string (UUID)",
      "role": "caller|callee",
      "joined_at": "2024-01-01T12:00:00Z",
      "connection_status": "connecting|connected|disconnected"
    }
  ]
}
```

### Active Calls

#### GET /calls/user/{user_id}/active
Get all active calls for a user.

**Response:**
```json
{
  "active_calls": [
    {
      "call_id": "string (UUID)",
      "type": "voice|video",
      "status": "connecting|active",
      "other_participants": ["string (UUID)"],
      "started_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

### Call History

#### GET /calls/user/{user_id}/history
Get call history for a user.

**Query Parameters:**
- `limit`: Maximum number of results (default: 50, max: 200)
- `offset`: Pagination offset (default: 0)
- `call_type`: Filter by call type (`voice|video`)
- `status`: Filter by call status (`completed|failed|missed`)

**Response:**
```json
{
  "calls": [
    {
      "call_id": "string (UUID)",
      "type": "voice|video",
      "status": "completed|failed|missed",
      "initiated_by": "string (UUID)",
      "participants": ["string (UUID)"],
      "started_at": "2024-01-01T12:00:00Z",
      "ended_at": "2024-01-01T12:03:00Z",
      "duration": 180
    }
  ],
  "pagination": {
    "limit": 50,
    "offset": 0,
    "total": 150,
    "has_more": true
  }
}
```

### Presence Management

#### GET /presence/status/{user_id}
Get current presence status for a user.

**Response:**
```json
{
  "user_id": "string (UUID)",
  "status": "online|offline",
  "in_call": true,
  "call_id": "string (UUID)",
  "last_seen": "2024-01-01T12:00:00Z"
}
```

#### PUT /presence/status
Update presence status for current user.

**Request Body:**
```json
{
  "status": "online|offline",
  "call_id": "string (UUID, optional)"
}
```

**Response:**
```json
{
  "user_id": "string (UUID)",
  "status": "online|offline",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

#### GET /presence/check/{user_id}
Check if a user is available for calling.

**Response:**
```json
{
  "user_id": "string (UUID)",
  "available": true,
  "status": "online|offline",
  "in_call": false,
  "last_seen": "2024-01-01T12:00:00Z"
}
```

## WebSocket Signaling

### Connection Endpoint
`ws://localhost:8093/ws/signaling/{call_id}`

### Authentication
Include JWT token in connection query parameter:
`ws://localhost:8093/ws/signaling/{call_id}?token={jwt_token}`

### Message Types

#### Join Room
```json
{
  "type": "join-room",
  "call_id": "string (UUID)",
  "user_id": "string (UUID)"
}
```

#### Leave Room
```json
{
  "type": "leave-room",
  "call_id": "string (UUID)",
  "user_id": "string (UUID)"
}
```

#### WebRTC Offer
```json
{
  "type": "offer",
  "call_id": "string (UUID)",
  "from": "string (UUID)",
  "to": "string (UUID)",
  "sdp": "string (SDP)"
}
```

#### WebRTC Answer
```json
{
  "type": "answer",
  "call_id": "string (UUID)",
  "from": "string (UUID)",
  "to": "string (UUID)",
  "sdp": "string (SDP)"
}
```

#### ICE Candidate
```json
{
  "type": "candidate",
  "call_id": "string (UUID)",
  "from": "string (UUID)",
  "to": "string (UUID)",
  "candidate": "string",
  "sdp_mid": "string",
  "sdp_m_line_index": "number"
}
```

#### User Joined
```json
{
  "type": "user-joined",
  "call_id": "string (UUID)",
  "user_id": "string (UUID)",
  "participants": ["string (UUID)"]
}
```

#### User Left
```json
{
  "type": "user-left",
  "call_id": "string (UUID)",
  "user_id": "string (UUID)",
  "participants": ["string (UUID)"]
}
```

#### Call Status Update
```json
{
  "type": "call-status",
  "call_id": "string (UUID)",
  "status": "connecting|active|ended|failed",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Health Check

#### GET /health
Service health check endpoint.

**Response:**
```json
{
  "status": "healthy",
  "service": "calling",
  "version": "1.0.0",
  "timestamp": "2024-01-01T12:00:00Z",
  "uptime": "2h30m15s"
}
```

#### GET /health/detailed
Detailed health check with dependency status.

**Response:**
```json
{
  "status": "healthy|degraded|unhealthy",
  "service": "calling",
  "version": "1.0.0",
  "timestamp": "2024-01-01T12:00:00Z",
  "uptime": "2h30m15s",
  "checks": {
    "database": {
      "status": "healthy",
      "message": "Database connection is healthy",
      "duration": "2.5ms"
    },
    "redis": {
      "status": "healthy",
      "message": "Redis connection is healthy",
      "duration": "1.2ms"
    },
    "webrtc": {
      "status": "healthy",
      "message": "WebRTC components are ready"
    }
  }
}
```

## Error Responses

All error responses follow this format:

```json
{
  "error": {
    "code": "string",
    "message": "string",
    "details": "object (optional)"
  },
  "timestamp": "2024-01-01T12:00:00Z",
  "request_id": "string (UUID)"
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `INVALID_REQUEST` | 400 | Invalid request parameters |
| `UNAUTHORIZED` | 401 | Authentication required |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `CONFLICT` | 409 | Resource conflict (e.g., user already in call) |
| `RATE_LIMITED` | 429 | Rate limit exceeded |
| `INTERNAL_ERROR` | 500 | Internal server error |
| `SERVICE_UNAVAILABLE` | 503 | Service temporarily unavailable |

### Specific Error Codes

| Code | Description |
|------|-------------|
| `USER_ALREADY_IN_CALL` | User is already participating in an active call |
| `CALL_NOT_FOUND` | Call session not found |
| `CALL_ALREADY_ENDED` | Cannot perform action on ended call |
| `USER_NOT_AVAILABLE` | Target user is not available for calling |
| `INVALID_CALL_TYPE` | Invalid call type specified |
| `WEBRTC_CONNECTION_FAILED` | WebRTC connection establishment failed |
| `SIGNALING_ERROR` | WebSocket signaling error |
| `PRESENCE_UPDATE_FAILED` | Failed to update user presence |

## Performance Characteristics

### Response Time Targets
- Call initiation: <5 seconds
- WebRTC signaling: <200ms latency
- Presence updates: <100ms
- History queries: <500ms

### Scalability Limits
- Concurrent calls: 1000+ simultaneous
- Signaling messages: 10,000+ messages/second
- History storage: 1M+ call records
- Memory usage: <100MB for long-duration calls

### Rate Limits
- Call initiation: 10 calls/minute per user
- Presence updates: 100 updates/minute per user
- History queries: 100 requests/minute per user
- WebSocket connections: 5 concurrent per user

## WebRTC Configuration

### STUN Servers
```json
{
  "iceServers": [
    {
      "urls": ["stun:stun.l.google.com:19302"]
    }
  ]
}
```

### TURN Servers (Optional)
```json
{
  "iceServers": [
    {
      "urls": ["turn:your-turn-server.com:3478"],
      "username": "username",
      "credential": "password"
    }
  ]
}
```

### Supported Codecs
- **Audio**: OPUS (48kHz), G.722, PCMU, PCMA
- **Video**: VP8, VP9, H.264 (if available)

## Security Considerations

### Authentication
- JWT tokens required for all endpoints
- WebSocket connections authenticated via query parameter
- Token validation against auth service

### Data Protection
- All signaling data encrypted in transit (WSS/HTTPS)
- Call metadata stored securely with encryption at rest
- User presence information access controlled

### Privacy
- Call history accessible only to participants
- Presence information limited to authorized users
- No call content stored or logged

## Testing

### Contract Tests
Pact contract tests validate API compatibility:
- Consumer: Web frontend, mobile apps
- Provider: Calling service
- Location: `backend/calling/tests/contract/`

### Performance Tests
Load and performance validation:
- Concurrent call testing: `tests/performance/load_test.go`
- WebRTC quality testing: `tests/performance/webrtc_quality_test.go`
- Memory usage testing: `tests/performance/memory_usage_test.go`

### Integration Tests
End-to-end testing with real dependencies:
- Database integration tests
- Redis integration tests
- WebSocket integration tests

## Monitoring and Observability

### Metrics
- Call success/failure rates
- Connection establishment times
- WebRTC connection quality
- Memory and CPU usage
- Signaling message throughput

### Logging
- Structured JSON logging
- Call session lifecycle events
- WebRTC connection events
- Error tracking and debugging

### Health Monitoring
- Service health endpoints
- Dependency health checks
- Performance metric alerts
- Real-time monitoring dashboards

---

**Last Updated**: $(date)
**API Version**: 1.0.0
**Service Version**: 1.0.0