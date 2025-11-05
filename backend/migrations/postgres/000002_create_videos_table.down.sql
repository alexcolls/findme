-- Remove foreign key from users table
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_video_id;

-- Drop trigger
DROP TRIGGER IF EXISTS update_videos_updated_at ON videos;

-- Drop indexes
DROP INDEX IF EXISTS idx_videos_verification_score;
DROP INDEX IF EXISTS idx_videos_metadata;
DROP INDEX IF EXISTS idx_videos_created_at;
DROP INDEX IF EXISTS idx_videos_status;
DROP INDEX IF EXISTS idx_videos_user_id;
DROP INDEX IF EXISTS idx_videos_one_per_user;

-- Drop table
DROP TABLE IF EXISTS videos;
