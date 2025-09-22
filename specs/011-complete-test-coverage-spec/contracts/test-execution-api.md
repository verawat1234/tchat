# Test Execution API Contract

## Overview
API contract for executing and managing test suites across Tchat microservices.

## Base URL
```
POST /api/v1/tests
GET  /api/v1/tests
```

## Test Suite Execution

### Execute Test Suite
```http
POST /api/v1/tests/suites/{suiteId}/execute
```

**Request Body:**
```json
{
  "environment": "integration",
  "filters": {
    "tags": ["critical", "smoke"],
    "priorities": ["critical", "high"],
    "testTypes": ["unit", "integration"]
  },
  "configuration": {
    "parallel": true,
    "timeout": "30m",
    "retryCount": 2,
    "coverage": true
  }
}
```

**Response 200:**
```json
{
  "executionId": "exec_123456789",
  "suiteId": "suite_content_service",
  "status": "running",
  "startedAt": "2025-09-22T10:00:00Z",
  "estimatedDuration": "15m",
  "progress": {
    "total": 45,
    "completed": 0,
    "running": 5,
    "pending": 40
  },
  "environment": "integration"
}
```

**Response 400:**
```json
{
  "error": "invalid_request",
  "message": "Invalid test suite configuration",
  "details": {
    "field": "configuration.timeout",
    "reason": "timeout must be between 1m and 60m"
  }
}
```

### Get Test Execution Status
```http
GET /api/v1/tests/executions/{executionId}
```

**Response 200:**
```json
{
  "executionId": "exec_123456789",
  "suiteId": "suite_content_service",
  "status": "completed",
  "startedAt": "2025-09-22T10:00:00Z",
  "completedAt": "2025-09-22T10:12:34Z",
  "duration": "12m34s",
  "results": {
    "total": 45,
    "passed": 42,
    "failed": 2,
    "skipped": 1,
    "errors": 0
  },
  "coverage": {
    "lines": 87.5,
    "functions": 92.1,
    "branches": 84.3
  },
  "failures": [
    {
      "testId": "test_content_create_validation",
      "name": "TestContentHandlers_CreateContent_InvalidData_ReturnsError",
      "error": "assertion failed: expected status 400, got 500",
      "duration": "125ms"
    }
  ]
}
```

## Coverage Reporting

### Get Coverage Report
```http
GET /api/v1/tests/coverage/{serviceId}
```

**Query Parameters:**
- `type`: unit | integration | combined
- `format`: json | html | xml
- `historical`: true | false

**Response 200:**
```json
{
  "serviceId": "content-service",
  "reportType": "combined",
  "generatedAt": "2025-09-22T10:15:00Z",
  "overall": {
    "lines": 87.5,
    "functions": 92.1,
    "branches": 84.3,
    "statements": 89.2
  },
  "byTestType": {
    "unit": {
      "lines": 82.1,
      "functions": 88.9,
      "branches": 79.6
    },
    "integration": {
      "lines": 91.2,
      "functions": 94.7,
      "branches": 87.8
    }
  },
  "files": [
    {
      "path": "handlers/content_handlers.go",
      "lines": 95.2,
      "functions": 100.0,
      "branches": 92.3,
      "uncoveredLines": [45, 67, 89]
    }
  ],
  "trends": {
    "lastWeek": 85.2,
    "lastMonth": 82.7,
    "direction": "increasing"
  }
}
```

## Performance Testing

### Execute Performance Test
```http
POST /api/v1/tests/performance/{serviceId}/execute
```

**Request Body:**
```json
{
  "testType": "load",
  "configuration": {
    "duration": "10m",
    "concurrentUsers": 1000,
    "rampUpTime": "2m",
    "regions": ["TH", "SG", "ID"],
    "scenarios": [
      {
        "name": "content_creation_flow",
        "weight": 60,
        "endpoints": ["/api/v1/content", "/api/v1/content/{id}"]
      }
    ]
  },
  "targets": {
    "responseTime95th": "200ms",
    "errorRate": "1%",
    "throughput": "500rps"
  }
}
```

**Response 200:**
```json
{
  "executionId": "perf_987654321",
  "status": "running",
  "progress": {
    "elapsedTime": "3m45s",
    "remainingTime": "6m15s",
    "currentUsers": 750,
    "totalRequests": 125000
  },
  "realTimeMetrics": {
    "responseTime": {
      "avg": "145ms",
      "p95": "189ms",
      "p99": "245ms"
    },
    "throughput": "520rps",
    "errorRate": "0.3%"
  }
}
```

## Security Testing

### Execute Security Test Suite
```http
POST /api/v1/tests/security/{serviceId}/execute
```

