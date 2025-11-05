-- Create video_calls table
CREATE TABLE video_calls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    caller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    callee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id VARCHAR(255) UNIQUE NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    duration INTEGER CHECK (duration >= 0),
    status VARCHAR(50) NOT NULL DEFAULT 'initiated' CHECK (status IN ('initiated', 'ringing', 'active', 'ended', 'missed', 'rejected', 'failed')),
    recording_url VARCHAR(500),
    quality_rating INTEGER CHECK (quality_rating >= 1 AND quality_rating <= 5),
    feedback TEXT,
    metadata JSONB DEFAULT '{}',
    ice_servers JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT no_self_call CHECK (caller_id != callee_id)
);

-- Create indexes
CREATE INDEX idx_video_calls_match ON video_calls(match_id);
CREATE INDEX idx_video_calls_caller ON video_calls(caller_id);
CREATE INDEX idx_video_calls_callee ON video_calls(callee_id);
CREATE INDEX idx_video_calls_participants ON video_calls(caller_id, callee_id);
CREATE INDEX idx_video_calls_session ON video_calls(session_id);
CREATE INDEX idx_video_calls_status ON video_calls(status);
CREATE INDEX idx_video_calls_created_at ON video_calls(created_at DESC);
CREATE INDEX idx_video_calls_started_at ON video_calls(started_at DESC) WHERE started_at IS NOT NULL;

-- Function to calculate duration on end
CREATE OR REPLACE FUNCTION calculate_call_duration()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'ended' AND NEW.started_at IS NOT NULL AND NEW.ended_at IS NOT NULL THEN
        NEW.duration := EXTRACT(EPOCH FROM (NEW.ended_at - NEW.started_at))::INTEGER;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER calculate_video_call_duration
    BEFORE UPDATE ON video_calls
    FOR EACH ROW
    WHEN (NEW.status = 'ended')
    EXECUTE FUNCTION calculate_call_duration();

COMMENT ON TABLE video_calls IS 'Video call sessions between matched users';
COMMENT ON COLUMN video_calls.duration IS 'Call duration in seconds';
COMMENT ON COLUMN video_calls.quality_rating IS 'User-provided call quality rating';
