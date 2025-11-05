-- Drop triggers and functions
DROP TRIGGER IF EXISTS check_mutual_match ON matches;
DROP FUNCTION IF EXISTS update_mutual_match();
DROP TRIGGER IF EXISTS update_matches_updated_at ON matches;

-- Drop indexes
DROP INDEX IF EXISTS idx_matches_created_at;
DROP INDEX IF EXISTS idx_matches_score;
DROP INDEX IF EXISTS idx_matches_expires_at;
DROP INDEX IF EXISTS idx_matches_mutual;
DROP INDEX IF EXISTS idx_matches_status;
DROP INDEX IF EXISTS idx_matches_matched_user_week;
DROP INDEX IF EXISTS idx_matches_user_week;

-- Drop table
DROP TABLE IF EXISTS matches;