**Request Body:**
```json
{
  "testCategories": ["authentication", "authorization", "input_validation"],
  "complianceStandards": ["GDPR", "PDPA_TH", "PDPA_SG"],
  "regions": ["TH", "SG", "ID", "MY", "PH", "VN"],
  "severity": ["critical", "high"],
  "configuration": {
    "deepScan": true,
    "includeCompliance": true,
    "generateReport": true
  }
}
```

**Response 200:**
```json
{
  "executionId": "sec_456789123",
  "status": "completed",
  "results": {
    "total": 28,
    "passed": 26,
    "failed": 2,
    "critical": 0,
    "high": 2,
    "medium": 0,
    "low": 0
  },
  "vulnerabilities": [
    {
      "id": "vuln_auth_001",
      "severity": "high",
      "category": "authentication",
      "title": "JWT Token Validation Bypass",
      "description": "Weak signature validation in JWT middleware",
      "impact": "Unauthorized access to protected endpoints",
      "remediation": "Strengthen JWT signature validation",
      "cwe": "CWE-287"
    }
  ],
  "compliance": {
    "GDPR": {
      "status": "compliant",
      "lastValidated": "2025-09-22T10:00:00Z",
      "evidence": ["test_data_deletion", "test_data_export"]
    },
    "PDPA_TH": {
      "status": "non_compliant",
      "issues": ["Data residency validation failed"],
      "requiredActions": ["Implement Thailand data storage validation"]
    }
  }
}
```

## Test Management

### List Test Suites
```http
GET /api/v1/tests/suites
```

**Query Parameters:**
- `serviceId`: Filter by service
- `type`: Filter by test type
- `status`: active | inactive
- `page`: Page number
- `limit`: Results per page

**Response 200:**
```json
{
  "suites": [
    {
      "id": "suite_content_service",
      "name": "Content Service Test Suite",
      "serviceId": "content-service",
      "type": "unit",
      "status": "active",
      "testCount": 45,
      "lastExecution": "2025-09-22T09:30:00Z",
      "lastResult": "passed",
      "coverage": 87.5
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 35,
    "pages": 2
  }
}
```

### Create Test Suite
```http
POST /api/v1/tests/suites
```

**Request Body:**
```json
{
  "name": "Payment Service Security Tests",
  "serviceId": "payment-service",
  "type": "security",
  "description": "Comprehensive security testing for payment service",
  "configuration": {
    "timeout": "30m",
    "retryPolicy": {
      "maxRetries": 3,
      "backoffStrategy": "exponential"
    },
    "environment": "staging",
    "tags": ["security", "payment", "compliance"]
  },
  "tests": [
    {
      "name": "test_payment_encryption",
      "function": "TestPaymentEncryption_ValidCard_EncryptsData",
      "priority": "critical",
      "tags": ["encryption", "pci_dss"]
    }
  ]
}
```

**Response 201:**
```json
{
  "id": "suite_payment_security",
  "name": "Payment Service Security Tests",
  "serviceId": "payment-service",
  "type": "security",
  "status": "active",
  "testCount": 1,
  "createdAt": "2025-09-22T10:20:00Z",
  "configuration": {
    "timeout": "30m",
    "retryPolicy": {
      "maxRetries": 3,
      "backoffStrategy": "exponential"
    }
  }
}
```

## Error Responses

### Standard Error Format
```json
{
  "error": "error_code",
  "message": "Human readable error message",
  "details": {
    "field": "specific_field",
    "reason": "detailed_reason"
  },
  "timestamp": "2025-09-22T10:00:00Z",
  "requestId": "req_123456789"
}
```

### Common Error Codes
- `invalid_request`: Request validation failed
- `suite_not_found`: Test suite does not exist
- `execution_failed`: Test execution encountered errors
- `insufficient_permissions`: User lacks required permissions
- `rate_limit_exceeded`: Too many requests
- `service_unavailable`: Testing service temporarily unavailable

## Authentication
All endpoints require Bearer token authentication:
```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Rate Limits
- Test execution: 10 requests per minute per service
- Coverage reports: 60 requests per minute
- General API: 1000 requests per minute

## Webhooks
Optional webhook notifications for test completion:
```json
{
  "event": "test.execution.completed",
  "executionId": "exec_123456789",
  "suiteId": "suite_content_service",
  "status": "completed",
  "results": {
    "passed": 42,
    "failed": 2,
    "coverage": 87.5
  },
  "timestamp": "2025-09-22T10:15:00Z"
}
```

This contract defines the API for comprehensive test execution, coverage reporting, and security validation across the Tchat microservices platform.