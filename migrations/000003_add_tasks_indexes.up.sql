CREATE INDEX IF NOT EXISTS tasks_title_idx ON tasks USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS tasks_description_idx ON tasks USING GIN (to_tsvector('simple', description));
CREATE INDEX IF NOT EXISTS tasks_status_idx ON tasks USING GIN (to_tsvector('simple', status));
