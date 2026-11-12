-- Assign a monk to every pending reading before rolling back.
ALTER TABLE reading MODIFY COLUMN monk_id BIGINT NOT NULL;
