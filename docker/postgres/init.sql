\c postgres

CREATE EXTENSION IF NOT EXISTS "pg_cron";

\c social_forge

CREATE EXTENSION IF NOT EXISTS "pg_trgm";   -- trigram indexes for host/target search
CREATE EXTENSION IF NOT EXISTS "btree_gin"; -- gin index support
CREATE EXTENSION IF NOT EXISTS "citext";    -- case-insensitive text

CREATE EXTENSION IF NOT EXISTS "btree_gist";
CREATE EXTENSION IF NOT EXISTS "vector";

SET timezone = 'UTC';
