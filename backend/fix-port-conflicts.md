# Backend Service Port Configuration Fix

## Issue Identified
The Tchat backend has port conflicts that prevent proper service communication:

### Current Port Conflicts
- **Gateway**: Running on port 8080 ✅
- **Commerce**: Also trying to use port 8080 ❌ (CONFLICT)
- **Auth**: Running on port 8081 ✅
- **Notification**: Running on port 8083 ✅
- **Content**: Running on port 8086 ✅

### Gateway Expected vs Actual Service Ports
The gateway (`/Users/weerawat/Tchat/backend/infrastructure/gateway/main.go`) expects:
- auth-service: port 8081 ✅ (matches actual)
- messaging-service: port 8082 ⚠️ (no service found)
- payment-service: port 8083 ❌ (notification is on 8083)
- commerce-service: port 8084 ❌ (commerce is on 8080, conflicts with gateway)
- notification-service: port 8085 ❌ (notification is on 8083)
- content-service: port 8086 ✅ (matches actual)

## Required Fixes

### 1. Fix Commerce Service Port
Change commerce service from port 8080 to 8084:
```bash
cd /Users/weerawat/Tchat/backend/commerce
# Set environment variable
export COMMERCE_SERVICE_PORT=8084
go run main.go
```

### 2. Fix Notification Service Port
Change notification service from port 8083 to 8085:
```bash
cd /Users/weerawat/Tchat/backend/notification
# Set environment variable
export NOTIFICATION_SERVICE_PORT=8085
go run main.go
```

### 3. Start Payment Service
The payment service should run on port 8083:
```bash
cd /Users/weerawat/Tchat/backend/payment
export PAYMENT_SERVICE_PORT=8083
go run main.go
```

### 4. Start Messaging Service
The messaging service should run on port 8082:
```bash
cd /Users/weerawat/Tchat/backend/messaging
export MESSAGING_SERVICE_PORT=8082
go run main.go
```

## Quick Fix Script
Create a script to start all services with correct ports:

```bash
#!/bin/bash
# Save as start-backend-services.sh

echo "Starting Tchat backend services with correct port configuration..."

# Kill existing processes
pkill -f "go run main.go"

# Start services in background with correct ports
cd /Users/weerawat/Tchat/backend/auth && AUTHE_SERVICE_PORT=8081 go run main.go &
cd /Users/weerawat/Tchat/backend/messaging && MESSAGING_SERVICE_PORT=8082 go run main.go &
cd /Users/weerawat/Tchat/backend/payment && PAYMENT_SERVICE_PORT=8083 go run main.go &
cd /Users/weerawat/Tchat/backend/commerce && COMMERCE_SERVICE_PORT=8084 go run main.go &
cd /Users/weerawat/Tchat/backend/notification && NOTIFICATION_SERVICE_PORT=8085 go run main.go &
cd /Users/weerawat/Tchat/backend/content && CONTENT_SERVICE_PORT=8086 go run main.go &

# Start gateway last
cd /Users/weerawat/Tchat/backend/infrastructure/gateway && GATEWAY_SERVICE_PORT=8080 go run main.go

echo "All services started. Gateway available at http://localhost:8080"
```

## Verification
After applying fixes, verify with:
```bash
lsof -i :8080 -i :8081 -i :8082 -i :8083 -i :8084 -i :8085 -i :8086
```

Expected output should show:
- gateway on 8080
- auth on 8081
- messaging on 8082
- payment on 8083
- commerce on 8084
- notification on 8085
- content on 8086