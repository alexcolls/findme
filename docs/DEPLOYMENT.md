# FindMe Deployment Guide

Complete guide for deploying FindMe to production environments.

## Overview

This guide covers deploying FindMe's three main components:
1. **Mobile Apps** (iOS & Android)
2. **Backend API** (Go server)
3. **Databases** (PostgreSQL, Qdrant, Redis)

## Prerequisites

- GitHub repository with CI/CD access
- Cloud provider account (AWS/GCP/Azure)
- Apple Developer Account ($99/year) for iOS
- Google Play Developer Account ($25 one-time) for Android
- Domain name and SSL certificates

---

## Mobile App Deployment

### iOS App Store

#### 1. Prepare for Release

```bash
cd mobile/ios

# Update version and build number in Xcode
# Update app icons and splash screens
# Configure app signing with provisioning profiles
```

#### 2. Build Release

```bash
# Clean build folder
xcodebuild clean

# Archive for App Store
xcodebuild archive \
  -workspace FindMe.xcworkspace \
  -scheme FindMe \
  -configuration Release \
  -archivePath ./build/FindMe.xcarchive

# Export IPA
xcodebuild -exportArchive \
  -archivePath ./build/FindMe.xcarchive \
  -exportPath ./build \
  -exportOptionsPlist ExportOptions.plist
```

#### 3. Submit to App Store

```bash
# Upload with Transporter or Xcode
# Or use fastlane
fastlane ios release
```

#### 4. App Store Connect

1. Create app in App Store Connect
2. Fill metadata (description, screenshots, keywords)
3. Set age rating and pricing
4. Submit for review
5. Typical review time: 1-2 days

### Android Play Store

#### 1. Prepare for Release

```bash
cd mobile/android

# Generate signing key (first time only)
keytool -genkeypair -v \
  -keystore findme-release.keystore \
  -alias findme \
  -keyalg RSA \
  -keysize 2048 \
  -validity 10000
```

#### 2. Build Release APK/AAB

```bash
# Build release bundle (recommended)
./gradlew bundleRelease

# Or build APK
./gradlew assembleRelease

# Output: android/app/build/outputs/bundle/release/app-release.aab
```

#### 3. Upload to Play Console

```bash
# Manual upload or use fastlane
fastlane android deploy
```

#### 4. Google Play Console

1. Create app in Play Console
2. Complete store listing
3. Set content rating
4. Upload screenshots and videos
5. Publish to production/beta
6. Review time: Hours to 1 day

---

## Backend Deployment

### Docker Containerization

**Dockerfile:**
```dockerfile
# Multi-stage build
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api cmd/api/main.go

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/bin/api .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./api"]
```

**docker-compose.yml (Production):**
```yaml
version: '3.8'

services:
  api:
    image: findme/api:latest
    restart: always
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
    env_file:
      - .env.production
    depends_on:
      - postgres
      - qdrant
      - redis

  postgres:
    image: postgres:15-alpine
    restart: always
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      POSTGRES_PASSWORD: ${DB_PASSWORD}

  qdrant:
    image: qdrant/qdrant:v1.7.0
    restart: always
    volumes:
      - qdrant_data:/qdrant/storage

  redis:
    image: redis:7-alpine
    restart: always
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  qdrant_data:
  redis_data:
```

### Cloud Deployment Options

#### AWS Deployment

**Services Used:**
- **ECS/Fargate**: Container orchestration
- **RDS**: PostgreSQL database
- **ElastiCache**: Redis
- **EC2**: Qdrant (or self-managed)
- **S3**: Video storage
- **CloudFront**: CDN
- **Route 53**: DNS
- **ALB**: Load balancer
- **ACM**: SSL certificates

**Deploy with ECS:**
```bash
# Build and push to ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 123456789.dkr.ecr.us-east-1.amazonaws.com
docker build -t findme-api .
docker tag findme-api:latest 123456789.dkr.ecr.us-east-1.amazonaws.com/findme-api:latest
docker push 123456789.dkr.ecr.us-east-1.amazonaws.com/findme-api:latest

# Deploy to ECS
aws ecs update-service --cluster findme-cluster --service findme-api --force-new-deployment
```

#### Google Cloud Platform

**Services Used:**
- **Cloud Run**: Serverless containers
- **Cloud SQL**: PostgreSQL
- **Memorystore**: Redis
- **GCE**: Qdrant
- **Cloud Storage**: Videos
- **Cloud CDN**: Content delivery
- **Cloud Load Balancing**: Traffic distribution

