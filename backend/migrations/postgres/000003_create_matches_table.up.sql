-- Create matches table
CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    matched_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    week_number INTEGER NOT NULL CHECK (week_number >= 1 AND week_number <= 53),
    year INTEGER NOT NULL CHECK (year >= 2025),
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'rejected', 'completed', 'expired')),
    match_score FLOAT NOT NULL CHECK (match_score >= 0 AND match_score <= 1),
    user_action VARCHAR(50) CHECK (user_action IN ('accepted', 'rejected')),
    matched_user_action VARCHAR(50) CHECK (matched_user_action IN ('accepted', 'rejected')),
    mutual_match BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    matched_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT unique_match_per_week UNIQUE(user_id, matched_user_id, week_number, year),
    CONSTRAINT no_self_match CHECK (user_id != matched_user_id)
);

-- Create indexes
CREATE INDEX idx_matches_user_week ON matches(user_id, week_number, year) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_matched_user_week ON matches(matched_user_id, week_number, year) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_status ON matches(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_mutual ON matches(mutual_match) WHERE mutual_match = TRUE AND deleted_at IS NULL;
CREATE INDEX idx_matches_expires_at ON matches(expires_at) WHERE status = 'pending';
CREATE INDEX idx_matches_score ON matches(match_score DESC);
CREATE INDEX idx_matches_created_at ON matches(created_at DESC);

-- Create trigger for auto-update
CREATE TRIGGER update_matches_updated_at
    BEFORE UPDATE ON matches
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to automatically set mutual_match flag
CREATE OR REPLACE FUNCTION update_mutual_match()
RETURNS TRIGGER AS $$
BEGIN
    -- If both users accepted, set mutual_match to true
    IF NEW.user_action = 'accepted' AND NEW.matched_user_action = 'accepted' THEN
        NEW.mutual_match := TRUE;
        NEW.matched_at := NOW();
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER check_mutual_match
    BEFORE UPDATE ON matches
    FOR EACH ROW
    EXECUTE FUNCTION update_mutual_match();

COMMENT ON TABLE matches IS 'Weekly AI-generated matches between users';
COMMENT ON COLUMN matches.match_score IS 'AI compatibility score (0-1)';
COMMENT ON COLUMN matches.mutual_match IS 'True when both users accept the match';
