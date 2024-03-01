CREATE TABLE IF NOT EXISTS tasks (
id bigserial PRIMARY KEY,
title text NOT NULL,
description text NOT NULL,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
Status text NOT NULL,
expired_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
expired BOOLEAN NOT NULL,
version integer NOT NULL DEFAULT 1
);