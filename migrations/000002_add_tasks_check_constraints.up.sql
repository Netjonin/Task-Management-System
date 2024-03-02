ALTER TABLE tasks ADD CONSTRAINT tasks_date_check CHECK (created_at <= expired_at);
