# Tchat Environment Variables Setup Commands

## Complete Railway CLI Commands for Environment Setup

Execute these commands to set up all environment variables for your Tchat microservices with dedicated databases.

### Prerequisites
```bash
npm install -g @railway/cli
railway login
```

---

## 🌐 Gateway Service (Port 8080)
```bash
railway link 0a1f3508-2150-4d0c-8ae9-878f74a607a0 --service gateway
railway variables set PORT="8080"
railway variables set GATEWAY_PORT="8080"
railway variables set ENVIRONMENT="production"
railway variables set LOG_LEVEL="info"
railway variables set DATABASE_URL="postgresql://postgres:TFMPTDlpHkWKQuCqmkCgopOMWxherEvA@postgres-7jdx.railway.internal:5432/railway"
railway variables set REDIS_URL="redis://default:jOFuvSfgpVbdzhbzzMyThAxjMfEHxsKT@redis.railway.internal:6379"
railway variables set AUTH_SERVICE_URL="http://auth.railway.internal:8081"
railway variables set CONTENT_SERVICE_URL="http://content.railway.internal:8086"
railway variables set MESSAGING_SERVICE_URL="http://messaging.railway.internal:8082"
railway variables set VIDEO_SERVICE_URL="http://video.railway.internal:8091"
railway variables set SOCIAL_SERVICE_URL="http://social.railway.internal:8092"
```

---

## 🔐 Auth Service (Port 8081)
```bash
railway link 0a1f3508-2150-4d0c-8ae9-878f74a607a0 --service auth
railway variables set PORT="8081"
railway variables set ENVIRONMENT="production"
railway variables set LOG_LEVEL="info"
railway variables set DATABASE_URL="postgresql://postgres:BpcMkwzFeULuAINVIRScuCBfNwQaqsyo@postgres.railway.internal:5432/railway"
railway variables set REDIS_URL="redis://default:jOFuvSfgpVbdzhbzzMyThAxjMfEHxsKT@redis.railway.internal:6379"
railway variables set JWT_SECRET="tchat-super-secret-jwt-key-change-in-production-2024"
railway variables set TOKEN_EXPIRY="24h"
railway variables set GATEWAY_URL="http://gateway.railway.internal:8080"
```

---

## 📄 Content Service (Port 8086)
```bash
railway link 0a1f3508-2150-4d0c-8ae9-878f74a607a0 --service content
railway variables set PORT="8086"
railway variables set ENVIRONMENT="production"
railway variables set LOG_LEVEL="info"
railway variables set DATABASE_URL="postgresql://postgres:hzvjxZBAemJrTedKCQhENgBWcjBwbOLR@postgres-rxgw.railway.internal:5432/railway"
railway variables set REDIS_URL="redis://default:jOFuvSfgpVbdzhbzzMyThAxjMfEHxsKT@redis.railway.internal:6379"
railway variables set AUTH_SERVICE_URL="http://auth.railway.internal:8081"
railway variables set GATEWAY_URL="http://gateway.railway.internal:8080"
```

---

## 💬 Messaging Service (Port 8082)
```bash
railway link 0a1f3508-2150-4d0c-8ae9-878f74a607a0 --service messaging
railway variables set PORT="8082"
railway variables set WEBSOCKET_PORT="8082"
railway variables set ENVIRONMENT="production"
railway variables set LOG_LEVEL="info"
railway variables set DATABASE_URL="postgresql://postgres:GDxpAYSmtkuiVLESvwscIJCDVqxhjnIw@postgres-kcfi.railway.internal:5432/railway"
railway variables set REDIS_URL="redis://default:jOFuvSfgpVbdzhbzzMyThAxjMfEHxsKT@redis.railway.internal:6379"
railway variables set AUTH_SERVICE_URL="http://auth.railway.internal:8081"
railway variables set GATEWAY_URL="http://gateway.railway.internal:8080"
```

---

