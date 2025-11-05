# FindMe Architecture

This document provides a comprehensive overview of FindMe's technical architecture, design decisions, and system components.

## Table of Contents

- [System Overview](#system-overview)
- [Architecture Principles](#architecture-principles)
- [Component Architecture](#component-architecture)
- [Data Flow](#data-flow)
- [Technology Stack](#technology-stack)
- [Security Architecture](#security-architecture)
- [Scalability Considerations](#scalability-considerations)
- [Future Improvements](#future-improvements)

## System Overview

FindMe is built as a modern three-tier application with a mobile-first approach:

```
┌─────────────────────────────────────────────────────────┐
│                    Mobile Clients                       │
│         (iOS / Android - React Native)                  │
└──────────────────┬──────────────────────────────────────┘
                   │
                   │ HTTPS/REST API + WebSocket
                   │
┌──────────────────▼──────────────────────────────────────┐
│                  API Gateway / Backend                  │
│                  (Go + Gin Framework)                   │
│                                                         │
│  ┌─────────────┬──────────────┬──────────────────┐    │
│  │   Auth      │  Video       │    Matching      │    │
│  │  Service    │  Service     │    Service       │    │
│  └─────────────┴──────────────┴──────────────────┘    │
└──────────────────┬──────────────────────────────────────┘
                   │
        ┌──────────┼──────────────────┐
        │          │                  │
┌───────▼───┐ ┌───▼────────┐  ┌─────▼──────┐
│ PostgreSQL│ │   Qdrant   │  │   Redis    │
│   (User   │ │  (Vector   │  │  (Cache/   │
│   Data)   │ │    DB)     │  │  Session)  │
└───────────┘ └────────────┘  └────────────┘
                   │
            ┌──────▼──────┐
            │     S3      │
            │   (Video    │
            │   Storage)  │
            └─────────────┘
```

### Architecture Principles

1. **Mobile-First**: Optimize for mobile experience and constraints
2. **API-First**: Backend as a service for mobile clients
3. **Scalability**: Horizontal scaling for all components
4. **Security**: End-to-end encryption for sensitive data
5. **Performance**: Sub-second response times for critical operations
6. **Privacy**: GDPR/CCPA compliant by design
7. **Reliability**: 99.9% uptime target

## Component Architecture

### Mobile Application (React Native)

The mobile app follows a feature-based architecture with clear separation of concerns.

#### Directory Structure

```
mobile/
├── src/
│   ├── components/           # Reusable UI components
│   │   ├── common/          # Buttons, inputs, cards
│   │   ├── video/           # Video player, recorder
│   │   └── profile/         # Profile-related components
│   ├── screens/             # Screen-level components
│   │   ├── auth/           # Login, signup, verification
│   │   ├── profile/        # Profile creation, editing
│   │   ├── matches/        # Match list, weekly picks
│   │   └── video-call/     # Video call interface
│   ├── navigation/          # Navigation configuration
│   ├── hooks/              # Custom React hooks
│   ├── services/           # API clients and business logic
│   │   ├── api/           # REST API client
│   │   ├── auth/          # Authentication service
│   │   ├── video/         # Video upload/processing
│   │   └── webrtc/        # WebRTC for video calls
│   ├── store/             # Redux state management
│   │   ├── slices/       # Redux Toolkit slices
│   │   └── middleware/   # Custom middleware
│   ├── types/            # TypeScript definitions
│   ├── utils/            # Utility functions
│   ├── constants/        # App constants
│   └── assets/           # Images, fonts, etc.
├── ios/                  # iOS native code
├── android/              # Android native code
└── __tests__/           # Test files
```

#### Key Mobile Components

**Authentication Flow**
- JWT-based authentication
- Biometric authentication support
- Secure token storage in device keychain
- Automatic token refresh

**Video Recording**
- Native camera integration
- Real-time video preview
- Maximum 2-minute recording
- Client-side compression before upload
- Progress tracking during upload

**WebRTC Video Calls**
- Peer-to-peer video/audio
- STUN/TURN server configuration
- Connection quality indicators
- Call recording (with consent)
- Screen sharing capability

**State Management**
- Redux Toolkit for global state
- Local component state for UI
- Persistent storage for offline capability
- Optimistic UI updates

### Backend API (Go + Gin)

The backend follows clean architecture principles with clear layer separation.

#### Directory Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/                     # HTTP layer
│   │   ├── handlers/           # Request handlers
│   │   ├── middleware/         # Middleware (auth, logging)
│   │   └── routes/             # Route definitions
│   ├── domain/                 # Business domain
│   │   ├── models/            # Domain models
│   │   └── errors/            # Domain errors
│   ├── service/               # Business logic layer
│   │   ├── auth/             # Authentication service
│   │   ├── user/             # User management
│   │   ├── video/            # Video processing
│   │   ├── matching/         # AI matching engine
│   │   └── notification/     # Notifications
│   ├── repository/           # Data access layer
│   │   ├── postgres/        # PostgreSQL repositories
│   │   ├── qdrant/          # Qdrant repositories
│   │   └── redis/           # Redis cache
│   ├── infrastructure/       # External services
│   │   ├── storage/         # S3/storage client
│   │   ├── ai/              # AI service integration
│   │   └── email/           # Email service
│   └── config/              # Configuration management
├── pkg/                     # Public packages
│   ├── jwt/                # JWT utilities
│   ├── validator/          # Input validation
│   └── logger/             # Structured logging
├── migrations/             # Database migrations
│   ├── postgres/
│   └── qdrant/
├── scripts/               # Utility scripts
└── docs/                  # API documentation
```

#### Key Backend Services

**Authentication Service**
```go
type AuthService interface {
    Register(ctx context.Context, req RegisterRequest) (*User, error)
    Login(ctx context.Context, req LoginRequest) (*TokenPair, error)
    RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
    VerifyEmail(ctx context.Context, token string) error
    ResetPassword(ctx context.Context, req ResetPasswordRequest) error
}
```

**Video Service**
```go
type VideoService interface {
    UploadVideo(ctx context.Context, userID string, video io.Reader) (*Video, error)
    ProcessVideo(ctx context.Context, videoID string) error
    VerifyVideo(ctx context.Context, videoID string) (*VerificationResult, error)
    DeleteVideo(ctx context.Context, videoID string) error
    GetVideoURL(ctx context.Context, videoID string) (string, error)
}
```

**Matching Service**
```go
type MatchingService interface {
    GenerateWeeklyMatches(ctx context.Context, userID string) ([]*Match, error)
    GetUserEmbedding(ctx context.Context, userID string) ([]float32, error)
    FindSimilarProfiles(ctx context.Context, embedding []float32, limit int) ([]*Profile, error)
    RecordInteraction(ctx context.Context, interaction *Interaction) error
    GetMatchScore(ctx context.Context, userID1, userID2 string) (float64, error)
}
```

**Video Call Service**
```go
type VideoCallService interface {
    InitiateCall(ctx context.Context, callerID, calleeID string) (*CallSession, error)
    AcceptCall(ctx context.Context, sessionID string) error
    RejectCall(ctx context.Context, sessionID string) error
    EndCall(ctx context.Context, sessionID string) (*CallSummary, error)
    GetICEServers(ctx context.Context) ([]*ICEServer, error)
}
```

### Database Layer

#### PostgreSQL Schema

**Users Table**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(50) NOT NULL,
    bio TEXT,
    video_id UUID REFERENCES videos(id),
    verified BOOLEAN DEFAULT FALSE,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_verified ON users(verified);
```

**Videos Table**
```sql
CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    storage_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    duration INTEGER NOT NULL, -- seconds
    status VARCHAR(50) NOT NULL, -- processing, verified, rejected
    verification_score FLOAT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_videos_user_id ON videos(user_id);
CREATE INDEX idx_videos_status ON videos(status);
```

**Matches Table**
```sql
CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    matched_user_id UUID REFERENCES users(id),
    week_number INTEGER NOT NULL,
    year INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL, -- pending, accepted, rejected, completed
    match_score FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, matched_user_id, week_number, year)
);

CREATE INDEX idx_matches_user_week ON matches(user_id, week_number, year);
CREATE INDEX idx_matches_status ON matches(status);
```

**Video Calls Table**
```sql
CREATE TABLE video_calls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID REFERENCES matches(id),
    caller_id UUID REFERENCES users(id),
    callee_id UUID REFERENCES users(id),
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    duration INTEGER, -- seconds
    status VARCHAR(50) NOT NULL, -- initiated, ringing, active, ended, missed
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_video_calls_match ON video_calls(match_id);
CREATE INDEX idx_video_calls_participants ON video_calls(caller_id, callee_id);
```

#### Qdrant Collections

**Profile Embeddings Collection**
```json
{
  "name": "profile_embeddings",
  "vectors": {
    "size": 512,
    "distance": "Cosine"
  },
  "payload_schema": {
    "user_id": "keyword",
    "gender": "keyword",
    "age_group": "keyword",
    "interests": "keyword[]",
    "location": "geo",
    "active": "bool",
    "verified": "bool"
  }
}
```

**Search Filters**
- Gender preferences
- Age range
- Location proximity
- Verification status
- Activity status
- Previous matches (exclusion)

### AI Matching Engine

#### Embedding Generation

Video profiles are processed to generate multi-modal embeddings:

```
Video Input
    │
    ├─> Visual Features (CNN)
    │   ├─> Facial features
    │   ├─> Body language
    │   └─> Environment/style
    │
    ├─> Audio Features (RNN)
    │   ├─> Voice tone
    │   ├─> Speech patterns
    │   └─> Confidence level
    │
    └─> Text Features (Transformer)
        ├─> Transcribed speech
        ├─> Keyword extraction
        └─> Sentiment analysis
            │
            ▼
    Fusion Layer (Multi-Head Attention)
            │
            ▼
    512-dimensional embedding vector
```

#### Matching Algorithm

```go
func (s *MatchingService) GenerateWeeklyMatches(ctx context.Context, userID string) ([]*Match, error) {
    // 1. Get user embedding
    userEmbedding, err := s.getUserEmbedding(ctx, userID)
    
    // 2. Apply filters (preferences, exclusions)
    filters := s.buildFilters(userID)
    
    // 3. Query Qdrant for similar profiles
    candidates, err := s.qdrant.Search(ctx, SearchParams{
        Vector: userEmbedding,
        Filter: filters,
        Limit:  100, // Get more candidates than needed
    })
    
    // 4. Apply secondary scoring (interaction history, activity)
    scored := s.scoreAndRank(candidates, userID)
    
    // 5. Select top 4 diverse matches
    matches := s.selectDiverseMatches(scored, 4)
    
    // 6. Store matches for the week
    return s.repository.SaveWeeklyMatches(ctx, userID, matches)
}
```

## Data Flow

### User Registration & Profile Creation

```
1. User signs up (email/password)
   │
   ▼
2. Email verification sent
   │
   ▼
3. User verifies email
   │
   ▼
4. User records video profile
   │
   ▼
5. Video uploaded to S3
   │
   ▼
6. Video processing pipeline:
   ├─> Transcoding
   ├─> Thumbnail generation
   ├─> AI verification
   └─> Embedding generation
   │
   ▼
7. Embedding stored in Qdrant
   │
   ▼
8. Profile activated
```

### Weekly Matching Cycle

```
Every Monday 00:00 UTC:
│
├─> For each active user:
│   │
│   ├─> Generate 4 new matches
│   ├─> Send push notification
│   └─> Update user's weekly matches
│
▼
User views matches → Selects one → Mutual match?
                          │              │
                      No match        Yes (Match!)
                          │              │
                          └──────────────┼─> Enable video calling
                                         │
                                         ▼
                                  Video calls during week
                                         │
                                         ▼
                              Sunday: Decision time
                                    │
                           ┌────────┴────────┐
                           │                 │
                    Continue match    Delete app or
                      Next week      Find new matches
```

### Video Call Flow

```
User A initiates call
    │
    ▼
Create call session
    │
    ▼
Send push to User B → User B accepts
    │                      │
    │                      ▼
    │              Exchange ICE candidates
    │                      │
    │                      ▼
    └──────────────> Establish P2P connection
                           │
                           ▼
                    WebRTC video/audio stream
                           │
                           ▼
                    Call ends → Save duration
```

## Security Architecture

### Authentication & Authorization

- **JWT-based authentication** with access and refresh tokens
- **Access tokens**: 15-minute expiry, stored in memory
- **Refresh tokens**: 7-day expiry, stored in secure storage
- **Token rotation** on each refresh

### Data Encryption

- **In transit**: TLS 1.3 for all API communications
- **At rest**: AES-256 encryption for sensitive data
- **Video storage**: Server-side encryption (S3)
- **Database**: Encrypted PostgreSQL fields for PII

### Privacy Controls

- **Data minimization**: Collect only necessary information
- **User consent**: Explicit consent for video processing
- **Right to deletion**: Complete data removal on request
- **Data portability**: Export user data in standard formats

### Rate Limiting

```go
// Example rate limits
var RateLimits = map[string]RateLimit{
    "auth/register":     {Requests: 5, Window: time.Hour},
    "auth/login":        {Requests: 10, Window: time.Hour},
    "video/upload":      {Requests: 3, Window: time.Day},
    "matches/generate":  {Requests: 1, Window: time.Week},
    "calls/initiate":    {Requests: 50, Window: time.Day},
}
```

## Scalability Considerations

### Horizontal Scaling

**API Servers**
- Stateless design enables unlimited horizontal scaling
- Load balancer distributes traffic (round-robin or least connections)
- Auto-scaling based on CPU/memory metrics

**Database Scaling**
- **PostgreSQL**: Read replicas for read-heavy operations
- **Qdrant**: Sharding for large vector collections
- **Redis**: Cluster mode for cache distribution

### Caching Strategy

**Redis Cache Layers**
```
L1: User session data (TTL: 7 days)
L2: API responses (TTL: 5 minutes)
L3: Computed embeddings (TTL: 24 hours)
L4: Match recommendations (TTL: 1 week)
```

### CDN Strategy

- Video content served via CloudFront/CDN
- Edge caching for frequently accessed videos
- Geo-distribution for global low latency

### Background Jobs

Using a job queue system for async operations:

```go
type Job interface {
    VideoProcessing    // Transcode, thumbnail, embedding
    WeeklyMatchGeneration
    NotificationDelivery
    DataBackup
    AnalyticsAggregation
}
```

## Technology Stack

### Mobile (React Native)

| Category | Technology | Purpose |
|----------|-----------|---------|
| Framework | React Native 0.73 | Cross-platform mobile |
| Language | TypeScript 5.0 | Type-safe development |
| State | Redux Toolkit | Global state management |
| Navigation | React Navigation 6 | Screen navigation |
| Video | react-native-camera | Camera & video recording |
| Calls | react-native-webrtc | Video calling |
| Storage | @react-native-async-storage | Local data persistence |
| Auth | react-native-keychain | Secure credential storage |
| UI | React Native Elements | UI component library |

### Backend (Go)

| Category | Technology | Purpose |
|----------|-----------|---------|
| Language | Go 1.21 | High-performance backend |
| Framework | Gin | HTTP web framework |
| ORM | GORM | Database ORM |
| Database | PostgreSQL 15 | Relational data |
| Vector DB | Qdrant 1.7 | Vector similarity search |
| Cache | Redis 7 | In-memory caching |
| Queue | Redis Queue | Background jobs |
| Storage | AWS S3 | Video/file storage |
| Auth | golang-jwt | JWT implementation |
| Validation | go-playground/validator | Input validation |
| Testing | testify | Unit & integration tests |

### Infrastructure

| Category | Technology | Purpose |
|----------|-----------|---------|
| Container | Docker | Containerization |
| Orchestration | Kubernetes | Container orchestration |
| CI/CD | GitHub Actions | Automated deployments |
| Monitoring | Prometheus + Grafana | Metrics & dashboards |
| Logging | ELK Stack | Centralized logging |
| Tracing | Jaeger | Distributed tracing |
| CDN | CloudFront | Content delivery |

## Future Improvements

### Short Term (Q1-Q2 2025)

- [ ] Implement video quality adaptive streaming
- [ ] Add in-app messaging between matches
- [ ] Enhance AI matching with user feedback loop
- [ ] Add profile verification levels (basic, photo ID, video KYC)
- [ ] Implement safety features (report, block, safety center)

### Medium Term (Q3-Q4 2025)

- [ ] Multi-language support and localization
- [ ] Advanced matching preferences (lifestyle, values)
- [ ] Group video calls / virtual events
- [ ] Premium features and subscription model
- [ ] Dating coach / AI assistant

### Long Term (2026+)

- [ ] AR filters and effects (maintaining authenticity)
- [ ] Voice-only matching option
- [ ] Success stories and testimonials platform
- [ ] API for third-party integrations
- [ ] Open-source community plugins

---

For questions about the architecture, please open a [Discussion](https://github.com/yourusername/findme/discussions) or contact the team at architecture@findme.ai.
