#!/bin/bash
set -e

# FindMe Deployment Script
# Usage: ./scripts/deploy.sh [staging|production]

ENVIRONMENT=${1:-staging}

echo "🚀 Deploying FindMe to $ENVIRONMENT..."

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if environment is valid
if [ "$ENVIRONMENT" != "staging" ] && [ "$ENVIRONMENT" != "production" ]; then
    echo -e "${RED}❌ Invalid environment. Use 'staging' or 'production'${NC}"
    exit 1
fi

# Check if docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed${NC}"
    exit 1
fi

# Check if docker-compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ Docker Compose is not installed${NC}"
    exit 1
fi

echo -e "${YELLOW}📦 Building Docker images...${NC}"
docker-compose -f backend/docker-compose.yml build

echo -e "${YELLOW}🔄 Running database migrations...${NC}"
docker-compose -f backend/docker-compose.yml run --rm api \
    sh -c "cd /app && go run cmd/migrate/main.go up" || true

echo -e "${YELLOW}🌱 Seeding database (if needed)...${NC}"
if [ "$ENVIRONMENT" = "staging" ]; then
    docker-compose -f backend/docker-compose.yml run --rm api \
        sh -c "cd /app && go run scripts/seed/main.go" || true
fi

echo -e "${YELLOW}🔄 Restarting services...${NC}"
docker-compose -f backend/docker-compose.yml down
docker-compose -f backend/docker-compose.yml up -d

echo -e "${YELLOW}⏳ Waiting for services to be healthy...${NC}"
sleep 5

# Health check
echo -e "${YELLOW}🏥 Running health check...${NC}"
MAX_RETRIES=10
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Health check passed${NC}"
        break
    fi
    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo -e "${YELLOW}⏳ Waiting for service... ($RETRY_COUNT/$MAX_RETRIES)${NC}"
    sleep 3
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo -e "${RED}❌ Health check failed after $MAX_RETRIES attempts${NC}"
    docker-compose -f backend/docker-compose.yml logs --tail=50
    exit 1
fi

echo -e "${GREEN}✅ Deployment to $ENVIRONMENT completed successfully!${NC}"
echo -e "${GREEN}❤️  FindMe is now running${NC}"

# Show logs
echo -e "${YELLOW}📝 Showing recent logs...${NC}"
docker-compose -f backend/docker-compose.yml logs --tail=20
