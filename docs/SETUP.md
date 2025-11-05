# FindMe Development Setup Guide

This guide will walk you through setting up the complete FindMe development environment on your local machine.

## Table of Contents

- [Prerequisites](#prerequisites)
- [System Requirements](#system-requirements)
- [Backend Setup](#backend-setup)
- [Mobile App Setup](#mobile-app-setup)
- [Database Setup](#database-setup)
- [Running the Application](#running-the-application)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Software

Before you begin, ensure you have the following installed:

| Software | Minimum Version | Download Link |
|----------|----------------|---------------|
| **Node.js** | 20.x LTS | [nodejs.org](https://nodejs.org/) |
| **npm** | 10.x | Comes with Node.js |
| **Go** | 1.21+ | [go.dev](https://go.dev/dl/) |
| **Docker** | 24.x | [docker.com](https://www.docker.com/products/docker-desktop/) |
| **Docker Compose** | 2.20+ | Comes with Docker Desktop |
| **Git** | 2.40+ | [git-scm.com](https://git-scm.com/) |

### Platform-Specific Requirements

#### macOS

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Install Homebrew (if not already installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Watchman (required for React Native)
brew install watchman

# Install CocoaPods (required for iOS)
sudo gem install cocoapods
```

#### Linux (Ubuntu/Debian)

```bash
# Update package lists
sudo apt update

# Install build essentials
sudo apt install -y build-essential git curl

# Install Watchman
cd /tmp
git clone https://github.com/facebook/watchman.git
cd watchman
git checkout v2023.09.04.00
./autogen.sh
./configure
make
sudo make install
```

#### Windows

```powershell
# Install Chocolatey (package manager)
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Install dependencies
choco install -y nodejs.install python2 jdk8 git
```

### Mobile Development Setup

#### iOS Development (macOS Only)

1. **Install Xcode** (from App Store)
   - Version 15.0 or later
   - Open Xcode and agree to license terms
   - Install iOS Simulator

2. **Configure Xcode Command Line Tools**
   ```bash
   sudo xcode-select --switch /Applications/Xcode.app/Contents/Developer
   ```

#### Android Development (All Platforms)

1. **Install Android Studio**
   - Download from [developer.android.com](https://developer.android.com/studio)
   - Install with default settings

2. **Configure Android SDK**
   ```bash
   # Add to ~/.bashrc, ~/.zshrc, or equivalent
   export ANDROID_HOME=$HOME/Library/Android/sdk  # macOS
   # export ANDROID_HOME=$HOME/Android/Sdk       # Linux
   # export ANDROID_HOME=%LOCALAPPDATA%\Android\Sdk  # Windows

   export PATH=$PATH:$ANDROID_HOME/emulator
   export PATH=$PATH:$ANDROID_HOME/platform-tools
   export PATH=$PATH:$ANDROID_HOME/tools
   export PATH=$PATH:$ANDROID_HOME/tools/bin
   ```

3. **Install SDK Components**
   ```bash
   # Open Android Studio
   # Go to: Preferences > Appearance & Behavior > System Settings > Android SDK
   # Install:
   # - Android SDK Platform 34 (Android 14)
   # - Android SDK Build-Tools 34.0.0
   # - Android SDK Platform-Tools
   # - Android Emulator
   # - Intel x86 Emulator Accelerator (HAXM)
   ```

4. **Create Android Emulator**
   ```bash
   # List available system images
   sdkmanager --list | grep system-images

   # Download system image
   sdkmanager "system-images;android-34;google_apis;x86_64"

   # Create AVD
   avdmanager create avd -n Pixel_5_API_34 -k "system-images;android-34;google_apis;x86_64" -d "pixel_5"
   ```

## System Requirements

### Minimum Hardware

- **CPU**: Quad-core processor
- **RAM**: 16GB (8GB minimum)
- **Storage**: 20GB free space
- **Network**: Broadband internet connection

### Recommended Hardware

- **CPU**: 8-core processor or better
- **RAM**: 32GB
- **Storage**: 50GB+ SSD
- **Network**: High-speed internet

## Backend Setup

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/findme.git
cd findme
```

### 2. Install Go Dependencies

```bash
cd backend

# Initialize Go modules
go mod download

# Verify installation
go version
```

### 3. Set Up Environment Variables

```bash
# Copy environment template
cp ../.env.sample .env

# Edit .env with your configuration
# See .env.sample for detailed descriptions
nano .env  # or use your preferred editor
```

### 4. Start Docker Services

```bash
# Start PostgreSQL, Qdrant, and Redis
docker-compose up -d

# Verify services are running
docker-compose ps

# View logs
docker-compose logs -f
```

**docker-compose.yml** (create if not exists):
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: findme-postgres
    environment:
      POSTGRES_USER: findme
      POSTGRES_PASSWORD: findme_dev_password
      POSTGRES_DB: findme_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - findme-network

  qdrant:
    image: qdrant/qdrant:v1.7.0
    container_name: findme-qdrant
    ports:
      - "6333:6333"
      - "6334:6334"
    volumes:
      - qdrant_data:/qdrant/storage
    networks:
      - findme-network

  redis:
    image: redis:7-alpine
    container_name: findme-redis
    command: redis-server --appendonly yes
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - findme-network

networks:
  findme-network:
    driver: bridge

volumes:
  postgres_data:
  qdrant_data:
  redis_data:
```

### 5. Run Database Migrations

```bash
# Run PostgreSQL migrations
go run cmd/api/main.go migrate up

# Verify migrations
go run cmd/api/main.go migrate status
```

### 6. Seed Development Data (Optional)

```bash
# Seed database with test data
go run cmd/api/main.go seed

# Or use custom seed script
go run scripts/seed/main.go
```

### 7. Start Backend Server

```bash
# Development mode with hot reload
air  # if you have air installed

# Or run directly
go run cmd/api/main.go

# Server should start on http://localhost:8080
```

### 8. Verify Backend

```bash
# Health check
curl http://localhost:8080/health

# Expected response:
# {"status":"ok","version":"0.1.0"}

# API documentation
open http://localhost:8080/swagger/index.html
```

## Mobile App Setup

### 1. Navigate to Mobile Directory

```bash
cd mobile  # from project root
```

### 2. Install Dependencies

```bash
# Install npm packages
npm install

# Or using Yarn
yarn install
```

### 3. Install iOS Pods (macOS Only)

```bash
cd ios
pod install
cd ..
```

### 4. Configure Environment

```bash
# Copy environment configuration
cp .env.sample .env

# Edit with your local backend URL
# API_URL=http://localhost:8080
nano .env
```

### 5. Start Metro Bundler

```bash
# Start Metro
npm start

# Or
npx react-native start

# Clear cache if needed
npm start -- --reset-cache
```

### 6. Run on iOS Simulator (macOS Only)

```bash
# In a new terminal window
npm run ios

# Or specify simulator
npm run ios -- --simulator="iPhone 15 Pro"

# List available simulators
xcrun simctl list devices
```

### 7. Run on Android Emulator

```bash
# Start Android emulator first
emulator -avd Pixel_5_API_34 &

# Wait for emulator to boot, then:
npm run android

# Or use Android Studio to start emulator
```

## Database Setup

### PostgreSQL Setup

#### Manual Installation (Alternative to Docker)

```bash
# macOS
brew install postgresql@15
brew services start postgresql@15

# Linux
sudo apt install postgresql-15
sudo systemctl start postgresql

# Create database and user
psql postgres
CREATE DATABASE findme_db;
CREATE USER findme WITH ENCRYPTED PASSWORD 'findme_dev_password';
GRANT ALL PRIVILEGES ON DATABASE findme_db TO findme;
\q
```

#### Connection String

```bash
# Format
postgresql://username:password@localhost:5432/database_name

# Example
postgresql://findme:findme_dev_password@localhost:5432/findme_db
```

### Qdrant Setup

#### Using Docker (Recommended)

Already covered in Backend Setup step 4.

#### Manual Installation

```bash
# Download Qdrant
curl -L https://github.com/qdrant/qdrant/releases/download/v1.7.0/qdrant-x86_64-unknown-linux-gnu.tar.gz | tar xz

# Run Qdrant
./qdrant
```

#### Initialize Collections

```bash
# Run migration script
go run migrations/qdrant/init.go

# Or use HTTP API
curl -X PUT 'http://localhost:6333/collections/profile_embeddings' \
  -H 'Content-Type: application/json' \
  -d '{
    "vectors": {
      "size": 512,
      "distance": "Cosine"
    }
  }'
```

### Redis Setup

#### Using Docker (Recommended)

Already covered in Backend Setup step 4.

#### Manual Installation

```bash
# macOS
brew install redis
brew services start redis

# Linux
sudo apt install redis-server
sudo systemctl start redis

# Test connection
redis-cli ping
# Expected: PONG
```

## Running the Application

### Full Stack Development

Open 3 terminal windows:

**Terminal 1: Backend**
```bash
cd backend
docker-compose up -d  # Start services
go run cmd/api/main.go
```

**Terminal 2: Metro Bundler**
```bash
cd mobile
npm start
```

**Terminal 3: Mobile App**
```bash
cd mobile
npm run ios     # for iOS
# OR
npm run android # for Android
```

### Useful Commands

```bash
# Backend
cd backend
go test ./...                    # Run tests
go run cmd/api/main.go --help   # View CLI options
go build -o bin/api cmd/api/main.go  # Build binary

# Mobile
cd mobile
npm test                         # Run tests
npm run lint                     # Lint code
npm run type-check              # TypeScript check
npm run ios -- --configuration Release  # Release build
```

## Troubleshooting

### Common Issues

#### Port Already in Use

```bash
# Find process using port
lsof -i :8080  # backend
lsof -i :8081  # metro bundler

# Kill process
kill -9 <PID>
```

#### Metro Bundler Issues

```bash
# Clear cache and restart
cd mobile
rm -rf node_modules
npm install
npm start -- --reset-cache
```

#### iOS Build Fails

```bash
# Clean build
cd mobile/ios
xcodebuild clean
pod deintegrate
pod install
cd ..
npm run ios
```

#### Android Build Fails

```bash
# Clean Gradle cache
cd mobile/android
./gradlew clean
cd ..

# Clear Android build
rm -rf android/app/build
npm run android
```

#### Docker Services Not Starting

```bash
# Stop all containers
docker-compose down

# Remove volumes
docker-compose down -v

# Rebuild and start
docker-compose up -d --build

# View logs
docker-compose logs -f
```

#### Database Connection Errors

```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Check connection
psql postgresql://findme:findme_dev_password@localhost:5432/findme_db -c "SELECT 1"

# Reset database
docker-compose down postgres
docker volume rm findme_postgres_data
docker-compose up -d postgres
go run cmd/api/main.go migrate up
```

#### Qdrant Connection Issues

```bash
# Check Qdrant health
curl http://localhost:6333/health

# View Qdrant logs
docker-compose logs qdrant

# Restart Qdrant
docker-compose restart qdrant
```

### React Native Specific Issues

#### Module Not Found

```bash
cd mobile
rm -rf node_modules
rm package-lock.json
npm install
```

#### iOS Simulator Not Opening

```bash
# List simulators
xcrun simctl list devices

# Boot simulator manually
xcrun simctl boot "iPhone 15 Pro"

# Reset simulator
xcrun simctl erase all
```

#### Android Emulator Slow

```bash
# Enable hardware acceleration (macOS/Linux)
# Add to ~/.bash_profile or ~/.zshrc:
export ANDROID_EMULATOR_USE_SYSTEM_LIBS=1

# Or use a physical device with USB debugging enabled
adb devices
npm run android
```

### Getting Help

If you encounter issues not covered here:

1. Check the [GitHub Issues](https://github.com/yourusername/findme/issues)
2. Search [Discussions](https://github.com/yourusername/findme/discussions)
3. Join our community chat
4. Email: support@findme.ai

## Next Steps

Once your development environment is set up:

1. Read the [Architecture Documentation](../ARCHITECTURE.md)
2. Review the [API Documentation](API.md)
3. Check out the [Contributing Guidelines](../CONTRIBUTING.md)
4. Start building! 🚀

---

**Happy coding!** If you found this guide helpful, please give us a ⭐ on GitHub!
