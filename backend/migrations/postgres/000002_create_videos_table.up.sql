-- Create videos table
CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    storage_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    duration INTEGER NOT NULL CHECK (duration > 0 AND duration <= 120),
    file_size BIGINT NOT NULL CHECK (file_size > 0),
    mime_type VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'processing' CHECK (status IN ('uploading', 'processing', 'verifying', 'verified', 'rejected', 'error')),
    verification_score FLOAT CHECK (verification_score >= 0 AND verification_score <= 1),
    rejection_reason TEXT,
    metadata JSONB DEFAULT '{}',
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes
CREATE INDEX idx_videos_user_id ON videos(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_videos_status ON videos(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_videos_created_at ON videos(created_at DESC);
CREATE INDEX idx_videos_metadata ON videos USING gin(metadata);
CREATE INDEX idx_videos_verification_score ON videos(verification_score) WHERE status = 'verified';

-- Only one active video per user constraint
CREATE UNIQUE INDEX idx_videos_one_per_user ON videos(user_id) 
    WHERE deleted_at IS NULL AND status = 'verified';

-- Create trigger for auto-update
CREATE TRIGGER update_videos_updated_at
    BEFORE UPDATE ON videos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add foreign key to users table
ALTER TABLE users ADD CONSTRAINT fk_users_video_id 
    FOREIGN KEY (video_id) REFERENCES videos(id) ON DELETE SET NULL;

COMMENT ON TABLE videos IS 'Video profiles for user verification and matching';
COMMENT ON COLUMN videos.verification_score IS 'AI verification score (0-1)';
COMMENT ON COLUMN videos.metadata IS 'Additional video metadata (resolution, codec, etc)';
