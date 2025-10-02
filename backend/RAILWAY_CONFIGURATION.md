# Railway Configuration Guide - Complete Fix

## Problem Summary

ALL services failed because:
1. ❌ Root directories were set to `backend/<service>` (can't access shared module)
2. ❌ Dockerfiles expected `COPY shared/` which doesn't exist in service-specific build context
3. ❌ Services were on `master` branch instead of `develop`

## Solution

✅ **All Dockerfiles fixed** - Now properly copy `shared/` and `<service>/` directories
✅ **Develop branch created** - Commit `7e0dedc` with all fixes
✅ **Configuration guide created** - Manual Railway UI configuration required

## Railway Configuration Required

### Step-by-Step for Each Service

**Access Railway Dashboard:**
https://railway.app/project/0a1f3508-2150-4d0c-8ae9-878f74a607a0

---

### Service 1: auth-service

1. Click on **"auth-service"**
2. Go to **Settings** → **Source**
3. Configure:
   - **Branch**: `develop` (change from `master`)
   - **Root Directory**: `backend` (change from `backend/auth`)
   - **Dockerfile Path**: `auth/Dockerfile`
4. Click **Save**
5. Railway will auto-deploy with commit `7e0dedc`

---

### Service 2: messaging-service

1. Click on **"messaging-service"**
2. Go to **Settings** → **Source**
3. Configure:
   - **Branch**: `develop`
   - **Root Directory**: `backend` (change from `backend/messaging`)
   - **Dockerfile Path**: `messaging/Dockerfile`
4. Click **Save**

---

### Service 3: social-service

1. Click on **"social-service"**
2. Go to **Settings** → **Source**
3. Configure:
   - **Branch**: `develop`
   - **Root Directory**: `backend` (change from `backend/social`)
   - **Dockerfile Path**: `social/Dockerfile`
4. Click **Save**

---

### Service 4: content-service

1. Click on **"content-service"**
2. Go to **Settings** → **Source**
3. Configure:
   - **Branch**: `develop`
   - **Root Directory**: `backend` (change from `backend/content`)
   - **Dockerfile Path**: `content/Dockerfile`
4. Click **Save**

---

### Service 5: commerce-service

1. Click on **"commerce-service"**
2. Go to **Settings** → **Source**
3. Configure:
   - **Branch**: `develop`
   - **Root Directory**: `backend` (change from `backend/commerce`)
   - **Dockerfile Path**: `commerce/Dockerfile`
4. Click **Save**

---

### Service 6: payment-service

1. Click on **"payment-service"**
2. Go to **Settings** → **Source**
3. Configure:
   - **Branch**: `develop`
   - **Root Directory**: `backend` (change from `backend/payment`)
   - **Dockerfile Path**: `payment/Dockerfile`
4. Click **Save**

---

### Service 7: notification-service

1. Click on **"notification-service"**
2. Go to **Settings** → **Source**
3. Configure:
   - **Branch**: `develop`
   - **Root Directory**: `backend` (change from `backend/notification`)
   - **Dockerfile Path**: `notification/Dockerfile`
4. Click **Save**

---

### Service 8: video-service

1. Click on **"video-service"**
2. Go to **Settings** → **Source**
3. Configure:
   - **Branch**: `develop`
   - **Root Directory**: `backend` (change from `backend/video`)
   - **Dockerfile Path**: `video/Dockerfile`
4. Click **Save**

---

### Service 9: calling-service

1. Click on **"calling-service"**
2. Go to **Settings** → **Source**
3. Configure:
   - **Branch**: `develop`
   - **Root Directory**: `backend` (change from `backend/calling`)
   - **Dockerfile Path**: `calling/Dockerfile`
4. Click **Save**

---

### Service 10: gateway-service

1. Click on **"gateway-service"**
2. Go to **Settings** → **Source**
3. Configure:
   - **Branch**: `develop`
   - **Root Directory**: `backend/infrastructure/gateway` (keep as is - has local shared copy)
   - **Dockerfile Path**: `Dockerfile`
4. Click **Save**

---

## Verification Checklist

After configuring all services, verify:

### 1. Build Success
Check each service's deployment logs:
```
✅ [build] FROM docker.io/library/golang:1.24-alpine
✅ [build] COPY shared/ ./shared/
✅ [build] COPY <service>/ ./<service>/
✅ [build] RUN go mod tidy
✅ [build] go: found tchat.dev/shared/config
✅ [build] RUN go build -o <service>-service ./main.go
✅ [build] Successfully built
```

### 2. Deployment Success
All services should show:
```
✅ Status: SUCCESS
✅ Health: Healthy
✅ Replicas: 1/1 running
```

### 3. Service URLs
Get the Railway URLs for:
- **gateway-service** (port 8080) - Main entry point
- **auth-service** (port 8081) - Authentication

## Testing Authentication Flow

Once all services are deployed:

```bash
# 1. Test gateway health
curl https://<gateway-url>/health

# 2. Request OTP
curl -X POST https://<gateway-url>/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+66812345678"}'

# 3. Verify OTP and get tokens
curl -X POST https://<gateway-url>/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+66812345678", "otp": "123456"}'

# 4. Test authenticated request
curl https://<gateway-url>/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>"
```

## Summary

**Fixed:**
- ✅ All 10 Dockerfiles now correctly copy shared module
- ✅ All services use Go 1.24-alpine
- ✅ Develop branch created with all fixes (commit: 7e0dedc)

**Required:**
- ⚠️ Manual Railway UI configuration for all 10 services
- ⚠️ Switch branch to `develop`
- ⚠️ Update root directories (9 services to `backend`, gateway to `backend/infrastructure/gateway`)

**Expected Result:**
- 🎯 All 10 services build and deploy successfully
- 🎯 Authentication flow working end-to-end
- 🎯 All microservices operational