## 🎥 Video Service (Port 8091)
```bash
railway link 0a1f3508-2150-4d0c-8ae9-878f74a607a0 --service video
railway variables set PORT="8091"
railway variables set ENVIRONMENT="production"
railway variables set LOG_LEVEL="info"
railway variables set DATABASE_URL="postgresql://postgres:gsQvUsgSEjARBJuHXzXuOFnHgtawXkGt@postgres-_oks.railway.internal:5432/railway"
railway variables set REDIS_URL="redis://default:jOFuvSfgpVbdzhbzzMyThAxjMfEHxsKT@redis.railway.internal:6379"
railway variables set VIDEO_UPLOAD_PATH="/uploads"
railway variables set MAX_FILE_SIZE="100MB"
railway variables set AUTH_SERVICE_URL="http://auth.railway.internal:8081"
railway variables set GATEWAY_URL="http://gateway.railway.internal:8080"
```

---

## 👥 Social Service (Port 8092)
```bash
railway link 0a1f3508-2150-4d0c-8ae9-878f74a607a0 --service social
railway variables set PORT="8092"
railway variables set SOCIAL_PORT="8092"
railway variables set ENVIRONMENT="production"
railway variables set LOG_LEVEL="info"
railway variables set DATABASE_URL="postgresql://postgres:YGLKJkIRWZSFupkCtNyaYzYfotaUZbjj@postgres-mo3g.railway.internal:5432/railway"
railway variables set REDIS_URL="redis://default:jOFuvSfgpVbdzhbzzMyThAxjMfEHxsKT@redis.railway.internal:6379"
railway variables set AUTH_SERVICE_URL="http://auth.railway.internal:8081"
railway variables set GATEWAY_URL="http://gateway.railway.internal:8080"
```

---

## 🚀 Deployment Commands

After setting environment variables, deploy each service:

```bash
# Deploy all services
railway deploy --service gateway
railway deploy --service auth
railway deploy --service content
railway deploy --service messaging
railway deploy --service video
railway deploy --service social
```

---

## 📊 Database-Service Mapping Summary

| Service    | Database      | Internal URL                              | Port |
|------------|---------------|-------------------------------------------|------|
| gateway    | postgres-7jdx | postgres-7jdx.railway.internal:5432      | 8080 |
| auth       | postgres      | postgres.railway.internal:5432           | 8081 |
| content    | postgres-rxgw | postgres-rxgw.railway.internal:5432      | 8086 |
| messaging  | postgres-kcfi | postgres-kcfi.railway.internal:5432      | 8082 |
| video      | postgres-_oks | postgres-_oks.railway.internal:5432      | 8091 |
| social     | postgres-mo3g | postgres-mo3g.railway.internal:5432      | 8092 |
| shared     | redis         | redis.railway.internal:6379              | -    |

---

## 🔍 Verification Commands

Check deployment status:
```bash
railway status
railway logs --service gateway
railway logs --service auth
railway logs --service content
railway logs --service messaging
railway logs --service video
railway logs --service social
```

Test connectivity:
```bash
# Test Gateway API
curl https://gateway-production-xxxx.up.railway.app/health

# Test individual services
curl http://auth.railway.internal:8081/health
curl http://content.railway.internal:8086/health
curl http://messaging.railway.internal:8082/health
curl http://video.railway.internal:8091/health
curl http://social.railway.internal:8092/health
```

---

## 🔧 Next Steps

1. **Execute Commands**: Run the commands above section by section
2. **Deploy Services**: Use the deployment commands
3. **Database Migrations**: Run migrations for each service
4. **Test Connectivity**: Verify all services are communicating
5. **Monitor Logs**: Check for any startup issues

---

## 📝 Environment Variables Reference

### Common Variables (All Services)
- `PORT`: Service port
- `ENVIRONMENT`: production
- `LOG_LEVEL`: info
- `DATABASE_URL`: Service-specific PostgreSQL URL
- `REDIS_URL`: Shared Redis cache

### Service-Specific Variables
- **Gateway**: Service URLs for routing
- **Auth**: JWT configuration
- **Messaging**: WebSocket configuration
- **Video**: Upload and file size limits
- **Social**: Social-specific port configuration

---

## 🔗 Useful Links

- **Railway Project**: https://railway.app/project/0a1f3508-2150-4d0c-8ae9-878f74a607a0
- **Railway Dashboard**: Access services, view logs, monitor deployments
- **Repository**: weerawat-tchat/Tchat