-- Drop trigger and function
DROP TRIGGER IF EXISTS calculate_video_call_duration ON video_calls;
DROP FUNCTION IF EXISTS calculate_call_duration();

-- Drop indexes
DROP INDEX IF EXISTS idx_video_calls_started_at;
DROP INDEX IF EXISTS idx_video_calls_created_at;
DROP INDEX IF EXISTS idx_video_calls_status;
DROP INDEX IF EXISTS idx_video_calls_session;
DROP INDEX IF EXISTS idx_video_calls_participants;
DROP INDEX IF EXISTS idx_video_calls_callee;
DROP INDEX IF EXISTS idx_video_calls_caller;
DROP INDEX IF EXISTS idx_video_calls_match;

-- Drop table
DROP TABLE IF EXISTS video_calls;
