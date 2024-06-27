BEGIN;

DROP INDEX IF EXISTS idx_users_id;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_cats_id;
DROP INDEX IF EXISTS idx_cats_owner_id;
DROP INDEX IF EXISTS idx_match_requests_id;
DROP INDEX IF EXISTS idx_match_requests_match_cat_id;
DROP INDEX IF EXISTS idx_match_requests_user_cat_id;

COMMIT;