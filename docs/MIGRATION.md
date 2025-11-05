# FindMe Database Migration Guide

Complete guide for managing database schema changes and migrations for PostgreSQL and Qdrant.

## Overview

FindMe uses two primary databases:
- **PostgreSQL**: Relational data (users, matches, video metadata)
- **Qdrant**: Vector database (embedding search, AI matching)

## Table of Contents

- [PostgreSQL Migrations](#postgresql-migrations)
- [Qdrant Collections](#qdrant-collections)
- [Migration Scripts](#migration-scripts)
- [Seeding Data](#seeding-data)
- [Backup & Restore](#backup--restore)
- [Version Management](#version-management)
- [Privacy & Compliance](#privacy--compliance)

---

## PostgreSQL Migrations

### Migration Tool

We use [golang-migrate](https://github.com/golang-migrate/migrate) for PostgreSQL schema migrations.

**Installation:**
```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Go install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Directory Structure

```
backend/
├── migrations/
│   ├── postgres/
│   │   ├── 000001_create_users_table.up.sql
│   │   ├── 000001_create_users_table.down.sql
│   │   ├── 000002_create_videos_table.up.sql
│   │   ├── 000002_create_videos_table.down.sql
│   │   ├── 000003_create_matches_table.up.sql
│   │   ├── 000003_create_matches_table.down.sql
│   │   └── ...
│   └── qdrant/
│       └── init.go
```

### Creating Migrations

```bash
# Create new migration
migrate create -ext sql -dir backend/migrations/postgres -seq create_users_table

# This creates:
# - 000001_create_users_table.up.sql
# - 000001_create_users_table.down.sql
```

### Example Migrations

**000001_create_users_table.up.sql:**
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(50) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    bio TEXT,
    video_id UUID,
    verified BOOLEAN DEFAULT FALSE,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_verified ON users(verified) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_active ON users(active) WHERE active = TRUE AND deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- Trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Full text search on bio
CREATE INDEX idx_users_bio_search ON users USING gin(to_tsvector('english', bio));
```

**000001_create_users_table.down.sql:**
```sql
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_users_bio_search;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_active;
DROP INDEX IF EXISTS idx_users_verified;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

**000002_create_videos_table.up.sql:**
```sql
CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    storage_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    duration INTEGER NOT NULL CHECK (duration > 0 AND duration <= 120),
    file_size BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'processing' CHECK (status IN ('uploading', 'processing', 'verifying', 'verified', 'rejected', 'error')),
    verification_score FLOAT CHECK (verification_score >= 0 AND verification_score <= 1),
    rejection_reason TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_videos_user_id ON videos(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_videos_status ON videos(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_videos_created_at ON videos(created_at DESC);
CREATE INDEX idx_videos_metadata ON videos USING gin(metadata);

CREATE TRIGGER update_videos_updated_at
    BEFORE UPDATE ON videos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

**000003_create_matches_table.up.sql:**
```sql
CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    matched_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    week_number INTEGER NOT NULL,
    year INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'rejected', 'completed', 'expired')),
    match_score FLOAT NOT NULL CHECK (match_score >= 0 AND match_score <= 1),
    user_action VARCHAR(50),
    matched_user_action VARCHAR(50),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(user_id, matched_user_id, week_number, year),
    CHECK (user_id != matched_user_id)
);

CREATE INDEX idx_matches_user_week ON matches(user_id, week_number, year) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_status ON matches(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_expires_at ON matches(expires_at) WHERE status = 'pending';
CREATE INDEX idx_matches_score ON matches(match_score DESC);

CREATE TRIGGER update_matches_updated_at
    BEFORE UPDATE ON matches
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

**000004_create_video_calls_table.up.sql:**
```sql
CREATE TABLE video_calls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID REFERENCES matches(id) ON DELETE CASCADE,
    caller_id UUID REFERENCES users(id) ON DELETE CASCADE,
    callee_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_id VARCHAR(255) UNIQUE NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    duration INTEGER,
    status VARCHAR(50) NOT NULL DEFAULT 'initiated' CHECK (status IN ('initiated', 'ringing', 'active', 'ended', 'missed', 'rejected')),
    recording_url VARCHAR(500),
    quality_rating INTEGER CHECK (quality_rating >= 1 AND quality_rating <= 5),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_video_calls_match ON video_calls(match_id);
CREATE INDEX idx_video_calls_participants ON video_calls(caller_id, callee_id);
CREATE INDEX idx_video_calls_session ON video_calls(session_id);
CREATE INDEX idx_video_calls_status ON video_calls(status);
CREATE INDEX idx_video_calls_created_at ON video_calls(created_at DESC);
```

### Running Migrations

```bash
# Set database URL
export DATABASE_URL="postgresql://findme:password@localhost:5432/findme_db?sslmode=disable"

# Apply all up migrations
migrate -path backend/migrations/postgres -database $DATABASE_URL up

# Apply specific number of migrations
migrate -path backend/migrations/postgres -database $DATABASE_URL up 2

# Rollback last migration
migrate -path backend/migrations/postgres -database $DATABASE_URL down 1

# Rollback all migrations
migrate -path backend/migrations/postgres -database $DATABASE_URL down

# Check migration version
migrate -path backend/migrations/postgres -database $DATABASE_URL version

# Force version (use with caution)
migrate -path backend/migrations/postgres -database $DATABASE_URL force 3
```

### In Application Code

```go
package main

import (
    "database/sql"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB) error {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return err
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations/postgres",
        "postgres",
        driver,
    )
    if err != nil {
        return err
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }

    return nil
}
```

---

## Qdrant Collections

### Profile Embeddings Collection

**Schema:**
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
    "age": "integer",
    "age_group": "keyword",
    "location": "geo",
    "interests": "keyword[]",
    "active": "bool",
    "verified": "bool",
    "created_at": "datetime",
    "video_features": "float[]"
  },
  "optimizers_config": {
    "indexing_threshold": 10000
  },
  "hnsw_config": {
    "m": 16,
    "ef_construct": 100
  }
}
```

### Creating Qdrant Collections

**Go Migration Script (backend/migrations/qdrant/init.go):**
```go
package main

import (
    "context"
    "log"
    
    "github.com/qdrant/go-client/qdrant"
)

func InitQdrantCollections(client *qdrant.Client) error {
    ctx := context.Background()

    // Create profile embeddings collection
    err := client.CreateCollection(ctx, &qdrant.CreateCollection{
        CollectionName: "profile_embeddings",
        VectorsConfig: &qdrant.VectorsConfig{
            Params: &qdrant.VectorParams{
                Size:     512,
                Distance: qdrant.Distance_Cosine,
            },
        },
        OptimizersConfig: &qdrant.OptimizersConfigDiff{
            IndexingThreshold: qdrant.PtrOf(uint64(10000)),
        },
        HnswConfig: &qdrant.HnswConfigDiff{
            M:              qdrant.PtrOf(uint64(16)),
            EfConstruct:    qdrant.PtrOf(uint64(100)),
            FullScanThreshold: qdrant.PtrOf(uint64(10000)),
        },
    })
    
    if err != nil {
        return err
    }

    // Create payload index for filtering
    err = client.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
        CollectionName: "profile_embeddings",
        FieldName:      "gender",
        FieldType:      qdrant.FieldType_FieldTypeKeyword,
    })
    
    if err != nil {
        return err
    }

    // Additional indexes
    indexes := []struct {
        field string
        ftype qdrant.FieldType
    }{
        {"age", qdrant.FieldType_FieldTypeInteger},
        {"age_group", qdrant.FieldType_FieldTypeKeyword},
        {"active", qdrant.FieldType_FieldTypeBool},
        {"verified", qdrant.FieldType_FieldTypeBool},
        {"location", qdrant.FieldType_FieldTypeGeo},
    }

    for _, idx := range indexes {
        err = client.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
            CollectionName: "profile_embeddings",
            FieldName:      idx.field,
            FieldType:      idx.ftype,
        })
        if err != nil {
            log.Printf("Warning: Failed to create index for %s: %v", idx.field, err)
        }
    }

    return nil
}

func main() {
    client, err := qdrant.NewClient(&qdrant.Config{
        Host: "localhost",
        Port: 6333,
    })
    if err != nil {
        log.Fatal(err)
    }

    if err := InitQdrantCollections(client); err != nil {
        log.Fatal(err)
    }

    log.Println("Qdrant collections initialized successfully")
}
```

### Inserting Vectors

```go
func InsertProfileEmbedding(client *qdrant.Client, userID string, embedding []float32, payload map[string]interface{}) error {
    ctx := context.Background()

    _, err := client.Upsert(ctx, &qdrant.UpsertPoints{
        CollectionName: "profile_embeddings",
        Points: []*qdrant.PointStruct{
            {
                Id:      &qdrant.PointId{PointIdOptions: &qdrant.PointId_Uuid{Uuid: userID}},
                Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: embedding}}},
                Payload: payload,
            },
        },
    })

    return err
}
```

### Searching Vectors

```go
func SearchSimilarProfiles(client *qdrant.Client, queryVector []float32, filters map[string]interface{}, limit uint64) ([]*qdrant.ScoredPoint, error) {
    ctx := context.Background()

    // Build filter conditions
    filter := &qdrant.Filter{
        Must: []*qdrant.Condition{
            {
                ConditionOneOf: &qdrant.Condition_Field{
                    Field: &qdrant.FieldCondition{
                        Key: "active",
                        Match: &qdrant.Match{MatchValue: &qdrant.Match_Boolean{Boolean: true}},
                    },
                },
            },
            {
                ConditionOneOf: &qdrant.Condition_Field{
                    Field: &qdrant.FieldCondition{
                        Key: "verified",
                        Match: &qdrant.Match{MatchValue: &qdrant.Match_Boolean{Boolean: true}},
                    },
                },
            },
        },
    }

    // Add gender filter if specified
    if gender, ok := filters["gender"].(string); ok {
        filter.Must = append(filter.Must, &qdrant.Condition{
            ConditionOneOf: &qdrant.Condition_Field{
                Field: &qdrant.FieldCondition{
                    Key: "gender",
                    Match: &qdrant.Match{MatchValue: &qdrant.Match_Keyword{Keyword: gender}},
                },
            },
        })
    }

    result, err := client.Search(ctx, &qdrant.SearchPoints{
        CollectionName: "profile_embeddings",
        Vector:         queryVector,
        Filter:         filter,
        Limit:          limit,
        WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
    })

    if err != nil {
        return nil, err
    }

    return result, nil
}
```

---

## Seeding Data

### Development Seed Data

**backend/scripts/seed/main.go:**
```go
package main

import (
    "context"
    "log"
    "time"
    
    "your-project/internal/repository"
    "your-project/internal/models"
)

func SeedDatabase() error {
    // Seed users
    users := []models.User{
        {
            Email:       "john@example.com",
            FullName:    "John Doe",
            DateOfBirth: time.Date(1995, 6, 15, 0, 0, 0, 0, time.UTC),
            Gender:      "male",
            Bio:         "Adventure seeker and coffee lover",
            Verified:    true,
            Active:      true,
        },
        {
            Email:       "jane@example.com",
            FullName:    "Jane Smith",
            DateOfBirth: time.Date(1993, 3, 20, 0, 0, 0, 0, time.UTC),
            Gender:      "female",
            Bio:         "Book worm and travel enthusiast",
            Verified:    true,
            Active:      true,
        },
    }

    for _, user := range users {
        if err := repo.CreateUser(context.Background(), &user); err != nil {
            return err
        }
        log.Printf("Created user: %s", user.Email)
    }

    return nil
}

func main() {
    if err := SeedDatabase(); err != nil {
        log.Fatal(err)
    }
    log.Println("Database seeded successfully")
}
```

---

## Backup & Restore

### PostgreSQL Backup

```bash
# Full backup
pg_dump -h localhost -U findme findme_db > backup_$(date +%Y%m%d_%H%M%S).sql

# Compressed backup
pg_dump -h localhost -U findme findme_db | gzip > backup_$(date +%Y%m%d_%H%M%S).sql.gz

# Backup specific tables
pg_dump -h localhost -U findme -t users -t videos findme_db > backup_users_videos.sql
```

### PostgreSQL Restore

```bash
# Restore from backup
psql -h localhost -U findme findme_db < backup.sql

# Restore from compressed backup
gunzip -c backup.sql.gz | psql -h localhost -U findme findme_db
```

### Qdrant Backup

```bash
# Create snapshot
curl -X POST 'http://localhost:6333/collections/profile_embeddings/snapshots'

# Download snapshot
curl 'http://localhost:6333/collections/profile_embeddings/snapshots/snapshot-2025-01-15.snapshot' \
  -o snapshot.snapshot
```

### Qdrant Restore

```bash
# Upload snapshot
curl -X PUT 'http://localhost:6333/collections/profile_embeddings/snapshots/upload' \
  --form 'snapshot=@snapshot.snapshot'

# Recover from snapshot
curl -X PUT 'http://localhost:6333/collections/profile_embeddings' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "profile_embeddings",
    "snapshot_path": "/snapshots/profile_embeddings/snapshot-2025-01-15.snapshot"
  }'
```

---

## Version Management

### Tracking Schema Version

```sql
CREATE TABLE schema_version (
    version INTEGER PRIMARY KEY,
    description TEXT NOT NULL,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Migration History

```bash
# View migration history
migrate -path backend/migrations/postgres -database $DATABASE_URL version

# Or query directly
psql -d findme_db -c "SELECT * FROM schema_migrations ORDER BY version DESC LIMIT 10;"
```

---

## Privacy & Compliance

### GDPR Compliance

**Right to be Forgotten:**
```sql
-- Soft delete (recommended)
UPDATE users SET deleted_at = NOW() WHERE id = 'user_uuid';

-- Hard delete (complete removal)
DELETE FROM users WHERE id = 'user_uuid'; -- Cascades to related data
```

**Data Export:**
```go
func ExportUserData(ctx context.Context, userID string) (map[string]interface{}, error) {
    // Collect all user data
    userData := make(map[string]interface{})
    
    // User profile
    user, _ := repo.GetUser(ctx, userID)
    userData["profile"] = user
    
    // Videos
    videos, _ := repo.GetUserVideos(ctx, userID)
    userData["videos"] = videos
    
    // Matches
    matches, _ := repo.GetUserMatches(ctx, userID)
    userData["matches"] = matches
    
    // Video calls
    calls, _ := repo.GetUserCalls(ctx, userID)
    userData["calls"] = calls
    
    return userData, nil
}
```

### Data Anonymization

```sql
-- Anonymize user data (retain for analytics)
UPDATE users 
SET 
    email = 'deleted_' || id || '@findme.ai',
    full_name = 'Deleted User',
    bio = NULL,
    password_hash = '',
    deleted_at = NOW()
WHERE id = 'user_uuid';
```

---

For migration support: database@findme.ai
