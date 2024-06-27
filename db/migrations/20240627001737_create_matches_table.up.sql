BEGIN;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'match_status') THEN
    CREATE TYPE match_status AS ENUM ('pending', 'approved', 'rejected');
  END IF;
END $$;

CREATE TABLE IF NOT EXISTS match_requests (
  id VARCHAR(26) PRIMARY KEY NOT NULL,
  match_cat_id VARCHAR(26) NOT NULL,
  user_cat_id VARCHAR(26) NOT NULL,
  message VARCHAR(200) NOT NULL,
  status match_status NOT NULL DEFAULT 'pending',
  is_deleted BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  FOREIGN KEY (match_cat_id) REFERENCES cats(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  FOREIGN KEY (user_cat_id) REFERENCES cats(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

COMMIT;