**Deploy to Cloud Run:**
```bash
# Build and deploy
gcloud builds submit --tag gcr.io/findme-project/api
gcloud run deploy findme-api \
  --image gcr.io/findme-project/api \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

#### Kubernetes Deployment

**deployment.yaml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: findme-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: findme-api
  template:
    metadata:
      labels:
        app: findme-api
    spec:
      containers:
      - name: api
        image: findme/api:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: findme-secrets
              key: database-url
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

---

## Database Deployment

### PostgreSQL

**Production Setup:**
```bash
# Use managed service (RDS, Cloud SQL, etc.)
# Enable automated backups
# Set up read replicas for scaling
# Enable SSL connections
```

### Qdrant Cluster

**Production Config:**
```yaml
# qdrant-config.yaml
service:
  host: 0.0.0.0
  grpc_port: 6334

storage:
  storage_path: /qdrant/storage
  snapshots_path: /qdrant/snapshots
  on_disk_payload: true

cluster:
  enabled: true
  p2p:
    port: 6335
  consensus:
    tick_period_ms: 100
```

### Redis

**High Availability Setup:**
```bash
# Use managed service or Redis Sentinel/Cluster
# Enable persistence (AOF + RDS)
# Configure replication
# Set appropriate eviction policy
```

---

## CI/CD Pipeline

### GitHub Actions

**.github/workflows/deploy.yml:**
```yaml
name: Deploy

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: cd backend && go test ./...

  test-mobile:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '20'
      - run: cd mobile && npm ci && npm test

  deploy-backend:
    needs: [test-backend]
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build and push Docker image
        run: |
          docker build -t findme-api:${{ github.sha }} .
          docker push findme-api:${{ github.sha }}
      - name: Deploy to production
        run: |
          # Deploy to your cloud provider

  deploy-mobile:
    needs: [test-mobile]
    if: github.ref == 'refs/heads/main'
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build and deploy iOS
        run: fastlane ios beta
      - name: Build and deploy Android
        run: fastlane android beta
```

---

## SSL/TLS Configuration

### Let's Encrypt (Free)

```bash
# Install certbot
sudo apt install certbot

# Generate certificate
sudo certbot certonly --standalone -d api.findme.ai

# Auto-renewal
sudo certbot renew --dry-run
```

### Nginx Configuration

```nginx
server {
    listen 443 ssl http2;
    server_name api.findme.ai;

    ssl_certificate /etc/letsencrypt/live/api.findme.ai/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.findme.ai/privkey.pem;
    
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## Monitoring & Logging

### Prometheus + Grafana

**docker-compose.monitoring.yml:**
```yaml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

**prometheus.yml:**
```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'findme-api'
    static_configs:
      - targets: ['api:8080']
```

### Logging with ELK Stack

```bash
# Elasticsearch + Logstash + Kibana
docker-compose -f docker-compose.elk.yml up -d
```

---

## Health Checks & Readiness

**Kubernetes Probes:**
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

---

## Backup Strategy

### Database Backups

```bash
# Automated daily backups
0 2 * * * pg_dump findme_db > /backups/findme_$(date +\%Y\%m\%d).sql

# Qdrant snapshots
curl -X POST 'http://localhost:6333/collections/profile_embeddings/snapshots'
```

### Disaster Recovery

1. **RTO**: Recovery Time Objective - 1 hour
2. **RPO**: Recovery Point Objective - 1 hour
3. Regular backup testing
4. Multi-region redundancy
5. Automated failover

---

## Scaling Considerations

### Horizontal Scaling

- Load balancer with multiple API instances
- Database read replicas
- Redis cluster for cache distribution
- Qdrant sharding for large vector collections

### Auto-scaling

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: findme-api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: findme-api
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

---

## Security Checklist

- [ ] Enable HTTPS/TLS everywhere
- [ ] Rotate all secrets and API keys
- [ ] Enable database encryption at rest
- [ ] Configure firewalls and security groups
- [ ] Enable DDoS protection
- [ ] Set up WAF (Web Application Firewall)
- [ ] Regular security audits
- [ ] Vulnerability scanning
- [ ] Rate limiting in place
- [ ] CORS properly configured

---

## Cost Optimization

- Use spot instances for non-critical workloads
- Implement auto-scaling
- Use CDN for static content
- Optimize database queries
- Archive old data
- Monitor and alert on costs

---

For deployment support: devops@findme.ai
