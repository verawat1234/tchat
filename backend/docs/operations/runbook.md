# Tchat Operations Runbook

This runbook provides step-by-step procedures for common operational tasks and incident response.

## Table of Contents

1. [Service Health Monitoring](#service-health-monitoring)
2. [Incident Response](#incident-response)
3. [Scaling Procedures](#scaling-procedures)
4. [Database Operations](#database-operations)
5. [Deployment Procedures](#deployment-procedures)
6. [Backup and Recovery](#backup-and-recovery)
7. [Performance Tuning](#performance-tuning)
8. [Security Procedures](#security-procedures)

## Service Health Monitoring

### Health Check Endpoints

All services expose standardized health endpoints:

- `GET /health` - Basic service health
- `GET /health/ready` - Service and dependencies health
- `GET /metrics` - Prometheus metrics

### Health Check Scripts

```bash
#!/bin/bash
# health-check.sh

SERVICES=("auth:8081" "messaging:8082" "payment:8083" "notification:8084")
BASE_URL="${BASE_URL:-http://localhost}"

for service in "${SERVICES[@]}"; do
    IFS=':' read -r name port <<< "$service"
    echo "Checking $name service..."

    response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL:$port/health")

    if [ "$response" = "200" ]; then
        echo "✅ $name service is healthy"
    else
        echo "❌ $name service is unhealthy (HTTP $response)"
        exit 1
    fi
done

echo "All services are healthy"
```

### Monitoring Alerts

#### Critical Alerts (Page immediately)

1. **Service Down**
   - Query: `up{job="tchat-services"} == 0`
   - Threshold: Any service down for >1 minute
   - Action: Page on-call engineer

2. **High Error Rate**
   - Query: `rate(http_requests_total{status=~"5.."}[5m]) > 0.05`
   - Threshold: >5% error rate for 5 minutes
   - Action: Page on-call engineer

3. **Database Connection Failures**
   - Query: `rate(db_connection_errors_total[5m]) > 0.1`
   - Threshold: >0.1 errors/second for 5 minutes
   - Action: Page on-call engineer

#### Warning Alerts (Notify during business hours)

1. **High Latency**
   - Query: `histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1`
   - Threshold: P95 latency >1 second
   - Action: Investigate during business hours

2. **Memory Usage**
   - Query: `container_memory_usage_bytes / container_spec_memory_limit_bytes > 0.8`
   - Threshold: >80% memory usage
   - Action: Consider scaling up

3. **Disk Usage**
   - Query: `(node_filesystem_size_bytes - node_filesystem_free_bytes) / node_filesystem_size_bytes > 0.85`
   - Threshold: >85% disk usage
   - Action: Clean up or scale storage

## Incident Response

### Incident Severity Levels

- **P0 (Critical)**: Complete service outage, data loss
- **P1 (High)**: Major functionality broken, significant user impact
- **P2 (Medium)**: Minor functionality issues, limited user impact
- **P3 (Low)**: Cosmetic issues, no user impact

### Incident Response Procedures

#### P0/P1 Incident Response

1. **Immediate Response (0-5 minutes)**
   ```bash
   # Check overall system health
   ./scripts/health-check.sh

   # Check service status
   kubectl get pods -n tchat
   # or
   docker-compose ps

   # Check recent deployments
   kubectl rollout history deployment -n tchat
   ```

2. **Investigation (5-15 minutes)**
   ```bash
   # Check logs for errors
   kubectl logs -f deployment/auth-service -n tchat --tail=100

   # Check metrics
   curl http://prometheus:9090/api/v1/query?query=up{job="tchat-services"}

   # Check external dependencies
   ./scripts/check-dependencies.sh
   ```

3. **Mitigation (15-30 minutes)**
   ```bash
   # Rollback if recent deployment
   kubectl rollout undo deployment/auth-service -n tchat

   # Scale up if resource issue
   kubectl scale deployment auth-service --replicas=5 -n tchat

   # Restart unhealthy pods
   kubectl delete pod -l app=auth-service -n tchat
   ```

4. **Communication**
   - Update status page
   - Notify stakeholders
   - Post in #incidents Slack channel

### Common Issues and Solutions

#### 1. Service Not Responding

**Symptoms**: Health checks failing, 502/503 errors

**Investigation**:
```bash
# Check pod status
kubectl describe pod <pod-name> -n tchat

# Check logs
kubectl logs <pod-name> -n tchat --previous

# Check resource usage
kubectl top pod <pod-name> -n tchat
```

**Solutions**:
```bash
# Restart pod
kubectl delete pod <pod-name> -n tchat

# Scale up if resource constrained
kubectl scale deployment <service> --replicas=<count> -n tchat

# Check and fix configuration
kubectl edit configmap tchat-config -n tchat
```

#### 2. Database Connection Issues

**Symptoms**: Connection timeout errors, high connection count

**Investigation**:
```bash
# Check PostgreSQL connections
psql -h postgres-host -U username -c "SELECT count(*) FROM pg_stat_activity;"

# Check Redis connections
redis-cli info clients

# Check ScyllaDB status
nodetool status
```

**Solutions**:
```bash
# Restart services to reset connection pools
kubectl rollout restart deployment/auth-service -n tchat

# Increase connection limits
kubectl patch configmap tchat-config -n tchat --patch '{"data":{"db_max_connections":"50"}}'

# Scale database if needed
kubectl scale statefulset postgres --replicas=2 -n tchat
```

#### 3. High Memory Usage

**Symptoms**: Pods getting OOMKilled, slow response times

**Investigation**:
```bash
# Check memory usage
kubectl top pods -n tchat

# Check memory leaks
curl http://service:port/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

**Solutions**:
```bash
# Increase memory limits
kubectl patch deployment auth-service -n tchat --patch '{"spec":{"template":{"spec":{"containers":[{"name":"auth-service","resources":{"limits":{"memory":"1Gi"}}}]}}}}'

# Restart services
kubectl rollout restart deployment -n tchat

# Enable garbage collection tuning
kubectl patch deployment auth-service -n tchat --patch '{"spec":{"template":{"spec":{"containers":[{"name":"auth-service","env":[{"name":"GOGC","value":"50"}]}]}}}}'
```

## Scaling Procedures

### Horizontal Scaling

#### Manual Scaling

```bash
# Scale individual service
kubectl scale deployment auth-service --replicas=5 -n tchat

# Scale all services
for service in auth-service messaging-service payment-service notification-service; do
    kubectl scale deployment $service --replicas=3 -n tchat
done
```

#### Auto-scaling Setup

```yaml
# hpa.yml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: auth-service-hpa
  namespace: tchat
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: auth-service
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### Vertical Scaling

```bash
# Increase CPU and memory limits
kubectl patch deployment auth-service -n tchat --patch '{
  "spec": {
    "template": {
      "spec": {
        "containers": [{
          "name": "auth-service",
          "resources": {
            "requests": {"cpu": "500m", "memory": "1Gi"},
            "limits": {"cpu": "1000m", "memory": "2Gi"}
          }
        }]
      }
    }
  }
}'
```

### Database Scaling

#### PostgreSQL Read Replicas

```yaml
# postgres-replica.yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres-replica
  namespace: tchat
spec:
  replicas: 2
  selector:
    matchLabels:
      app: postgres-replica
  template:
    metadata:
      labels:
        app: postgres-replica
    spec:
      containers:
      - name: postgres
        image: postgres:15-alpine
        env:
        - name: PGUSER
          value: postgres
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: password
        - name: PGPASSWORD
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: password
        command:
        - bash
        - -c
        - |
          pg_basebackup -h postgres-master -D /var/lib/postgresql/data -U postgres -v -P -W
          echo "hot_standby = on" >> /var/lib/postgresql/data/postgresql.conf
          echo "primary_conninfo = 'host=postgres-master port=5432 user=postgres'" > /var/lib/postgresql/data/recovery.conf
          postgres
```

#### ScyllaDB Cluster Expansion

```bash
# Add new ScyllaDB node
kubectl scale statefulset scylla --replicas=4 -n tchat

# Check cluster status
kubectl exec -it scylla-0 -n tchat -- nodetool status

# Repair data
kubectl exec -it scylla-0 -n tchat -- nodetool repair
```

## Database Operations

### PostgreSQL Maintenance

#### Backup

```bash
#!/bin/bash
# postgres-backup.sh

BACKUP_DIR="/backups/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
DATABASES=("tchat_auth" "tchat_messaging" "tchat_payment" "tchat_notification")

mkdir -p $BACKUP_DIR

for db in "${DATABASES[@]}"; do
    echo "Backing up $db..."
    pg_dump -h $POSTGRES_HOST -U $POSTGRES_USER $db | gzip > $BACKUP_DIR/${db}_${DATE}.sql.gz

    if [ $? -eq 0 ]; then
        echo "✅ Backup completed for $db"
    else
        echo "❌ Backup failed for $db"
        exit 1
    fi
done

# Clean up old backups (keep last 7 days)
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete
```

#### Database Migration

```bash
#!/bin/bash
# migrate.sh

MIGRATION_DIR="migrations/postgres"
DATABASES=("tchat_auth" "tchat_messaging" "tchat_payment" "tchat_notification")

for db in "${DATABASES[@]}"; do
    echo "Running migrations for $db..."
    migrate -path $MIGRATION_DIR -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:5432/$db?sslmode=require" up

    if [ $? -eq 0 ]; then
        echo "✅ Migrations completed for $db"
    else
        echo "❌ Migrations failed for $db"
        exit 1
    fi
done
```

### ScyllaDB Maintenance

#### Backup

```bash
#!/bin/bash
# scylla-backup.sh

KEYSPACES=("tchat_auth" "tchat_messaging" "tchat_payment" "tchat_notification")
DATE=$(date +%Y%m%d_%H%M%S)

for keyspace in "${KEYSPACES[@]}"; do
    echo "Creating snapshot for $keyspace..."
    kubectl exec -it scylla-0 -n tchat -- nodetool snapshot $keyspace -t snapshot_$DATE
done

echo "Snapshots created successfully"
```

#### Repair

```bash
#!/bin/bash
# scylla-repair.sh

KEYSPACES=("tchat_auth" "tchat_messaging" "tchat_payment" "tchat_notification")

for keyspace in "${KEYSPACES[@]}"; do
    echo "Repairing $keyspace..."
    kubectl exec -it scylla-0 -n tchat -- nodetool repair $keyspace
done
```

### Redis Maintenance

#### Backup

```bash
#!/bin/bash
# redis-backup.sh

BACKUP_DIR="/backups/redis"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

echo "Creating Redis backup..."
kubectl exec -it redis-0 -n tchat -- redis-cli --rdb - > $BACKUP_DIR/redis_${DATE}.rdb

# Clean up old backups
find $BACKUP_DIR -name "*.rdb" -mtime +7 -delete
```

## Deployment Procedures

### Blue-Green Deployment

```bash
#!/bin/bash
# blue-green-deploy.sh

SERVICE=$1
NEW_IMAGE=$2

if [ -z "$SERVICE" ] || [ -z "$NEW_IMAGE" ]; then
    echo "Usage: $0 <service> <image>"
    exit 1
fi

echo "Starting blue-green deployment for $SERVICE"

# Create green deployment
kubectl patch deployment $SERVICE -n tchat --patch "{
    \"metadata\": {\"labels\": {\"version\": \"green\"}},
    \"spec\": {
        \"template\": {
            \"metadata\": {\"labels\": {\"version\": \"green\"}},
            \"spec\": {\"containers\": [{\"name\": \"$SERVICE\", \"image\": \"$NEW_IMAGE\"}]}
        }
    }
}"

# Wait for green deployment to be ready
kubectl rollout status deployment/$SERVICE -n tchat --timeout=300s

# Run health checks
./scripts/health-check.sh

if [ $? -eq 0 ]; then
    echo "✅ Green deployment successful"

    # Update service selector to point to green
    kubectl patch service $SERVICE -n tchat --patch '{"spec": {"selector": {"version": "green"}}}'

    echo "✅ Traffic switched to green deployment"
else
    echo "❌ Health checks failed, rolling back"
    kubectl rollout undo deployment/$SERVICE -n tchat
    exit 1
fi
```

### Canary Deployment

```bash
#!/bin/bash
# canary-deploy.sh

SERVICE=$1
NEW_IMAGE=$2
CANARY_PERCENTAGE=${3:-10}

echo "Starting canary deployment for $SERVICE (${CANARY_PERCENTAGE}% traffic)"

# Scale up current deployment
CURRENT_REPLICAS=$(kubectl get deployment $SERVICE -n tchat -o jsonpath='{.spec.replicas}')
CANARY_REPLICAS=$((CURRENT_REPLICAS * CANARY_PERCENTAGE / 100))

# Create canary deployment
kubectl create deployment ${SERVICE}-canary -n tchat --image=$NEW_IMAGE --replicas=$CANARY_REPLICAS

# Add service selector for canary
kubectl patch service $SERVICE -n tchat --patch "{
    \"spec\": {
        \"selector\": {
            \"app\": \"$SERVICE\"
        }
    }
}"

# Monitor canary metrics
echo "Monitoring canary deployment for 10 minutes..."
sleep 600

# Check error rates
ERROR_RATE=$(curl -s "http://prometheus:9090/api/v1/query?query=rate(http_requests_total{status=~\"5..\",deployment=\"${SERVICE}-canary\"}[5m])" | jq -r '.data.result[0].value[1]')

if (( $(echo "$ERROR_RATE < 0.01" | bc -l) )); then
    echo "✅ Canary deployment successful, promoting to production"

    # Update main deployment
    kubectl set image deployment/$SERVICE $SERVICE=$NEW_IMAGE -n tchat

    # Delete canary
    kubectl delete deployment ${SERVICE}-canary -n tchat
else
    echo "❌ Canary deployment failed, rolling back"
    kubectl delete deployment ${SERVICE}-canary -n tchat
    exit 1
fi
```

## Backup and Recovery

### Automated Backup Script

```bash
#!/bin/bash
# full-backup.sh

BACKUP_ROOT="/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="$BACKUP_ROOT/$DATE"

mkdir -p $BACKUP_DIR

echo "Starting full system backup..."

# PostgreSQL backup
./scripts/postgres-backup.sh $BACKUP_DIR

# ScyllaDB backup
./scripts/scylla-backup.sh $BACKUP_DIR

# Redis backup
./scripts/redis-backup.sh $BACKUP_DIR

# Configuration backup
kubectl get configmaps -n tchat -o yaml > $BACKUP_DIR/configmaps.yaml
kubectl get secrets -n tchat -o yaml > $BACKUP_DIR/secrets.yaml

# Create backup manifest
cat > $BACKUP_DIR/manifest.json << EOF
{
    "timestamp": "$(date -Iseconds)",
    "version": "$(git rev-parse HEAD)",
    "services": {
        "auth": "$(kubectl get deployment auth-service -n tchat -o jsonpath='{.spec.template.spec.containers[0].image}')",
        "messaging": "$(kubectl get deployment messaging-service -n tchat -o jsonpath='{.spec.template.spec.containers[0].image}')",
        "payment": "$(kubectl get deployment payment-service -n tchat -o jsonpath='{.spec.template.spec.containers[0].image}')",
        "notification": "$(kubectl get deployment notification-service -n tchat -o jsonpath='{.spec.template.spec.containers[0].image}')"
    }
}
EOF

# Upload to S3 or backup storage
aws s3 sync $BACKUP_DIR s3://tchat-backups/$DATE/

echo "✅ Full backup completed: $DATE"
```

### Disaster Recovery

```bash
#!/bin/bash
# disaster-recovery.sh

BACKUP_DATE=$1

if [ -z "$BACKUP_DATE" ]; then
    echo "Usage: $0 <backup_date>"
    echo "Available backups:"
    aws s3 ls s3://tchat-backups/ | grep "PRE"
    exit 1
fi

echo "Starting disaster recovery from backup: $BACKUP_DATE"

# Download backup
aws s3 sync s3://tchat-backups/$BACKUP_DATE/ /tmp/recovery/

# Restore databases
./scripts/restore-postgres.sh /tmp/recovery/
./scripts/restore-scylla.sh /tmp/recovery/
./scripts/restore-redis.sh /tmp/recovery/

# Restore configurations
kubectl apply -f /tmp/recovery/configmaps.yaml
kubectl apply -f /tmp/recovery/secrets.yaml

# Deploy services with backed up versions
MANIFEST=$(cat /tmp/recovery/manifest.json)

for service in auth messaging payment notification; do
    IMAGE=$(echo $MANIFEST | jq -r ".services.$service")
    kubectl set image deployment/${service}-service ${service}-service=$IMAGE -n tchat
done

# Wait for all deployments
kubectl rollout status deployment -n tchat --timeout=600s

echo "✅ Disaster recovery completed"
```

## Performance Tuning

### Database Performance

#### PostgreSQL Tuning

```sql
-- Connection tuning
ALTER SYSTEM SET max_connections = 200;
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';

-- Query tuning
ALTER SYSTEM SET random_page_cost = 1.1;
ALTER SYSTEM SET seq_page_cost = 1.0;

-- WAL tuning
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;

SELECT pg_reload_conf();
```

#### ScyllaDB Tuning

```bash
# Compaction strategy
kubectl exec -it scylla-0 -n tchat -- cqlsh -e "
ALTER TABLE tchat_messaging.messages
WITH compaction = {
    'class': 'SizeTieredCompactionStrategy',
    'bucket_high': 2.0,
    'bucket_low': 0.5
};"

# Read/write timeout
kubectl exec -it scylla-0 -n tchat -- cqlsh -e "
ALTER KEYSPACE tchat_messaging
WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': 3
} AND durable_writes = true;"
```

### Service Performance

#### Go Runtime Tuning

```bash
# Set environment variables for optimal performance
kubectl patch deployment auth-service -n tchat --patch '{
    "spec": {
        "template": {
            "spec": {
                "containers": [{
                    "name": "auth-service",
                    "env": [
                        {"name": "GOGC", "value": "100"},
                        {"name": "GOMAXPROCS", "value": "2"},
                        {"name": "GODEBUG", "value": "gctrace=1"}
                    ]
                }]
            }
        }
    }
}'
```

#### Connection Pool Tuning

```yaml
# configmap-performance.yml
apiVersion: v1
kind: ConfigMap
metadata:
  name: tchat-performance-config
  namespace: tchat
data:
  # Database connections
  db_max_open_conns: "50"
  db_max_idle_conns: "10"
  db_conn_max_lifetime: "5m"

  # Redis connections
  redis_pool_size: "50"
  redis_min_idle_conns: "10"
  redis_dial_timeout: "5s"
  redis_read_timeout: "3s"
  redis_write_timeout: "3s"

  # HTTP client
  http_client_timeout: "30s"
  http_client_max_idle_conns: "100"
  http_client_max_conns_per_host: "10"
```

## Security Procedures

### Certificate Management

```bash
#!/bin/bash
# cert-renewal.sh

# Let's Encrypt certificate renewal
certbot certonly --webroot -w /var/www/html -d api.tchat.sea

# Update Kubernetes secret
kubectl create secret tls tchat-tls \
    --cert=/etc/letsencrypt/live/api.tchat.sea/fullchain.pem \
    --key=/etc/letsencrypt/live/api.tchat.sea/privkey.pem \
    -n tchat --dry-run=client -o yaml | kubectl apply -f -

# Restart ingress controller
kubectl rollout restart deployment nginx-ingress-controller -n ingress-nginx
```

### Security Audit

```bash
#!/bin/bash
# security-audit.sh

echo "Running security audit..."

# Check for exposed secrets
kubectl get secrets -n tchat -o json | jq -r '.items[].data | keys[]' | sort | uniq

# Check RBAC permissions
kubectl auth can-i --list --as=system:serviceaccount:tchat:default -n tchat

# Check network policies
kubectl get networkpolicies -n tchat

# Check pod security policies
kubectl get podsecuritypolicies

# Check for privileged containers
kubectl get pods -n tchat -o jsonpath='{range .items[*]}{.metadata.name}{": "}{.spec.containers[*].securityContext.privileged}{"\n"}{end}'

echo "Security audit completed"
```

### Incident Response Checklist

#### Security Incident Response

1. **Immediate Actions** (0-15 minutes)
   - [ ] Identify the scope of the incident
   - [ ] Isolate affected systems
   - [ ] Preserve evidence
   - [ ] Notify security team

2. **Investigation** (15-60 minutes)
   - [ ] Analyze logs for suspicious activity
   - [ ] Check for data exfiltration
   - [ ] Identify attack vectors
   - [ ] Document findings

3. **Containment** (1-4 hours)
   - [ ] Block malicious IPs
   - [ ] Rotate compromised credentials
   - [ ] Patch vulnerabilities
   - [ ] Update security rules

4. **Recovery** (4-24 hours)
   - [ ] Restore from clean backups
   - [ ] Verify system integrity
   - [ ] Implement additional security measures
   - [ ] Monitor for reoccurrence

5. **Post-Incident** (1-7 days)
   - [ ] Conduct post-mortem
   - [ ] Update procedures
   - [ ] Implement lessons learned
   - [ ] Report to stakeholders

## Emergency Contacts

- **On-call Engineer**: +1-xxx-xxx-xxxx
- **Tech Lead**: +1-xxx-xxx-xxxx
- **SRE Team**: sre@tchat.sea
- **Security Team**: security@tchat.sea
- **Management**: management@tchat.sea

## Useful Commands Reference

```bash
# Kubernetes
kubectl get pods -n tchat
kubectl describe pod <pod-name> -n tchat
kubectl logs -f <pod-name> -n tchat
kubectl exec -it <pod-name> -n tchat -- /bin/bash

# Docker Compose
docker-compose ps
docker-compose logs -f <service>
docker-compose restart <service>

# Database
psql -h host -U user -d database
redis-cli -h host -p port
cqlsh host port

# Monitoring
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```