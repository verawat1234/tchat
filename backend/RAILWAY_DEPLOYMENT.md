# Railway Deployment Guide for Tchat Microservices

Comprehensive guide to deploy your 6 microservices to Railway using a cost-optimized 3-project strategy.

## Architecture Overview

**Consolidated Projects:**
- **tchat-core**: Auth Service (8081) + API Gateway
- **tchat-messaging**: Messaging (8082) + Notification (8084) Services
- **tchat-commerce**: Payment (8083) + Commerce (8085) + Content (8086) Services

**Cost:** ~$15/month vs $30/month for individual services

## Prerequisites

1. **Railway Account**: [Sign up at railway.app](https://railway.app)
2. **Railway CLI**: Install globally
   ```bash
   npm install -g @railway/cli
   # OR
   curl -fsSL https://railway.app/install.sh | sh
   ```
3. **Git Repository**: Your code must be in a Git repository (GitHub/GitLab)

## Step 1: Railway CLI Setup

```bash
# Login to Railway
railway login

# Verify authentication
railway whoami
```

## Step 2: Deploy Project 1 - Core Services (Auth)

```bash
# Navigate to your backend directory
cd backend

# Create new Railway project
railway create tchat-core

# Link to your git repository
railway link

# Add PostgreSQL plugin for user data
railway add postgresql

# Add Redis plugin for sessions
railway add redis

# Set environment variables
railway variables set SERVICE_NAME=auth
railway variables set SERVICE_PORT=8081
railway variables set ENVIRONMENT=production
railway variables set LOG_LEVEL=info
railway variables set DB_NAME=tchat_auth

# Deploy the service
railway up --detach
```

**Environment Variables to Set in Railway Dashboard:**
```bash
# Database (auto-populated by PostgreSQL plugin)
DB_HOST=${{Postgres.PGHOST}}
DB_PORT=${{Postgres.PGPORT}}
DB_NAME=tchat_auth
DB_USER=${{Postgres.PGUSER}}
DB_PASSWORD=${{Postgres.PGPASSWORD}}
DB_SSL_MODE=require

# Redis (auto-populated by Redis plugin)
REDIS_HOST=${{Redis.REDIS_HOST}}
REDIS_PORT=${{Redis.REDIS_PORT}}
REDIS_PASSWORD=${{Redis.REDIS_PASSWORD}}
REDIS_DB=0

# External services (set manually)
TWILIO_ACCOUNT_SID=your_twilio_sid
TWILIO_AUTH_TOKEN=your_twilio_token
```

## Step 3: Deploy Project 2 - Messaging Services

```bash
# Create second Railway project
railway create tchat-messaging

# Switch to messaging context
railway environment

# Add Redis plugin for real-time data
railway add redis

# Add PostgreSQL plugin (for user references)
railway add postgresql

# Set environment variables
railway variables set SERVICE_NAME=messaging
railway variables set SERVICE_PORT=8082
railway variables set ENVIRONMENT=production
railway variables set LOG_LEVEL=info
railway variables set DB_NAME=tchat_messaging

# Deploy messaging service
railway up --detach
```

**Additional Environment Variables:**
```bash
# Auth service connection (use Railway URL from Project 1)
AUTH_SERVICE_URL=https://tchat-core.up.railway.app

# For real-time messaging (external services)
KAFKA_BROKERS=your_kafka_cluster:9092
KAFKA_GROUP_ID=tchat-messaging
KAFKA_TOPIC_PREFIX=tchat

# ScyllaDB (external or PostgreSQL fallback)
SCYLLA_HOSTS=your_scylla_cluster:9042
SCYLLA_KEYSPACE=tchat_messaging
SCYLLA_REPLICATION_FACTOR=3
```

## Step 4: Deploy Project 3 - Commerce Services

```bash
# Create third Railway project
railway create tchat-commerce

# Add PostgreSQL plugin for transactional data
railway add postgresql

# Add Redis plugin for caching
railway add redis

# Set environment variables
railway variables set SERVICE_NAME=payment
railway variables set SERVICE_PORT=8083
railway variables set ENVIRONMENT=production
railway variables set LOG_LEVEL=info
railway variables set DB_NAME=tchat_payment

# Deploy commerce services
railway up --detach
```

**Additional Environment Variables:**
```bash
# Payment gateway credentials
STRIPE_SECRET_KEY=your_stripe_key
OMISE_SECRET_KEY=your_omise_key

# Auth service connection
AUTH_SERVICE_URL=https://tchat-core.up.railway.app

# Internal service communication
COMMERCE_SERVICE_URL=internal
CONTENT_SERVICE_URL=internal
```

## Step 5: Database Migrations

Each project needs database migrations. Run these after deployment:

```bash
# Project 1 (Auth)
railway run --service tchat-core -- ./auth-service migrate up

# Project 2 (Messaging)
railway run --service tchat-messaging -- ./messaging-service migrate up

# Project 3 (Commerce)
railway run --service tchat-commerce -- ./payment-service migrate up
```

## Step 6: Configure Custom Domains (Optional)

```bash
# Add custom domains for each project
railway domain add api.tchat.sea --service tchat-core
railway domain add messaging.tchat.sea --service tchat-messaging
railway domain add commerce.tchat.sea --service tchat-commerce
```

## Step 7: Configure Service Networking

Update your service configurations to use Railway URLs for inter-service communication:

**In Railway Dashboard Environment Variables:**

**tchat-messaging project:**
```bash
AUTH_SERVICE_URL=https://tchat-core.up.railway.app
```

**tchat-commerce project:**
```bash
AUTH_SERVICE_URL=https://tchat-core.up.railway.app
MESSAGING_SERVICE_URL=https://tchat-messaging.up.railway.app
```

## Step 8: Health Check Validation

Test your deployed services:

```bash
# Test auth service
curl https://tchat-core.up.railway.app/health

# Test messaging service
curl https://tchat-messaging.up.railway.app/health

# Test commerce service
curl https://tchat-commerce.up.railway.app/health
```

## Step 9: Configure External Services

**For ScyllaDB/Kafka (if needed):**

Option 1: **Use Railway External Services**
- Connect external ScyllaDB/Kafka clusters via environment variables

Option 2: **PostgreSQL Fallback** (Recommended for getting started)
- Modify your messaging service to use PostgreSQL instead of ScyllaDB
- Use Railway's built-in Redis for real-time messaging instead of Kafka

## Production Considerations

### Security
- Enable Railway's **Environment Secrets** for sensitive data
- Use **Custom Domains** with SSL certificates
- Configure **CORS** properly for web client access

### Monitoring
- Enable Railway's **Usage Metrics**
- Set up **Health Check Alerts**
- Monitor **Memory and CPU Usage**

### Scaling
- Use Railway's **Auto-scaling** for high traffic
- Monitor **Response Times** and scale accordingly
- Consider **Load Balancer** for multiple replicas

## Troubleshooting

### Common Issues

1. **Service Won't Start**
   ```bash
   # Check logs
   railway logs --service tchat-core

   # Check environment variables
   railway variables
   ```

2. **Database Connection Failed**
   ```bash
   # Verify PostgreSQL plugin is added
   railway plugins

   # Check database connection string
   railway shell -- echo $DATABASE_URL
   ```

3. **Service Communication Issues**
   ```bash
   # Verify internal URLs are correct
   railway variables | grep SERVICE_URL

   # Test connectivity
   railway run -- curl $AUTH_SERVICE_URL/health
   ```

## Cost Optimization

**Monthly Estimates:**
- **3 Projects**: ~$15/month (3 × $5)
- **Database Plugins**: ~$5/month each (PostgreSQL + Redis per project)
- **Total**: ~$30/month for full microservice deployment

**Scaling Costs:**
- Auto-scaling triggers additional charges
- Monitor usage to avoid unexpected costs
- Use Railway's **Usage Alerts** to track spending

## Next Steps

1. Deploy projects in order: Core → Messaging → Commerce
2. Configure environment variables via Railway dashboard
3. Test each service after deployment
4. Set up monitoring and alerts
5. Configure custom domains if needed

Your microservices will be accessible at:
- **Auth**: `https://tchat-core.up.railway.app`
- **Messaging**: `https://tchat-messaging.up.railway.app`
- **Commerce**: `https://tchat-commerce.up.railway.app`